// @File     : manager
// @Author   : Ville
// @Time     : 19-9-24 上午9:41
// manager
package mng

import (
	"encoding/binary"
	"errors"
	"io"
	"ray-seep/ray-seep/common/pkg"
	"ray-seep/ray-seep/conn"
	"sync"
	"time"
	"vilgo/vlog"
)

const (
	MaxLinkNumber     = 100         // 客户端的最大连接数
	maxBytesCachePool = 1024 * 1024 // 接收消息的最大缓存 1M
)

// 连接管理
type ConnManager struct {
	cuMap      map[int64]conn.Conn
	cntLinkNum uint32
	maxLinkNum uint32
	sync.RWMutex
}

// client 连接管理
func NewConnManager() (cm *ConnManager) {
	cm = new(ConnManager)
	cm.cuMap = make(map[int64]conn.Conn)
	cm.cntLinkNum = 0
	cm.maxLinkNum = MaxLinkNumber
	return
}

func (cm *ConnManager) Put(key int64, cu conn.Conn) error {
	if cm.cntLinkNum >= cm.maxLinkNum {
		return errors.New("connect number is full")
	}
	cm.Lock()
	defer cm.Unlock()
	cm.cuMap[key] = cu
	cm.cntLinkNum++
	return nil
}

func (cm *ConnManager) Get(key int64) (conn.Conn, bool) {
	cm.RLock()
	defer cm.RUnlock()
	cu, ok := cm.cuMap[key]
	return cu, ok
}

func (cm *ConnManager) Delete(key int64) {
	cm.Lock()
	defer cm.Unlock()
	if cu, ok := cm.cuMap[key]; ok {
		cu.Close()
		delete(cm.cuMap, key)
		cm.cntLinkNum--
		vlog.LogD("当前连接数：%v", cm.cntLinkNum)
	}
}

type Receiver interface {
	RecvMsg(p *pkg.Package) (err error)
	RecvMsgWithChan(wait *sync.WaitGroup, mCh chan<- pkg.Package, cancel chan<- pkg.Package)
}

type receiver struct {
	r io.Reader
}

func (c *receiver) recv(buf []byte) (int, error) {
	return c.r.Read(buf)
}

func (c *receiver) Recv(buf []byte) (int, error) {
	return c.recv(buf)
}

// recvForPkg 读取一个定义好格式的消息体
func (c *receiver) recvForPkg() (buf []byte, err error) {
	var l int64

	// 先读取消息长度， 必须写入 消息的时候有写入消息长度
	err = binary.Read(c.r, binary.LittleEndian, &l)
	if err != nil {
		return
	}
	buf = make([]byte, l)

	// 再读取消息体
	n, err := c.r.Read(buf)
	if err != nil {
		return
	}

	buf = buf[:n]
	return
}

func (c *receiver) RecvMsg(p *pkg.Package) (err error) {
	buf, err := c.recvForPkg()
	if err != nil {
		return
	}
	// 解包
	return pkg.UnPack(buf, p)

}

// RecvMsgWithChan 开启一个协程 使用 chan 来接收定义好格式的消息
func (c *receiver) RecvMsgWithChan(wait *sync.WaitGroup, mCh chan<- pkg.Package, cancel chan<- pkg.Package) {
	wait.Add(1)
	go func() {
		wait.Done()
		for {
			var m pkg.Package
			if err := c.RecvMsg(&m); err != nil {
				if err == io.EOF {
					cancel <- pkg.Package{}
					return
				}
				continue
			}
			mCh <- m
		}
	}()
	return
}

type Sender interface {
	SendMsg(p *pkg.Package) (err error)
	SendMsgWithChan(wait *sync.WaitGroup, mch <-chan pkg.Package, t time.Duration)
}

type sender struct {
	w io.Writer
}

func (c *sender) send(m []byte) (err error) {
	_, err = c.w.Write(m)
	return
}

func (c *sender) sendForPkg(data []byte) (err error) {
	// 先发送消息的长度
	err = binary.Write(c.w, binary.LittleEndian, int64(len(data)))
	if err != nil {
		return
	}
	// 再发消息据体
	_, err = c.w.Write(data)
	return nil
}

// SendMsg 发送消息
func (c *sender) SendMsg(p *pkg.Package) (err error) {
	data, err := pkg.Pack(p)
	if err != nil {
		vlog.ERROR("Pack message error %v", err)
		return
	}
	return c.sendForPkg(data)
}

// SendMsgWithChan 开启一个协程 使用 chan 发送定义好格式的消息
func (c *sender) SendMsgWithChan(wait *sync.WaitGroup, mch <-chan pkg.Package, t time.Duration) {
	wait.Add(1)
	go func() {
		tk := time.NewTicker(t)
		wait.Done()
		for {
			select {
			case <-tk.C:
				return
			case mch, ok := <-mch:
				if !ok {
					return
				}
				if err := c.SendMsg(&mch); err != nil {
					continue
				}
			}
		}
	}()
}

type MsgTransfer interface {
	Receiver
	Sender
}

// 消息管理
type msgTransfer struct {
	Receiver
	Sender
	readTimeout  time.Duration
	writeTimeout time.Duration
	bytePool     sync.Pool
	isRecvStart  bool
	isSendStart  bool
}

func NewMsgTransfer(c conn.Conn) *msgTransfer {
	return &msgTransfer{
		Receiver:     &receiver{r: c},
		Sender:       &sender{w: c},
		readTimeout:  time.Second * 10,
		writeTimeout: time.Second * 10,
		bytePool: sync.Pool{
			New: func() interface{} {
				b := make([]byte, maxBytesCachePool)
				return &b
			},
		},
	}
}

type MsgManager struct {
	cuMap map[int64]MsgTransfer
	sync.RWMutex
}

func NewMsgManager() *MsgManager {
	return &MsgManager{
		cuMap: make(map[int64]MsgTransfer),
	}
}

func (cm *MsgManager) Put(key int64, cu MsgTransfer) error {
	cm.Lock()
	defer cm.Unlock()
	cm.cuMap[key] = cu
	return nil
}

func (cm *MsgManager) Get(key int64) (MsgTransfer, bool) {
	cm.RLock()
	defer cm.RUnlock()
	cu, ok := cm.cuMap[key]
	return cu, ok
}

func (cm *MsgManager) Delete(key int64) {
	cm.Lock()
	defer cm.Unlock()
	if _, ok := cm.cuMap[key]; ok {
		delete(cm.cuMap, key)
	}
}
