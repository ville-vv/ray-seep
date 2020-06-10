package control

import (
	"github.com/vilsongwei/vilgo/vlog"
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/msg"
	"sync"
)

type ConnectCenter struct {
	mu       sync.Mutex
	pods     map[int64]*Pod
	podHd    *PodHandler
	cNum     int
	cfg      *conf.Server
	register *RegisterCenter
	runner   *Runner
}

func (c *ConnectCenter) OnConnect(id int64, msgCtr *msg.MessageCenter) (err error) {
	pod := NewPod(id, msgCtr)
	msgCtr.SetRouter(pod.OnMessage)
	// 认证成功加入到管理服务中
	return c.addPod(id, pod)
}

func (c *ConnectCenter) OnDisConnect(id int64) {
	vlog.DEBUG("[%d] disconnect current number:%d", id, c.cNum)
}

func (c *ConnectCenter) addPod(id int64, pd *Pod) error {
	// 认证成功加入到管理服务中
	c.mu.Lock()
	defer c.mu.Unlock()
	c.pods[id] = pd
	c.cNum += 1
	vlog.DEBUG("current number [%d]", c.cNum)
	return nil
}

func (c *ConnectCenter) Register() {
}
