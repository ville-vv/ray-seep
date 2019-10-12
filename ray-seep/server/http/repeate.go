// @File     : repeate
// @Author   : Ville
// @Time     : 19-10-12 下午3:22
// http
package http

import (
	"net"
)

// Repeater 是一个中继器，用于转发 conn 的数据
type Repeater interface {
	// 转发
	Transmit(host string, c net.Conn)
}

type repaterHttp struct {
}

func (sel *repaterHttp) Transmit(host string, c net.Conn) {
}
