// @File     : client
// @Author   : Ville
// @Time     : 19-9-24 下午4:13 
// client 
package main

import (
	"encoding/binary"
	"net"
	"ray-seep/ray-seep/msg"
	"time"
	"vilgo/vlog"
)
func WriteMsg(pkg *msg.Message, conn net.Conn)(err error){
	data, err := msg.Pack(pkg)
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


func main(){
	vlog.DefaultLogger()
	conn , err:= net.Dial("tcp","127.0.0.1:30080")
	if err != nil{
		vlog.LogE("connect server fail %v", err)
		return
	}
	defer conn.Close()

	go func() {
		for{
			var length int64
			binary.Read(conn, binary.LittleEndian, &length)
			buf := make([]byte, length)
			vlog.INFO("等待接受消息：")
			n, err := conn.Read(buf)
			if err != nil{
				vlog.LogE("发生错误：%v", err)
				return
			}
			if n > 0{
				pkg, err := msg.UnPack(buf)
				if err != nil{
					vlog.LogE("收到消息解析错误")
					return
				}
				vlog.INFO("收到的消息：%v",pkg)
			}
		}

	}()
	auth :=  &msg.Message{
		Cmd:"auth",
	}

	err = WriteMsg(auth, conn)
	if err != nil{
		vlog.ERROR("Write auth message error %v", err)
		return
	}

	ping :=  &msg.Message{
		Cmd:"ping",
	}

	for  {
		sendT := time.NewTicker(time.Second*3)
		select {
		case <- sendT.C:
			vlog.INFO("定时发送ping")
			err = WriteMsg(ping, conn)
			if err != nil{
				vlog.ERROR("Write ping message error %v", err)
				return
			}
		}
	}
}
