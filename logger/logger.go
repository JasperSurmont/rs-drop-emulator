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
	env := os.Getenv("RS_DROP_simulator_ENV")

	if env == "PROD" {
		zl, err := zd.NewProductionWithCore(zd.WrapCore(
			zd.ReportAllErrors(true),
			zd.ServiceName("rs-drop-simulator"),
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

	zl, _ := zap.NewDevelopment(zap.AddCallerSkip(1))
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
	newValues = append(newValues, zd.SourceLocation(runtime.Caller(2)))
	return newValues
}

func (l wrappedLogger) Info(msg string, fields ...interface{}) {
	if l.dev {
		if len(fields) > 0 {
			l.sugared.Infow(msg, fields...)
		} else {
			l.sugared.Info(msg)
		}
	} else {
		l.base.Info(msg, wrapWithLabel(fields...)...)
	}
}

func (l wrappedLogger) Warn(msg string, fields ...interface{}) {
	if l.dev {
		if len(fields) > 0 {
			l.sugared.Warnw(msg, fields...)
		} else {
			l.sugared.Warn(msg)
		}
	} else {
		l.base.Warn(msg, wrapWithLabel(fields...)...)
	}
}

func (l wrappedLogger) Debug(msg string, fields ...interface{}) {
	if l.dev {
		if len(fields) > 0 {
			l.sugared.Debugw(msg, fields...)
		} else {
			l.sugared.Debug(msg)
		}
	} else {
		l.base.Debug(msg, wrapWithLabel(fields...)...)
	}
}

func (l wrappedLogger) Error(msg string, fields ...interface{}) {
	if l.dev {
		if len(fields) > 0 {
			l.sugared.Errorw(msg, fields...)
		} else {
			l.sugared.Error(msg)
		}
	} else {
		l.base.Error(msg, wrapWithLabel(fields...)...)
	}
}

func (l wrappedLogger) Fatal(msg string, fields ...interface{}) {
	if l.dev {
		if len(fields) > 0 {
			l.sugared.Fatalw(msg, fields...)
		} else {
			l.sugared.Fatal(msg)
		}
	} else {
		l.base.Fatal(msg, wrapWithLabel(fields...)...)
	}
}

func (l wrappedLogger) Sync() {
	l.base.Sync()
}
