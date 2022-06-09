package db

import (
	"Open_IM/pkg/common/config"

	"fmt"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type mysqlDB struct {
	sync.RWMutex
	dbMap map[string]*gorm.DB
}

func initMysqlDB() {
	//When there is no open IM database, connect to the mysql built-in database to create openIM database
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.Config.Mysql.DBUserName, config.Config.Mysql.DBPassword, config.Config.Mysql.DBAddress[0], "mysql")
	var db *gorm.DB
	var err1 error
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		fmt.Println("0", "Open failed ", err.Error(), dsn)
	}
	if err != nil {
		time.Sleep(time.Duration(30) * time.Second)
		db, err1 = gorm.Open("mysql", dsn)
		if err1 != nil {
			fmt.Println("0", "Open failed ", err1.Error(), dsn)
			panic(err1.Error())
		}
	}

	//Check the database and table during initialization
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s default charset utf8 COLLATE utf8_general_ci;", config.Config.Mysql.DBDatabaseName)
	err = db.Exec(sql).Error
	if err != nil {
		fmt.Println("0", "Exec failed ", err.Error(), sql)
		panic(err.Error())
	}
	db.Close()

	dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.Config.Mysql.DBUserName, config.Config.Mysql.DBPassword, config.Config.Mysql.DBAddress[0], config.Config.Mysql.DBDatabaseName)
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		fmt.Println("0", "Open failed ", err.Error(), dsn)
		panic(err.Error())
	}

	fmt.Println("open db ok ", dsn)
	db.AutoMigrate(
		&Register{},
		&Friend{},
		&FriendRequest{},
		&Group{},
		&GroupMember{},
		&GroupRequest{},
		&User{},
		&Black{}, &ChatLog{}, &Register{}, &Conversation{}, &AppVersion{}, &Department{})
	db.Set("gorm:table_options", "CHARSET=utf8")
	db.Set("gorm:table_options", "collation=utf8_unicode_ci")

	if !db.HasTable(&Friend{}) {
		fmt.Println("CreateTable Friend")
		db.CreateTable(&Friend{})
	}

	if !db.HasTable(&FriendRequest{}) {
		fmt.Println("CreateTable FriendRequest")
		db.CreateTable(&FriendRequest{})
	}

	if !db.HasTable(&Group{}) {
		fmt.Println("CreateTable Group")
		db.CreateTable(&Group{})
	}

	if !db.HasTable(&GroupMember{}) {
		fmt.Println("CreateTable GroupMember")
		db.CreateTable(&GroupMember{})
	}
	if !db.HasTable(&GroupRequest{}) {
		fmt.Println("CreateTable GroupRequest")
		db.CreateTable(&GroupRequest{})
	}
	if !db.HasTable(&User{}) {
		fmt.Println("CreateTable User")
		db.CreateTable(&User{})
	}
	if !db.HasTable(&Black{}) {
		fmt.Println("CreateTable Black")
		db.CreateTable(&Black{})
	}
	if !db.HasTable(&ChatLog{}) {
		fmt.Println("CreateTable ChatLog")
		db.CreateTable(&ChatLog{})
	}
	if !db.HasTable(&Register{}) {
		fmt.Println("CreateTable Register")
		db.CreateTable(&Register{})
	}
	if !db.HasTable(&Conversation{}) {
		fmt.Println("CreateTable Conversation")
		db.CreateTable(&Conversation{})
	}

	if !db.HasTable(&Department{}) {
		fmt.Println("CreateTable Department")
		db.CreateTable(&Department{})
	}
	if !db.HasTable(&OrganizationUser{}) {
		fmt.Println("CreateTable OrganizationUser")
		db.CreateTable(&OrganizationUser{})
	}
	if !db.HasTable(&DepartmentMember{}) {
		fmt.Println("CreateTable DepartmentMember")
		db.CreateTable(&DepartmentMember{})
	}
	if !db.HasTable(&AppVersion{}) {
		fmt.Println("CreateTable DepartmentMember")
		db.CreateTable(&AppVersion{})
	}
	return
}

func (m *mysqlDB) DefaultGormDB() (*gorm.DB, error) {
	return m.GormDB(config.Config.Mysql.DBAddress[0], config.Config.Mysql.DBDatabaseName)
}

func (m *mysqlDB) GormDB(dbAddress, dbName string) (*gorm.DB, error) {
	m.Lock()
	defer m.Unlock()

	k := key(dbAddress, dbName)
	if _, ok := m.dbMap[k]; !ok {
		if err := m.open(dbAddress, dbName); err != nil {
			return nil, err
		}
	}
	return m.dbMap[k], nil
}

func (m *mysqlDB) open(dbAddress, dbName string) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.Config.Mysql.DBUserName, config.Config.Mysql.DBPassword, dbAddress, dbName)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return err
	}

	db.SingularTable(true)
	db.DB().SetMaxOpenConns(config.Config.Mysql.DBMaxOpenConns)
	db.DB().SetMaxIdleConns(config.Config.Mysql.DBMaxIdleConns)
	db.DB().SetConnMaxLifetime(time.Duration(config.Config.Mysql.DBMaxLifeTime) * time.Second)

	if m.dbMap == nil {
		m.dbMap = make(map[string]*gorm.DB)
	}
	k := key(dbAddress, dbName)
	m.dbMap[k] = db
	return nil
}
