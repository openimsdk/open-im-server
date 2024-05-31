package group

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/internal/rpc/incrversion"
	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	pbgroup "github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/sdkws"
)

func (s *groupServer) SearchGroupMember(ctx context.Context, req *pbgroup.SearchGroupMemberReq) (*pbgroup.SearchGroupMemberResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s *groupServer) GetIncrementalGroupMember(ctx context.Context, req *pbgroup.GetIncrementalGroupMemberReq) (*pbgroup.GetIncrementalGroupMemberResp, error) {
	opt := incrversion.Option[*sdkws.GroupMemberFullInfo, pbgroup.GetIncrementalGroupMemberResp]{
		VersionID: req.VersionID,
		Version: func() (*model.VersionLog, error) {
			return s.db.FindMemberIncrVersion(ctx, req.GroupID, uint(req.Version), incrversion.Limit(s.config.RpcConfig.GroupSyncCount, req.Version))
		},
		AllID: func() ([]string, error) {
			return s.db.FindSortGroupMemberUserIDs(ctx, req.GroupID)
		},
		Find: func(ids []string) ([]*sdkws.GroupMemberFullInfo, error) {
			return s.getGroupMembersInfo(ctx, req.GroupID, ids)
		},
		ID: func(elem *sdkws.GroupMemberFullInfo) string {
			return elem.UserID
		},
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
		VersionID: req.VersionID,
		Version: func() (*model.VersionLog, error) {
			return s.db.FindJoinIncrVersion(ctx, req.UserID, uint(req.Version), incrversion.Limit(s.config.RpcConfig.GroupSyncCount, req.Version))
		},
		AllID: func() ([]string, error) {
			return s.db.FindSortJoinGroupIDs(ctx, req.UserID)
		},
		Find: func(ids []string) ([]*sdkws.GroupInfo, error) {
			return s.getGroupsInfo(ctx, ids)
		},
		ID: func(elem *sdkws.GroupInfo) string {
			return elem.GroupID
		},
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
