package controller

import (
	"Open_IM/pkg/common/db/relation"
	"context"
	"gorm.io/gorm"
)

type FriendInterface interface {
	// CheckIn 检查fromUserID是否在toUserID的好友列表中(inTo==true) 检查toUserID是否在fromUserID的好友列表中(inFrom==true)
	CheckIn(ctx context.Context, fromUserID, toUserID string) (err error, inTo bool, inFrom bool)
	// AddFriendRequest 增加或者更新好友申请
	AddFriendRequest(ctx context.Context, fromUserID, toUserID string, reqMsg string, ex string) (err error)
	// 先判断是否在好友表，如果在则不插入
	BecomeFriend(ctx context.Context, friends []*relation.Friend) (err error)
	//拒绝好友申请
	RefuseFriendRequest(ctx context.Context, friendRequest *relation.FriendRequest) (err error)
	// 同意好友申请
	AgreeFriendRequest(ctx context.Context, friendRequest *relation.FriendRequest) (err error)
	Create(ctx context.Context, friends []*relation.Friend) (err error)
	Delete(ctx context.Context, ownerUserID string, friendUserIDs string) (err error)
	UpdateByMap(ctx context.Context, ownerUserID string, args map[string]interface{}) (err error)
	Update(ctx context.Context, friends []*relation.Friend) (err error)
	UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) (err error)
	//获取ownerUserID的好友列表
	FindOwnerFriends(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (friends []*relation.Friend, err error)
	//friendUserID在哪些人的好友列表中
	FindInWhoseFriends(ctx context.Context, friendUserID string, pageNumber, showNumber int32) (friends []*relation.Friend, err error)
	Take(ctx context.Context, ownerUserID, friendUserID string) (friend *relation.Friend, err error)
	FindUserState(ctx context.Context, userID1, userID2 string) (friends []*relation.Friend, err error)
}

type FriendController struct {
	database FriendDatabaseInterface
}

func NewFriendController(db *gorm.DB) *FriendController {
	return &FriendController{database: NewFriendDatabase(db)}
}

func (f *FriendController) Create(ctx context.Context, friends []*relation.Friend) (err error) {
	return f.database.Create(ctx, friends)
}
func (f *FriendController) Delete(ctx context.Context, ownerUserID string, friendUserIDs string) (err error) {
	return f.database.Delete(ctx, ownerUserID, friendUserIDs)
}
func (f *FriendController) UpdateByMap(ctx context.Context, ownerUserID string, args map[string]interface{}) (err error) {
	return f.database.UpdateByMap(ctx, ownerUserID, args)
}
func (f *FriendController) Update(ctx context.Context, friends []*relation.Friend) (err error) {
	return f.database.Update(ctx, friends)
}
func (f *FriendController) UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) (err error) {
	return f.database.UpdateRemark(ctx, ownerUserID, friendUserID, remark)
}
func (f *FriendController) FindOwnerUserID(ctx context.Context, ownerUserID string) (friends []*relation.Friend, err error) {
	return f.database.FindOwnerUserID(ctx, ownerUserID)
}
func (f *FriendController) FindFriendUserID(ctx context.Context, friendUserID string) (friends []*relation.Friend, err error) {
	return f.database.FindFriendUserID(ctx, friendUserID)
}
func (f *FriendController) Take(ctx context.Context, ownerUserID, friendUserID string) (friend *relation.Friend, err error) {
	return f.database.Take(ctx, ownerUserID, friendUserID)
}
func (f *FriendController) FindUserState(ctx context.Context, userID1, userID2 string) (friends []*relation.Friend, err error) {
	return f.database.FindUserState(ctx, userID1, userID2)
}

type FriendDatabaseInterface interface {
	Create(ctx context.Context, friends []*relation.Friend) (err error)
	Delete(ctx context.Context, ownerUserID string, friendUserIDs string) (err error)
	UpdateByMap(ctx context.Context, ownerUserID string, args map[string]interface{}) (err error)
	Update(ctx context.Context, friends []*relation.Friend) (err error)
	UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) (err error)

	FindOwnerUserID(ctx context.Context, ownerUserID string) (friends []*relation.Friend, err error)
	FindFriendUserID(ctx context.Context, friendUserID string) (friends []*relation.Friend, err error)
	Take(ctx context.Context, ownerUserID, friendUserID string) (friend *relation.Friend, err error)
	FindUserState(ctx context.Context, userID1, userID2 string) (friends []*relation.Friend, err error)
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

func (f *FriendDatabase) Create(ctx context.Context, friends []*relation.Friend) (err error) {
	return f.sqlDB.Create(ctx, friends)
}
func (f *FriendDatabase) Delete(ctx context.Context, ownerUserID string, friendUserIDs string) (err error) {
	return f.sqlDB.Delete(ctx, ownerUserID, friendUserIDs)
}
func (f *FriendDatabase) UpdateByMap(ctx context.Context, ownerUserID string, args map[string]interface{}) (err error) {
	return f.sqlDB.UpdateByMap(ctx, ownerUserID, args)
}
func (f *FriendDatabase) Update(ctx context.Context, friends []*relation.Friend) (err error) {
	return f.sqlDB.Update(ctx, friends)
}
func (f *FriendDatabase) UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) (err error) {
	return f.sqlDB.UpdateRemark(ctx, ownerUserID, friendUserID, remark)
}
func (f *FriendDatabase) FindOwnerUserID(ctx context.Context, ownerUserID string) (friends []*relation.Friend, err error) {
	return f.sqlDB.FindOwnerUserID(ctx, ownerUserID)
}
func (f *FriendDatabase) FindFriendUserID(ctx context.Context, friendUserID string) (friends []*relation.Friend, err error) {
	return f.sqlDB.FindFriendUserID(ctx, friendUserID)
}
func (f *FriendDatabase) Take(ctx context.Context, ownerUserID, friendUserID string) (friend *relation.Friend, err error) {
	return f.sqlDB.Take(ctx, ownerUserID, friendUserID)
}
func (f *FriendDatabase) FindUserState(ctx context.Context, userID1, userID2 string) (friends []*relation.Friend, err error) {
	return f.sqlDB.FindUserState(ctx, userID1, userID2)
}
