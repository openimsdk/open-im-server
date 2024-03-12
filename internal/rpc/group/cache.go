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

	pbgroup "github.com/OpenIMSDK/protocol/group"

	"github.com/openimsdk/open-im-server/v3/pkg/common/convert"
)

func (s *groupServer) GetGroupInfoCache(
	ctx context.Context,
	req *pbgroup.GetGroupInfoCacheReq,
) (resp *pbgroup.GetGroupInfoCacheResp, err error) {
	group, err := s.db.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	resp = &pbgroup.GetGroupInfoCacheResp{GroupInfo: convert.Db2PbGroupInfo(group, "", 0)}
	return resp, nil
}

func (s *groupServer) GetGroupMemberCache(
	ctx context.Context,
	req *pbgroup.GetGroupMemberCacheReq,
) (resp *pbgroup.GetGroupMemberCacheResp, err error) {
	members, err := s.db.TakeGroupMember(ctx, req.GroupID, req.GroupMemberID)
	if err != nil {
		return nil, err
	}
	resp = &pbgroup.GetGroupMemberCacheResp{Member: convert.Db2PbGroupMember(members)}
	return resp, nil
}
