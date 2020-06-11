package conn

import (
	"ray-seep/ray-seep/common/errs"
	"sync/atomic"
)

// 代理链接的缓存池
type poolV2 struct {
	cntCacheNum int64 // 当前缓存数量
	maxCacheNum int
	caches      chan Conn
}

func NewPoolV2(cacheNum int) Pool {
	p := &poolV2{
		caches:      make(chan Conn, cacheNum),
		maxCacheNum: cacheNum,
	}
	return p
}

func (p *poolV2) Push(key int64, c Conn) error {
	p.caches <- c
	atomic.AddInt64(&p.cntCacheNum, 1)
	return nil
}

func (p *poolV2) Get(key int64) (Conn, error) {
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

func (p *poolV2) WaitGet() <-chan Conn {
	atomic.AddInt64(&p.cntCacheNum, -1)
	return p.caches
}

func (p *poolV2) Drop(key int64) {
	close(p.caches)
	for v := range p.caches {
		_ = v.Close()
	}
	p.caches = nil
	p.caches = make(chan Conn, p.maxCacheNum)
}

func (p *poolV2) Size() int {
	l := atomic.LoadInt64(&p.cntCacheNum)
	return int(l)
}
func (p *poolV2) Close() {
	close(p.caches)
	for v := range p.caches {
		_ = v.Close()
	}
}
