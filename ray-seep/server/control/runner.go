package control

import (
	"ray-seep/ray-seep/common/repeat"
	"ray-seep/ray-seep/server/http"
	"sync"
)

type IRunner interface {
	Start() error
	Stop()
}

type JoinItem struct {
	Name   string
	ConnId int64
	Addr   string
	Err    chan error
}

type LeaveItem struct {
	Name   string
	ConnId string
}

type Runner struct {
	join   chan JoinItem
	leave  chan LeaveItem
	items  map[string]IRunner
	gainer repeat.NetConnGainer
}

func NewRunner() *Runner {
	return &Runner{
		join:  make(chan JoinItem, 100),
		leave: make(chan LeaveItem, 100),
		items: make(map[string]IRunner),
	}
}

func (sel *Runner) SetGainer(gainer repeat.NetConnGainer) {
	sel.gainer = gainer
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
	run := http.NewServerWithAddr(item.Addr, sel.gainer)
	errCh := make(chan error)
	go func() {
		errCh <- run.Start()
	}()
	err := <-errCh
	if err == nil {
		sel.items[item.Name] = run
	}
	item.Err <- err
	return
}

func (sel *Runner) delItem(item *LeaveItem) {
	if pxy, ok := sel.items[item.Name]; ok {
		pxy.Stop()
		delete(sel.items, item.Name)
	}
}

func (sel *Runner) Close() {
	close(sel.join)
	close(sel.leave)
}
