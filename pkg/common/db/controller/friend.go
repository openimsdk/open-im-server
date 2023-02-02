package controller

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/relation"
	"Open_IM/pkg/common/db/table"
	"context"
	"gorm.io/gorm"
)

type FriendInterface interface {
	// CheckIn 检查user2是否在user1的好友列表中(inUser1Friends==true) 检查user1是否在user2的好友列表中(inUser2Friends==true)
	CheckIn(ctx context.Context, user1, user2 string) (inUser1Friends bool, inUser2Friends bool, err error)
	// AddFriendRequest 增加或者更新好友申请
	AddFriendRequest(ctx context.Context, fromUserID, toUserID string, reqMsg string, ex string) (err error)
	// BecomeFriend 先判断是否在好友表，如果在则不插入
	BecomeFriend(ctx context.Context, friends []*table.FriendModel, revFriends []*table.FriendModel) (err error)
	// RefuseFriendRequest 拒绝好友申请
	RefuseFriendRequest(ctx context.Context, friendRequest *table.FriendRequestModel) (err error)
	// AgreeFriendRequest 同意好友申请
	AgreeFriendRequest(ctx context.Context, friendRequest *table.FriendRequestModel) (err error)
	// Delete 删除好友
	Delete(ctx context.Context, ownerUserID string, friendUserIDs string) (err error)
	// UpdateRemark 更新好友备注
	UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) (err error)
	// FindOwnerFriends 获取ownerUserID的好友列表
	FindOwnerFriends(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (friends []*table.FriendModel, total int64, err error)
	// FindInWhoseFriends friendUserID在哪些人的好友列表中
	FindInWhoseFriends(ctx context.Context, friendUserID string, pageNumber, showNumber int32) (friends []*table.FriendModel, total int64, err error)
	// FindFriendRequestFromMe 获取我发出去的好友申请
	FindFriendRequestFromMe(ctx context.Context, userID string, pageNumber, showNumber int32) (friends []*table.FriendRequestModel, total int64, err error)
	// FindFriendRequestToMe 获取我收到的的好友申请
	FindFriendRequestToMe(ctx context.Context, userID string, pageNumber, showNumber int32) (friends []*table.FriendRequestModel, total int64, err error)
	// FindFriends 获取某人指定好友的信息 如果有一个不存在也返回错误
	FindFriends(ctx context.Context, ownerUserID string, friendUserIDs []string) (friends []*table.FriendModel, err error)
}

type FriendController struct {
	database FriendDatabaseInterface
}

func NewFriendController(db *gorm.DB) *FriendController {
	return &FriendController{database: NewFriendDatabase(db)}
}

// CheckIn 检查user2是否在user1的好友列表中(inUser1Friends==true) 检查user1是否在user2的好友列表中(inUser2Friends==true)
func (f *FriendController) CheckIn(ctx context.Context, user1, user2 string) (err error, inUser1Friends bool, inUser2Friends bool) {
}

// AddFriendRequest 增加或者更新好友申请
func (f *FriendController) AddFriendRequest(ctx context.Context, fromUserID, toUserID string, reqMsg string, ex string) (err error) {
}

// BecomeFriend 先判断是否在好友表，如果在则不插入
func (f *FriendController) BecomeFriend(ctx context.Context, friends []*table.FriendModel) (err error) {
}

// RefuseFriendRequest 拒绝好友申请
func (f *FriendController) RefuseFriendRequest(ctx context.Context, friendRequest *table.FriendRequestModel) (err error) {
}

// AgreeFriendRequest 同意好友申请
func (f *FriendController) AgreeFriendRequest(ctx context.Context, friendRequest *table.FriendRequestModel) (err error) {
}

// Delete 删除好友
func (f *FriendController) Delete(ctx context.Context, ownerUserID string, friendUserIDs string) (err error) {
}

// UpdateRemark 更新好友备注
func (f *FriendController) UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) (err error) {
}

