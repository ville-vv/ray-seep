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
	ErrNoticeProxyRunErr     ErrCode = "notice proxy run error "
	ErrWaitProxyRunTimeout   ErrCode = "wait proxy run timeout"
	ErrConnPoolIsFull        ErrCode = "proxy pool is full"
	ErrNoCmdRouterNot        ErrCode = "router not found"
	ErrClientControlNotExist ErrCode = "client control is not exist"
	ErrProxyConnNotExist     ErrCode = "proxy connect is not exist"
	ErrSecretIsInValid       ErrCode = "secret key is invalid"
)
