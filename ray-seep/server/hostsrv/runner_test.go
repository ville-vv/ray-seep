package hostsrv

import (
	"github.com/vilsongwei/vilgo/vlog"
	"testing"
	"time"
)

func TestRunner_Join(t *testing.T) {
	vlog.DefaultLogger()
	runner := NewRunnerMng()
	runner.Start()
	time.Sleep(time.Second * 1)

	rj := runner.Join()

	j1 := JoinItem{
		Name:   "ddd",
		ConnId: 0,
		Run:    newHttpRunner(":23455", nil),
		Err:    make(chan error),
	}
	rj <- j1
	<-j1.Err

	j2 := JoinItem{
		Name:   "dddddd",
		ConnId: 5,
		Run:    newHttpRunner(":23456", nil),
		Err:    make(chan error),
	}
	rj <- j2
	<-j2.Err

	j3 := JoinItem{
		Name:   "657hth",
		ConnId: 9,
		Run:    newHttpRunner(":23457", nil),
		Err:    make(chan error),
	}
	rj <- j3
	<-j3.Err
	select {}
}
