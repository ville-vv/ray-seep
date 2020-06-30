package ifc

import (
	"ray-seep/ray-seep/common/conn"
)

type MessageNotice interface {
	NoticeRunProxy(data []byte) error
	NoticeRunProxyRsp(data []byte) error
}

type NoticeGetter interface {
	GetNotice(id int64) (MessageNotice, error)
}

type ExitDevice interface {
	Logout(name string, id int64)
}

type Register interface {
	// 注册
	Register(name string, id int64, cc conn.Conn) error
	// 注销
	Logout(name string, id int64)
}

type PodHandler interface {

	// 用户登录操作
	OnLogin(connId, userId int64, user string, appKey string) (token string, port string, err error)
	// 用户创建服务主机操作
	OnCreateHost(connId int64, user string, token string) error
	// 用户登出操作
	OnLogout(name string, id int64) error
	KeepLive(userName string, connID int64)
}
