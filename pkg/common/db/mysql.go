package db

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type mysqlDB struct {
	//sync.RWMutex
	db *gorm.DB
}

type Writer struct{}

func (w Writer) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func initMysqlDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.Config.Mysql.DBUserName, config.Config.Mysql.DBPassword, config.Config.Mysql.DBAddress[0], "mysql")
	var db *gorm.DB
	var err1 error
	db, err := gorm.Open(mysql.Open(dsn), nil)
	if err != nil {
		time.Sleep(time.Duration(30) * time.Second)
		db, err1 = gorm.Open(mysql.Open(dsn), nil)
		if err1 != nil {
			panic(err1.Error() + " open failed " + dsn)
		}
	}

	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s default charset utf8 COLLATE utf8_general_ci;", config.Config.Mysql.DBDatabaseName)
	err = db.Exec(sql).Error
	if err != nil {
		panic(err.Error() + " Exec failed " + sql)
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
		panic(err.Error() + " db.DB() failed ")
	}

	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(config.Config.Mysql.DBMaxLifeTime))
	sqlDB.SetMaxOpenConns(config.Config.Mysql.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(config.Config.Mysql.DBMaxIdleConns)

	db.AutoMigrate(
		&im_mysql_model.Register{},
		&im_mysql_model.Friend{},
		&im_mysql_model.FriendRequest{},
		&im_mysql_model.Group{},
		&im_mysql_model.GroupMember{},
		&im_mysql_model.GroupRequest{},
		&im_mysql_model.User{},
		&im_mysql_model.Black{}, &im_mysql_model.ChatLog{}, &im_mysql_model.Register{}, &im_mysql_model.Conversation{}, &im_mysql_model.AppVersion{}, &im_mysql_model.Department{}, &im_mysql_model.BlackList{}, &im_mysql_model.IpLimit{}, &im_mysql_model.UserIpLimit{}, &im_mysql_model.Invitation{}, &im_mysql_model.RegisterAddFriend{},
		&im_mysql_model.ClientInitConfig{}, &im_mysql_model.UserIpRecord{})
	db.Set("gorm:table_options", "CHARSET=utf8")
	db.Set("gorm:table_options", "collation=utf8_unicode_ci")

	if !db.Migrator().HasTable(&im_mysql_model.Friend{}) {
		db.Migrator().CreateTable(&im_mysql_model.Friend{})
	}
	if !db.Migrator().HasTable(&im_mysql_model.FriendRequest{}) {
		db.Migrator().CreateTable(&im_mysql_model.FriendRequest{})
	}
	if !db.Migrator().HasTable(&im_mysql_model.Group{}) {
		db.Migrator().CreateTable(&im_mysql_model.Group{})
	}
	if !db.Migrator().HasTable(&im_mysql_model.GroupMember{}) {
		db.Migrator().CreateTable(&im_mysql_model.GroupMember{})
	}
	if !db.Migrator().HasTable(&im_mysql_model.GroupRequest{}) {
		db.Migrator().CreateTable(&im_mysql_model.GroupRequest{})
	}
	if !db.Migrator().HasTable(&im_mysql_model.User{}) {
		db.Migrator().CreateTable(&im_mysql_model.User{})
	}
	if !db.Migrator().HasTable(&im_mysql_model.Black{}) {
		db.Migrator().CreateTable(&im_mysql_model.Black{})
	}
	if !db.Migrator().HasTable(&im_mysql_model.ChatLog{}) {
		db.Migrator().CreateTable(&im_mysql_model.ChatLog{})
	}
	if !db.Migrator().HasTable(&im_mysql_model.Register{}) {
		db.Migrator().CreateTable(&im_mysql_model.Register{})
	}
	if !db.Migrator().HasTable(&im_mysql_model.Conversation{}) {
		db.Migrator().CreateTable(&im_mysql_model.Conversation{})
	}
	if !db.Migrator().HasTable(&im_mysql_model.Department{}) {
		db.Migrator().CreateTable(&im_mysql_model.Department{})
	}
	if !db.Migrator().HasTable(&im_mysql_model.OrganizationUser{}) {
		db.Migrator().CreateTable(&im_mysql_model.OrganizationUser{})
	}
	if !db.Migrator().HasTable(&im_mysql_model.DepartmentMember{}) {
		db.Migrator().CreateTable(&im_mysql_model.DepartmentMember{})
	}
	if !db.Migrator().HasTable(&im_mysql_model.AppVersion{}) {
		db.Migrator().CreateTable(&im_mysql_model.AppVersion{})
	}
	if !db.Migrator().HasTable(&im_mysql_model.BlackList{}) {
		db.Migrator().CreateTable(&im_mysql_model.BlackList{})
	}
	if !db.Migrator().HasTable(&im_mysql_model.IpLimit{}) {
		db.Migrator().CreateTable(&im_mysql_model.IpLimit{})
	}
	if !db.Migrator().HasTable(&im_mysql_model.UserIpLimit{}) {
		db.Migrator().CreateTable(&im_mysql_model.UserIpLimit{})
	}
	if !db.Migrator().HasTable(&im_mysql_model.RegisterAddFriend{}) {
		db.Migrator().CreateTable(&im_mysql_model.RegisterAddFriend{})
	}
	if !db.Migrator().HasTable(&im_mysql_model.Invitation{}) {
		db.Migrator().CreateTable(&im_mysql_model.Invitation{})
	}
	if !db.Migrator().HasTable(&im_mysql_model.ClientInitConfig{}) {
		db.Migrator().CreateTable(&im_mysql_model.ClientInitConfig{})
	}
	if !db.Migrator().HasTable(&im_mysql_model.UserIpRecord{}) {
		db.Migrator().CreateTable(&im_mysql_model.UserIpRecord{})
	}
	DB.MysqlDB.db = db
}

func (m *mysqlDB) DefaultGormDB() *gorm.DB {
	return DB.MysqlDB.db
}
