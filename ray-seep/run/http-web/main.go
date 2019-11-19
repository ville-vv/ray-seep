package main

import (
	"fmt"
	"net"
	"ray-seep/ray-seep/common/rayhttp"
)

func main() {
	ls, err := net.Listen("tcp", ":23455")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		c, err := ls.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go func(c net.Conn) {
			defer c.Close()
			fmt.Println("有请求接入：")
			for {
				buf := make([]byte, 1024*2)
				n, err := c.Read(buf)
				if err != nil {
					fmt.Println(err)
					return
				}
				buf = buf[:n]
				fmt.Println(string(buf))
				n, err = c.Write([]byte(rayhttp.Success))
				if err != nil {
					fmt.Println(err)
				}
				break
			}
			fmt.Println("有请求断开：")
		}(c)
	}
}
