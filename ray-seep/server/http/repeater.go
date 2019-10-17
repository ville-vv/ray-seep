// @File     : repeate
// @Author   : Ville
// @Time     : 19-10-12 下午3:22
// http
package http

import (
	"net"
	"vilgo/vlog"
)

// Repeater 是一个中继器，用于转发 conn 的数据
type Repeater interface {
	// 转发
	Transmit(host string, c net.Conn)
}

// ProxyGainer 代理连接获取器
type ProxyGainer interface {
	GetProxy(identify string) (net.Conn, error)
}

// repeaterHttp HTTP 请求的使用的中续器
type RepeaterHttp struct {
	regCenter ProxyGainer // 注册中心
}

func NewRepeaterHttp(regCenter ProxyGainer) *RepeaterHttp {
	return &RepeaterHttp{regCenter: regCenter}
}

func (sel *RepeaterHttp) Transmit(host string, c net.Conn) {
	pxyConn, err := sel.regCenter.GetProxy(host)
	if err != nil {
		vlog.ERROR("%v", err)
		return
	}
	sel.copy(c, pxyConn)
}

func (sel *RepeaterHttp) copy(dst net.Conn, src net.Conn) {

}
