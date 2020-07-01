package control

import (
	"fmt"
	"github.com/vilsongwei/vilgo/vlog"
	"net"
	"ray-seep/ray-seep/common/conn"
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/msg"
	"time"
)

type Router interface {
	OnConnect(sender msg.ResponseSender) error
	OnMessage(req *msg.Request) error
	OnDisconnect(id int64)
}

type ClientManager struct {
	cfg  *conf.ControlCli
	addr string

	connCh    chan net.Conn
	errCh     chan error
	runCh     chan int
	stopCh    chan int
	isClose   bool
	waitClose chan bool

	route          Router
	isReconnect    bool  // 是否启动自动重连
	reConnEndTime  int64 // 重连持续时间（断开多久就不再重连了）
	reConnInternal int64 // 重连间隔时间（多久重连一次）
}

func (sel *ClientManager) Start() {
	go sel.process()
	sel.runCh <- 1
}

func (sel *ClientManager) WaitClose() <-chan bool {
	return sel.waitClose
}

func (sel *ClientManager) close(ok bool) {
	close(sel.runCh)
	close(sel.errCh)
	close(sel.connCh)
	close(sel.waitClose)
	sel.isClose = true
	if ok {
		close(sel.stopCh)
	}
}

func (sel *ClientManager) Stop() {
	if sel.isClose {
		return
	}
	sel.stopCh <- 1
}

func NewClientManager(cfg *conf.ControlCli, hd Handler) *ClientManager {
	cli := &ClientManager{
		cfg:   cfg,
		addr:  fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		route: NewRouteControl(hd),

		connCh:    make(chan net.Conn),
		errCh:     make(chan error),
		runCh:     make(chan int),
		stopCh:    make(chan int),
		waitClose: make(chan bool),

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

func (sel *ClientManager) process() {
	for {
		select {
		case cn, ok := <-sel.connCh:
			if ok {
				go sel.dealConn(cn)
			}
		case err, ok := <-sel.errCh:
			if ok {
				go sel.dealErr(err)
			}
		case _, ok := <-sel.runCh:
			if ok {
				go sel.connect()
			}
		case _, ok := <-sel.stopCh:
			sel.close(ok)
			return
		}
	}
}

func (sel *ClientManager) connect() {
	c, err := net.Dial("tcp", sel.addr)
	if err != nil {
		vlog.LogE("connect node server fail %v", err)
		sel.errCh <- err
		return
	}
	sel.connCh <- c
}

func (sel *ClientManager) reconnect() {
	tm := time.NewTicker(time.Second * time.Duration(sel.reConnInternal))
	endTm := time.NewTicker(time.Second * time.Duration(sel.reConnEndTime))
	for {
		select {
		case <-tm.C:
			sel.runCh <- 1
			return
		case <-endTm.C:
			vlog.WARN("重连超时")
			return
		}
	}
}

func (sel *ClientManager) dealConn(c net.Conn) {
	defer c.Close()
	var err error
	msgMng := msg.NewMessageCenter(conn.TurnConn(c))
	if err = sel.route.OnConnect(msgMng); err != nil {
		vlog.ERROR("server connect error:%s", err.Error())
		return
	}
	defer func() {
		sel.errCh <- err
		sel.route.OnDisconnect(0)
	}()
	msgMng.Run(sel.route.OnMessage)
	err = msgMng.Err
}

func (sel *ClientManager) dealErr(err error) {
	if err == nil {
		sel.Stop()
		return
	}
	vlog.ERROR("收到错误：%v", err)
	sel.reconnect()
}
