// @File     : package_test
// @Author   : Ville
// @Time     : 19-9-24 下午1:35 
// msg 
package msg

import (
	"testing"
)

func TestPack(t *testing.T) {
	orgPkg := &Message{Cmd:"Login",Body:"hello"}
	data, err := Pack(orgPkg)
	if err != nil{
		t.Error(err)
		return
	}
	dtPkg , err:= UnPack(data)
	if err != nil{
		t.Error(err)
		return
	}
	if dtPkg.Body != orgPkg.Body{
		t.Error("unpack err Body not right")
		return
	}
	if dtPkg.Cmd != orgPkg.Cmd{
		t.Error("unpack err Cmd not right")
		return
	}
}

func TestFrame_Pack(t *testing.T) {
	frame  := Frame{
		Cmd:0xffffffff,
		Body:[]byte("hello"),
	}
	data := frame.Pack()

	unFrame := &Frame{}
	dtBuf,err := unFrame.UnPack(data)
	if err != nil{
		t.Error(err)
		return
	}

	if string(dtBuf) != string(frame.Body){
		t.Error("failed")
		return
	}

}

