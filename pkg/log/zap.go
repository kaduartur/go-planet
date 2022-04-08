package log

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger
var ZapLogger *zap.Logger

type loggerKeyType int

const loggerKey loggerKeyType = iota

type Event map[string]interface{}

func init() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.LevelKey = "log_level"
	encoderConfig.MessageKey = "message"
	encoderConfig.StacktraceKey = "stacktrace"
	encoderConfig.TimeKey = "timestamp"

	encoding := "json"
	if os.Getenv("ENV") != "prod" {
		encoding = "console"
	}

	cfg := zap.Config{
		Encoding:         encoding,
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
		EncoderConfig:    encoderConfig,
		DisableCaller:    true,
	}

	var err error
	ZapLogger, err = cfg.Build()
	if err != nil {
		panic(err)
	}

	logger = ZapLogger.Sugar()
	logger.Sync()
}

func addEvent(event ...Event) zap.Field {
	if len(event) > 0 {
		return zap.Any("event", event[0])
	}
	return zap.Skip()
}

func Debug(context context.Context, msg string, event ...Event) {
	WithContext(context).Debug(msg, addEvent(event...))
}

func Info(context context.Context, msg string, event ...Event) {
	WithContext(context).Info(msg, addEvent(event...))
}

func Warn(context context.Context, msg string, event ...Event) {
	WithContext(context).Warn(msg, addEvent(event...))
}

func Error(context context.Context, msg string, err error, event ...Event) {
	WithContext(context).Error(msg, zap.Error(err), addEvent(event...))
}

func Fatal(context context.Context, msg string, err error, event ...Event) {
	WithContext(context).Fatal(msg, zap.Error(err), addEvent(event...))
}

func Sync() {
	logger.Desugar().Sync()
}

func NewLogContext(ctx context.Context, fields map[string]interface{}) context.Context {
	var zapFields []zap.Field

	for index, value := range fields {
		zapFields = append(zapFields, zap.Any(index, value))
	}

	return context.WithValue(ctx, loggerKey, WithContext(ctx).With(zapFields...))
}

func WithExtra(ctx context.Context, fields map[string]interface{}) context.Context {
	return context.WithValue(ctx, loggerKey, WithContext(ctx).With(zap.Any("extra", fields)))
}

func WithContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return ZapLogger
	}

	if ctxLogger, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return ctxLogger
	}

	return ZapLogger
}