// FindOwnerFriends 获取ownerUserID的好友列表
func (f *FriendController) FindOwnerFriends(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (friends []*table.FriendModel, err error) {
}

// FindInWhoseFriends friendUserID在哪些人的好友列表中
func (f *FriendController) FindInWhoseFriends(ctx context.Context, friendUserID string, pageNumber, showNumber int32) (friends []*table.FriendModel, err error) {
}

// FindFriendRequestFromMe 获取我发出去的好友申请
func (f *FriendController) FindFriendRequestFromMe(ctx context.Context, userID string, pageNumber, showNumber int32) (friends []*table.FriendRequestModel, err error) {
}

// FindFriendRequestToMe 获取我收到的的好友申请
func (f *FriendController) FindFriendRequestToMe(ctx context.Context, userID string, pageNumber, showNumber int32) (friends []*table.FriendRequestModel, err error) {
}

// FindFriends 获取某人指定好友的信息
func (f *FriendController) FindFriends(ctx context.Context, ownerUserID string, friendUserIDs []string) (friends []*table.FriendModel, err error) {
}

type FriendDatabaseInterface interface {
	// CheckIn 检查user2是否在user1的好友列表中(inUser1Friends==true) 检查user1是否在user2的好友列表中(inUser2Friends==true)
	CheckIn(ctx context.Context, user1, user2 string) (inUser1Friends bool, inUser2Friends bool, err error)
	// AddFriendRequest 增加或者更新好友申请
	AddFriendRequest(ctx context.Context, fromUserID, toUserID string, reqMsg string, ex string) (err error)
	// BecomeFriend 先判断是否在好友表，如果在则不插入
	BecomeFriend(ctx context.Context, friends []*table.FriendModel) (err error)
	// RefuseFriendRequest 拒绝好友申请
	RefuseFriendRequest(ctx context.Context, friendRequest *table.FriendRequestModel) (err error)
	// AgreeFriendRequest 同意好友申请
	AgreeFriendRequest(ctx context.Context, friendRequest *table.FriendRequestModel) (err error)
	// Delete 删除好友
	Delete(ctx context.Context, ownerUserID string, friendUserIDs string) (err error)
	// UpdateRemark 更新好友备注
	UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) (err error)
	// FindOwnerFriends 获取ownerUserID的好友列表
	FindOwnerFriends(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (friends []*table.FriendModel, err error)
	// FindInWhoseFriends friendUserID在哪些人的好友列表中
	FindInWhoseFriends(ctx context.Context, friendUserID string, pageNumber, showNumber int32) (friends []*table.FriendModel, err error)
	// FindFriendRequestFromMe 获取我发出去的好友申请
	FindFriendRequestFromMe(ctx context.Context, userID string, pageNumber, showNumber int32) (friends []*table.FriendRequestModel, err error)
	// FindFriendRequestToMe 获取我收到的的好友申请
	FindFriendRequestToMe(ctx context.Context, userID string, pageNumber, showNumber int32) (friends []*table.FriendRequestModel, err error)
	// FindFriends 获取某人指定好友的信息
	FindFriends(ctx context.Context, ownerUserID string, friendUserIDs []string) (friends []*table.FriendModel, err error)
}

type FriendDatabase struct {
	friend        *relation.FriendGorm
	friendRequest *relation.FriendRequestGorm
}

func NewFriendDatabase(db *gorm.DB) *FriendDatabase {
	return &FriendDatabase{friend: relation.NewFriendGorm(db), friendRequest: relation.NewFriendRequestGorm(db)}
}

// CheckIn 检查user2是否在user1的好友列表中(inUser1Friends==true) 检查user1是否在user2的好友列表中(inUser2Friends==true)
func (f *FriendDatabase) CheckIn(ctx context.Context, userID1, userID2 string) (inUser1Friends bool, inUser2Friends bool, err error) {
	friends, err := f.friend.FindUserState(ctx, userID1, userID2)
	for _, v := range friends {
		if v.OwnerUserID == userID1 && v.FriendUserID == userID2 {
			inUser1Friends = true
		}
		if v.OwnerUserID == userID2 && v.FriendUserID == userID1 {
			inUser2Friends = true
		}
	}
	return
}

// AddFriendRequest 增加或者更新好友申请
func (f *FriendDatabase) AddFriendRequest(ctx context.Context, fromUserID, toUserID string, reqMsg string, ex string) (err error) {

}

// BecomeFriend 先判断是否在好友表，如果在则不插入
func (f *FriendDatabase) BecomeFriend(ctx context.Context, ownerUserID string, friends []*table.FriendModel) (err error) {
	return f.friend.DB.Transaction(func(tx *gorm.DB) error {
		//先find 找出重复的 去掉重复的
		friendUserIDs := make([]string, 0, len(friends))
		for _, v := range friends {
			friendUserIDs = append(friendUserIDs, v.FriendUserID)
		}
		fs1, err := f.friend.FindFriends(ctx, ownerUserID, friendUserIDs, tx)
		if err != nil {
			return err
		}
		fs2, err := f.friend.FindReversalFriends(ctx, ownerUserID, friendUserIDs, tx)
		if err != nil {
			return err
		}

		return nil
	})

}

// RefuseFriendRequest 拒绝好友申请
func (f *FriendDatabase) RefuseFriendRequest(ctx context.Context, friendRequest *table.FriendRequestModel) (err error) {
}

// AgreeFriendRequest 同意好友申请
func (f *FriendDatabase) AgreeFriendRequest(ctx context.Context, friendRequest *table.FriendRequestModel) (err error) {
}

// Delete 删除好友
func (f *FriendDatabase) Delete(ctx context.Context, ownerUserID string, friendUserIDs string) (err error) {
}

// UpdateRemark 更新好友备注
func (f *FriendDatabase) UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) (err error) {
}

// FindOwnerFriends 获取ownerUserID的好友列表
func (f *FriendDatabase) FindOwnerFriends(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (friends []*table.FriendModel, err error) {
}

// FindInWhoseFriends friendUserID在哪些人的好友列表中
func (f *FriendDatabase) FindInWhoseFriends(ctx context.Context, friendUserID string, pageNumber, showNumber int32) (friends []*table.FriendModel, err error) {
}

// FindFriendRequestFromMe 获取我发出去的好友申请
func (f *FriendDatabase) FindFriendRequestFromMe(ctx context.Context, userID string, pageNumber, showNumber int32) (friends []*table.FriendRequestModel, err error) {
}

// FindFriendRequestToMe 获取我收到的的好友申请
func (f *FriendDatabase) FindFriendRequestToMe(ctx context.Context, userID string, pageNumber, showNumber int32) (friends []*table.FriendRequestModel, err error) {
}

// FindFriends 获取某人指定好友的信息 如果有一个不存在也返回错误
func (f *FriendDatabase) FindFriends(ctx context.Context, ownerUserID string, friendUserIDs []string) (friends []*table.FriendModel, err error) {
	friends, err = f.friend.Find(ctx, ownerUserID, friendUserIDs)
	if err != nil {
		return
	}
	if len(friends) != len(friendUserIDs) {
		err = constant.ErrRecordNotFound.Wrap()
	}
	return
}
