// @File     : def_cfg
// @Author   : Ville
// @Time     : 19-10-15 上午10:02
// conf
package conf

var clientDefaultConfig = `
[Control]
Host="193.112.47.13"
#Host=""
Port=4301
Name="rayseep"
UserId=123
Secret="4a35022cb0af2bc8471a1345d162575d"
AppKey="b753c6ad848e19ddd36c529430d262d5"

[Proxy]
Address = "193.112.47.13:4303"

[Web]
Address = ":9090"
`

//--------------------------------------------------------------------

var serverDefaultConfig = `
[Log]
ProgramName="ray_seep"
OutPutFile=["stdout","./log/rayseep.log"]
OutPutErrFile=[""]
Level="INFO"
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
Domain = "example.com"

[Proxy]
Host=""
Port=4303

[DataBase]
OpenRedis=false
OpenMysql=false

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
UserId=123
Secret="4a35022cb0af2bc8471a1345d162575d"
AppKey="b753c6ad848e19ddd36c529430d262d5"

`
