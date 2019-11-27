// @File     : def_cfg
// @Author   : Ville
// @Time     : 19-10-15 上午10:02
// conf
package conf

var clientDefaultConfig = `
[Control]
Host=""
Port=32201
Name="test"

[Proxy]
Address = ":32202"

[Web]
Address = ":9090"
`

//--------------------------------------------------------------------

var serverDefaultConfig = `
[Control]
Host = ""
Port = 32201
# 连接的超时时间 单位/毫秒
ReadMsgTimeout=10000
#最大客户端连接数
MaxUserNum = 5
# 一个客户端最大的代理数
UserMaxProxyNum = 10

[Http]
Host = ""
Port = 32203
Domain = "example.com"

[Proxy]
Host=""
Port=32202
`
