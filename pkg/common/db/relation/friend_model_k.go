package relation

import (
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/utils"
	"context"
	"gorm.io/gorm"
	"time"
)

type Friend struct {
	OwnerUserID    string    `gorm:"column:owner_user_id;primary_key;size:64"`
	FriendUserID   string    `gorm:"column:friend_user_id;primary_key;size:64"`
	Remark         string    `gorm:"column:remark;size:255"`
	CreateTime     time.Time `gorm:"column:create_time"`
	AddSource      int32     `gorm:"column:add_source"`
	OperatorUserID string    `gorm:"column:operator_user_id;size:64"`
	Ex             string    `gorm:"column:ex;size:1024"`
	DB             *gorm.DB  `gorm:"-"`
}

type FriendUser struct {
	Friend
	Nickname string `gorm:"column:name;size:255"`
}

func NewFriendDB(db *gorm.DB) *Friend {
	var friend Friend
	friend.DB = initModel(db, friend)
	return &friend
}

func (f *Friend) Create(ctx context.Context, friends []*Friend) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "friends", friends)
	}()
	return utils.Wrap(f.DB.Create(&friends).Error, "")
}

func (f *Friend) Delete(ctx context.Context, ownerUserID string, friendUserIDs string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "ownerUserID", ownerUserID, "friendUserIDs", friendUserIDs)
	}()
	err = utils.Wrap(f.DB.Where("owner_user_id = ? and friend_user_id = ?", ownerUserID, friendUserIDs).Delete(&Friend{}).Error, "")
	return err
}

func (f *Friend) UpdateByMap(ctx context.Context, ownerUserID string, args map[string]interface{}) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "ownerUserID", ownerUserID, "args", args)
	}()
	return utils.Wrap(f.DB.Where("owner_user_id = ?", ownerUserID).Updates(args).Error, "")
}

func (f *Friend) Update(ctx context.Context, friends []*Friend) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "friends", friends)
	}()
	return utils.Wrap(f.DB.Updates(&friends).Error, "")
}

func (f *Friend) UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "ownerUserID", ownerUserID, "friendUserID", friendUserID, "remark", remark)
	}()
	return utils.Wrap(f.DB.Model(f).Where("owner_user_id = ? and friend_user_id = ?", ownerUserID, friendUserID).Update("remark", remark).Error, "")
}

func (f *Friend) FindOwnerUserID(ctx context.Context, ownerUserID string) (friends []*Friend, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "ownerUserID", ownerUserID, "friends", friends)
	}()
	return friends, utils.Wrap(f.DB.Where("owner_user_id = ?", ownerUserID).Find(&friends).Error, "")
}

func (f *Friend) FindFriendUserID(ctx context.Context, friendUserID string) (friends []*Friend, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "friendUserID", friendUserID, "friends", friends)
	}()
	return friends, utils.Wrap(f.DB.Where("friend_user_id = ?", friendUserID).Find(&friends).Error, "")
}

func (f *Friend) Take(ctx context.Context, ownerUserID, friendUserID string) (friend *Friend, err error) {
	friend = &Friend{}
	defer tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "ownerUserID", ownerUserID, "friendUserID", friendUserID, "friend", friend)
	return friend, utils.Wrap(f.DB.Where("owner_user_id = ? and friend_user_id", ownerUserID, friendUserID).Take(friend).Error, "")
}

func (f *Friend) FindUserState(ctx context.Context, userID1, userID2 string) (friends []*Friend, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "userID1", userID1, "userID2", userID2)
	}()
	return friends, utils.Wrap(f.DB.Where("(owner_user_id = ? and friend_user_id = ?) or (owner_user_id = ? and friend_user_id = ?)", userID1, userID2, userID2, userID1).Find(&friends).Error, "")
}
