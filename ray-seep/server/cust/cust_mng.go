// @File     : CustomerManage
// @Author   : Ville
// @Time     : 19-9-24 下午4:46
// server
package cust

import (
	"errors"
	"ray-seep/ray-seep/conn"
	"ray-seep/ray-seep/mng"
	"ray-seep/ray-seep/msg"
	"sync"
	"time"
	"vilgo/vlog"
)

type Manager interface {
	Connect(conn.Conn) error
	Handler(conn.Conn)
	DisConnect(int64)
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
	authMsg := msg.Message{}
	if err := msgMng.RecvMsg(&authMsg); err != nil {
		vlog.ERROR("receive client auth message error：%v", err)
		return err
	}
	if authMsg.Cmd != "auth" {
		return errors.New("client authentication fail")
	}

	// 权限认证成功即可放入连接池中
	if err := cs.connMng.Put(conn.Id(), conn); err != nil {
		return err
	}
	// 建立消息管理器
	//cs.msgMng.Put(conn.Id(), msgMng)
	vlog.INFO("客户端链接成功：%v", conn.RemoteAddr())
	return nil
}

func (cs *CustomerMng) DisConnect(id int64) {
	vlog.INFO("客户端断开 cid：%d", id)
	cs.msgMng.Delete(id)
	cs.connMng.Delete(id)
}

func (cs *CustomerMng) Handler(c conn.Conn) {
	connMsg, ok := cs.connMng.Get(c.Id())
	if !ok {
		t := time.NewTicker(time.Second * 2)
		select {
		case <-t.C:
			connMsg, ok = cs.connMng.Get(c.Id())
			if !ok {
				vlog.ERROR("无法获取到链接管理 cid：%d , addr %s", c.Id(), c.RemoteAddr())
				return
			}
		}
	}

	msgMng := mng.NewMsgTransfer(connMsg)
	wait := sync.WaitGroup{}
	recvMsg := make(chan msg.Message)
	cancel := make(chan msg.Message)
	msgMng.RecvMsgWithChan(&wait, recvMsg, cancel)
	wait.Wait()
	for {
		select {
		case m := <-recvMsg:
			vlog.INFO("收到客户端[%d]消息：%v", msgMng.Id(), m)
		case m := <-cancel:
			vlog.INFO("客户端[%d]退出：%d", msgMng.Id(), m)
			return
		}
	}
}
