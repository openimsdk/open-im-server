// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package relation

import (
	"fmt"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"

	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/mw/specialerror"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	maxRetry = 100 // number of retries
)

// newMysqlGormDB Initialize the database connection.
func newMysqlGormDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.Config.Mysql.Username, config.Config.Mysql.Password, config.Config.Mysql.Address[0], "mysql")

	db, err := connectToDatabase(dsn, maxRetry)
	if err != nil {
		panic(err.Error() + " Open failed " + dsn)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer sqlDB.Close()
	sql := fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS %s default charset utf8mb4 COLLATE utf8mb4_unicode_ci;",
		config.Config.Mysql.Database,
	)
	err = db.Exec(sql).Error
	if err != nil {
		return nil, fmt.Errorf("init db %w", err)
	}
	dsn = fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.Config.Mysql.Username,
		config.Config.Mysql.Password,
		config.Config.Mysql.Address[0],
		config.Config.Mysql.Database,
	)
	sqlLogger := log.NewSqlLogger(
		logger.LogLevel(config.Config.Mysql.LogLevel),
		true,
		time.Duration(config.Config.Mysql.SlowThreshold)*time.Millisecond,
	)
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
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(config.Config.Mysql.MaxLifeTime))
	sqlDB.SetMaxOpenConns(config.Config.Mysql.MaxOpenConn)
	sqlDB.SetMaxIdleConns(config.Config.Mysql.MaxIdleConn)
	return db, nil
}

// connectToDatabase Connection retry for mysql.
func connectToDatabase(dsn string, maxRetry int) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	for i := 0; i <= maxRetry; i++ {
		db, err = gorm.Open(mysql.Open(dsn), nil)
		if err == nil {
			return db, nil
		}
		if mysqlErr, ok := err.(*mysqldriver.MySQLError); ok && mysqlErr.Number == 1045 {
			return nil, err
		}
		time.Sleep(time.Duration(1) * time.Second)
	}
	return nil, err
}

// NewGormDB gorm mysql.
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
	if mysqlErr, ok := err.(*mysqldriver.MySQLError); ok {
		return mysqlErr.Number == 1062
	}
	return false
}
