package main

import (
	"flag"
	"fmt"
	static2 "ray-seep/run/http-web/static"
)

func main() {
	var (
		rootPath string
		address  string
		help     bool
	)
	flag.StringVar(&rootPath, "root", "./", "文件服务的根目录，默认为当前目录")
	flag.StringVar(&address, "server", "0.0.0.0:12345", "服务的地址，默认为 0.0.0.0:12345")
	flag.BoolVar(&help, "addr", false, "帮助提示")
	flag.Parse()
	if help {
		flag.PrintDefaults()
		return
	}
	fmt.Println("本地web服务启动，请使用 http://localhost:12345 在浏览器中访问\r\n如果配套 RaySeep 使用请在浏览器中打开 ray-seep-cli 中输出的 http 地址")
	static2.NewFileServer(rootPath).Start(address)
}
