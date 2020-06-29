package node

import (
	"github.com/vilsongwei/vilgo/vlog"
	"ray-seep/ray-seep/common/conn"
	"ray-seep/ray-seep/common/errs"
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/msg"
	"ray-seep/ray-seep/server/hostsrv"
	"ray-seep/ray-seep/server/ifc"
	"sync"
)

type ConnectCenter struct {
	mu    sync.Mutex
	pods  map[int64]*Pod
	podHd *PodHandler
	cNum  int
	cfg   *conf.Server
	hsr   hostsrv.HostServer
	exit  ifc.ExitDevice
}

func NewConnectCenter(cfg *conf.Server, runner hostsrv.HostServer, podHd *PodHandler, exit ifc.ExitDevice) *ConnectCenter {
	return &ConnectCenter{
		mu:    sync.Mutex{},
		pods:  make(map[int64]*Pod),
		podHd: podHd,
		cNum:  0,
		cfg:   cfg,
		hsr:   runner,
		exit:  exit,
	}
}

func (c *ConnectCenter) OnConnect(cancel chan interface{}, cn conn.Conn) error {
	id := cn.Id()
	msgCtr := msg.NewMessageCenter(cn)
	pod := NewPod(id, c.cfg, msgCtr, c.hsr, c.podHd)
	if err := c.addPod(id, pod); err != nil {
		return err
	}
	msgCtr.Run(pod.OnMessage)
	_ = cn.Close()
	return nil
}

func (c *ConnectCenter) OnDisConnect(id int64) {
	// 认证成功加入到管理服务中
	var name string
	c.mu.Lock()
	if pd, ok := c.pods[id]; ok {
		name = pd.HttpAddr()
		_ = pd.LogoutReq(nil)
		delete(c.pods, id)
		c.cNum--
	}
	c.mu.Unlock()
	c.exit.Logout(name, id)
	vlog.DEBUG("[%d] disconnect current number:%d", id, c.cNum)
}

func (c *ConnectCenter) addPod(id int64, pd *Pod) error {
	// 认证成功加入到管理服务中
	c.mu.Lock()
	defer c.mu.Unlock()
	c.pods[id] = pd
	c.cNum += 1
	vlog.DEBUG("current client connect number [%d]", c.cNum)
	return nil
}

func (c *ConnectCenter) GetNotice(id int64) (ifc.MessageNotice, error) {
	c.mu.Lock()
	pod, ok := c.pods[id]
	if !ok {
		c.mu.Unlock()
		return nil, errs.ErrClientControlNotExist
	}
	c.mu.Unlock()
	tmp := *pod
	return &tmp, nil
}
