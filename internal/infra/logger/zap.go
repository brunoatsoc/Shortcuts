package logger

import (
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

	if err != nil {
		return nil, err
	}

	return &zapLogger{log: logger.Sugar()}, nil
}

func (l *zapLogger) Info(msg string, args ...any) {
	if len(args)%2 != 0 {
		args = append(args, "invalid_args", true)
	}

	l.log.Infow(msg, args...)
}
func (l *zapLogger) Warn(msg string, args ...any) {
	if len(args)%2 != 0 {
		args = append(args, "invalid_args", true)
	}

	l.log.Warnw(msg, args...)
}
func (l *zapLogger) Error(msg string, args ...any) {
	if len(args)%2 != 0 {
		args = append(args, "invalid_args", true)
	}

	l.log.Errorw(msg, args...)
}
func (l *zapLogger) Debug(msg string, args ...any) {
	if len(args)%2 != 0 {
		args = append(args, "invalid_args", true)
	}

	l.log.Debugw(msg, args...)
}
func (l *zapLogger) Fatal(msg string, args ...any) {
	if len(args)%2 != 0 {
		args = append(args, "invalid_args", true)
	}

	l.log.Fatalw(msg, args...)
}

func (l *zapLogger) With(args ...any) log.Logger {
	if len(args)%2 != 0 {
		args = append(args, "invalid_args", true)
	}

	return &zapLogger{log: l.log.With(args...)}
}
