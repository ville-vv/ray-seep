// @File     : http
// @Author   : Ville
// @Time     : 19-9-25 下午4:33
// http
package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"ray-seep/ray-seep/common/rayhttp"
	"ray-seep/ray-seep/conf"
	"time"
	"vilgo/vlog"
)

type Server struct {
	addr   string
	repeat Repeater // 请求中续器

}

func (s *Server) Stop() {}

func (s *Server) Scheme() string {
	return "http server"
}

// NewServer http 请求服务
// repeat 用于 http 请求转发
func NewServer(c *conf.HttpSrv, pxyGainer ProxyGainer) *Server {
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	return &Server{addr: addr, repeat: NewNetRepeater(pxyGainer)}
}

// Start 启动http服务
func (s *Server) Start() error {
	lin, err := net.Listen("tcp", s.addr)
	if err != nil {
		vlog.ERROR("%v", err)
		return nil
	}
	vlog.INFO("HttpServer start [%s]", s.addr)
	for {
		c, err := lin.Accept()
		if err != nil {
			vlog.ERROR("http accept error %s", err.Error())
		}
		go s.dealConn(c)
	}
}

// dealConn 处理 http 请求链接
func (s *Server) dealConn(c net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			vlog.ERROR("http transfer error %v", err)
			return
		}
	}()
	defer c.Close()
	_ = c.(*net.TCPConn).SetKeepAlive(true)
	vlog.INFO("http request from： %s", c.RemoteAddr())
	// 请求连接转为http协议
	copyHttp, err := rayhttp.ToHttp(c)
	if err != nil {
		vlog.ERROR("tcp connect  to http request error %v", err.Error())
		SayBackText(c, 400, []byte("Bad Request"))
		return
	}
	// 获取请求的地址（主要是子域名有用）
	host := copyHttp.Host()
	//vlog.DEBUG("request proxy host is [%s]", host)
	// 根据host 获取  proxy 转发
	s.repeat.Transfer(host, copyHttp)
}

func SayBackText(c net.Conn, status int, body []byte) {
	resp := http.Response{
		StatusCode:    status,
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        http.Header{},
		ContentLength: int64(len(body)),
		Body:          ioutil.NopCloser(bytes.NewBuffer(body)),
	}
	resp.Header.Add("Content-Type", "text/html;charset=utf-8")
	resp.Header.Add("date", time.Now().Format(time.RFC1123))
	resp.Write(c)
}
