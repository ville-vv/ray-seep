// @File     : package
// @Author   : Ville
// @Time     : 19-9-24 下午3:01
// msg
package msg

type LoginReq struct {
	Customer string `json:"customer"`
	Password string `json:"password"`
}

type LoginRsp struct {
	Id int64
}

type Ping struct{}
type Pong struct{}
