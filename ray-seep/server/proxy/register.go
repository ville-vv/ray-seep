// @File     : register
// @Author   : Ville
// @Time     : 19-10-12 下午3:32
// proxy
package proxy

import (
	"net"
	"ray-seep/ray-seep/common/conn"
	"ray-seep/ray-seep/common/errs"
	"ray-seep/ray-seep/proto"
	"ray-seep/ray-seep/server/online"
	"strings"
	"sync"
	"time"
	"vilgo/vlog"
)

type MessagePusher interface {
	PushMsg(id int64, p *proto.Package) error
}

// RegisterCenter 注册中心，记录用户启动的本地服务id与用户使用的域名映射
// 记录用户启动的服务的代理池
type RegisterCenter struct {
	lock     sync.RWMutex
	userMng  *online.UserManager
	pxyPools map[string]conn.Pool // 记录用户本地服务的代理 tcp 链接，使用 cid 获取链接
	pushMsg  MessagePusher
	caches   int
}

func NewRegisterCenter(caches int, ph MessagePusher, userMng *online.UserManager) *RegisterCenter {
	return &RegisterCenter{
		pxyPools: make(map[string]conn.Pool),
		pushMsg:  ph,
		caches:   caches, // 一个节点需能缓存的数量
		userMng:  userMng,
	}
}

// 注册用户链接
func (sel *RegisterCenter) Register(name string, id int64, cc conn.Conn) error {
	// 把tcp连接放到代理池中
	if err := sel.addProxy(name, id, cc); err != nil {
		vlog.ERROR("[%d]register proxy error %s", id, err.Error())
		return err
	}
	return sel.pushMsg.PushMsg(id, &proto.Package{Cmd: proto.CmdRunProxyRsp, Body: []byte("{}")})
}

func (sel *RegisterCenter) addProxy(name string, id int64, cc conn.Conn) error {
	sel.lock.Lock()
	defer sel.lock.Unlock()
	if p, ok := sel.pxyPools[name]; ok {
		return p.Push(id, cc)
	}
	pl := conn.NewPool(sel.caches)
	if err := pl.Push(id, cc); err != nil {
		return err
	}
	sel.pxyPools[name] = pl
	vlog.INFO("当前代理数量%d", pl.Size())
	return nil
}

func (sel *RegisterCenter) delProxy(name string, cid int64) {
	if pl, ok := sel.pxyPools[name]; ok {
		pl.Drop(cid)
		delete(sel.pxyPools, name)
	}
}

// GetProxy 获取代理tcp连接
func (sel *RegisterCenter) GetProxy(name string) (net.Conn, error) {
	name = strings.Split(name, ".")[0]
	sel.lock.RLock()
	pl, ok := sel.pxyPools[name]
	sel.lock.RUnlock()
	if ok {
		if cn, err := pl.Get(0); err == nil {
			vlog.INFO("获得代理连接：%d", cn.Id())
			return &registerConn{cn}, nil
		}
	}
	if pl == nil {
		return nil, errs.ErrProxySrvNotExist
	}

	id := sel.userMng.GetId(name)
	notice := &proto.Package{Cmd: proto.CmdNoticeRunProxy, Body: []byte("{}")}
	if err := sel.pushMsg.PushMsg(id, notice); err != nil {
		vlog.ERROR("[%d]push notice run proxy error %s", id, err.Error())
		return nil, errs.ErrProxySrvNotExist
	}
	// 如果没有取到就发送重置消息，请求连接一个代理
	tm := time.NewTicker(time.Second * 5)
	select {
	case cn, ok := <-pl.WaitGet():
		if !ok {
			return nil, errs.ErrProxySrvNotExist
		}
		return cn, nil
	case <-tm.C:
		vlog.WARN("wait get proxy timeout")
	}
	return nil, errs.ErrProxySrvNotExist
}

// LogOff 注销用户的代理
func (sel *RegisterCenter) LogOff(name string, id int64) {
	sel.delProxy(name, id)
}

type registerConn struct {
	conn.Conn
}

func (sel *registerConn) Read(buf []byte) (int, error) {
	n, err := sel.Conn.Read(buf)
	if err != nil {
		vlog.ERROR("读取消息错误了%v, error : %v", sel.Conn.RemoteAddr(), err)
	}
	return n, err
}
