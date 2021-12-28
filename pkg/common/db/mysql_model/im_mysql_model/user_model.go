package im_mysql_model

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/utils"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

func init() {
	//init managers
	for k, v := range config.Config.Manager.AppManagerUid {
		user, err := GetUserByUserID(v)
		if err != nil {
			fmt.Println("GetUserByUserID failed ", err.Error(), v, user)
			continue
		}
		var appMgr User
		appMgr.Nickname = "AppManager" + utils.IntToString(k+1)
		err = UserRegister(appMgr)
		if err != nil {
			fmt.Println("AppManager insert error", err.Error())
		}

	}
}

func UserRegister(user User) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	user.CreateTime = time.Now()
	if user.Birth != 0 {
		user.Birth = utils.UnixSecondToTime(user.Birth)
	}

	err = dbConn.Table("user").Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

//type User struct {
//	UserID      string    `gorm:"column:user_id;primaryKey;"`
//	Nickname    string    `gorm:"column:name"`
//	FaceUrl     string    `gorm:"column:icon"`
//	Gender      int32     `gorm:"column:gender"`
//	PhoneNumber string    `gorm:"column:phone_number"`
//	Birth       string    `gorm:"column:birth"`
//	Email       string    `gorm:"column:email"`
//	Ex          string    `gorm:"column:ex"`
//	CreateTime  time.Time `gorm:"column:create_time"`
//}

func DeleteUser(userID string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	err = dbConn.Table("user").Where("user_id=?", userID).Delete(User{}).Error
	return err
}

func GetUserByUserID(userID string) (*User, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var user User
	err = dbConn.Table("user").Where("user_id=?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUserInfo(user User) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	if user.Birth != 0 {
		user.Birth = utils.UnixSecondToTime(user.Birth)
	}
	err = dbConn.Table("user").Where("user_id=?", user.UserID).Update(&user).Error
	return err
}

func SelectAllUserID() ([]string, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var resultArr []string
	err = dbConn.Table("user").Select([]string{"user_id"}).Scan(&resultArr).Error
	if err != nil {
		return nil, err
	}
	return resultArr, nil
}

func SelectSomeUserID(userIDList []string) ([]string, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var userList []User
	err = dbConn.Table("user").Where("(user_id) IN ? ", userIDList).Find(&userList).Error
	if err != nil {
		return nil, err
	}
	return userIDList, nil
}
