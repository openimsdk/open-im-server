package database

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/pagination"
)

type FriendRequest interface {
	// Insert multiple records
	Create(ctx context.Context, friendRequests []*model.FriendRequest) (err error)
	// Delete record
	Delete(ctx context.Context, fromUserID, toUserID string) (err error)
	// Update with zero values
	UpdateByMap(ctx context.Context, formUserID string, toUserID string, args map[string]any) (err error)
	// Update multiple records (non-zero values)
	Update(ctx context.Context, friendRequest *model.FriendRequest) (err error)
	// Get friend requests sent to a specific user, no error returned if not found
	Find(ctx context.Context, fromUserID, toUserID string) (friendRequest *model.FriendRequest, err error)
	Take(ctx context.Context, fromUserID, toUserID string) (friendRequest *model.FriendRequest, err error)
	// Get list of friend requests received by toUserID
	FindToUserID(ctx context.Context, toUserID string, handleResults []int, pagination pagination.Pagination) (total int64, friendRequests []*model.FriendRequest, err error)
	// Get list of friend requests sent by fromUserID
	FindFromUserID(ctx context.Context, fromUserID string, handleResults []int, pagination pagination.Pagination) (total int64, friendRequests []*model.FriendRequest, err error)
	FindBothFriendRequests(ctx context.Context, fromUserID, toUserID string) (friends []*model.FriendRequest, err error)
	GetUnhandledCount(ctx context.Context, userID string, ts int64) (int64, error)
}
