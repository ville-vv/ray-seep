// @File     : proxy
// @Author   : Ville
// @Time     : 19-9-27 下午4:59
// proxy
package conn

import (
	"github.com/vilsongwei/vilgo/vlog"
	"ray-seep/ray-seep/common/errs"
	"sync"
	"sync/atomic"
	"time"
)

type Pool interface {
	Push(key int64, c Conn) error
	Get(key int64) (Conn, error)
	WaitGet() (Conn, error)
	Size() int
	Inc() int64
	Drop(key int64)
	Close()
}

type element struct {
	ct time.Time
	id int64
	c  Conn
}

// 代理链接的缓存池
type pool struct {
	cntCacheNum int64 // 当前缓存数量
	maxCacheNum int
	expire      time.Duration // 链接到期时间
	lock        sync.Mutex
	pxyConn     map[int64]*element
	caches      chan Conn
	expCh       chan int64
}

func NewPool(cacheNum int) Pool {
	p := &pool{
		expire:      time.Second * time.Duration(30),
		pxyConn:     make(map[int64]*element),
		caches:      make(chan Conn, cacheNum),
		expCh:       make(chan int64, 10000),
		maxCacheNum: cacheNum,
	}
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
					//vlog.DEBUG("当前时间：%s 到期时间%s", time.Now().Format("2006-01-01 15:04:05"), v.ct.Add(p.expire).Format("2006-01-01 15:04:05"))
					// 如果判断是否过期了
					if v.ct.Before(time.Now().Add(-1 * p.expire)) {
						vlog.WARN("[%d] expire", v.id)
						p.expCh <- v.id
						p.Drop(k)
						_ = v.c.Close()
					}
				}
			}
		}
	}()
}

func (p *pool) Push(key int64, c Conn) error {
	p.caches <- c
	atomic.AddInt64(&p.cntCacheNum, 1)
	return nil
}

func (p *pool) Get(key int64) (Conn, error) {
	select {
	case c, ok := <-p.caches:
		if !ok {
			return nil, errs.ErrProxyConnNotExist
		}
		atomic.AddInt64(&p.cntCacheNum, -1)
		return c, nil
	default:
		return nil, errs.ErrProxyConnNotExist
	}
}

func (p *pool) WaitGet() (Conn, error) {
	tm := time.NewTicker(time.Second * 5)
	select {
	case cn, ok := <-p.caches:
		if !ok {
			return nil, errs.ErrProxySrvNotExist
		}
		atomic.AddInt64(&p.cntCacheNum, -1)
		return cn, nil
	case <-tm.C:
		return nil, errs.ErrWaitProxyRunTimeout
	}
}
func (p *pool) Inc() int64 {
	l := atomic.LoadInt64(&p.cntCacheNum)
	return l
}
func (p *pool) Drop(key int64) {
	close(p.caches)
	for v := range p.caches {
		_ = v.Close()
	}
	p.caches = nil
	p.caches = make(chan Conn, p.maxCacheNum)
}

func (p *pool) Size() int {
	l := atomic.LoadInt64(&p.cntCacheNum)
	return int(l)
}
func (p *pool) Close() {
	close(p.caches)
	for v := range p.caches {
		_ = v.Close()
	}
	close(p.expCh)
}
