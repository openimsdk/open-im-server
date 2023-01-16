package mysql

import (
	"gorm.io/gorm"
)

var RegisterDB *gorm.DB

type Register struct {
	Account        string `gorm:"column:account;primary_key;type:char(255)" json:"account"`
	Password       string `gorm:"column:password;type:varchar(255)" json:"password"`
	Ex             string `gorm:"column:ex;size:1024" json:"ex"`
	UserID         string `gorm:"column:user_id;type:varchar(255)" json:"userID"`
	AreaCode       string `gorm:"column:area_code;type:varchar(255)"`
	InvitationCode string `gorm:"column:invitation_code;type:varchar(255)"`
	RegisterIP     string `gorm:"column:register_ip;type:varchar(255)"`
}

func GetRegister(account, areaCode, userID string) (*Register, error) {
	var r Register
	return &r, RegisterDB.Table("registers").Where("user_id = ? and user_id != ? or account = ? or account =? and area_code=?",
		userID, "", account, account, areaCode).Take(&r).Error
}

func SetPassword(account, password, ex, userID, areaCode, ip string) error {
	r := Register{
		Account:    account,
		Password:   password,
		Ex:         ex,
		UserID:     userID,
		RegisterIP: ip,
		AreaCode:   areaCode,
	}
	return RegisterDB.Table("registers").Create(&r).Error
}

func ResetPassword(account, password string) error {
	r := Register{
		Password: password,
	}
	return RegisterDB.Table("registers").Where("account = ?", account).Updates(&r).Error
}
