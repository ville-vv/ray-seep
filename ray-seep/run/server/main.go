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
	"ray-seep/ray-seep/common/util"
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/server"
	"ray-seep/ray-seep/server/env_init"
	"syscall"
	"time"
	"vilgo/vlog"
)

var (
	configPath string
	help       bool
	genCfgFile string
	dbInit     bool
)

func printServerInfo(cfg *conf.Server) {
	fmt.Printf(" ==================================================================\n")
	fmt.Printf("\t control server address is [%s:%d]\n", cfg.Ctl.Host, cfg.Ctl.Port)
	fmt.Printf("\t   proxy server address is [%s:%d]\n", cfg.Pxy.Host, cfg.Pxy.Port)
	fmt.Printf(" ==================================================================\n")
}

func argsParse() {
	flag.StringVar(&configPath, "c", "", "the config file")
	flag.BoolVar(&help, "h", false, "the tool use help")
	flag.BoolVar(&dbInit, "db-init", false, "create database and table if not exist, must to do with -c point config file")
	flag.StringVar(&genCfgFile, "gen", "", "generate default config file")

	flag.Parse()
	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}
	if genCfgFile != "" {
		conf.GenDefServerConfigFile(genCfgFile)
		os.Exit(0)
	}
	if dbInit {
		vlog.DefaultLogger()
		env_init.InitDb(conf.InitServer(configPath))
		os.Exit(0)
	}

}

func main() {
	argsParse()
	cfg := conf.InitServer(configPath)
	vlog.DefaultLogger()
	if cfg.Log != nil {
		fmt.Println("reset logger:", *cfg.Log)
		vlog.SetLogger(vlog.NewGoLogger(cfg.Log))
	}
	printServerInfo(cfg)
	srv := server.NewRaySeepServer(cfg)
	go srv.Start()
	_ = util.WritePid("process_id")
	// 获取系统信号
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	srv.Stop()
	time.Sleep(time.Millisecond)
}
