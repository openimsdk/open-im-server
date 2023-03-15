package relation

import (
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/log"
	"OpenIM/pkg/common/mw/specialerror"
	"OpenIM/pkg/errs"
	"context"
	"errors"
	"fmt"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm/utils"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newMysqlGormDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.Config.Mysql.DBUserName, config.Config.Mysql.DBPassword, config.Config.Mysql.DBAddress[0], "mysql")
	db, err := gorm.Open(mysql.Open(dsn), nil)
	if err != nil {
		time.Sleep(time.Duration(30) * time.Second)
		db, err = gorm.Open(mysql.Open(dsn), nil)
		if err != nil {
			panic(err.Error() + " open failed " + dsn)
		}
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer sqlDB.Close()
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s default charset utf8 COLLATE utf8_general_ci;", config.Config.Mysql.DBDatabaseName)
	err = db.Exec(sql).Error
	if err != nil {
		return nil, fmt.Errorf("init db %w", err)
	}
	dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.Config.Mysql.DBUserName, config.Config.Mysql.DBPassword, config.Config.Mysql.DBAddress[0], config.Config.Mysql.DBDatabaseName)
	//newLogger := logger.New(
	//	Writer{},
	//	logger.Config{
	//		SlowThreshold:             time.Duration(config.Config.Mysql.SlowThreshold) * time.Millisecond, // Slow SQL threshold
	//		LogLevel:                  logger.LogLevel(config.Config.Mysql.LogLevel),                       // Log level
	//		IgnoreRecordNotFoundError: true,                                                                // Ignore ErrRecordNotFound error for logger
	//	},
	//)
	sqlLogger := NewSqlLogger(logger.LogLevel(config.Config.Mysql.LogLevel), true, time.Duration(config.Config.Mysql.SlowThreshold)*time.Millisecond)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: sqlLogger,
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err = db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(config.Config.Mysql.DBMaxLifeTime))
	sqlDB.SetMaxOpenConns(config.Config.Mysql.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(config.Config.Mysql.DBMaxIdleConns)
	return db, nil
}

// gorm mysql
func NewGormDB() (*gorm.DB, error) {
	specialerror.AddReplace(gorm.ErrRecordNotFound, errs.ErrRecordNotFound)
	specialerror.AddErrHandler(replaceDuplicateKey)
	return newMysqlGormDB()
}

func replaceDuplicateKey(err error) errs.CodeError {
	if IsMysqlDuplicateKey(err) {
		return errs.ErrDuplicateKey
	}
	return nil
}

func IsMysqlDuplicateKey(err error) bool {
	if mysqlErr, ok := err.(*mysqlDriver.MySQLError); ok {
		return mysqlErr.Number == 1062
	}
	return false
}

type SqlLogger struct {
	LogLevel                  logger.LogLevel
	IgnoreRecordNotFoundError bool
	SlowThreshold             time.Duration
}

func NewSqlLogger(logLevel logger.LogLevel, ignoreRecordNotFoundError bool, slowThreshold time.Duration) *SqlLogger {
	return &SqlLogger{
		LogLevel:                  logLevel,
		IgnoreRecordNotFoundError: ignoreRecordNotFoundError,
		SlowThreshold:             slowThreshold,
	}
}

func (l *SqlLogger) LogMode(logLevel logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = logLevel
	return &newLogger
}

func (SqlLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	log.ZInfo(ctx, msg, args)
}

func (SqlLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	log.ZWarn(ctx, msg, nil, args)
}

func (SqlLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	log.ZError(ctx, msg, nil, args)
}

func (l SqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			log.ZError(ctx, utils.FileWithLineNum(), err, "time", float64(elapsed.Nanoseconds())/1e6, "sql", sql)
		} else {
			log.ZError(ctx, utils.FileWithLineNum(), err, "time", float64(elapsed.Nanoseconds())/1e6, "rows", rows, "sql", sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			log.ZWarn(ctx, utils.FileWithLineNum(), nil, "slow sql", slowLog, "time", float64(elapsed.Nanoseconds())/1e6, "sql", sql)
		} else {
			log.ZWarn(ctx, utils.FileWithLineNum(), nil, "slow sql", slowLog, "time", float64(elapsed.Nanoseconds())/1e6, "rows", rows, "sql", sql)
		}
	case l.LogLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			log.ZDebug(ctx, utils.FileWithLineNum(), "time", float64(elapsed.Nanoseconds())/1e6, "sql", sql)
		} else {
			log.ZDebug(ctx, utils.FileWithLineNum(), "time", float64(elapsed.Nanoseconds())/1e6, "rows", rows, "sql", sql)

		}
	}
}

type Writer struct{}

func (w Writer) Printf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	l := strings.Split(s, "\n")
	if len(l) == 2 {
		log.ZDebug(context.Background(), "sql exec detail", "gorm", l[0], "sql", l[1])
	} else {
		log.ZDebug(context.Background(), "sql exec detail", "sql", s)
	}
}
