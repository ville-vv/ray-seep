// @File     : main_test
// @Author   : Ville
// @Time     : 19-9-24 下午4:43
// server
package server

import (
	"github.com/vilsongwei/vilgo/vlog"
	"ray-seep/ray-seep/conf"
	"testing"
)

func TestStart(t *testing.T) {
	//vlog.DefaultLogger()
	//node := NewControlServer()
	//node.Start()
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
