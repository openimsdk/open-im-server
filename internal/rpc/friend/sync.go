package friend

import (
	"context"
	"github.com/openimsdk/protocol/sdkws"

	"github.com/openimsdk/open-im-server/v3/internal/rpc/incrversion"
	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/protocol/relation"
)

func (s *friendServer) GetIncrementalFriends(ctx context.Context, req *relation.GetIncrementalFriendsReq) (*relation.GetIncrementalFriendsResp, error) {
	if err := authverify.CheckAccessV3(ctx, req.UserID, s.config.Share.IMAdminUserID); err != nil {
		return nil, err
	}
	opt := incrversion.Option[*sdkws.FriendInfo, relation.GetIncrementalFriendsResp]{
		Ctx:             ctx,
		VersionKey:      req.UserID,
		VersionID:       req.VersionID,
		VersionNumber:   req.Version,
		Version:         s.friendDatabase.FindFriendIncrVersion,
		CacheMaxVersion: s.friendDatabase.FindMaxFriendVersionCache,
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
