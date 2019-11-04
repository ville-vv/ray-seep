// @File     : connServer
// @Author   : Ville
// @Time     : 19-9-26 下午5:40
// node
package node

import (
	"fmt"
	"ray-seep/ray-seep/common/conn"
	"ray-seep/ray-seep/common/pkg"
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/mng"
	"runtime/debug"
	"sync"
	"time"
	"vilgo/vlog"
)

type IControlHandler interface {
	OnConnect(id int64, tr mng.MsgTransfer) error
	OnDisConnect(id int64)
	OnMessage(id int64, p *pkg.Package) (pkg.Package, error)
}

// 处理用户的通信，接收和发送用户的操作信息
type ControlServer struct {
	addr    string
	timeout time.Duration
	ish     IControlHandler
	pushCh  chan pkg.Package
}

func NewControlServer(ctlCnf *conf.ControlSrv, handler IControlHandler) *ControlServer {
	timeout := time.Millisecond * time.Duration(ctlCnf.Timeout)
	addr := fmt.Sprintf("%s:%d", ctlCnf.Host, ctlCnf.Port)
	if timeout == 0 {
		timeout = 5000
	}
	return &ControlServer{
		timeout: timeout,
		addr:    addr,
		ish:     handler,
		pushCh:  make(chan pkg.Package, 1000),
	}
}

func (sel *ControlServer) Start() error {
	lis, err := conn.Listen(sel.addr)
	if err != nil {
		panic(err)
	}
	vlog.INFO("ConnServer start [%s]", sel.addr)
	for c := range lis.Conn {
		go sel.dealConn(c)
	}
	return nil
}

// dealConn 处理连接
func (sel *ControlServer) dealConn(c conn.Conn) {
	vlog.DEBUG("customer [%d] connecting ", c.Id())

	defer func() {
		if r := recover(); r != nil {
			vlog.LogE("customer listener failed with error %v: %s", r, debug.Stack())
		}
	}()
	defer c.Close()
	// 刚刚建立连接需要设置超时时间
	_ = c.SetReadDeadline(time.Now().Add(sel.timeout))
	msgMng := mng.NewMsgTransfer(c)

	// 通知有用户连接
	if err := sel.ish.OnConnect(c.Id(), msgMng); err != nil {
		vlog.ERROR("customer [%d] connect fail %s", c.Id(), err.Error())
		return
	}
	vlog.DEBUG("customer [%d] connect success", c.Id())
	// 通知有连接断开
	defer sel.ish.OnDisConnect(c.Id())
	_ = c.SetReadDeadline(time.Time{})

	wg := sync.WaitGroup{}
	recvMsg := make(chan pkg.Package)
	cancel := make(chan pkg.Package)
	wg.Add(1)
	// 开启一个协程 接收消息
	msgMng.RecvMsgWithChan(&wg, recvMsg, cancel)
	wg.Wait()
	for {
		select {
		case req := <-recvMsg:
			rsp, err := sel.ish.OnMessage(c.Id(), &req)
			if err != nil {
				// 执行消息出现错误
				rsp = pkg.Package{Cmd: pkg.CmdError, Body: []byte(err.Error())}
			}
			_ = msgMng.SendMsg(&rsp)
		case <-cancel:
			return
		}
	}
}
