package log

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	gormUtils "gorm.io/gorm/utils"
)

type SqlLogger struct {
	LogLevel                  gormLogger.LogLevel
	IgnoreRecordNotFoundError bool
	SlowThreshold             time.Duration
}

func NewSqlLogger(logLevel gormLogger.LogLevel, ignoreRecordNotFoundError bool, slowThreshold time.Duration) *SqlLogger {
	return &SqlLogger{
		LogLevel:                  logLevel,
		IgnoreRecordNotFoundError: ignoreRecordNotFoundError,
		SlowThreshold:             slowThreshold,
	}
}

func (l *SqlLogger) LogMode(logLevel gormLogger.LogLevel) gormLogger.Interface {
	newLogger := *l
	newLogger.LogLevel = logLevel
	return &newLogger
}

func (SqlLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	ZInfo(ctx, msg, args)
}

func (SqlLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	ZWarn(ctx, msg, nil, args)
}

func (SqlLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	ZError(ctx, msg, nil, args)
}

func (l *SqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= gormLogger.Silent {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= gormLogger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			ZError(ctx, "sql exec detail", err, "gorm", gormUtils.FileWithLineNum(), "elapsed time", fmt.Sprintf("%f(ms)", float64(elapsed.Nanoseconds())/1e6), "sql", sql)
		} else {
			ZError(ctx, "sql exec detail", err, "gorm", gormUtils.FileWithLineNum(), "elapsed time", fmt.Sprintf("%f(ms)", float64(elapsed.Nanoseconds())/1e6), "rows", rows, "sql", sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormLogger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			ZWarn(ctx, "sql exec detail", nil, "gorm", gormUtils.FileWithLineNum(), nil, "slow sql", slowLog, "elapsed time", fmt.Sprintf("%f(ms)", float64(elapsed.Nanoseconds())/1e6), "sql", sql)
		} else {
			ZWarn(ctx, "sql exec detail", nil, "gorm", gormUtils.FileWithLineNum(), nil, "slow sql", slowLog, "elapsed time", fmt.Sprintf("%f(ms)", float64(elapsed.Nanoseconds())/1e6), "rows", rows, "sql", sql)
		}
	case l.LogLevel == gormLogger.Info:
		sql, rows := fc()
		if rows == -1 {
			ZDebug(ctx, "sql exec detail", "gorm", gormUtils.FileWithLineNum(), "elapsed time", fmt.Sprintf("%f(ms)", float64(elapsed.Nanoseconds())/1e6), "sql", sql)
		} else {
			ZDebug(ctx, "sql exec detail", "gorm", gormUtils.FileWithLineNum(), "elapsed time", fmt.Sprintf("%f(ms)", float64(elapsed.Nanoseconds())/1e6), "rows", rows, "sql", sql)
		}
	}
}
