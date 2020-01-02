package static

import (
	"net"
	"net/http"
)

type FileServer struct {
	serve http.Server
	fSys  *FileSystem
}

func NewFileServer(root string) *FileServer {
	return &FileServer{fSys: NewFileSystem(root)}
}

func (f *FileServer) Start(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
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
