package relation

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"gorm.io/gorm"
)

type FriendRequestGorm struct {
	*MetaDB
}

func NewFriendRequestGorm(db *gorm.DB) relation.FriendRequestModelInterface {
	return &FriendRequestGorm{NewMetaDB(db, &relation.FriendRequestModel{})}
}

func (f *FriendRequestGorm) NewTx(tx any) relation.FriendRequestModelInterface {
	return &FriendRequestGorm{NewMetaDB(tx.(*gorm.DB), &relation.FriendRequestModel{})}
}

// 插入多条记录
func (f *FriendRequestGorm) Create(ctx context.Context, friendRequests []*relation.FriendRequestModel) (err error) {
	return utils.Wrap(f.db(ctx).Create(&friendRequests).Error, "")
}

// 删除记录
func (f *FriendRequestGorm) Delete(ctx context.Context, fromUserID, toUserID string) (err error) {
	return utils.Wrap(f.db(ctx).Where("from_user_id = ? AND to_user_id = ?", fromUserID, toUserID).Delete(&relation.FriendRequestModel{}).Error, "")
}

// 更新零值
func (f *FriendRequestGorm) UpdateByMap(ctx context.Context, fromUserID string, toUserID string, args map[string]interface{}) (err error) {
	return utils.Wrap(f.db(ctx).Model(&relation.FriendRequestModel{}).Where("from_user_id = ? AND to_user_id =?", fromUserID, toUserID).Updates(args).Error, "")
}

// 更新记录 （非零值）
func (f *FriendRequestGorm) Update(ctx context.Context, friendRequest *relation.FriendRequestModel) (err error) {
	return utils.Wrap(f.db(ctx).Where("from_user_id = ? AND to_user_id =?", friendRequest.FromUserID, friendRequest.ToUserID).Updates(friendRequest).Error, "")
}

// 获取来指定用户的好友申请  未找到 不返回错误
func (f *FriendRequestGorm) Find(ctx context.Context, fromUserID, toUserID string) (friendRequest *relation.FriendRequestModel, err error) {
	friendRequest = &relation.FriendRequestModel{}
	err = utils.Wrap(f.db(ctx).Where("from_user_id = ? and to_user_id = ?", fromUserID, toUserID).Find(friendRequest).Error, "")
	return friendRequest, err
}

func (f *FriendRequestGorm) Take(ctx context.Context, fromUserID, toUserID string) (friendRequest *relation.FriendRequestModel, err error) {
	friendRequest = &relation.FriendRequestModel{}
	err = utils.Wrap(f.db(ctx).Where("from_user_id = ? and to_user_id = ?", fromUserID, toUserID).Take(friendRequest).Error, "")
	return friendRequest, err
}

// 获取toUserID收到的好友申请列表
func (f *FriendRequestGorm) FindToUserID(ctx context.Context, toUserID string, pageNumber, showNumber int32) (friendRequests []*relation.FriendRequestModel, total int64, err error) {
	err = f.db(ctx).Model(&relation.FriendRequestModel{}).Where("to_user_id = ? ", toUserID).Count(&total).Error
	if err != nil {
		return nil, 0, utils.Wrap(err, "")
	}
	err = utils.Wrap(f.db(ctx).Where("to_user_id = ? ", toUserID).Limit(int(showNumber)).Offset(int(pageNumber-1)*int(showNumber)).Find(&friendRequests).Error, "")
	return
}

// 获取fromUserID发出去的好友申请列表
func (f *FriendRequestGorm) FindFromUserID(ctx context.Context, fromUserID string, pageNumber, showNumber int32) (friendRequests []*relation.FriendRequestModel, total int64, err error) {
	err = f.db(ctx).Model(&relation.FriendRequestModel{}).Where("from_user_id = ? ", fromUserID).Count(&total).Error
	if err != nil {
		return nil, 0, utils.Wrap(err, "")
	}
	err = utils.Wrap(f.db(ctx).Where("from_user_id = ? ", fromUserID).Limit(int(showNumber)).Offset(int(pageNumber-1)*int(showNumber)).Find(&friendRequests).Error, "")
	return
}
