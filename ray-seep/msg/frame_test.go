// @File     : package_test
// @Author   : Ville
// @Time     : 19-9-24 下午1:35
// msg
package msg

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
