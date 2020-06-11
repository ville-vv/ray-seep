package server_v2

import (
	"fmt"
	"github.com/vilsongwei/vilgo/vlog"
	"ray-seep/ray-seep/common/conn"
	"ray-seep/ray-seep/conf"
	"runtime/debug"
	"time"
)

type ServerHandler interface {
	OnConnect(cancel chan interface{}, cn conn.Conn) error
	OnDisConnect(id int64)
}

// 处理用户的通信，接收和发送用户的操作信息
type ControlServer struct {
	scheme  string
	addr    string
	timeout int64
	ish     ServerHandler
}

func (sel *ControlServer) Stop() {}

func (sel *ControlServer) Scheme() string {
	return sel.scheme
}

func NewControlServer(src *conf.ControlSrv, handler ServerHandler) *ControlServer {
	ctlCnf := src
	timeout := ctlCnf.ReadMsgTimeout
	addr := fmt.Sprintf("%s:%d", ctlCnf.Host, ctlCnf.Port)
	if timeout == 0 {
		timeout = 5000
	}
	return &ControlServer{
		timeout: timeout,
		addr:    addr,
		ish:     handler,
	}
}

func (sel *ControlServer) Start() error {
	lis, err := conn.Listen(sel.addr)
	if err != nil {
		vlog.ERROR("node listen error %v", err)
		return err
	}
	for c := range lis.Conn {
		go sel.dealConn(c)
	}
	return nil
}

// dealConn 处理连接
func (sel *ControlServer) dealConn(c conn.Conn) {
	defer func() {
		if r := recover(); r != nil {
			vlog.LogE("customer listener failed with error %v: %s", r, debug.Stack())
		}
	}()

	cancel := make(chan interface{})
	defer func() {
		close(cancel)
		// 通知有连接断开
		sel.ish.OnDisConnect(c.Id())
	}()
	// 刚刚建立连接需要设置超时时间
	_ = c.SetReadDeadline(time.Now().Add(time.Duration(sel.timeout) * time.Millisecond))
	// 通知有用户连接
	if err := sel.ish.OnConnect(cancel, c); err != nil {
		vlog.DEBUG("[%d] connect fail %s", c.Id(), err.Error())
		return
	}
}
