// @File     : main_test
// @Author   : Ville
// @Time     : 19-9-24 下午4:43
// server
package server

import (
	"ray-seep/ray-seep/conf"
	"testing"
	"github.com/vilsongwei/vilgo/vlog"
)

func TestStart(t *testing.T) {
	//vlog.DefaultLogger()
	//control := NewControlServer()
	//control.Start()
}

type MockServer struct {
}

func (m *MockServer) Start() error {
	return nil
}

func (m *MockServer) Stop() {
}

func (m *MockServer) Scheme() string {
	return "mock"
}

func TestRaySeepServer_Start(t *testing.T) {
	t.Skip()
	vlog.DefaultLogger()
	srv := NewRaySeepServer(&conf.Server{})
	go srv.Start()
	srv.Stop()
}
