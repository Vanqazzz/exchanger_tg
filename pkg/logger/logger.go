package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LevelDebug = "DEBUG"
	LevelProd  = "PROD"
)

type Logger struct {
	z *zap.Logger
}

func (l *Logger) Get() *zap.Logger {
	return l.z

}

func New(level string) *Logger {
	var cfg zap.Config

	switch level {
	case LevelProd:
		cfg = zap.NewProductionConfig()
	default:
		cfg = zap.NewDevelopmentConfig()
	}

	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	sink := zapcore.AddSync(os.Stdout)

	encoder := zapcore.NewConsoleEncoder(cfg.EncoderConfig)

	core := zapcore.NewCore(encoder, sink, cfg.Level)

	return &Logger{z: zap.New(core)}
}
