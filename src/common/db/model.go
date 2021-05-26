package db

import (
	"gopkg.in/mgo.v2"
)

var DB DataBases

type DataBases struct {
	MgoDB   mongoDB
	RedisDB redisDB
	MysqlDB mysqlDB
}

func key(dbAddress, dbName string) string {
	return dbAddress + "_" + dbName
}

//type Config struct {
//	Mongo struct {
//		DBAddress     []string `yaml:"dbAddress"`
//		DBDirect      bool     `yaml:"dbDirect"`
//		DBTimeout     int      `yaml:"dbTimeout"`
//		DBDatabase    []string `yaml:"dbDatabase"`
//		DBSource      string   `yaml:"dbSource"`
//		DBUserName    string   `yaml:"dbUserName"`
//		DBPassword    string   `yaml:"dbPassword"`
//		DBMaxPoolSize int      `yaml:"dbMaxPoolSize"`
//	}
//	Mysql struct {
//		DBAddress      []string `yaml:"dbAddress"`
//		DBPort         int      `yaml:"dbPort"`
//		DBUserName     string   `yaml:"dbUserName"`
//		DBPassword     string   `yaml:"dbPassword"`
//		DBDatabaseName     string   `yaml:"dbChatName"` // 默认使用DBAddress[0]
//		DBTableName      string   `yaml:"dbMsgName"`
//		DBMsgTableNum  int      `yaml:"dbMsgTableNum"`
//		DBCharset      string   `yaml:"dbCharset"`
//		DBMaxOpenConns int      `yaml:"dbMaxOpenConns"`
//		DBMaxIdleConns int      `yaml:"dbMaxIdleConns"`
//		DBMaxLifeTime  int      `yaml:"dbMaxLifeTime"`
//	}
//	Redis struct {
//		DBAddress     string `yaml:"dbAddress"`
//		DBPort        int    `yaml:"dbPort"`
//		DBMaxIdle     int    `yaml:"dbMaxIdle"`
//		DBMaxActive   int    `yaml:"dbMaxActive"`
//		DBIdleTimeout int    `yaml:"dbIdleTimeout"`
//	}
//}

//func init() {
//	bytes, err := ioutil.ReadFile("config/db.yaml")
//	if err != nil {
//		log.Error("", "", "read db.yaml config fail! err = %s", err.Error())
//		return
//	}
//
//	if err = yaml.Unmarshal(bytes, &DB.Config); err != nil {
//		log.Error("", "", "unmarshal db.yaml config fail! err = %s", err.Error())
//		return
//	}
//
//	DB.RedisDB.newPool(DB.Config)
//	//DB.MysqlDB.sqlxDB(DB.Config.Mysql.DBName[0], DB.Config)
//}
func init() {
	DB.RedisDB.newPool()
}
func (d *DataBases) session(dbName string) *mgo.Session {
	return d.MgoDB.mgoSession(dbName)
}
