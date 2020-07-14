// Package err_msg
package errs

type ErrCode string

func (s ErrCode) Error() string {
	return string(s)
}
func (s ErrCode) String() string {
	return string(s)
}

const (
	StatusOK                 ErrCode = "200"
	ErrServerNumFull         ErrCode = "you server number is full "
	ErrServerNotExist        ErrCode = "server not exist"
	ErrProxySrvNotExist      ErrCode = "proxy server not exist"
	ErrNoticeProxyRunErr     ErrCode = "notice proxy apps error "
	ErrWaitProxyRunTimeout   ErrCode = "wait proxy apps timeout"
	ErrConnPoolIsFull        ErrCode = "proxy pool is full"
	ErrNoCmdRouterNot        ErrCode = "router not found"
	ErrClientControlNotExist ErrCode = "client node is not exist"
	ErrProxyConnNotExist     ErrCode = "proxy connect is not exist"
	ErrProxyWaitCacheErr     ErrCode = "proxy wait caches err"
	ErrProxyHaveRegister     ErrCode = "proxy have be registered"
	ErrUserInfoValidFail     ErrCode = "store information validation error"
	ErrHttpPortIsInValid     ErrCode = "http port is invalid"
	ErrNoThisUser            ErrCode = "no this store"
	ErrNoLogin               ErrCode = "store have not login"
	ErrUserHaveLogin         ErrCode = "user have login"
)
