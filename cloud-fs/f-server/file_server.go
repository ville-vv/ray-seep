package f_server

import (
	"github.com/vilsongwei/vilgo/vlog"
	"net"
	"net/http"
	f_system "ray-seep/cloud-fs/f-system"
)

type FileServer struct {
	serve http.Server
	fSys  *f_system.FileSystem
}

func NewFileServer(root string) *FileServer {
	return &FileServer{fSys: f_system.NewFileSystem(root)}
}

func (f *FileServer) Start(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		vlog.ERROR("file server start failure error is %s", err.Error())
		return
	}
	f.serve.Handler = f
	errCh := make(chan error)
	go func() {
		errS := f.serve.Serve(lis)
		errCh <- errS
	}()
	<-errCh
}

func (f *FileServer) ServeHTTP(rsp http.ResponseWriter, req *http.Request) {
	f.fSys.Display(rsp, req)
}
