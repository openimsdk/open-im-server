package group

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/internal/rpc/incrversion"
	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/open-im-server/v3/pkg/util/hashutil"
	"github.com/openimsdk/protocol/constant"
	pbgroup "github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/sdkws"
)

const versionSyncLimit = 500

func (g *groupServer) GetFullGroupMemberUserIDs(ctx context.Context, req *pbgroup.GetFullGroupMemberUserIDsReq) (*pbgroup.GetFullGroupMemberUserIDsResp, error) {
	userIDs, err := g.db.FindGroupMemberUserID(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if err := authverify.CheckAccessIn(ctx, userIDs...); err != nil {
		return nil, err
	}
	vl, err := g.db.FindMaxGroupMemberVersionCache(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	idHash := hashutil.IdHash(userIDs)
	if req.IdHash == idHash {
		userIDs = nil
	}
	return &pbgroup.GetFullGroupMemberUserIDsResp{
		Version:   idHash,
		VersionID: vl.ID.Hex(),
		Equal:     req.IdHash == idHash,
		UserIDs:   userIDs,
	}, nil
}

func (g *groupServer) GetFullJoinGroupIDs(ctx context.Context, req *pbgroup.GetFullJoinGroupIDsReq) (*pbgroup.GetFullJoinGroupIDsResp, error) {
	if err := authverify.CheckAccess(ctx, req.UserID); err != nil {
		return nil, err
	}
	vl, err := g.db.FindMaxJoinGroupVersionCache(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	groupIDs, err := g.db.FindJoinGroupID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	idHash := hashutil.IdHash(groupIDs)
	if req.IdHash == idHash {
		groupIDs = nil
	}
	return &pbgroup.GetFullJoinGroupIDsResp{
		Version:   idHash,
		VersionID: vl.ID.Hex(),
		Equal:     req.IdHash == idHash,
		GroupIDs:  groupIDs,
	}, nil
}

func (g *groupServer) GetIncrementalGroupMember(ctx context.Context, req *pbgroup.GetIncrementalGroupMemberReq) (*pbgroup.GetIncrementalGroupMemberResp, error) {
	if err := g.checkAdminOrInGroup(ctx, req.GroupID); err != nil {
		return nil, err
	}
	group, err := g.db.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if group.Status == constant.GroupStatusDismissed {
		return nil, servererrs.ErrDismissedAlready.Wrap()
	}
	var (
		hasGroupUpdate bool
		sortVersion    uint64
	)
	opt := incrversion.Option[*sdkws.GroupMemberFullInfo, pbgroup.GetIncrementalGroupMemberResp]{
		Ctx:           ctx,
		VersionKey:    req.GroupID,
		VersionID:     req.VersionID,
		VersionNumber: req.Version,
		Version: func(ctx context.Context, groupID string, version uint, limit int) (*model.VersionLog, error) {
			vl, err := g.db.FindMemberIncrVersion(ctx, groupID, version, limit)
			if err != nil {
				return nil, err
			}
			logs := make([]model.VersionLogElem, 0, len(vl.Logs))
			for i, log := range vl.Logs {
				switch log.EID {
				case model.VersionGroupChangeID:
					vl.LogLen--
					hasGroupUpdate = true
				case model.VersionSortChangeID:
					vl.LogLen--
					sortVersion = uint64(log.Version)
				default:
					logs = append(logs, vl.Logs[i])
				}
			}
			vl.Logs = logs
			if vl.LogLen > 0 {
				hasGroupUpdate = true
			}
			return vl, nil
		},
		CacheMaxVersion: g.db.FindMaxGroupMemberVersionCache,
		Find: func(ctx context.Context, ids []string) ([]*sdkws.GroupMemberFullInfo, error) {
			return g.getGroupMembersInfo(ctx, req.GroupID, ids)
		},
		Resp: func(version *model.VersionLog, delIDs []string, insertList, updateList []*sdkws.GroupMemberFullInfo, full bool) *pbgroup.GetIncrementalGroupMemberResp {
			return &pbgroup.GetIncrementalGroupMemberResp{
				VersionID:   version.ID.Hex(),
				Version:     uint64(version.Version),
				Full:        full,
				Delete:      delIDs,
				Insert:      insertList,
				Update:      updateList,
				SortVersion: sortVersion,
			}
		},
	}
	resp, err := opt.Build()
	if err != nil {
		return nil, err
	}
	if resp.Full || hasGroupUpdate {
		count, err := g.db.FindGroupMemberNum(ctx, group.GroupID)
		if err != nil {
			return nil, err
		}
		owner, err := g.db.TakeGroupOwner(ctx, group.GroupID)
		if err != nil {
			return nil, err
		}
		resp.Group = g.groupDB2PB(group, owner.UserID, count)
	}
	return resp, nil
}

func (g *groupServer) GetIncrementalJoinGroup(ctx context.Context, req *pbgroup.GetIncrementalJoinGroupReq) (*pbgroup.GetIncrementalJoinGroupResp, error) {
	if err := authverify.CheckAccess(ctx, req.UserID); err != nil {
		return nil, err
	}
	opt := incrversion.Option[*sdkws.GroupInfo, pbgroup.GetIncrementalJoinGroupResp]{
		Ctx:             ctx,
		VersionKey:      req.UserID,
		VersionID:       req.VersionID,
		VersionNumber:   req.Version,
		Version:         g.db.FindJoinIncrVersion,
		CacheMaxVersion: g.db.FindMaxJoinGroupVersionCache,
		Find:            g.getGroupsInfo,
		Resp: func(version *model.VersionLog, delIDs []string, insertList, updateList []*sdkws.GroupInfo, full bool) *pbgroup.GetIncrementalJoinGroupResp {
			return &pbgroup.GetIncrementalJoinGroupResp{
				VersionID: version.ID.Hex(),
				Version:   uint64(version.Version),
				Full:      full,
				Delete:    delIDs,
				Insert:    insertList,
				Update:    updateList,
			}
		},
	}
	return opt.Build()
}

func (g *groupServer) BatchGetIncrementalGroupMember(ctx context.Context, req *pbgroup.BatchGetIncrementalGroupMemberReq) (*pbgroup.BatchGetIncrementalGroupMemberResp, error) {
	var num int
	resp := make(map[string]*pbgroup.GetIncrementalGroupMemberResp)
	for _, memberReq := range req.ReqList {
		if _, ok := resp[memberReq.GroupID]; ok {
			continue
		}
		memberResp, err := g.GetIncrementalGroupMember(ctx, memberReq)
		if err != nil {
			return nil, err
		}
		resp[memberReq.GroupID] = memberResp
		num += len(memberResp.Insert) + len(memberResp.Update) + len(memberResp.Delete)
		if num >= versionSyncLimit {
			break
		}
	}
	return &pbgroup.BatchGetIncrementalGroupMemberResp{RespList: resp}, nil
}
