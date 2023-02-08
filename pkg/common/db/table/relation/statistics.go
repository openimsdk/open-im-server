package relation

import "time"

// these two is virtual table just for cms
type ActiveGroup struct {
	Name       string
	ID         string `gorm:"column:recv_id"`
	MessageNum int    `gorm:"column:message_num"`
}

type ActiveUser struct {
	Name       string
	ID         string `gorm:"column:send_id"`
	MessageNum int    `gorm:"column:message_num"`
}

type StatisticsInterface interface {
	GetActiveUserNum(from, to time.Time) (num int64, err error)
	GetIncreaseUserNum(from, to time.Time) (num int64, err error)
	GetTotalUserNum() (num int64, err error)
	GetTotalUserNumByDate(to time.Time) (num int64, err error)
	GetSingleChatMessageNum(from, to time.Time) (num int64, err error)
	GetGroupMessageNum(from, to time.Time) (num int64, err error)
	GetIncreaseGroupNum(from, to time.Time) (num int64, err error)
	GetTotalGroupNum() (num int64, err error)
	GetGroupNum(to time.Time) (num int64, err error)
	GetActiveGroups(from, to time.Time, limit int) ([]*ActiveGroup, error)
	GetActiveUsers(from, to time.Time, limit int) (activeUsers []*ActiveUser, err error)
}
