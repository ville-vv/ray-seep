// @File     : main
// @Author   : Ville
// @Time     : 19-9-24 下午3:10
// server
package server

import (
	"os"
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/dao"
	"ray-seep/ray-seep/server/control"
	"ray-seep/ray-seep/server/http"
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
	cfg         *conf.Server
	proxy       *proxy.ProxyServer
	http        *http.Server
	control     *control.NodeServer
	start       []string
	stopCh      chan int
	db          dao.BaseDao
	proxyRunner *control.Runner
}

func NewRaySeepServer(cfg *conf.Server) *RaySeepServer {

	rds := dao.NewDao(cfg)
	runner := control.NewRunner()
	msgAdopter := control.NewMessageControl(cfg, control.NewPodHandler(rds), runner)

	return &RaySeepServer{
		cfg:         cfg,
		stopCh:      make(chan int, 1),
		http:        http.NewServer(cfg.Proto, msgAdopter),
		proxy:       proxy.NewProxyServer(cfg.Pxy, msgAdopter),
		control:     control.NewNodeServer(cfg.Ctl, msgAdopter),
		db:          rds,
		proxyRunner: runner,
	}
}

func (r *RaySeepServer) toGo(name string, f func() error) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				r.Stop()
				debug.PrintStack()
				os.Exit(2)
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
	r.toGo("runner", r.proxyRunner.Start)
	vlog.INFO("server have started success")
	<-r.stopCh
	vlog.INFO("server have stop success")
}

func (r *RaySeepServer) Stop() {
	r.db.Close()
	// 停止已启动的服务
	r.control.Stop()
	r.proxy.Stop()
	r.http.Stop()
	close(r.stopCh)
}
