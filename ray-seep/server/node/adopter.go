// @File     : adopter
// @Author   : Ville
// @Time     : 19-9-27 下午3:27
// node
package node

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"ray-seep/ray-seep/common/pkg"
	"ray-seep/ray-seep/common/util"
	"ray-seep/ray-seep/mng"
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
	author Author
}

func NewMessageAdopter() *MessageAdopter {
	return &MessageAdopter{
		pods: make(map[int64]*Pod),
	}
}

func (sel *MessageAdopter) identify(p pkg.Package) error {
	if p.Cmd != pkg.CmdIdentifyReq {
		return errors.New("identify authentication fail")
	}
	if sel.author != nil {
		return sel.author.Identity("", "", "")
	}
	return nil
}

func (sel *MessageAdopter) OnConnect(id int64, tr mng.MsgTransfer) (err error) {
	// 建立连接的首要任务就是获取认证信息，如果认证失败就直接断开连接
	var authMsg pkg.Package
	if err = tr.RecvMsg(&authMsg); err != nil {
		return err
	}
	if err = sel.identify(authMsg); err != nil {
		return
	}
	authMsg.Cmd = pkg.CmdIdentifyRsp
	authMsg.Body, err = jsoniter.Marshal(pkg.IdentifyRsp{Id: id, Token: util.RandToken()})
	if err != nil {
		return
	}
	if err = tr.SendMsg(&authMsg); err != nil {
		vlog.ERROR("response auth message error %s", err.Error())
		return
	}

	// 认证成功加入到管理服务中
	sel.mu.Lock()
	sel.pods[id] = NewPod(id, tr)
	sel.mu.Unlock()
	return nil
}

func (sel *MessageAdopter) OnDisConnect(id int64) {
	vlog.DEBUG("Pod disconnect :%d", id)
}

// OnHandler 这里传入 sender 是因为不用每次都
func (sel *MessageAdopter) OnMessage(id int64, req *pkg.Package) (rsp pkg.Package, err error) {

	vlog.DEBUG("Pod %d msg [cmd:%v][body:%s]", id, req.Cmd, string(req.Body))
	// 心跳直接返回
	if req.Cmd == pkg.CmdPing {
		rsp.Cmd = pkg.CmdPong
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

func (sel *MessageAdopter) PushMsg(id int64, p *pkg.Package) error {
	sel.mu.Lock()
	pod, ok := sel.pods[id]
	sel.mu.Unlock()
	if !ok {
		return nil
	}
	return pod.PushMsg(p)
}
