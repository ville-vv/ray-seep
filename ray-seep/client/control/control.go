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
	OnMessage(req *proto.Package) error
	OnDisconnect(id int64)
}

type ClientControl struct {
	cfg            *conf.ControlCli
	addr           string
	haveConn       bool
	isReconnect    bool  // 是否启动自动重连
	reConnEndTime  int64 // 重连持续时间（断开多久就不再重连了）
	reConnInternal int64 // 重连间隔时间（多久重连一次）
	hd             Router
	msgMng         proto.MsgTransfer
	offCh          chan int
	onCh           chan net.Conn
	stopCh         chan int
}

func NewClientControl(cfg *conf.ControlCli, hd Handler) *ClientControl {
	cli := &ClientControl{
		cfg:            cfg,
		addr:           fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		hd:             NewRouteControl(hd),
		offCh:          make(chan int),
		onCh:           make(chan net.Conn),
		stopCh:         make(chan int),
		isReconnect:    cfg.CanReconnect,
		reConnInternal: cfg.ReconnectInternal,
		reConnEndTime:  cfg.ReconnectEndTime,
	}

	if cli.reConnInternal <= 0 {
		cli.reConnInternal = 3
	}
	if cli.reConnEndTime <= 0 {
		cli.reConnEndTime = 60
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
		vlog.LogE("connect control server fail %v", err)
		return
	}
	sel.onCh <- c
}

// 检测有打开链接
func (sel *ClientControl) onDial() {
	for v := range sel.onCh {
		go sel.dealConn(v)
	}
	vlog.WARN("client control exit")
}

// 检测有断开链接
func (sel *ClientControl) offDial() {
	if !sel.isReconnect {
		return
	}
	for range sel.offCh {
		go sel.reconnect()
	}
}

func (sel *ClientControl) reconnect() {
	vlog.DEBUG("reConnInternal %s", sel.reConnInternal)
	vlog.DEBUG("reConnEndTime %s", sel.reConnEndTime)
	tm := time.NewTicker(time.Second * time.Duration(sel.reConnInternal))
	endTm := time.NewTicker(time.Second * time.Duration(sel.reConnEndTime))
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
	recvCh := make(chan proto.Package)
	cancel := make(chan interface{})
	defer func() {
		close(recvCh)
		close(cancel)
		sel.hd.OnDisconnect(0)
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	sel.msgMng.AsyncRecvMsg(&wg, recvCh, cancel)
	wg.Wait()
	needReConn := false
	for {
		select {
		case ms := <-recvCh:
			if err := sel.hd.OnMessage(&ms); err != nil {
				needReConn = false
				return
			}
		case err := <-cancel:
			vlog.ERROR("disconnect：%v", err)
			needReConn = true
		}
		// 只有打开了断线重连才能发送重连消息
		if needReConn && sel.isReconnect {
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
