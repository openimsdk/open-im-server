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

	pbGroup "github.com/OpenIMSDK/protocol/group"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/convert"
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
