package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"ray-seep/ray-seep/common/rayhttp"
)

func main() {
	HttpServer()
}

func HttpServer() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		//files := ""
		dep := 0
		err := filepath.Walk("./", func(path string, f os.FileInfo, err error) error {
			if f == nil && dep > 1 {
				return err
			}
			dep += 1
			if f.IsDir() {
				//files = files + path + "\n"
				writer.Write([]byte(path + "\n"))
				return nil
			}
			//files = files + path + "\n"
			writer.Write([]byte(path + "\n"))
			return nil
		})
		if err != nil {
			writer.Write([]byte(err.Error()))
		}
		//files, err := ioutil.ReadDir(".")
		//if err != nil {
		//	writer.Write([]byte(err.Error()))
		//}
		//for _, f := range files {
		//	fmt.Println(f.Name())
		//	writer.Write([]byte(f.Name()))
		//}

	})
	http.ListenAndServe(":12345", nil)
	fmt.Println("结束了？")
}

func tcpServer() {
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
