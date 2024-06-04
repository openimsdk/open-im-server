package friend

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/internal/rpc/incrversion"
	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	pbfriend "github.com/openimsdk/protocol/friend"
	"github.com/openimsdk/tools/errs"
)

func (s *friendServer) NotificationUserInfoUpdate(ctx context.Context, req *pbfriend.NotificationUserInfoUpdateReq) (*pbfriend.NotificationUserInfoUpdateResp, error) {
	if req.NewUserInfo == nil {
		var err error
		req.NewUserInfo, err = s.userRpcClient.GetUserInfo(ctx, req.UserID)
		if err != nil {
			return nil, err
		}
	}
	if req.UserID != req.NewUserInfo.UserID {
		return nil, errs.ErrArgs.WrapMsg("req.UserID != req.NewUserInfo.UserID")
	}
	userIDs, err := s.friendDatabase.FindFriendUserID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if len(userIDs) > 0 {
		if err := s.friendDatabase.UpdateFriendUserInfo(ctx, req.UserID, userIDs, req.NewUserInfo.Nickname, req.NewUserInfo.FaceURL); err != nil {
			return nil, err
		}
		s.notificationSender.FriendsInfoUpdateNotification(ctx, req.UserID, userIDs)
	}
	return &pbfriend.NotificationUserInfoUpdateResp{}, nil
}

func (s *friendServer) SearchFriends(ctx context.Context, req *pbfriend.SearchFriendsReq) (*pbfriend.SearchFriendsResp, error) {
	if err := s.userRpcClient.Access(ctx, req.UserID); err != nil {
		return nil, err
	}
	if req.Keyword == "" {
		total, friends, err := s.friendDatabase.PageOwnerFriends(ctx, req.UserID, req.Pagination)
		if err != nil {
			return nil, err
		}
		return &pbfriend.SearchFriendsResp{
			Total:   total,
			Friends: friendsDB2PB(friends),
		}, nil
	}
	total, friends, err := s.friendDatabase.SearchFriend(ctx, req.UserID, req.Keyword, req.Pagination)
	if err != nil {
		return nil, err
	}
	return &pbfriend.SearchFriendsResp{
		Total:   total,
		Friends: friendsDB2PB(friends),
	}, nil
}

func (s *friendServer) GetIncrementalFriends(ctx context.Context, req *pbfriend.GetIncrementalFriendsReq) (*pbfriend.GetIncrementalFriendsResp, error) {
	if err := authverify.CheckAccessV3(ctx, req.UserID, s.config.Share.IMAdminUserID); err != nil {
		return nil, err
	}
	opt := incrversion.Option[*pbfriend.FriendInfo, pbfriend.GetIncrementalFriendsResp]{
		Ctx:             ctx,
		VersionKey:      req.UserID,
		VersionID:       req.VersionID,
		VersionNumber:   req.Version,
		SyncLimit:       s.config.RpcConfig.FriendSyncCount,
		Version:         s.friendDatabase.FindFriendIncrVersion,
		CacheMaxVersion: s.friendDatabase.FindMaxFriendVersionCache,
		SortID:          s.friendDatabase.FindSortFriendUserIDs,
		Find: func(ctx context.Context, ids []string) ([]*pbfriend.FriendInfo, error) {
			friends, err := s.friendDatabase.FindFriendsWithError(ctx, req.UserID, ids)
			if err != nil {
				return nil, err
			}
			return friendsDB2PB(friends), nil
		},
		ID: func(elem *pbfriend.FriendInfo) string { return elem.FriendUserID },
		Resp: func(version *model.VersionLog, delIDs []string, list []*pbfriend.FriendInfo, full bool) *pbfriend.GetIncrementalFriendsResp {
			return &pbfriend.GetIncrementalFriendsResp{
				VersionID:     version.ID.Hex(),
				Version:       uint64(version.Version),
				Full:          full,
				SyncCount:     uint32(s.config.RpcConfig.FriendSyncCount),
				DeleteUserIds: delIDs,
				Changes:       list,
			}
		},
	}
	return opt.Build()
}
