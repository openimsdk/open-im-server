package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	_ "gorm.io/gorm"
)

func GetRegister(account, areaCode, userID string) (*db.Register, error) {
	var r db.Register
	return &r, db.DB.MysqlDB.DefaultGormDB().Table("registers").Where("user_id = ? and user_id != ? or account = ? or account =? and area_code=?",
		userID, "", account, account, areaCode).Take(&r).Error
}

func SetPassword(account, password, ex, userID, areaCode string) error {
	r := db.Register{
		Account:  account,
		Password: password,
		Ex:       ex,
		UserID:   userID,
		AreaCode: areaCode,
	}
	return db.DB.MysqlDB.DefaultGormDB().Table("registers").Create(&r).Error
}

func ResetPassword(account, password string) error {
	r := db.Register{
		Password: password,
	}
	return db.DB.MysqlDB.DefaultGormDB().Table("registers").Where("account = ?", account).Updates(&r).Error
}
