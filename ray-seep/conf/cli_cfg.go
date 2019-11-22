// @File     : cli_cfg
// @Author   : Ville
// @Time     : 19-10-14 下午3:14
// config
package conf

import "vilgo/vcnf"

type Client struct {
	Pxy     *ProxyCli     `json:"proxy" toml:"Proxy"`
	Control *ControlCli   `json:"control" toml:"Control"`
	Web     *WebServerCli `json:"web" toml:"Web"`
}

type WebServerCli struct {
	Address string `json:"address"`
}

type ProxyCli struct {
	Address string `json:"address"`
}

type ControlCli struct {
	Host string `json:"host"`
	Port int64  `json:"port"`
	Name string `json:"name" toml:"Name"` // 子域名
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
