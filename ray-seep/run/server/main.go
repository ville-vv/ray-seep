// @File     : main.go
// @Author   : Ville
// @Time     : 19-9-23 下午6:23
// main
package main

import (
	"flag"
	"os"
	"os/signal"
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/server"
	"syscall"
	"time"
	"vilgo/vlog"
)

var (
	configPath = ""
	help       bool
)

func printServerInfo(cfg *conf.Server) {
	vlog.INFO("\t ==========================================================================")
	vlog.INFO("\t\t control server address is [%s:%d] ", cfg.Ctl.Host, cfg.Ctl.Port)
	vlog.INFO("\t\t   proxy server address is [%s:%d]", cfg.Pxy.Host, cfg.Pxy.Port)
	vlog.INFO("\t\t    http server address is [%s:%d]", cfg.Proto.Host, cfg.Proto.Port)
	vlog.INFO("\t ==========================================================================")
}

func main() {
	flag.StringVar(&configPath, "c", "", "the config file")
	flag.BoolVar(&help, "h", false, "the tool use help")
	flag.Parse()
	if help {
		flag.PrintDefaults()
		return
	}
	cfg := conf.InitServer(configPath)
	vlog.DefaultLogger()
	printServerInfo(cfg)
	srv := server.NewRaySeepServer(cfg)
	go srv.Start()
	// 获取系统信号
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	srv.Stop()
	time.Sleep(time.Millisecond)
}
