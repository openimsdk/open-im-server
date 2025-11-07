package group

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/convert"
	pbgroup "github.com/openimsdk/protocol/group"
)

// GetGroupInfoCache get group info from cache.
func (g *groupServer) GetGroupInfoCache(ctx context.Context, req *pbgroup.GetGroupInfoCacheReq) (*pbgroup.GetGroupInfoCacheResp, error) {
	group, err := g.db.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	return &pbgroup.GetGroupInfoCacheResp{
		GroupInfo: convert.Db2PbGroupInfo(group, "", 0),
	}, nil
}

func (g *groupServer) GetGroupMemberCache(ctx context.Context, req *pbgroup.GetGroupMemberCacheReq) (*pbgroup.GetGroupMemberCacheResp, error) {
	if err := g.checkAdminOrInGroup(ctx, req.GroupID); err != nil {
		return nil, err
	}
	members, err := g.db.TakeGroupMember(ctx, req.GroupID, req.GroupMemberID)
	if err != nil {
		return nil, err
	}
	return &pbgroup.GetGroupMemberCacheResp{
		Member: convert.Db2PbGroupMember(members),
	}, nil
}
