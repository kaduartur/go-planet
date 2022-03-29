package log

import (
	"context"

	"github.com/kaduartur/go-planet/pkg/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var LoggerZap *zap.Logger

type zapLogger struct {
	*zap.SugaredLogger
}

func ZapLogger(config *env.Config) Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.LevelKey = "log_level"
	encoderConfig.MessageKey = "message"
	encoderConfig.StacktraceKey = "stacktrace"
	encoderConfig.TimeKey = "timestamp_app"

	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
		EncoderConfig:    encoderConfig,
	}

	if config.App.Env == "dev" {
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		cfg.Encoding = "console"
		cfg.DisableCaller = false
	}

	zlog, err := cfg.Build()
	LoggerZap = zlog
	if err != nil {
		panic(err)
	}
	logger := zlog.Sugar()
	logger.Sync()

	return zapLogger{
		SugaredLogger: logger,
	}
}

func (l zapLogger) With(v ...interface{}) Logger {
	return zapLogger{
		SugaredLogger: l.SugaredLogger.With(v...),
	}
}

func (l zapLogger) WithContext(ctx context.Context) Logger {
	fields, ok := ctx.Value(FieldsContextKey).([]interface{})
	if !ok {
		return l
	}

	return l.With(fields...)
}
