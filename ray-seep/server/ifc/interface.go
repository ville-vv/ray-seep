package ifc

import "ray-seep/ray-seep/common/conn"

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
