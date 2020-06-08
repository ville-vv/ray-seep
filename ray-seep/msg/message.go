package msg

import (
	"ray-seep/ray-seep/common/conn"
	"sync"
	"time"
)

type MessageCenter struct {
	recv   chan Package     // 接收消息 chan
	send   chan Package     // 发送消息 chan
	cancel chan interface{} //
	c      conn.Conn        //
	router Router           // 消息路由器
	msgTr  Transfer         // 消息运输器， 用于发送和接收消息
	pkgMng PackerManager    // 消息包管理器 用于打包和解包消息
}

func (m *MessageCenter) Run(c conn.Conn) {
	msgTr := NewMsgTransfer(c)
	wg := sync.WaitGroup{}
	wg.Add(2)
	msgTr.AsyncRecvMsg(&wg, m.recv, m.cancel)
	msgTr.AsyncSendMsg(&wg, m.send, time.Minute)
	wg.Wait()
}

func (m *MessageCenter) Close() {
	close(m.cancel)
	close(m.recv)
	close(m.send)
}

func (m *MessageCenter) Recv() <-chan Package {
	return m.recv
}

func (m *MessageCenter) Send(send Package) {
	m.send <- send
	return
}
