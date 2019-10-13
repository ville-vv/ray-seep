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
	StatusOK          ErrCode = "200"
	ErrServerNumFull  ErrCode = "you server number is full "
	ErrServerNotExist ErrCode = "server not exist"
	ErrProxyNotExist  ErrCode = "proxy server not exist"
)
