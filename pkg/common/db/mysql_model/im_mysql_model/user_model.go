package im_mysql_model

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db"
	pbAuth "Open_IM/pkg/proto/auth"
	"Open_IM/pkg/utils"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

func init() {
	//init managers
	var pb pbAuth.UserRegisterReq
	for k, v := range config.Config.Manager.AppManagerUid {
		if !IsExistUser(v) {
			pb.UID = v
			pb.Name = "AppManager" + utils.IntToString(k+1)
			err := UserRegister(&pb)
			if err != nil {
				fmt.Println("AppManager insert error", err.Error())
			}
		}
	}
}
func UserRegister(pb *pbAuth.UserRegisterReq) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	addUser := User{
		UID:        pb.UID,
		Name:       pb.Name,
		Icon:       pb.Icon,
		Gender:     pb.Gender,
		Mobile:     pb.Mobile,
		Birth:      pb.Birth,
		Email:      pb.Email,
		Ex:         pb.Ex,
		CreateTime: time.Now(),
	}
	err = dbConn.Table("user").Create(&addUser).Error
	if err != nil {
		return err
	}
	return nil
}

func UserDelete(uid string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	err = dbConn.Table("user").Where("uid=?", uid).Delete(User{}).Error
	if err != nil {
		return err
	}
	return nil
}
func FindUserByUID(uid string) (*User, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var user User
	err = dbConn.Table("user").Where("uid=?", uid).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpDateUserInfo(uid, name, headUrl, mobilePhoneNum, birth, email, extendInfo string, gender int32) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	if name != "" {
		if err = dbConn.Exec("update user set name=? where uid=?", name, uid).Error; err != nil {
			return err
		}
	}
	if headUrl != "" {
		if err = dbConn.Exec("update user set icon=? where uid=?", headUrl, uid).Error; err != nil {
			return err
		}
	}
	if mobilePhoneNum != "" {
		if err = dbConn.Exec("update user set mobile=? where uid=?", mobilePhoneNum, uid).Error; err != nil {
			return err
		}
	}
	if birth != "" {
		if err = dbConn.Exec("update user set birth=? where uid=?", birth, uid).Error; err != nil {
			return err
		}
	}
	if email != "" {
		if err = dbConn.Exec("update user set email=? where uid=?", email, uid).Error; err != nil {
			return err
		}
	}
	if extendInfo != "" {
		if err = dbConn.Exec("update user set ex=? where uid=?", extendInfo, uid).Error; err != nil {
			return err
		}
	}
	if gender != 0 {
		if err = dbConn.Exec("update user set gender=? where uid=?", gender, uid).Error; err != nil {
			return err
		}
	}
	return nil
}

func SelectAllUID() ([]string, error) {
	var uid []string

	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return uid, err
	}
	rows, _ := dbConn.Raw("select uid from user").Rows()
	defer rows.Close()
	var strUID string
	for rows.Next() {
		rows.Scan(&strUID)
		uid = append(uid, strUID)
	}
	return uid, nil
}

func IsExistUser(uid string) bool {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return false
	}
	var number int32
	err = dbConn.Raw("select count(*) from `user` where  uid = ?", uid).Count(&number).Error
	if err != nil {
		return false
	}
	if number != 1 {
		return false
	}
	return true
}
