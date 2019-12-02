// @File     : main.go
// @Author   : Ville
// @Time     : 19-9-23 下午6:23
// main
package main

import (
	"flag"
	"fmt"
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
	genCfgFile = ""
)

func printServerInfo(cfg *conf.Server) {
	fmt.Printf(" ==================================================================\n")
	fmt.Printf("\t control server address is [%s:%d]\n", cfg.Ctl.Host, cfg.Ctl.Port)
	fmt.Printf("\t   proxy server address is [%s:%d]\n", cfg.Pxy.Host, cfg.Pxy.Port)
	fmt.Printf("\t    http server address is [%s:%d]\n", cfg.Proto.Host, cfg.Proto.Port)
	fmt.Printf(" ==================================================================\n")
}

func main() {
	flag.StringVar(&configPath, "c", "", "the config file")
	flag.BoolVar(&help, "h", false, "the tool use help")
	flag.StringVar(&genCfgFile, "gen", "", "generate default config file")
	flag.Parse()
	if help {
		flag.PrintDefaults()
		return
	}
	if genCfgFile != "" {
		conf.GenDefServerConfigFile(genCfgFile)
		return
	}

	cfg := conf.InitServer(configPath)
	vlog.DefaultLogger()
	if cfg.Log != nil {
		fmt.Println("reset logger:", *cfg.Log)
		vlog.SetLogger(vlog.NewGoLogger(cfg.Log))
	}
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
