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

package relation

import (
	"context"
	"time"

	"github.com/OpenIMSDK/tools/pagination"
)

type GroupMemberModel struct {
	GroupID        string    `bson:"group_id"`
	UserID         string    `bson:"user_id"`
	Nickname       string    `bson:"nickname"`
	FaceURL        string    `bson:"face_url"`
	RoleLevel      int32     `bson:"role_level"`
	JoinTime       time.Time `bson:"join_time"`
	JoinSource     int32     `bson:"join_source"`
	InviterUserID  string    `bson:"inviter_user_id"`
	OperatorUserID string    `bson:"operator_user_id"`
	MuteEndTime    time.Time `bson:"mute_end_time"`
	Ex             string    `bson:"ex"`
}

type GroupMemberModelInterface interface {
	//NewTx(tx any) GroupMemberModelInterface
	Create(ctx context.Context, groupMembers []*GroupMemberModel) (err error)
	Delete(ctx context.Context, groupID string, userIDs []string) (err error)
	//DeleteGroup(ctx context.Context, groupIDs []string) (err error)
	Update(ctx context.Context, groupID string, userID string, data map[string]any) (err error)
	UpdateRoleLevel(ctx context.Context, groupID string, userID string, roleLevel int32) error
	FindMemberUserID(ctx context.Context, groupID string) (userIDs []string, err error)
	Take(ctx context.Context, groupID string, userID string) (groupMember *GroupMemberModel, err error)
	TakeOwner(ctx context.Context, groupID string) (groupMember *GroupMemberModel, err error)
	SearchMember(ctx context.Context, keyword string, groupID string, pagination pagination.Pagination) (total int64, groupList []*GroupMemberModel, err error)
	FindRoleLevelUserIDs(ctx context.Context, groupID string, roleLevel int32) ([]string, error)
	//MapGroupMemberNum(ctx context.Context, groupIDs []string) (count map[string]uint32, err error)
	//FindJoinUserID(ctx context.Context, groupIDs []string) (groupUsers map[string][]string, err error)
	FindUserJoinedGroupID(ctx context.Context, userID string) (groupIDs []string, err error)
	TakeGroupMemberNum(ctx context.Context, groupID string) (count int64, err error)
	//FindUsersJoinedGroupID(ctx context.Context, userIDs []string) (map[string][]string, error)
	FindUserManagedGroupID(ctx context.Context, userID string) (groupIDs []string, err error)
	IsUpdateRoleLevel(data map[string]any) bool
}
