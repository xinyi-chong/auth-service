package logger

import (
	"go.uber.org/zap"
	"os"
)

var (
	Log *zap.Logger
)

func Init() error {
	var err error
	env := os.Getenv("APP_ENV")

	if env == "development" {
		Log, err = zap.NewDevelopment()
	} else {
		Log, err = zap.NewProduction()
	}

	if err != nil {
		return err
	}

	Log.Info("Logger: Init", zap.String("env", env))

	return nil
}

func Debug(msg string, fields ...zap.Field) {
	Log.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Log.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Log.Fatal(msg, fields...)
}
