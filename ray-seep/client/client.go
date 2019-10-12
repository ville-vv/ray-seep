// @File     : client
// @Author   : Ville
// @Time     : 19-9-24 下午4:13
// client
package main

import (
	"io"
	"net"
	"ray-seep/ray-seep/common/conn"
	pkg2 "ray-seep/ray-seep/common/pkg"
	"ray-seep/ray-seep/mng"
	"time"
	"vilgo/vlog"
)

func main() {
	vlog.DefaultLogger()
	c, err := net.Dial("tcp", "127.0.0.1:30080")
	if err != nil {
		vlog.LogE("connect server fail %v", err)
		return
	}
	defer c.Close()

	msgMng := mng.NewMsgTransfer(conn.TurnConn(c))

	go func() {
		for {
			var pkg pkg2.Package
			vlog.INFO("等待接受消息：")
			if err := msgMng.RecvMsg(&pkg); err != nil {
				if err == io.EOF {
					vlog.INFO("收到关闭链接：")
					return
				}
				vlog.LogE("发生错误：%v", err)
				continue
			}
			vlog.INFO("收到的消息 cmd[%d] body:%s", pkg.Cmd, string(pkg.Body))
		}

	}()
	auth := pkg2.New(pkg2.CmdIdentifyReq, []byte(""))
	if err := msgMng.SendMsg(auth); err != nil {
		vlog.ERROR("Write auth message error %v", err)
		return
	}

	ping := pkg2.New(pkg2.CmdPing, []byte("星跳"))

	for {
		sendT := time.NewTicker(time.Second * 3)
		select {
		case <-sendT.C:
			vlog.INFO("定时发送ping")
			err = msgMng.SendMsg(ping)
			if err != nil {
				vlog.ERROR("Write ping message error %v", err)
				return
			}
		}
	}
}
