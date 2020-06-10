package msg

import (
	"context"
	"ray-seep/ray-seep/common/conn"
	"sync"
	"time"
)

type MessageCenter struct {
	recvCh chan Package     // 接收消息 chan
	sendCh chan Package     // 发送消息 chan
	cancel chan interface{} //
	c      conn.Conn        //
	router RouterFunc       // 消息路由器
	msgTr  Transfer         // 消息运输器， 用于发送和接收消息
	pkgMng PackerManager    // 消息包管理器 用于打包和解包消息
}

func NewMessageCenter(c conn.Conn) *MessageCenter {
	return &MessageCenter{
		c:      c,
		recvCh: make(chan Package, 100),
		sendCh: make(chan Package, 100),
		cancel: make(chan interface{}),
		msgTr:  NewMsgTransfer(c),
		pkgMng: &packerManager01{},
	}
}

// 设置路由
func (m *MessageCenter) SetRouter(router RouterFunc) {
	m.router = router
	return
}

func (m *MessageCenter) Run() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	m.AsyncRecvMsg(&wg)
	m.AsyncSendMsg(&wg, time.Minute)
	wg.Wait()
}

func (m *MessageCenter) Cancel() <-chan interface{} {
	return m.cancel
}

func (m *MessageCenter) Close() {
	close(m.cancel)
	close(m.recvCh)
	close(m.sendCh)
}

// 使用 chan 接收消息
//func (m *MessageCenter) RecvCh() <-chan Package {
//	return m.recvCh
//}

// 使用 chan 推送消息
func (m *MessageCenter) SendCh() chan<- Package {
	return m.sendCh
}

// 启动一个异步接收消息的协程，消息会发送到 router 处理
// 如果要使用它，需要调用 SetRouter 设置处理数据的 router
func (m *MessageCenter) AsyncRecvMsg(wait *sync.WaitGroup) {
	go func() {
		wait.Done()
		ctx, cel := context.WithCancel(context.Background())
		for {
			pg := new(Package)
			if err := m.recvMsg(pg); err != nil {
				m.cancel <- err
				cel()
				return
			}
			//m.recvCh <- *pg
			if err := m.router(&Request{ctx: ctx, Body: pg}, m); err != nil {
				return
			}
		}
	}()
	return
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
func (m *MessageCenter) AsyncSendMsg(wait *sync.WaitGroup, t time.Duration) {
	go func() {
		wait.Done()
		for {
			select {
			case mch, ok := <-m.sendCh:
				if !ok {
					return
				}
				if err := m.sendMsg(&mch); err != nil {
					return
				}
			}
		}
	}()
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
