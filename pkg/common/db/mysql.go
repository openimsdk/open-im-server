package db

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"sync"
	"time"
)

type mysqlDB struct {
	sync.RWMutex
	dbMap map[string]*gorm.DB
}

func initMysqlDB() {
	//When there is no open IM database, connect to the mysql built-in database to create openIM database
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.Config.Mysql.DBUserName, config.Config.Mysql.DBPassword, config.Config.Mysql.DBAddress[0], "mysql")

	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Error("", "", dsn)
		panic(err)
	}

	//Check the database and table during initialization
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s ;", config.Config.Mysql.DBDatabaseName)
	err = db.Exec(sql).Error
	if err != nil {
		panic(err)
	}
	db.Close()

	dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.Config.Mysql.DBUserName, config.Config.Mysql.DBPassword, config.Config.Mysql.DBAddress[0], config.Config.Mysql.DBDatabaseName)
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	sqlTable := "CREATE TABLE IF NOT EXISTS `user` ("+
		" `uid` varchar(64) NOT NULL,"+
		" `name` varchar(64) DEFAULT NULL,"+
		" `icon` varchar(1024) DEFAULT NULL,"+
		" `gender` tinyint(4) unsigned zerofill DEFAULT NULL,"+
		" `mobile` varchar(32) DEFAULT NULL,"+
		" `birth` varchar(16) DEFAULT NULL,"+
		" `email` varchar(64) DEFAULT NULL,"+
		" `ex` varchar(1024) DEFAULT NULL,"+
		" `create_time` datetime DEFAULT NULL,"+
		" PRIMARY KEY (`uid`),"+
		" UNIQUE KEY `uk_uid` (`uid`)"+
		" ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;"
	err = db.Exec(sqlTable).Error
	if err != nil {
		panic(err)
	}

	sqlTable = "CREATE TABLE IF NOT EXISTS `friend` ("+
		" `owner_id` varchar(64) NOT NULL,"+
		" `friend_id` varchar(64) NOT NULL,"+
		" `comment` varchar(255) DEFAULT NULL,"+
		" `friend_flag` int(11) NOT NULL,"+
		" `create_time` datetime NOT NULL,"+
		" PRIMARY KEY (`owner_id`,`friend_id`) USING BTREE"+
		" ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;"
	err = db.Exec(sqlTable).Error
	if err != nil {
		panic(err)
	}

	sqlTable = "CREATE TABLE IF NOT EXISTS  `friend_request` ("+
		" `req_id` varchar(64) NOT NULL,"+
		" `user_id` varchar(64) NOT NULL,"+
		" `flag` int(11) NOT NULL DEFAULT '0',"+
		" `req_message` varchar(255) DEFAULT NULL,"+
		" `create_time` datetime NOT NULL,"+
		" PRIMARY KEY (`user_id`,`req_id`) USING BTREE"+
		" ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;"
	err = db.Exec(sqlTable).Error
	if err != nil {
		panic(err)
	}

	sqlTable = "CREATE TABLE IF NOT EXISTS `black_list` ("+
		" `uid` varchar(32) NOT NULL COMMENT 'uid',"+
		" `begin_disable_time` datetime DEFAULT NULL,"+
		" `end_disable_time` datetime DEFAULT NULL,"+
		" `ex` varchar(1024) DEFAULT NULL,"+
		" PRIMARY KEY (`uid`) USING BTREE"+
		" ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;"
	err = db.Exec(sqlTable).Error
	if err != nil {
		panic(err)
	}

	sqlTable = "CREATE TABLE IF NOT EXISTS `user_black_list` ("+
		" `owner_id` varchar(64) NOT NULL,"+
		" `block_id` varchar(64) NOT NULL,"+
		" `create_time` datetime NOT NULL,"+
		" PRIMARY KEY (`owner_id`,`block_id`) USING BTREE"+
		" ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;"
	err = db.Exec(sqlTable).Error
	if err != nil {
		panic(err)
	}

	sqlTable = "CREATE TABLE IF NOT EXISTS `group` ("+
		" `group_id` varchar(64) NOT NULL,"+
		" `name` varchar(255) DEFAULT NULL,"+
		" `introduction` varchar(255) DEFAULT NULL,"+
		" `notification` varchar(255) DEFAULT NULL,"+
		" `face_url` varchar(255) DEFAULT NULL,"+
		" `create_time` datetime DEFAULT NULL,"+
		" `ex` varchar(255) DEFAULT NULL,"+
		" PRIMARY KEY (`group_id`)"+
		" ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;"
	err = db.Exec(sqlTable).Error
	if err != nil {
		panic(err)
	}

	sqlTable = "CREATE TABLE IF NOT EXISTS `group_member` ("+
		" `group_id` varchar(64) NOT NULL,"+
		" `uid` varchar(64) NOT NULL,"+
		" `nickname` varchar(255) DEFAULT NULL,"+
		" `user_group_face_url` varchar(255) DEFAULT NULL,"+
		" `administrator_level` int(11) NOT NULL,"+
		" `join_time` datetime NOT NULL,"+
		"  PRIMARY KEY (`group_id`,`uid`) USING BTREE"+
		" ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;"
	err = db.Exec(sqlTable).Error
	if err != nil {
		panic(err)
	}

	sqlTable = "CREATE TABLE IF NOT EXISTS `group_request` ("+
		" `id` int(11) NOT NULL AUTO_INCREMENT,"+
		" `group_id` varchar(64) NOT NULL,"+
		" `from_user_id` varchar(255) NOT NULL,"+
		" `to_user_id` varchar(255) NOT NULL,"+
		" `flag` int(10) NOT NULL DEFAULT '0',"+
		" `req_msg` varchar(255) DEFAULT '',"+
		" `handled_msg` varchar(255) DEFAULT '',"+
		" `create_time` datetime NOT NULL,"+
		" `from_user_nickname` varchar(255) DEFAULT '',"+
		" `to_user_nickname` varchar(255) DEFAULT NULL,"+
		" `from_user_face_url` varchar(255) DEFAULT '',"+
		" `to_user_face_url` varchar(255) DEFAULT '',"+
		" `handled_user` varchar(255) DEFAULT '',"+
		" PRIMARY KEY (`id`)"+
		" ) ENGINE=InnoDB AUTO_INCREMENT=38 DEFAULT CHARSET=utf8mb4;"
	err = db.Exec(sqlTable).Error
	if err != nil {
		panic(err)
	}

	sqlTable = "CREATE TABLE IF NOT EXISTS  `chat_log` ("+
		" `msg_id` varchar(128) NOT NULL,"+
		" `send_id` varchar(255) NOT NULL,"+
		" `session_type` int(11) NOT NULL,"+
		" `recv_id` varchar(255) NOT NULL,"+
		" `content_type` int(11) NOT NULL,"+
		" `msg_from` int(11) NOT NULL,"+
		" `content` varchar(1000) NOT NULL,"+
		" `remark` varchar(100) DEFAULT NULL,"+
		" `sender_platform_id` int(11) NOT NULL,"+
		" `send_time` datetime NOT NULL,"+
		" PRIMARY KEY (`msg_id`) USING BTREE"+
		" ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;"
	err = db.Exec(sqlTable).Error
	if err != nil {
		panic(err)
	}

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
