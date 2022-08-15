package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"github.com/jinzhu/gorm"
	"time"
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

func IsLimitUserLoginIp(userID string, LoginIp string) (bool, error) {
	//如果已经存在则放行
	var count int64
	if err := db.DB.MysqlDB.DefaultGormDB().Table("user_ip_limits").Where("ip=? and user_id=?", LoginIp, userID).Count(&count).Error; err != nil {
		return false, err
	}
	return count == 0, nil
}

func QueryIPLimits(ip string) (*db.IpLimit, error) {
	var ipLimit db.IpLimit
	err := db.DB.MysqlDB.DefaultGormDB().Model(&db.IpLimit{}).Where("ip=?", ip).Take(&ip).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &ipLimit, err
}

func QueryUserIPLimits(ip string) ([]db.UserIpLimit, error) {
	var ips []db.UserIpLimit
	err := db.DB.MysqlDB.DefaultGormDB().Model(&db.UserIpLimit{}).Where("ip=?", ip).Find(&ips).Error
	return ips, err
}

func InsertOneIntoIpLimits(ipLimits db.IpLimit) error {
	return db.DB.MysqlDB.DefaultGormDB().Model(&db.IpLimit{}).Create(ipLimits).Error
}

func DeleteOneFromIpLimits(ip string) error {
	ipLimits := &db.IpLimit{Ip: ip}
	return db.DB.MysqlDB.DefaultGormDB().Model(ipLimits).Where("ip=?", ip).Delete(ipLimits).Error
}

func GetIpLimitsLoginByUserID(userID string) ([]db.UserIpLimit, error) {
	var ips []db.UserIpLimit
	err := db.DB.MysqlDB.DefaultGormDB().Model(&db.UserIpLimit{}).Where("user_id=?", userID).Find(&ips).Error
	return ips, err
}

func InsertUserIpLimitsLogin(userIp *db.UserIpLimit) error {
	userIp.CreateTime = time.Now()
	return db.DB.MysqlDB.DefaultGormDB().Model(&db.UserIpLimit{}).Create(userIp).Error
}

func DeleteUserIpLimitsLogin(userID, ip string) error {
	userIp := db.UserIpLimit{UserID: userID, Ip: ip}
	return db.DB.MysqlDB.DefaultGormDB().Model(&db.UserIpLimit{}).Delete(&userIp).Error
}

func GetRegisterUserNum(ip string) ([]string, error) {
	var userIDList []string
	err := db.DB.MysqlDB.DefaultGormDB().Model(&db.Register{}).Where("register_ip=?", ip).Pluck("user_id", &userIDList).Error
	return userIDList, err
}
