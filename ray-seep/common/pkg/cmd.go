// @File     : const
// @Author   : Ville
// @Time     : 19-9-27 下午2:15
// pkg
package pkg

type Command int32

const (
	CmdError            Command = 999999
	CmdPing             Command = 100000
	CmdPong             Command = 100001
	CmdIdentifyReq      Command = 100002
	CmdIdentifyRsp      Command = 100003
	CmdCreateHostReq    Command = 100004 //
	CmdCreateHostRsp    Command = 100005
	CmdRegisterProxyReq Command = 100006 // 注册代理请求
	CmdRegisterProxyRsp Command = 100007 // 注册代理返回
)
