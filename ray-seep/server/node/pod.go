// @File     : pod
// @Author   : Ville
// @Time     : 19-9-26 下午4:40
// node
package node

import (
	"ray-seep/ray-seep/proto"
)

// Pod 是一个 代理服务 的管理器连接，包括代理和控制连接
type Pod struct {
	domain string
	sender proto.Sender
	id     int64
}

func NewPod(id int64, sender proto.Sender) *Pod {
	p := &Pod{id: id, sender: sender}
	return p
}

func (p *Pod) Id() int64 {
	return p.id
}

func (p *Pod) OnMessage(cmd int32, body []byte) ([]byte, error) {
	switch cmd {
	case proto.CmdRegisterProxyReq:
	case proto.CmdCreateHostReq:
	case proto.CmdError:
	}
	return nil, nil
}

func (p *Pod) PushMsg(msgPkg *proto.Package) (err error) {
	return p.sender.SendMsg(msgPkg)
}
