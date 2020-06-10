package msg

import (
	"encoding/binary"
	"io"
	"ray-seep/ray-seep/common/conn"
	"sync"
	"time"
)

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
func (c *receiver) RecvMsg() (buf []byte, err error) {
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

type sender struct {
	w io.Writer
}

func (c *sender) send(m []byte) (err error) {
	_, err = c.w.Write(m)
	return
}

func (c *sender) SendMsg(data []byte) (err error) {
	// 先发送消息的长度
	err = binary.Write(c.w, binary.LittleEndian, int32(len(data)))
	if err != nil {
		return
	}
	// 再发消息据体
	_, err = c.w.Write(data)
	return nil
}

// 消息运输器，包含一个接收器和一个发送器
type Transfer interface {
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

func NewMsgTransfer(c conn.Conn) Transfer {
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
