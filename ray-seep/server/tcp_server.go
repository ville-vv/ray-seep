// @File     : tcp_server
// @Author   : Ville
// @Time     : 19-9-24 下午3:15
// server
package server

import (
	"ray-seep/ray-seep/conn"
	"ray-seep/ray-seep/mng"
	"ray-seep/ray-seep/server/cust"
	"runtime/debug"
	"time"
	"vilgo/vlog"
)

// 管理客户端的连接
type ControlServer struct {
	addr             string
	timeout          time.Duration
	startConnTimeout time.Duration
	csMng            cust.Handler
	msgTran          *mng.MsgManager
}

func NewControlServer() *ControlServer {
	return &ControlServer{
		startConnTimeout: time.Second * 10,
		timeout:          time.Second * 15,
		csMng:            cust.NewCustomerMng(),
		addr:             ":30080",
		msgTran:          mng.NewMsgManager(),
	}
}

func (f *ControlServer) Start() {
	lis, err := conn.Listen(f.addr)
	if err != nil {
		panic(err)
	}
	vlog.INFO("ControlServer start [%s]", f.addr)
	for c := range lis.Conn {
		f.dealConn(c)
	}
}

// dealConn 处理连接
func (f *ControlServer) dealConn(conn conn.Conn) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				vlog.LogE("customer listener failed with error %v: %s", r, debug.Stack())
			}
		}()
		// 刚刚建立连接需要设置超时时间
		_ = conn.SetReadDeadline(time.Now().Add(f.startConnTimeout))
		vlog.INFO("client [%d] connecting ", conn.Id())
		f.msgTran.Put(conn.Id(), mng.NewMsgTransfer(conn))
		if err := f.csMng.Connect(conn); err != nil {
			vlog.LogE("client [%d] connect fail %v", conn.Id(), err)
			conn.Close()
			return
		}

		conn.SetReadDeadline(time.Time{})
		f.csMng.Handler(conn)
		f.csMng.DisConnect(conn.Id())
	}()
}
