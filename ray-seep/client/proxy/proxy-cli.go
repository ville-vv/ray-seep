package proxy

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/vilsongwei/vilgo/vlog"
	"net"
	"ray-seep/ray-seep/common/conn"
	"ray-seep/ray-seep/common/repeat"
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/msg"
	"ray-seep/ray-seep/proto"
)

type ClientProxy struct {
	cfg        *conf.ProxyCli
	addr       string
	cid        int64
	token      string
	name       string
	httpDomain string
	stopCh     chan int
	netRet     *repeat.NetRepeater
}

func NewClientProxy(stopCh chan int, cfg *conf.Client) *ClientProxy {
	return &ClientProxy{
		stopCh: stopCh,
		netRet: repeat.NewNetRepeater(NewTunnel(cfg.Web.Address)),
	}
}

func (sel *ClientProxy) RunProxy(id int64, token string, httpDomain string, pxyAddr string) error {
	sel.cid = id
	sel.token = token
	sel.name = httpDomain
	sel.addr = pxyAddr
	go func() {
		defer func() {
			if err := recover(); err != nil {
				vlog.ERROR("%v", err)
			}
		}()
		cn, err := sel.dial()
		if err != nil {
			vlog.ERROR("connect to proxy server error %s", err.Error())
			return
		}
		defer cn.Close()
		reqLength, respLength, err := sel.netRet.Transfer(httpDomain, cn)
		if err != nil {
			vlog.ERROR("forwarding data to local service failed %s", err.Error())
			return
		}
		vlog.INFO("request size：[%d]. response size：[%d]", reqLength, respLength)
	}()
	return nil
}

func (sel *ClientProxy) dial() (net.Conn, error) {
	cn, err := net.Dial("tcp", sel.addr)
	if err != nil {
		vlog.ERROR("connect your local server fail：%s", sel.addr)
		return nil, err
	}
	msgMng := msg.MessagePusher{ResponseSender: msg.NewMessageCenter(conn.TurnConn(cn))}
	data, _ := jsoniter.Marshal(&proto.RunProxyReq{
		Cid:   sel.cid,
		Token: sel.token,
		Name:  sel.name,
	})
	return cn, msgMng.Send(&msg.Package{
		Cmd:  msg.CmdRunProxyReq,
		Body: data,
	})
}
