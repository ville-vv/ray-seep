package proxy

import (
	"fmt"
	"io"
	"net"
	"ray-seep/ray-seep/common/conn"
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/proto"
	"sync"
	"vilgo/vlog"
)

type ClientProxy struct {
	cfg    *conf.ProxyCli
	host   string
	port   int64
	cid    int64
	token  string
	name   string
	stopCh chan int
}

func NewClientProxy(stopCh chan int, host string, port int64) *ClientProxy {
	return &ClientProxy{
		host:   host,
		port:   port,
		stopCh: stopCh,
	}
}

func (sel *ClientProxy) RunProxy(id int64, token string, name string) error {
	sel.cid = id
	sel.token = token
	sel.name = name
	go func() {
		defer func() {
			if err := recover(); err != nil {
				vlog.ERROR("%v", err)
			}
		}()
		sel.dial()
	}()
	return nil
}
func (sel *ClientProxy) dial() {
	cn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", sel.host, sel.port))

	msgMng := proto.NewMsgTransfer(conn.TurnConn(cn))
	if err != nil {
		vlog.ERROR("connect to proxy server error %s", err.Error())
		return
	}
	defer cn.Close()
	runProxyReq := &proto.RunProxyReq{
		Cid:   sel.cid,
		Token: sel.token,
		Name:  sel.name,
	}
	err = msgMng.SendMsg(proto.NewPackage(proto.CmdRunProxyReq, runProxyReq))
	if err != nil {
		vlog.ERROR("send register proxy error %s", err.Error())
		return
	}
	//tm := time.NewTicker(time.Second * 15)
	for {
		//buf := make([]byte, 1024*4)
		//n, err := cn.Read(buf)
		//if err != nil {
		//	vlog.DEBUG("proxy 错误：%s", err.Error())
		//	return
		//}
		//buf = buf[:n]
		//vlog.DEBUG("proxy 收到消息：%s", string(buf))
		//select {
		//case <-tm.C:
		//	vlog.WARN("代理连接超时自动退出")
		//	return
		//case <-sel.stopCh:
		//	return
		//}
		cn2, err := net.Dial("tcp", ":3000") //23455
		if err != nil {
			vlog.ERROR("代理服务连接错误%v", err)
			return
		}
		defer cn2.Close()
		n, n2, err := sel.join(cn2, &registerConn{cn})
		if err != nil {
			vlog.WARN("代理访问%v", err)
			return
		}
		vlog.WARN("请求数据：%d, 返回数据：%d", n, n2)
	}
}
func (sel *ClientProxy) join(dst net.Conn, src net.Conn) (int64, int64, error) {
	var wait sync.WaitGroup
	var err error
	pipe := func(dst net.Conn, src net.Conn, bytesCopied *int64) {
		defer wait.Done()
		*bytesCopied, err = io.Copy(dst, src)
	}
	wait.Add(2)
	var fromBytes, toBytes int64
	go pipe(src, dst, &fromBytes)
	go pipe(dst, src, &toBytes)
	wait.Wait()
	return fromBytes, toBytes, err
}

type registerConn struct {
	net.Conn
}

func (sel *registerConn) Read(buf []byte) (int, error) {
	n, err := sel.Conn.Read(buf)
	if err != nil {
		vlog.ERROR("读取消息错误了%v, error : %v", sel.Conn.RemoteAddr(), err)
	}
	return n, err
}
