// @File     : msg
// @Author   : Ville
// @Time     : 19-9-24 下午2:59 
// msg 
package msg

import (
	"encoding/binary"
	jsoniter "github.com/json-iterator/go"
	"ray-seep/ray-seep/mng"
	"vilgo/vlog"
)

type Message struct {
	Cmd  string `json:"cmd"`
	Body interface{} `json:"body"`
}

func UnPack(data []byte)(pkg Message, err error){
	frame  := Frame{}
	if _, err = frame.UnPack(data); err != nil{
		return
	}

	if len(frame.Body)> 0{
		if err = jsoniter.Unmarshal(frame.Body, &pkg);err != nil{
			return
		}
	}

	return
}

func Pack(pkg *Message)(data []byte ,  err error){
	frame  := Frame{}
	frame.Body , err = jsoniter.Marshal(pkg)
	if err != nil{
		return nil, err
	}
	return frame.Pack(),nil
}

// 接收消息
func RecvMsg(c mng.Conn)(pkg Message, err error){
	var l int64
	err = binary.Read(c, binary.LittleEndian, &l)
	if err != nil{
		return
	}
	buf := make([]byte, l)
	n , err := c.Read(buf)
	if err !=nil{
		return
	}
	return UnPack(buf[:n])
}

// 发送消息
func SendMsg(pkg *Message, conn mng.Conn)(err error){
	data, err := Pack(pkg)
	if err != nil{
		vlog.ERROR("Pack message error %v", err)
		return
	}

	err = binary.Write(conn, binary.LittleEndian, int64(len(data)))
	if err != nil {
		return
	}

	_, err = conn.Write(data)
	return
}

