package control

import (
	"ray-seep/ray-seep/common/errs"
	"ray-seep/ray-seep/proto"
	"vilgo/vlog"
)

type ResponsePush interface {
	PushEvent(cmd int32, dt []byte) error
}

type Handler interface {
	Pong(req *proto.Package, push ResponsePush) error
	Login(req *proto.Package, push ResponsePush) error
	LoginRsp(req *proto.Package, push ResponsePush) error
	RegisterProxyRsp(req *proto.Package, push ResponsePush) error
	LogoutRsp(req *proto.Package, push ResponsePush) error
}

type HandlerFun func(req *proto.Package, push ResponsePush) (err error)

type router struct {
	hds map[int32]HandlerFun
}

func (r *router) route(req *proto.Package, push ResponsePush) error {
	hd, ok := r.hds[req.Cmd]
	if !ok {
		return errs.ErrNoCmdRouterNot
	}
	return hd(req, push)
}
func (r *router) Add(cmd int32, fun HandlerFun) {
	r.hds[cmd] = fun
	return
}

type RouteControl struct {
	route    *router
	hd       Handler
	sender   proto.Sender
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

func (r *RouteControl) OnConnect(sender proto.Sender) error {
	r.sender = sender
	return r.hd.Login(&proto.Package{Cmd: proto.CmdLoginReq, Body: []byte{}}, r)
}

func (r *RouteControl) OnMessage(req *proto.Package) {
	//vlog.INFO("收到服务器的信息：cmd[%d]", req.Cmd)
	if err := r.route.route(req, r); err != nil {
		vlog.ERROR("on message exec route [%d] error  %s", req.Cmd, err.Error())
	}
}

func (r *RouteControl) OnDisconnect(localId int64) {
	_ = r.hd.LogoutRsp(&proto.Package{Cmd: proto.CmdLogoutReq, Body: []byte{}}, r)
	return
}

func (r *RouteControl) initRouter() {
	r.route.Add(proto.CmdLoginReq, r.hd.Login)
	r.route.Add(proto.CmdLoginRsp, r.hd.LoginRsp)
	r.route.Add(proto.CmdRegisterProxyRsp, r.hd.RegisterProxyRsp)
	r.route.Add(proto.CmdPong, r.hd.Pong)
}

//
func (r *RouteControl) PushEvent(cmd int32, dt []byte) error {
	return r.pushEvent(cmd, dt)
}

//
func (r *RouteControl) pushEvent(cmd int32, dt []byte) error {
	return r.sender.SendMsg(&proto.Package{Cmd: cmd, Body: dt})
}
