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
		var appMgr db.Users
		appMgr.UserID = v
		appMgr.Nickname = "AppManager" + utils.IntToString(k+1)
		appMgr.AppMangerLevel = constant.AppAdmin
		err = UserRegister(appMgr)
		if err != nil {
			fmt.Println("AppManager insert error", err.Error(), appMgr, "time: ", appMgr.Birth.Unix())
		}

	}
}

func UserRegister(user db.Users) error {
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
	i = dbConn.Table("users").Where("user_id=?", userID).Delete(db.Users{}).RowsAffected
	return i
}

func GetUserByUserID(userID string) (*db.Users, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var user db.Users
	err = dbConn.Table("users").Where("user_id=?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUserInfo(user db.Users) error {
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

func GetUsers(showNumber, pageNumber int32) ([]db.Users, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	var users []db.Users
	if err != nil {
		return users, err
	}
	dbConn.LogMode(true)
	err = dbConn.Limit(showNumber).Offset(showNumber*(pageNumber-1)).Find(&users).Error
	if err != nil {
		return users, err
	}
	return users, err
}

func GetUsersNumCount() (int, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return 0, err
	}
	dbConn.LogMode(true)
	var count int
	if err := dbConn.Model(&db.Users{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func AddUser(userId, phoneNumber, name string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	user := db.Users{
		PhoneNumber:phoneNumber,
		Birth:time.Now(),
		CreateTime:time.Now(),
		UserID: userId,
		Nickname:name,
	}
	result := dbConn.Create(&user)
	return  result.Error
}

func UserIsBlock(userId string) (bool, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return false, err
	}
	var user db.BlackList
	rows := dbConn.Table("black_list").Where("uid=?", userId).First(&user).RowsAffected
	if rows >= 1 {
		return true, nil
	}
	return false, nil
}

func BlockUser(userId, endDisableTime string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	end, err := time.Parse("2006-01-02 15:04:05", endDisableTime)
	if err != nil {
		return err
	}
	if end.Before(time.Now()) {
		return constant.ErrDB
	}
	var user db.BlackList
	dbConn.Table("black_list").Where("uid=?", userId).First(&user)
	if user.UserId != "" {
		dbConn.Model(&user).Where("uid=?", user.UserId).Update("end_disable_time", end)
		return nil
	}
 	user = db.BlackList{
		UserId: userId,
		BeginDisableTime: time.Now(),
		EndDisableTime: end,
	}
	result := dbConn.Create(&user)
	return result.Error
}

func UnBlockUser(userId string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	dbConn.LogMode(true)
	fmt.Println(userId)
	result := dbConn.Where("uid=?", userId).Delete(&db.BlackList{})
	return result.Error
}

func GetBlockUsersID(showNumber, pageNumber int32) ([]string, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	var blockUsers []db.BlackList
	var blockUserIds []string
	if err != nil {
		return blockUserIds, err
	}
	dbConn.LogMode(true)
	err = dbConn.Limit(showNumber).Offset(showNumber*(pageNumber-1)).Find(&blockUsers).Error
	if err != nil {
		return blockUserIds, err
	}
	for _, v := range blockUsers {
		blockUserIds = append(blockUserIds, v.UserId)
	}
	return blockUserIds, err
}

func GetBlockUsers(userIds []string) ([]db.Users, error){
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	var blockUsers []db.Users
	if err != nil {
		return blockUsers, err
	}
	dbConn.LogMode(true)
	dbConn.Find(&blockUsers,userIds)
	return blockUsers, err
}

func GetBlockUsersNumCount() (int, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return 0, err
	}
	dbConn.LogMode(true)
	var count int
	if err := dbConn.Model(&db.BlackList{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}