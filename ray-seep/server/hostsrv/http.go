package hostsrv

import (
	"bytes"
	"github.com/vilsongwei/vilgo/vlog"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"ray-seep/ray-seep/common/rayhttp"
	"ray-seep/ray-seep/common/repeat"
	"runtime/debug"
	"time"
)

// Repeater 是一个中继器，用于转发 conn 的数据
type Repeater interface {
	// 转发
	Transfer(host string, c net.Conn) (int64, int64, error)
}

type httpRunner struct {
	clientId int64
	lis      net.Listener
	isStop   bool
	addr     string
	repeat   Repeater // 请求中续器
}

func (s *httpRunner) Id() int64 {
	return s.clientId
}

func (s *httpRunner) Stop() {
	s.isStop = true
	if s.lis != nil {
		_ = s.lis.Close()
	}
}

func (s *httpRunner) Scheme() string {
	return "http server"
}

// NewServer http 请求服务
// repeat 用于 http 请求转发
func newHttpRunner(id int64, addr string, pxyGainer repeat.NetConnGainer) *httpRunner {
	return &httpRunner{clientId: id, addr: addr, repeat: repeat.NewNetRepeater(pxyGainer)}
}

// Start 启动http服务
func (s *httpRunner) Start() error {
	lin, err := net.Listen("tcp", s.addr)
	if err != nil {
		vlog.ERROR("http listen error %v", err)
		return err
	}
	s.lis = lin
	go func() {
		for !s.isStop {
			c, err := lin.Accept()
			// 上报连接数
			if err != nil {
				operr, ok := err.(*net.OpError)
				if !(ok && operr.Err.Error() == "use of closed network connection") {
					vlog.ERROR("http accept error %s", err.Error())
				}
				return
			}
			go s.dealConn(c)
		}
	}()
	return nil
}

// dealConn 处理 http 请求链接
func (s *httpRunner) dealConn(c net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			return
		}
	}()
	defer c.Close()
	_ = c.(*net.TCPConn).SetKeepAlive(true)
	vlog.DEBUG("http request from： %s", c.RemoteAddr())
	// 请求连接转为http协议
	copyHttp, err := rayhttp.ToHttp(c)
	if err != nil {
		if err == io.EOF {
			return
		}
		vlog.ERROR("tcp connect  to http request error %v", err.Error())
		SayBackText(c, 400, []byte("Bad Request"))
		return
	}
	// 获取请求的地址（主要是子域名有用）
	host := copyHttp.Host()
	//vlog.DEBUG("request proxy host is [%s]", host)
	// 根据host 获取  proxy 转发
	reqLength, respLength, err := s.repeat.Transfer(host, copyHttp)
	if err != nil {
		vlog.ERROR("%s", err.Error())
		return
	}
	vlog.INFO("request size：[%d]. response size：[%d]", reqLength, respLength)
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
