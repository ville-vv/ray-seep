// @File     : cli_cfg
// @Author   : Ville
// @Time     : 19-10-14 下午3:14
// config
package conf

import "vilgo/vcnf"

type Client struct {
	Pxy  *ProxyCli `json:"pxy" toml:"Proxy"`
	Node *NodeCli  `json:"node"`
}

// 节点信息配置， 一个用户一个节点，一个 node 可以有多个 Pod, Pod 想到
type NodeCli struct {
	Host      string `json:"host"`
	Port      int64  `json:"port"`
	Domain    string `json:"domain"`     // 域名
	SubDomain string `json:"sub_domain"` // 子域名
}

type ProxyCli struct {
	Host string `json:"host"`
	Port int64  `json:"port"`
}

type ControlCli struct {
	Host string `json:"host"`
	Port int64  `json:"port"`
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
