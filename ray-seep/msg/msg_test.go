// @File     : msg_test
// @Author   : Ville
// @Time     : 19-9-25 下午3:49
// msg
package msg

import "testing"

func TestPack(t *testing.T) {
	orgPkg := &Message{Cmd: "Login", Body: "hello"}
	data, err := Pack(orgPkg)
	if err != nil {
		t.Error(err)
		return
	}
	var dtPkg Message
	if err := UnPack(data, &dtPkg); err != nil {
		t.Error(err)
		return
	}
	if dtPkg.Body != orgPkg.Body {
		t.Error("unpack err Body not right")
		return
	}
	if dtPkg.Cmd != orgPkg.Cmd {
		t.Error("unpack err Cmd not right")
		return
	}
}
