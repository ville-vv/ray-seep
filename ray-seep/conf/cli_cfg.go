// @File     : cli_cfg
// @Author   : Ville
// @Time     : 19-10-14 下午3:14
// config
package conf

import (
	"github.com/vilsongwei/vilgo/vcnf"
	"os"
)

type Client struct {
	Pxy     *ProxyCli     `json:"proxy" toml:"Proxy"`
	Control *ControlCli   `json:"control" toml:"Control"`
	Web     *WebServerCli `json:"web" toml:"Web"`
}

type WebServerCli struct {
	Address string `json:"address"`
}

type ProxyCli struct {
	Host string `json:"Host"`
}

type ControlCli struct {
	Host              string `json:"host"`
	Port              int64  `json:"port"`
	Name              string `json:"name" toml:"Name"` // 子域名
	UserId            int64  `json:"user_id"`
	Secret            string `json:"secret" toml:"Secret"`
	AppKey            string `json:"app_key" toml:"AppKey"`
	HttpPort          string `json:"http_port" toml:"HttpPort"`
	CanReconnect      bool   `json:"can_reconnect"`      // 是否能够重连
	ReconnectEndTime  int64  `json:"reconnect_end_time"` // 重连结束时间
	ReconnectInternal int64  `json:"reconnect_internal"` // 重连一次的间隔时间
}

// 初始化客户端配置
func InitClient(fileName ...string) *Client {
	fn := ""
	if len(fileName) > 0 {
		fn = fileName[0]
	}
	cliCnf := &Client{}
	reader := vcnf.NewReader(fn)
	if defReader, ok := reader.(*vcnf.DefaultReader); ok {
		defReader.SetInfo(clientDefaultConfig, "toml")
	}
	if err := reader.CnfRead(cliCnf); err != nil {
		panic(err)
	}
	return cliCnf
}

func GenDefClientConfigFile(fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	_, err = file.Write([]byte(clientDefaultConfig))
	if err != nil {
		panic(err)
	}
}
