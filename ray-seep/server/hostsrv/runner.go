// 一个启动器，用于启动外部加入的服务，外部服务需实现接口 IRunner
// 加入服务使用 Join 函数， 销毁使用 Leave 函数

package hostsrv

import (
	"github.com/vilsongwei/vilgo/vlog"
	"sync"
)

type IRunner interface {
	Id() int64
	Start() error
	Stop()
}

type JoinItem struct {
	Name   string
	ConnId int64
	Run    IRunner
	Err    chan error
}

type LeaveItem struct {
	Name   string
	ConnId int64
}

type RunnerMng struct {
	join  chan JoinItem
	leave chan LeaveItem
	items map[string]IRunner
}

func NewRunnerMng() *RunnerMng {
	return &RunnerMng{
		join:  make(chan JoinItem, 100),
		leave: make(chan LeaveItem, 100),
		items: make(map[string]IRunner),
	}
}

func (sel *RunnerMng) Join() chan<- JoinItem {
	return sel.join
}

func (sel *RunnerMng) Leave() chan<- LeaveItem {
	return sel.leave
}

func (sel *RunnerMng) Start() error {
	w := sync.WaitGroup{}
	w.Add(1)
	go func() {
		w.Done()
		for {
			select {
			case jv, ok := <-sel.join:
				if !ok {
					return
				}
				sel.addItem(&jv)
			case lv, ok := <-sel.leave:
				if !ok {
					return
				}
				sel.delItem(&lv)
			}
		}
	}()
	w.Wait()
	return nil
}

func (sel *RunnerMng) addItem(item *JoinItem) {
	if _, ok := sel.items[item.Name]; ok {
		item.Err <- nil
		return
	}
	errCh := make(chan error)
	go func() {
		errCh <- item.Run.Start()
	}()
	err := <-errCh
	if err == nil {
		// vlog.DEBUG("启动服务：%s", item.Name)
		sel.items[item.Name] = item.Run
	}
	item.Err <- err
	return
}

func (sel *RunnerMng) delItem(item *LeaveItem) {
	if pxy, ok := sel.items[item.Name]; ok {
		if pxy.Id() == item.ConnId {
			vlog.DEBUG("delete the runner [%s]", item.Name)
			pxy.Stop()
			delete(sel.items, item.Name)
		}
	}
}

func (sel *RunnerMng) Close() {
	close(sel.join)
	close(sel.leave)
}
