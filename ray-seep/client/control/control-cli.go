package control

import (
	"fmt"
	"io"
	"net"
	"ray-seep/ray-seep/common/conn"
	"ray-seep/ray-seep/common/pkg"
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/mng"
	"sync"
	"time"
	"vilgo/vlog"
)

type MessageHandler interface {
	DealMessage(req *pkg.Package, sender mng.Sender) error
}

type ClientControl struct {
	cfg     *conf.ControlCli
	addr    string
	dealMsg MessageHandler
}

func NewClientControl(cfg *conf.ControlCli, msgHandler MessageHandler) *ClientControl {
	cli := &ClientControl{
		cfg:     cfg,
		addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		dealMsg: msgHandler,
	}
	return cli
}

func (sel *ClientControl) Start() {
	c, err := net.Dial("tcp", sel.addr)
	if err != nil {
		vlog.LogE("connect server fail %v", err)
		return
	}
	defer c.Close()
	msgMng := mng.NewMsgTransfer(conn.TurnConn(c))

	var wg sync.WaitGroup
	recvCh := make(chan pkg.Package)
	cancel := make(chan pkg.Package)
	msgMng.RecvMsgWithChan(&wg, recvCh, cancel)
	for {
		select {
		case ms, ok := <-recvCh:
			if !ok {
				close(cancel)
				return
			}
			if err := sel.dealMsg.DealMessage(&ms, msgMng); err != nil {
				vlog.ERROR("发送消息失败：%s", err.Error())
				close(cancel)
				return
			}
		case <-cancel:
			return
		}
	}
}

func Start() {
	c, err := net.Dial("tcp", "127.0.0.1:30080")
	if err != nil {
		vlog.LogE("connect server fail %v", err)
		return
	}
	defer c.Close()

	msgMng := mng.NewMsgTransfer(conn.TurnConn(c))

	go func() {
		for {
			var msgPkg pkg.Package
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
	auth := pkg.New(pkg.CmdIdentifyReq, []byte(""))
	if err := msgMng.SendMsg(auth); err != nil {
		vlog.ERROR("Write auth message error %v", err)
		return
	}

	ping := pkg.New(pkg.CmdPing, []byte("星跳"))

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
