// @File     : connServer
// @Author   : Ville
// @Time     : 19-9-26 下午5:40
// node
package control

import (
	"fmt"
	"ray-seep/ray-seep/common/conn"
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/monitor"
	"ray-seep/ray-seep/proto"
	"runtime/debug"
	"sync"
	"time"
	"vilgo/vlog"
)

type ServerMsgHandler interface {
	OnConnect(id int64, in, out chan proto.Package) error
	OnDisConnect(id int64)
	OnMessage(id int64, p *proto.Package) (proto.Package, error)
}

// 处理用户的通信，接收和发送用户的操作信息
type NodeServer struct {
	mtr     monitor.Monitor
	addr    string
	timeout int64
	ish     ServerMsgHandler
	pushCh  chan proto.Package
}

func (sel *NodeServer) Stop() {}

func (sel *NodeServer) Scheme() string {
	return "control server"
}

func NewNodeServer(ctlCnf *conf.ControlSrv, handler ServerMsgHandler) *NodeServer {
	timeout := ctlCnf.ReadMsgTimeout
	addr := fmt.Sprintf("%s:%d", ctlCnf.Host, ctlCnf.Port)
	if timeout == 0 {
		timeout = 5000
	}
	return &NodeServer{
		mtr:     monitor.NewMonitor("node-server", "counter"),
		timeout: timeout,
		addr:    addr,
		ish:     handler,
		pushCh:  make(chan proto.Package, 1000),
	}
}

func (sel *NodeServer) Start() error {
	lis, err := conn.Listen(sel.addr)
	if err != nil {
		vlog.ERROR("control listen error %v", err)
		return err
	}
	for c := range lis.Conn {
		go sel.dealConn(c)
	}
	return nil
}

// dealConn 处理连接
func (sel *NodeServer) dealConn(c conn.Conn) {
	//vlog.DEBUG("[%d] connecting ", c.Id())

	defer func() {
		if r := recover(); r != nil {
			vlog.LogE("customer listener failed with error %v: %s", r, debug.Stack())
		}
	}()
	defer c.Close()
	// 刚刚建立连接需要设置超时时间
	_ = c.SetReadDeadline(time.Now().Add(time.Duration(sel.timeout) * time.Millisecond))
	msgMng := proto.NewMsgTransfer(c)

	recvMsg := make(chan proto.Package, 10)
	sendMsg := make(chan proto.Package, 10)
	cancel := make(chan interface{})
	wg := sync.WaitGroup{}
	wg.Add(2)

	msgMng.AsyncRecvMsg(&wg, recvMsg, cancel)
	msgMng.AsyncSendMsg(&wg, sendMsg, time.Minute)

	// 通知有用户连接
	if err := sel.ish.OnConnect(c.Id(), recvMsg, sendMsg); err != nil {
		vlog.ERROR("[%d] connect fail %s", c.Id(), err.Error())
		return
	}
	sel.mtr.Inc(1)

	vlog.DEBUG("[%d] connect success", c.Id())
	// 通知有连接断开
	defer func() {
		sel.mtr.Dec(1)
		sel.ish.OnDisConnect(c.Id())
		close(recvMsg)
		close(sendMsg)
		close(cancel)
	}()

	wg.Wait()
	for {
		select {
		case req := <-recvMsg:
			_ = c.SetReadDeadline(time.Now().Add(time.Duration(sel.timeout) * time.Millisecond))
			// 刚刚建立连接需要设置超时时间
			rsp, err := sel.ish.OnMessage(c.Id(), &req)
			if err != nil {
				// 执行消息出现错误
				rsp = proto.Package{Cmd: proto.CmdError, Body: []byte(err.Error())}
			}
			sendMsg <- rsp
		case err := <-cancel:
			vlog.INFO("%v", err)
			return
		}
	}
}
