package control

import (
	"fmt"
	"io"
	"net"
	"ray-seep/ray-seep/common/conn"
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/proto"
	"sync"
	"time"
	"vilgo/vlog"
)

type Router interface {
	OnConnect(sender proto.Sender) error
	OnMessage(req *proto.Package)
	OnDisconnect(id int64)
}

type ClientControl struct {
	cfg    *conf.ControlCli
	addr   string
	hd     Router
	msgMng proto.MsgTransfer
}

func NewClientControl(cfg *conf.ControlCli, hd Handler) *ClientControl {
	cli := &ClientControl{
		cfg:  cfg,
		addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		hd:   NewRouteControl(hd),
	}
	return cli
}

func (sel *ClientControl) Start() {
	c, err := net.Dial("tcp", sel.addr)
	if err != nil {
		vlog.LogE("connect server fail %v", err)
		return
	}
	sel.dealConn(conn.TurnConn(c))
}

func (sel *ClientControl) dealConn(c conn.Conn) {
	defer c.Close()
	sel.msgMng = proto.NewMsgTransfer(c)
	if err := sel.hd.OnConnect(sel.msgMng); err != nil {
		vlog.ERROR("server connect error:%s", err.Error())
		return
	}
	defer sel.hd.OnDisconnect(c.Id())
	var wg sync.WaitGroup
	recvCh := make(chan proto.Package)
	cancel := make(chan interface{})
	wg.Add(1)
	sel.msgMng.RecvMsgWithChan(&wg, recvCh, cancel)
	wg.Wait()
	for {
		select {
		case ms, ok := <-recvCh:
			if !ok {
				vlog.ERROR("关闭连接：%d", c.Id())
				return
			}
			sel.hd.OnMessage(&ms)
		case err := <-cancel:
			vlog.ERROR("被服务器断开：%v", err)
			return
		}
	}
}

func (sel *ClientControl) PushEvent(cmd int32, dt []byte) error {
	return sel.pushEvent(&proto.Package{Cmd: cmd, Body: dt})
}

func (sel *ClientControl) pushEvent(p *proto.Package) error {
	return sel.msgMng.SendMsg(p)
}

func Start() {
	c, err := net.Dial("tcp", "127.0.0.1:30080")
	if err != nil {
		vlog.LogE("connect server fail %v", err)
		return
	}
	defer c.Close()

	msgMng := proto.NewMsgTransfer(conn.TurnConn(c))

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
	auth := proto.New(proto.CmdLoginReq, []byte(""))
	if err := msgMng.SendMsg(auth); err != nil {
		vlog.ERROR("Write auth message error %v", err)
		return
	}

	ping := proto.New(proto.CmdPing, []byte("星跳"))

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
