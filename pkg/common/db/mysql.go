package db

import (
	"Open_IM/pkg/common/config"
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
		&Register{},
		&Friend{},
		&FriendRequest{},
		&Group{},
		&GroupMember{},
		&GroupRequest{},
		&User{},
		&Black{}, &ChatLog{}, &Register{}, &Conversation{}, &AppVersion{}, &Department{}, &BlackList{}, &IpLimit{}, &UserIpLimit{}, &Invitation{}, &RegisterAddFriend{},
		&ClientInitConfig{}, &UserIpRecord{})
	db.Set("gorm:table_options", "CHARSET=utf8")
	db.Set("gorm:table_options", "collation=utf8_unicode_ci")

	if !db.Migrator().HasTable(&Friend{}) {
		db.Migrator().CreateTable(&Friend{})
	}
	if !db.Migrator().HasTable(&FriendRequest{}) {
		db.Migrator().CreateTable(&FriendRequest{})
	}
	if !db.Migrator().HasTable(&Group{}) {
		db.Migrator().CreateTable(&Group{})
	}
	if !db.Migrator().HasTable(&GroupMember{}) {
		db.Migrator().CreateTable(&GroupMember{})
	}
	if !db.Migrator().HasTable(&GroupRequest{}) {
		db.Migrator().CreateTable(&GroupRequest{})
	}
	if !db.Migrator().HasTable(&User{}) {
		db.Migrator().CreateTable(&User{})
	}
	if !db.Migrator().HasTable(&Black{}) {
		db.Migrator().CreateTable(&Black{})
	}
	if !db.Migrator().HasTable(&ChatLog{}) {
		db.Migrator().CreateTable(&ChatLog{})
	}
	if !db.Migrator().HasTable(&Register{}) {
		db.Migrator().CreateTable(&Register{})
	}
	if !db.Migrator().HasTable(&Conversation{}) {
		db.Migrator().CreateTable(&Conversation{})
	}
	if !db.Migrator().HasTable(&Department{}) {
		db.Migrator().CreateTable(&Department{})
	}
	if !db.Migrator().HasTable(&OrganizationUser{}) {
		db.Migrator().CreateTable(&OrganizationUser{})
	}
	if !db.Migrator().HasTable(&DepartmentMember{}) {
		db.Migrator().CreateTable(&DepartmentMember{})
	}
	if !db.Migrator().HasTable(&AppVersion{}) {
		db.Migrator().CreateTable(&AppVersion{})
	}
	if !db.Migrator().HasTable(&BlackList{}) {
		db.Migrator().CreateTable(&BlackList{})
	}
	if !db.Migrator().HasTable(&IpLimit{}) {
		db.Migrator().CreateTable(&IpLimit{})
	}
	if !db.Migrator().HasTable(&UserIpLimit{}) {
		db.Migrator().CreateTable(&UserIpLimit{})
	}
	if !db.Migrator().HasTable(&RegisterAddFriend{}) {
		db.Migrator().CreateTable(&RegisterAddFriend{})
	}
	if !db.Migrator().HasTable(&Invitation{}) {
		db.Migrator().CreateTable(&Invitation{})
	}
	if !db.Migrator().HasTable(&ClientInitConfig{}) {
		db.Migrator().CreateTable(&ClientInitConfig{})
	}
	if !db.Migrator().HasTable(&UserIpRecord{}) {
		db.Migrator().CreateTable(&UserIpRecord{})
	}
	DB.MysqlDB.db = db
}

func (m *mysqlDB) DefaultGormDB() *gorm.DB {
	return DB.MysqlDB.db
}
