// @File     : srv_cfg
// @Author   : Ville
// @Time     : 19-10-14 下午3:14
// config
package conf

import (
	"github.com/vilsongwei/vilgo/vcnf"
	"github.com/vilsongwei/vilgo/vlog"
	"github.com/vilsongwei/vilgo/vredis"
	"github.com/vilsongwei/vilgo/vsql"
	"os"
)

type Server struct {
	Host     string           `json:"host"`
	Domain   string           `json:"domain"` //服务的域名
	Log      *vlog.LogCnf     `json:"log"`
	Ctl      *SubServer       `json:"ctl" toml:"Control"`
	Pxy      *SubServer       `json:"ctl" toml:"Proxy"`
	DataBase *DataBaseSrv     `json:"database" toml:"DataBase"`
	User     map[string]*User `json:"store" toml:"User"`
}

// SubServer 用户控制器，用于与用户客户端消息通信，和控制命令的处理
type SubServer struct {
	Host            string `json:"host"`
	Port            int64  `json:"port"`
	ReadMsgTimeout  int64  `json:"read_msg_timeout"`   // 连接的超时时间毫秒,如果在时间内没有收到任何消息就会断开
	MaxUserNum      int    `json:"max_user_num"`       // 最大客户端连接数
	UserMaxProxyNum int    `json:"user_max_proxy_num"` // 一个用户下面允许连接的最大代理pod服务数量
}

type DataBaseSrv struct {
	OpenDb bool             `json:"open_db"` //是否打开Databases服务
	Redis  *vredis.RedisCnf `json:"redis" toml:"Redis"`
	Mysql  *vsql.MySqlCnf   `json:"mysql" toml:"Mysql"`
}

type User struct {
	UserId   int64  `json:"user_id" toml:"UserId"`
	UserName string `json:"user_name" toml:"UserName"`
	Secret   string `json:"secret" toml:"Secret"`
	AppKey   string `json:"app_key" toml:"AppKey"`
	HttpPort string `json:"http_port" toml:"HttpPort"`
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

func GenDefServerConfigFile(fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	_, err = file.Write([]byte(serverDefaultConfig))
	if err != nil {
		panic(err)
	}
}
