package proxy

import (
	"net"
	"ray-seep/ray-seep/common/conn"
	"ray-seep/ray-seep/common/pkg"
	"ray-seep/ray-seep/mng"
	"vilgo/vlog"
)

func Start() {

	cn, err := net.Dial("tcp", ":32202")
	if err != nil {
		vlog.ERROR("connect to proxy server error %s", err.Error())
		return
	}
	msgMng := mng.NewMsgTransfer(conn.TurnConn(cn))
	msgMng.SendMsg(pkg.NewWithObj(pkg.CmdRegisterProxyReq, &pkg.RegisterProxyReq{Cid: 89797, SubDomain: "test"}))
}
