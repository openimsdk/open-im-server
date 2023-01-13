package im_mysql_model

import (
	"Open_IM/pkg/common/trace_log"
	"Open_IM/pkg/utils"
	"context"
	"gorm.io/gorm"
	"time"
)

var FriendDB *gorm.DB

type Friend struct {
	OwnerUserID    string    `gorm:"column:owner_user_id;primary_key;size:64"`
	FriendUserID   string    `gorm:"column:friend_user_id;primary_key;size:64"`
	Remark         string    `gorm:"column:remark;size:255"`
	CreateTime     time.Time `gorm:"column:create_time"`
	AddSource      int32     `gorm:"column:add_source"`
	OperatorUserID string    `gorm:"column:operator_user_id;size:64"`
	Ex             string    `gorm:"column:ex;size:1024"`
}

func (*Friend) Create(ctx context.Context, friends []*Friend) (err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "friends", friends)
	}()
	err = utils.Wrap(FriendDB.Create(&friends).Error, "")
	return err
}

func (*Friend) Delete(ctx context.Context, ownerUserID string, friendUserIDs []string) (err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "ownerUserID", ownerUserID, "friendUserIDs", friendUserIDs)
	}()
	err = utils.Wrap(FriendDB.Where("owner_user_id = ? and friend_user_id in (?)", ownerUserID, friendUserIDs).Delete(&Friend{}).Error, "")
	return err
}

func (*Friend) UpdateByMap(ctx context.Context, ownerUserID string, args map[string]interface{}) (err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "ownerUserID", ownerUserID, "args", args)
	}()
	return utils.Wrap(FriendDB.Where("owner_user_id = ?", ownerUserID).Updates(args).Error, "")
}

func (*Friend) Update(ctx context.Context, friends []*Friend) (err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "friends", friends)
	}()
	return utils.Wrap(FriendDB.Updates(&friends).Error, "")
}

func (*Friend) Find(ctx context.Context, ownerUserID string) (friends []*Friend, err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "ownerUserID", ownerUserID, "friends", friends)
	}()
	err = utils.Wrap(FriendDB.Where("owner_user_id = ?", ownerUserID).Find(&friends).Error, "")
	return friends, err
}

func (*Friend) Take(ctx context.Context, ownerUserID, friendUserID string) (friend *Friend, err error) {
	friend = &Friend{}
	defer trace_log.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "ownerUserID", ownerUserID, "friendUserID", friendUserID, "group", *friend)
	err = utils.Wrap(FriendDB.Where("owner_user_id = ? and friend_user_id", ownerUserID, friendUserID).Take(friend).Error, "")
	return friend, err
}
