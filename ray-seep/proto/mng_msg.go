// @File     : manager
// @Author   : Ville
// @Time     : 19-9-24 上午9:41
// manager
package proto

import (
	"encoding/binary"
	"github.com/vilsongwei/vilgo/vlog"
	"io"
	"ray-seep/ray-seep/common/conn"
	"sync"
	"time"
)

const (
	maxBytesCachePool = 1024 * 1024 // 接收消息的最大缓存 1M
)

// 消息发送器
type Receiver interface {
	RecvMsg(p *Package) (err error)
	AsyncRecvMsg(wait *sync.WaitGroup, mCh chan<- Package, cancel chan<- interface{})
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
	var l int32

	// 先读取消息长度， 必须写入 消息的时候有写入消息长度
	err = binary.Read(c.r, binary.LittleEndian, &l)
	if err != nil {
		return
	}
	if l < 0 || l > 1024*1024 {
		l = 1024 * 1024
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

func (c *receiver) RecvMsg(p *Package) (err error) {
	buf, err := c.recvForPkg()
	if err != nil {
		return
	}
	// 解包
	return UnPack(buf, p)

}

// RecvMsgWithChan 开启一个协程 使用 chan 来接收定义好格式的消息
func (c *receiver) AsyncRecvMsg(wait *sync.WaitGroup, mCh chan<- Package, cancel chan<- interface{}) {
	go func() {
		wait.Done()
		for {
			var m Package
			if err := c.RecvMsg(&m); err != nil {
				cancel <- err
				return
			}
			mCh <- m
		}
	}()
	return
}

// 消息发送器
type Sender interface {
	SendMsg(p *Package) (err error)
	AsyncSendMsg(wait *sync.WaitGroup, mch <-chan Package, t time.Duration)
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
	err = binary.Write(c.w, binary.LittleEndian, int32(len(data)))
	if err != nil {
		return
	}
	// 再发消息据体
	_, err = c.w.Write(data)
	return nil
}

// SendMsg 发送消息
func (c *sender) SendMsg(p *Package) (err error) {
	data, err := Pack(p)
	if err != nil {
		vlog.ERROR("Pack message error %v", err)
		return
	}
	return c.sendForPkg(data)
}

// SendMsgWithChan 开启一个协程 使用 chan 发送定义好格式的消息
func (c *sender) AsyncSendMsg(wait *sync.WaitGroup, mch <-chan Package, t time.Duration) {
	go func() {
		wait.Done()
		for {
			select {
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

// 消息运输器，包含一个接收器和一个发送器
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

func NewMsgTransfer(c conn.Conn) MsgTransfer {
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
