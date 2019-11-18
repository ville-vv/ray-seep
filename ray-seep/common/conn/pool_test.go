package conn

import (
	"fmt"
	"time"

	//"fmt"
	"testing"
	"vilgo/vlog"
)

type MockConn struct {
	Conn
}

func (m *MockConn) Close() error {
	return nil
}

func TestPool_Expire(t *testing.T) {
	vlog.DefaultLogger()
	pl := NewPool(2)
	_ = pl.Push(int64(132432), &MockConn{})
	_ = pl.Push(int64(132433), &MockConn{})
	_ = pl.Push(int64(132434), &MockConn{})
	_ = pl.Push(int64(132435), &MockConn{})
	_ = pl.Push(int64(132436), &MockConn{})
	num := pl.Size()
	time.Sleep(time.Second * 3)
	fmt.Println("开始")
	for {
		if num <= 0 {
			return
		}
		select {
		case id := <-pl.Expire():
			fmt.Println("到期:", id)
			num--
		}
	}
}
