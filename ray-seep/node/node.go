// @File     : node
// @Author   : Ville
// @Time     : 19-9-26 下午4:38
// node
package node

import (
	"sync"
)

type PodPolicyMaker interface {
	Make([]Pod) Pod
}

// Node 是一个用户的管理器，一个用户对应一个 domain, 一个node 下面存在多个 pod
// pod 代表一个服务
type Node struct {
	lock   sync.RWMutex
	name   string
	domain string
	maxPod int // pod 最大数量
	pods   []Pod
	policy PodPolicyMaker // 决策器，用来决定Pod调用方法
}

func NewNode() *Node {
	return &Node{}
}

func (sel *Node) Make(ps []Pod) Pod {
	return ps[0]
}

// 获取pod
func (sel *Node) GetPod() Pod {
	sel.lock.RLock()
	pod := sel.policy.Make(sel.pods)
	sel.lock.RUnlock()
	return pod
}

// 添加一个pod
func (sel *Node) AddPod(p Pod) {
	sel.pods = append(sel.pods, p)
}
