// @File     : main
// @Author   : Ville
// @Time     : 19-9-24 下午3:10
// server
package server

import (
	"ray-seep/ray-seep/common/conn"
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/server/http"
	"ray-seep/ray-seep/server/node"
	"ray-seep/ray-seep/server/proxy"
	"sync"
	"vilgo/vlog"
)

type Server struct {
	srvCnf *conf.Server
}

func Start() {
	vlog.DefaultLogger()
	cfg := conf.InitServer()

	regCenter := proxy.NewRegisterCenter(conn.NewPool())

	wait := sync.WaitGroup{}
	wait.Add(1)
	go func() {
		wait.Done()
		control := node.NewControlServer(cfg.Ctl)
		control.Start()
	}()
	wait.Add(1)
	go func() {
		wait.Done()
		pxy := proxy.NewProxyServer(cfg.Pxy, regCenter)
		pxy.Start()
	}()
	wait.Wait()
	hserver := http.NewServer(cfg.Http, regCenter)
	hserver.Start()
}
