// @File     : main
// @Author   : Ville
// @Time     : 19-9-24 下午3:10
// server
package server

import (
	"ray-seep/ray-seep/server/http"
	"ray-seep/ray-seep/server/node"
	"vilgo/vlog"
)

type Server struct {
}

func Start() {
	vlog.DefaultLogger()
	go func() {
		control := node.NewConnServer()
		control.Start()
	}()
	hserver := http.NewServer()
	hserver.Start()
}
