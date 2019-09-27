// @File     : connServer
// @Author   : Ville
// @Time     : 19-9-26 下午5:40
// node
package node

import (
	"ray-seep/ray-seep/common/pkg"
	"ray-seep/ray-seep/conn"
	"ray-seep/ray-seep/mng"
	"runtime/debug"
	"sync"
	"time"
	"vilgo/vlog"
)

type IConnServerHandler interface {
	OnConnect(id int64, tr mng.MsgTransfer) error
	OnDisConnect(id int64)
	OnMessage(id int64, p *pkg.Package) (pkg.Package, error)
}

type ConnServer struct {
	addr    string
	timeout time.Duration
	ish     IConnServerHandler
}

func NewConnServer() *ConnServer {
	return &ConnServer{
		timeout: time.Second * 15,
		addr:    ":30080",
		ish:     NewAdopterPod(),
	}
}

func (sel *ConnServer) Start() error {
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
func (sel *ConnServer) dealConn(c conn.Conn) {
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
			msgMng.SendMsg(&rsp)
		case <-cancel:
			return
		}
	}
}
