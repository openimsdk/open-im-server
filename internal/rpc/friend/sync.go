package friend

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/internal/rpc/incrversion"
	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/tools/errs"
)

func (s *friendServer) NotificationUserInfoUpdate(ctx context.Context, req *relation.NotificationUserInfoUpdateReq) (*relation.NotificationUserInfoUpdateResp, error) {
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
	return &relation.NotificationUserInfoUpdateResp{}, nil
}

func (s *friendServer) SearchFriends(ctx context.Context, req *relation.SearchFriendsReq) (*relation.SearchFriendsResp, error) {
	if err := s.userRpcClient.Access(ctx, req.UserID); err != nil {
		return nil, err
	}
	if req.Keyword == "" {
		total, friends, err := s.friendDatabase.PageOwnerFriends(ctx, req.UserID, req.Pagination)
		if err != nil {
			return nil, err
		}
		return &relation.SearchFriendsResp{
			Total:   total,
			Friends: friendsDB2PB(friends),
		}, nil
	}
	total, friends, err := s.friendDatabase.SearchFriend(ctx, req.UserID, req.Keyword, req.Pagination)
	if err != nil {
		return nil, err
	}
	return &relation.SearchFriendsResp{
		Total:   total,
		Friends: friendsDB2PB(friends),
	}, nil
}

func (s *friendServer) GetIncrementalFriends(ctx context.Context, req *relation.GetIncrementalFriendsReq) (*relation.GetIncrementalFriendsResp, error) {
	if err := authverify.CheckAccessV3(ctx, req.UserID, s.config.Share.IMAdminUserID); err != nil {
		return nil, err
	}
	opt := incrversion.Option[*relation.FriendInfo, relation.GetIncrementalFriendsResp]{
		VersionID: req.VersionID,
		Version: func() (*model.VersionLog, error) {
			return s.friendDatabase.FindFriendIncrVersion(ctx, req.UserID, uint(req.Version), incrversion.Limit(s.config.RpcConfig.FriendSyncCount, req.Version))
		},
		AllID: func() ([]string, error) {
			return s.friendDatabase.FindSortFriendUserIDs(ctx, req.UserID)
		},
		Find: func(ids []string) ([]*relation.FriendInfo, error) {
			friends, err := s.friendDatabase.FindFriendsWithError(ctx, req.UserID, ids)
			if err != nil {
				return nil, err
			}
			return friendsDB2PB(friends), nil
		},
		ID: func(elem *relation.FriendInfo) string {
			return elem.FriendUserID
		},
		Resp: func(version *model.VersionLog, delIDs []string, list []*relation.FriendInfo, full bool) *relation.GetIncrementalFriendsResp {
			return &relation.GetIncrementalFriendsResp{
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

//func (s *friendServer) GetIncrementalFriends(ctx context.Context, req *relation.GetIncrementalFriendsReq) (*relation.GetIncrementalFriendsResp, error) {
//	if err := authverify.CheckAccessV3(ctx, req.UserID, s.config.Share.IMAdminUserID); err != nil {
//		return nil, err
//	}
//	var limit int
//	if req.Version > 0 {
//		limit = s.config.RpcConfig.FriendSyncCount
//	}
//	incrVer, err := s.friendDatabase.FindFriendIncrVersion(ctx, req.UserID, uint(req.Version), limit)
//	if err != nil {
//		return nil, err
//	}
//	var (
//		deleteUserIDs []string
//		changeUserIDs []string
//	)
//	if incrVer.Full() {
//		changeUserIDs, err = s.friendDatabase.FindSortFriendUserIDs(ctx, req.UserID)
//		if err != nil {
//			return nil, err
//		}
//	} else {
//		deleteUserIDs, changeUserIDs = incrVer.DeleteAndChangeIDs()
//	}
//	var friends []*model.Friend
//	if len(changeUserIDs) > 0 {
//		friends, err = s.friendDatabase.FindFriendsWithError(ctx, req.UserID, changeUserIDs)
//		if err != nil {
//			return nil, err
//		}
//	}
//	return &relation.GetIncrementalFriendsResp{
//		Version:       uint64(incrVer.Version),
//		VersionID:     incrVer.ID.Hex(),
//		Full:          incrVer.Full(),
//		SyncCount:     uint32(s.config.RpcConfig.FriendSyncCount),
//		DeleteUserIds: deleteUserIDs,
//		Changes:       friendsDB2PB(friends),
//	}, nil
//}
