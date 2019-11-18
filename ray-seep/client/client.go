// @File     : client
// @Author   : Ville
// @Time     : 19-9-24 下午4:13
// client
package client

import (
	"fmt"
	"os"
	"os/signal"
	"ray-seep/ray-seep/client/control"
	"ray-seep/ray-seep/conf"
	"syscall"
	"vilgo/vlog"
)

type RaySeepClient struct {
	ctl *control.ClientControl
}

func Main() {
	vlog.DefaultLogger()
	// 初始化配置
	cfg := conf.InitClient()
	ctrCli := control.NewClientControl(cfg.Control, control.NewClientControlHandler(cfg))

	go func() {
		ctrCli.Start()
	}()

	//go func() {
	//	proxy.Start()
	//}()

	sgn := make(chan os.Signal, 1)
	signal.Notify(sgn, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	fmt.Println(<-sgn)
}
