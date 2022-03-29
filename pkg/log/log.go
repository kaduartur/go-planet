package log

import "context"

type Logger interface {
	WithContext(ctx context.Context) Logger
	With(v ...interface{}) Logger
	Info(v ...interface{})
	Error(v ...interface{})
	Warn(v ...interface{})
	Fatal(v ...interface{})
}

var FieldsContextKey struct{}
