package control

import (
	jsoniter "github.com/json-iterator/go"
	"ray-seep/ray-seep/client/proxy"
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/proto"
	"time"
	"vilgo/vlog"
)

type ClientControlHandler struct {
	cId          int64
	token        string
	domain       string
	name         string // 子域名
	push         ResponsePush
	cliPxy       *proxy.ClientProxy
	cliPxyStopCh chan int
}

func NewClientControlHandler(cfg *conf.Client) *ClientControlHandler {
	stopCh := make(chan int)
	return &ClientControlHandler{
		name:         cfg.Control.Name,
		cliPxyStopCh: stopCh,
		cliPxy:       proxy.NewClientProxy(stopCh, cfg),
	}
}

func (c *ClientControlHandler) Ping() {
	go func() {
		tm := time.NewTicker(time.Second * 3)
		for {
			select {
			case <-tm.C:
				if err := c.push.PushEvent(proto.CmdPing, nil); err != nil {
					return
				}
			}
		}
	}()
	return
}

func (c *ClientControlHandler) Pong(req *proto.Package) (err error) {
	//vlog.INFO("server message  pong [%d]", req.Cmd)
	return
}

// 登录服务器
func (c *ClientControlHandler) Login(push ResponsePush) (err error) {
	dt, err := jsoniter.Marshal(&proto.LoginReq{UserId: 1234, Name: c.name, AppId: ""})
	if err != nil {
		vlog.ERROR("push event json marshal error %s", err.Error())
		return err
	}
	c.push = push
	return c.push.PushEvent(proto.CmdLoginReq, dt)
}

func (c *ClientControlHandler) LoginRsp(req *proto.Package) (err error) {
	//vlog.INFO("login success body:%s", string(req.Body))
	rsp := &proto.LoginRsp{}
	if err := jsoniter.Unmarshal(req.Body, rsp); err != nil {
		return err
	}
	vlog.INFO("login success ID:%d  Token:%s", rsp.Id, rsp.Token)
	c.cId = rsp.Id
	c.token = rsp.Token
	c.Ping()
	return c.CreateHostReq()
}

//  CreateHostReq 创建服务主机
func (c *ClientControlHandler) CreateHostReq() error {
	reqData, err := jsoniter.Marshal(proto.CreateHostReq{SubDomain: c.name})
	if err != nil {
		return err
	}
	return c.push.PushEvent(proto.CmdCreateHostReq, reqData)
}

// CreateHostRsp 创建服务主机返回
func (c *ClientControlHandler) CreateHostRsp(req *proto.Package) (err error) {
	//vlog.INFO("收到 [CreateHostRsp]Cmd:%d Body:%s", req.Cmd, string(req.Body))
	ctInfo := &proto.CreateHostRsp{}
	if err = jsoniter.Unmarshal(req.Body, ctInfo); err != nil {
		vlog.ERROR("create host response json un parse error %s", err.Error())
		return
	}
	vlog.INFO("[%d] create host success, domain is [%s]", c.cId, ctInfo.Domain)
	c.domain = ctInfo.Domain
	// 收到创建主机的返回信息就可 运行代理了
	return c.RunProxyReq()
}

// NoticeRunProxy 通知创建代理
func (c *ClientControlHandler) NoticeRunProxy(req *proto.Package) error {
	//vlog.INFO("收到 [NoticeRunProxy]Cmd:%d Body:%s", req.Cmd, string(req.Body))
	return c.RunProxyReq()
}

func (c *ClientControlHandler) RunProxyReq() (err error) {
	return c.cliPxy.RunProxy(c.cId, c.token, c.name)
}

func (c *ClientControlHandler) RunProxyRsp(req *proto.Package) (err error) {
	//vlog.INFO("收到 [RunProxyRsp]Cmd:%d Body:%s", req.Cmd, string(req.Body))
	return nil
}

func (c *ClientControlHandler) LogoutRsp(req *proto.Package) (err error) {
	vlog.INFO("disconnect cid:%d", c.cId)
	return nil
}
