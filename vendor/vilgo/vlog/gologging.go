package vlog

import (
	"github.com/op/go-logging"
	"io"
	"os"
)

// 日志级别
const (
	critical int = iota
	err_or
	warning
	notice
	info
	debug
)

var (
	// 日志输出格式
	logFormat = []string{
		`%{shortfunc} ▶ %{level:.4s} %{message}`,
		`%{color}%{time:15:04:05.00} %{shortfile} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
		`%{color}%{time:15:04:05.00} %{shortfunc} %{shortfile} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	}

	// 日志级别与 string 类型映射
	LogLevelMap = map[string]int{
		"CRITICAL":     critical,
		LogLevelError:  err_or,
		LogLevelWarn:   warning,
		LogLevelNotice: notice,
		LogLevelInfo:   info,
		LogLevelDebug:  debug,
	}
)

type GoLogger struct {
	log    *logging.Logger
	format int
	logCnf *LogCnf
}

func NewGoLogger(cnf *LogCnf) *GoLogger {
	gl := newLog(cnf)
	return gl
}

func newLog(cnf *LogCnf) *GoLogger {
	log := new(GoLogger)
	if cnf == nil {
		log.logCnf = &LogCnf{
			OutPutFile:    []string{"./log/" + ProgramName + "_log.log"},
			OutPutErrFile: []string{},
			Level:         LogLevelDebug,
			ProgramName:   ProgramName,
		}
	} else {
		log.logCnf = cnf
	}
	log.log = logging.MustGetLogger(log.logCnf.ProgramName)
	log.format = 2
	log.AddLogBackend()
	return log
}

// 添加日志输出终端，可以是文件，控制台，还有网络输出等。
func (l *GoLogger) AddLogBackend() {
	l.log.ExtraCalldepth = 2
	// 打开文件输出终端
	var backend []logging.Backend
	for _, v := range l.logCnf.OutPutFile {
		switch v {
		case "stdout":
			backend = append(backend, l.getStdOutBackend())
		default:
			backend = append(backend, l.getFileBackend(v))
		}
	}
	l.log.SetBackend(logging.SetBackend(backend...))
	return
}

// 文件输出终端
func (l *GoLogger) getFileBackend(filePath string) logging.LeveledBackend {

	// 创建目录
	CreateLogPath(filePath)

	// 打开一个文件
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	backend := l.getLogBackend(file, LogLevelMap[l.logCnf.Level])
	return backend
}

// 控制台输出终端
func (l *GoLogger) getStdOutBackend() logging.LeveledBackend {
	bked := l.getLogBackend(os.Stderr, LogLevelMap[l.logCnf.Level])
	return bked
}

func (l *GoLogger) getLogBackend(out io.Writer, level int) logging.LeveledBackend {
	backend := logging.NewLogBackend(out, "", 1)
	format := logging.MustStringFormatter(logFormat[1])
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backendFormatter)
	backendLeveled.SetLevel(logging.Level(level), "")
	return backendLeveled
}

func (l *GoLogger) LogI(infmt string, args ...interface{}) {
	l.log.Infof(infmt, args...)
	return
}

func (l *GoLogger) LogE(infmt string, args ...interface{}) {
	l.log.Errorf(infmt, args...)
	return
}

func (l *GoLogger) LogW(infmt string, args ...interface{}) {
	l.log.Warningf(infmt, args...)
	return
}

func (l *GoLogger) LogD(infmt string, args ...interface{}) {
	l.log.Debugf(infmt, args...)
	return
}
