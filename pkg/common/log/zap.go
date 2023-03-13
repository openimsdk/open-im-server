package log

import (
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/tracelog"
	"context"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	pkgLogger Logger = &ZapLogger{}
	sp               = string(filepath.Separator)
)

// InitFromConfig initializes a Zap-based logger
func InitFromConfig(name string) error {
	l, err := NewZapLogger()
	if err != nil {
		return err
	}
	pkgLogger = l.WithCallDepth(2).WithName(name)
	return nil
}

func ZDebug(ctx context.Context, msg string, keysAndValues ...interface{}) {
	pkgLogger.Debug(ctx, msg, keysAndValues...)
}

func ZInfo(ctx context.Context, msg string, keysAndValues ...interface{}) {
	pkgLogger.Info(ctx, msg, keysAndValues...)
}

func ZWarn(ctx context.Context, msg string, err error, keysAndValues ...interface{}) {
	pkgLogger.Warn(ctx, msg, err, keysAndValues...)
}

func ZError(ctx context.Context, msg string, err error, keysAndValues ...interface{}) {
	pkgLogger.Error(ctx, msg, err, keysAndValues...)
}

type ZapLogger struct {
	zap *zap.SugaredLogger
	// store original logger without sampling to avoid multiple samplers
	SampleDuration time.Duration
	SampleInitial  int
	SampleInterval int
}

func NewZapLogger() (*ZapLogger, error) {
	zapConfig := zap.Config{
		Level:         zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Development:   true,
		Encoding:      "json",
		EncoderConfig: zap.NewProductionEncoderConfig(),
	}
	zl := &ZapLogger{}
	if config.Config.Log.Stderr {
		zapConfig.OutputPaths = append(zapConfig.OutputPaths, "stderr")
	}
	l, err := zapConfig.Build(zl.cores())
	if err != nil {
		return nil, err
	}
	zl.zap = l.Sugar()
	return zl, nil
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func (l *ZapLogger) cores() zap.Option {
	fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	writer := l.getWriter()
	var cores []zapcore.Core
	if config.Config.Log.StorageLocation != "" {
		cores = []zapcore.Core{
			zapcore.NewCore(fileEncoder, writer, zapcore.DebugLevel),
		}
	}
	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	})
}

func (l *ZapLogger) getWriter() zapcore.WriteSyncer {
	logf, _ := rotatelogs.New(config.Config.Log.StorageLocation+sp+"openIM"+"-"+".%Y_%m%d_%H",
		rotatelogs.WithLinkName(config.Config.Log.StorageLocation+sp+"openIM"+"-"),
		rotatelogs.WithMaxAge(2*24*time.Hour),
		rotatelogs.WithRotationTime(time.Minute),
	)
	return zapcore.AddSync(logf)
}

func (l *ZapLogger) ToZap() *zap.SugaredLogger {
	return l.zap
}

func (l *ZapLogger) Debug(ctx context.Context, msg string, keysAndValues ...interface{}) {
	keysAndValues = append([]interface{}{constant.OperationID, tracelog.GetOperationID(ctx)}, keysAndValues...)
	l.zap.Debugw(msg, keysAndValues...)
}

func (l *ZapLogger) Info(ctx context.Context, msg string, keysAndValues ...interface{}) {
	keysAndValues = append([]interface{}{constant.OperationID, tracelog.GetOperationID(ctx)}, keysAndValues...)
	l.zap.Infow(msg, keysAndValues...)
}

func (l *ZapLogger) Warn(ctx context.Context, msg string, err error, keysAndValues ...interface{}) {
	if err != nil {
		keysAndValues = append(keysAndValues, "error", err)
	}
	keysAndValues = append([]interface{}{constant.OperationID, tracelog.GetOperationID(ctx)}, keysAndValues...)
	l.zap.Warnw(msg, keysAndValues...)
}

func (l *ZapLogger) Error(ctx context.Context, msg string, err error, keysAndValues ...interface{}) {
	if err != nil {
		keysAndValues = append(keysAndValues, "error", err)
	}
	keysAndValues = append([]interface{}{constant.OperationID, tracelog.GetOperationID(ctx)}, keysAndValues...)
	l.zap.Errorw(msg, keysAndValues...)
}

func (l *ZapLogger) WithValues(keysAndValues ...interface{}) Logger {
	dup := *l
	dup.zap = l.zap.With(keysAndValues...)
	return &dup
}

func (l *ZapLogger) WithName(name string) Logger {
	dup := *l
	dup.zap = l.zap.Named(name)
	return &dup
}

func (l *ZapLogger) WithCallDepth(depth int) Logger {
	dup := *l
	dup.zap = l.zap.WithOptions(zap.AddCallerSkip(depth))
	return &dup
}
