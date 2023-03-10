package log

import (
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/tracelog"
	"context"
	"time"

	"github.com/go-logr/logr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	discardLogger        = logr.Discard()
	defaultLogger Logger = LogRLogger(discardLogger)
	pkgLogger     Logger = LogRLogger(discardLogger)
)

// InitFromConfig initializes a Zap-based logger
func InitFromConfig(name string) {
	//var c zap.Config
	//file, _ := os.Create(config.Config.Log.StorageLocation)
	//writeSyncer := zapcore.AddSync(file)

	l, err := NewZapLogger()
	if err == nil {
		setLogger(l, name)
	}
}

// GetLogger returns the logger that was set with SetLogger with an extra depth of 1
func GetLogger() Logger {
	return defaultLogger
}

// SetLogger lets you use a custom logger. Pass in a logr.Logger with default depth
func setLogger(l Logger, name string) {
	defaultLogger = l.WithCallDepth(1).WithName(name)
	// pkg wrapper needs to drop two levels of depth
	pkgLogger = l.WithCallDepth(2).WithName(name)
}

func Debug(ctx context.Context, msg string, keysAndValues ...interface{}) {
	pkgLogger.Debug(ctx, msg, keysAndValues...)
}

func Info(ctx context.Context, msg string, keysAndValues ...interface{}) {
	pkgLogger.Info(ctx, msg, keysAndValues...)
}

func Warn(ctx context.Context, msg string, err error, keysAndValues ...interface{}) {
	pkgLogger.Warn(ctx, msg, err, keysAndValues...)
}

func Error(ctx context.Context, msg string, err error, keysAndValues ...interface{}) {
	pkgLogger.Error(ctx, msg, err, keysAndValues...)
}

func ParseZapLevel(level string) zapcore.Level {
	lvl := zapcore.InfoLevel
	if level != "" {
		_ = lvl.UnmarshalText([]byte(level))
	}
	return lvl
}

type Logger interface {
	Debug(ctx context.Context, msg string, keysAndValues ...interface{})
	Info(ctx context.Context, msg string, keysAndValues ...interface{})
	Warn(ctx context.Context, msg string, err error, keysAndValues ...interface{})
	Error(ctx context.Context, msg string, err error, keysAndValues ...interface{})
	WithValues(keysAndValues ...interface{}) Logger
	WithName(name string) Logger
	WithCallDepth(depth int) Logger
	WithItemSampler() Logger
	// WithoutSampler returns the original logger without sampling
	WithoutSampler() Logger
}

type ZapLogger struct {
	zap *zap.SugaredLogger
	// store original logger without sampling to avoid multiple samplers
	unsampled      *zap.SugaredLogger
	SampleDuration time.Duration
	SampleInitial  int
	SampleInterval int
}

func NewZapLogger() (*ZapLogger, error) {
	zapConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Development:      true,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{config.Config.Log.StorageLocation},
		ErrorOutputPaths: []string{"stderr"},
	}
	l, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}
	zl := &ZapLogger{
		unsampled: l.Sugar(),
		//SampleDuration: time.Duration(conf.ItemSampleSeconds) * time.Second,
		//SampleInitial:  conf.ItemSampleInitial,
		//SampleInterval: conf.ItemSampleInterval,
	}

	//if conf.Sample {
	//	// use a sampling logger for the main logger
	//	samplingConf := &zap.SamplingConfig{
	//		Initial:    conf.SampleInitial,
	//		Thereafter: conf.SampleInterval,
	//	}
	//	// sane defaults
	//	if samplingConf.Initial == 0 {
	//		samplingConf.Initial = 20
	//	}
	//	if samplingConf.Thereafter == 0 {
	//		samplingConf.Thereafter = 100
	//	}
	//	zl.zap = l.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
	//		return zapcore.NewSamplerWithOptions(
	//			core,
	//			time.Second,
	//			samplingConf.Initial,
	//			samplingConf.Thereafter,
	//		)
	//	})).Sugar()
	//} else {
	//	zl.zap = zl.unsampled
	//}
	return zl, nil
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
	// mirror unsampled logger too
	if l.unsampled == l.zap {
		dup.unsampled = dup.zap
	} else {
		dup.unsampled = l.unsampled.With(keysAndValues...)
	}
	return &dup
}

func (l *ZapLogger) WithName(name string) Logger {
	dup := *l
	dup.zap = l.zap.Named(name)
	if l.unsampled == l.zap {
		dup.unsampled = dup.zap
	} else {
		dup.unsampled = l.unsampled.Named(name)
	}
	return &dup
}

func (l *ZapLogger) WithCallDepth(depth int) Logger {
	dup := *l
	dup.zap = l.zap.WithOptions(zap.AddCallerSkip(depth))
	if l.unsampled == l.zap {
		dup.unsampled = dup.zap
	} else {
		dup.unsampled = l.unsampled.WithOptions(zap.AddCallerSkip(depth))
	}
	return &dup
}

func (l *ZapLogger) WithItemSampler() Logger {
	if l.SampleDuration == 0 {
		return l
	}
	dup := *l
	dup.zap = l.unsampled.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewSamplerWithOptions(
			core,
			l.SampleDuration,
			l.SampleInitial,
			l.SampleInterval,
		)
	}))
	return &dup
}

func (l *ZapLogger) WithoutSampler() Logger {
	if l.SampleDuration == 0 {
		return l
	}
	dup := *l
	dup.zap = l.unsampled
	return &dup
}

type LogRLogger logr.Logger

func (l LogRLogger) toLogr() logr.Logger {
	if logr.Logger(l).GetSink() == nil {
		return discardLogger
	}
	return logr.Logger(l)
}

func (l LogRLogger) Debug(ctx context.Context, msg string, keysAndValues ...interface{}) {
	keysAndValues = append([]interface{}{constant.OperationID, tracelog.GetOperationID(ctx)}, keysAndValues...)
	l.toLogr().V(1).Info(msg, keysAndValues...)
}

func (l LogRLogger) Info(ctx context.Context, msg string, keysAndValues ...interface{}) {
	keysAndValues = append([]interface{}{constant.OperationID, tracelog.GetOperationID(ctx)}, keysAndValues...)
	l.toLogr().Info(msg, keysAndValues...)
}

func (l LogRLogger) Warn(ctx context.Context, msg string, err error, keysAndValues ...interface{}) {
	if err != nil {
		keysAndValues = append(keysAndValues, "error", err)
	}
	keysAndValues = append([]interface{}{constant.OperationID, tracelog.GetOperationID(ctx)}, keysAndValues...)
	l.toLogr().Info(msg, keysAndValues...)
}

func (l LogRLogger) Error(ctx context.Context, msg string, err error, keysAndValues ...interface{}) {
	keysAndValues = append([]interface{}{constant.OperationID, tracelog.GetOperationID(ctx)}, keysAndValues...)
	l.toLogr().Error(err, msg, keysAndValues...)
}

func (l LogRLogger) WithValues(keysAndValues ...interface{}) Logger {
	return LogRLogger(l.toLogr().WithValues(keysAndValues...))
}

func (l LogRLogger) WithName(name string) Logger {
	return LogRLogger(l.toLogr().WithName(name))
}

func (l LogRLogger) WithCallDepth(depth int) Logger {
	return LogRLogger(l.toLogr().WithCallDepth(depth))
}

func (l LogRLogger) WithItemSampler() Logger {
	return l
}

func (l LogRLogger) WithoutSampler() Logger {
	return l
}
