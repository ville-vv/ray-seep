// @File     : proxy
// @Author   : Ville
// @Time     : 19-9-27 下午4:59
// proxy
package proxy

type Server struct {
	addr string
}

func NewServer() *Server {
	return &Server{
		addr: ":39990",
	}
}

func Start() {
}
