// @File     : const
// @Author   : Ville
// @Time     : 19-9-27 下午2:15
// pkg
package proto

const (
	CmdError            int32 = 999999
	CmdPing             int32 = 100000
	CmdPong             int32 = 100001
	CmdLoginReq         int32 = 100002
	CmdLoginRsp         int32 = 100003
	CmdCreateHostReq    int32 = 100004 //
	CmdCreateHostRsp    int32 = 100005
	CmdRegisterProxyReq int32 = 100006 // 注册代理请求
	CmdRegisterProxyRsp int32 = 100007 // 注册代理返回
	CmdLogoutReq        int32 = 100038
	CmdLogoutRsp        int32 = 100039
)
