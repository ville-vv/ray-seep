// @File     : CustomerManage
// @Author   : Ville
// @Time     : 19-9-24 下午4:46 
// server 
package cust

import (
	"ray-seep/ray-seep/mng"
	"vilgo/vlog"
)

// 客户连接管理
type CustomerMng struct {
	connMng *mng.ConnManager
}

func NewCustomerMng()*CustomerMng{
	return &CustomerMng{
		connMng:mng.NewConnManager(),
	}
}

func (cs *CustomerMng) Connect(conn mng.Conn) error {
	vlog.INFO("客户端链接：%v", conn.RemoteAddr())
	return nil
}

func (cs *CustomerMng) DisConnect(id int64) {
	conn, ok := cs.connMng.Get(id)
	if ok{
		vlog.INFO("客户端断开：%v", conn.RemoteAddr())
	}
}


