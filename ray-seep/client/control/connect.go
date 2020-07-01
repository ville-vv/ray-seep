package control

import (
	"github.com/vilsongwei/vilgo/vlog"
	"net"
)

type ClientConnect struct {
	addr   string
	connCh chan net.Conn
	errCh  chan error
}

func (sel *ClientConnect) ConnCh() <-chan net.Conn {
	return sel.connCh
}

func (sel *ClientConnect) Connect() {
	c, err := net.Dial("tcp", sel.addr)
	if err != nil {
		vlog.LogE("connect node server fail %v", err)
		sel.errCh <- err
		return
	}
	sel.connCh <- c
}

func (sel *ClientConnect) Reconnect() {
}
