package relation

import (
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/utils"
	"context"
	"gorm.io/gorm"
	"time"
)

type FriendDB interface {
	Create(ctx context.Context, friends []*FriendModel) (err error)
	Delete(ctx context.Context, ownerUserID string, friendUserIDs string) (err error)
	UpdateByMap(ctx context.Context, ownerUserID string, args map[string]interface{}) (err error)
	Update(ctx context.Context, friends []*FriendModel) (err error)
	UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) (err error)
	FindOwnerUserID(ctx context.Context, ownerUserID string) (friends []*FriendModel, err error)
}

type FriendModel struct {
	OwnerUserID    string    `gorm:"column:owner_user_id;primary_key;size:64"`
	FriendUserID   string    `gorm:"column:friend_user_id;primary_key;size:64"`
	Remark         string    `gorm:"column:remark;size:255"`
	CreateTime     time.Time `gorm:"column:create_time"`
	AddSource      int32     `gorm:"column:add_source"`
	OperatorUserID string    `gorm:"column:operator_user_id;size:64"`
	Ex             string    `gorm:"column:ex;size:1024"`
}
type FriendGorm struct {
	DB *gorm.DB `gorm:"-"`
}

func NewFriendGorm(DB *gorm.DB) *FriendGorm {
	return &FriendGorm{DB: DB}
}

type FriendUser struct {
	FriendGorm
	Nickname string `gorm:"column:name;size:255"`
}

func (f *FriendGorm) Create(ctx context.Context, friends []*FriendModel) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "friends", friends)
	}()
	return utils.Wrap(f.DB.Create(&friends).Error, "")
}

func (f *FriendGorm) Delete(ctx context.Context, ownerUserID string, friendUserIDs string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "ownerUserID", ownerUserID, "friendUserIDs", friendUserIDs)
	}()
	err = utils.Wrap(f.DB.Where("owner_user_id = ? and friend_user_id = ?", ownerUserID, friendUserIDs).Delete(&FriendModel{}).Error, "")
	return err
}

func (f *FriendGorm) UpdateByMap(ctx context.Context, ownerUserID string, args map[string]interface{}) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "ownerUserID", ownerUserID, "args", args)
	}()
	return utils.Wrap(f.DB.Where("owner_user_id = ?", ownerUserID).Updates(args).Error, "")
}

func (f *FriendGorm) Update(ctx context.Context, friends []*FriendModel) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "friends", friends)
	}()
	return utils.Wrap(f.DB.Updates(&friends).Error, "")
}

func (f *FriendGorm) UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "ownerUserID", ownerUserID, "friendUserID", friendUserID, "remark", remark)
	}()
	return utils.Wrap(f.DB.Model(&FriendModel{}).Where("owner_user_id = ? and friend_user_id = ?", ownerUserID, friendUserID).Update("remark", remark).Error, "")
}

func (f *FriendGorm) FindOwnerUserID(ctx context.Context, ownerUserID string) (friends []*FriendModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "ownerUserID", ownerUserID, "friends", friends)
	}()
	return friends, utils.Wrap(f.DB.Where("owner_user_id = ?", ownerUserID).Find(&friends).Error, "")
}

func (f *FriendGorm) FindFriendUserID(ctx context.Context, friendUserID string) (friends []*FriendModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "friendUserID", friendUserID, "friends", friends)
	}()
	return friends, utils.Wrap(f.DB.Where("friend_user_id = ?", friendUserID).Find(&friends).Error, "")
}

func (f *FriendGorm) Take(ctx context.Context, ownerUserID, friendUserID string) (friend *FriendModel, err error) {
	friend = &FriendModel{}
	defer tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "ownerUserID", ownerUserID, "friendUserID", friendUserID, "friend", friend)
	return friend, utils.Wrap(f.DB.Where("owner_user_id = ? and friend_user_id", ownerUserID, friendUserID).Take(friend).Error, "")
}

func (f *FriendGorm) FindUserState(ctx context.Context, userID1, userID2 string) (friends []*FriendModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "userID1", userID1, "userID2", userID2)
	}()
	return friends, utils.Wrap(f.DB.Where("(owner_user_id = ? and friend_user_id = ?) or (owner_user_id = ? and friend_user_id = ?)", userID1, userID2, userID2, userID1).Find(&friends).Error, "")
}
