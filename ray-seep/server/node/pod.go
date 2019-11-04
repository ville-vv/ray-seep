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
	id int64
}

func NewPod(id int64, sender mng.Sender) *Pod {
	p := &Pod{id: id, Sender: sender}

	return p
}

func (p *Pod) Id() int64 {
	return p.id
}

func (p *Pod) OnMessage(cmd pkg.Command, body []byte) ([]byte, error) {
	switch cmd {
	case pkg.CmdRegisterProxyReq:
	case pkg.CmdCreateHostReq:
	case pkg.CmdError:
	}
	return nil, nil
}

func (p *Pod) PushMsg(msgPkg *pkg.Package) (err error) {
	return p.Sender.SendMsg(msgPkg)
}
