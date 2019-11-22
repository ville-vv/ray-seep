package proxy

import "net"

type Tunnel struct {
	pAddr string // 目标地址
}

func NewTunnel(pAddr string) *Tunnel {
	return &Tunnel{pAddr: pAddr}
}

func (sel *Tunnel) GetNetConn(name string) (net.Conn, error) {
	cn, err := net.Dial("tcp", sel.pAddr)
	if err != nil {
		return nil, err
	}
	return cn, nil
}
