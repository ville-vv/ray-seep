package main

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"ray-seep/ray-seep/common/rayhttp"
	"text/template"
	"time"
)

var htmlPageTemp = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Ray Seep Test</title>
</head>
<body>

<div style="margin: 0px 100px 200px 300px">
	<div style="">
		<h1> Ray Seep(射线渗透) Test</h1>
		<span> 你可以再文件中添加几个文件然后刷新网页看看有什么变化。</span>
		<p>1、这个是一个部署在本地的内网web服务网站，此服务网站试用以体验 Ray Seep 工具的内网穿透能力。</p>
		<p>2、转发次网站的数据技术是由 Ray-Seep 工具实现。</p>
		<p>3、部署此工具可以学生共享外网服务器，在宿舍自己电脑学习开发web应用，并且可以在公网访问。</p>
		<p>4、部署此工具可以提供企业，在本地开发软件，并可以与第三方公司联调本地HTTP API接口。</p>
		<h4>这个网站的示例功能是显示当前运行环境目录的所有文件名称，以下是你当前软件目录下的文件：<h4>
	</div>
	
	<div style="background-color: #c6d5e9">
	{{.Context}}
	<div>
</div>

</body>
</html>
`

func main() {
	HttpServer()
}

func HttpServer() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(time.Now().String(), "access from remoter address : ", request.RemoteAddr)
		fmt.Println(time.Now().String(), "access User-Agent : ", request.Header.Get("User-Agent"))

		buf := bytes.NewBufferString("")
		dep := 0
		err := filepath.Walk("./", func(path string, f os.FileInfo, err error) error {
			if f == nil && dep > 1 {
				return err
			}
			dep += 1
			if f.IsDir() {
				//files = files + path + "\n"
				buf.Write([]byte(path + "</br>"))
				return nil
			}
			//files = files + path + "\n"
			//writer.Write([]byte(path + "\n"))
			buf.Write([]byte(path + "</br>"))
			return nil
		})
		if err != nil {
			writer.Write([]byte(err.Error()))
		}
		tmpl := template.New("test")
		tmpl, err = tmpl.Parse(htmlPageTemp)
		if err != nil {
			writer.Write(buf.Bytes())
			return
		}
		tmpl.Execute(writer, map[string]string{"Context": buf.String()})
	})
	fmt.Println("本地web服务启动，请使用 http://localhost:12345 在浏览器中访问\r\n如果配套 RaySeep 使用请在浏览器中打开 ray-seep-cli 中输出的 http 地址")
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
