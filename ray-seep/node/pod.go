// @File     : pod
// @Author   : Ville
// @Time     : 19-9-26 下午4:40
// node
package node

import (
	"ray-seep/ray-seep/common/pkg"
	"ray-seep/ray-seep/mng"
)

// Pod 是一个 代理服务 的管理器连接，包括代理和控制连接
//type Pod interface {
//}

type Pod struct {
	domain string
	mng.Sender
	proxy Proxy
	id    int64
}

func NewPod(id int64, sender mng.Sender) *Pod {
	return &Pod{id: id, Sender: sender}
}

func (p *Pod) Id() int64 {
	return p.id
}

func (p *Pod) RegisterProxy() {
}

func (p *Pod) Operate(cmd pkg.Command, body []byte) ([]byte, error) {
	return nil, nil
}

func (p *Pod) CreateHost(req *pkg.CreateHostReq) {
}
