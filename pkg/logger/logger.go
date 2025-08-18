package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

var (
	once     sync.Once
	instance *zap.Logger
)

func Init() error {
	var initErr error
	once.Do(func() {
		env := os.Getenv("APP_ENV")
		var cfg zap.Config
		if env == "production" {
			cfg = zap.NewProductionConfig()
		} else {
			cfg = zap.NewDevelopmentConfig()
		}
		cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		instance, initErr = cfg.Build()
	})
	return initErr
}

func Get() *zap.Logger {
	if instance == nil {
		panic("logger not initialized - call logger.Init() first")
	}
	return instance
}

func Sync() error {
	return instance.Sync()
}

func Debug(msg string, fields ...zap.Field) {
	instance.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	instance.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	instance.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	instance.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	instance.Fatal(msg, fields...)
}
