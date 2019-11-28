// @File     : proto
// @Author   : Ville
// @Time     : 19-9-27 下午2:08
// pkg
package proto

type LoginReq struct {
	UserId int64  `json:"user_id"`
	Name   string `json:"name"`
	AppId  string `json:"app_id"`
}

type LoginRsp struct {
	Id    int64
	Token string
}

type Ping struct{}
type Pong struct{}

type CreateHostReq struct {
	Token     string `json:"token"`
	SubDomain string `json:"sub_domain"`
}
type CreateHostRsp struct {
	Domain string `json:"domain"`
}

type NoticeRunProxy struct {
	Cid  int64  `json:"cid"`
	Name string `json:"name"`
}

type RunProxyReq struct {
	Cid   int64  `json:"cid"`
	Token string `json:"token"`
	Name  string `json:"name"`
}

type RunProxyRsp struct {
	Cid int64 `json:"cid"`
}
