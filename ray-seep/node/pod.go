// @File     : pod
// @Author   : Ville
// @Time     : 19-9-26 下午4:40
// node
package node

import "ray-seep/ray-seep/conn"

// Pod 是一个 代理服务 的管理器连接，包括代理和控制连接
type Pod interface {
}

type ConnPod struct {
	clint conn.Conn
	proxy Proxy
	id    int64
}

func (p *ConnPod) Id() int64 {
	return p.id
}
