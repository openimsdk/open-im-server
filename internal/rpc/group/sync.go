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
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
)

func (g *groupServer) GetFullGroupMemberUserIDs(ctx context.Context, req *pbgroup.GetFullGroupMemberUserIDsReq) (*pbgroup.GetFullGroupMemberUserIDsResp, error) {
	vl, err := g.db.FindMaxGroupMemberVersionCache(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	userIDs, err := g.db.FindGroupMemberUserID(ctx, req.GroupID)
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

func (g *groupServer) BatchGetIncrementalGroupMember(ctx context.Context, req *pbgroup.BatchGetIncrementalGroupMemberReq) (resp *pbgroup.BatchGetIncrementalGroupMemberResp, err error) {
	type VersionInfo struct {
		GroupID       string
		VersionID     string
		VersionNumber uint64
	}

	var groupIDs []string

	groupsVersionMap := make(map[string]*VersionInfo)
	groupsMap := make(map[string]*model.Group)
	hasGroupUpdateMap := make(map[string]bool)
	sortVersionMap := make(map[string]uint64)

	var targetKeys, versionIDs []string
	var versionNumbers []uint64

	var requestBodyLen int

	for _, group := range req.ReqList {
		groupsVersionMap[group.GroupID] = &VersionInfo{
			GroupID:       group.GroupID,
			VersionID:     group.VersionID,
			VersionNumber: group.Version,
		}

		groupIDs = append(groupIDs, group.GroupID)
	}

	groups, err := g.db.FindGroup(ctx, groupIDs)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	for _, group := range groups {
		if group.Status == constant.GroupStatusDismissed {
			err = servererrs.ErrDismissedAlready.Wrap()
			log.ZError(ctx, "This group is Dismissed Already", err, "group is", group.GroupID)

			delete(groupsVersionMap, group.GroupID)
		} else {
			groupsMap[group.GroupID] = group
		}
	}

	for groupID, vInfo := range groupsVersionMap {
		targetKeys = append(targetKeys, groupID)
		versionIDs = append(versionIDs, vInfo.VersionID)
		versionNumbers = append(versionNumbers, vInfo.VersionNumber)
	}

	opt := incrversion.BatchOption[[]*sdkws.GroupMemberFullInfo, pbgroup.BatchGetIncrementalGroupMemberResp]{
		Ctx:            ctx,
		TargetKeys:     targetKeys,
		VersionIDs:     versionIDs,
		VersionNumbers: versionNumbers,
		Versions: func(ctx context.Context, groupIDs []string, versions []uint64, limits []int) (map[string]*model.VersionLog, error) {
			vLogs, err := g.db.BatchFindMemberIncrVersion(ctx, groupIDs, versions, limits)
			if err != nil {
				return nil, errs.Wrap(err)
			}

			for groupID, vlog := range vLogs {
				vlogElems := make([]model.VersionLogElem, 0, len(vlog.Logs))
				for i, log := range vlog.Logs {
					switch log.EID {
					case model.VersionGroupChangeID:
						vlog.LogLen--
						hasGroupUpdateMap[groupID] = true
					case model.VersionSortChangeID:
						vlog.LogLen--
						sortVersionMap[groupID] = uint64(log.Version)
					default:
						vlogElems = append(vlogElems, vlog.Logs[i])
					}
				}
				vlog.Logs = vlogElems
				if vlog.LogLen > 0 {
					hasGroupUpdateMap[groupID] = true
				}
			}

			return vLogs, nil
		},
		CacheMaxVersions: g.db.BatchFindMaxGroupMemberVersionCache,
		Find: func(ctx context.Context, groupID string, ids []string) ([]*sdkws.GroupMemberFullInfo, error) {
			memberInfo, err := g.getGroupMembersInfo(ctx, groupID, ids)
			if err != nil {
				return nil, err
			}

			return memberInfo, err
		},
		Resp: func(versions map[string]*model.VersionLog, deleteIdsMap map[string][]string, insertListMap, updateListMap map[string][]*sdkws.GroupMemberFullInfo, fullMap map[string]bool) *pbgroup.BatchGetIncrementalGroupMemberResp {
			resList := make(map[string]*pbgroup.GetIncrementalGroupMemberResp)

			for groupID, versionLog := range versions {
				resList[groupID] = &pbgroup.GetIncrementalGroupMemberResp{
					VersionID:   versionLog.ID.Hex(),
					Version:     uint64(versionLog.Version),
					Full:        fullMap[groupID],
					Delete:      deleteIdsMap[groupID],
					Insert:      insertListMap[groupID],
					Update:      updateListMap[groupID],
					SortVersion: sortVersionMap[groupID],
				}

				requestBodyLen += len(insertListMap[groupID]) + len(updateListMap[groupID]) + len(deleteIdsMap[groupID])
				if requestBodyLen > 200 {
					break
				}
			}

			return &pbgroup.BatchGetIncrementalGroupMemberResp{
				RespList: resList,
			}
		},
	}

	resp, err = opt.Build()
	if err != nil {
		return nil, errs.Wrap(err)
	}

	for groupID, val := range resp.RespList {
		if val.Full || hasGroupUpdateMap[groupID] {
			count, err := g.db.FindGroupMemberNum(ctx, groupID)
			if err != nil {
				return nil, err
			}

			owner, err := g.db.TakeGroupOwner(ctx, groupID)
			if err != nil {
				return nil, err
			}

			resp.RespList[groupID].Group = g.groupDB2PB(groupsMap[groupID], owner.UserID, count)
		}
	}

	return resp, nil

}

func (g *groupServer) GetIncrementalJoinGroup(ctx context.Context, req *pbgroup.GetIncrementalJoinGroupReq) (*pbgroup.GetIncrementalJoinGroupResp, error) {
	if err := authverify.CheckAccessV3(ctx, req.UserID, g.config.Share.IMAdminUserID); err != nil {
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
