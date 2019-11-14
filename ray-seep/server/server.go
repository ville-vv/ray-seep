// @File     : main
// @Author   : Ville
// @Time     : 19-9-24 下午3:10
// server
package server

import (
	"ray-seep/ray-seep/conf"
	"vilgo/vlog"
)

type Server interface {
	Start() error
	Stop()
	Scheme() string
}

type RaySeepServer struct {
	srvCnf *conf.Server
	srvs   map[string]Server
	start  []string
	stopCh chan int
}

func NewRaySeepServer(srvCnf *conf.Server) *RaySeepServer {
	return &RaySeepServer{srvCnf: srvCnf, srvs: make(map[string]Server), stopCh: make(chan int, 1)}
}

func (r *RaySeepServer) Start() {
	for k, v := range r.srvs {
		go func(sv Server) {
			defer func() {
				if err := recover(); err != nil {
					r.Stop()
					panic(err)
				}
			}()
			// 记录启动了的服务
			r.start = append(r.start, k)
			vlog.INFO("server [%s] starting", sv.Scheme())
			if err := sv.Start(); err != nil {
				panic(err)
			}
		}(v)
	}
	<-r.stopCh
	vlog.INFO("server [all] have stop success")
}

func (r *RaySeepServer) Stop() {
	// 停止已启动的服务
	for i := range r.start {
		r.srvs[r.start[i]].Stop()
	}
	close(r.stopCh)
}

func (r *RaySeepServer) Use(s ...Server) {
	for i := range s {
		r.srvs[s[i].Scheme()] = s[i]
	}
}
