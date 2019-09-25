// @File     : http
// @Author   : Ville
// @Time     : 19-9-25 下午4:33
// http
package http

import (
	"ray-seep/ray-seep/conn"
	"vilgo/vlog"
)

type Server struct {
	addr string
}

func NewServer() *Server {
	return &Server{addr: ":40090"}
}

func (s *Server) Start() {
	lis, err := conn.Listen(s.addr)
	if err != nil {
		vlog.ERROR("%v", err)
		return
	}
	vlog.INFO("HttpServer start [%s]", s.addr)
	for c := range lis.Conn {
		go s.dealConn(c)
	}
}

func (s *Server) dealConn(c conn.Conn) {
	vlog.INFO("收到请求：%s", c.RemoteAddr())
	headConn, err := NewCopyHttp(c)
	if err != nil {
		headConn.SayBack(400, []byte("error "))
		return
	}
	headConn.SayBack(200, []byte("success"))
}
