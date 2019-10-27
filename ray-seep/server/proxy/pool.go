// @File     : proxy
// @Author   : Ville
// @Time     : 19-9-27 下午4:59
// proxy
package proxy

import (
	"ray-seep/ray-seep/common/conn"
	"ray-seep/ray-seep/common/errs"
	"sync"
	"time"
)

type Pool interface {
	Push(key int64, c conn.Conn) error
	Get(key int64) (conn.Conn, error)
}

type element struct {
	ct time.Time
	c  conn.Conn
}

// 代理链接的缓存池
type pool struct {
	cntCache int64         // 当前缓存数量
	expire   time.Duration // 链接到期时间
	lock     sync.Mutex
	pxyConn  map[int64]*element
}

func NewPool() Pool {
	return &pool{
		expire:  time.Second * 5,
		pxyConn: make(map[int64]*element),
	}
}

// 循环查询到期时间，到期后自动销毁
func (p *pool) loopExpire() {
	go func() {
		for {
			for k, v := range p.pxyConn {
				// 如果判断是否过期了
				if v.ct.Before(time.Now().Add(-1 * p.expire)) {
					p.lock.Lock()
					delete(p.pxyConn, k)
					p.lock.Unlock()
				}
			}
		}
	}()
}

func (p *pool) Push(key int64, c conn.Conn) error {
	p.lock.Lock()
	p.pxyConn[key] = &element{ct: time.Now(), c: c}
	p.cntCache++
	p.lock.Unlock()
	return nil
}

func (p *pool) Get(key int64) (conn.Conn, error) {
	p.lock.Lock()
	if c, ok := p.pxyConn[key]; ok {
		c.ct = time.Now()
		p.lock.Unlock()
		return c.c, nil
	}
	p.lock.Unlock()
	// 如果没找到
	return nil, errs.ErrProxyNotExist
}
