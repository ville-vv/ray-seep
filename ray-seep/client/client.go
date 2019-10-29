// @File     : client
// @Author   : Ville
// @Time     : 19-9-24 下午4:13
// client
package main

import (
	"fmt"
	"os"
	"os/signal"
	"ray-seep/ray-seep/client/proxy"
	"syscall"
	"vilgo/vlog"
)

func main() {
	vlog.DefaultLogger()
	sgn := make(chan os.Signal, 1)

	proxy.Start()

	signal.Notify(sgn, syscall.SIGABRT, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	fmt.Println(<-sgn)
}
