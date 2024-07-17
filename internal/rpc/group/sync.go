package group

import (
	"context"
	"slices"

	"github.com/openimsdk/open-im-server/v3/internal/rpc/incrversion"
	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/open-im-server/v3/pkg/util/hashutil"
	"github.com/openimsdk/protocol/constant"
	pbgroup "github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
)

func (s *groupServer) GetFullGroupMemberUserIDs(ctx context.Context, req *pbgroup.GetFullGroupMemberUserIDsReq) (*pbgroup.GetFullGroupMemberUserIDsResp, error) {
	vl, err := s.db.FindMaxGroupMemberVersionCache(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	userIDs, err := s.db.FindGroupMemberUserID(ctx, req.GroupID)
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

func (s *groupServer) GetFullJoinGroupIDs(ctx context.Context, req *pbgroup.GetFullJoinGroupIDsReq) (*pbgroup.GetFullJoinGroupIDsResp, error) {
	vl, err := s.db.FindMaxJoinGroupVersionCache(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	groupIDs, err := s.db.FindJoinGroupID(ctx, req.UserID)
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

func (s *groupServer) GetIncrementalGroupMember(ctx context.Context, req *pbgroup.GetIncrementalGroupMemberReq) (*pbgroup.GetIncrementalGroupMemberResp, error) {
	group, err := s.db.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if group.Status == constant.GroupStatusDismissed {
		return nil, servererrs.ErrDismissedAlready.Wrap()
	}
	var hasGroupUpdate bool
	opt := incrversion.Option[*sdkws.GroupMemberFullInfo, pbgroup.GetIncrementalGroupMemberResp]{
		Ctx:           ctx,
		VersionKey:    req.GroupID,
		VersionID:     req.VersionID,
		VersionNumber: req.Version,
		Version: func(ctx context.Context, groupID string, version uint, limit int) (*model.VersionLog, error) {
			vl, err := s.db.FindMemberIncrVersion(ctx, groupID, version, limit)
			if err != nil {
				return nil, err
			}
			vl.Logs = slices.DeleteFunc(vl.Logs, func(elem model.VersionLogElem) bool {
				if elem.EID == "" {
					vl.LogLen--
					hasGroupUpdate = true
					return true
				}
				return false
			})
			if vl.LogLen > 0 {
				hasGroupUpdate = true
			}
			return vl, nil
		},
		CacheMaxVersion: s.db.FindMaxGroupMemberVersionCache,
		Find: func(ctx context.Context, ids []string) ([]*sdkws.GroupMemberFullInfo, error) {
			return s.getGroupMembersInfo(ctx, req.GroupID, ids)
		},
		Resp: func(version *model.VersionLog, delIDs []string, insertList, updateList []*sdkws.GroupMemberFullInfo, full bool) *pbgroup.GetIncrementalGroupMemberResp {
			return &pbgroup.GetIncrementalGroupMemberResp{
				VersionID: version.ID.Hex(),
				Version:   uint64(version.Version),
				Full:      full,
				Delete:    delIDs,
				Insert:    insertList,
				Update:    updateList,
			}
		},
	}
	resp, err := opt.Build()
	if err != nil {
		return nil, err
	}
	if resp.Full || hasGroupUpdate {
		count, err := s.db.FindGroupMemberNum(ctx, group.GroupID)
		if err != nil {
			return nil, err
		}
		owner, err := s.db.TakeGroupOwner(ctx, group.GroupID)
		if err != nil {
			return nil, err
		}
		resp.Group = s.groupDB2PB(group, owner.UserID, count)
	}
	return resp, nil
}

func (s *groupServer) BatchGetIncrementalGroupMember(ctx context.Context, req *pbgroup.BatchGetIncrementalGroupMemberReq) (resp *pbgroup.BatchGetIncrementalGroupMemberResp, err error) {
	type VersionInfo struct {
		GroupID       string
		VersionID     string
		VersionNumber uint64
	}

	var groupIDs []string
	groupVersionMap := make(map[string]*VersionInfo)
	groupsMap := make(map[string]*model.Group)

	var targetKeys, versionIDs []string
	var versionNumbers []uint64

	// var requestBodyLen int

	for _, group := range req.ReqList {
		groupVersionMap[group.GroupID] = &VersionInfo{
			GroupID:       group.GroupID,
			VersionID:     group.VersionID,
			VersionNumber: group.Version,
		}

		groupIDs = append(groupIDs, group.GroupID)
	}

	groups, err := s.db.FindGroup(ctx, groupIDs)
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		if group.Status == constant.GroupStatusDismissed {
			err = servererrs.ErrDismissedAlready.Wrap()
			log.ZError(ctx, "This group is Dismissed Already", err, "group is", group.GroupID)
			delete(groupVersionMap, group.GroupID)
		} else {
			groupsMap[group.GroupID] = group
			// truegroupIDs = append(truegroupIDs, group.GroupID)
		}
	}
	for key, val := range groupVersionMap {
		targetKeys = append(targetKeys, key)
		versionIDs = append(versionIDs, val.VersionID)
		versionNumbers = append(versionNumbers, val.VersionNumber)
	}

	var hasGroupUpdate map[string]bool
	opt := incrversion.BatchOption[[]*sdkws.GroupMemberFullInfo, pbgroup.BatchGetIncrementalGroupMemberResp]{
		Ctx:            ctx,
		TargetKeys:     targetKeys,
		VersionIDs:     versionIDs,
		VersionNumbers: versionNumbers,
		Versions: func(ctx context.Context, groupIDs []string, versions []uint64, limits []int) (map[string]*model.VersionLog, error) {
			vLogs, err := s.db.BatchFindMemberIncrVersion(ctx, groupIDs, versions, limits)
			if err != nil {
				return nil, err
			}
			for key, vlog := range vLogs {
				vlog.Logs = slices.DeleteFunc(vlog.Logs, func(elem model.VersionLogElem) bool {
					if elem.EID == "" {
						vlog.LogLen--
						hasGroupUpdate[key] = true
						return true
					}
					return false
				})
				if vlog.LogLen > 0 {
					hasGroupUpdate[key] = true
				}
			}

			return vLogs, nil
		},
		CacheMaxVersions: s.db.BatchFindMaxGroupMemberVersionCache,
		Find: func(ctx context.Context, groupID string, ids []string) ([]*sdkws.GroupMemberFullInfo, error) {
			// memberInfoMap := make(map[string][]*sdkws.GroupMemberFullInfo)
			// for _, groupID := range groupIDs {
			memberInfo, err := s.getGroupMembersInfo(ctx, groupID, ids)
			if err != nil {
				return nil, err
			}
			// memberInfoMap:=datautil.SliceToMap(memberInfo, func(e *sdkws.GroupMemberFullInfo) string {
			// 	return e.GroupID
			// })
			// // memberInfoMap[groupID] = memberInfo
			// // }
			return memberInfo, err
		},
		Resp: func(versions map[string]*model.VersionLog, deleteIdsMap map[string][]string, insertListMap, updateListMap map[string][]*sdkws.GroupMemberFullInfo, fullMap map[string]bool) *pbgroup.BatchGetIncrementalGroupMemberResp {
			resList := make(map[string]*pbgroup.GetIncrementalGroupMemberResp)

			for key, version := range versions {
				resList[key] = &pbgroup.GetIncrementalGroupMemberResp{
					VersionID: version.ID.Hex(),
					Version:   uint64(version.Version),
					Full:      fullMap[key],
					Delete:    deleteIdsMap[key],
					Insert:    insertListMap[key],
					Update:    updateListMap[key],
				}
			}

			return &pbgroup.BatchGetIncrementalGroupMemberResp{
				RespList: resList,
			}
		},
	}

	resp, err = opt.Build()
	if err != nil {
		return nil, err
	}
	for key, val := range resp.RespList {
		if val.Full || hasGroupUpdate[key] {
			count, err := s.db.FindGroupMemberNum(ctx, key)
			if err != nil {
				return nil, err
			}
			owner, err := s.db.TakeGroupOwner(ctx, key)
			if err != nil {
				return nil, err
			}
			resp.RespList[key].Group = s.groupDB2PB(groupsMap[key], owner.UserID, count)
		}
	}

	return resp, nil

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
		Version:         s.db.FindJoinIncrVersion,
		CacheMaxVersion: s.db.FindMaxJoinGroupVersionCache,
		Find:            s.getGroupsInfo,
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
