// @File     : tcp_server
// @Author   : Ville
// @Time     : 19-9-24 下午3:15
// server
package server

import (
	"ray-seep/ray-seep/conn"
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
	csMng            cust.Manager
}

func NewControlServer() *ControlServer {
	return &ControlServer{
		startConnTimeout: time.Second * 15,
		timeout:          time.Second * 15,
		csMng:            cust.NewCustomerMng(),
		addr:             ":30080",
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
				vlog.LogE("tunnelListener failed with error %v: %s", r, debug.Stack())
			}
		}()
		_ = conn.SetReadDeadline(time.Now().Add(f.startConnTimeout))
		vlog.INFO("客户端正在链接：%v", conn.RemoteAddr())
		if err := f.csMng.Connect(conn); err != nil {
			vlog.LogE("client %v connect fail %v", conn.RemoteAddr(), err)
			conn.Close()
			return
		}
		conn.SetReadDeadline(time.Time{})
		f.csMng.Handler(conn)
		f.csMng.DisConnect(conn.Id())
	}()
}
