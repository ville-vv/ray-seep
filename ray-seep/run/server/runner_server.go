// @File     : main.go
// @Author   : Ville
// @Time     : 19-9-23 下午6:23
// main
package main

import (
	"os"
	"os/signal"
	"ray-seep/ray-seep/common/conn"
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/server"
	"ray-seep/ray-seep/server/http"
	"ray-seep/ray-seep/server/node"
	"ray-seep/ray-seep/server/proxy"
	"syscall"
	"time"
	"vilgo/vlog"
)

func main() {

	vlog.DefaultLogger()
	cfg := conf.InitServer()

	controlHandler := node.NewMessageAdopter()
	regCenter := proxy.NewRegisterCenter(conn.NewPool(), controlHandler)

	srv := server.NewRaySeepServer(cfg)
	srv.Use(
		node.NewControlServer(cfg.Ctl, controlHandler),
		proxy.NewProxyServer(cfg.Pxy, regCenter),
		http.NewServer(cfg.Http, regCenter),
	)

	go srv.Start()
	// 获取系统信号
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	srv.Stop()
	time.Sleep(time.Millisecond)
}
