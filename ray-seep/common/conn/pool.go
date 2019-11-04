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
	p := &pool{
		expire:  time.Second * 5,
		pxyConn: make(map[int64]*element),
		addCh:   make(chan element, 100),
	}
	p.loopExpire()
	return p
}

// 循环查询到期时间，到期后自动销毁
func (p *pool) loopExpire() {
	go func() {
		tk := time.NewTicker(time.Second * 1)
		for {
			select {
			case <-tk.C:
				for k, v := range p.pxyConn {
					// 如果判断是否过期了
					if v.ct.Before(time.Now().Add(-1 * p.expire)) {
						vlog.ERROR("超时 %d", v.id)
						_ = v.c.Close()
						p.Drop(k)
					}
				}
			}
		}
	}()
}

func (p *pool) Push(key int64, c Conn) error {
	p.lock.Lock()
	p.cntCache++
	p.pxyConn[key] = &element{ct: time.Now(), id: key, c: c}
	p.lock.Unlock()
	return nil
}

func (p *pool) Get(key int64) (Conn, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if c, ok := p.pxyConn[key]; ok {
		// 有获取这个链接就重置时间
		c.ct = time.Now()
		return c.c, nil
	}
	// 如果没找到
	return nil, errs.ErrProxyNotExist
}

func (p *pool) Drop(key int64) {
	p.lock.Lock()
	delete(p.pxyConn, key)
	p.lock.Unlock()
}
