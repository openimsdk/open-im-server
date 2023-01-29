package controller

import (
	"Open_IM/pkg/common/db/cache"
	"Open_IM/pkg/common/db/relation"
	"context"
)

type FriendRequestModel struct {
	db    *relation.FriendRequest
	cache *cache.GroupCache
}

func (f *FriendRequestModel) Create(ctx context.Context, friends []*relation.FriendRequest) (err error) {
	return f.db.Create(ctx, friends)
}

func (f *FriendRequestModel) Delete(ctx context.Context, fromUserID, toUserID string) (err error) {
	return f.db.Delete(ctx, fromUserID, toUserID)
}

func (f *FriendRequestModel) UpdateByMap(ctx context.Context, ownerUserID string, args map[string]interface{}) (err error) {
	return f.db.UpdateByMap(ctx, ownerUserID, args)
}

func (f *FriendRequestModel) Update(ctx context.Context, friends []*relation.FriendRequest) (err error) {
	return f.db.Update(ctx, friends)
}

func (f *FriendRequestModel) Find(ctx context.Context, ownerUserID string) (friends []*relation.FriendRequest, err error) {
	return f.db.Find(ctx, ownerUserID)
}

func (f *FriendRequestModel) Take(ctx context.Context, fromUserID, toUserID string) (friend *relation.FriendRequest, err error) {
	return f.db.Take(ctx, fromUserID, toUserID)
}

func (f *FriendRequestModel) FindToUserID(ctx context.Context, toUserID string) (friends []*relation.FriendRequest, err error) {
	return f.db.FindToUserID(ctx, toUserID)
}

func (f *FriendRequestModel) FindFromUserID(ctx context.Context, fromUserID string) (friends []*relation.FriendRequest, err error) {
	return f.db.FindFromUserID(ctx, fromUserID)
}
