// @File     : repeate
// @Author   : Ville
// @Time     : 19-10-12 下午3:22
// http
package http

import (
	"net"
	"ray-seep/ray-seep/server/proxy"
	"vilgo/vlog"
)

// Repeater 是一个中继器，用于转发 conn 的数据
type Repeater interface {
	// 转发
	Transmit(host string, c net.Conn)
}

// repeaterHttp HTTP 请求的使用的中续器
type repeaterHttp struct {
	regCenter *proxy.RegisterCenter // 注册中心
}

func (sel *repeaterHttp) Transmit(host string, c net.Conn) {
	pxyConn, err := sel.regCenter.GetProxy(host)
	if err != nil {
		vlog.ERROR("%v", err)
		return
	}
	sel.Copy(c, pxyConn)
}

func (sel *repeaterHttp) Copy(dst net.Conn, src net.Conn) {

}
