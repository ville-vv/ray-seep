// @File     : client
// @Author   : Ville
// @Time     : 19-9-24 下午4:13
// client
package client

import (
	"flag"
	"fmt"
	"github.com/vilsongwei/vilgo/vlog"
	"os"
	"os/signal"
	"ray-seep/ray-seep/client/control"
	"ray-seep/ray-seep/conf"
	"syscall"
)

var (
	configPath string
	help       bool
	genCfgFile string
	dbInit     bool
)

func argsParse() {
	flag.StringVar(&configPath, "c", "", "the config file")
	flag.BoolVar(&help, "h", false, "the tool use help")
	flag.StringVar(&genCfgFile, "gen", "", "generate default config file")

	flag.Parse()
	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}
	if genCfgFile != "" {
		conf.GenDefClientConfigFile(genCfgFile)
		os.Exit(0)
	}
}

func Main() {
	argsParse()
	// 初始化配置
	cfg := conf.InitClient(configPath)
	vlog.DefaultLogger()
	if cfg.Log != nil {
		fmt.Println("reset logger:", *cfg.Log)
		vlog.SetLogger(vlog.NewGoLogger(cfg.Log))
	}
	ctrCli := control.NewClientManager(cfg.Control, control.NewClientControlHandler(cfg))

	go func() {
		ctrCli.Start()
	}()
	sgn := make(chan os.Signal, 1)
	signal.Notify(sgn, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	select {
	case <-sgn:
		ctrCli.Stop()
	case <-ctrCli.WaitClose():
	}
}
