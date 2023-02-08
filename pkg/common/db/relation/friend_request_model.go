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

// 插入多条记录
func (f *FriendRequestGorm) Create(ctx context.Context, friendRequests []*relation.FriendRequestModel, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "friendRequests", friendRequests)
	}()
	return utils.Wrap(getDBConn(f.DB, tx).Create(&friendRequests).Error, "")
}

// 删除记录
func (f *FriendRequestGorm) Delete(ctx context.Context, fromUserID, toUserID string, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "fromUserID", fromUserID, "toUserID", toUserID)
	}()
	return utils.Wrap(getDBConn(f.DB, tx).Where("from_user_id = ? AND to_user_id = ?", fromUserID, toUserID).Delete(&relation.FriendRequestModel{}).Error, "")
}

// 更新零值
func (f *FriendRequestGorm) UpdateByMap(ctx context.Context, formUserID string, toUserID string, args map[string]interface{}, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "formUserID", formUserID, "toUserID", toUserID, "args", args)
	}()
	return utils.Wrap(getDBConn(f.DB, tx).Model(&relation.FriendRequestModel{}).Where("from_user_id = ? AND to_user_id ", formUserID, toUserID).Updates(args).Error, "")
}

// 更新多条记录 （非零值）
func (f *FriendRequestGorm) Update(ctx context.Context, friendRequests []*relation.FriendRequestModel, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "friendRequests", friendRequests)
	}()
	return utils.Wrap(getDBConn(f.DB, tx).Updates(&friendRequests).Error, "")
}

// 获取来指定用户的好友申请  未找到 不返回错误
func (f *FriendRequestGorm) Find(ctx context.Context, fromUserID, toUserID string, tx ...*gorm.DB) (friendRequest *relation.FriendRequestModel, err error) {
	friendRequest = &relation.FriendRequestModel{}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "fromUserID", fromUserID, "toUserID", toUserID, "friendRequest", *friendRequest)
	}()
	utils.Wrap(getDBConn(f.DB, tx).Where("from_user_id = ? and to_user_id", fromUserID, toUserID).Find(friendRequest).Error, "")
	return
}

func (f *FriendRequestGorm) Take(ctx context.Context, fromUserID, toUserID string, tx ...*gorm.DB) (friendRequest *relation.FriendRequestModel, err error) {
	friendRequest = &relation.FriendRequestModel{}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "fromUserID", fromUserID, "toUserID", toUserID, "friendRequest", *friendRequest)
	}()
	utils.Wrap(getDBConn(f.DB, tx).Where("from_user_id = ? and to_user_id", fromUserID, toUserID).Take(friendRequest).Error, "")
	return
}

// 获取toUserID收到的好友申请列表
func (f *FriendRequestGorm) FindToUserID(ctx context.Context, toUserID string, pageNumber, showNumber int32, tx ...*gorm.DB) (friendRequests []*relation.FriendRequestModel, total int64, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "toUserID", toUserID, "friendRequests", friendRequests)
	}()

	err = getDBConn(f.DB, tx).Model(&relation.FriendRequestModel{}).Where("to_user_id = ? ", toUserID).Count(&total).Error
	if err != nil {
		return nil, 0, utils.Wrap(err, "")
	}
	err = utils.Wrap(getDBConn(f.DB, tx).Where("to_user_id = ? ", toUserID).Limit(int(showNumber)).Offset(int(pageNumber*showNumber)).Find(&friendRequests).Error, "")
	return
}

// 获取fromUserID发出去的好友申请列表
func (f *FriendRequestGorm) FindFromUserID(ctx context.Context, fromUserID string, pageNumber, showNumber int32, tx ...*gorm.DB) (friendRequests []*relation.FriendRequestModel, total int64, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "fromUserID", fromUserID, "friendRequests", friendRequests)
	}()
	err = getDBConn(f.DB, tx).Model(&relation.FriendRequestModel{}).Where("from_user_id = ? ", fromUserID).Count(&total).Error
	if err != nil {
		return nil, 0, utils.Wrap(err, "")
	}
	err = utils.Wrap(getDBConn(f.DB, tx).Where("from_user_id = ? ", fromUserID).Limit(int(showNumber)).Offset(int(pageNumber*showNumber)).Find(&friendRequests).Error, "")
	return
}
