package control

import (
	"ray-seep/ray-seep/common/errs"
	"ray-seep/ray-seep/proto"
	"vilgo/vlog"
)

type HandlerFun func(req *proto.Package) (rsp *proto.Package, err error)

type Router struct {
	hds map[int32]HandlerFun
}

func (r *Router) route(req *proto.Package) (rsp *proto.Package, err error) {
	hd, ok := r.hds[int32(req.Cmd)]
	if !ok {
		return nil, errs.ErrNoCmdRouterNot
	}
	return hd(req)
}
func (r *Router) Add(cmd int32, fun HandlerFun) {
	r.hds[cmd] = fun
	return
}

type RouteControl struct {
	route *Router
}

func NewRouteControl() *RouteControl {
	r := &RouteControl{route: &Router{hds: make(map[int32]HandlerFun)}}
	r.initRouter()
	return r
}

func (r *RouteControl) OnConnect(sender proto.Sender) error {
	return nil
}

func (r *RouteControl) OnMessage(req *proto.Package) (rsp *proto.Package, err error) {
	return
}

func (r *RouteControl) OnDisconnect(id int64) {
	vlog.INFO("disconnect %d", id)
	return
}

func (r *RouteControl) initRouter() {
	r.route.Add(proto.CmdRegisterProxyReq, r.RegisterProxy)
	r.route.Add(proto.CmdPing, r.Ping)
}
func (r *RouteControl) RegisterProxy(p *proto.Package) (rsp *proto.Package, err error) {
	return
}
func (r *RouteControl) Ping(p *proto.Package) (rsp *proto.Package, err error) {
	return
}
