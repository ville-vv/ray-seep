// @File     : package_test
// @Author   : Ville
// @Time     : 19-9-24 下午1:35
// msg
package pkg

import (
	"testing"
)

func TestFrame_Pack(t *testing.T) {
	frame := Frame{
		Cmd:  0xffffffff,
		Body: []byte("hello"),
	}
	data := frame.Pack()

	unFrame := &Frame{}
	if err := unFrame.UnPack(data); err != nil {
		t.Error(err)
		return
	}

	if string(unFrame.Body) != string(frame.Body) {
		t.Error("failed")
		return
	}

}

func TestPack(t *testing.T) {
	orgPkg := &Package{Cmd: 1, Body: []byte("hello")}
	data, err := Pack(orgPkg)
	if err != nil {
		t.Error(err)
		return
	}
	var dtPkg Package
	if err := UnPack(data, &dtPkg); err != nil {
		t.Error(err)
		return
	}
	if string(orgPkg.Body) != string(dtPkg.Body) {
		t.Error("unpack err Body not right")
		return
	}
	if orgPkg.Cmd != dtPkg.Cmd {
		t.Error("unpack err Cmd not right")
		return
	}
}
