package relation

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"gorm.io/gorm"
)

type FriendGorm struct {
	*MetaDB
}

func NewFriendGorm(db *gorm.DB) relation.FriendModelInterface {
	return &FriendGorm{NewMetaDB(db, &relation.FriendModel{})}
}

func (f *FriendGorm) NewTx(tx any) relation.FriendModelInterface {
	return &FriendGorm{NewMetaDB(tx.(*gorm.DB), &relation.FriendModel{})}
}

// 插入多条记录
func (f *FriendGorm) Create(ctx context.Context, friends []*relation.FriendModel) (err error) {
	return utils.Wrap(f.db(ctx).Create(&friends).Error, "")
}

// 删除ownerUserID指定的好友
func (f *FriendGorm) Delete(ctx context.Context, ownerUserID string, friendUserIDs []string) (err error) {
	err = utils.Wrap(f.db(ctx).Where("owner_user_id = ? AND friend_user_id in ( ?)", ownerUserID, friendUserIDs).Delete(&relation.FriendModel{}).Error, "")
	return err
}

// 更新ownerUserID单个好友信息 更新零值
func (f *FriendGorm) UpdateByMap(ctx context.Context, ownerUserID string, friendUserID string, args map[string]interface{}) (err error) {
	return utils.Wrap(f.db(ctx).Where("owner_user_id = ? AND friend_user_id = ? ", ownerUserID, friendUserID).Updates(args).Error, "")
}

// 更新好友信息的非零值
func (f *FriendGorm) Update(ctx context.Context, friends []*relation.FriendModel) (err error) {
	return utils.Wrap(f.db(ctx).Updates(&friends).Error, "")
}

// 更新好友备注（也支持零值 ）
func (f *FriendGorm) UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) (err error) {
	if remark != "" {
		return utils.Wrap(f.db(ctx).Where("owner_user_id = ? and friend_user_id = ?", ownerUserID, friendUserID).Update("remark", remark).Error, "")
	}
	m := make(map[string]interface{}, 1)
	m["remark"] = ""
	return utils.Wrap(f.db(ctx).Where("owner_user_id = ?", ownerUserID).Updates(m).Error, "")
}

// 获取单个好友信息，如没找到 返回错误
func (f *FriendGorm) Take(ctx context.Context, ownerUserID, friendUserID string) (friend *relation.FriendModel, err error) {
	friend = &relation.FriendModel{}
	return friend, utils.Wrap(f.db(ctx).Where("owner_user_id = ? and friend_user_id", ownerUserID, friendUserID).Take(friend).Error, "")
}

// 查找好友关系，如果是双向关系，则都返回
func (f *FriendGorm) FindUserState(ctx context.Context, userID1, userID2 string) (friends []*relation.FriendModel, err error) {
	return friends, utils.Wrap(f.db(ctx).Where("(owner_user_id = ? and friend_user_id = ?) or (owner_user_id = ? and friend_user_id = ?)", userID1, userID2, userID2, userID1).Find(&friends).Error, "")
}

// 获取 owner指定的好友列表 如果有friendUserIDs不存在，也不返回错误
func (f *FriendGorm) FindFriends(ctx context.Context, ownerUserID string, friendUserIDs []string) (friends []*relation.FriendModel, err error) {
	return friends, utils.Wrap(f.db(ctx).Where("owner_user_id = ? AND friend_user_id in (?)", ownerUserID, friendUserIDs).Find(&friends).Error, "")
}

// 获取哪些人添加了friendUserID 如果有ownerUserIDs不存在，也不返回错误
func (f *FriendGorm) FindReversalFriends(ctx context.Context, friendUserID string, ownerUserIDs []string) (friends []*relation.FriendModel, err error) {
	return friends, utils.Wrap(f.db(ctx).Where("friend_user_id = ? AND owner_user_id in (?)", friendUserID, ownerUserIDs).Find(&friends).Error, "")
}

// 获取ownerUserID好友列表 支持翻页
func (f *FriendGorm) FindOwnerFriends(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (friends []*relation.FriendModel, total int64, err error) {
	err = f.DB.Model(&relation.FriendModel{}).Where("owner_user_id = ? ", ownerUserID).Count(&total).Error
	if err != nil {
		return nil, 0, utils.Wrap(err, "")
	}
	err = utils.Wrap(f.db(ctx).Where("owner_user_id = ? ", ownerUserID).Limit(int(showNumber)).Offset(int((pageNumber-1)*showNumber)).Find(&friends).Error, "")
	return
}

// 获取哪些人添加了friendUserID 支持翻页
func (f *FriendGorm) FindInWhoseFriends(ctx context.Context, friendUserID string, pageNumber, showNumber int32) (friends []*relation.FriendModel, total int64, err error) {
	err = f.DB.Model(&relation.FriendModel{}).Where("friend_user_id = ? ", friendUserID).Count(&total).Error
	if err != nil {
		return nil, 0, utils.Wrap(err, "")
	}
	err = utils.Wrap(f.db(ctx).Where("friend_user_id = ? ", friendUserID).Limit(int(showNumber)).Offset(int((pageNumber-1)*showNumber)).Find(&friends).Error, "")
	return
}

func (f *FriendGorm) FindFriendUserIDs(ctx context.Context, ownerUserID string) (friendUserIDs []string, err error) {
	return friendUserIDs, utils.Wrap(f.db(ctx).Model(&relation.FriendModel{}).Where("owner_user_id = ? ", ownerUserID).Pluck("friend_user_id", &friendUserIDs).Error, "")
}
