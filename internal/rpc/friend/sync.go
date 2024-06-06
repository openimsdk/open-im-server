package friend

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/internal/rpc/incrversion"
	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/tools/errs"
)

//func (s *friendServer) SearchFriends(ctx context.Context, req *pbfriend.SearchFriendsReq) (*pbfriend.SearchFriendsResp, error) {
//	if err := s.userRpcClient.Access(ctx, req.UserID); err != nil {
//		return nil, err
//	}
//	if req.Keyword == "" {
//		total, friends, err := s.friendDatabase.PageOwnerFriends(ctx, req.UserID, req.Pagination)
//		if err != nil {
//			return nil, err
//		}
//		return &pbfriend.SearchFriendsResp{
//			Total:   total,
//			Friends: friendsDB2PB(friends),
//		}, nil
//	}
//	total, friends, err := s.friendDatabase.SearchFriend(ctx, req.UserID, req.Keyword, req.Pagination)
//	if err != nil {
//		return nil, err
//	}
//	return &pbfriend.SearchFriendsResp{
//		Total:   total,
//		Friends: friendsDB2PB(friends),
//	}, nil
//}

func (s *friendServer) GetIncrementalFriends(ctx context.Context, req *relation.GetIncrementalFriendsReq) (*relation.GetIncrementalFriendsResp, error) {
	if err := authverify.CheckAccessV3(ctx, req.UserID, s.config.Share.IMAdminUserID); err != nil {
		return nil, err
	}
	opt := incrversion.Option[*sdkws.FriendInfo, pbfriend.GetIncrementalFriendsResp]{
		Ctx:             ctx,
		VersionKey:      req.UserID,
		VersionID:       req.VersionID,
		VersionNumber:   req.Version,
		SyncLimit:       s.config.RpcConfig.FriendSyncCount,
		Version:         s.friendDatabase.FindFriendIncrVersion,
		CacheMaxVersion: s.friendDatabase.FindMaxFriendVersionCache,
		SortID:          s.friendDatabase.FindSortFriendUserIDs,
		Find: func(ctx context.Context, ids []string) ([]*sdkws.FriendInfo, error) {
			return s.getFriend(ctx, req.UserID, ids)
		},
		ID: func(elem *sdkws.FriendInfo) string { return elem.FriendUser.UserID },
		Resp: func(version *model.VersionLog, delIDs []string, list []*sdkws.FriendInfo, full bool) *pbfriend.GetIncrementalFriendsResp {
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
