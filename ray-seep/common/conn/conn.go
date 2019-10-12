// @File     : conn
// @Author   : Ville
// @Time     : 19-9-24 上午9:39
// manager
package conn

import (
	"net"
	"vilgo/vuid"
)

type Conn interface {
	net.Conn
	Id() int64
}

type ConnUnit struct {
	net.Conn
	Cid int64
}

func TurnConn(c net.Conn) Conn {
	return newConnUnit(c, vuid.GenUUid())
}

func newConnUnit(c net.Conn, cid int64) *ConnUnit {
	return &ConnUnit{
		Conn: c,
		Cid:  cid,
	}
}

func (c *ConnUnit) Id() int64 {
	return c.Cid
}

// Listener 监听连接
type Listener struct {
	net.Addr
	Conn chan *ConnUnit
}

// 监听链接
func Listen(addr string) (l *Listener, err error) {
	lin, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}
	l = &Listener{
		Addr: lin.Addr(),
		Conn: make(chan *ConnUnit),
	}
	go func() {
		for {
			cn, err := lin.Accept()
			if err != nil {
				continue
			}
			l.Conn <- newConnUnit(cn, vuid.GenUUid())
		}
	}()
	return
}
