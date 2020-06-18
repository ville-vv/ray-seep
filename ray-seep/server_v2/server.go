// @File     : main
// @Author   : Ville
// @Time     : 19-9-24 下午3:10
// server
package server_v2

import (
	"github.com/vilsongwei/vilgo/vlog"
	"os"
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/databus"
	"ray-seep/ray-seep/server_v2/hostsrv"
	"ray-seep/ray-seep/server_v2/node"
	"ray-seep/ray-seep/server_v2/proxy"
	"runtime/debug"
)

type Server interface {
	Start() error
	Stop()
	Scheme() string
}

type RaySeepServer struct {
	cfg         *conf.Server
	nodeServer  *ControlServer
	pxyServer   *ControlServer
	start       []string
	stopCh      chan int
	db          databus.BaseDao
	proxyRunner hostsrv.HostServer
}

func NewRaySeepServer(cfg *conf.Server) *RaySeepServer {
	rds := databus.NewDao(cfg)
	hsr := hostsrv.NewHostService()
	register := proxy.NewRegisterCenter(cfg.Pxy.UserMaxProxyNum)
	hsr.SetDstConn(register)

	pxySrv := proxy.NewPxyManager(cfg.Pxy, register)
	cctSrv := node.NewConnectCenter(cfg, hsr, node.NewPodHandler(rds))
	register.SetNoticeGetter(cctSrv)

	return &RaySeepServer{
		cfg:         cfg,
		stopCh:      make(chan int, 1),
		nodeServer:  NewControlServer(cfg.Ctl, cctSrv),
		pxyServer:   NewControlServer(cfg.Pxy, pxySrv),
		db:          rds,
		proxyRunner: hsr,
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
	r.toGo(r.nodeServer.Scheme(), r.nodeServer.Start)
	r.toGo(r.pxyServer.Scheme(), r.pxyServer.Start)
	r.toGo("hostsrv", r.proxyRunner.Start)
	vlog.INFO("server have started success")
	<-r.stopCh
	vlog.INFO("server have stop success")
}

func (r *RaySeepServer) Stop() {
	r.db.Close()
	r.nodeServer.Stop()
	r.pxyServer.Stop()
	r.proxyRunner.Stop()
	// 停止已启动的服务
	close(r.stopCh)
}
