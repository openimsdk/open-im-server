package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	_ "github.com/jinzhu/gorm"
)

func GetRegister(account string) (*db.Register, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var r db.Register
	return &r, dbConn.Table("registers").Where("account = ?",
		account).Take(&r).Error
}

func SetPassword(account, password, ex string) error {
	r := db.Register{
		Account:  account,
		Password: password,
		Ex:       ex,
	}
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	return dbConn.Table("registers").Create(&r).Error
}

func ResetPassword(account, password string) error {
	r := db.Register{
		Password: password,
	}
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	dbConn.LogMode(true)
	if err != nil {
		return err
	}
	return dbConn.Table("registers").Where("account = ?", account).Update(&r).Error
}
