package control

import (
	"ray-seep/ray-seep/common/pkg"
)

type HandlerFun func(req *pkg.Package) (rsp *pkg.Package, err error)

type RouteControl struct {
	route *Router
}

func (r *RouteControl) InitRouter() {
	r.route.Cmd(int32(pkg.CmdRegisterProxyReq), r.RegisterProxy)
}
func (r *RouteControl) RegisterProxy(p *pkg.Package) (rsp *pkg.Package, err error) {
	return
}
