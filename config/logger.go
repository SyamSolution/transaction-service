package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerConf struct {
	dep *zap.Logger
}

type Logger interface {
	Info(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
}

func SetupLogger() Logger {
	lg, _ := zap.NewProduction()
	defer func() {
		_ = lg.Sync()
	}()
	return &LoggerConf{
		dep: lg,
	}
}

func (l *LoggerConf) Info(msg string, fields ...zapcore.Field) {
	l.dep.Info(msg, fields...)
}

func (l *LoggerConf) Error(msg string, fields ...zapcore.Field) {
	l.dep.Error(msg, fields...)
}
