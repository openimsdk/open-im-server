package relation

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"gorm.io/gorm"
)

func NewFriendRequestGorm(db *gorm.DB) relation.FriendRequestModelInterface {
	return &FriendRequestGorm{db}
}

type FriendRequestGorm struct {
	DB *gorm.DB
}

func (f *FriendRequestGorm) NewTx(tx any) relation.FriendRequestModelInterface {
	return &FriendRequestGorm{DB: tx.(*gorm.DB)}
}

// 插入多条记录
func (f *FriendRequestGorm) Create(ctx context.Context, friendRequests []*relation.FriendRequestModel) (err error) {
	return utils.Wrap(f.DB.Create(&friendRequests).Error, "")
}

// 删除记录
func (f *FriendRequestGorm) Delete(ctx context.Context, fromUserID, toUserID string) (err error) {
	return utils.Wrap(f.DB.Where("from_user_id = ? AND to_user_id = ?", fromUserID, toUserID).Delete(&relation.FriendRequestModel{}).Error, "")
}

// 更新零值
func (f *FriendRequestGorm) UpdateByMap(ctx context.Context, formUserID string, toUserID string, args map[string]interface{}) (err error) {
	return utils.Wrap(f.DB.Model(&relation.FriendRequestModel{}).Where("from_user_id = ? AND to_user_id ", formUserID, toUserID).Updates(args).Error, "")
}

// 更新多条记录 （非零值）
func (f *FriendRequestGorm) Update(ctx context.Context, friendRequests []*relation.FriendRequestModel) (err error) {
	return utils.Wrap(f.DB.Updates(&friendRequests).Error, "")
}

// 获取来指定用户的好友申请  未找到 不返回错误
func (f *FriendRequestGorm) Find(ctx context.Context, fromUserID, toUserID string) (friendRequest *relation.FriendRequestModel, err error) {
	friendRequest = &relation.FriendRequestModel{}
	utils.Wrap(f.DB.Where("from_user_id = ? and to_user_id", fromUserID, toUserID).Find(friendRequest).Error, "")
	return
}

func (f *FriendRequestGorm) Take(ctx context.Context, fromUserID, toUserID string) (friendRequest *relation.FriendRequestModel, err error) {
	friendRequest = &relation.FriendRequestModel{}
	utils.Wrap(f.DB.Where("from_user_id = ? and to_user_id", fromUserID, toUserID).Take(friendRequest).Error, "")
	return
}

// 获取toUserID收到的好友申请列表
func (f *FriendRequestGorm) FindToUserID(ctx context.Context, toUserID string, pageNumber, showNumber int32) (friendRequests []*relation.FriendRequestModel, total int64, err error) {
	err = f.DB.Model(&relation.FriendRequestModel{}).Where("to_user_id = ? ", toUserID).Count(&total).Error
	if err != nil {
		return nil, 0, utils.Wrap(err, "")
	}
	err = utils.Wrap(f.DB.Where("to_user_id = ? ", toUserID).Limit(int(showNumber)).Offset(int(pageNumber*showNumber)).Find(&friendRequests).Error, "")
	return
}

// 获取fromUserID发出去的好友申请列表
func (f *FriendRequestGorm) FindFromUserID(ctx context.Context, fromUserID string, pageNumber, showNumber int32) (friendRequests []*relation.FriendRequestModel, total int64, err error) {
	err = f.DB.Model(&relation.FriendRequestModel{}).Where("from_user_id = ? ", fromUserID).Count(&total).Error
	if err != nil {
		return nil, 0, utils.Wrap(err, "")
	}
	err = utils.Wrap(f.DB.Where("from_user_id = ? ", fromUserID).Limit(int(showNumber)).Offset(int(pageNumber*showNumber)).Find(&friendRequests).Error, "")
	return
}
