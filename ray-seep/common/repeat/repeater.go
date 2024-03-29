// @File     : repeate
// @Author   : Ville
// @Time     : 19-10-12 下午3:22
// http
package repeat

import (
	"io"
	"net"
	"time"
)

// ProxyGainer 代理连接获取器
type NetConnGainer interface {
	GetNetConn(name string) (net.Conn, error)
}

// NetRepeater 网络请求的使用的中续器
type NetRepeater struct {
	pxyGainer NetConnGainer // 注册中心
}

func NewNetRepeater(pxyGainer NetConnGainer) *NetRepeater {
	return &NetRepeater{pxyGainer: pxyGainer}
}

func (sel *NetRepeater) Transfer(host string, c net.Conn) (int64, int64, error) {
	pxyConn, err := sel.pxyGainer.GetNetConn(host)
	if err != nil {
		return 0, 0, err
	}
	defer pxyConn.Close()
	_ = pxyConn.SetDeadline(time.Time{})
	reqLength, respLength, err := Relay(pxyConn, c)
	if err != nil {
		if netErr, ok := err.(net.Error); !(ok && netErr.Timeout()) {
			if err != io.EOF {
				return 0, 0, err
			}
		}
	}
	return reqLength, respLength, nil
}

type result struct {
	N   int64
	Err error
}

func relay(dst net.Conn, src net.Conn, resCh chan result) {
	n, err := io.Copy(dst, src)
	_ = src.SetDeadline(time.Now())
	_ = dst.SetDeadline(time.Now())
	resCh <- result{N: n, Err: err}
}

// Relay 请求数据转播
// @dst 是目标请求网络连接
// @src 是发起请求的连接
// @return 1 : 请求者发送的数据长度
// @return 2 : 被请求者返回的数据长度
// @return err : 错误
func Relay(dst net.Conn, src net.Conn) (int64, int64, error) {
	rspCh := make(chan result)
	// 返回数据转发
	//go relay(src, dst, rspCh)

	go func() {
		n, err := io.Copy(src, dst)
		_ = src.SetDeadline(time.Now())
		_ = dst.SetDeadline(time.Now())
		rspCh <- result{N: n, Err: err}
	}()

	n, err := io.Copy(dst, src)
	_ = src.SetDeadline(time.Now())
	_ = dst.SetDeadline(time.Now())
	rs := <-rspCh
	if err == nil {
		err = rs.Err
	}
	return n, rs.N, err
}
