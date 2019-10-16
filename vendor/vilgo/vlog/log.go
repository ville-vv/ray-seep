package vlog

var (
	log ILogger
)

const (
	LogLevelInfo   = "INFO"
	LogLevelDebug  = "DEBUG"
	LogLevelWarn   = "WARN"
	LogLevelError  = "ERROR"
	LogLevelNotice = "NOTICE"
)

var (
	ProgramName = "vlog"
)

type LogCnf struct {
	ProgramName   string
	OutPutFile    []string
	OutPutErrFile []string
	Level         string
}

// mo
//go:generate mockgen -destination mock/log_mock.go tstl/src/common/vlog ILogger

// logger 模块接口
type ILogger interface {
	LogD(string, ...interface{})
	LogE(string, ...interface{})
	LogI(string, ...interface{})
	LogW(string, ...interface{})
}

// 初始化默认日志模块
func DefaultLogger(lf ...string) ILogger {
	cnf := &LogCnf{
		OutPutErrFile: []string{},
		ProgramName:   ProgramName,
		OutPutFile:    []string{"stdout"},
		Level:         LogLevelDebug,
	}
	if len(lf) > 0 {
		cnf.OutPutFile = append(cnf.OutPutFile, lf[0])
	}
	log = NewGoLogger(cnf)
	return log
}

func SetLogger(logger ILogger) {
	log = logger
}

func LogD(format string, args ...interface{}) {
	log.LogD(format, args...)
}
func LogE(format string, args ...interface{}) {
	log.LogE(format, args...)
}
func LogI(format string, args ...interface{}) {
	log.LogI(format, args...)
}
func LogW(format string, args ...interface{}) {
	log.LogW(format, args...)
}


func DEBUG(format string, args ...interface{}) {
	log.LogD(format, args...)
}
func ERROR(format string, args ...interface{}) {
	log.LogE(format, args...)
}
func INFO(format string, args ...interface{}) {
	log.LogI(format, args...)
}
func WARN(format string, args ...interface{}) {
	log.LogW(format, args...)
}
