package log

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init_V1(serviceName, logKind string) {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		}, // ISO8601 UTC 时间格式
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
		EncodeCaller: zapcore.ShortCallerEncoder, // 不要全路径编码器
	}
	// 设置日志级别
	atom := zap.NewAtomicLevelAt(zap.DebugLevel)
	config := zap.Config{
		Level:            atom,                                               // 日志级别
		Development:      false,                                              // 开发模式，堆栈跟踪
		Encoding:         "json",                                             // 输出格式 console 或 json
		EncoderConfig:    encoderConfig,                                      // 编码器配置
		InitialFields:    map[string]interface{}{"serviceName": serviceName}, // 初始化字段，如：添加一个服务器名称
		OutputPaths:      []string{"./logs/" + logKind + ".log"},             // 输出到指定文件 stdout（标准输出，正常颜色） stderr（错误输出，红色）
		ErrorOutputPaths: []string{"./logs/stderr.log"},
	}
	// 构建日志
	logger, err := config.Build()
	if err != nil {
		panic(fmt.Sprintf("log 初始化失败: %v", err))
	}
	switch logKind {
	case "access":
		AccessLogger = logger
	case "debug":
		DebugLogger = logger
	case "info":
		InfoLogger = logger
	case "warnning":
		WarnningLogger = logger
	case "error":
		ErrorLogger = logger
	default:
	}
}
