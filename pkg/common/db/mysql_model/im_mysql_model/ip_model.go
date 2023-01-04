package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/utils"
	"time"

	"gorm.io/gorm"
)

func IsLimitRegisterIp(RegisterIp string) (bool, error) {
	//如果已经存在则限制
	var count int64
	if err := db.DB.MysqlDB.DefaultGormDB().Table("ip_limits").Where("ip=? and limit_register=? and limit_time>now()", RegisterIp, 1).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func IsLimitLoginIp(LoginIp string) (bool, error) {
	//如果已经存在则限制
	var count int64
	if err := db.DB.MysqlDB.DefaultGormDB().Table("ip_limits").Where("ip=? and limit_login=? and limit_time>now()", LoginIp, 1).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func IsLimitUserLoginIp(userID string, loginIp string) (limit bool, err error) {
	//如果已经存在则放行
	var count int64
	result := db.DB.MysqlDB.DefaultGormDB().Table("user_ip_limits").Where("user_id=?", userID).Count(&count)
	if err := result.Error; err != nil {
		return true, err
	}
	if count < 1 {
		return false, nil
	}
	result = db.DB.MysqlDB.DefaultGormDB().Table("user_ip_limits").Where("user_id=? and ip = ?", userID, loginIp).Count(&count)
	if err := result.Error; err != nil {
		return true, err
	}

	return count > 0, nil
}

func QueryIPLimits(ip string) (*IpLimit, error) {
	var ipLimit IpLimit
	err := db.DB.MysqlDB.DefaultGormDB().Model(&IpLimit{}).Where("ip=?", ip).First(&ipLimit).Error
	return &ipLimit, err
}

func QueryUserIPLimits(ip string) ([]UserIpLimit, error) {
	var ips []UserIpLimit
	err := db.DB.MysqlDB.DefaultGormDB().Model(&UserIpLimit{}).Where("ip=?", ip).Find(&ips).Error
	return ips, err
}

func InsertOneIntoIpLimits(ipLimits IpLimit) error {
	return db.DB.MysqlDB.DefaultGormDB().Model(&IpLimit{}).Create(ipLimits).Error
}

func DeleteOneFromIpLimits(ip string) error {
	ipLimits := &IpLimit{Ip: ip}
	return db.DB.MysqlDB.DefaultGormDB().Model(ipLimits).Where("ip=?", ip).Delete(ipLimits).Error
}

func GetIpLimitsLoginByUserID(userID string) ([]UserIpLimit, error) {
	var ips []UserIpLimit
	err := db.DB.MysqlDB.DefaultGormDB().Model(&UserIpLimit{}).Where("user_id=?", userID).Find(&ips).Error
	return ips, err
}

func InsertUserIpLimitsLogin(userIp *UserIpLimit) error {
	userIp.CreateTime = time.Now()
	return db.DB.MysqlDB.DefaultGormDB().Model(&UserIpLimit{}).Create(userIp).Error
}

func DeleteUserIpLimitsLogin(userID, ip string) error {
	userIp := UserIpLimit{UserID: userID, Ip: ip}
	return db.DB.MysqlDB.DefaultGormDB().Model(&UserIpLimit{}).Delete(&userIp).Error
}

func GetRegisterUserNum(ip string) ([]string, error) {
	var userIDList []string
	err := db.DB.MysqlDB.DefaultGormDB().Model(&Register{}).Where("register_ip=?", ip).Pluck("user_id", &userIDList).Error
	return userIDList, err
}

func InsertIpRecord(userID, createIp string) error {
	record := &UserIpRecord{UserID: userID, CreateIp: createIp, LastLoginTime: time.Now(), LoginTimes: 1}
	err := db.DB.MysqlDB.DefaultGormDB().Model(&UserIpRecord{}).Create(record).Error
	return err
}

func UpdateIpReocord(userID, ip string) (err error) {
	record := &UserIpRecord{UserID: userID, LastLoginIp: ip, LastLoginTime: time.Now()}
	result := db.DB.MysqlDB.DefaultGormDB().Model(&UserIpRecord{}).Where("user_id=?", userID).Updates(record).Update("login_times", gorm.Expr("login_times+?", 1))
	if result.Error != nil {
		return utils.Wrap(result.Error, "")
	}
	if result.RowsAffected == 0 {
		err = InsertIpRecord(userID, ip)
	}
	return utils.Wrap(err, "")
}
