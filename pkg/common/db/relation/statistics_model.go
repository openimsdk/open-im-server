package relation

import (
	"Open_IM/pkg/common/constant"
	"gorm.io/gorm"
	"time"
)

type Statistics struct {
	DB *gorm.DB
}

func NewStatistics(db *gorm.DB) *Statistics {
	return &Statistics{DB: db}
}

func (s *Statistics) getUserModel() *gorm.DB {
	return s.DB.Model(&User{})
}

func (s *Statistics) getChatLogModel() *gorm.DB {
	return s.DB.Model(&ChatLog{})
}

func (s *Statistics) getGroupModel() *gorm.DB {
	return s.DB.Model(&Group{})
}

func (s *Statistics) GetActiveUserNum(from, to time.Time) (num int64, err error) {
	err = s.getChatLogModel().Select("count(distinct(send_id))").Where("send_time >= ? and send_time <= ?", from, to).Count(&num).Error
	return num, err
}

func (s *Statistics) GetIncreaseUserNum(from, to time.Time) (num int64, err error) {
	err = s.getUserModel().Where("create_time >= ? and create_time <= ?", from, to).Count(&num).Error
	return num, err
}

func (s *Statistics) GetTotalUserNum() (num int64, err error) {
	err = s.getUserModel().Count(&num).Error
	return num, err
}

func (s *Statistics) GetTotalUserNumByDate(to time.Time) (num int64, err error) {
	err = s.getUserModel().Where("create_time <= ?", to).Count(&num).Error
	return num, err
}

func (s *Statistics) GetSingleChatMessageNum(from, to time.Time) (num int64, err error) {
	err = s.getChatLogModel().Where("send_time >= ? and send_time <= ? and session_type = ?", from, to, constant.SingleChatType).Count(&num).Error
	return num, err
}

func (s *Statistics) GetGroupMessageNum(from, to time.Time) (num int64, err error) {
	err = s.getChatLogModel().Where("send_time >= ? and send_time <= ? and session_type in (?)", from, to, []int{constant.GroupChatType, constant.SuperGroupChatType}).Count(&num).Error
	return num, err
}

func (s *Statistics) GetIncreaseGroupNum(from, to time.Time) (num int64, err error) {
	err = s.getGroupModel().Where("create_time >= ? and create_time <= ?", from, to).Count(&num).Error
	return num, err
}

func (s *Statistics) GetTotalGroupNum() (num int64, err error) {
	err = s.getGroupModel().Count(&num).Error
	return num, err
}

func (s *Statistics) GetGroupNum(to time.Time) (num int64, err error) {
	err = s.getGroupModel().Where("create_time <= ?", to).Count(&num).Error
	return num, err
}

type ActiveGroup struct {
	Name       string
	ID         string `gorm:"column:recv_id"`
	MessageNum int    `gorm:"column:message_num"`
}

func (s *Statistics) GetActiveGroups(from, to time.Time, limit int) ([]*ActiveGroup, error) {
	var activeGroups []*ActiveGroup
	err := s.getChatLogModel().Select("recv_id, count(*) as message_num").Where("send_time >= ? and send_time <= ? and session_type in (?)", from, to, []int{constant.GroupChatType, constant.SuperGroupChatType}).Group("recv_id").Limit(limit).Order("message_num DESC").Find(&activeGroups).Error
	for _, activeGroup := range activeGroups {
		group := Group{
			GroupID: activeGroup.ID,
		}
		s.getGroupModel().Where("group_id= ? ", group.GroupID).Find(&group)
		activeGroup.Name = group.GroupName
	}
	return activeGroups, err
}

type ActiveUser struct {
	Name       string
	ID         string `gorm:"column:send_id"`
	MessageNum int    `gorm:"column:message_num"`
}

func (s *Statistics) GetActiveUsers(from, to time.Time, limit int) (activeUsers []*ActiveUser, err error) {
	err = s.getChatLogModel().Select("send_id, count(*) as message_num").Where("send_time >= ? and send_time <= ? and session_type in (?)", from, to, []int{constant.SingleChatType, constant.GroupChatType, constant.SuperGroupChatType}).Group("send_id").Limit(limit).Order("message_num DESC").Find(&activeUsers).Error
	for _, activeUser := range activeUsers {
		user := User{
			UserID: activeUser.ID,
		}
		err = s.getUserModel().Select("user_id, name").Find(&user).Error
		if err != nil {
			return nil, err
		}
		activeUser.Name = user.Nickname
		activeUser.ID = user.UserID
	}
	return activeUsers, err
}
