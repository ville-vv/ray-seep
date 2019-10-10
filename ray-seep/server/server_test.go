// @File     : main_test
// @Author   : Ville
// @Time     : 19-9-24 下午4:43 
// server 
package server

import (
	"testing"
	"vilgo/vlog"
)

func TestStart(t *testing.T) {
	vlog.DefaultLogger()
	control := NewControlServer()
	control.Start()
}
