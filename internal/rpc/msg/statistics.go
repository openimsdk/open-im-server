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

package msg

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"time"
)

func (m *msgServer) GetActiveUser(ctx context.Context, req *msg.GetActiveUserReq) (*msg.GetActiveUserResp, error) {
	msgCount, userCount, users, dateCount, err := m.MsgDatabase.RangeUserSendCount(ctx, time.UnixMilli(req.Start), time.UnixMilli(req.End), req.Group, req.Ase, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	var pbUsers []*msg.ActiveUser
	if len(users) > 0 {
		userIDs := utils.Slice(users, func(e *unrelation.UserCount) string { return e.UserID })
		userMap, err := m.User.GetUsersInfoMap(ctx, userIDs)
		if err != nil {
			return nil, err
		}
		pbUsers = make([]*msg.ActiveUser, 0, len(users))
		for _, user := range users {
			pbUser := userMap[user.UserID]
			if pbUser == nil {
				pbUser = &sdkws.UserInfo{
					UserID:   user.UserID,
					Nickname: user.UserID,
				}
			}
			pbUsers = append(pbUsers, &msg.ActiveUser{
				User:  pbUser,
				Count: user.Count,
			})
		}
	}
	return &msg.GetActiveUserResp{
		MsgCount:  msgCount,
		UserCount: userCount,
		DateCount: dateCount,
		Users:     pbUsers,
	}, nil
}

func (m *msgServer) GetActiveGroup(ctx context.Context, req *msg.GetActiveGroupReq) (*msg.GetActiveGroupResp, error) {
	msgCount, groupCount, groups, dateCount, err := m.MsgDatabase.RangeGroupSendCount(ctx, time.UnixMilli(req.Start), time.UnixMilli(req.End), req.Ase, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	var pbGroups []*msg.ActiveGroup
	if len(groups) > 0 {
		groupIDs := utils.Slice(groups, func(e *unrelation.GroupCount) string { return e.GroupID })
		resp, err := m.Group.GetGroupInfos(ctx, groupIDs, false)
		if err != nil {
			return nil, err
		}
		groupMap := make(map[string]*sdkws.GroupInfo, len(groups))
		for i, group := range groups {
			groupMap[group.GroupID] = resp[i]
		}
		pbGroups = make([]*msg.ActiveGroup, 0, len(groups))
		for _, group := range groups {
			pbGroup := groupMap[group.GroupID]
			if pbGroup == nil {
				pbGroup = &sdkws.GroupInfo{
					GroupID:   group.GroupID,
					GroupName: group.GroupID,
				}
			}
			pbGroups = append(pbGroups, &msg.ActiveGroup{
				Group: pbGroup,
				Count: group.Count,
			})
		}
	}
	return &msg.GetActiveGroupResp{
		MsgCount:   msgCount,
		GroupCount: groupCount,
		DateCount:  dateCount,
		Groups:     pbGroups,
	}, nil
}
