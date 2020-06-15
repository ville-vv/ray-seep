package proxy

import (
	"encoding/json"
	"errors"
	"github.com/vilsongwei/vilgo/vlog"
	"ray-seep/ray-seep/common/conn"
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/msg"
	"ray-seep/ray-seep/proto"
	"ray-seep/ray-seep/server_v2/ifc"
)

type PxyManager struct {
	register ifc.Register
	router   msg.RouterFunc
}

func NewPxyManager(cft *conf.ControlSrv, reg ifc.Register) *PxyManager {
	return &PxyManager{register: reg}
}

func (p *PxyManager) OnConnect(cancel chan interface{}, cn conn.Conn) error {
	msgCtr := msg.NewMessageCenter(cn)
	var req msg.Package
	if err := msgCtr.Recv(&req); err != nil {
		return err
	}
	//
	if req.Cmd != msg.CmdRunProxyReq {
		_ = msgCtr.Send(&msg.Package{
			Cmd:  msg.CmdError,
			Body: []byte(""),
		})
		return errors.New("")
	}
	regData := proto.RunProxyReq{}
	if err := json.Unmarshal(req.Body, &regData); err != nil {
		vlog.ERROR("parse register proxy request data fail %s , data is %s ", err.Error(), string(req.Body))
		return err
	}

	if err := p.register.Register(regData.Name, regData.Cid, cn); err != nil {
		_ = msgCtr.Send(&msg.Package{Cmd: msg.CmdError, Body: []byte("")})
		return err
	}
	return nil
}

func (p *PxyManager) OnDisConnect(id int64) {
	return
}
