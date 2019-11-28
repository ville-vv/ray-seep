// @File     : pod
// @Author   : Ville
// @Time     : 19-9-26 下午4:40
// node
package control

import (
	jsoniter "github.com/json-iterator/go"
	"ray-seep/ray-seep/common/errs"
	"ray-seep/ray-seep/common/util"
	"ray-seep/ray-seep/proto"
	"ray-seep/ray-seep/server/online"
	"vilgo/vlog"
)

type PodRouterFun func([]byte) ([]byte, error)

// Pod 是一个 代理服务 的管理器连接，包括代理和控制连接
type Pod struct {
	connId  int64
	userId  int64
	name    string
	secret  string
	domain  string
	out     chan proto.Package
	route   map[int32]PodRouterFun
	userMng *online.UserManager
	podHd   *PodHandler
}

func NewPod(id int64, domain string, podHd *PodHandler, out chan proto.Package) *Pod {
	p := &Pod{connId: id, domain: domain, podHd: podHd, out: out}
	p.initRoute()
	return p
}

func (p *Pod) initRoute() {
	p.route = make(map[int32]PodRouterFun)
	p.route[proto.CmdLoginReq] = p.LoginReq
	p.route[proto.CmdCreateHostReq] = p.CreateHostReq
	p.route[proto.CmdRunProxyRsp] = p.RunProxyReq
	p.route[proto.CmdLogoutReq] = p.LogoutReq
}
func (p *Pod) ConnId() int64 {
	return p.connId
}
func (p *Pod) PushMsg(msgPkg *proto.Package) (err error) {
	p.out <- *msgPkg
	return
}

func (p *Pod) OnMessage(cmd int32, body []byte) ([]byte, error) {
	vlog.INFO("[%d] message cmd:[%d] body:%s", p.connId, cmd, string(body))
	if rt, ok := p.route[cmd]; ok {
		return rt(body)
	}
	return nil, errs.ErrNoCmdRouterNot
}

func (p *Pod) LoginReq(req []byte) (rsp []byte, err error) {
	reqLogin := proto.LoginReq{}
	if err = jsoniter.Unmarshal(req, &reqLogin); err != nil {
		vlog.ERROR("[%d] login Unmarshal error:%s", p.connId, err.Error())
		return
	}
	token := util.RandToken()
	vlog.INFO("login request userId=%d,appId=%s,token=%s", reqLogin.UserId, reqLogin.AppId, token)
	secret, err := p.podHd.OnLogin(p.connId, reqLogin.UserId, reqLogin.AppId, token)
	if err != nil {
		vlog.ERROR("[%d] login store token error:%s", p.connId, err.Error())
		return
	}
	rsp, err = jsoniter.Marshal(&proto.LoginRsp{
		Id:    p.connId,
		Token: token,
	})
	if err != nil {
		vlog.ERROR("[%d] login Marshal error:%s", err.Error())
	}
	p.name = reqLogin.Name
	p.secret = secret
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
		vlog.ERROR("[%d] notice run proxy error %s", p.connId, err.Error())
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

func (p *Pod) LogoutReq(req []byte) (rsp []byte, err error) {
	return nil, nil
}
