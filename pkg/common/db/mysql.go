package db

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/utils"
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
		&im_mysql_model.Friend{},
		&im_mysql_model.FriendRequest{},
		&im_mysql_model.Group{},
		&im_mysql_model.GroupMember{},
		&im_mysql_model.GroupRequest{},
		&im_mysql_model.User{},
		&im_mysql_model.Black{}, &im_mysql_model.ChatLog{}, &im_mysql_model.Conversation{}, &im_mysql_model.AppVersion{}, &im_mysql_model.BlackList{},
	)
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

	DB.MysqlDB.db = db
	im_mysql_model.GroupDB = db.Table("groups")
	im_mysql_model.GroupMemberDB = db.Table("group_members")
	im_mysql_model.UserDB = db.Table("users")
	im_mysql_model.ChatLogDB = db.Table("chat_logs")
	im_mysql_model.BlackListDB = db.Table("black_lists")
	im_mysql_model.BlackDB = db.Table("blacks")
	im_mysql_model.AppDB = db.Table("app_version")
	im_mysql_model.BlackDB = db.Table("blacks")
	im_mysql_model.ConversationDB = db.Table("conversations")
	im_mysql_model.FriendDB = db.Table("friends")
	im_mysql_model.FriendRequestDB = db.Table("friend_requests")
	im_mysql_model.GroupRequestDB = db.Table("group_requests")
	InitManager()
}

func InitManager() {
	for k, v := range config.Config.Manager.AppManagerUid {
		_, err := im_mysql_model.GetUserByUserID(v)
		if err != nil {
		} else {
			continue
		}
		var appMgr im_mysql_model.User
		appMgr.UserID = v
		if k == 0 {
			appMgr.Nickname = config.Config.Manager.AppSysNotificationName
		} else {
			appMgr.Nickname = "AppManager" + utils.IntToString(k+1)
		}
		appMgr.AppMangerLevel = constant.AppAdmin
		err = im_mysql_model.UserRegister(appMgr)
		if err != nil {
			fmt.Println("AppManager insert error ", err.Error(), appMgr)
		} else {
			fmt.Println("AppManager insert ", appMgr)
		}
	}
}

func (m *mysqlDB) DefaultGormDB() *gorm.DB {
	return DB.MysqlDB.db
}
