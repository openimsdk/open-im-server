package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"time"
)

func GetActiveUserNum(from, to time.Time) (int32, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return 0, err
	}
	dbConn.LogMode(true)
	var num int32
	err = dbConn.Table("chat_logs").Select("count(distinct(send_id))").Where("create_time >= ? and create_time <= ?", from, to).Count(&num).Error
	return num, err
}

func GetIncreaseUserNum(from, to time.Time) (int32, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return 0, err
	}
	dbConn.LogMode(true)
	var num int32
	err = dbConn.Table("users").Where("create_time >= ? and create_time <= ?", from, to).Count(&num).Error
	return num, err
}

func GetTotalUserNum() (int32, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return 0, err
	}
	dbConn.LogMode(true)
	var num int32
	err = dbConn.Table("users").Count(&num).Error
	return num, err
}

func GetTotalUserNumByDate(to time.Time) (int32, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return 0, err
	}
	dbConn.LogMode(true)
	var num int32
	err = dbConn.Table("users").Where("create_time <= ?", to).Count(&num).Error
	return num, err
}

func GetPrivateMessageNum(from, to time.Time) (int32, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return 0, err
	}
	dbConn.LogMode(true)
	var num int32
	err = dbConn.Table("chat_logs").Where("create_time >= ? and create_time <= ? and session_type = ?", from, to, 1).Count(&num).Error
	return num, err
}

func GetGroupMessageNum(from, to time.Time) (int32, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return 0, err
	}
	dbConn.LogMode(true)
	var num int32
	err = dbConn.Table("chat_logs").Where("create_time >= ? and create_time <= ? and session_type = ?", from, to, 2).Count(&num).Error
	return num, err
}

func GetIncreaseGroupNum(from, to time.Time) (int32, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return 0, err
	}
	dbConn.LogMode(true)
	var num int32
	err = dbConn.Table("groups").Where("create_time >= ? and create_time <= ?", from, to).Count(&num).Error
	return num, err
}

func GetTotalGroupNum() (int32, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return 0, err
	}
	dbConn.LogMode(true)
	var num int32
	err = dbConn.Table("groups").Count(&num).Error
	return num, err
}

func GetGroupNum(to time.Time) (int32, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return 0, err
	}
	dbConn.LogMode(true)
	var num int32
	err = dbConn.Table("groups").Where("create_time <= ?", to).Count(&num).Error
	return num, err
}

type activeGroup struct {
	Name       string
	Id         string `gorm:"column:recv_id"`
	MessageNum int    `gorm:"column:message_num"`
}

func GetActiveGroups(from, to time.Time, limit int) ([]*activeGroup, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	var activeGroups []*activeGroup
	if err != nil {
		return activeGroups, err
	}
	dbConn.LogMode(true)
	err = dbConn.Table("chat_logs").Select("recv_id, count(*) as message_num").Where("create_time >= ? and create_time <= ? and session_type = ?", from, to, 2).Group("recv_id").Limit(limit).Order("message_num DESC").Find(&activeGroups).Error
	for _, activeGroup := range activeGroups {
		group := db.Group{
			GroupID: activeGroup.Id,
		}
		dbConn.Table("groups").Where("group_id= ? ", group.GroupID).Find(&group)
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
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	var activeUsers []*activeUser
	if err != nil {
		return activeUsers, err
	}
	dbConn.LogMode(true)
	err = dbConn.Table("chat_logs").Select("send_id, count(*) as message_num").Where("create_time >= ? and create_time <= ? and session_type = ?", from, to, 1).Group("send_id").Limit(limit).Order("message_num DESC").Find(&activeUsers).Error
	for _, activeUser := range activeUsers {
		user := db.User{
			UserID: activeUser.Id,
		}
		dbConn.Table("users").Select("user_id, name").Find(&user)
		activeUser.Name = user.Nickname
	}
	return activeUsers, err
}
