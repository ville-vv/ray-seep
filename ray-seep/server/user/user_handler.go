package user

import (
	"fmt"
	"net"
	"ray-seep/ray-seep/common/errs"
	"ray-seep/ray-seep/common/util"
	"ray-seep/ray-seep/databus"
)

type Handler struct {
	db databus.BaseDao
}

func NewHandler(db databus.BaseDao) *Handler {
	return &Handler{db: db}
}

func (sel *Handler) randPort() (int, error) {
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

func (sel *Handler) OnLogin(connId, userId int64, user string, appKey string) (token, srvPort string, err error) {
	// 查询用户是否已经登录
	if sel.db.GetToken(connId, user) != "" && user != "test" {
		return "", "", errs.ErrUserHaveLogin
	}
	// 查询用户信息
	loginDao, err := sel.db.UserLogin(connId, userId, user, appKey)
	if err != nil {
		return
	}
	srvPort = loginDao.HttpPort
	if srvPort == "" {
		return "", "", errs.ErrHttpPortIsInValid
	}
	if user == "test" {
		port, err := sel.randPort()
		if err != nil {
			return "", "", err
		}
		srvPort = fmt.Sprintf("%d", port)
	}
	// 没有登录，随机生成一个 token
	token = util.RandToken()
	return token, srvPort, sel.db.SaveToken(connId, user, token)
}

func (sel *Handler) OnLogout(name string, id int64) error {
	sel.db.DelToken(id, name)
	return nil
}

// 创建主机判断是否登录
func (sel *Handler) OnCreateHost(connId int64, user string, token string) error {
	if token != sel.db.GetToken(connId, user) || token == "" {
		return errs.ErrNoLogin
	}
	return nil
}
