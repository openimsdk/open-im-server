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
	userIDs, err := s.db.FindFriendUserIDs(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if len(userIDs) > 0 {
		friendUserIDs := []string{req.UserID}
		noCancelCtx := context.WithoutCancel(ctx)
		err := s.queue.PushCtx(ctx, func() {
			for _, userID := range userIDs {
				if err := s.db.OwnerIncrVersion(noCancelCtx, userID, friendUserIDs, model.VersionStateUpdate); err != nil {
					log.ZError(ctx, "OwnerIncrVersion", err, "userID", userID, "friendUserIDs", friendUserIDs)
				}
			}
			for _, userID := range userIDs {
				s.notificationSender.FriendInfoUpdatedNotification(noCancelCtx, req.UserID, userID)
			}
		})
		if err != nil {
			log.ZError(ctx, "NotificationUserInfoUpdate timeout", err, "userID", req.UserID)
		}
	}
	return &relation.NotificationUserInfoUpdateResp{}, nil
}

func (s *friendServer) GetFullFriendUserIDs(ctx context.Context, req *relation.GetFullFriendUserIDsReq) (*relation.GetFullFriendUserIDsResp, error) {
	if err := authverify.CheckAccess(ctx, req.UserID); err != nil {
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
		Version:   idHash,
		VersionID: vl.ID.Hex(),
		Equal:     req.IdHash == idHash,
		UserIDs:   userIDs,
	}, nil
}

func (s *friendServer) GetIncrementalFriends(ctx context.Context, req *relation.GetIncrementalFriendsReq) (*relation.GetIncrementalFriendsResp, error) {
	if err := authverify.CheckAccess(ctx, req.UserID); err != nil {
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
