package group

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/internal/rpc/incrversion"
	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	pbgroup "github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/sdkws"
)

func (s *groupServer) GetIncrementalGroupMember(ctx context.Context, req *pbgroup.GetIncrementalGroupMemberReq) (*pbgroup.GetIncrementalGroupMemberResp, error) {
	opt := incrversion.Option[*sdkws.GroupMemberFullInfo, pbgroup.GetIncrementalGroupMemberResp]{
		Ctx:             ctx,
		VersionKey:      req.GroupID,
		VersionID:       req.VersionID,
		VersionNumber:   req.Version,
		SyncLimit:       s.config.RpcConfig.GroupSyncCount,
		Version:         s.db.FindMemberIncrVersion,
		CacheMaxVersion: s.db.FindMaxGroupMemberVersionCache,
		SortID:          s.db.FindSortGroupMemberUserIDs,
		Find: func(ctx context.Context, ids []string) ([]*sdkws.GroupMemberFullInfo, error) {
			return s.getGroupMembersInfo(ctx, req.GroupID, ids)
		},
		ID: func(elem *sdkws.GroupMemberFullInfo) string { return elem.UserID },
		Resp: func(version *model.VersionLog, delIDs []string, list []*sdkws.GroupMemberFullInfo, full bool) *pbgroup.GetIncrementalGroupMemberResp {
			return &pbgroup.GetIncrementalGroupMemberResp{
				VersionID:     version.ID.Hex(),
				Version:       uint64(version.Version),
				Full:          full,
				SyncCount:     uint32(s.config.RpcConfig.GroupSyncCount),
				DeleteUserIds: delIDs,
				Changes:       list,
			}
		},
	}
	return opt.Build()
}

func (s *groupServer) GetIncrementalJoinGroup(ctx context.Context, req *pbgroup.GetIncrementalJoinGroupReq) (*pbgroup.GetIncrementalJoinGroupResp, error) {
	if err := authverify.CheckAccessV3(ctx, req.UserID, s.config.Share.IMAdminUserID); err != nil {
		return nil, err
	}
	opt := incrversion.Option[*sdkws.GroupInfo, pbgroup.GetIncrementalJoinGroupResp]{
		Ctx:             ctx,
		VersionKey:      req.UserID,
		VersionID:       req.VersionID,
		VersionNumber:   req.Version,
		SyncLimit:       s.config.RpcConfig.GroupSyncCount,
		Version:         s.db.FindJoinIncrVersion,
		CacheMaxVersion: s.db.FindMaxJoinGroupVersionCache,
		SortID:          s.db.FindSortJoinGroupIDs,
		Find:            s.getGroupsInfo,
		ID:              func(elem *sdkws.GroupInfo) string { return elem.GroupID },
		Resp: func(version *model.VersionLog, delIDs []string, list []*sdkws.GroupInfo, full bool) *pbgroup.GetIncrementalJoinGroupResp {
			return &pbgroup.GetIncrementalJoinGroupResp{
				VersionID:      version.ID.Hex(),
				Version:        uint64(version.Version),
				Full:           full,
				SyncCount:      uint32(s.config.RpcConfig.GroupSyncCount),
				DeleteGroupIds: delIDs,
				Changes:        list,
			}
		},
	}
	return opt.Build()
}
