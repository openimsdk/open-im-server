package friend

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/util/hashutil"
	"github.com/openimsdk/protocol/sdkws"

	"github.com/openimsdk/open-im-server/v3/internal/rpc/incrversion"
	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/protocol/relation"
)

func (s *friendServer) NotificationUserInfoUpdate(ctx context.Context, req *relation.NotificationUserInfoUpdateReq) (*relation.NotificationUserInfoUpdateResp, error) {
	userIDs, err := s.db.FindFriendUserIDs(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	for _, userID := range userIDs {
		if err := s.db.OwnerIncrVersion(ctx, userID, []string{req.UserID}, model.VersionStateUpdate); err != nil {
			return nil, err
		}
	}
	for _, userID := range userIDs {
		s.notificationSender.FriendInfoUpdatedNotification(ctx, req.UserID, userID)
	}
	return &relation.NotificationUserInfoUpdateResp{}, nil
}

func (s *friendServer) GetFullFriendUserIDs(ctx context.Context, req *relation.GetFullFriendUserIDsReq) (*relation.GetFullFriendUserIDsResp, error) {
	vl, err := s.db.FindMaxFriendVersionCache(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	userIDs, err := s.db.FindFriendUserIDs(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	idHash := hashutil.IdHash(userIDs)
	if req.IdHash == idHash {
		userIDs = nil
	}
	return &relation.GetFullFriendUserIDsResp{
		Version:   idHash,
		VersionID: vl.ID.Hex(),
		Equal:     req.IdHash == idHash,
		UserIDs:   userIDs,
	}, nil
}

func (s *friendServer) GetIncrementalFriends(ctx context.Context, req *relation.GetIncrementalFriendsReq) (*relation.GetIncrementalFriendsResp, error) {
	if err := authverify.CheckAccessV3(ctx, req.UserID, s.config.Share.IMAdminUserID); err != nil {
		return nil, err
	}
	opt := incrversion.Option[*sdkws.FriendInfo, relation.GetIncrementalFriendsResp]{
		Ctx:             ctx,
		VersionKey:      req.UserID,
		VersionID:       req.VersionID,
		VersionNumber:   req.Version,
		Version:         s.db.FindFriendIncrVersion,
		CacheMaxVersion: s.db.FindMaxFriendVersionCache,
		Find: func(ctx context.Context, ids []string) ([]*sdkws.FriendInfo, error) {
			return s.getFriend(ctx, req.UserID, ids)
		},
		ID: func(elem *sdkws.FriendInfo) string { return elem.FriendUser.UserID },
		Resp: func(version *model.VersionLog, deleteIds []string, insertList, updateList []*sdkws.FriendInfo, full bool) *relation.GetIncrementalFriendsResp {
			return &relation.GetIncrementalFriendsResp{
				VersionID: version.ID.Hex(),
				Version:   uint64(version.Version),
				Full:      full,
				Delete:    deleteIds,
				Insert:    insertList,
				Update:    updateList,
			}
		},
	}
	return opt.Build()
}
