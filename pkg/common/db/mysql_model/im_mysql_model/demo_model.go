package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	_ "github.com/jinzhu/gorm"
)

type GetRegisterParams struct {
	Account string `json:"account"`
}
type Register struct {
	Account  string `gorm:"column:account"`
	Password string `gorm:"column:password"`
}

func GetRegister(params *GetRegisterParams) (Register, error, int64) {
	var r Register
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return r, err, 0
	}
	result := dbConn.
		Where(&Register{Account: params.Account}).
		Find(&r)
	return r, result.Error, result.RowsAffected
}

type SetPasswordParams struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

func SetPassword(params *SetPasswordParams) (Register, error) {
	r := Register{
		Account:  params.Account,
		Password: params.Password,
	}
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return r, err
	}

	result := dbConn.Create(&r)

	return r, result.Error
}

func Login(params *Register) int64 {
	var r Register
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return 3
	}
	result := dbConn.
		Where(&Register{Account: params.Account}).
		Find(&r)
	if result.Error != nil && result.RowsAffected == 0 {
		return 1
	}
	if r.Password != params.Password {
		return 2
	}
	return 0
}
