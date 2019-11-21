// @File     : main.go
// @Author   : Ville
// @Time     : 19-9-23 下午6:23
// main
package main

import (
	"os"
	"os/signal"
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/server"
	"ray-seep/ray-seep/server/http"
	"ray-seep/ray-seep/server/node"
	"ray-seep/ray-seep/server/online"
	"ray-seep/ray-seep/server/proxy"
	"syscall"
	"time"
	"vilgo/vlog"
)

func main() {

	vlog.DefaultLogger()
	cfg := conf.InitServer()

	userMng := online.NewUserManager()
	msgAdopter := node.NewMessageAdopter(cfg, userMng)
	//regCenter := node.NewRegisterCenter(100, msgAdopter, userMng)

	srv := server.NewRaySeepServer(cfg)
	srv.Use(
		node.NewControlServer(cfg.Ctl, msgAdopter),
		proxy.NewProxyServer(cfg.Pxy, msgAdopter),
		http.NewServer(cfg.Http, msgAdopter),
	)

	go srv.Start()
	// 获取系统信号
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	srv.Stop()
	time.Sleep(time.Millisecond)
}
