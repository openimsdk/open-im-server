package relation

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/utils"
	"fmt"
	"gorm.io/gorm"
	"time"
)

var (
	BlackListDB *gorm.DB
	UserDB      *gorm.DB
)

func InitManager() {
	for k, v := range config.Config.Manager.AppManagerUid {
		_, err := GetUserByUserID(v)
		if err != nil {
		} else {
			continue
		}
		var appMgr User
		appMgr.UserID = v
		if k == 0 {
			appMgr.Nickname = config.Config.Manager.AppSysNotificationName
		} else {
			appMgr.Nickname = "AppManager" + utils.IntToString(k+1)
		}
		appMgr.AppMangerLevel = constant.AppAdmin
		err = UserRegister(appMgr)
		if err != nil {
			fmt.Println("AppManager insert error ", err.Error(), appMgr)
		} else {
			fmt.Println("AppManager insert ", appMgr)
		}
	}
}

func UserRegister(user User) error {
	user.CreateTime = time.Now()
	if user.AppMangerLevel == 0 {
		user.AppMangerLevel = constant.AppOrdinaryUsers
	}
	if user.Birth.Unix() < 0 {
		user.Birth = utils.UnixSecondToTime(0)
	}
	err := UserDB.Table("users").Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func GetAllUser() ([]User, error) {
	var userList []User
	err := UserDB.Table("users").Find(&userList).Error
	return userList, err
}

func TakeUserByUserID(userID string) (*User, error) {
	var user User
	err := UserDB.Table("users").Where("user_id=?", userID).Take(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByUserID(userID string) (*User, error) {
	var user User
	err := UserDB.Table("users").Where("user_id=?", userID).Take(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUsersByUserIDList(userIDList []string) ([]*User, error) {
	var userList []*User
	err := UserDB.Table("users").Where("user_id in (?)", userIDList).Find(&userList).Error
	return userList, err
}

func GetUserNameByUserID(userID string) (string, error) {
	var user User
	err := UserDB.Table("users").Select("name").Where("user_id=?", userID).First(&user).Error
	if err != nil {
		return "", err
	}
	return user.Nickname, nil
}

func UpdateUserInfo(user User) error {
	return UserDB.Where("user_id=?", user.UserID).Updates(&user).Error
}

func UpdateUserInfoByMap(user User, m map[string]interface{}) error {
	err := UserDB.Where("user_id=?", user.UserID).Updates(m).Error
	return err
}

func SelectAllUserID() ([]string, error) {
	var resultArr []string
	err := UserDB.Pluck("user_id", &resultArr).Error
	if err != nil {
		return nil, err
	}
	return resultArr, nil
}

func SelectSomeUserID(userIDList []string) ([]string, error) {
	var resultArr []string
	err := UserDB.Pluck("user_id", &resultArr).Error
	if err != nil {
		return nil, err
	}
	return resultArr, nil
}

func GetUsers(showNumber, pageNumber int32) ([]User, error) {
	var users []User
	err := UserDB.Limit(int(showNumber)).Offset(int(showNumber * (pageNumber - 1))).Find(&users).Error
	if err != nil {
		return users, err
	}
	return users, err
}

func AddUser(userID string, phoneNumber string, name string, email string, gender int32, faceURL string, birth string) error {
	_birth, err := utils.TimeStringToTime(birth)
	if err != nil {
		return err
	}
	user := User{
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
	result := UserDB.Create(&user)
	return result.Error
}

func UsersIsBlock(userIDList []string) (inBlockUserIDList []string, err error) {
	err = BlackListDB.Where("uid in (?) and end_disable_time > now()", userIDList).Pluck("uid", &inBlockUserIDList).Error
	return inBlockUserIDList, err
}

type BlockUserInfo struct {
	User             User
	BeginDisableTime time.Time
	EndDisableTime   time.Time
}

func GetUserByName(userName string, showNumber, pageNumber int32) ([]User, error) {
	var users []User
	err := UserDB.Where(" name like ?", fmt.Sprintf("%%%s%%", userName)).Limit(int(showNumber)).Offset(int(showNumber * (pageNumber - 1))).Find(&users).Error
	return users, err
}

func GetUsersByNameAndID(content string, showNumber, pageNumber int32) ([]User, int64, error) {
	var users []User
	var count int64
	db := UserDB.Where(" name like ? or user_id = ? ", fmt.Sprintf("%%%s%%", content), content)
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	err := db.Limit(int(showNumber)).Offset(int(showNumber * (pageNumber - 1))).Find(&users).Error
	return users, count, err
}

func GetUserIDsByEmailAndID(phoneNumber, email string) ([]string, error) {
	if phoneNumber == "" && email == "" {
		return nil, nil
	}
	db := UserDB
	if phoneNumber != "" {
		db = db.Where("phone_number = ? ", phoneNumber)
	}
	if email != "" {
		db = db.Where("email = ? ", email)
	}
	var userIDList []string
	err := db.Pluck("user_id", &userIDList).Error
	return userIDList, err
}

func GetUsersCount(userName string) (int32, error) {
	var count int64
	if err := UserDB.Where(" name like ? ", fmt.Sprintf("%%%s%%", userName)).Count(&count).Error; err != nil {
		return 0, err
	}
	return int32(count), nil
}

func GetBlockUsersNumCount() (int32, error) {
	var count int64
	if err := BlackListDB.Count(&count).Error; err != nil {
		return 0, err
	}
	return int32(count), nil
}
