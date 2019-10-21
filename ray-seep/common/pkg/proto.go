// @File     : proto
// @Author   : Ville
// @Time     : 19-9-27 下午2:08
// pkg
package pkg

type IdentifyReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type IdentifyRsp struct {
	Id    int64
	Token string
}

type Ping struct{}
type Pong struct{}

type CreateHostReq struct {
	SubDomain string `json:"sub_domain"`
}
type CreateHostRsp struct {
	Domain string `json:"domain"`
}

type RegisterProxyReq struct {
	Cid       int64  `json:"cid"`
	SubDomain string `json:"sub_domain"`
}

type RegisterProxyRsp struct {
}
