// @File     : register
// @Author   : Ville
// @Time     : 19-10-12 下午3:32
// proxy
package proxy

import (
	"net"
	"ray-seep/ray-seep/common/errs"
	"sync"
)

type IdRuler interface {
	Opt([]int64) int64
}

type domainMap struct {
	Domain string
	IdList []int64 // 一个域名下可能存在多个服务，用于实现用户服务负载均衡功能使用特定算法获取一条 cid对应的conn
}

// RegisterCenter 注册中心，记录用户启动的本地服务id与用户使用的域名映射
// 记录用户启动的服务的代理池
type RegisterCenter struct {
	lock      sync.RWMutex
	dmpIdsNum int
	dmp       map[string]*domainMap // 域名映射 一个域名对多个
	pxyPool   Pool                  // 记录用户本地服务的代理 tcp 链接，使用 cid 获取链接
	rule      IdRuler               // 选择器
}

// 注册用户链接
func (sel *RegisterCenter) Register(domain string, cid int64) error {
	return sel.addDmp(domain, cid)
}

// addDmp 根据域名添加一个 客户端的链接ID
func (sel *RegisterCenter) addDmp(domain string, cid int64) error {
	sel.lock.Lock()
	if d, ok := sel.dmp[domain]; ok {
		if len(d.IdList) >= sel.dmpIdsNum {
			return errs.ErrServerNumFull
		}
		d.IdList = append(d.IdList, cid)
		return nil
	}
	idList := make([]int64, 0, sel.dmpIdsNum)
	idList = append(idList, cid)
	sel.dmp[domain] = &domainMap{Domain: domain, IdList: idList}
	sel.lock.Unlock()
	return nil
}

// 获取代理tcp连接
func (sel *RegisterCenter) GetProxy(domain string) (net.Conn, error) {
	sel.lock.RLock()
	dmp, ok := sel.dmp[domain]
	if !ok {
		return nil, errs.ErrServerNotExist
	}
	sel.lock.RUnlock()
	cid := sel.rule.Opt(dmp.IdList)
	return sel.pxyPool.Get(cid)
}
