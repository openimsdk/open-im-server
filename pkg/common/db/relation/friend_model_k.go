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

// 插入多条记录
func (f *FriendGorm) Create(ctx context.Context, friends []*relation.FriendModel, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "friends", friends)
	}()
	return utils.Wrap(getDBConn(f.DB, tx).Create(&friends).Error, "")
}

// 删除ownerUserID指定的好友
func (f *FriendGorm) Delete(ctx context.Context, ownerUserID string, friendUserIDs []string, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "friendUserIDs", friendUserIDs)
	}()
	err = utils.Wrap(getDBConn(f.DB, tx).Where("owner_user_id = ? AND friend_user_id in ( ?)", ownerUserID, friendUserIDs).Delete(&relation.FriendModel{}).Error, "")
	return err
}

// 更新ownerUserID单个好友信息 更新零值
func (f *FriendGorm) UpdateByMap(ctx context.Context, ownerUserID string, friendUserID string, args map[string]interface{}, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "friendUserID", friendUserID, "args", args)
	}()
	return utils.Wrap(getDBConn(f.DB, tx).Model(&relation.FriendModel{}).Where("owner_user_id = ? AND friend_user_id = ? ", ownerUserID, friendUserID).Updates(args).Error, "")
}

// 更新好友信息的非零值
func (f *FriendGorm) Update(ctx context.Context, friends []*relation.FriendModel, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "friends", friends)
	}()
	return utils.Wrap(getDBConn(f.DB, tx).Updates(&friends).Error, "")
}

// 更新好友备注（也支持零值 ）
func (f *FriendGorm) UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "friendUserID", friendUserID, "remark", remark)
	}()
	if remark != "" {
		return utils.Wrap(getDBConn(f.DB, tx).Model(&relation.FriendModel{}).Where("owner_user_id = ? and friend_user_id = ?", ownerUserID, friendUserID).Update("remark", remark).Error, "")
	}
	m := make(map[string]interface{}, 1)
	m["remark"] = ""
	return utils.Wrap(getDBConn(f.DB, tx).Model(&relation.FriendModel{}).Where("owner_user_id = ?", ownerUserID).Updates(m).Error, "")

}

// 获取单个好友信息，如没找到 返回错误
func (f *FriendGorm) Take(ctx context.Context, ownerUserID, friendUserID string, tx ...*gorm.DB) (friend *relation.FriendModel, err error) {
	friend = &relation.FriendModel{}
	defer tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "friendUserID", friendUserID, "friend", *friend)
	return friend, utils.Wrap(getDBConn(f.DB, tx).Where("owner_user_id = ? and friend_user_id", ownerUserID, friendUserID).Take(friend).Error, "")
}

// 查找好友关系，如果是双向关系，则都返回
func (f *FriendGorm) FindUserState(ctx context.Context, userID1, userID2 string, tx ...*gorm.DB) (friends []*relation.FriendModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userID1", userID1, "userID2", userID2, "friends", friends)
	}()
	return friends, utils.Wrap(getDBConn(f.DB, tx).Where("(owner_user_id = ? and friend_user_id = ?) or (owner_user_id = ? and friend_user_id = ?)", userID1, userID2, userID2, userID1).Find(&friends).Error, "")
}

// 获取 owner指定的好友列表 如果有friendUserIDs不存在，也不返回错误
func (f *FriendGorm) FindFriends(ctx context.Context, ownerUserID string, friendUserIDs []string, tx ...*gorm.DB) (friends []*relation.FriendModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "friendUserIDs", friendUserIDs, "friends", friends)
	}()
	return friends, utils.Wrap(getDBConn(f.DB, tx).Where("owner_user_id = ? AND friend_user_id in (?)", ownerUserID, friendUserIDs).Find(&friends).Error, "")
}

// 获取哪些人添加了friendUserID 如果有ownerUserIDs不存在，也不返回错误
func (f *FriendGorm) FindReversalFriends(ctx context.Context, friendUserID string, ownerUserIDs []string, tx ...*gorm.DB) (friends []*relation.FriendModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "friendUserID", friendUserID, "ownerUserIDs", ownerUserIDs, "friends", friends)
	}()
	return friends, utils.Wrap(getDBConn(f.DB, tx).Where("friend_user_id = ? AND owner_user_id in (?)", friendUserID, ownerUserIDs).Find(&friends).Error, "")
}

// 获取ownerUserID好友列表 支持翻页
func (f *FriendGorm) FindOwnerFriends(ctx context.Context, ownerUserID string, pageNumber, showNumber int32, tx ...*gorm.DB) (friends []*relation.FriendModel, total int64, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "pageNumber", pageNumber, "showNumber", showNumber, "friends", friends, "total", total)
	}()
	err = getDBConn(f.DB, tx).Model(&relation.FriendModel{}).Where("owner_user_id = ? ", ownerUserID).Count(&total).Error
	if err != nil {
		return nil, 0, utils.Wrap(err, "")
	}
	err = utils.Wrap(getDBConn(f.DB, tx).Where("owner_user_id = ? ", ownerUserID).Limit(int(showNumber)).Offset(int(pageNumber*showNumber)).Find(&friends).Error, "")
	return
}

// 获取哪些人添加了friendUserID 支持翻页
func (f *FriendGorm) FindInWhoseFriends(ctx context.Context, friendUserID string, pageNumber, showNumber int32, tx ...*gorm.DB) (friends []*relation.FriendModel, total int64, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "friendUserID", friendUserID, "pageNumber", pageNumber, "showNumber", showNumber, "friends", friends, "total", total)
	}()
	err = getDBConn(f.DB, tx).Model(&relation.FriendModel{}).Where("friend_user_id = ? ", friendUserID).Count(&total).Error
	if err != nil {
		return nil, 0, utils.Wrap(err, "")
	}
	err = utils.Wrap(getDBConn(f.DB, tx).Where("friend_user_id = ? ", friendUserID).Limit(int(showNumber)).Offset(int(pageNumber*showNumber)).Find(&friends).Error, "")
	return
}
