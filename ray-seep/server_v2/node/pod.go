// @File     : pod
// @Author   : Ville
// @Time     : 19-9-26 下午4:40
// node
package node

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/vilsongwei/vilgo/vlog"
	"ray-seep/ray-seep/common/util"
	"ray-seep/ray-seep/msg"
	"ray-seep/ray-seep/proto"
	"ray-seep/ray-seep/server_v2/hostsrv"
)

type PodRouterFun func([]byte) (interface{}, error)

// Pod 是一个 代理服务 的管理器连接，包括代理和控制连接
type Pod struct {
	connId   int64 // 连接ID号
	sender   msg.ResponseSender
	hsr      hostsrv.HostServer
	appKey   string
	name     string
	secret   string
	httpPort string
	httpAddr string
}

func NewPod(id int64, sender msg.ResponseSender, hsr hostsrv.HostServer) *Pod {
	p := &Pod{
		connId: id,
		sender: sender,
		hsr:    hsr,
	}
	return p
}

func (p *Pod) ConnId() int64 {
	return p.connId
}

// PushMsg 用来主动推送消息
func (p *Pod) PushInJson(cmd int32, obj interface{}) (err error) {
	body, err := jsoniter.Marshal(obj)
	if err != nil {
		return err
	}
	p.sender.SendCh() <- msg.Package{Cmd: cmd, Body: body}
	return
}

func (p *Pod) PushInByte(cmd int32, data []byte) (err error) {
	p.sender.SendCh() <- msg.Package{Cmd: cmd, Body: data[:]}
	return
}

// OnMessage 有消息接收就会发送到这里来
func (p *Pod) OnMessage(req *msg.Request) {
	var err error
	switch req.Body.Cmd {
	case msg.CmdLoginReq:
		err = p.LoginReq(req.Body.Cmd, req.Body.Body)
	case msg.CmdCreateHostReq:
		err = p.CreateHostReq(req.Body.Cmd, req.Body.Body)
	}
	if err != nil {
		vlog.ERROR("on message error %s", err.Error())
		_ = p.PushInJson(msg.CmdError, &proto.ErrorNotice{ErrMsg: err.Error()})
	}
	return
}

func (p *Pod) LoginReq(cmd int32, body []byte) (err error) {
	reqLogin := proto.LoginReq{}
	resp := &proto.LoginRsp{Id: p.connId, Token: util.RandToken()}
	if err := jsoniter.Unmarshal(body, &reqLogin); err != nil {
		vlog.ERROR("[%d] login Unmarshal error:%s", p.connId, err.Error())
		return err
	}
	// 登录操作处理
	p.appKey = reqLogin.AppKey
	p.name = reqLogin.Name
	vlog.INFO("[%d] login token =%s name=%s", p.connId, p.appKey, reqLogin.Name)
	vlog.INFO("[%d] login request userId=%d ", p.connId, reqLogin.UserId)
	return p.PushInJson(msg.CmdLoginRsp, resp)
}

func (p *Pod) CreateHostReq(cmd int32, body []byte) (err error) {
	vlog.INFO("create host request message json ")
	reqObj := &proto.CreateHostReq{}
	if err = jsoniter.Unmarshal(body, reqObj); err != nil {
		vlog.ERROR("create host request message json unmarshal fail", err)
		return
	}
	// 创建主机需要检验是否已经登录了
	//if err = p.podHd.OnCreateHost(p.connId, p.name, reqObj.Token); err != nil {
	//	vlog.ERROR("on create host fail", err)
	//	return
	//}

	if err := p.hsr.Create(&hostsrv.Option{
		Id:     0,
		Kind:   "",
		Addr:   "",
		SendCh: p.sender.SendCh(),
	}); err != nil {
		return err
	}

	return p.PushInJson(cmd+1, proto.CreateHostRsp{})
}

// NoticeRunProxy 通知用户启动代理服务服务
func (p *Pod) NoticeRunProxy(data []byte) error {
	if err := p.PushInJson(proto.CmdNoticeRunProxy, nil); err != nil {
		vlog.ERROR("[%d] notice run proxy error %s", p.connId, err.Error())
	}
	return nil
}
func (p *Pod) NoticeRunProxyRsp(data []byte) error {
	return nil
}

func (p *Pod) LogoutReq(req []byte) (rsp interface{}, err error) {
	return nil, nil
}
