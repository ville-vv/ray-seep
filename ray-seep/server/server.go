// @File     : main
// @Author   : Ville
// @Time     : 19-9-24 下午3:10
// server
package server

import (
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/server/control"
	"ray-seep/ray-seep/server/http"
	"ray-seep/ray-seep/server/online"
	"ray-seep/ray-seep/server/proxy"
	"runtime/debug"
	"vilgo/vlog"
)

type Server interface {
	Start() error
	Stop()
	Scheme() string
}

type RaySeepServer struct {
	cfg     *conf.Server
	proxy   *proxy.ProxyServer
	http    *http.Server
	control *control.NodeServer
	start   []string
	stopCh  chan int
}

func NewRaySeepServer(cfg *conf.Server) *RaySeepServer {

	msgAdopter := control.NewMessageControl(cfg, online.NewUserManager())

	return &RaySeepServer{
		cfg:     cfg,
		stopCh:  make(chan int, 1),
		http:    http.NewServer(cfg.Http, msgAdopter),
		proxy:   proxy.NewProxyServer(cfg.Pxy, msgAdopter),
		control: control.NewNodeServer(cfg.Ctl, msgAdopter),
	}
}

func (r *RaySeepServer) toGo(name string, f func() error) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				r.Stop()
				debug.PrintStack()
			}
		}()
		if err := f(); err != nil {
			panic(err)
		}
	}()
}

func (r *RaySeepServer) Start() {
	r.toGo(r.control.Scheme(), r.control.Start)
	r.toGo(r.proxy.Scheme(), r.proxy.Start)
	r.toGo(r.http.Scheme(), r.http.Start)
	<-r.stopCh
	vlog.INFO("server [all] have stop success")
}

func (r *RaySeepServer) Stop() {
	// 停止已启动的服务
	r.control.Stop()
	r.proxy.Stop()
	r.http.Stop()
	close(r.stopCh)
}
