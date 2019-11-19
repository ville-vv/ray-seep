package online

import (
	"sync"
	"vilgo/vlog"
)

// 登录管理
type OnLiner interface {
	IsOnLine(cid int64, token string) bool
}

type pod struct {
	id    int64
	token string
	user  string
}

type User struct {
	index int
	name  string
	pods  []*pod
}

func (sel *User) add(id int64, name, token string) {
	sel.pods = append(sel.pods, &pod{id: id, token: token, user: name})
}
func (sel *User) drop(id int64) {
	for i := 0; i < len(sel.pods); i++ {
		if id == sel.pods[i].id {
			sel.pods = append(sel.pods[:i], sel.pods[i+1:]...)
		}
	}
	vlog.DEBUG("Del当前代理数[%s][%d]", sel.name, len(sel.pods))
}
func (sel *User) Len() int {
	return len(sel.pods)
}

// Choose 获取 Id
func (sel *User) Choose() int64 {
	if len(sel.pods) <= 0 {
		return 0
	}
	pod := sel.pods[sel.index]
	sel.index++
	if sel.index >= len(sel.pods) {
		sel.index = 0
	}
	return pod.id
}

type UserManager struct {
	lock  sync.RWMutex
	users map[string]*User
}

func NewUserManager() *UserManager {
	return &UserManager{
		users: make(map[string]*User),
		lock:  sync.RWMutex{},
	}
}

func (l *UserManager) GetId(name string) int64 {
	l.lock.RLock()
	defer l.lock.RUnlock()
	if user, ok := l.users[name]; ok {
		return user.Choose()
	}
	return 0
}
func (l *UserManager) Login(id int64, name, token string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if user, ok := l.users[name]; ok {
		user.add(id, name, token)
		return
	}
	user := &User{index: 0, name: name}
	user.add(id, name, token)
	l.users[name] = user
	return
}
func (l *UserManager) Logout(id int64, name string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if user, ok := l.users[name]; ok {
		user.drop(id)
		if user.Len() == 0 {
			delete(l.users, name)
		}
	}
}
