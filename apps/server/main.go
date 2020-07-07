// @File     : main.go
// @Author   : Ville
// @Time     : 19-9-23 下午6:23
// main
package main

import (
	"flag"
	"fmt"
	"github.com/vilsongwei/vilgo/vlog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"ray-seep/ray-seep/common/util"
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/server"
	"ray-seep/ray-seep/server/env_init"
	"syscall"
	"time"
)

var (
	configPath string
	help       bool
	genCfgFile string
	dbInit     bool
	pprofOpen  bool
)

func printServerInfo(cfg *conf.Server) {
	fmt.Printf(" ==================================================================\n")
	fmt.Printf("\t node server address is [%s:%d]\n", cfg.Host, cfg.Ctl.Port)
	fmt.Printf("\t   proxy server address is [%s:%d]\n", cfg.Host, cfg.Pxy.Port)
	fmt.Printf(" ==================================================================\n")
}

func argsParse() {
	flag.StringVar(&configPath, "c", "", "the config file")
	flag.BoolVar(&help, "h", false, "the tool use help")
	flag.BoolVar(&dbInit, "db-init", false, "create database and table if not exist, must to do with -c point config file")
	flag.StringVar(&genCfgFile, "gen", "", "generate default config file")
	flag.BoolVar(&pprofOpen, "pprof", false, "open the pprof tool")

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

	if pprofOpen {
		pprof()
	}

}

func pprof() {
	go http.ListenAndServe(":8078", nil)
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
