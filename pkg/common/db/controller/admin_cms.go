package controller

import (
	"Open_IM/pkg/common/db/relation"
	"gorm.io/gorm"
	"time"
)

type AdminCMSInterface interface {
	GetActiveUserNum(from, to time.Time) (num int64, err error)
	GetIncreaseUserNum(from, to time.Time) (num int64, err error)
	GetTotalUserNum() (num int64, err error)
	GetTotalUserNumByDate(to time.Time) (num int64, err error)
	GetSingleChatMessageNum(from, to time.Time) (num int64, err error)
	GetGroupMessageNum(from, to time.Time) (num int64, err error)
	GetIncreaseGroupNum(from, to time.Time) (num int64, err error)
	GetTotalGroupNum() (num int64, err error)
	GetGroupNum(to time.Time) (num int64, err error)
	GetActiveGroups(from, to time.Time, limit int) (activeGroups []*relation.ActiveGroup, err error)
	GetActiveUsers(from, to time.Time, limit int) (activeUsers []*relation.ActiveUser, err error)
}

type AdminCMSController struct {
	database AdminCMSDatabaseInterface
}

func NewAdminCMSController(db *gorm.DB) AdminCMSInterface {
	adminCMSController := &AdminCMSController{
		database: newAdminCMSDatabase(db),
	}
	return adminCMSController
}

func newAdminCMSDatabase(db *gorm.DB) AdminCMSDatabaseInterface {
	return &AdminCMSDatabase{Statistics: relation.NewStatistics(db)}
}

func (admin *AdminCMSController) GetActiveUserNum(from, to time.Time) (num int64, err error) {
	return admin.database.GetActiveUserNum(from, to)
}

func (admin *AdminCMSController) GetIncreaseUserNum(from, to time.Time) (num int64, err error) {
	return admin.database.GetIncreaseUserNum(from, to)
}

func (admin *AdminCMSController) GetTotalUserNum() (num int64, err error) {
	return admin.database.GetTotalUserNum()
}

func (admin *AdminCMSController) GetTotalUserNumByDate(to time.Time) (num int64, err error) {
	return admin.database.GetTotalUserNumByDate(to)
}

func (admin *AdminCMSController) GetSingleChatMessageNum(from, to time.Time) (num int64, err error) {
	return admin.GetSingleChatMessageNum(from, to)
}

func (admin *AdminCMSController) GetGroupMessageNum(from, to time.Time) (num int64, err error) {
	return admin.database.GetGroupMessageNum(from, to)
}

func (admin *AdminCMSController) GetIncreaseGroupNum(from, to time.Time) (num int64, err error) {
	return admin.database.GetIncreaseGroupNum(from, to)
}

func (admin *AdminCMSController) GetTotalGroupNum() (num int64, err error) {
	return admin.database.GetTotalGroupNum()
}

func (admin *AdminCMSController) GetGroupNum(to time.Time) (num int64, err error) {
	return admin.database.GetGroupNum(to)
}

func (admin *AdminCMSController) GetActiveGroups(from, to time.Time, limit int) ([]*relation.ActiveGroup, error) {
	return admin.database.GetActiveGroups(from, to, limit)
}

func (admin *AdminCMSController) GetActiveUsers(from, to time.Time, limit int) (activeUsers []*relation.ActiveUser, err error) {
	return admin.database.GetActiveUsers(from, to, limit)
}

type AdminCMSDatabaseInterface interface {
	GetActiveUserNum(from, to time.Time) (num int64, err error)
	GetIncreaseUserNum(from, to time.Time) (num int64, err error)
	GetTotalUserNum() (num int64, err error)
	GetTotalUserNumByDate(to time.Time) (num int64, err error)
	GetSingleChatMessageNum(from, to time.Time) (num int64, err error)
	GetGroupMessageNum(from, to time.Time) (num int64, err error)
	GetIncreaseGroupNum(from, to time.Time) (num int64, err error)
	GetTotalGroupNum() (num int64, err error)
	GetGroupNum(to time.Time) (num int64, err error)
	GetActiveGroups(from, to time.Time, limit int) ([]*relation.ActiveGroup, error)
	GetActiveUsers(from, to time.Time, limit int) (activeUsers []*relation.ActiveUser, err error)
}

type AdminCMSDatabase struct {
	Statistics *relation.Statistics
}

func (admin *AdminCMSDatabase) GetActiveUserNum(from, to time.Time) (num int64, err error) {
	return admin.Statistics.GetActiveUserNum(from, to)
}

func (admin *AdminCMSDatabase) GetIncreaseUserNum(from, to time.Time) (num int64, err error) {
	return admin.Statistics.GetIncreaseUserNum(from, to)
}

func (admin *AdminCMSDatabase) GetTotalUserNum() (num int64, err error) {
	return admin.Statistics.GetTotalUserNum()
}

func (admin *AdminCMSDatabase) GetTotalUserNumByDate(to time.Time) (num int64, err error) {
	return admin.Statistics.GetTotalUserNumByDate(to)
}

func (admin *AdminCMSDatabase) GetSingleChatMessageNum(from, to time.Time) (num int64, err error) {
	return admin.Statistics.GetSingleChatMessageNum(from, to)
}

func (admin *AdminCMSDatabase) GetGroupMessageNum(from, to time.Time) (num int64, err error) {
	return admin.Statistics.GetGroupMessageNum(from, to)
}

func (admin *AdminCMSDatabase) GetIncreaseGroupNum(from, to time.Time) (num int64, err error) {
	return admin.Statistics.GetIncreaseGroupNum(from, to)
}

func (admin *AdminCMSDatabase) GetTotalGroupNum() (num int64, err error) {
	return admin.Statistics.GetTotalGroupNum()
}

func (admin *AdminCMSDatabase) GetGroupNum(to time.Time) (num int64, err error) {
	return admin.Statistics.GetGroupNum(to)
}

func (admin *AdminCMSDatabase) GetActiveGroups(from, to time.Time, limit int) ([]*relation.ActiveGroup, error) {
	return admin.Statistics.GetActiveGroups(from, to, limit)
}

func (admin *AdminCMSDatabase) GetActiveUsers(from, to time.Time, limit int) (activeUsers []*relation.ActiveUser, err error) {
	return admin.Statistics.GetActiveUsers(from, to, limit)
}
