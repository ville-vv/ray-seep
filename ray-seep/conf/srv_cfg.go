// @File     : srv_cfg
// @Author   : Ville
// @Time     : 19-10-14 下午3:14
// config
package conf

import "vilgo/vcnf"

type Server struct {
	Ctl  *ControlSrv `json:"ctl" toml:"Control"`
	Http *HttpSrv    `json:"http" toml:"Http"`
	Pxy  *ProxySrv   `json:"pxy" toml:"Proxy"`
}

// ProxySrv 代理服务， 用户建立客户端连接后，需要建立代理的连接
type ProxySrv struct {
	Host string `json:"host"`
	Port int64  `json:"port"`
}

// HttpSrv 服务程序对外的 http 服务信息
type HttpSrv struct {
	Host   string `json:"host"`
	Port   int64  `json:"port"`
	Domain string `json:"domain"` //服务的域名
}

// ControlSrv 用户控制器，用于与用户客户端消息通信，和控制命令的处理
type ControlSrv struct {
	Host            string `json:"host"`
	Port            int64  `json:"port"`
	ReadMsgTimeout  int64  `json:"read_msg_timeout"`   // 连接的超时时间毫秒
	MaxUserNum      int    `json:"max_user_num"`       // 最大客户端连接数
	UserMaxProxyNum int    `json:"user_max_proxy_num"` // 一个客户的最大代理数
}

//--------------------------------------------------------------------
// 初始化服务端配置
func InitServer(fileName ...string) *Server {
	fn := ""
	if len(fileName) > 0 {
		fn = fileName[0]
	}
	srvCnf := &Server{}
	reader := vcnf.NewReader(fn)
	if defReader, ok := reader.(*vcnf.DefaultReader); ok {
		defReader.SetInfo(serverDefaultConfig, "toml")
	}
	if err := reader.CnfRead(srvCnf); err != nil {
		panic(err)
	}
	return srvCnf
}
