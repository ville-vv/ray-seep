package control

import (
	"sync"
	"github.com/vilsongwei/vilgo/vlog"
)

type IRunner interface {
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
	ConnId string
}

type Runner struct {
	join  chan JoinItem
	leave chan LeaveItem
	items map[string]IRunner
}

func NewRunner() *Runner {
	return &Runner{
		join:  make(chan JoinItem, 100),
		leave: make(chan LeaveItem, 100),
		items: make(map[string]IRunner),
	}
}

func (sel *Runner) Join() chan<- JoinItem {
	return sel.join
}

func (sel *Runner) Leave() chan<- LeaveItem {
	return sel.leave
}

func (sel *Runner) Start() error {
	w := sync.WaitGroup{}
	w.Add(2)
	go func() {
		w.Done()
		for v := range sel.leave {
			sel.delItem(&v)
		}
	}()

	go func() {
		w.Done()
		for v := range sel.join {
			sel.addItem(&v)
		}
	}()
	w.Wait()
	return nil
}

func (sel *Runner) addItem(item *JoinItem) {
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
		vlog.DEBUG("启动服务：%s", item.Name)
		sel.items[item.Name] = item.Run
	}
	item.Err <- err
	return
}

func (sel *Runner) delItem(item *LeaveItem) {
	if pxy, ok := sel.items[item.Name]; ok {
		vlog.DEBUG("清理服务：%s", item.Name)
		pxy.Stop()
		delete(sel.items, item.Name)
	}
}

func (sel *Runner) Close() {
	close(sel.join)
	close(sel.leave)
}
