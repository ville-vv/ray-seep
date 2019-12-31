package control

import (
	"fmt"
	"net"
	"ray-seep/ray-seep/common/errs"
	"ray-seep/ray-seep/databus"
	"ray-seep/ray-seep/model"
)

type PodHandler struct {
	db databus.BaseDao
}

func NewPodHandler(db databus.BaseDao) *PodHandler {
	return &PodHandler{db: db}
}

func (sel *PodHandler) randPort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func (sel *PodHandler) OnLogin(connId, userId int64, user string, appKey string, token string) (loginDao *model.UserLoginDao, err error) {
	loginDao, err = sel.db.UserLogin(connId, userId, user, appKey, token)
	if err != nil {
		return
	}
	if loginDao.HttpPort == "" {
		return nil, errs.ErrHttpPortIsInValid
	}
	if user == "test" {
		port, err := sel.randPort()
		if err != nil {
			return nil, err
		}
		loginDao.HttpPort = fmt.Sprintf("%d", port)
	}
	return
}

func (sel *PodHandler) OnLogout(name string, id int64, isClean bool) error {
	sel.db.DelToken(id, name, isClean)
	return nil
}

// 创建主机判断是否登录
func (sel *PodHandler) OnCreateHost(connId int64, user string, token string) error {
	if token != sel.db.GetToken(connId, user) || token == "" {
		return errs.ErrNoLogin
	}
	return nil
}
