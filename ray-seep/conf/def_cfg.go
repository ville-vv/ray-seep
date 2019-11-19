// @File     : def_cfg
// @Author   : Ville
// @Time     : 19-10-15 上午10:02
// conf
package conf

var clientDefaultConfig = `
[Control]
Host=""
Port=32201
Domain="exampletest.cn"
SubDomain="test"

[Proxy]
Host=""
Port=32202
`

//--------------------------------------------------------------------

var serverDefaultConfig = `
[Control]
Host = ""
Port = 32201
# 连接的超时时间 单位/毫秒
Timeout=10000
MaxClientNumber = 5
MaxClientProxyNumber = 5

[Http]
Host = ""
Port = 32203
Domain = "example.com"

[Proxy]
Host=""
Port=32202
`
