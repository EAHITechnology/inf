package log

import (
	"fmt"
	"io"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/net/context"
)

var myLogger *zap.Logger

func Debug(msg string) {
	myLogger.Debug(msg)
}

/*
直接拼接字符串,如果缓冲变得太大，Write会采用错误值ErrTooLarge引发panic
*/
func write(strs ...string) (string, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Write 捕获到的错误：%s\n", r)
		}
	}()
	var strBuild strings.Builder
	for _, str := range strs {
		strBuild.WriteString(str)
	}
	return strBuild.String(), nil
}

func DebugCtx(ctx context.Context, msg string) {
	if ctx.Value(LOGID) != nil {
		msg, _ = write(msg, ",", LOGID, ":"+ctx.Value(LOGID).(string))
	}
	myLogger.Debug(msg)
}

func Debugf(template string, args ...interface{}) {
	msg := template
	if msg == "" && len(args) > 0 {
		msg = fmt.Sprint(args...)
	} else if msg != "" && len(args) > 0 {
		msg = fmt.Sprintf(template, args...)
	}
	myLogger.Debug(msg)
}

func DebugfCtx(ctx context.Context, template string, args ...interface{}) {
	msg := template
	if msg == "" && len(args) > 0 {
		msg = fmt.Sprint(args...)
	} else if msg != "" && len(args) > 0 {
		msg = fmt.Sprintf(template, args...)
	}
	if ctx.Value(LOGID) != nil {
		msg = msg + "," + LOGID + ":" + ctx.Value(LOGID).(string)
	}
	myLogger.Debug(msg)
}

func Info(msg string) {
	myLogger.Info(msg)
}

func InfoCtx(ctx context.Context, msg string) {
	if ctx.Value(LOGID) != nil {
		msg = msg + "," + LOGID + ":" + ctx.Value(LOGID).(string)
	}
	myLogger.Info(msg)
}

func Infof(template string, args ...interface{}) {
	msg := template
	if msg == "" && len(args) > 0 {
		msg = fmt.Sprint(args...)
	} else if msg != "" && len(args) > 0 {
		msg = fmt.Sprintf(template, args...)
	}
	myLogger.Info(msg)
}

func InfofCtx(ctx context.Context, template string, args ...interface{}) {
	msg := template
	if msg == "" && len(args) > 0 {
		msg = fmt.Sprint(args...)
	} else if msg != "" && len(args) > 0 {
		msg = fmt.Sprintf(template, args...)
	}
	if ctx.Value(LOGID) != nil {
		msg = msg + "," + LOGID + ":" + ctx.Value(LOGID).(string)
	}
	myLogger.Info(msg)
}

func Warn(msg string) {
	myLogger.Warn(msg)
}

func WarnCtx(ctx context.Context, msg string) {
	if ctx.Value(LOGID) != nil {
		msg = msg + "," + LOGID + ":" + ctx.Value(LOGID).(string)
	}
	myLogger.Warn(msg)
}

func Warnf(template string, args ...interface{}) {
	msg := template
	if msg == "" && len(args) > 0 {
		msg = fmt.Sprint(args...)
	} else if msg != "" && len(args) > 0 {
		msg = fmt.Sprintf(template, args...)
	}
	myLogger.Warn(msg)
}

func WarnfCtx(ctx context.Context, template string, args ...interface{}) {
	msg := template
	if msg == "" && len(args) > 0 {
		msg = fmt.Sprint(args...)
	} else if msg != "" && len(args) > 0 {
		msg = fmt.Sprintf(template, args...)
	}
	if ctx.Value(LOGID) != nil {
		msg = msg + "," + LOGID + ":" + ctx.Value(LOGID).(string)
	}
	myLogger.Warn(msg)
}

func Error(msg string) {
	myLogger.Error(msg)
}

func ErrorCtx(ctx context.Context, msg string) {
	if ctx.Value(LOGID) != nil {
		msg = msg + "," + LOGID + ":" + ctx.Value(LOGID).(string)
	}
	myLogger.Error(msg)
}

func Errorf(template string, args ...interface{}) {
	msg := template
	if msg == "" && len(args) > 0 {
		msg = fmt.Sprint(args...)
	} else if msg != "" && len(args) > 0 {
		msg = fmt.Sprintf(template, args...)
	}
	myLogger.Error(msg)
}

func ErrorfCtx(ctx context.Context, template string, args ...interface{}) {
	msg := template
	if msg == "" && len(args) > 0 {
		msg = fmt.Sprint(args...)
	} else if msg != "" && len(args) > 0 {
		msg = fmt.Sprintf(template, args...)
	}
	if ctx.Value(LOGID) != nil {
		msg = msg + "," + LOGID + ":" + ctx.Value(LOGID).(string)
	}
	myLogger.Error(msg)
}

func DPanic(msg string) {
	myLogger.DPanic(msg)
}

func DPanicf(template string, args ...interface{}) {
	msg := template
	if msg == "" && len(args) > 0 {
		msg = fmt.Sprint(args...)
	} else if msg != "" && len(args) > 0 {
		msg = fmt.Sprintf(template, args...)
	}
	myLogger.DPanic(msg)
}

func getWriter(filename string) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 xxx.log.YYmmddHH
	// xxxx.log是指向最新日志的链接
	// 保存X天内的日志，每1小时(整点)分割一次日志
	hook, err := rotatelogs.New(
		strings.Replace(filename, ".log", "", -1)+"-%Y%m%d%H.log",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*30),
		rotatelogs.WithRotationTime(time.Hour),
	)
	if err != nil {
		panic(err)
	}
	return hook
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func Init(serviceName, path string) {
	encoderConfig := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		NameKey:        "Name",
		LevelKey:       "level",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "trace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 不要全路径编码器
	})
	// 实现两个判断日志等级的interface
	debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel && lvl < zapcore.InfoLevel
	})
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel && lvl < zapcore.WarnLevel
	})
	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel && lvl < zapcore.ErrorLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	debugWriter := getWriter(path + "/debug.log")
	infoWriter := getWriter(path + "/info.log")
	warnWriter := getWriter(path + "/warning.log")
	errorWriter := getWriter(path + "/error.log")
	// 最后创建具体的Logger
	core := zapcore.NewTee(
		zapcore.NewCore(encoderConfig, zapcore.AddSync(debugWriter), debugLevel),
		zapcore.NewCore(encoderConfig, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoderConfig, zapcore.AddSync(warnWriter), warnLevel),
		zapcore.NewCore(encoderConfig, zapcore.AddSync(errorWriter), errorLevel),
	)
	myLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	myLogger.Sync()
}
