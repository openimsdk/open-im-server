package cache

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

type UserCache interface {
	BatchDeleter
	CloneUserCache() UserCache
	GetUserInfo(ctx context.Context, userID string) (userInfo *model.User, err error)
	GetUsersInfo(ctx context.Context, userIDs []string) ([]*model.User, error)
	DelUsersInfo(userIDs ...string) UserCache
	GetUserGlobalRecvMsgOpt(ctx context.Context, userID string) (opt int, err error)
	DelUsersGlobalRecvMsgOpt(userIDs ...string) UserCache
	//GetUserStatus(ctx context.Context, userIDs []string) ([]*user.OnlineStatus, error)
	//SetUserStatus(ctx context.Context, userID string, status, platformID int32) error
}
