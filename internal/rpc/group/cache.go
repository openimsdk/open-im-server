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

	"github.com/openimsdk/open-im-server/v3/pkg/common/convert"
	pbgroup "github.com/openimsdk/protocol/group"
)

// GetGroupInfoCache get group info from cache.
func (s *groupServer) GetGroupInfoCache(ctx context.Context, req *pbgroup.GetGroupInfoCacheReq) (*pbgroup.GetGroupInfoCacheResp, error) {
	group, err := s.db.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	return &pbgroup.GetGroupInfoCacheResp{
		GroupInfo: convert.Db2PbGroupInfo(group, "", 0),
	}, nil
}

func (s *groupServer) GetGroupMemberCache(ctx context.Context, req *pbgroup.GetGroupMemberCacheReq) (*pbgroup.GetGroupMemberCacheResp, error) {
	if err := s.checkAdminOrInGroup(ctx, req.GroupID); err != nil {
		return nil, err
	}
	members, err := s.db.TakeGroupMember(ctx, req.GroupID, req.GroupMemberID)
	if err != nil {
		return nil, err
	}
	return &pbgroup.GetGroupMemberCacheResp{
		Member: convert.Db2PbGroupMember(members),
	}, nil
}
