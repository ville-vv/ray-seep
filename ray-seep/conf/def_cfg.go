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

[Proto]
Proto="http"
Host = ""
Port = 32203
Domain = "example.com"

[Proxy]
Host=""
Port=32202

[DataBase]
[DataBase.Redis]
Address = "127.0.0.1:6379"    
Password =""    
UserName = ""     
MaxIdles = 100    
MaxOpens = 1000    
IdleTimeout = 10
IsMaxConnWait = false
Db=0

[DataBase.Mysql]
Version="8"
UserName="root"
Password="Root123"
Address="192.168.3.8:3306"
Default="information_schema"
Databases=["ray_seep"]

`
