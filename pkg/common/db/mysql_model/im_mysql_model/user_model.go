package im_mysql_model

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/utils"
	"fmt"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func init() {
	//init managers
	for k, v := range config.Config.Manager.AppManagerUid {
		user, err := GetUserByUserID(v)
		if err != nil {
			fmt.Println("GetUserByUserID failed ", err.Error(), v, user)
		} else {
			continue
		}
		var appMgr db.User
		appMgr.UserID = v
		appMgr.Nickname = "AppManager" + utils.IntToString(k+1)
		appMgr.AppMangerLevel = constant.AppAdmin
		err = UserRegister(appMgr)
		if err != nil {
			fmt.Println("AppManager insert error", err.Error(), appMgr, "time: ", appMgr.Birth.Unix())
		}

	}
}

func UserRegister(user db.User) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	user.CreateTime = time.Now()
	if user.AppMangerLevel == 0 {
		user.AppMangerLevel = constant.AppOrdinaryUsers
	}
	if user.Birth.Unix() < 0 {
		user.Birth = utils.UnixSecondToTime(0)
	}
	err = dbConn.Table("users").Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

type User struct {
	UserID      string    `gorm:"column:user_id;primaryKey;"`
	Nickname    string    `gorm:"column:name"`
	FaceUrl     string    `gorm:"column:icon"`
	Gender      int32     `gorm:"column:gender"`
	PhoneNumber string    `gorm:"column:phone_number"`
	Birth       string    `gorm:"column:birth"`
	Email       string    `gorm:"column:email"`
	Ex          string    `gorm:"column:ex"`
	CreateTime  time.Time `gorm:"column:create_time"`
}

func DeleteUser(userID string) (i int64) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return 0
	}
	i = dbConn.Table("users").Where("user_id=?", userID).Delete(db.User{}).RowsAffected
	return i
}

func GetUserByUserID(userID string) (*db.User, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var user db.User
	err = dbConn.Table("users").Where("user_id=?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUserInfo(user db.User) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}

	err = dbConn.Table("users").Where("user_id=?", user.UserID).Update(&user).Error
	return err
}

func SelectAllUserID() ([]string, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var resultArr []string
	err = dbConn.Table("users").Pluck("user_id", &resultArr).Error
	if err != nil {
		return nil, err
	}
	return resultArr, nil
}

func SelectSomeUserID(userIDList []string) ([]string, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	dbConn.LogMode(true)
	if err != nil {
		return nil, err
	}
	var resultArr []string
	err = dbConn.Table("users").Where("user_id IN (?) ", userIDList).Pluck("user_id", &resultArr).Error

	if err != nil {
		return nil, err
	}
	return resultArr, nil
}

func GetUsers(showNumber, pageNumber int32) ([]User, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	dbConn.LogMode(true)
	var users []User
	if err != nil {
		return users, err
	}
	err = dbConn.Limit(showNumber).Offset(pageNumber).Find(&users).Error
	if err != nil {
		return users, err
	}
	return users, nil
}
