// @File     : mng_conn
// @Author   : Ville
// @Time     : 19-10-12 下午2:58
// mng
package mng

import (
	"errors"
	"ray-seep/ray-seep/common/conn"
	"sync"
	"vilgo/vlog"
)

const (
	MaxLinkNumber = 100 // 客户端的最大连接数
)

// 连接管理
type ConnManager struct {
	cuMap      map[int64]conn.Conn
	cntLinkNum uint32
	maxLinkNum uint32
	sync.RWMutex
}

// client 连接管理
func NewConnManager() (cm *ConnManager) {
	cm = new(ConnManager)
	cm.cuMap = make(map[int64]conn.Conn)
	cm.cntLinkNum = 0
	cm.maxLinkNum = MaxLinkNumber
	return
}

func (cm *ConnManager) Put(key int64, cu conn.Conn) error {
	if cm.cntLinkNum >= cm.maxLinkNum {
		return errors.New("connect number is full")
	}
	cm.Lock()
	defer cm.Unlock()
	cm.cuMap[key] = cu
	cm.cntLinkNum++
	return nil
}

func (cm *ConnManager) Get(key int64) (conn.Conn, bool) {
	cm.RLock()
	defer cm.RUnlock()
	cu, ok := cm.cuMap[key]
	return cu, ok
}

func (cm *ConnManager) Delete(key int64) {
	cm.Lock()
	defer cm.Unlock()
	if cu, ok := cm.cuMap[key]; ok {
		cu.Close()
		delete(cm.cuMap, key)
		cm.cntLinkNum--
		vlog.LogD("当前连接数：%v", cm.cntLinkNum)
	}
}
