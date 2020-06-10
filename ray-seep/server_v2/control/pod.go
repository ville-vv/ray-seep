// @File     : pod
// @Author   : Ville
// @Time     : 19-9-26 下午4:40
// node
package control

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/vilsongwei/vilgo/vlog"
	"ray-seep/ray-seep/common/util"
	"ray-seep/ray-seep/msg"
	"ray-seep/ray-seep/proto"
	"ray-seep/ray-seep/server/http"
)

type PodRouterFun func([]byte) (interface{}, error)

// Pod 是一个 代理服务 的管理器连接，包括代理和控制连接
type Pod struct {
	connId int64 // 连接ID号
	msgCtr *msg.MessageCenter
}

func NewPod(id int64, msgCtr *msg.MessageCenter) *Pod {
	p := &Pod{
		connId: id,
		msgCtr: msgCtr,
	}
	return p
}

func (p *Pod) ConnId() int64 {
	return p.connId
}

// PushMsg 用来主动推送消息
func (p *Pod) PushMsg(msgPkg *msg.Package) (err error) {
	return p.msgCtr.Send(msgPkg)
}

// OnMessage 有消息接收就会发送到这里来
func (p *Pod) OnMessage(req *msg.Request, sender msg.ResponseSender) error {
	return nil
}

func (p *Pod) LoginReq(req []byte) (interface{}, error) {
	resp := &proto.LoginRsp{Id: 0, Token: util.RandToken()}
	return resp, nil
}

func (p *Pod) CreateHostReq(req []byte) (rsp interface{}, err error) {
	reqObj := &proto.CreateHostReq{}
	if err = jsoniter.Unmarshal(req, reqObj); err != nil {
		vlog.ERROR("create host request message json unmarshal fail", err)
		return
	}
	// 创建主机需要检验是否已经登录了
	if err = p.podHd.OnCreateHost(p.connId, p.name, reqObj.Token); err != nil {
		vlog.ERROR("on create host fail", err)
		return
	}

	join := JoinItem{
		Name:   p.httpAddr,
		ConnId: p.connId,
		Run:    http.NewServerWithAddr(":"+p.httpPort, p.gainer),
		Err:    make(chan error),
	}

	p.runner.Join() <- join
	if err = <-join.Err; err != nil {
		vlog.ERROR("[%d] http join error %s", p.connId, err.Error())
		return
	}

	rspObj := proto.CreateHostRsp{}
	return rspObj, nil
}

// NoticeRunProxy 通知用户启动代理服务服务
func (p *Pod) NoticeRunProxy() {
}

func (p *Pod) RunProxyReq(req []byte) (rsp interface{}, err error) {
}

func (p *Pod) LogoutReq(req []byte) (rsp interface{}, err error) {
}
