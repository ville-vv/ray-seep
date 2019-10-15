package vlog

import (
	"net"
	"os"
	"path"
	"strings"
)

// 获取本机正在使用的IP
func GetServerIP() (ip string, err error) {
	conn, err := net.Dial("udp", "google.com:80")
	if err != nil {
		return
	}
	defer conn.Close()
	ip = strings.Split(conn.LocalAddr().String(), ":")[0]

	return
}

// 获取本机名字
func GetHostName() (hostname string, err error) {
	hostname, err = os.Hostname()
	return
}

// 检查文件路径是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// 穿件文件路径，支持多个文件路径的创建
func CreateLogPath(a ...string) {
	for _, v := range a {
		p := path.Dir(v)
		if p == "" {
			continue
		}
		if PathExists(p) {
			continue
		}
		if err := os.MkdirAll(p, os.ModePerm); err != nil {
			panic(err)
		}
	}
}
