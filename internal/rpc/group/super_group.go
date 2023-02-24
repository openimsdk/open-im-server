package group

import (
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/db/table/relation"
	"OpenIM/pkg/common/db/table/unrelation"
	pbGroup "OpenIM/pkg/proto/group"
	sdkws "OpenIM/pkg/proto/sdkws"
	"OpenIM/pkg/utils"
	"context"
	"fmt"
	"strings"
)

func (s *groupServer) GetJoinedSuperGroupList(ctx context.Context, req *pbGroup.GetJoinedSuperGroupListReq) (*pbGroup.GetJoinedSuperGroupListResp, error) {
	resp := &pbGroup.GetJoinedSuperGroupListResp{}
	joinSuperGroup, err := s.GroupDatabase.FindJoinSuperGroup(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if len(joinSuperGroup.GroupIDs) == 0 {
		return resp, nil
	}
	owners, err := s.GroupDatabase.FindGroupMember(ctx, joinSuperGroup.GroupIDs, nil, []int32{constant.GroupOwner})
	if err != nil {
		return nil, err
	}
	ownerMap := utils.SliceToMap(owners, func(e *relation.GroupMemberModel) string {
		return e.GroupID
	})
	if ids := utils.Single(joinSuperGroup.GroupIDs, utils.Keys(ownerMap)); len(ids) > 0 {
		return nil, constant.ErrData.Wrap(fmt.Sprintf("super group %s not owner", strings.Join(ids, ",")))
	}
	groups, err := s.GroupDatabase.FindGroup(ctx, joinSuperGroup.GroupIDs)
	if err != nil {
		return nil, err
	}
	groupMap := utils.SliceToMap(groups, func(e *relation.GroupModel) string {
		return e.GroupID
	})
	if ids := utils.Single(joinSuperGroup.GroupIDs, utils.Keys(groupMap)); len(ids) > 0 {
		return nil, constant.ErrData.Wrap(fmt.Sprintf("super group info %s not found", strings.Join(ids, ",")))
	}
	superGroupMembers, err := s.GroupDatabase.FindSuperGroup(ctx, joinSuperGroup.GroupIDs)
	if err != nil {
		return nil, err
	}
	superGroupMemberMap := utils.SliceToMapAny(superGroupMembers, func(e *unrelation.SuperGroupModel) (string, []string) {
		return e.GroupID, e.MemberIDs
	})
	resp.Groups = utils.Slice(joinSuperGroup.GroupIDs, func(groupID string) *sdkws.GroupInfo {
		return DbToPbGroupInfo(groupMap[groupID], ownerMap[groupID].UserID, uint32(len(superGroupMemberMap)))
	})
	return resp, nil
}

func (s *groupServer) GetSuperGroupsInfo(ctx context.Context, req *pbGroup.GetSuperGroupsInfoReq) (resp *pbGroup.GetSuperGroupsInfoResp, err error) {
	resp = &pbGroup.GetSuperGroupsInfoResp{}
	if len(req.GroupIDs) == 0 {
		return nil, constant.ErrArgs.Wrap("groupIDs empty")
	}
	groups, err := s.GroupDatabase.FindGroup(ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}
	superGroupMembers, err := s.GroupDatabase.FindSuperGroup(ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}
	superGroupMemberMap := utils.SliceToMapAny(superGroupMembers, func(e *unrelation.SuperGroupModel) (string, []string) {
		return e.GroupID, e.MemberIDs
	})
	owners, err := s.GroupDatabase.FindGroupMember(ctx, req.GroupIDs, nil, []int32{constant.GroupOwner})
	if err != nil {
		return nil, err
	}
	ownerMap := utils.SliceToMap(owners, func(e *relation.GroupMemberModel) string {
		return e.GroupID
	})
	resp.GroupInfos = utils.Slice(groups, func(e *relation.GroupModel) *sdkws.GroupInfo {
		return DbToPbGroupInfo(e, ownerMap[e.GroupID].UserID, uint32(len(superGroupMemberMap[e.GroupID])))
	})
	return resp, nil
}
