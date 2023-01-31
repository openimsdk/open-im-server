package controller

import (
	"Open_IM/pkg/common/db/relation"
	"context"
	"gorm.io/gorm"
)

type FriendInterface interface {
	// CheckIn 检查user2是否在user1的好友列表中(inUser1Friends==true) 检查user1是否在user2的好友列表中(inUser2Friends==true)
	CheckIn(ctx context.Context, user1, user2 string) (err error, inUser1Friends bool, inUser2Friends bool)
	// AddFriendRequest 增加或者更新好友申请
	AddFriendRequest(ctx context.Context, fromUserID, toUserID string, reqMsg string, ex string) (err error)
	// BecomeFriend 先判断是否在好友表，如果在则不插入
	BecomeFriend(ctx context.Context, friends []*relation.Friend) (err error)
	// RefuseFriendRequest 拒绝好友申请
	RefuseFriendRequest(ctx context.Context, friendRequest *relation.FriendRequest) (err error)
	// AgreeFriendRequest 同意好友申请
	AgreeFriendRequest(ctx context.Context, friendRequest *relation.FriendRequest) (err error)
	// Delete 删除好友
	Delete(ctx context.Context, ownerUserID string, friendUserIDs string) (err error)
	// UpdateRemark 更新好友备注
	UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) (err error)
	// FindOwnerFriends 获取ownerUserID的好友列表
	FindOwnerFriends(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (friends []*relation.Friend, err error)
	// FindInWhoseFriends friendUserID在哪些人的好友列表中
	FindInWhoseFriends(ctx context.Context, friendUserID string, pageNumber, showNumber int32) (friends []*relation.Friend, err error)
	// FindFriendRequestFromMe 获取我发出去的好友申请
	FindFriendRequestFromMe(ctx context.Context, userID string, pageNumber, showNumber int32) (friends []*relation.FriendRequest, err error)
	// FindFriendRequestToMe 获取我收到的的好友申请
	FindFriendRequestToMe(ctx context.Context, userID string, pageNumber, showNumber int32) (friends []*relation.FriendRequest, err error)
	// FindFriends 获取某人指定好友的信息
	FindFriends(ctx context.Context, ownerUserID string, friendUserIDs []string) (friends []*relation.Friend, err error)
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
func (f *FriendController) BecomeFriend(ctx context.Context, friends []*relation.Friend) (err error) {
}

// RefuseFriendRequest 拒绝好友申请
func (f *FriendController) RefuseFriendRequest(ctx context.Context, friendRequest *relation.FriendRequest) (err error) {
}

// AgreeFriendRequest 同意好友申请
func (f *FriendController) AgreeFriendRequest(ctx context.Context, friendRequest *relation.FriendRequest) (err error) {
}

// Delete 删除好友
func (f *FriendController) Delete(ctx context.Context, ownerUserID string, friendUserIDs string) (err error) {
}

// UpdateRemark 更新好友备注
func (f *FriendController) UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) (err error) {
}

