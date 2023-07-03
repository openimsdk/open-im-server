package group

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/convert"
	pbGroup "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/group"
)

func (s *groupServer) GetGroupInfoCache(
	ctx context.Context,
	req *pbGroup.GetGroupInfoCacheReq,
) (resp *pbGroup.GetGroupInfoCacheResp, err error) {
	group, err := s.GroupDatabase.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	resp = &pbGroup.GetGroupInfoCacheResp{GroupInfo: convert.Db2PbGroupInfo(group, "", 0)}
	return resp, nil
}

func (s *groupServer) GetGroupMemberCache(
	ctx context.Context,
	req *pbGroup.GetGroupMemberCacheReq,
) (resp *pbGroup.GetGroupMemberCacheResp, err error) {
	members, err := s.GroupDatabase.TakeGroupMember(ctx, req.GroupID, req.GroupMemberID)
	if err != nil {
		return nil, err
	}
	resp = &pbGroup.GetGroupMemberCacheResp{Member: convert.Db2PbGroupMember(members)}
	return resp, nil
}
