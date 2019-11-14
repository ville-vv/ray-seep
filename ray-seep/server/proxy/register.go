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
	"sync"
	"vilgo/vlog"
)

type MessagePusher interface {
	PushMsg(id int64, p *proto.Package) error
}

type IDChooseRuler interface {
	Choose([]int64) int64
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

// 节点ID存放列表
type nodeIdList struct {
	name   string
	cho    IDChooseRuler
	idList []int64 // 一个域名下可能存在多个服务，用于实现用户服务负载均衡功能使用特定算法获取一条 cid对应的conn
}

func newNodeIdList(name string) *nodeIdList {
	return &nodeIdList{name: name, cho: &idChooser{}}
}
func (sel *nodeIdList) Name() string {
	return sel.name
}
func (sel *nodeIdList) Add(id int64) {
	for i := range sel.idList {
		if sel.idList[i] == id {
			return
		}
	}
	sel.idList = append(sel.idList, id)
}
func (sel *nodeIdList) Del(id int64) {
	for i := 0; i < len(sel.idList); i++ {
		if id == sel.idList[i] {
			sel.idList = append(sel.idList[:i], sel.idList[i+1:]...)
		}
	}
}
func (sel *nodeIdList) Len() int {
	return len(sel.idList)
}
func (sel *nodeIdList) Get() int64 {
	return sel.cho.Choose(sel.idList)
}

// RegisterCenter 注册中心，记录用户启动的本地服务id与用户使用的域名映射
// 记录用户启动的服务的代理池
type RegisterCenter struct {
	lock       sync.RWMutex
	nodeIdsNum int                    // 节点 Id 数量
	nodes      map[string]*nodeIdList // 域名映射 一个域名对多个服务节点
	pxyPool    conn.Pool              // 记录用户本地服务的代理 tcp 链接，使用 cid 获取链接
	pushMsg    MessagePusher
}

func NewRegisterCenter(pl conn.Pool, ph MessagePusher) *RegisterCenter {
	return &RegisterCenter{
		nodes:   make(map[string]*nodeIdList),
		pxyPool: pl,
		pushMsg: ph,
	}
}

// 注册用户链接
func (sel *RegisterCenter) Register(domain string, id int64, cc conn.Conn) error {
	if err := sel.pxyPool.Push(id, cc); err != nil {
		vlog.ERROR("register add node id[%d] error %s", cc.Id(), err.Error())
		return err
	}
	return sel.addNodeId(domain, id)
}

// addDmp 根据域名添加一个 客户端的链接ID
func (sel *RegisterCenter) addNodeId(domain string, cid int64) error {
	sel.lock.Lock()
	defer sel.lock.Unlock()
	if d, ok := sel.nodes[domain]; ok {
		d.Add(cid)
		return nil
	}
	idList := newNodeIdList(domain)
	idList.Add(cid)
	sel.nodes[domain] = idList
	return nil
}

// GetProxy 获取代理tcp连接
func (sel *RegisterCenter) GetProxy(domain string) (net.Conn, error) {
	sel.lock.RLock()
	ids, ok := sel.nodes[domain]
	if !ok {
		sel.lock.RUnlock()
		return nil, errs.ErrServerNotExist
	}
	sel.lock.RUnlock()
	return sel.pxyPool.Get(ids.Get())
}
