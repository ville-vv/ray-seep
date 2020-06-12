package server_v2

import (
	"github.com/vilsongwei/vilgo/vlog"
	"ray-seep/ray-seep/conf"
	"testing"
)

func TestRaySeepServer_Start(t *testing.T) {
	t.Skip("")
	vlog.DefaultLogger()
	srv := NewRaySeepServer(conf.InitServer(""))
	srv.Start()
}
