// @File     : register
// @Author   : Ville
// @Time     : 19-10-12 下午3:32
// proxy
package proxy

import (
	"github.com/vilsongwei/vilgo/vlog"
	"net"
	"ray-seep/ray-seep/common/conn"
	"ray-seep/ray-seep/common/errs"
	"ray-seep/ray-seep/server_v2/ifc"
	"sync"
	"time"
)

type registerItem struct {
	Name string
	Id   int64
	conn.Pool
}

// RegisterCenter 注册中心，记录用户启动的本地服务id与用户使用的域名映射
// 记录用户启动的服务的代理池
type RegisterCenter struct {
	lock     sync.RWMutex
	pxyPools map[string]*registerItem // 记录用户本地服务的代理 tcp 链接，使用 cid 获取链接
	ntc      ifc.NoticeGetter
	caches   int
}

func NewRegisterCenter(caches int) *RegisterCenter {
	return &RegisterCenter{
		pxyPools: make(map[string]*registerItem),
		caches:   caches, // 一个节点需能缓存的数量
	}
}

func (sel *RegisterCenter) SetNoticeGetter(ph ifc.NoticeGetter) {
	sel.ntc = ph
}

// 注册用户链接
func (sel *RegisterCenter) Register(name string, id int64, cc conn.Conn) error {
	// 把tcp连接放到代理池中
	if err := sel.addProxy(name, id, cc); err != nil {
		vlog.ERROR("[%s][%d]register proxy error %s", name, id, err.Error())
		return err
	}
	notice, err := sel.ntc.GetNotice(id)
	if err != nil {
		vlog.ERROR("[%s][%d]register proxy notice error %s", name, id, err.Error())
		return err
	}
	return notice.NoticeRunProxyRsp([]byte("{}"))
}

func (sel *RegisterCenter) addProxy(name string, id int64, cc conn.Conn) error {
	sel.lock.Lock()
	defer sel.lock.Unlock()
	if p, ok := sel.pxyPools[name]; ok {
		return p.Push(id, cc)
	}
	pl := &registerItem{Name: name, Id: id, Pool: conn.NewPoolV2(sel.caches)}
	if err := pl.Push(id, cc); err != nil {
		return err
	}
	sel.pxyPools[name] = pl
	// vlog.INFO("[%s][%d]当前代理数量%d", name, id, pl.Size())
	return nil
}

func (sel *RegisterCenter) delProxy(name string, cid int64) (clean bool) {
	if pl, ok := sel.pxyPools[name]; ok {
		pl.Drop(cid)
		if pl.Size() == 0 {
			pl.Close()
			delete(sel.pxyPools, name)
			clean = true
		}
	}
	return
}

// GetNetConn 获取代理tcp连接
func (sel *RegisterCenter) GetNetConn(name string) (net.Conn, error) {
	sel.lock.RLock()
	pl, ok := sel.pxyPools[name]
	sel.lock.RUnlock()
	if !ok {
		return nil, errs.ErrProxySrvNotExist
	}
	cn, err := pl.Get(0)
	if err != nil {
		return sel.getAndRunProxy(name, pl)
	}
	return cn, err
}

func (sel *RegisterCenter) getAndRunProxy(name string, pl *registerItem) (net.Conn, error) {
	id := pl.Id
	if err := sel.noticeRunProxy(name, id); err != nil {
		vlog.ERROR("[%s][%d] push notice run proxy message error %s", name, id, err.Error())
		return nil, errs.ErrNoticeProxyRunErr
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
		vlog.WARN("[%s][%d] wait get proxy timeout", name, id)
	}
	return nil, errs.ErrWaitProxyRunTimeout
}

func (sel *RegisterCenter) noticeRunProxy(name string, id int64) error {
	notice, err := sel.ntc.GetNotice(id)
	if err != nil {
		vlog.ERROR("[%s][%d]notice run proxy error %s", name, id, err.Error())
		return err
	}
	return notice.NoticeRunProxy([]byte("{}"))
}

// LogOff 注销用户的代理
func (sel *RegisterCenter) LogOff(name string, id int64) (clean bool) {
	return sel.delProxy(name, id)
}
