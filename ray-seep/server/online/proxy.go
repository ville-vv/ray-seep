package online

import (
	"ray-seep/ray-seep/common/conn"
	"vilgo/vlog"
)

type ProxyPool struct {
	index     int
	name      string
	cIds      []int64
	cntUseCid int64 // 当前使用的ID
	conn.Pool
}

func NewProxyPool(name string, pool conn.Pool) *ProxyPool {
	return &ProxyPool{name: name, Pool: pool}
}

func (p *ProxyPool) Push(key int64, c conn.Conn) error {
	p.cIds = append(p.cIds, key)
	return p.Pool.Push(key, c)
}

func (p *ProxyPool) Get(key int64) (conn.Conn, error) {
	return p.Pool.Get(key)
}

func (p *ProxyPool) Drop(id int64) {
	for i := 0; i < len(p.cIds); i++ {
		if id == p.cIds[i] {
			p.cIds = append(p.cIds[:i], p.cIds[i+1:]...)
		}
	}
	vlog.DEBUG(" drop 当前代理数 [%s][%d]", p.name, len(p.cIds))
	p.Pool.Drop(id)
}

func (p *ProxyPool) Size() int {
	return len(p.cIds)
}

func (p *ProxyPool) GetId() int64 {
	if len(p.cIds) <= 0 {
		return 0
	}
	id := p.cIds[p.index]
	p.index++
	if p.index >= len(p.cIds) {
		p.index = 0
	}
	p.cntUseCid = id
	return id
}
