package mysql

import (
	"Open_IM/pkg/common/constant"
	"time"
)

func GetActiveUserNum(from, to time.Time) (int32, error) {
	var num int64
	err := ChatLogDB.Table("chat_logs").Select("count(distinct(send_id))").Where("send_time >= ? and send_time <= ?", from, to).Count(&num).Error
	return int32(num), err
}

func GetIncreaseUserNum(from, to time.Time) (int32, error) {
	var num int64
	err := UserDB.Where("create_time >= ? and create_time <= ?", from, to).Count(&num).Error
	return int32(num), err
}

func GetTotalUserNum() (int32, error) {
	var num int64
	err := UserDB.Count(&num).Error
	return int32(num), err
}

func GetTotalUserNumByDate(to time.Time) (int32, error) {
	var num int64
	err := UserDB.Where("create_time <= ?", to).Count(&num).Error
	return int32(num), err
}

func GetPrivateMessageNum(from, to time.Time) (int32, error) {
	var num int64
	err := ChatLogDB.Where("send_time >= ? and send_time <= ? and session_type = ?", from, to, 1).Count(&num).Error
	return int32(num), err
}

func GetGroupMessageNum(from, to time.Time) (int32, error) {
	var num int64
	err := ChatLogDB.Where("send_time >= ? and send_time <= ? and session_type = ?", from, to, 2).Count(&num).Error
	return int32(num), err
}

func GetIncreaseGroupNum(from, to time.Time) (int32, error) {
	var num int64
	err := GroupDB.Where("create_time >= ? and create_time <= ?", from, to).Count(&num).Error
	return int32(num), err
}

func GetTotalGroupNum() (int32, error) {
	var num int64
	err := GroupDB.Count(&num).Error
	return int32(num), err
}

func GetGroupNum(to time.Time) (int32, error) {
	var num int64
	err := GroupDB.Where("create_time <= ?", to).Count(&num).Error
	return int32(num), err
}

type activeGroup struct {
	Name       string
	Id         string `gorm:"column:recv_id"`
	MessageNum int    `gorm:"column:message_num"`
}

func GetActiveGroups(from, to time.Time, limit int) ([]*activeGroup, error) {
	var activeGroups []*activeGroup
	err := ChatLogDB.Select("recv_id, count(*) as message_num").Where("send_time >= ? and send_time <= ? and session_type in (?)", from, to, []int{constant.GroupChatType, constant.SuperGroupChatType}).Group("recv_id").Limit(limit).Order("message_num DESC").Find(&activeGroups).Error
	for _, activeGroup := range activeGroups {
		group := Group{
			GroupID: activeGroup.Id,
		}
		GroupDB.Where("group_id= ? ", group.GroupID).Find(&group)
		activeGroup.Name = group.GroupName
	}
	return activeGroups, err
}

type activeUser struct {
	Name       string
	ID         string `gorm:"column:send_id"`
	MessageNum int    `gorm:"column:message_num"`
}

func GetActiveUsers(from, to time.Time, limit int) ([]*activeUser, error) {
	var activeUsers []*activeUser
	err := ChatLogDB.Select("send_id, count(*) as message_num").Where("send_time >= ? and send_time <= ? and session_type = ?", from, to, constant.SingleChatType).Group("send_id").Limit(limit).Order("message_num DESC").Find(&activeUsers).Error
	for _, activeUser := range activeUsers {
		user := User{
			UserID: activeUser.ID,
		}
		err = UserDB.Select("user_id, name").Find(&user).Error
		if err != nil {
			continue
		}
		activeUser.Name = user.Nickname
		activeUser.ID = user.UserID
	}
	return activeUsers, err
}
