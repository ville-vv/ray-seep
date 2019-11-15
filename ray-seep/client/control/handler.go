package control

import (
	jsoniter "github.com/json-iterator/go"
	"ray-seep/ray-seep/proto"
	"time"
	"vilgo/vlog"
)

type ClientControlHandler struct {
	cId int64
}

func (c *ClientControlHandler) Ping(push ResponsePush) {
	go func() {
		tm := time.NewTicker(time.Millisecond * 3)
		for {
			select {
			case <-tm.C:
				if err := push.PushEvent(proto.CmdPing, nil); err != nil {
					return
				}
			}
		}
	}()
	return
}

func (c *ClientControlHandler) Pong(req *proto.Package, push ResponsePush) (err error) {
	vlog.INFO("server message  pong [%d]", req.Cmd)
	return
}

// 登录服务器
func (c *ClientControlHandler) Login(req *proto.Package, push ResponsePush) (err error) {
	dt, err := jsoniter.Marshal(&proto.LoginReq{Name: "", Password: ""})
	if err != nil {
		vlog.ERROR("push event json marshal error %s", err.Error())
		return err
	}
	return push.PushEvent(proto.CmdLoginReq, dt)
}

func (c *ClientControlHandler) LoginRsp(req *proto.Package, push ResponsePush) (err error) {
	vlog.INFO("login success")
	rsp := &proto.LoginRsp{}
	if err := jsoniter.Unmarshal(req.Body, rsp); err != nil {
		return err
	}
	vlog.INFO("当前被分配的 ID：%d  Token:%s", rsp.Id, rsp.Token)
	c.cId = rsp.Id
	c.Ping(push)
	return
}

func (c *ClientControlHandler) RegisterProxyRsp(req *proto.Package, push ResponsePush) (err error) {
	return nil
}

func (c *ClientControlHandler) LogoutRsp(req *proto.Package, push ResponsePush) (err error) {
	vlog.INFO("disconnect cid:%d", c.cId)
	return nil
}
