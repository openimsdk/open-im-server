package controller

import (
	"Open_IM/pkg/common/db/cache"
	"Open_IM/pkg/common/db/relation"
	"context"
	"errors"
	"gorm.io/gorm"
)

type FriendModel struct {
	db    *relation.Friend
	cache *cache.GroupCache
}

func (f *FriendModel) Create(ctx context.Context, friends []*relation.Friend) (err error) {
	return f.db.Create(ctx, friends)
}

func (f *FriendModel) Delete(ctx context.Context, ownerUserID string, friendUserIDs string) (err error) {
	return f.db.Delete(ctx, ownerUserID, friendUserIDs)
}

func (f *FriendModel) UpdateByMap(ctx context.Context, ownerUserID string, args map[string]interface{}) (err error) {
	return f.db.UpdateByMap(ctx, ownerUserID, args)
}

func (f *FriendModel) Update(ctx context.Context, friends []*relation.Friend) (err error) {
	return f.db.Update(ctx, friends)
}

func (f *FriendModel) UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) (err error) {
	return f.db.UpdateRemark(ctx, ownerUserID, friendUserID, remark)
}

func (f *FriendModel) FindOwnerUserID(ctx context.Context, ownerUserID string) (friends []*relation.Friend, err error) {
	return f.db.FindOwnerUserID(ctx, ownerUserID)
}

func (f *FriendModel) FindFriendUserID(ctx context.Context, friendUserID string) (friends []*relation.Friend, err error) {
	return f.db.FindFriendUserID(ctx, friendUserID)
}

func (f *FriendModel) Take(ctx context.Context, ownerUserID, friendUserID string) (friend *relation.Friend, err error) {
	return f.db.Take(ctx, ownerUserID, friendUserID)
}

func (f *FriendModel) FindUserState(ctx context.Context, userID1, userID2 string) (friends []*relation.Friend, err error) {
	return f.db.FindUserState(ctx, userID1, userID2)
}

func (f *FriendModel) IsExist(ctx context.Context, ownerUserID, friendUserID string) (bool, error) {
	if _, err := f.Take(ctx, ownerUserID, friendUserID); err == nil {
		return true, nil
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else {
		return false, err
	}
}
