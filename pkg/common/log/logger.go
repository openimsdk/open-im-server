package log

import "context"

type Logger interface {
	Debug(ctx context.Context, msg string, keysAndValues ...interface{})
	Info(ctx context.Context, msg string, keysAndValues ...interface{})
	Warn(ctx context.Context, msg string, err error, keysAndValues ...interface{})
	Error(ctx context.Context, msg string, err error, keysAndValues ...interface{})
	WithValues(keysAndValues ...interface{}) LogrusLogger
	WithName(name string) LogrusLogger
	WithCallDepth(depth int) LogrusLogger
}
