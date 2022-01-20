package logger

import (
	"fmt"
	"os"

	zd "github.com/blendle/zapdriver"
	"go.uber.org/zap"
)

func CreateLogger(name string) LoggerWrapper {
	env := os.Getenv("RS_DROP_EMULATOR_ENV")

	if env == "PROD" {
		zl, _ := zd.NewProduction()
		return &wrappedLogger{
			base: zl.Sugar().Named(name),
			dev:  false,
		}
	}
	zl, _ := zap.NewDevelopment()
	return &wrappedLogger{
		base: zl.Sugar().Named(name),
		dev:  true,
	}
}

type LoggerWrapper interface {
	Info(...interface{})
	Infow(string, ...interface{})
	Infof(string, ...interface{})
	Debug(...interface{})
	Debugw(string, ...interface{})
	Debugf(string, ...interface{})
	Error(...interface{})
	Errorw(string, ...interface{})
	Errorf(string, ...interface{})
	Warn(...interface{})
	Warnw(string, ...interface{})
	Warnf(string, ...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalw(string, ...interface{})
	Sync()
}

type wrappedLogger struct {
	base *zap.SugaredLogger
	dev  bool
}

func (l wrappedLogger) Info(msg ...interface{}) {
	l.base.Info(msg)
}

func (l wrappedLogger) Infow(msg string, values ...interface{}) {
	if l.dev {
		l.base.Infow(msg, values)
		return
	}
	var newValues []zap.Field
	for i := 0; i < len(values)-1; i += 2 {
		newValues = append(newValues, zd.Label(fmt.Sprint(values[i]), fmt.Sprint(values[i+1])))
	}
	l.base.Infow(msg, newValues)
}

func (l wrappedLogger) Infof(template string, values ...interface{}) {
	l.base.Infof(template, values)
}

func (l wrappedLogger) Error(msg ...interface{}) {
	l.base.Error(msg)
}

func (l wrappedLogger) Errorw(msg string, values ...interface{}) {
	if l.dev {
		l.base.Errorw(msg, values)
		return
	}
	var newValues []zap.Field
	for i := 0; i < len(values)-1; i += 2 {
		newValues = append(newValues, zd.Label(fmt.Sprint(values[i]), fmt.Sprint(values[i+1])))
	}
	l.base.Errorw(msg, newValues)
}

func (l wrappedLogger) Errorf(template string, values ...interface{}) {
	l.base.Errorf(template, values)
}

func (l wrappedLogger) Debug(msg ...interface{}) {
	l.base.Debug(msg)
}

func (l wrappedLogger) Debugw(msg string, values ...interface{}) {
	if l.dev {
		l.base.Debugw(msg, values)
		return
	}
	var newValues []zap.Field
	for i := 0; i < len(values)-1; i += 2 {
		newValues = append(newValues, zd.Label(fmt.Sprint(values[i]), fmt.Sprint(values[i+1])))
	}
	l.base.Debugw(msg, newValues)
}

func (l wrappedLogger) Debugf(template string, values ...interface{}) {
	l.base.Debugf(template, values)
}

func (l wrappedLogger) Warn(msg ...interface{}) {
	l.base.Warn(msg)
}

func (l wrappedLogger) Warnw(msg string, values ...interface{}) {
	if l.dev {
		l.base.Warnw(msg, values)
		return
	}
	var newValues []zap.Field
	for i := 0; i < len(values); i++ {
		newValues = append(newValues, zd.Label(fmt.Sprint(values[i]), fmt.Sprint(values[i+1])))
	}
	l.base.Warnw(msg, newValues)
}

func (l wrappedLogger) Warnf(template string, values ...interface{}) {
	l.base.Warnf(template, values)
}

func (l wrappedLogger) Fatal(msg ...interface{}) {
	l.base.Fatal(msg)
}

func (l wrappedLogger) Fatalw(msg string, values ...interface{}) {
	if l.dev {
		l.base.Fatalw(msg, values)
		return
	}
	var newValues []zap.Field
	for i := 0; i < len(values)-1; i += 2 {
		newValues = append(newValues, zd.Label(fmt.Sprint(values[i]), fmt.Sprint(values[i+1])))
	}
	l.base.Fatalw(msg, newValues)
}

func (l wrappedLogger) Fatalf(template string, values ...interface{}) {
	l.base.Fatalf(template, values)
}

func (l wrappedLogger) Sync() {
	l.base.Sync()
}
