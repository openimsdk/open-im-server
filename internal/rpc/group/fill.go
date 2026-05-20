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
	relationtb "github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/datautil"
)

func (s *groupServer) PopulateGroupMember(ctx context.Context, members ...*relationtb.GroupMember) error {
	return s.notification.PopulateGroupMember(ctx, members...)
}

// applyMemberDisplayNicknames 按当前用户视角重设群成员 Nickname：
// 好友用 remark；非好友用 firstName+lastName，为空则 fallback 到 nickname。
func (s *groupServer) applyMemberDisplayNicknames(ctx context.Context, members []*sdkws.GroupMemberFullInfo) error {
	if len(members) == 0 {
		return nil
	}
	viewerID := mcontext.GetOpUserID(ctx)
	if viewerID == "" {
		return nil
	}
	userIDs := datautil.Slice(members, func(m *sdkws.GroupMemberFullInfo) string { return m.UserID })
	users, err := s.userClient.GetUsersInfoMap(ctx, userIDs)
	if err != nil {
		return err
	}
	friendInfos, err := s.relationClient.GetFriendsInfo(ctx, viewerID, userIDs)
	if err != nil {
		return err
	}
	remarkMap := make(map[string]string, len(friendInfos))
	for _, f := range friendInfos {
		if f != nil {
			remarkMap[f.FriendUserID] = f.Remark
		}
	}
	for _, m := range members {
		if remark, ok := remarkMap[m.UserID]; ok && remark != "" {
			m.Nickname = remark
			continue
		}
		if name := convert.MemberDisplayNickname(users[m.UserID]); name != "" {
			m.Nickname = name
		}
	}
	return nil
}
