// @File     : main
// @Author   : Ville
// @Time     : 19-9-24 下午3:10
// server
package server

import (
	"github.com/vilsongwei/vilgo/vlog"
	"os"
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/databus"
	"ray-seep/ray-seep/server/proxy"
	"ray-seep/ray-seep/server_v2/hostsrv"
	"ray-seep/ray-seep/server_v2/node"
	"runtime/debug"
)

type Server interface {
	Start() error
	Stop()
	Scheme() string
}

type RaySeepServer struct {
	cfg         *conf.Server
	proxy       *proxy.ProxyServer
	control     *node.ControlServer
	start       []string
	stopCh      chan int
	db          databus.BaseDao
	proxyRunner *hostsrv.Runner
}

func NewRaySeepServer(cfg *conf.Server) *RaySeepServer {
	rds := databus.NewDao(cfg)
	runner := hostsrv.NewRunner()
	msgAdopter := node.NewMessageControl(cfg, node.NewPodHandler(rds), runner)
	return &RaySeepServer{
		cfg:         cfg,
		stopCh:      make(chan int, 1),
		proxy:       proxy.NewProxyServer(cfg.Pxy, msgAdopter),
		control:     node.NewNodeServer(cfg.Ctl, msgAdopter),
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
	r.toGo("hostsrv", r.proxyRunner.Start)
	vlog.INFO("server have started success")
	<-r.stopCh
	vlog.INFO("server have stop success")
}

func (r *RaySeepServer) Stop() {
	r.db.Close()
	// 停止已启动的服务
	r.control.Stop()
	r.proxy.Stop()
	close(r.stopCh)
}