// FindOwnerFriends 获取ownerUserID的好友列表
func (f *FriendController) FindOwnerFriends(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (friends []*relation.Friend, err error) {
}

// FindInWhoseFriends friendUserID在哪些人的好友列表中
func (f *FriendController) FindInWhoseFriends(ctx context.Context, friendUserID string, pageNumber, showNumber int32) (friends []*relation.Friend, err error) {
}

// FindFriendRequestFromMe 获取我发出去的好友申请
func (f *FriendController) FindFriendRequestFromMe(ctx context.Context, userID string, pageNumber, showNumber int32) (friends []*relation.FriendRequest, err error) {
}

// FindFriendRequestToMe 获取我收到的的好友申请
func (f *FriendController) FindFriendRequestToMe(ctx context.Context, userID string, pageNumber, showNumber int32) (friends []*relation.FriendRequest, err error) {
}

// FindFriends 获取某人指定好友的信息
func (f *FriendController) FindFriends(ctx context.Context, ownerUserID string, friendUserIDs []string) (friends []*relation.Friend, err error) {
}

type FriendDatabaseInterface interface {
	// CheckIn 检查user2是否在user1的好友列表中(inUser1Friends==true) 检查user1是否在user2的好友列表中(inUser2Friends==true)
	CheckIn(ctx context.Context, user1, user2 string) (err error, inUser1Friends bool, inUser2Friends bool)
	// AddFriendRequest 增加或者更新好友申请
	AddFriendRequest(ctx context.Context, fromUserID, toUserID string, reqMsg string, ex string) (err error)
	// BecomeFriend 先判断是否在好友表，如果在则不插入
	BecomeFriend(ctx context.Context, friends []*relation.Friend) (err error)
	// RefuseFriendRequest 拒绝好友申请
	RefuseFriendRequest(ctx context.Context, friendRequest *relation.FriendRequest) (err error)
	// AgreeFriendRequest 同意好友申请
	AgreeFriendRequest(ctx context.Context, friendRequest *relation.FriendRequest) (err error)
	// Delete 删除好友
	Delete(ctx context.Context, ownerUserID string, friendUserIDs string) (err error)
	// UpdateRemark 更新好友备注
	UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) (err error)
	// FindOwnerFriends 获取ownerUserID的好友列表
	FindOwnerFriends(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (friends []*relation.Friend, err error)
	// FindInWhoseFriends friendUserID在哪些人的好友列表中
	FindInWhoseFriends(ctx context.Context, friendUserID string, pageNumber, showNumber int32) (friends []*relation.Friend, err error)
	// FindFriendRequestFromMe 获取我发出去的好友申请
	FindFriendRequestFromMe(ctx context.Context, userID string, pageNumber, showNumber int32) (friends []*relation.FriendRequest, err error)
	// FindFriendRequestToMe 获取我收到的的好友申请
	FindFriendRequestToMe(ctx context.Context, userID string, pageNumber, showNumber int32) (friends []*relation.FriendRequest, err error)
	// FindFriends 获取某人指定好友的信息
	FindFriends(ctx context.Context, ownerUserID string, friendUserIDs []string) (friends []*relation.Friend, err error)
}

type FriendDatabase struct {
	sqlDB *relation.Friend
}

func NewFriendDatabase(db *gorm.DB) *FriendDatabase {
	sqlDB := relation.NewFriendDB(db)
	database := &FriendDatabase{
		sqlDB: sqlDB,
	}
	return database
}

// CheckIn 检查user2是否在user1的好友列表中(inUser1Friends==true) 检查user1是否在user2的好友列表中(inUser2Friends==true)
func (f *FriendDatabase) CheckIn(ctx context.Context, user1, user2 string) (err error, inUser1Friends bool, inUser2Friends bool) {
}

// AddFriendRequest 增加或者更新好友申请
func (f *FriendDatabase) AddFriendRequest(ctx context.Context, fromUserID, toUserID string, reqMsg string, ex string) (err error) {
}

// BecomeFriend 先判断是否在好友表，如果在则不插入
func (f *FriendDatabase) BecomeFriend(ctx context.Context, friends []*relation.Friend) (err error) {
}

// RefuseFriendRequest 拒绝好友申请
func (f *FriendDatabase) RefuseFriendRequest(ctx context.Context, friendRequest *relation.FriendRequest) (err error) {
}

// AgreeFriendRequest 同意好友申请
func (f *FriendDatabase) AgreeFriendRequest(ctx context.Context, friendRequest *relation.FriendRequest) (err error) {
}

// Delete 删除好友
func (f *FriendDatabase) Delete(ctx context.Context, ownerUserID string, friendUserIDs string) (err error) {
}

// UpdateRemark 更新好友备注
func (f *FriendDatabase) UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) (err error) {
}

// FindOwnerFriends 获取ownerUserID的好友列表
func (f *FriendDatabase) FindOwnerFriends(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (friends []*relation.Friend, err error) {
}

// FindInWhoseFriends friendUserID在哪些人的好友列表中
func (f *FriendDatabase) FindInWhoseFriends(ctx context.Context, friendUserID string, pageNumber, showNumber int32) (friends []*relation.Friend, err error) {
}

// FindFriendRequestFromMe 获取我发出去的好友申请
func (f *FriendDatabase) FindFriendRequestFromMe(ctx context.Context, userID string, pageNumber, showNumber int32) (friends []*relation.FriendRequest, err error) {
}

// FindFriendRequestToMe 获取我收到的的好友申请
func (f *FriendDatabase) FindFriendRequestToMe(ctx context.Context, userID string, pageNumber, showNumber int32) (friends []*relation.FriendRequest, err error) {
}

// FindFriends 获取某人指定好友的信息
func (f *FriendDatabase) FindFriends(ctx context.Context, ownerUserID string, friendUserIDs []string) (friends []*relation.Friend, err error) {
}
