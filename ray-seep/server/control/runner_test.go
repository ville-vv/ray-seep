package control

import (
	"github.com/vilsongwei/vilgo/vlog"
	"ray-seep/ray-seep/server/http"
	"testing"
)

func TestRunner_Join(t *testing.T) {
	vlog.DefaultLogger()
	runner := NewRunner()
	runner.Start()

	runner.Join() <- JoinItem{
		Name:   "ddd",
		ConnId: 0,
		Run:    http.NewServerWithAddr(":23455", nil),
	}

	runner.Join() <- JoinItem{
		Name:   "dddddd",
		ConnId: 5,
		Run:    http.NewServerWithAddr(":23456", nil),
	}

	runner.Join() <- JoinItem{
		Name:   "dddddd",
		ConnId: 5,
		Run:    http.NewServerWithAddr(":23456", nil),
	}

	select {}
}
