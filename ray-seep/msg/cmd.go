// @File     : const
// @Author   : Ville
// @Time     : 19-9-27 下午2:15
// pkg
package msg

const (
	maxBytesCachePool = 1024 * 1024 // 接收消息的最大缓存 1M
)

const (
	CmdError          int32 = 999999
	CmdPing           int32 = 100000
	CmdPong           int32 = 100001
	CmdLoginReq       int32 = 100002
	CmdLoginRsp       int32 = 100003
	CmdCreateHostReq  int32 = 100004 //
	CmdCreateHostRsp  int32 = 100005
	CmdRunProxyReq    int32 = 100006 // 注册代理请求
	CmdRunProxyRsp    int32 = 100007 // 注册代理返回
	CmdLogoutReq      int32 = 100038
	CmdLogoutRsp      int32 = 100039
	CmdNoticeRunProxy int32 = 200006 // 通知用户启动代理服务
)
