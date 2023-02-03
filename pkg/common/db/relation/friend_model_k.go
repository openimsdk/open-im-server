package relation

import (
	"Open_IM/pkg/common/db/table/relation"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/utils"
	"context"
	"gorm.io/gorm"
)

type FriendDB interface {
	Create(ctx context.Context, friends []*relation.FriendModel) (err error)
	Delete(ctx context.Context, ownerUserID string, friendUserIDs string) (err error)
	UpdateByMap(ctx context.Context, ownerUserID string, args map[string]interface{}) (err error)
	Update(ctx context.Context, friends []*relation.FriendModel) (err error)
	UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) (err error)
	FindOwnerUserID(ctx context.Context, ownerUserID string) (friends []*relation.FriendModel, err error)
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

func (f *FriendGorm) Create(ctx context.Context, friends []*relation.FriendModel, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "friends", friends)
	}()
	return utils.Wrap(getDBConn(f.DB, tx).Model(&table.FriendModel{}).Create(&friends).Error, "")
}

func (f *FriendGorm) Delete(ctx context.Context, ownerUserID string, friendUserIDs string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "ownerUserID", ownerUserID, "friendUserIDs", friendUserIDs)
	}()
	err = utils.Wrap(f.DB.Model(&table.FriendModel{}).Where("owner_user_id = ? and friend_user_id = ?", ownerUserID, friendUserIDs).Delete(&table.FriendModel{}).Error, "")
	return err
}

func (f *FriendGorm) UpdateByMap(ctx context.Context, ownerUserID string, args map[string]interface{}) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "ownerUserID", ownerUserID, "args", args)
	}()
	return utils.Wrap(f.DB.Model(&table.FriendModel{}).Where("owner_user_id = ?", ownerUserID).Updates(args).Error, "")
}

func (f *FriendGorm) Update(ctx context.Context, friends []*relation.FriendModel, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "friends", friends)
	}()
	return utils.Wrap(f.DB.Model(&table.FriendModel{}).Updates(&friends).Error, "")
}

func (f *FriendGorm) UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "ownerUserID", ownerUserID, "friendUserID", friendUserID, "remark", remark)
	}()
	return utils.Wrap(f.DB.Model(&table.FriendModel{}).Where("owner_user_id = ? and friend_user_id = ?", ownerUserID, friendUserID).Update("remark", remark).Error, "")
}

func (f *FriendGorm) FindOwnerUserID(ctx context.Context, ownerUserID string) (friends []*relation.FriendModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "ownerUserID", ownerUserID, "friends", friends)
	}()
	return friends, utils.Wrap(f.DB.Model(&table.FriendModel{}).Where("owner_user_id = ?", ownerUserID).Find(&friends).Error, "")
}

func (f *FriendGorm) FindFriendUserID(ctx context.Context, friendUserID string) (friends []*relation.FriendModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "friendUserID", friendUserID, "friends", friends)
	}()
	return friends, utils.Wrap(f.DB.Model(&table.FriendModel{}).Where("friend_user_id = ?", friendUserID).Find(&friends).Error, "")
}

func (f *FriendGorm) Take(ctx context.Context, ownerUserID, friendUserID string) (friend *relation.FriendModel, err error) {
	friend = &table.FriendModel{}
	defer tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "ownerUserID", ownerUserID, "friendUserID", friendUserID, "friend", friend)
	return friend, utils.Wrap(f.DB.Model(&table.FriendModel{}).Where("owner_user_id = ? and friend_user_id", ownerUserID, friendUserID).Take(friend).Error, "")
}

func (f *FriendGorm) FindUserState(ctx context.Context, userID1, userID2 string) (friends []*relation.FriendModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "userID1", userID1, "userID2", userID2)
	}()
	return friends, utils.Wrap(f.DB.Model(&table.FriendModel{}).Where("(owner_user_id = ? and friend_user_id = ?) or (owner_user_id = ? and friend_user_id = ?)", userID1, userID2, userID2, userID1).Find(&friends).Error, "")
}

// 获取 owner的好友列表 如果不存在也不返回错误
func (f *FriendGorm) FindFriends(ctx context.Context, ownerUserID string, friendUserIDs []string, tx ...*gorm.DB) (friends []*relation.FriendModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "friendUserIDs", friendUserIDs, "friends", friends)
	}()
	return friends, utils.Wrap(getDBConn(f.DB, tx).Where("owner_user_id = ? AND friend_user_id in (?)", ownerUserID, friendUserIDs).Find(&friends).Error, "")
}

// 获取哪些人添加了friendUserID 如果不存在也不返回错误
func (f *FriendGorm) FindReversalFriends(ctx context.Context, friendUserID string, ownerUserIDs []string, tx ...*gorm.DB) (friends []*relation.FriendModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "friendUserID", friendUserID, "friends", friends)
	}()
	return friends, utils.Wrap(getDBConn(f.DB, tx).Where("friend_user_id = ? AND owner_user_id in (?)", friendUserID, ownerUserIDs).Find(&friends).Error, "")
}

func (f *FriendGorm) FindOwnerFriends(ctx context.Context, ownerUserID string, pageNumber, showNumber int32, tx ...*gorm.DB) (friends []*relation.FriendModel, total int64, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "ownerUserID", ownerUserID, "pageNumber", pageNumber, "showNumber", showNumber, "friends", friends, "total", total)
	}()
	err = getDBConn(f.DB, tx).Model(f).Where("owner_user_id = ? ", ownerUserID).Count(&total).Error
	if err != nil {
		return nil, 0, utils.Wrap(err, "")
	}
	err = utils.Wrap(getDBConn(f.DB, tx).Model(f).Where("owner_user_id = ? ", ownerUserID).Limit(int(showNumber)).Offset(int(pageNumber*showNumber)).Find(&friends).Error, "")
	return
}

func (f *FriendGorm) FindInWhoseFriends(ctx context.Context, friendUserID string, pageNumber, showNumber int32, tx ...*gorm.DB) (friends []*relation.FriendModel, total int64, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "friendUserID", friendUserID, "pageNumber", pageNumber, "showNumber", showNumber, "friends", friends, "total", total)
	}()
	err = getDBConn(f.DB, tx).Model(f).Where("friend_user_id = ? ", friendUserID).Count(&total).Error
	if err != nil {
		return nil, 0, utils.Wrap(err, "")
	}
	err = utils.Wrap(getDBConn(f.DB, tx).Model(f).Where("friend_user_id = ? ", friendUserID).Limit(int(showNumber)).Offset(int(pageNumber*showNumber)).Find(&friends).Error, "")
	return
}
