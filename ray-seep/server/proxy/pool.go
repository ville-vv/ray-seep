// @File     : proxy
// @Author   : Ville
// @Time     : 19-9-27 下午4:59
// proxy
package proxy

import (
	"ray-seep/ray-seep/common/conn"
	"sync"
	"time"
)

type Pool interface {
	Push(key int64, c conn.Conn)
	Get(key int64) (conn.Conn, bool)
}

type element struct {
	ct time.Time
	c  conn.Conn
}

// 代理链接的缓存池
type pool struct {
	maxCache int64         // 最大缓存数量
	cntCache int64         // 当前缓存数量
	expire   time.Duration // 链接到期时间
	sync.Mutex
	pxyConn map[int64]*element
}

// 循环查询到期时间，到期后自动销毁
func (p *pool) loopExpire() {
	go func() {
		for {
			for k, v := range p.pxyConn {
				// 如果判断是否过期了
				if v.ct.Before(time.Now().Add(-1 * p.expire)) {
					p.Lock()
					delete(p.pxyConn, k)
					p.Unlock()
				}
			}
		}
	}()
}

func (p *pool) Push(key int64, c conn.Conn) {
	p.Lock()
	if p.maxCache < p.cntCache {
		return
	}
	p.pxyConn[key] = &element{ct: time.Now(), c: c}
	p.Unlock()
}

func (p *pool) Get(key int64) (conn.Conn, bool) {
	p.Lock()
	if c, ok := p.pxyConn[key]; ok {
		c.ct = time.Now()
		p.Unlock()
		return c.c, ok
	}
	p.Unlock()
	// 如果没找到
	return nil, false
}
