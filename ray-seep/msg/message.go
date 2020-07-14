package msg

import (
	"context"
	"errors"
	"github.com/vilsongwei/vilgo/vlog"
	"io"
	"ray-seep/ray-seep/common/conn"
	"sync"
	"time"
)

type MessageCenter struct {
	recvCh      chan Package  // 接收消息 chan
	sendCh      chan Package  // 发送消息 chan
	stop        chan int      //
	c           conn.Conn     //
	router      RouterFunc    // 消息路由器
	msgTr       Transfer      // 消息运输器， 用于发送和接收消息
	pkgMng      PackerManager // 消息包管理器 用于打包和解包消息
	readTimeOut int64
	isTimeout   bool
	Err         error
}

func NewMessageCenter(c conn.Conn) *MessageCenter {
	return &MessageCenter{
		readTimeOut: 60,
		c:           c,
		recvCh:      make(chan Package, 10000),
		sendCh:      make(chan Package, 10000),
		stop:        make(chan int),
		msgTr:       NewMsgTransfer(c),
		pkgMng:      &packerManager01{},
	}
}

func (m *MessageCenter) IsTimeout() bool {
	return m.isTimeout
}

// 设置路由
func (m *MessageCenter) SetRouter(router RouterFunc) {
	m.router = router
	return
}

func (m *MessageCenter) Run(router RouterFunc) {
	m.router = router
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		wg.Done()
		m.SendMsg(time.Minute)
	}()
	wg.Wait()
	m.RecvMsg()
	m.close()
}

func (m *MessageCenter) close() {
	close(m.sendCh)
	close(m.recvCh)
	close(m.stop)
}

// 使用 chan 推送消息
func (m *MessageCenter) SendCh(p *Package) error {
	select {
	case m.sendCh <- *p:
	default:
		return errors.New("send message chan is error")
	}
	return nil
}

// 启动一个异步接收消息的协程，消息会发送到 router 处理
// 如果要使用它，需要调用 SetRouter 设置处理数据的 router
func (m *MessageCenter) RecvMsg() {
	ctx, cel := context.WithCancel(context.Background())
	for {
		select {
		case <-m.stop:
			cel()
			return
		default:
		}
		pg := new(Package)
		_ = m.c.SetReadDeadline(time.Now().Add(time.Duration(m.readTimeOut) * time.Second))
		if err := m.recvMsg(pg); err != nil {
			cel()
			if err != io.EOF {
				vlog.ERROR("break off connect: %s", err.Error())
				m.Err = err
			} else {
				vlog.WARN("break off connect")
			}
			return
		}
		if err := m.router(&Request{Ctx: ctx, Body: pg}); err != nil {
			// 这里 不等一会，错误消息就发送不出去
			time.Sleep(time.Millisecond)
			return
		}
	}
}

// 接收消息，会阻塞流程直到收到消息才向下走
func (m *MessageCenter) recvMsg(pg *Package) error {
	data, err := m.msgTr.RecvMsg()
	if err != nil {
		return err
	}
	if err = m.pkgMng.UnPack(data, pg); err != nil {
		return err
	}
	return nil
}

// 提供给外部接收消息，会阻塞流程直到收到消息才向下走,
// 不会调用 router
func (m *MessageCenter) Recv(pg *Package) error {
	return m.recvMsg(pg)
}

// AsyncSendMsg 开启一个协程 使用 chan 发送定义好格式的消息
func (m *MessageCenter) SendMsg(t time.Duration) {
	for {
		select {
		case mch, ok := <-m.sendCh:
			if !ok {
				return
			}
			if err := m.sendMsg(&mch); err != nil {
				continue
			}
		}
	}
}

// 发送想消息
func (m *MessageCenter) sendMsg(pg *Package) error {
	data, err := m.pkgMng.Pack(pg)
	if err != nil {
		return err
	}
	return m.msgTr.SendMsg(data)
}

// 发送消息-提供给外部使用
func (m *MessageCenter) Send(pg *Package) error {
	return m.sendMsg(pg)
}
