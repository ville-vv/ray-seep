// @File     : CustomerManage
// @Author   : Ville
// @Time     : 19-9-24 下午4:46
// server
package cust

import (
	"errors"
	"ray-seep/ray-seep/common/pkg"
	"ray-seep/ray-seep/conn"
	"ray-seep/ray-seep/mng"
	"ray-seep/ray-seep/session"
	"sync"
	"vilgo/vlog"
)

type Handler interface {
	Connect(conn.Conn) error
	Handler(conn.Conn)
	DisConnect(id int64)
}

// 客户连接管理
type CustomerMng struct {
	connMng *mng.ConnManager
	msgMng  *mng.MsgManager
}

func NewCustomerMng() *CustomerMng {
	return &CustomerMng{
		connMng: mng.NewConnManager(),
		msgMng:  mng.NewMsgManager(),
	}
}

func (cs *CustomerMng) Connect(conn conn.Conn) error {
	msgMng := mng.NewMsgTransfer(conn)
	authMsg := pkg.Package{}
	if err := msgMng.RecvMsg(&authMsg); err != nil {
		vlog.ERROR("receive client auth message error：%v", err)
		return err
	}

	if authMsg.Cmd != pkg.CmdIdentifyReq {
		return errors.New("client authentication fail")
	}
	// 发送验证通过消息

	if err := msgMng.SendMsg(pkg.NewPackage(pkg.CmdIdentifyRsp, &session.LoginRsp{Id: conn.Id(), Token: "abc"})); err != nil {
		return errors.New("client authentication result send fail")
	}

	// 权限认证成功即可放入连接池中
	if err := cs.connMng.Put(conn.Id(), conn); err != nil {
		return err
	}
	// 建立消息管理器
	//cs.msgMng.Put(conn.Id(), msgMng)
	vlog.INFO("client [%d] connect success", conn.Id())
	return nil
}

func (cs *CustomerMng) DisConnect(id int64) {
	vlog.INFO("client exit [%d]", id)
	cs.connMng.Delete(id)
}

func (cs *CustomerMng) Handler(c conn.Conn) {
	if _, ok := cs.connMng.Get(c.Id()); !ok {
		cs.connMng.Put(c.Id(), c)
	}
	msgMng := mng.NewMsgTransfer(c)
	wait := sync.WaitGroup{}
	recvMsg := make(chan pkg.Package)
	cancel := make(chan pkg.Package)
	msgMng.RecvMsgWithChan(&wait, recvMsg, cancel)
	wait.Wait()
	for {
		select {
		case m := <-recvMsg:
			vlog.DEBUG("[%d] [action] [%v]", c.Id(), m.Cmd)
		case <-cancel:
			return
		}
	}
}

func (cs *CustomerMng) CmdHandler(pkg pkg.Package) {
}
