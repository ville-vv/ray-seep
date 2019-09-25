// @File     : manager
// @Author   : Ville
// @Time     : 19-9-24 上午9:41
// manager
package mng

import (
	"encoding/binary"
	"errors"
	"io"
	"ray-seep/ray-seep/conn"
	"ray-seep/ray-seep/msg"
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

func (c *receiver) RecvMsg(pkg *msg.Message) (err error) {
	buf, err := c.recvForPkg()
	if err != nil {
		return
	}
	// 解包
	return msg.UnPack(buf, pkg)
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
func (c *sender) SendMsg(pkg *msg.Message) (err error) {
	data, err := msg.Pack(pkg)
	if err != nil {
		vlog.ERROR("Pack message error %v", err)
		return
	}
	return c.sendForPkg(data)
}

// 消息管理
type MsgTransfer struct {
	cn conn.Conn
	receiver
	sender
	cid          int64
	sequence     uint32
	stopCh       chan int32
	readTimeout  time.Duration
	writeTimeout time.Duration
	bytePool     sync.Pool
	isRecvStart  bool
	isSendStart  bool
}

func NewMsgTransfer(c conn.Conn) *MsgTransfer {
	return &MsgTransfer{
		cn:           c,
		receiver:     receiver{r: c},
		sender:       sender{w: c},
		stopCh:       make(chan int32),
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

func (c *MsgTransfer) SetReader(reader io.Reader) {
	c.r = reader
}

func (c *MsgTransfer) SetWriter(writer io.Writer) {
	c.w = writer
}

func (c *MsgTransfer) Id() int64 {
	return c.cn.Id()
}

// RecvMsgWithChan 开启一个协程 使用 chan 来接收定义好格式的消息
func (c *MsgTransfer) RecvMsgWithChan(wait *sync.WaitGroup, mCh chan<- msg.Message, cancel chan<- msg.Message) {
	wait.Add(1)
	go func() {
		wait.Done()
		for {
			select {
			case <-c.stopCh:
				return
			default:
				var m msg.Message
				if err := c.RecvMsg(&m); err != nil {
					if err == io.EOF {
						cancel <- msg.Message{Cmd: "", Body: c.cid}
						return
					}
					continue
				}
				select {
				case <-c.stopCh:
					return
				default:
					mCh <- m
				}
			}
		}
	}()
	return
}

// SendMsgWithChan 开启一个协程 使用 chan 发送定义好格式的消息
func (c *MsgTransfer) SendMsgWithChan(wait *sync.WaitGroup, mch <-chan msg.Message) {
	wait.Add(1)
	go func() {
		wait.Done()
		for {
			select {
			case <-c.stopCh:
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

type MsgManager struct {
	cuMap map[int64]*MsgTransfer
	sync.RWMutex
}

func NewMsgManager() *MsgManager {
	return &MsgManager{
		cuMap: make(map[int64]*MsgTransfer),
	}
}

func (cm *MsgManager) Put(key int64, cu *MsgTransfer) error {
	cm.Lock()
	defer cm.Unlock()
	cm.cuMap[key] = cu
	return nil
}

func (cm *MsgManager) Get(key int64) (*MsgTransfer, bool) {
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
