package proxy

import (
	"fmt"
	"github.com/vilsongwei/vilgo/vlog"
	"ray-seep/ray-seep/common/conn"
	"ray-seep/ray-seep/conf"
	"runtime/debug"
	"time"
)

type IRegister interface {
	Register(domain string, id int64, cc conn.Conn) error
}

type ProxyServer struct {
	addr      string
	proxyConn chan conn.Conn
	register  IRegister //
}

func (s *ProxyServer) Stop() {
}

func (s *ProxyServer) Scheme() string {
	return "proxy server"
}

func NewProxyServer(c *conf.ProxySrv, reg IRegister) *ProxyServer {
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	return &ProxyServer{
		addr:      addr,
		proxyConn: make(chan conn.Conn),
		register:  reg,
	}
}

func (s *ProxyServer) Start() error {
	ls, err := conn.Listen(s.addr)
	if err != nil {
		vlog.ERROR("proxy listen error %v", err)
		return err
	}
	go func() {
		for c := range ls.Conn {
			go s.dealConn(c)
		}
	}()

	return nil
}

func (s *ProxyServer) dealConn(cn conn.Conn) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			return
		}
	}()
	_ = cn.SetDeadline(time.Time{})
	// 把代理连接都注册到注册器里面
	//if err := s.register.Register(regData.Name, regData.Cid, cn); err != nil {
	//	vlog.ERROR("%s proxy is registered fail %s", cn.RemoteAddr().String(), err.Error())
	//	_ = cn.Close()
	//	return
	//}
}
