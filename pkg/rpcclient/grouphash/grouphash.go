package grouphash

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"

	"github.com/OpenIMSDK/protocol/group"
	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/utils"
)

func NewGroupHashFromGroupClient(x group.GroupClient) *GroupHash {
	return &GroupHash{
		getGroupAllUserIDs: func(ctx context.Context, groupID string) ([]string, error) {
			resp, err := x.GetGroupMemberUserIDs(ctx, &group.GetGroupMemberUserIDsReq{GroupID: groupID})
			if err != nil {
				return nil, err
			}
			return resp.UserIDs, nil
		},
		getGroupMemberInfo: func(ctx context.Context, groupID string, userIDs []string) ([]*sdkws.GroupMemberFullInfo, error) {
			resp, err := x.GetGroupMembersInfo(ctx, &group.GetGroupMembersInfoReq{GroupID: groupID, UserIDs: userIDs})
			if err != nil {
				return nil, err
			}
			return resp.Members, nil
		},
	}
}

func NewGroupHashFromGroupServer(x group.GroupServer) *GroupHash {
	return &GroupHash{
		getGroupAllUserIDs: func(ctx context.Context, groupID string) ([]string, error) {
			resp, err := x.GetGroupMemberUserIDs(ctx, &group.GetGroupMemberUserIDsReq{GroupID: groupID})
			if err != nil {
				return nil, err
			}
			return resp.UserIDs, nil
		},
		getGroupMemberInfo: func(ctx context.Context, groupID string, userIDs []string) ([]*sdkws.GroupMemberFullInfo, error) {
			resp, err := x.GetGroupMembersInfo(ctx, &group.GetGroupMembersInfoReq{GroupID: groupID, UserIDs: userIDs})
			if err != nil {
				return nil, err
			}
			return resp.Members, nil
		},
	}
}

type GroupHash struct {
	getGroupAllUserIDs func(ctx context.Context, groupID string) ([]string, error)
	getGroupMemberInfo func(ctx context.Context, groupID string, userIDs []string) ([]*sdkws.GroupMemberFullInfo, error)
}

func (gh *GroupHash) GetGroupHash(ctx context.Context, groupID string) (uint64, error) {
	userIDs, err := gh.getGroupAllUserIDs(ctx, groupID)
	if err != nil {
		return 0, err
	}
	var members []*sdkws.GroupMemberFullInfo
	if len(userIDs) > 0 {
		members, err = gh.getGroupMemberInfo(ctx, groupID, userIDs)
		if err != nil {
			return 0, err
		}
		utils.Sort(userIDs, true)
	}
	memberMap := utils.SliceToMap(members, func(e *sdkws.GroupMemberFullInfo) string {
		return e.UserID
	})
	res := make([]*sdkws.GroupMemberFullInfo, 0, len(members))
	for _, userID := range userIDs {
		member, ok := memberMap[userID]
		if !ok {
			continue
		}
		member.AppMangerLevel = 0
		res = append(res, member)
	}
	data, err := json.Marshal(res)
	if err != nil {
		return 0, err
	}
	sum := md5.Sum(data)
	return binary.BigEndian.Uint64(sum[:]), nil
}
