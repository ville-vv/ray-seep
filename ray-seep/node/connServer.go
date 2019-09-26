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
	OnConnect(id int64, sender mng.Sender)
	OnDisConnect(id int64)
	OnHandler(id int64, p pkg.Package, sender mng.Sender)
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
		ish:     NewAdopterNode(),
	}
}

func (sel *ConnServer) Start() error {
	lis, err := conn.Listen(sel.addr)
	if err != nil {
		panic(err)
	}
	vlog.INFO("ConnServer start [%s]", sel.addr)
	for c := range lis.Conn {
		sel.dealConn(c)
	}
	return nil
}

// dealConn 处理连接
func (sel *ConnServer) dealConn(c conn.Conn) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				vlog.LogE("customer listener failed with error %v: %s", r, debug.Stack())
			}
		}()

		// 刚刚建立连接需要设置超时时间
		_ = c.SetReadDeadline(time.Now().Add(sel.timeout))

		cid := c.Id()

		vlog.INFO("customer [%d] connecting ", cid)
		msgMng := mng.NewMsgTransfer(c)
		// 通知有连接进来
		sel.ish.OnConnect(cid, msgMng)
		_ = c.SetReadDeadline(time.Time{})
		// 通知有连接断开
		defer sel.ish.OnDisConnect(cid)

		wg := sync.WaitGroup{}
		recvMsg := make(chan pkg.Package)
		cancel := make(chan pkg.Package)
		msgMng.RecvMsgWithChan(&wg, recvMsg, cancel)
		wg.Wait()
		for {
			select {
			case m := <-recvMsg:
				sel.ish.OnHandler(cid, m, msgMng)
			case <-cancel:
				return
			}
		}
	}()
}
