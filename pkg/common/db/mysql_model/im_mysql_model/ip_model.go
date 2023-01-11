package im_mysql_model

import (
	"Open_IM/pkg/utils"
	"time"

	"gorm.io/gorm"
)

var IPDB *gorm.DB

type UserIpRecord struct {
	UserID        string    `gorm:"column:user_id;primary_key;size:64"`
	CreateIp      string    `gorm:"column:create_ip;size:15"`
	LastLoginTime time.Time `gorm:"column:last_login_time"`
	LastLoginIp   string    `gorm:"column:last_login_ip;size:15"`
	LoginTimes    int32     `gorm:"column:login_times"`
}

// ip limit login
type IpLimit struct {
	Ip            string    `gorm:"column:ip;primary_key;size:15"`
	LimitRegister int32     `gorm:"column:limit_register;size:1"`
	LimitLogin    int32     `gorm:"column:limit_login;size:1"`
	CreateTime    time.Time `gorm:"column:create_time"`
	LimitTime     time.Time `gorm:"column:limit_time"`
}

// ip login
type UserIpLimit struct {
	UserID     string    `gorm:"column:user_id;primary_key;size:64"`
	Ip         string    `gorm:"column:ip;primary_key;size:15"`
	CreateTime time.Time `gorm:"column:create_time"`
}

func IsLimitRegisterIp(RegisterIp string) (bool, error) {
	//如果已经存在则限制
	var count int64
	if err := IPDB.Table("ip_limits").Where("ip=? and limit_register=? and limit_time>now()", RegisterIp, 1).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func IsLimitLoginIp(LoginIp string) (bool, error) {
	//如果已经存在则限制
	var count int64
	if err := IPDB.Table("ip_limits").Where("ip=? and limit_login=? and limit_time>now()", LoginIp, 1).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func IsLimitUserLoginIp(userID string, loginIp string) (limit bool, err error) {
	//如果已经存在则放行
	var count int64
	result := IPDB.Table("user_ip_limits").Where("user_id=?", userID).Count(&count)
	if err := result.Error; err != nil {
		return true, err
	}
	if count < 1 {
		return false, nil
	}
	result = IPDB.Table("user_ip_limits").Where("user_id=? and ip = ?", userID, loginIp).Count(&count)
	if err := result.Error; err != nil {
		return true, err
	}

	return count > 0, nil
}

func QueryIPLimits(ip string) (*IpLimit, error) {
	var ipLimit IpLimit
	err := IPDB.Model(&IpLimit{}).Where("ip=?", ip).First(&ipLimit).Error
	return &ipLimit, err
}

func QueryUserIPLimits(ip string) ([]UserIpLimit, error) {
	var ips []UserIpLimit
	err := IPDB.Model(&UserIpLimit{}).Where("ip=?", ip).Find(&ips).Error
	return ips, err
}

func InsertOneIntoIpLimits(ipLimits IpLimit) error {
	return IPDB.Model(&IpLimit{}).Create(ipLimits).Error
}

func DeleteOneFromIpLimits(ip string) error {
	ipLimits := &IpLimit{Ip: ip}
	return IPDB.Model(ipLimits).Where("ip=?", ip).Delete(ipLimits).Error
}

func GetIpLimitsLoginByUserID(userID string) ([]UserIpLimit, error) {
	var ips []UserIpLimit
	err := IPDB.Model(&UserIpLimit{}).Where("user_id=?", userID).Find(&ips).Error
	return ips, err
}

func InsertUserIpLimitsLogin(userIp *UserIpLimit) error {
	userIp.CreateTime = time.Now()
	return IPDB.Model(&UserIpLimit{}).Create(userIp).Error
}

func DeleteUserIpLimitsLogin(userID, ip string) error {
	userIp := UserIpLimit{UserID: userID, Ip: ip}
	return IPDB.Model(&UserIpLimit{}).Delete(&userIp).Error
}

func GetRegisterUserNum(ip string) ([]string, error) {
	var userIDList []string
	err := IPDB.Model(&Register{}).Where("register_ip=?", ip).Pluck("user_id", &userIDList).Error
	return userIDList, err
}

func InsertIpRecord(userID, createIp string) error {
	record := &UserIpRecord{UserID: userID, CreateIp: createIp, LastLoginTime: time.Now(), LoginTimes: 1}
	err := IPDB.Model(&UserIpRecord{}).Create(record).Error
	return err
}

func UpdateIpReocord(userID, ip string) (err error) {
	record := &UserIpRecord{UserID: userID, LastLoginIp: ip, LastLoginTime: time.Now()}
	result := IPDB.Model(&UserIpRecord{}).Where("user_id=?", userID).Updates(record).Update("login_times", gorm.Expr("login_times+?", 1))
	if result.Error != nil {
		return utils.Wrap(result.Error, "")
	}
	if result.RowsAffected == 0 {
		err = InsertIpRecord(userID, ip)
	}
	return utils.Wrap(err, "")
}
