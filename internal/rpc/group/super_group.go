package group

import (
	"context"
	"fmt"
	"strings"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/convert"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	pbGroup "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/group"
	sdkws "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
)

func (s *groupServer) GetJoinedSuperGroupList(ctx context.Context, req *pbGroup.GetJoinedSuperGroupListReq) (*pbGroup.GetJoinedSuperGroupListResp, error) {
	resp := &pbGroup.GetJoinedSuperGroupListResp{}
	groupIDs, err := s.GroupDatabase.FindJoinSuperGroup(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if len(groupIDs) == 0 {
		return resp, nil
	}
	owners, err := s.FindGroupMember(ctx, groupIDs, nil, []int32{constant.GroupOwner})
	if err != nil {
		return nil, err
	}
	ownerMap := utils.SliceToMap(owners, func(e *relation.GroupMemberModel) string {
		return e.GroupID
	})
	if ids := utils.Single(groupIDs, utils.Keys(ownerMap)); len(ids) > 0 {
		return nil, errs.ErrData.Wrap(fmt.Sprintf("super group %s not owner", strings.Join(ids, ",")))
	}
	groups, err := s.GroupDatabase.FindGroup(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	groupMap := utils.SliceToMap(groups, func(e *relation.GroupModel) string {
		return e.GroupID
	})
	if ids := utils.Single(groupIDs, utils.Keys(groupMap)); len(ids) > 0 {
		return nil, errs.ErrData.Wrap(fmt.Sprintf("super group info %s not found", strings.Join(ids, ",")))
	}
	superGroupMembers, err := s.GroupDatabase.FindSuperGroup(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	superGroupMemberMap := utils.SliceToMapAny(superGroupMembers, func(e *unrelation.SuperGroupModel) (string, []string) {
		return e.GroupID, e.MemberIDs
	})
	resp.Groups = utils.Slice(groupIDs, func(groupID string) *sdkws.GroupInfo {
		return convert.Db2PbGroupInfo(groupMap[groupID], ownerMap[groupID].UserID, uint32(len(superGroupMemberMap)))
	})
	return resp, nil
}

func (s *groupServer) GetSuperGroupsInfo(ctx context.Context, req *pbGroup.GetSuperGroupsInfoReq) (resp *pbGroup.GetSuperGroupsInfoResp, err error) {
	resp = &pbGroup.GetSuperGroupsInfoResp{}
	if len(req.GroupIDs) == 0 {
		return nil, errs.ErrArgs.Wrap("groupIDs empty")
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
	owners, err := s.FindGroupMember(ctx, req.GroupIDs, nil, []int32{constant.GroupOwner})
	if err != nil {
		return nil, err
	}
	ownerMap := utils.SliceToMap(owners, func(e *relation.GroupMemberModel) string {
		return e.GroupID
	})
	resp.GroupInfos = utils.Slice(groups, func(e *relation.GroupModel) *sdkws.GroupInfo {
		return convert.Db2PbGroupInfo(e, ownerMap[e.GroupID].UserID, uint32(len(superGroupMemberMap[e.GroupID])))
	})
	return resp, nil
}
