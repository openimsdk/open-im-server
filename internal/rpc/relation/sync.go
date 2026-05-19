package relation

import (
	"context"
	"slices"

	"github.com/openimsdk/open-im-server/v3/pkg/util/hashutil"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"

	"github.com/openimsdk/open-im-server/v3/internal/rpc/incrversion"
	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/protocol/relation"
)

func (s *friendServer) NotificationUserInfoUpdate(ctx context.Context, req *relation.NotificationUserInfoUpdateReq) (*relation.NotificationUserInfoUpdateResp, error) {
	// 单向好友：仅通知「好友列表中包含 req.UserID」的用户（owner -> req.UserID），
	// 而非 req.UserID 自己好友列表中的对端。
	ownerUserIDs, err := s.db.FindFriendUserID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if len(ownerUserIDs) > 0 {
		friendUserIDs := []string{req.UserID}
		noCancelCtx := context.WithoutCancel(ctx)
		err := s.queue.PushCtx(ctx, func() {
			for _, ownerUserID := range ownerUserIDs {
				if err := s.db.OwnerIncrVersion(noCancelCtx, ownerUserID, friendUserIDs, model.VersionStateUpdate); err != nil {
					log.ZError(ctx, "OwnerIncrVersion", err, "userID", ownerUserID, "friendUserIDs", friendUserIDs)
				}
			}
			for _, ownerUserID := range ownerUserIDs {
				s.notificationSender.FriendInfoUpdatedNotification(noCancelCtx, req.UserID, ownerUserID)
			}
		})
		if err != nil {
			log.ZError(ctx, "NotificationUserInfoUpdate timeout", err, "userID", req.UserID)
		}
	}
	return &relation.NotificationUserInfoUpdateResp{}, nil
}

func (s *friendServer) GetFullFriendUserIDs(ctx context.Context, req *relation.GetFullFriendUserIDsReq) (*relation.GetFullFriendUserIDsResp, error) {
	if err := authverify.CheckAccessV3(ctx, req.UserID, s.config.Share.IMAdminUserID); err != nil {
		return nil, err
	}
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
		Version:   uint64(vl.Version),
		VersionID: vl.ID.Hex(),
		Equal:     req.IdHash == idHash,
		UserIDs:   userIDs,
	}, nil
}

func (s *friendServer) GetIncrementalFriends(ctx context.Context, req *relation.GetIncrementalFriendsReq) (*relation.GetIncrementalFriendsResp, error) {
	if err := authverify.CheckAccessV3(ctx, req.UserID, s.config.Share.IMAdminUserID); err != nil {
		return nil, err
	}
	var sortVersion uint64
	opt := incrversion.Option[*sdkws.FriendInfo, relation.GetIncrementalFriendsResp]{
		Ctx:           ctx,
		VersionKey:    req.UserID,
		VersionID:     req.VersionID,
		VersionNumber: req.Version,
		Version: func(ctx context.Context, ownerUserID string, version uint, limit int) (*model.VersionLog, error) {
			vl, err := s.db.FindFriendIncrVersion(ctx, ownerUserID, version, limit)
			if err != nil {
				return nil, err
			}
			vl.Logs = slices.DeleteFunc(vl.Logs, func(elem model.VersionLogElem) bool {
				if elem.EID == model.VersionSortChangeID {
					vl.LogLen--
					sortVersion = uint64(elem.Version)
					return true
				}
				return false
			})
			return vl, nil
		},
		CacheMaxVersion: s.db.FindMaxFriendVersionCache,
		Find: func(ctx context.Context, ids []string) ([]*sdkws.FriendInfo, error) {
			return s.getFriend(ctx, req.UserID, ids)
		},
		Resp: func(version *model.VersionLog, deleteIds []string, insertList, updateList []*sdkws.FriendInfo, full bool) *relation.GetIncrementalFriendsResp {
			return &relation.GetIncrementalFriendsResp{
				VersionID:   version.ID.Hex(),
				Version:     uint64(version.Version),
				Full:        full,
				Delete:      deleteIds,
				Insert:      insertList,
				Update:      updateList,
				SortVersion: sortVersion,
			}
		},
	}
	return opt.Build()
}
