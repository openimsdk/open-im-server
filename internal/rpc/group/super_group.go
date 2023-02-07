package group

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/table/relation"
	pbGroup "Open_IM/pkg/proto/group"
	sdk_ws "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
)

func (s *groupServer) GetJoinedSuperGroupList(ctx context.Context, req *pbGroup.GetJoinedSuperGroupListReq) (*pbGroup.GetJoinedSuperGroupListResp, error) {
	resp := &pbGroup.GetJoinedSuperGroupListResp{}
	total, groupIDs, err := s.GroupInterface.GetJoinSuperGroupByID(ctx, req.UserID, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	resp.Total = total
	if len(groupIDs) == 0 {
		return resp, nil
	}
	numMap, err := s.GroupInterface.GetSuperGroupMemberNum(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	ownerIdMap, err := s.GroupInterface.GetGroupOwnerUserID(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	groups, err := s.GroupInterface.FindGroupsByID(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	groupMap := utils.SliceToMap(groups, func(e *relation.GroupModel) string {
		return e.GroupID
	})
	resp.Groups = utils.Slice(groupIDs, func(groupID string) *sdk_ws.GroupInfo {
		return DbToPbGroupInfo(groupMap[groupID], ownerIdMap[groupID], numMap[groupID])
	})
	return resp, nil
}

func (s *groupServer) GetSuperGroupsInfo(ctx context.Context, req *pbGroup.GetSuperGroupsInfoReq) (resp *pbGroup.GetSuperGroupsInfoResp, err error) {
	resp = &pbGroup.GetSuperGroupsInfoResp{}
	if len(req.GroupIDs) == 0 {
		return nil, constant.ErrArgs.Wrap("groupIDs empty")
	}
	groups, err := s.GroupInterface.FindGroupsByID(ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}
	numMap, err := s.GroupInterface.GetSuperGroupMemberNum(ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}
	ownerIdMap, err := s.GroupInterface.GetGroupOwnerUserID(ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}
	resp.GroupInfos = utils.Slice(groups, func(e *relation.GroupModel) *sdk_ws.GroupInfo {
		return DbToPbGroupInfo(e, ownerIdMap[e.GroupID], numMap[e.GroupID])
	})
	return resp, nil
}
