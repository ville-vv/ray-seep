// @File     : proxy
// @Author   : Ville
// @Time     : 19-9-27 下午4:59
// proxy
package proxy

import (
	"errors"
	"ray-seep/ray-seep/common/pkg"
	"ray-seep/ray-seep/conn"
	"ray-seep/ray-seep/mng"
	"time"
	"vilgo/vlog"
)

type Server struct {
	addr      string
	proxyConn chan conn.Conn
}

func NewServer() *Server {
	return &Server{
		addr: ":39990",
	}
}

func (s *Server) Start() {
	ls, err := conn.Listen(s.addr)
	if err != nil {
		return
	}
	for c := range ls.Conn {
		go s.dealConn(c)
	}
}

func (s *Server) dealConn(cn conn.Conn) {
	defer func() {
		if err := recover(); err != nil {
			vlog.DEBUG("")
			return
		}
	}()
	tr := mng.NewMsgTransfer(cn)
	var regProxy pkg.Package

	if err := tr.RecvMsg(&regProxy); err != nil {
		cn.Close()
		return
	}

	if regProxy.Cmd != pkg.CmdRegisterProxyReq {
		cn.Close()
		return
	}
}
func (s *Server) SetProxy(cn conn.Conn) {
	cn.SetDeadline(time.Now().Add(time.Second * 15))
	select {
	case s.proxyConn <- cn:
	default:
		vlog.WARN("Proxies buffer is full, discarding.")
	}
}

func (s *Server) GetProxy() (cn conn.Conn, err error) {
	var ok bool
	select {
	case cn, ok = <-s.proxyConn:
		if !ok {
			err = errors.New("No proxy connections available, control is closing")
			return
		}
	}
	return
}
