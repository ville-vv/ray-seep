package proxy

import (
	"io"
	"net"
	"ray-seep/ray-seep/common/conn"
	"ray-seep/ray-seep/proto"
	"vilgo/vlog"
)

func Start() {

	cn, err := net.Dial("tcp", ":32202")
	if err != nil {
		vlog.ERROR("connect to proxy server error %s", err.Error())
		return
	}
	msgMng := proto.NewMsgTransfer(conn.TurnConn(cn))
	go func() {
		for {
			var msgPkg proto.Package
			vlog.INFO("等待接受消息：")
			if err := msgMng.RecvMsg(&msgPkg); err != nil {
				if err == io.EOF {
					vlog.INFO("收到关闭链接：")
					return
				}
				vlog.LogE("发生错误：%v", err)
				continue
			}
			vlog.INFO("收到的消息 cmd[%d] body:%s", msgPkg.Cmd, string(msgPkg.Body))
		}

	}()

	err = msgMng.SendMsg(proto.NewWithObj(proto.CmdRegisterProxyReq, &proto.RegisterProxyReq{Cid: 89797, SubDomain: "test"}))
	if err != nil {
		vlog.ERROR("send register proxy error %s", err.Error())
		return
	}

}
