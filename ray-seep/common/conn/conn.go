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
	SetId(int64)
	IsClose() bool
}

type connUnit struct {
	net.Conn
	Cid     int64
	isClose bool
}

func TurnConn(c net.Conn) Conn {
	return newConnUnit(c, vuid.GenUUid())
}

func newConnUnit(c net.Conn, cid int64) *connUnit {
	return &connUnit{
		Conn: c,
		Cid:  cid,
	}
}

func (c *connUnit) Id() int64 {
	return c.Cid
}
func (c *connUnit) SetId(id int64) {
	c.Cid = id
}

func (c *connUnit) IsClose() bool {
	return c.isClose
}
func (c *connUnit) Close() error {
	c.isClose = true
	return c.Conn.Close()
}

// Listener 监听连接
type Listener struct {
	net.Addr
	Conn chan *connUnit
}

// 监听链接
func Listen(addr string) (l *Listener, err error) {
	lin, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}
	l = &Listener{
		Addr: lin.Addr(),
		Conn: make(chan *connUnit),
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
