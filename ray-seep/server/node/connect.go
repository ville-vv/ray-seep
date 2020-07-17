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
	"sync/atomic"
)

// 使用 sync.Map 代替 map+Mutex 因为pods读数据比较多，使用sync.Map更合适
// pods 在客户端的代理链接断开后，如果有用户的访问请求，查询到register中没有client proxy 通道，
// 这个时候 register 会调用 GetNotice 方法获取pod 发送数据，所以pods 会存在大量的 read操作
type ConnectCenter struct {
	mu sync.Mutex
	//pods  map[int64]*Pod
	pods  sync.Map
	podHd ifc.PodHandler
	cNum  int32
	cfg   *conf.Server
	hsr   hostsrv.HostServer
	exit  ifc.ExitDevice
}

func NewConnectCenter(cfg *conf.Server, runner hostsrv.HostServer, podHd ifc.PodHandler, exit ifc.ExitDevice) *ConnectCenter {
	return &ConnectCenter{
		mu: sync.Mutex{},
		//pods:  make(map[int64]*Pod),
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
	if val, ok := c.pods.Load(id); ok {
		pd := val.(*Pod)
		name = pd.HttpAddr()
		_ = pd.LogoutReq(nil)
		atomic.AddInt32(&c.cNum, -1)
		c.pods.Delete(id)
	}
	c.exit.Logout(name, id)
	vlog.DEBUG("[%d] disconnect current number:%d", id, c.cNum)
}

func (c *ConnectCenter) addPod(id int64, pd *Pod) error {
	atomic.AddInt32(&c.cNum, +1)
	c.pods.Store(id, pd)
	vlog.DEBUG("current client connect number [%d]", c.cNum)
	return nil
}

func (c *ConnectCenter) GetNotice(id int64) (ifc.MessageNotice, error) {
	val, ok := c.pods.Load(id)
	if !ok {
		return nil, errs.ErrClientControlNotExist
	}
	pod := val.(*Pod)
	tmp := *pod
	return &tmp, nil
}
