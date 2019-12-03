// @File     : proto
// @Author   : Ville
// @Time     : 19-9-27 下午2:08
// pkg
package proto

type LoginReq struct {
	UserId int64  `json:"user_id"`
	Name   string `json:"name"`
	AppKey string `json:"app_key"`
}

type LoginRsp struct {
	Id    int64  `json:"id"`
	Token string `json:"token"`
}

type Ping struct{}
type Pong struct{}

type CreateHostReq struct {
	Token     string `json:"token"`
	SubDomain string `json:"sub_domain"`
}
type CreateHostRsp struct {
	ProxyHost string `json:"proxy_host"`
	Domain    string `json:"domain"`
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
