package hostsrv

import (
	"ray-seep/ray-seep/common/repeat"
)

type HostServer interface {
	Start() error
	Stop()
	Create(id int64, kind, addr string) error
	Destroy(id int64, addr string)
}

type NetConnGainer interface {
	GainConn()
}

type HostService struct {
	runner  *Runner
	dstConn repeat.NetConnGainer
}

func NewHostService() *HostService {
	return &HostService{runner: NewRunner()}
}

func (h *HostService) Start() error {
	return h.runner.Start()
}

func (h *HostService) Stop() {
	h.runner.Close()
}

func (h *HostService) SetDstConn(dstConn repeat.NetConnGainer) {
	h.dstConn = dstConn
}

func (h *HostService) Create(id int64, kind, addr string) error {
	join := JoinItem{
		Name:   addr,
		ConnId: id,
		Err:    make(chan error),
	}
	join.Run = NewServerWithAddr(addr, h.dstConn)
	h.runner.Join() <- join
	err := <-join.Err
	return err
}

func (h *HostService) Destroy(id int64, addr string) {
	h.runner.Leave() <- LeaveItem{
		Name:   addr,
		ConnId: id,
	}
	return
}
