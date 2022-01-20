package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"

	zd "github.com/blendle/zapdriver"
	"go.uber.org/zap"
)

func CreateLogger(name string) LoggerWrapper {
	env := os.Getenv("RS_DROP_EMULATOR_ENV")

	if env == "PROD" {
		zl, err := zd.NewProductionWithCore(zd.WrapCore(
			zd.ReportAllErrors(true),
			zd.ServiceName("rs-drop-emulator"),
		))

		if err != nil {
			log.Fatalf("couldn't configure logger: %v", err)
		}

		return &wrappedLogger{
			base:    zl,
			sugared: zl.Sugar().Named(name),
			dev:     false,
		}
	}

	zl, _ := zap.NewDevelopment()
	return &wrappedLogger{
		base:    zl,
		sugared: zl.Sugar().Named(name),
		dev:     true,
	}
}

type LoggerWrapper interface {
	Info(string, ...interface{})
	Debug(string, ...interface{})
	Error(string, ...interface{})
	Warn(string, ...interface{})
	Fatal(string, ...interface{})
	Sync()
}

type wrappedLogger struct {
	base    *zap.Logger
	sugared *zap.SugaredLogger
	dev     bool
}

func wrapWithLabel(fields ...interface{}) []zap.Field {
	var newValues []zap.Field
	for i := 0; i < len(fields)-1; i += 2 {
		newValues = append(newValues, zd.Label(fmt.Sprint(fields[i]), fmt.Sprint(fields[i+1])))
	}
	newValues = append(newValues, zd.SourceLocation(runtime.Caller(2))) // Skip 2 cause there are 2 calls here
	return newValues
}

func (l wrappedLogger) Info(msg string, fields ...interface{}) {
	if l.dev {
		l.sugared.Info(msg, fields)
	} else {
		l.base.Info(msg, wrapWithLabel(fields...)...)
	}
}

func (l wrappedLogger) Warn(msg string, fields ...interface{}) {
	if l.dev {
		l.sugared.Warnw(msg, fields)
	} else {
		l.base.Warn(msg, wrapWithLabel(fields...)...)
	}
}

func (l wrappedLogger) Debug(msg string, fields ...interface{}) {
	if l.dev {
		l.sugared.Debugw(msg, fields)
	} else {
		l.base.Debug(msg, wrapWithLabel(fields...)...)
	}
}

func (l wrappedLogger) Error(msg string, fields ...interface{}) {
	if l.dev {
		l.sugared.Error(msg, fields)
	} else {
		l.base.Error(msg, wrapWithLabel(fields...)...)
	}
}

func (l wrappedLogger) Fatal(msg string, fields ...interface{}) {
	if l.dev {
		l.sugared.Fatal(msg, fields)
	} else {
		l.base.Fatal(msg, wrapWithLabel(fields...)...)
	}
}

func (l wrappedLogger) Sync() {
	l.base.Sync()
}