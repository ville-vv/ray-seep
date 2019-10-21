package proxy

import (
	"encoding/json"
	"errors"
	"fmt"
	"ray-seep/ray-seep/common/conn"
	"ray-seep/ray-seep/common/pkg"
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/mng"
	"vilgo/vlog"
)

type IRegister interface {
	Register(domain string, cid int64) error
}

type Server struct {
	addr      string
	proxyConn chan conn.Conn
	register  IRegister //
}

func NewServer(c *conf.ProxySrv, reg IRegister) *Server {
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	return &Server{
		addr:      addr,
		proxyConn: make(chan conn.Conn),
		register:  reg,
	}
}

func (s *Server) Start() {
	ls, err := conn.Listen(s.addr)
	if err != nil {
		return
	}
	vlog.INFO("ProxyServer start [%s]", s.addr)
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
	defer cn.Close()
	tr := mng.NewMsgTransfer(cn)
	var regProxy pkg.Package

	if err := tr.RecvMsg(&regProxy); err != nil {
		return
	}

	if regProxy.Cmd != pkg.CmdRegisterProxyReq {
		return
	}
	regData := pkg.RegisterProxyReq{}
	if err := json.Unmarshal(regProxy.Body, &regData); err != nil {
		return
	}
	if err := s.register.Register(regData.SubDomain, regData.Cid); err != nil {
		return
	}

}
func (s *Server) SetProxy(cn conn.Conn) {
	//_= cn.SetDeadline(time.Now().Add(time.Second * 15))
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
