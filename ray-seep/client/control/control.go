package control

import (
	"fmt"
	"net"
	"ray-seep/ray-seep/common/conn"
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/proto"
	"sync"
	"time"
	"vilgo/vlog"
)

type Router interface {
	OnConnect(sender proto.Sender) error
	OnMessage(req *proto.Package)
	OnDisconnect(id int64)
}

type ClientControl struct {
	cfg    *conf.ControlCli
	addr   string
	hd     Router
	msgMng proto.MsgTransfer
	offCh  chan int
	onCh   chan net.Conn
	stopCh chan int
}

func NewClientControl(cfg *conf.ControlCli, hd Handler) *ClientControl {
	cli := &ClientControl{
		cfg:    cfg,
		addr:   fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		hd:     NewRouteControl(hd),
		offCh:  make(chan int),
		onCh:   make(chan net.Conn),
		stopCh: make(chan int),
	}
	return cli
}
func (sel *ClientControl) shutdown() {
	close(sel.onCh)
	close(sel.offCh)
	close(sel.stopCh)
}
func (sel *ClientControl) Stop() {
}
func (sel *ClientControl) Start() {
	go sel.offDial()
	go sel.onDial()

	c, err := net.Dial("tcp", sel.addr)
	if err != nil {
		vlog.LogE("connect server fail %v", err)
		return
	}
	sel.onCh <- c
}

func (sel *ClientControl) onDial() {
	for v := range sel.onCh {
		go sel.dealConn(v)
	}
	vlog.WARN("client control exit")
}
func (sel *ClientControl) offDial() {
	for range sel.offCh {
		go sel.reconnect()
	}
}

func (sel *ClientControl) reconnect() {
	tm := time.NewTicker(time.Second)
	endTm := time.NewTimer(time.Minute * 3)
	for {
		select {
		case <-tm.C:
			c, err := net.Dial("tcp", sel.addr)
			if err != nil {
				break
			}
			sel.onCh <- c
			return
		case <-endTm.C:
			vlog.WARN("重连超时")
			sel.shutdown()
			return
		}
	}
}

func (sel *ClientControl) dealConn(c net.Conn) {
	defer c.Close()
	sel.msgMng = proto.NewMsgTransfer(conn.TurnConn(c))
	if err := sel.hd.OnConnect(sel.msgMng); err != nil {
		vlog.ERROR("server connect error:%s", err.Error())
		return
	}
	defer sel.hd.OnDisconnect(0)
	var wg sync.WaitGroup
	recvCh := make(chan proto.Package)
	cancel := make(chan interface{})
	wg.Add(1)
	sel.msgMng.AsyncRecvMsg(&wg, recvCh, cancel)
	wg.Wait()
	isOff := false
	for {
		select {
		case ms, ok := <-recvCh:
			if !ok {
				isOff = true
			}
			sel.hd.OnMessage(&ms)
		case err := <-cancel:
			vlog.ERROR("disconnect：%v", err)
			isOff = true
		}
		if isOff {
			sel.offCh <- 1
			return
		}
	}
}

func (sel *ClientControl) PushEvent(cmd int32, dt []byte) error {
	return sel.pushEvent(&proto.Package{Cmd: cmd, Body: dt})
}

func (sel *ClientControl) pushEvent(p *proto.Package) error {
	return sel.msgMng.SendMsg(p)
}
