// @File     : adopter
// @Author   : Ville
// @Time     : 19-9-27 下午3:27
// node
package node

import (
	"errors"
	"ray-seep/ray-seep/common/util"
	"ray-seep/ray-seep/proto"
	"sync"
	"vilgo/vlog"
)

// 身份认证
type Author interface {
	Identity(name string, password string, token string) error
}

type MessageAdopter struct {
	mu     sync.Mutex
	pods   map[int64]*Pod
	cNum   int
	author Author
}

func NewMessageAdopter() *MessageAdopter {
	return &MessageAdopter{
		pods: make(map[int64]*Pod),
	}
}

func (sel *MessageAdopter) identify(p proto.Package) error {
	if p.Cmd != proto.CmdIdentifyReq {
		return errors.New("identify authentication fail")
	}
	if sel.author != nil {
		return sel.author.Identity("", "", "")
	}
	return nil
}

// OnConnect 有用户连接上来会出发这个事件
func (sel *MessageAdopter) OnConnect(id int64, tr proto.MsgTransfer) (err error) {
	// 建立连接的首要任务就是获取认证信息，如果认证失败就直接断开连接
	var authMsg proto.Package
	if err = tr.RecvMsg(&authMsg); err != nil {
		vlog.ERROR("get auth message error %s", err.Error())
		return err
	}
	if err = sel.identify(authMsg); err != nil {
		vlog.ERROR("identify check error %s", err.Error())
		return
	}

	authMsgRsp := proto.NewWithObj(proto.CmdIdentifyRsp, proto.IdentifyRsp{Id: id, Token: util.RandToken()})
	if err = tr.SendMsg(authMsgRsp); err != nil {
		vlog.ERROR("response auth message error %s", err.Error())
		return
	}
	// 认证成功加入到管理服务中
	sel.mu.Lock()
	sel.pods[id] = NewPod(id, tr)
	sel.cNum++
	sel.mu.Unlock()
	return nil
}

// OnDisConnect 有用户断开连接的时候会触发这个事件
func (sel *MessageAdopter) OnDisConnect(id int64) {
	// 认证成功加入到管理服务中
	sel.mu.Lock()
	for k := range sel.pods {
		delete(sel.pods, k)
		sel.cNum--
	}
	vlog.DEBUG("Pod disconnect current number[%d]:%d", sel.cNum, id)
	sel.mu.Unlock()
}

// OnMessage 客户端发送消息过来的时候会触发该事件
func (sel *MessageAdopter) OnMessage(id int64, req *proto.Package) (rsp proto.Package, err error) {

	vlog.DEBUG("Pod %d msg [cmd:%v][body:%s]", id, req.Cmd, string(req.Body))
	// 心跳直接返回
	if req.Cmd == proto.CmdPing {
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

// PushMsg 主动消息推送
func (sel *MessageAdopter) PushMsg(id int64, p *proto.Package) error {
	return sel.pushMsg(id, p)
}

func (sel *MessageAdopter) pushMsg(id int64, p *proto.Package) error {
	sel.mu.Lock()
	pod, ok := sel.pods[id]
	sel.mu.Unlock()
	if !ok {
		return nil
	}
	return pod.PushMsg(p)
}
