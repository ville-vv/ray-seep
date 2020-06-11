package node

import (
	"github.com/vilsongwei/vilgo/vlog"
	"ray-seep/ray-seep/common/conn"
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/msg"
	"ray-seep/ray-seep/server_v2/proxy"
	"sync"
)

type ConnectCenter struct {
	mu       sync.Mutex
	pods     map[int64]*Pod
	podHd    *PodHandler
	cNum     int
	cfg      *conf.Server
	register *proxy.RegisterCenter
	runner   *Runner
}

func (c *ConnectCenter) OnConnect(cancel chan interface{}, cn conn.Conn) error {
	id := cn.Id()
	msgCtr := msg.NewMessageCenter(cn)
	pod := NewPod(id, msgCtr)
	msgCtr.SetRouter(pod.OnMessage)
	if err := c.addPod(id, pod); err != nil {
		return err
	}
	msgCtr.Run()
	cn.Close()
	// 认证成功加入到管理服务中
	return nil
}

func (c *ConnectCenter) OnDisConnect(id int64) {
	// 认证成功加入到管理服务中
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.pods[id]; ok {
		delete(c.pods, id)
		c.cNum--
	}
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
