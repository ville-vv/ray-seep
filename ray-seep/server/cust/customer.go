// @File     : customer
// @Author   : Ville
// @Time     : 19-9-26 下午2:09
// cust
package cust

import (
	"ray-seep/ray-seep/conn"
	"ray-seep/ray-seep/mng"
)

type Customer struct {
	conn   conn.Conn
	msgMng mng.MsgManager
}
