package relation

import (
	"Open_IM/pkg/common/db/table/relation"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/utils"
	"context"
	"gorm.io/gorm"
)

//var FriendRequestDB *gorm.DB

func NewFriendRequestGorm(db *gorm.DB) *FriendRequestGorm {
	var fr FriendRequestGorm
	fr.DB = db
	return &fr
}

type FriendRequestGorm struct {
	DB *gorm.DB `gorm:"-"`
}

func (f *FriendRequestGorm) Create(ctx context.Context, friends []*relation.FriendRequestModel, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "friends", friends)
	}()
	return utils.Wrap(f.DB.Model(&relation.FriendRequestModel{}).Create(&friends).Error, "")
}

func (f *FriendRequestGorm) Delete(ctx context.Context, fromUserID, toUserID string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "fromUserID", fromUserID, "toUserID", toUserID)
	}()
	return utils.Wrap(f.DB.Model(&relation.FriendRequestModel{}).Where("from_user_id = ? and to_user_id = ?", fromUserID, toUserID).Delete(&relation.FriendRequestModel{}).Error, "")
}

func (f *FriendRequestGorm) UpdateByMap(ctx context.Context, formUserID string, toUserID string, args map[string]interface{}, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "formUserID", formUserID, "toUserID", toUserID, "args", args)
	}()
	return utils.Wrap(getDBConn(f.DB, tx).Model(&relation.FriendRequestModel{}).Where("from_user_id = ? AND to_user_id ", formUserID, toUserID).Updates(args).Error, "")
}

func (f *FriendRequestGorm) Update(ctx context.Context, friendRequests []*relation.FriendRequestModel, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "friendRequests", friendRequests)
	}()
	return utils.Wrap(getDBConn(f.DB, tx).Model(&relation.FriendRequestModel{}).Updates(&friendRequests).Error, "")
}

func (f *FriendRequestGorm) Take(ctx context.Context, fromUserID, toUserID string) (friend *relation.FriendRequestModel, err error) {
	friend = &relation.FriendRequestModel{}
	defer tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "fromUserID", fromUserID, "toUserID", toUserID, "friend", friend)
	return friend, utils.Wrap(f.DB.Model(&relation.FriendRequestModel{}).Where("from_user_id = ? and to_user_id", fromUserID, toUserID).Take(friend).Error, "")
}

func (f *FriendRequestGorm) Find(ctx context.Context, fromUserID, toUserID string, tx ...*gorm.DB) (friend *relation.FriendRequestModel, err error) {
	friend = &relation.FriendRequestModel{}
	defer tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "fromUserID", fromUserID, "toUserID", toUserID, "friend", friend)
	return friend, utils.Wrap(getDBConn(f.DB, tx).Model(&relation.FriendRequestModel{}).Where("from_user_id = ? and to_user_id", fromUserID, toUserID).Find(friend).Error, "")
}

func (f *FriendRequestGorm) FindToUserID(ctx context.Context, toUserID string, pageNumber, showNumber int32, tx ...*gorm.DB) (friends []*relation.FriendRequestModel, total int64, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "toUserID", toUserID, "friends", friends)
	}()

	err = getDBConn(f.DB, tx).Model(&relation.FriendRequestModel{}).Where("to_user_id = ? ", toUserID).Count(&total).Error
	if err != nil {
		return nil, 0, utils.Wrap(err, "")
	}
	err = utils.Wrap(getDBConn(f.DB, tx).Model(&relation.FriendRequestModel{}).Where("to_user_id = ? ", toUserID).Limit(int(showNumber)).Offset(int(pageNumber*showNumber)).Find(&friends).Error, "")
	return
}

func (f *FriendRequestGorm) FindFromUserID(ctx context.Context, fromUserID string, pageNumber, showNumber int32, tx ...*gorm.DB) (friends []*relation.FriendRequestModel, total int64, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "fromUserID", fromUserID, "friends", friends)
	}()

	err = getDBConn(f.DB, tx).Model(&relation.FriendRequestModel{}).Where("from_user_id = ? ", fromUserID).Count(&total).Error
	if err != nil {
		return nil, 0, utils.Wrap(err, "")
	}
	err = utils.Wrap(getDBConn(f.DB, tx).Model(&relation.FriendRequestModel{}).Where("from_user_id = ? ", fromUserID).Limit(int(showNumber)).Offset(int(pageNumber*showNumber)).Find(&friends).Error, "")
	return
}
