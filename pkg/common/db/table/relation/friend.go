package relation

import (
	"context"
	"time"
)

const (
	FriendModelTableName = "friends"
)

type FriendModel struct {
	OwnerUserID    string    `gorm:"column:owner_user_id;primary_key;size:64"`
	FriendUserID   string    `gorm:"column:friend_user_id;primary_key;size:64"`
	Remark         string    `gorm:"column:remark;size:255"`
	CreateTime     time.Time `gorm:"column:create_time;autoCreateTime"`
	AddSource      int32     `gorm:"column:add_source"`
	OperatorUserID string    `gorm:"column:operator_user_id;size:64"`
	Ex             string    `gorm:"column:ex;size:1024"`
}

func (FriendModel) TableName() string {
	return FriendModelTableName
}

type FriendModelInterface interface {
	// 插入多条记录
	Create(ctx context.Context, friends []*FriendModel) (err error)
	// 删除ownerUserID指定的好友
	Delete(ctx context.Context, ownerUserID string, friendUserIDs []string) (err error)
	// 更新ownerUserID单个好友信息 更新零值
	UpdateByMap(ctx context.Context, ownerUserID string, friendUserID string, args map[string]interface{}) (err error)
	// 更新好友信息的非零值
	Update(ctx context.Context, friends []*FriendModel) (err error)
	// 更新好友备注（也支持零值 ）
	UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) (err error)
	// 获取单个好友信息，如没找到 返回错误
	Take(ctx context.Context, ownerUserID, friendUserID string) (friend *FriendModel, err error)
	// 查找好友关系，如果是双向关系，则都返回
	FindUserState(ctx context.Context, userID1, userID2 string) (friends []*FriendModel, err error)
	// 获取 owner指定的好友列表 如果有friendUserIDs不存在，也不返回错误
	FindFriends(ctx context.Context, ownerUserID string, friendUserIDs []string) (friends []*FriendModel, err error)
	// 获取哪些人添加了friendUserID 如果有ownerUserIDs不存在，也不返回错误
	FindReversalFriends(ctx context.Context, friendUserID string, ownerUserIDs []string) (friends []*FriendModel, err error)
	// 获取ownerUserID好友列表 支持翻页
	FindOwnerFriends(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (friends []*FriendModel, total int64, err error)
	// 获取哪些人添加了friendUserID 支持翻页
	FindInWhoseFriends(ctx context.Context, friendUserID string, pageNumber, showNumber int32) (friends []*FriendModel, total int64, err error)
	// 获取好友UserID列表
	FindFriendUserIDs(ctx context.Context, ownerUserID string) (friendUserIDs []string, err error)
	NewTx(tx any) FriendModelInterface
}
