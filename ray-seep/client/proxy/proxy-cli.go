package proxy

import (
	"net"
	"ray-seep/ray-seep/common/conn"
	"ray-seep/ray-seep/common/repeat"
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/proto"
	"vilgo/vlog"
)

type ClientProxy struct {
	cfg    *conf.ProxyCli
	addr   string
	cid    int64
	token  string
	name   string
	stopCh chan int
	netRet *repeat.NetRepeater
}

func NewClientProxy(stopCh chan int, cfg *conf.Client) *ClientProxy {
	return &ClientProxy{
		addr:   cfg.Pxy.Address,
		stopCh: stopCh,
		netRet: repeat.NewNetRepeater(NewTunnel(cfg.Web.Address)),
	}
}

func (sel *ClientProxy) RunProxy(id int64, token string, name string, addr string) error {
	sel.cid = id
	sel.token = token
	sel.name = name
	sel.addr = addr
	go func() {
		defer func() {
			if err := recover(); err != nil {
				vlog.ERROR("%v", err)
			}
		}()
		cn, err := sel.dial()
		if err != nil {
			vlog.ERROR("connect to proxy server error %s", err.Error())
		}
		defer cn.Close()
		sel.netRet.Transfer(name, cn)
	}()
	return nil
}

func (sel *ClientProxy) dial() (net.Conn, error) {
	vlog.INFO("连接服务：%s", sel.addr)
	cn, err := net.Dial("tcp", sel.addr)
	msgMng := proto.NewMsgTransfer(conn.TurnConn(cn))
	if err != nil {
		return nil, err
	}
	runProxyReq := &proto.RunProxyReq{
		Cid:   sel.cid,
		Token: sel.token,
		Name:  sel.name,
	}
	err = msgMng.SendMsg(proto.NewPackage(proto.CmdRunProxyReq, runProxyReq))
	if err != nil {
		return nil, err
	}
	return cn, nil
}
