// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package grouphash

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"

	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/utils/datautil"
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
		datautil.Sort(userIDs, true)
	}
	memberMap := datautil.SliceToMap(members, func(e *sdkws.GroupMemberFullInfo) string {
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
