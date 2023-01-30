package controller

import (
	"Open_IM/pkg/common/db/relation"
	"context"
	"gorm.io/gorm"
)

type FriendRequestInterface interface {
	Create(ctx context.Context, friends []*relation.FriendRequest) (err error)
	Delete(ctx context.Context, fromUserID, toUserID string) (err error)
	UpdateByMap(ctx context.Context, ownerUserID string, args map[string]interface{}) (err error)
	Update(ctx context.Context, friends []*relation.FriendRequest) (err error)
	Find(ctx context.Context, ownerUserID string) (friends []*relation.FriendRequest, err error)
	Take(ctx context.Context, fromUserID, toUserID string) (friend *relation.FriendRequest, err error)
	FindToUserID(ctx context.Context, toUserID string) (friends []*relation.FriendRequest, err error)
	FindFromUserID(ctx context.Context, fromUserID string) (friends []*relation.FriendRequest, err error)
}

type FriendRequestController struct {
	database FriendRequestInterface
}

func NewFriendRequestController(db *gorm.DB) *FriendRequestController {
	return &FriendRequestController{database: NewFriendRequestDatabase(db)}
}

func (f *FriendRequestController) Create(ctx context.Context, friends []*relation.FriendRequest) (err error) {
	return f.database.Create(ctx, friends)
}
func (f *FriendRequestController) Delete(ctx context.Context, fromUserID, toUserID string) (err error) {
	return f.database.Delete(ctx, fromUserID, toUserID)
}
func (f *FriendRequestController) UpdateByMap(ctx context.Context, ownerUserID string, args map[string]interface{}) (err error) {
	return f.database.UpdateByMap(ctx, ownerUserID, args)
}
func (f *FriendRequestController) Update(ctx context.Context, friends []*relation.FriendRequest) (err error) {
	return f.database.Update(ctx, friends)
}
func (f *FriendRequestController) Find(ctx context.Context, ownerUserID string) (friends []*relation.FriendRequest, err error) {
	return f.database.Find(ctx, ownerUserID)
}
func (f *FriendRequestController) Take(ctx context.Context, fromUserID, toUserID string) (friend *relation.FriendRequest, err error) {
	return f.database.Take(ctx, fromUserID, toUserID)
}
func (f *FriendRequestController) FindToUserID(ctx context.Context, toUserID string) (friends []*relation.FriendRequest, err error) {
	return f.database.FindToUserID(ctx, toUserID)
}
func (f *FriendRequestController) FindFromUserID(ctx context.Context, fromUserID string) (friends []*relation.FriendRequest, err error) {
	return f.database.FindFromUserID(ctx, fromUserID)
}

type FriendRequestDatabaseInterface interface {
	Create(ctx context.Context, friends []*relation.FriendRequest) (err error)
	Delete(ctx context.Context, fromUserID, toUserID string) (err error)
	UpdateByMap(ctx context.Context, ownerUserID string, args map[string]interface{}) (err error)
	Update(ctx context.Context, friends []*relation.FriendRequest) (err error)
	Find(ctx context.Context, ownerUserID string) (friends []*relation.FriendRequest, err error)
	Take(ctx context.Context, fromUserID, toUserID string) (friend *relation.FriendRequest, err error)
	FindToUserID(ctx context.Context, toUserID string) (friends []*relation.FriendRequest, err error)
	FindFromUserID(ctx context.Context, fromUserID string) (friends []*relation.FriendRequest, err error)
}

type FriendRequestDatabase struct {
	sqlDB  *relation.FriendRequest
	friend *FriendDatabase
}

func (f *FriendRequestDatabase) Update(ctx context.Context, friends []*relation.FriendRequest) (err error) {
	return f.sqlDB.DB.Transaction(func(tx *gorm.DB) error {
		if err := f.sqlDB.Update(ctx, friends); err != nil {
			return err
		}
		if err := f.friend.Update(); err != nil {
			return err
		}
		return nil
	})
}

func NewFriendRequestDatabase(db *gorm.DB) *FriendRequestDatabase {
	sqlDB := relation.NewFriendRequest(db)
	database := &FriendRequestDatabase{
		sqlDB: sqlDB,
	}
	return database
}

func (f *FriendRequestDatabase) Create(ctx context.Context, friends []*relation.FriendRequest) (err error) {
	return f.sqlDB.Create(ctx, friends)
}
func (f *FriendRequestDatabase) Delete(ctx context.Context, fromUserID, toUserID string) (err error) {
	return f.sqlDB.Delete(ctx, fromUserID, toUserID)
}
func (f *FriendRequestDatabase) UpdateByMap(ctx context.Context, ownerUserID string, args map[string]interface{}) (err error) {
	return f.sqlDB.UpdateByMap(ctx, ownerUserID, args)
}

func (f *FriendRequestDatabase) Find(ctx context.Context, ownerUserID string) (friends []*relation.FriendRequest, err error) {
	return f.sqlDB.Find(ctx, ownerUserID)
}
func (f *FriendRequestDatabase) Take(ctx context.Context, fromUserID, toUserID string) (friend *relation.FriendRequest, err error) {
	return f.sqlDB.Take(ctx, fromUserID, toUserID)
}
func (f *FriendRequestDatabase) FindToUserID(ctx context.Context, toUserID string) (friends []*relation.FriendRequest, err error) {
	return f.sqlDB.FindToUserID(ctx, toUserID)
}
func (f *FriendRequestDatabase) FindFromUserID(ctx context.Context, fromUserID string) (friends []*relation.FriendRequest, err error) {
	return f.sqlDB.FindFromUserID(ctx, fromUserID)
}
