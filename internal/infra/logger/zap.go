package logger

import (
	"fmt"
	"intellisys-backend/internal/domain/log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	log *zap.SugaredLogger
}

func NewLogger() (log.Logger, error) {
	var cfg zap.Config

	cfg = zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := cfg.Build(zap.AddCallerSkip(1)) // CallerSkip para mostrar o caller correto
	fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	if err != nil {
		return nil, err
	}

	return &zapLogger{log: logger.Sugar()}, nil
}

func (l *zapLogger) Info(msg string, args ...any)  { l.log.Infow(msg, args...) }
func (l *zapLogger) Warn(msg string, args ...any)  { l.log.Warnw(msg, args...) }
func (l *zapLogger) Error(msg string, args ...any) { l.log.Errorw(msg, args...) }
func (l *zapLogger) Debug(msg string, args ...any) { l.log.Debugw(msg, args...) }
func (l *zapLogger) Fatal(msg string, args ...any) { l.log.Fatalw(msg, args...) }

func (l *zapLogger) With(args ...any) log.Logger {
	return &zapLogger{log: l.log.With(args...)}
}
