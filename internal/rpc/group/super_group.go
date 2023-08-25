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

package group

import (
	"context"
	"fmt"
	"strings"

	"github.com/OpenIMSDK/protocol/constant"
	pbgroup "github.com/OpenIMSDK/protocol/group"
	sdkws "github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/utils"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/convert"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/unrelation"
)

func (s *groupServer) GetJoinedSuperGroupList(
	ctx context.Context,
	req *pbgroup.GetJoinedSuperGroupListReq,
) (*pbgroup.GetJoinedSuperGroupListResp, error) {
	resp := &pbgroup.GetJoinedSuperGroupListResp{}
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
	superGroupMemberMap := utils.SliceToMapAny(
		superGroupMembers,
		func(e *unrelation.SuperGroupModel) (string, []string) {
			return e.GroupID, e.MemberIDs
		},
	)
	resp.Groups = utils.Slice(groupIDs, func(groupID string) *sdkws.GroupInfo {
		return convert.Db2PbGroupInfo(groupMap[groupID], ownerMap[groupID].UserID, uint32(len(superGroupMemberMap)))
	})
	return resp, nil
}

func (s *groupServer) GetSuperGroupsInfo(
	ctx context.Context,
	req *pbgroup.GetSuperGroupsInfoReq,
) (resp *pbgroup.GetSuperGroupsInfoResp, err error) {
	resp = &pbgroup.GetSuperGroupsInfoResp{}
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
	superGroupMemberMap := utils.SliceToMapAny(
		superGroupMembers,
		func(e *unrelation.SuperGroupModel) (string, []string) {
			return e.GroupID, e.MemberIDs
		},
	)
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
