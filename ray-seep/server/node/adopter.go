// @File     : adopter
// @Author   : Ville
// @Time     : 19-9-27 下午3:27
// node
package node

import (
	"fmt"
	"net"
	"ray-seep/ray-seep/common/conn"
	"ray-seep/ray-seep/common/errs"
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/proto"
	"ray-seep/ray-seep/server/online"
	"strings"
	"sync"
	"vilgo/vlog"
)

type MessageAdopter struct {
	mu       sync.Mutex
	pods     map[int64]*Pod
	userMng  *online.UserManager
	cNum     int
	cfg      *conf.Server
	register *RegisterCenter
}

func NewMessageAdopter(cfg *conf.Server, uMng *online.UserManager) *MessageAdopter {
	m := &MessageAdopter{
		pods:    make(map[int64]*Pod),
		cfg:     cfg,
		userMng: uMng,
	}
	m.register = NewRegisterCenter(cfg.Ctl.UserMaxProxyNum, m)

	return m
}

func (sel *MessageAdopter) Domain() string {
	domain := sel.cfg.Http.Domain
	if sel.cfg.Http.Port != 80 {
		domain = fmt.Sprintf("%s:%d", domain, sel.cfg.Http.Port)
	}
	return domain
}

// OnConnect 有用户连接上来会出发这个事件
func (sel *MessageAdopter) OnConnect(id int64, tr proto.MsgTransfer) (err error) {
	pd := NewPod(id, tr, sel.Domain(), sel.userMng)
	// 建立连接的首要任务就是获取认证信息，如果认证失败就直接断开连接
	var req proto.Package
	if err = tr.RecvMsg(&req); err != nil {
		vlog.ERROR("[%d] get auth message error %s", id, err.Error())
		return err
	}
	rsp := &proto.Package{
		Cmd: proto.CmdLoginRsp,
	}
	rsp.Body, err = pd.LoginReq(req.Body)
	if err != nil {
		vlog.ERROR("[%d] on connect deal message error:%d", id, err.Error())
		return
	}
	if err = tr.SendMsg(rsp); err != nil {
		vlog.ERROR("[%d] response auth message error %s", id, err.Error())
		return
	}
	// 认证成功加入到管理服务中
	return sel.addPod(id, pd)
}

func (sel *MessageAdopter) addPod(id int64, pd *Pod) error {
	// 认证成功加入到管理服务中
	sel.mu.Lock()
	defer sel.mu.Unlock()
	sel.pods[id] = pd
	sel.cNum += 1
	vlog.DEBUG("[%d] pod disconnect current number[%d]", id, sel.cNum)
	return nil
}

// OnDisConnect 有用户断开连接的时候会触发这个事件
func (sel *MessageAdopter) OnDisConnect(id int64) {
	// 认证成功加入到管理服务中
	sel.mu.Lock()
	defer sel.mu.Unlock()
	if pd, ok := sel.pods[id]; ok {
		sel.register.LogOff(pd.name, id)
		_, _ = pd.LogoutReq([]byte{})
		delete(sel.pods, id)
		sel.cNum--
	}
	vlog.DEBUG("[%d] disconnect current number:%d", id, sel.cNum)
}

// OnMessage 客户端发送消息过来的时候会触发该事件
func (sel *MessageAdopter) OnMessage(id int64, req *proto.Package) (rsp proto.Package, err error) {
	// 心跳直接返回
	if req.Cmd == proto.CmdPing {
		//vlog.DEBUG("[%d] Ping", id)
		rsp.Cmd = proto.CmdPong
		return
	}

	sel.mu.Lock()
	pod, ok := sel.pods[id]
	sel.mu.Unlock()
	if !ok {
		return
	}

	rsp.Cmd = req.Cmd + 1
	if rsp.Body, err = pod.OnMessage(req.Cmd, req.Body); err != nil {
		vlog.ERROR("pod operate error %s", err.Error())
		return
	}

	return
}

func (sel *MessageAdopter) Register(name string, id int64, cc conn.Conn) (err error) {
	return sel.register.Register(name, id, cc)
}

func (sel *MessageAdopter) GetNetConn(name string) (net.Conn, error) {
	return sel.register.GetNetConn(strings.Split(name, ".")[0])
}

// PushMsg 主动消息推送
func (sel *MessageAdopter) PushMsg(id int64, p *proto.Package) error {
	return sel.pushMsg(id, p)
}

func (sel *MessageAdopter) pushMsg(id int64, p *proto.Package) error {
	sel.mu.Lock()
	pod, ok := sel.pods[id]
	sel.mu.Unlock()
	if !ok {
		return errs.ErrClientControlNotExist
	}
	return pod.PushMsg(p)
}
