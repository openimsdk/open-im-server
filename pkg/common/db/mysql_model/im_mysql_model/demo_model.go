package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"errors"
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

func GetRegisterAddFriendList(showNumber, pageNumber int32) ([]string, error) {
	var IDList []string
	var err error
	model := db.DB.MysqlDB.DefaultGormDB().Model(&db.RegisterAddFriend{})
	if showNumber == 0 {
		err = model.Pluck("user_id", &IDList).Error
	} else {
		err = model.Limit(int(showNumber)).Offset(int(showNumber*(pageNumber-1))).Pluck("user_id", &IDList).Error
	}
	return IDList, err
}

func AddUserRegisterAddFriendIDList(userIDList ...string) error {
	var list []db.RegisterAddFriend
	for _, v := range userIDList {
		list = append(list, db.RegisterAddFriend{UserID: v})
	}
	result := db.DB.MysqlDB.DefaultGormDB().Create(list)
	if int(result.RowsAffected) < len(userIDList) {
		return errors.New("some line insert failed")
	}
	err := result.Error
	return err
}

func ReduceUserRegisterAddFriendIDList(userIDList ...string) error {
	var list []db.RegisterAddFriend
	for _, v := range userIDList {
		list = append(list, db.RegisterAddFriend{UserID: v})
	}
	err := db.DB.MysqlDB.DefaultGormDB().Delete(list).Error
	return err
}

func DeleteAllRegisterAddFriendIDList() error {
	err := db.DB.MysqlDB.DefaultGormDB().Where("1 = 1").Delete(&db.RegisterAddFriend{}).Error
	return err
}
