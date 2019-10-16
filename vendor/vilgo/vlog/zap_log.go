package vlog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	infoLog *zap.Logger
	logCnf  *LogCnf
}

func NewZapLogger(cnf *LogCnf) *ZapLogger {
	l := new(ZapLogger)
	l.logCnf = cnf

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	switch l.logCnf.Level {
	case LogLevelError:
		atomicLevel.SetLevel(zap.ErrorLevel)
	case LogLevelInfo:
		atomicLevel.SetLevel(zap.InfoLevel)
	case LogLevelWarn:
		atomicLevel.SetLevel(zap.WarnLevel)
	default:
		atomicLevel.SetLevel(zap.DebugLevel)
	}

	//encoderConfig = zap.NewProductionEncoderConfig()
	//zapCore := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(os.Stdout), atomicLevel)
	CreateLogPath(cnf.OutPutFile...)
	CreateLogPath(cnf.OutPutErrFile...)
	zapConfig := zap.Config{
		Level:             atomicLevel,
		Development:       true,
		DisableCaller:     true,
		EncoderConfig:     encoderConfig,
		DisableStacktrace: true,
		Encoding:          "json",
		OutputPaths:       l.logCnf.OutPutFile[:],
		ErrorOutputPaths:  l.logCnf.OutPutErrFile[:],
	}

	logger, err := zapConfig.Build()
	if err != nil {
		panic(err)
	}
	//logger := zap.New(zapCore)
	l.infoLog = logger
	return l
}

func (l *ZapLogger) LogD(format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)
	l.infoLog.Debug("=>", zap.String("desc", str))
}
func (l *ZapLogger) LogE(format string, args ...interface{}) {
	if len(args) == 1{
		l.infoLog.Error("=>", zap.Any("desc", args[0]))
		return
	}
	str := fmt.Sprintf(format, args...)
	l.infoLog.Error("=>", zap.String("desc", str))
}

func (l *ZapLogger) LogW(format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)
	l.infoLog.Warn("=>", zap.String("desc", str))
}
func (l *ZapLogger) LogI(format string, args ...interface{}) {
	if len(args) == 1{
		l.infoLog.Info("=>", zap.Any("desc", args[0]))
		return
	}
	str := fmt.Sprintf(format, args...)
	l.infoLog.Info("=>", zap.String("desc", str))
}
