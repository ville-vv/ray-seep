// @File     : register
// @Author   : Ville
// @Time     : 19-10-12 下午3:32
// proxy
package proxy

import (
	"net"
	"ray-seep/ray-seep/common/conn"
	"ray-seep/ray-seep/common/errs"
	"sync"
	"vilgo/vlog"
)

type IDChooseRuler interface {
	Choose([]int64) int64
}

type nodeIdList struct {
	Name   string
	IdList []int64 // 一个域名下可能存在多个服务，用于实现用户服务负载均衡功能使用特定算法获取一条 cid对应的conn
}

// RegisterCenter 注册中心，记录用户启动的本地服务id与用户使用的域名映射
// 记录用户启动的服务的代理池
type RegisterCenter struct {
	lock       sync.RWMutex
	nodeIdsNum int                    // 节点 Id 数量
	nodes      map[string]*nodeIdList // 域名映射 一个域名对多个服务节点
	pxyPool    Pool                   // 记录用户本地服务的代理 tcp 链接，使用 cid 获取链接
	rule       IDChooseRuler          // 选择器
}

func NewRegisterCenter(pl Pool) *RegisterCenter {
	return &RegisterCenter{
		nodes:   make(map[string]*nodeIdList),
		pxyPool: pl,
		rule:    &idChooser{},
	}
}

// 注册用户链接
func (sel *RegisterCenter) Register(domain string, cc conn.Conn) error {
	if err := sel.addNodeId(domain, cc.Id()); err != nil {
		vlog.ERROR("register add node id[%d] error %s", cc.Id(), err.Error())
		return err
	}
	return sel.pxyPool.Push(cc.Id(), cc)
}

// addDmp 根据域名添加一个 客户端的链接ID
func (sel *RegisterCenter) addNodeId(domain string, cid int64) error {
	sel.lock.Lock()
	defer sel.lock.Unlock()
	if d, ok := sel.nodes[domain]; ok {
		if len(d.IdList) >= sel.nodeIdsNum {
			return errs.ErrServerNumFull
		}
		d.IdList = append(d.IdList, cid)
		return nil
	}
	idList := make([]int64, 0, sel.nodeIdsNum)
	idList = append(idList, cid)
	sel.nodes[domain] = &nodeIdList{Name: domain, IdList: idList}

	return nil
}

// GetProxy 获取代理tcp连接
func (sel *RegisterCenter) GetProxy(domain string) (net.Conn, error) {
	sel.lock.RLock()
	dmp, ok := sel.nodes[domain]
	if !ok {
		sel.lock.RUnlock()
		return nil, errs.ErrServerNotExist
	}
	sel.lock.RUnlock()
	return sel.pxyPool.Get(sel.rule.Choose(dmp.IdList))
}

// Id 获取器
type idChooser struct {
	index int
}

// Choose 获取 Id
func (sel *idChooser) Choose(ids []int64) int64 {
	if len(ids) <= 0 {
		return 0
	}
	id := ids[sel.index]
	sel.index++
	if sel.index >= len(ids) {
		sel.index = 0
	}
	return id
}
