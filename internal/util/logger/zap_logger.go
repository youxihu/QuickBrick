package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"runtime"
	"strings"
	"time"
)

var Logger *zap.Logger

func InitLogger() {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     "\n",
		EncodeLevel:    levelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.AddSync(os.Stdout),
		zap.NewAtomicLevelAt(zapcore.DebugLevel),
	)

	Logger = zap.New(core, zap.AddCaller())
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02T15:04:05.000-0700"))
}

func levelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(level.CapitalString())
}

func log(level zapcore.Level, msg, ip string, status int, fields ...zap.Field) {
	var extraFields []string
	for _, field := range fields {
		var value string
		switch field.Type {
		case zapcore.StringType:
			value = field.String
		case zapcore.Int64Type, zapcore.Int32Type, zapcore.Int16Type, zapcore.Int8Type:
			value = fmt.Sprintf("%d", field.Integer)
		case zapcore.ErrorType:
			if field.Interface != nil {
				value = fmt.Sprintf("%v", field.Interface)
			} else {
				value = "nil"
			}
		default:
			if field.Interface != nil {
				value = fmt.Sprintf("%v", field.Interface)
			} else {
				value = "nil"
			}
		}
		extraFields = append(extraFields, fmt.Sprintf("%s: %s", field.Key, value))
	}

	extraStr := strings.Join(extraFields, ", ")

	// 加这一段来获取调用文件和行号
	_, file, line, ok := runtime.Caller(2) // 调用链向上 2 层，即 Info → log → runtime.Caller
	fileLine := ""
	if ok {
		shortFile := file
		if lastSlash := strings.LastIndex(file, "/"); lastSlash != -1 {
			shortFile = file[lastSlash+1:]
		}
		fileLine = fmt.Sprintf("%s:%d", shortFile, line)
	}

	fullMsg := fmt.Sprintf("%s\t%s\t%s\t%d\t%s\tmsg: %s\t%s",
		time.Now().Format("2006-01-02T15:04:05.000-0700"),
		ip,
		level.CapitalString(),
		status,
		fileLine,
		msg,
		extraStr,
	)

	os.Stdout.WriteString(fullMsg + "\n")
}

func Info(msg string, ip string, status int, fields ...zap.Field) {
	log(zapcore.InfoLevel, msg, ip, status, fields...)
}

func Warn(msg string, ip string, status int, fields ...zap.Field) {
	log(zapcore.WarnLevel, msg, ip, status, fields...)
}

func Error(msg string, ip string, status int, fields ...zap.Field) {
	log(zapcore.ErrorLevel, msg, ip, status, fields...)
}

func Fatal(msg string, ip string, status int, fields ...zap.Field) {
	log(zapcore.FatalLevel, msg, ip, status, fields...)
}
