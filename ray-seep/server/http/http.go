// @File     : http
// @Author   : Ville
// @Time     : 19-9-25 下午4:33
// http
package http

import (
	"fmt"
	"ray-seep/ray-seep/common/conn"
	"ray-seep/ray-seep/common/rayhttp"
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/server/proxy"
	"vilgo/vlog"
)

type Server struct {
	addr   string
	repeat Repeater // 请求中续器

}

// NewServer http 请求服务
// repeat 用于 http 请求转发
func NewServer(c *conf.HttpSrv, reg *proxy.RegisterCenter) *Server {
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	return &Server{addr: addr, repeat: NewNetRepeater(reg)}
}

// Start 启动http服务
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

// dealConn 处理 http 请求链接
func (s *Server) dealConn(c conn.Conn) {
	defer func() {
		if err := recover(); err != nil {
			vlog.ERROR("http transfer error ", err)
			return
		}
	}()
	defer c.Close()
	vlog.DEBUG("client request： %s", c.RemoteAddr())
	// 请求连接转为http协议
	copyHttp, err := rayhttp.NewCopyHttp(c)
	if err != nil {
		vlog.ERROR("%v", err)
		copyHttp.SayBackText(400, []byte("Bad Request"))
		return
	}

	// 获取请求的地址（主要是子域名有用）
	host := copyHttp.Host()
	vlog.DEBUG("request host is [%s]", host)
	//copyHttp.SayBackText(200, []byte("收到请求，请求转发尚未完成开发......"))
	// 这里转会成 Conn
	//c = conn.TurnConn(copyHttp)
	// 根据host 获取  proxy 转发
	s.repeat.Transfer(host, copyHttp)
}
