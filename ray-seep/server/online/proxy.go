package online

import (
	"ray-seep/ray-seep/common/conn"
)

type ProxyPool struct {
	index     int
	name      string
	cIds      []int64
	cIdCheck  map[int64]int
	cntUseCid int64 // 当前使用的ID
	conn.Pool
}

func NewProxyPool(name string, pool conn.Pool) *ProxyPool {
	return &ProxyPool{name: name, Pool: pool, cIdCheck: make(map[int64]int)}
}

func (p *ProxyPool) Push(key int64, c conn.Conn) error {
	if _, ok := p.cIdCheck[key]; !ok {
		p.cIds = append(p.cIds, key)
		p.cIdCheck[key] = len(p.cIds) - 1
	}
	return p.Pool.Push(key, c)
}

func (p *ProxyPool) Get(key int64) (conn.Conn, error) {
	return p.Pool.Get(key)
}

func (p *ProxyPool) Drop(id int64) {
	if _, ok := p.cIdCheck[id]; ok {
		for i, v := range p.cIds {
			if v == id {
				p.cIds = append(p.cIds[:i], p.cIds[i+1:]...)
			}
		}
		delete(p.cIdCheck, id)
	}
	p.Pool.Drop(0)
}

func (p *ProxyPool) Size() int {
	return len(p.cIdCheck)
}

func (p *ProxyPool) GetId() int64 {
	if len(p.cIds) <= 0 {
		return 0
	}
	if p.index >= len(p.cIdCheck) {
		p.index = 0
	}
	id := p.cIds[p.index]
	p.index++
	p.cntUseCid = id
	return id
}
