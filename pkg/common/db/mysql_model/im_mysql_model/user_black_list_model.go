package im_mysql_model

import (
	"Open_IM/pkg/common/trace_log"
	"Open_IM/pkg/utils"
	"context"
	"gorm.io/gorm"
	"time"
)

var BlackDB *gorm.DB

type Black struct {
	OwnerUserID    string    `gorm:"column:owner_user_id;primary_key;size:64"`
	BlockUserID    string    `gorm:"column:block_user_id;primary_key;size:64"`
	CreateTime     time.Time `gorm:"column:create_time"`
	AddSource      int32     `gorm:"column:add_source"`
	OperatorUserID string    `gorm:"column:operator_user_id;size:64"`
	Ex             string    `gorm:"column:ex;size:1024"`
}

func InsertInToUserBlackList(ctx context.Context, black Black) (err error) {
	defer trace_log.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "black", black)
	black.CreateTime = time.Now()
	err = BlackDB.Create(black).Error
	return err
}

func CheckBlack(ownerUserID, blockUserID string) error {
	var black Black
	return BlackDB.Where("owner_user_id=? and block_user_id=?", ownerUserID, blockUserID).Find(&black).Error
}

func RemoveBlackList(ownerUserID, blockUserID string) error {
	err := BlackDB.Where("owner_user_id=? and block_user_id=?", ownerUserID, blockUserID).Delete(Black{}).Error
	return utils.Wrap(err, "RemoveBlackList failed")
}

func GetBlackListByUserID(ownerUserID string) ([]Black, error) {
	var blackListUsersInfo []Black
	err := BlackDB.Where("owner_user_id=?", ownerUserID).Find(&blackListUsersInfo).Error
	if err != nil {
		return nil, err
	}
	return blackListUsersInfo, nil
}

func GetBlackIDListByUserID(ownerUserID string) ([]string, error) {
	var blackIDList []string
	err := BlackDB.Where("owner_user_id=?", ownerUserID).Pluck("block_user_id", &blackIDList).Error
	if err != nil {
		return nil, err
	}
	return blackIDList, nil
}
