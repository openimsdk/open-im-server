package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"time"
)

func GetActiveUserNum(from, to time.Time) (int32, error) {
	var num int64
	err := db.DB.MysqlDB.DefaultGormDB().Table("chat_logs").Select("count(distinct(send_id))").Where("create_time >= ? and create_time <= ?", from, to).Count(&num).Error
	return int32(num), err
}

func GetIncreaseUserNum(from, to time.Time) (int32, error) {
	var num int64
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").Where("create_time >= ? and create_time <= ?", from, to).Count(&num).Error
	return int32(num), err
}

func GetTotalUserNum() (int32, error) {
	var num int64
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").Count(&num).Error
	return int32(num), err
}

func GetTotalUserNumByDate(to time.Time) (int32, error) {
	var num int64
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").Where("create_time <= ?", to).Count(&num).Error
	return int32(num), err
}

func GetPrivateMessageNum(from, to time.Time) (int32, error) {
	var num int64
	err := db.DB.MysqlDB.DefaultGormDB().Table("chat_logs").Where("create_time >= ? and create_time <= ? and session_type = ?", from, to, 1).Count(&num).Error
	return int32(num), err
}

func GetGroupMessageNum(from, to time.Time) (int32, error) {
	var num int64
	err := db.DB.MysqlDB.DefaultGormDB().Table("chat_logs").Where("create_time >= ? and create_time <= ? and session_type = ?", from, to, 2).Count(&num).Error
	return int32(num), err
}

func GetIncreaseGroupNum(from, to time.Time) (int32, error) {
	var num int64
	err := db.DB.MysqlDB.DefaultGormDB().Table("groups").Where("create_time >= ? and create_time <= ?", from, to).Count(&num).Error
	return int32(num), err
}

func GetTotalGroupNum() (int32, error) {
	var num int64
	err := db.DB.MysqlDB.DefaultGormDB().Table("groups").Count(&num).Error
	return int32(num), err
}

func GetGroupNum(to time.Time) (int32, error) {
	var num int64
	err := db.DB.MysqlDB.DefaultGormDB().Table("groups").Where("create_time <= ?", to).Count(&num).Error
	return int32(num), err
}

type activeGroup struct {
	Name       string
	Id         string `gorm:"column:recv_id"`
	MessageNum int    `gorm:"column:message_num"`
}

func GetActiveGroups(from, to time.Time, limit int) ([]*activeGroup, error) {
	var activeGroups []*activeGroup
	err := db.DB.MysqlDB.DefaultGormDB().Table("chat_logs").Select("recv_id, count(*) as message_num").Where("create_time >= ? and create_time <= ? and session_type = ?", from, to, 2).Group("recv_id").Limit(limit).Order("message_num DESC").Find(&activeGroups).Error
	for _, activeGroup := range activeGroups {
		group := db.Group{
			GroupID: activeGroup.Id,
		}
		db.DB.MysqlDB.DefaultGormDB().Table("groups").Where("group_id= ? ", group.GroupID).Find(&group)
		activeGroup.Name = group.GroupName
	}
	return activeGroups, err
}

type activeUser struct {
	Name       string
	Id         string `gorm:"column:send_id"`
	MessageNum int    `gorm:"column:message_num"`
}

func GetActiveUsers(from, to time.Time, limit int) ([]*activeUser, error) {
	var activeUsers []*activeUser
	err := db.DB.MysqlDB.DefaultGormDB().Table("chat_logs").Select("send_id, count(*) as message_num").Where("create_time >= ? and create_time <= ? and session_type = ?", from, to, 1).Group("send_id").Limit(limit).Order("message_num DESC").Find(&activeUsers).Error
	for _, activeUser := range activeUsers {
		user := db.User{
			UserID: activeUser.Id,
		}
		db.DB.MysqlDB.DefaultGormDB().Table("users").Select("user_id, name").Find(&user)
		activeUser.Name = user.Nickname
	}
	return activeUsers, err
}
