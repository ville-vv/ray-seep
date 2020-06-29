package control

import (
	"ray-seep/ray-seep/common/errs"
	"ray-seep/ray-seep/msg"
	"ray-seep/ray-seep/proto"
)

type ResponsePush interface {
	PushEvent(cmd int32, dt []byte) error
}

type Handler interface {
	Pong(req *msg.Package) error
	Login(push ResponsePush) error
	LoginRsp(req *msg.Package) error
	CreateHostRsp(req *msg.Package) (err error)
	RunProxyRsp(req *msg.Package) error
	NoticeRunProxy(req *msg.Package) error
	LogoutRsp(req *msg.Package) error
	NoticeError(req *msg.Package) (err error)
}

type HandlerFun func(req *msg.Package) (err error)

type router struct {
	hds map[int32]HandlerFun
}

func (r *router) route(req *msg.Package, push ResponsePush) error {
	hd, ok := r.hds[req.Cmd]
	if !ok {
		return errs.ErrNoCmdRouterNot
	}
	return hd(req)
}
func (r *router) Add(cmd int32, fun HandlerFun) {
	r.hds[cmd] = fun
	return
}

type RouteControl struct {
	route    *router
	hd       Handler
	sender   msg.ResponseSender
	remoteId int64
	token    string
}

func NewRouteControl(hd Handler) *RouteControl {
	r := &RouteControl{
		route: &router{hds: make(map[int32]HandlerFun)},
		hd:    hd,
	}
	r.initRouter()
	return r
}

func (r *RouteControl) OnConnect(sender msg.ResponseSender) error {
	r.sender = sender
	return r.hd.Login(r)
}

func (r *RouteControl) OnMessage(req *msg.Request) error {
	//vlog.INFO("收到服务器的信息：cmd[%d]", req.head)
	return r.route.route(req.Body, r)
}

func (r *RouteControl) OnDisconnect(localId int64) {
	_ = r.hd.LogoutRsp(&msg.Package{Cmd: proto.CmdLogoutReq, Body: []byte{}})
	return
}

func (r *RouteControl) initRouter() {
	r.route.Add(proto.CmdPong, r.hd.Pong)
	// 登录返回
	r.route.Add(proto.CmdLoginRsp, r.hd.LoginRsp)
	//
	r.route.Add(proto.CmdCreateHostRsp, r.hd.CreateHostRsp)
	//
	r.route.Add(proto.CmdRunProxyRsp, r.hd.RunProxyRsp)

	r.route.Add(proto.CmdNoticeRunProxy, r.hd.NoticeRunProxy)

	r.route.Add(proto.CmdError, r.hd.NoticeError)

}

//
func (r *RouteControl) PushEvent(cmd int32, dt []byte) error {
	return r.pushEvent(cmd, dt)
}

//
func (r *RouteControl) pushEvent(cmd int32, dt []byte) error {
	return r.sender.Send(&msg.Package{Cmd: cmd, Body: dt})
}
