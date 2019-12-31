// @File     : def_cfg
// @Author   : Ville
// @Time     : 19-10-15 上午10:02
// conf
package conf

var clientDefaultConfig = `
[Control]
Host="127.0.0.1"
Port=4301
Name="test"
UserId=100
Secret="example"
AppKey="example"
#是否打开重连
CanReconnect=false
#重连结束时间/秒
ReconnectEndTime=60
#重连一次的间隔时间/秒
ReconnectInternal=3

[Proxy]
# 这里只需要指定host不用指定port, host暂时与Control配置共用
#Host = "127.0.0.1"

[Web]
Address = ":12345"
`

//--------------------------------------------------------------------

var serverDefaultConfig = `
[Log]
ProgramName="ray_seep"
OutPutFile=["stdout"]
OutPutErrFile=[""]
Level="DEBUG"

[Control]
Host = ""
Port = 4301
# 连接的超时时间 单位/毫秒
ReadMsgTimeout=10000
#最大客户端连接数
MaxUserNum = 5
# 一个客户端最大的代理数
UserMaxProxyNum = 10

[Proto]
Proto="http"
Host = ""
Port = 4302
Domain = "rayseep.example.com"

[Proxy]
Host=""
Port=4303

[DataBase]
OpenDb=false

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
Address="127.0.0.1:3306"
Default="information_schema"
Databases=["ray_seep"]

[User]
[User.test]
UserName="test"
UserId=100
Secret="example"
AppKey="example"
HttpPort="4900"
`
