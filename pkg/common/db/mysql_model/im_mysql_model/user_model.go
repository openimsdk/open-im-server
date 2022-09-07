package im_mysql_model

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/utils"
	"errors"
	"fmt"
	"strconv"
	"time"
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
		if k == 0 {
			appMgr.Nickname = config.Config.Manager.AppSysNotificationName
		} else {
			appMgr.Nickname = "AppManager" + utils.IntToString(k+1)
		}
		appMgr.AppMangerLevel = constant.AppAdmin
		err = UserRegister(appMgr)
		if err != nil {
			fmt.Println("AppManager insert error", err.Error(), appMgr, "time: ", appMgr.Birth.Unix())
		}

	}
}

func UserRegister(user db.User) error {
	user.CreateTime = time.Now()
	if user.AppMangerLevel == 0 {
		user.AppMangerLevel = constant.AppOrdinaryUsers
	}
	if user.Birth.Unix() < 0 {
		user.Birth = utils.UnixSecondToTime(0)
	}
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func GetAllUser() ([]db.User, error) {
	var userList []db.User
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").Find(&userList).Error
	return userList, err
}

func GetUserByUserID(userID string) (*db.User, error) {
	var user db.User
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").Where("user_id=?", userID).Take(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUsersByUserIDList(userIDList []string) ([]*db.User, error) {
	var userList []*db.User
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").Where("user_id in (?)", userIDList).Find(&userList).Error
	return userList, err
}

func GetUserNameByUserID(userID string) (string, error) {
	var user db.User
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").Select("name").Where("user_id=?", userID).First(&user).Error
	if err != nil {
		return "", err
	}
	return user.Nickname, nil
}

func UpdateUserInfo(user db.User) error {
	return db.DB.MysqlDB.DefaultGormDB().Table("users").Where("user_id=?", user.UserID).Updates(&user).Error
}

func UpdateUserInfoByMap(user db.User, m map[string]interface{}) error {
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").Where("user_id=?", user.UserID).Updates(m).Error
	return err
}

func SelectAllUserID() ([]string, error) {
	var resultArr []string
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").Pluck("user_id", &resultArr).Error
	if err != nil {
		return nil, err
	}
	return resultArr, nil
}

func SelectSomeUserID(userIDList []string) ([]string, error) {
	var resultArr []string
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").Where("user_id IN (?) ", userIDList).Pluck("user_id", &resultArr).Error
	if err != nil {
		return nil, err
	}
	return resultArr, nil
}

func GetUsers(showNumber, pageNumber int32) ([]db.User, error) {
	var users []db.User
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").Limit(int(showNumber)).Offset(int(showNumber * (pageNumber - 1))).Find(&users).Error
	if err != nil {
		return users, err
	}
	return users, err
}

func AddUser(userID string, phoneNumber string, name string, email string, gender int32, faceURL string, birth uint32) error {
	_birth, _err := time.ParseInLocation("2006-01-02", strconv.Itoa(int(birth)), time.Local)
	if _err != nil {
		_birth = time.Now()
	}
	user := db.User{
		UserID:      userID,
		Nickname:    name,
		FaceURL:     faceURL,
		Gender:      gender,
		PhoneNumber: phoneNumber,
		Birth:       _birth,
		Email:       email,
		Ex:          "",
		CreateTime:  time.Now(),
	}
	result := db.DB.MysqlDB.DefaultGormDB().Table("users").Create(&user)
	return result.Error
}

func UserIsBlock(userId string) (bool, error) {
	var user db.BlackList
	rows := db.DB.MysqlDB.DefaultGormDB().Table("black_lists").Where("uid=?", userId).First(&user).RowsAffected
	if rows >= 1 {
		return user.EndDisableTime.After(time.Now()), nil
	}
	return false, nil
}

func UsersIsBlock(userIDList []string) (inBlockUserIDList []string, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("black_lists").Where("uid in (?) and end_disable_time > now()", userIDList).Pluck("uid", &inBlockUserIDList).Error
	return inBlockUserIDList, err
}

func BlockUser(userID, endDisableTime string) error {
	user, err := GetUserByUserID(userID)
	if err != nil || user.UserID == "" {
		return err
	}
	end, err := time.Parse("2006-01-02 15:04:05", endDisableTime)
	if err != nil {
		return err
	}
	if end.Before(time.Now()) {
		return errors.New("endDisableTime is before now")
	}
	var blockUser db.BlackList
	db.DB.MysqlDB.DefaultGormDB().Table("black_lists").Where("uid=?", userID).First(&blockUser)
	if blockUser.UserId != "" {
		db.DB.MysqlDB.DefaultGormDB().Model(&blockUser).Where("uid=?", blockUser.UserId).Update("end_disable_time", end)
		return nil
	}
	blockUser = db.BlackList{
		UserId:           userID,
		BeginDisableTime: time.Now(),
		EndDisableTime:   end,
	}
	err = db.DB.MysqlDB.DefaultGormDB().Create(&blockUser).Error
	return err
}

func UnBlockUser(userID string) error {
	return db.DB.MysqlDB.DefaultGormDB().Where("uid=?", userID).Delete(&db.BlackList{}).Error
}

type BlockUserInfo struct {
	User             db.User
	BeginDisableTime time.Time
	EndDisableTime   time.Time
}

func GetBlockUserByID(userId string) (BlockUserInfo, error) {
	var blockUserInfo BlockUserInfo
	blockUser := db.BlackList{
		UserId: userId,
	}
	if err := db.DB.MysqlDB.DefaultGormDB().Table("black_lists").Where("uid=?", userId).Find(&blockUser).Error; err != nil {
		return blockUserInfo, err
	}
	user := db.User{
		UserID: blockUser.UserId,
	}
	if err := db.DB.MysqlDB.DefaultGormDB().Find(&user).Error; err != nil {
		return blockUserInfo, err
	}
	blockUserInfo.User.UserID = user.UserID
	blockUserInfo.User.FaceURL = user.FaceURL
	blockUserInfo.User.Nickname = user.Nickname
	blockUserInfo.User.Birth = user.Birth
	blockUserInfo.User.PhoneNumber = user.PhoneNumber
	blockUserInfo.User.Email = user.Email
	blockUserInfo.User.Gender = user.Gender
	blockUserInfo.BeginDisableTime = blockUser.BeginDisableTime
	blockUserInfo.EndDisableTime = blockUser.EndDisableTime
	return blockUserInfo, nil
}

func GetBlockUsers(showNumber, pageNumber int32) ([]BlockUserInfo, error) {
	var blockUserInfos []BlockUserInfo
	var blockUsers []db.BlackList
	if err := db.DB.MysqlDB.DefaultGormDB().Limit(int(showNumber)).Offset(int(showNumber * (pageNumber - 1))).Find(&blockUsers).Error; err != nil {
		return blockUserInfos, err
	}
	for _, blockUser := range blockUsers {
		var user db.User
		if err := db.DB.MysqlDB.DefaultGormDB().Table("users").Where("user_id=?", blockUser.UserId).First(&user).Error; err == nil {
			blockUserInfos = append(blockUserInfos, BlockUserInfo{
				User: db.User{
					UserID:      user.UserID,
					Nickname:    user.Nickname,
					FaceURL:     user.FaceURL,
					Birth:       user.Birth,
					PhoneNumber: user.PhoneNumber,
					Email:       user.Email,
					Gender:      user.Gender,
				},
				BeginDisableTime: blockUser.BeginDisableTime,
				EndDisableTime:   blockUser.EndDisableTime,
			})
		}
	}
	return blockUserInfos, nil
}

func GetUserByName(userName string, showNumber, pageNumber int32) ([]db.User, error) {
	var users []db.User
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").Where(" name like ?", fmt.Sprintf("%%%s%%", userName)).Limit(int(showNumber)).Offset(int(showNumber * (pageNumber - 1))).Find(&users).Error
	return users, err
}

func GetUsersByNameAndID(content string, showNumber, pageNumber int32) ([]db.User, int64, error) {
	var users []db.User
	var count int64
	db := db.DB.MysqlDB.DefaultGormDB().Table("users").Where(" name like ? or user_id = ? ", fmt.Sprintf("%%%s%%", content), content)
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	err := db.Limit(int(showNumber)).Offset(int(showNumber * (pageNumber - 1))).Find(&users).Error
	return users, count, err
}

func GetUsersCount(userName string) (int32, error) {
	var count int64
	if err := db.DB.MysqlDB.DefaultGormDB().Table("users").Where(" name like ? ", fmt.Sprintf("%%%s%%", userName)).Count(&count).Error; err != nil {
		return 0, err
	}
	return int32(count), nil
}

func GetBlockUsersNumCount() (int32, error) {
	var count int64
	if err := db.DB.MysqlDB.DefaultGormDB().Model(&db.BlackList{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return int32(count), nil
}
