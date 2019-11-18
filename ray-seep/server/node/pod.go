// @File     : pod
// @Author   : Ville
// @Time     : 19-9-26 下午4:40
// node
package node

import (
	jsoniter "github.com/json-iterator/go"
	"ray-seep/ray-seep/common/errs"
	"ray-seep/ray-seep/common/util"
	"ray-seep/ray-seep/proto"
	"vilgo/vlog"
)

type PodRouterFun func([]byte) ([]byte, error)

// Pod 是一个 代理服务 的管理器连接，包括代理和控制连接
type Pod struct {
	domain string
	sender proto.Sender
	id     int64
	route  map[int32]PodRouterFun
}

func NewPod(id int64, sender proto.Sender, domain string) *Pod {
	p := &Pod{id: id, sender: sender, domain: domain}
	p.initRoute()
	return p
}

func (p *Pod) initRoute() {
	p.route = make(map[int32]PodRouterFun)
	p.route[proto.CmdLoginReq] = p.Login
	p.route[proto.CmdCreateHostReq] = p.CreateHostReq
	p.route[proto.CmdRunProxyRsp] = p.RunProxyReq
}
func (p *Pod) Id() int64 {
	return p.id
}
func (p *Pod) PushMsg(msgPkg *proto.Package) (err error) {
	return p.sender.SendMsg(msgPkg)
}

func (p *Pod) OnMessage(cmd int32, body []byte) ([]byte, error) {
	vlog.INFO("[%d] message cmd:[%d] body:%s", p.id, cmd, string(body))
	if rt, ok := p.route[cmd]; ok {
		return rt(body)
	}
	return nil, errs.ErrNoCmdRouterNot
}

func (p *Pod) Login(req []byte) (rsp []byte, err error) {
	rsp, err = jsoniter.Marshal(&proto.LoginRsp{
		Id:    p.id,
		Token: util.RandToken(),
	})
	if err != nil {
		vlog.ERROR("[%d] login error:%s", err.Error())
	}
	return
}

func (p *Pod) CreateHostReq(req []byte) (rsp []byte, err error) {
	reqObj := &proto.CreateHostReq{}
	if err = jsoniter.Unmarshal(req, reqObj); err != nil {
		vlog.ERROR("")
		return
	}
	rspObj := proto.CreateHostRsp{
		Domain: reqObj.SubDomain + "." + p.domain,
	}
	return jsoniter.Marshal(rspObj)
}

// NoticeRunProxy 通知用户启动代理服务服务
func (p *Pod) NoticeRunProxy() {
	notice := &proto.Package{
		Cmd:  proto.CmdNoticeRunProxy,
		Body: nil,
	}

	if err := p.PushMsg(notice); err != nil {
		vlog.ERROR("[%d] notice run proxy error %s", p.id, err.Error())
	}
}

func (p *Pod) RunProxyReq(req []byte) (rsp []byte, err error) {
	reqObj := &proto.RunProxyReq{}
	if err = jsoniter.Unmarshal(req, reqObj); err != nil {
		vlog.ERROR("run proxy request error %s", err.Error())
		return
	}
	return jsoniter.Marshal(proto.RunProxyRsp{})
}
