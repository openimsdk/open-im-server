package relation

import (
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mw/specialerror"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
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
	sqlLogger := log.NewSqlLogger(logger.LogLevel(config.Config.Mysql.LogLevel), true, time.Duration(config.Config.Mysql.SlowThreshold)*time.Millisecond)
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
