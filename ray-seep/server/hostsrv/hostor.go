package hostsrv

import (
	"fmt"
	"ray-seep/ray-seep/common/repeat"
)

type HostServer interface {
	Start() error
	Stop()
	Create(id int64, kind, addr string) error
	Destroy(id int64, addr string)
}

type HostService struct {
	runner  *RunnerMng
	dstConn repeat.NetConnGainer
}

func NewHostService() *HostService {
	return &HostService{runner: NewRunnerMng()}
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
	run, err := RunnerFactory(id, kind, addr, h.dstConn)
	if err != nil {
		return err
	}
	join := JoinItem{
		Name:   addr,
		ConnId: id,
		Err:    make(chan error),
		Run:    run,
	}
	h.runner.Join() <- join
	err = <-join.Err
	return err
}

func (h *HostService) Destroy(id int64, addr string) {
	h.runner.Leave() <- LeaveItem{
		Name:   addr,
		ConnId: id,
	}
	return
}

// 对外代理类型工厂， kind 可以为 http, tcp, ssh, 等等
func RunnerFactory(id int64, kind, addr string, pxyGainer repeat.NetConnGainer) (IRunner, error) {
	// 可以根据 kind 不同类型启动 不同的服务
	switch kind {
	case "http":
		return newHttpRunner(id, addr, pxyGainer), nil
	default:
		return nil, fmt.Errorf("server kind of %s not support", kind)
	}
}
