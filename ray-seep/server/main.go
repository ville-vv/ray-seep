// @File     : main
// @Author   : Ville
// @Time     : 19-9-24 下午3:10 
// server 
package server

import "vilgo/vlog"

type Server struct {

}

func Start(){
	vlog.DefaultLogger()
	control := NewControlServer()
	control.Start()
}
