package proxy

import (
	"fmt"
	"net"
	"ray-seep/ray-seep/common/conn"
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/proto"
	"vilgo/vlog"
)

type ClientProxy struct {
	cfg       *conf.ProxyCli
	host      string
	port      int64
	cid       int64
	token     string
	subDomain string
	stopCh    chan int
}

func NewClientProxy(stopCh chan int, sdm string, host string, port int64) *ClientProxy {
	return &ClientProxy{
		subDomain: sdm,
		host:      host,
		port:      port,
		stopCh:    stopCh,
	}
}

func (sel *ClientProxy) RunProxy(id int64, token string) error {
	sel.cid = id
	sel.token = token
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
		Cid:       sel.cid,
		Token:     sel.token,
		SubDomain: sel.subDomain,
	}
	err = msgMng.SendMsg(proto.NewPackage(proto.CmdRunProxyReq, runProxyReq))
	if err != nil {
		vlog.ERROR("send register proxy error %s", err.Error())
		return
	}
	//tm := time.NewTicker(time.Second * 15)
	for {
		buf := make([]byte, 1024*4)
		n, err := cn.Read(buf)
		if err != nil {
			vlog.DEBUG("proxy 错误：%s", err.Error())
			return
		}
		buf = buf[:n]
		vlog.DEBUG("proxy 收到消息：%s", string(buf))
		//select {
		//case <-tm.C:
		//	vlog.WARN("代理连接超时自动退出")
		//	return
		//case <-sel.stopCh:
		//	return
		//}
	}
}
