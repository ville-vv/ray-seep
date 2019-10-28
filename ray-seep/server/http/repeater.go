// @File     : repeate
// @Author   : Ville
// @Time     : 19-10-12 下午3:22
// http
package http

import (
	"io"
	"net"
	"sync"
	"vilgo/vlog"
)

// Repeater 是一个中继器，用于转发 conn 的数据
type Repeater interface {
	// 转发
	Transfer(host string, c net.Conn)
}

// ProxyGainer 代理连接获取器
type ProxyGainer interface {
	GetProxy(identify string) (net.Conn, error)
}

// NetRepeater 网络请求的使用的中续器
type NetRepeater struct {
	regCenter ProxyGainer // 注册中心
}

func NewNetRepeater(regCenter ProxyGainer) *NetRepeater {
	return &NetRepeater{regCenter: regCenter}
}

func (sel *NetRepeater) Transfer(host string, c net.Conn) {
	pxyConn, err := sel.regCenter.GetProxy(host)
	if err != nil {
		vlog.ERROR("%v", err)
		return
	}
	sel.relay(c, pxyConn)
}

// exchange 请求数据转播
// @dst 是目标请求网络连接
// @src 是发起请求的连接
// @return 1 : 请求者发送的数据长度
// @return 2 : 被请求者返回的数据长度
// @return err : 错误
func (sel *NetRepeater) relay(dst net.Conn, src net.Conn) (int64, int64, error) {
	defer dst.Close()
	defer src.Close()
	type res struct {
		N   int64
		Err error
	}
	ch := make(chan res)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		wg.Done()
		// 先启动
		n, err := io.Copy(dst, src)
		ch <- res{N: n, Err: err}

	}()
	wg.Wait()
	n, err := io.Copy(src, dst)
	rs := <-ch
	if err == nil {
		err = rs.Err
	}
	return n, rs.N, err
}

// pip
func (sel *NetRepeater) join(dst net.Conn, src net.Conn) (int64, int64, error) {
	var wait sync.WaitGroup
	var err error
	pipe := func(dst net.Conn, src net.Conn, bytesCopied *int64) {
		defer wait.Done()
		*bytesCopied, err = io.Copy(dst, src)
	}
	wait.Add(2)
	var fromBytes, toBytes int64
	go pipe(src, dst, &fromBytes)
	go pipe(dst, src, &toBytes)
	wait.Wait()
	return fromBytes, toBytes, err
}
