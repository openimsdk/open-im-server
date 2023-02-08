package relation

import (
	"Open_IM/pkg/common/config"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Mysql struct {
	gormConn *gorm.DB
}

func (m *Mysql) GormConn() *gorm.DB {
	return m.gormConn
}

func (m *Mysql) SetGormConn(gormConn *gorm.DB) {
	m.gormConn = gormConn
}

func (m *Mysql) InitConn() *Mysql {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.Config.Mysql.DBUserName, config.Config.Mysql.DBPassword, config.Config.Mysql.DBAddress[0], "mysql")
	var db *gorm.DB
	db, err := gorm.Open(mysql.Open(dsn), nil)
	if err != nil {
		time.Sleep(time.Duration(30) * time.Second)
		db, err = gorm.Open(mysql.Open(dsn), nil)
		if err != nil {
			panic(err.Error() + " open failed " + dsn)
		}
	}
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s default charset utf8 COLLATE utf8_general_ci;", config.Config.Mysql.DBDatabaseName)
	err = db.Exec(sql).Error
	if err != nil {
		panic(err.Error() + " Exec failed:" + sql)
	}
	dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.Config.Mysql.DBUserName, config.Config.Mysql.DBPassword, config.Config.Mysql.DBAddress[0], config.Config.Mysql.DBDatabaseName)
	newLogger := logger.New(
		Writer{},
		logger.Config{
			SlowThreshold:             time.Duration(config.Config.Mysql.SlowThreshold) * time.Millisecond, // Slow SQL threshold
			LogLevel:                  logger.LogLevel(config.Config.Mysql.LogLevel),                       // Log level
			IgnoreRecordNotFoundError: true,                                                                // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,                                                                // Disable color
		},
	)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err.Error() + " Open failed " + dsn)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err.Error() + " DB.DB() failed ")
	}
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(config.Config.Mysql.DBMaxLifeTime))
	sqlDB.SetMaxOpenConns(config.Config.Mysql.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(config.Config.Mysql.DBMaxIdleConns)
	m.SetGormConn(db)
	return m
}

//models := []interface{}{&Friend{}, &FriendRequest{}, &Group{}, &GroupMember{}, &GroupRequest{},
//	&User{}, &Black{}, &ChatLog{}, &Conversation{}, &AppVersion{}}

func (m *Mysql) AutoMigrateModel(model interface{}) error {
	err := m.gormConn.AutoMigrate(model)
	if err != nil {
		return err
	}
	m.gormConn.Set("gorm:table_options", "CHARSET=utf8")
	m.gormConn.Set("gorm:table_options", "collation=utf8_unicode_ci")
	_ = m.gormConn.Migrator().CreateTable(model)
	return nil
}

type Writer struct{}

func (w Writer) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func getDBConn(db *gorm.DB, tx []any) *gorm.DB {
	if len(tx) > 0 {
		if txDb, ok := tx[0].(*gorm.DB); ok {
			return txDb
		}
	}
	return db
}
