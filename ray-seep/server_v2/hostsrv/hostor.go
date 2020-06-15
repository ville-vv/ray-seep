package hostsrv

import (
	"ray-seep/ray-seep/common/repeat"
	"ray-seep/ray-seep/msg"
)

type Option struct {
	Id     int64              // ID号
	Kind   string             // 类型
	SendCh chan<- msg.Package //
	Addr   string
}

type HostServer interface {
	Create(opt *Option) error
	Destroy(opt *Option)
}

type NetConnGainer interface {
	GainConn()
}

type HostService struct {
	runner  *Runner
	dstConn repeat.NetConnGainer
}

func NewHostService(runner *Runner) *HostService {
	return &HostService{runner: runner}
}

func (h *HostService) Start() error {
	return h.runner.Start()
}

func (h *HostService) SetDstConn(dstConn repeat.NetConnGainer) {
	h.dstConn = dstConn
}

func (h *HostService) Create(opt *Option) error {
	join := JoinItem{
		Name:   opt.Addr,
		ConnId: opt.Id,
		Err:    make(chan error),
	}
	join.Run = NewServerWithAddr(opt.Addr, h.dstConn)
	h.runner.Join() <- join
	return <-join.Err
}

func (h *HostService) Destroy(opt *Option) {
	return
}
