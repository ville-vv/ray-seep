// @File     : pod
// @Author   : Ville
// @Time     : 19-9-26 下午4:40
// node
package node

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/vilsongwei/vilgo/vlog"
	"ray-seep/ray-seep/common/util"
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/msg"
	"ray-seep/ray-seep/proto"
	"ray-seep/ray-seep/server_v2/hostsrv"
)

type PodRouterFun func([]byte) (interface{}, error)

// Pod 是一个 代理服务 的管理器连接，包括代理和控制连接
type Pod struct {
	srvCfg   *conf.Server
	sender   msg.ResponseSender
	hsr      hostsrv.HostServer
	podHd    *PodHandler
	connId   int64 // 连接ID号
	appKey   string
	name     string
	secret   string
	httpPort string
	httpAddr string
}

func NewPod(id int64, srv *conf.Server, sender msg.ResponseSender, hsr hostsrv.HostServer, podHd *PodHandler) *Pod {
	p := &Pod{
		connId: id,
		sender: sender,
		hsr:    hsr,
		podHd:  podHd,
		srvCfg: srv,
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
	vlog.INFO("发送消息%d：%s", cmd, string(body))
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
	connId := p.connId
	reqLogin := proto.LoginReq{}
	resp := &proto.LoginRsp{Id: connId, Token: util.RandToken()}
	if err := jsoniter.Unmarshal(body, &reqLogin); err != nil {
		vlog.ERROR("[%d] login Unmarshal error:%s", connId, err.Error())
		return err
	}

	vlog.INFO("[%d] login request userId=%d ", p.connId, reqLogin.UserId)
	ul, err := p.podHd.OnLogin(p.connId, reqLogin.UserId, reqLogin.Name, reqLogin.AppKey, resp.Token)
	if err != nil {
		vlog.ERROR("[%d] login store token error:%s", p.connId, err.Error())
		return err
	}
	// 登录操作处理
	p.appKey = reqLogin.AppKey
	p.name = reqLogin.Name
	p.secret = ul.Secret
	p.httpPort = ul.HttpPort
	p.httpAddr = fmt.Sprintf("%s:%s", p.srvCfg.Domain, ul.HttpPort)

	vlog.INFO("[%d] login token =%s name=%s addr:%s", connId, p.appKey, reqLogin.Name, p.httpAddr)
	return p.PushInJson(msg.CmdLoginRsp, resp)
}

func (p *Pod) CreateHostReq(cmd int32, body []byte) (err error) {
	vlog.INFO("create host request message: %s", string(body))
	reqObj := &proto.CreateHostReq{}
	if err = jsoniter.Unmarshal(body, reqObj); err != nil {
		vlog.ERROR("create host request message json unmarshal fail", err)
		return
	}
	// 创建主机需要检验是否已经登录了
	//if err = p.podHd.OnCreateHost(p.ConnId, p.name, reqObj.Token); err != nil {
	//	vlog.ERROR("on create host fail", err)
	//	return
	//}

	if err := p.hsr.Create(p.connId, "http", fmt.Sprintf(":%s", p.httpPort)); err != nil {
		return err
	}

	rspObj := proto.CreateHostRsp{
		ProxyPort:  p.srvCfg.Pxy.Port,
		HttpDomain: p.httpAddr,
	}

	return p.PushInJson(cmd+1, rspObj)
}

// NoticeRunProxy 通知用户启动代理服务服务
func (p *Pod) NoticeRunProxy(data []byte) error {
	if err := p.PushInByte(proto.CmdNoticeRunProxy, data); err != nil {
		vlog.ERROR("[%d] notice run proxy error %s", p.connId, err.Error())
	}
	return nil
}

func (p *Pod) NoticeRunProxyRsp(data []byte) error {
	// 通知客户端，启动一个代理链接的回应
	if err := p.PushInByte(proto.CmdNoticeRunProxy, data); err != nil {
		vlog.ERROR("[%d] notice run proxy error %s", p.connId, err.Error())
	}
	return nil
}

func (p *Pod) LogoutReq(req []byte) (err error) {
	vlog.INFO("停止服务")
	p.hsr.Destroy(p.connId, fmt.Sprintf(":%s", p.httpPort))
	return nil
}
