// @File     : proxy
// @Author   : Ville
// @Time     : 19-9-27 下午4:59
// proxy
package conn

import (
	"ray-seep/ray-seep/common/errs"
	"sync"
	"time"
	"vilgo/vlog"
)

type Pool interface {
	Push(key int64, c Conn) error
	Get(key int64) (Conn, error)
	Drop(key int64)
}

type element struct {
	ct time.Time
	id int64
	c  Conn
}

// 代理链接的缓存池
type pool struct {
	cntCache int64         // 当前缓存数量
	expire   time.Duration // 链接到期时间
	lock     sync.Mutex
	pxyConn  map[int64]*element
	addCh    chan element
}

func NewPool() Pool {
	return &pool{
		expire:  time.Second * 5,
		pxyConn: make(map[int64]*element),
		addCh:   make(chan element, 100),
	}
}

// 循环查询到期时间，到期后自动销毁
func (p *pool) loopExpire() {
	go func() {
		for {
			for k, v := range p.pxyConn {
				// 如果判断是否过期了
				if v.ct.Before(time.Now().Add(-1 * p.expire)) {
					p.Drop(k)
				}
			}
		}
	}()
}

func (p *pool) Push(key int64, c Conn) error {
	select {
	case p.addCh <- element{ct: time.Now(), c: c, id: key}:
		vlog.WARN("pool push fail, proxies is full")
	}
	return nil
}

func (p *pool) Get(key int64) (Conn, error) {
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

func (p *pool) Drop(key int64) {
	p.lock.Lock()
	delete(p.pxyConn, key)
	p.lock.Unlock()
}

// 选好检测 add chan 添加链接，
func (p *pool) loopCheckAdd() {
	for ele := range p.addCh {
		p.lock.Lock()
		p.cntCache++
		p.pxyConn[ele.id] = &ele
		p.lock.Unlock()
	}
}
