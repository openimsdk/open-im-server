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

package database

import (
	"context"

	"github.com/openimsdk/tools/db/pagination"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

type GroupMember interface {
	Create(ctx context.Context, groupMembers []*model.GroupMember) (err error)
	Delete(ctx context.Context, groupID string, userIDs []string) (err error)
	Update(ctx context.Context, groupID string, userID string, data map[string]any) (err error)
	UpdateRoleLevel(ctx context.Context, groupID string, userID string, roleLevel int32) error
	UpdateUserRoleLevels(ctx context.Context, groupID string, firstUserID string, firstUserRoleLevel int32, secondUserID string, secondUserRoleLevel int32) error
	FindMemberUserID(ctx context.Context, groupID string) (userIDs []string, err error)
	Take(ctx context.Context, groupID string, userID string) (groupMember *model.GroupMember, err error)
	Find(ctx context.Context, groupID string, userIDs []string) ([]*model.GroupMember, error)
	FindInGroup(ctx context.Context, userID string, groupIDs []string) ([]*model.GroupMember, error)
	TakeOwner(ctx context.Context, groupID string) (groupMember *model.GroupMember, err error)
	SearchMember(ctx context.Context, keyword string, groupID string, pagination pagination.Pagination) (total int64, groupList []*model.GroupMember, err error)
	FindRoleLevelUserIDs(ctx context.Context, groupID string, roleLevel int32) ([]string, error)
	FindUserJoinedGroupID(ctx context.Context, userID string) (groupIDs []string, err error)
	TakeGroupMemberNum(ctx context.Context, groupID string) (count int64, err error)
	FindUserManagedGroupID(ctx context.Context, userID string) (groupIDs []string, err error)
	IsUpdateRoleLevel(data map[string]any) bool
	JoinGroupIncrVersion(ctx context.Context, userID string, groupIDs []string, state int32) error
	MemberGroupIncrVersion(ctx context.Context, groupID string, userIDs []string, state int32) error
	FindMemberIncrVersion(ctx context.Context, groupID string, version uint, limit int) (*model.VersionLog, error)
	BatchFindMemberIncrVersion(ctx context.Context, groupIDs []string, versions []uint, limits []int) ([]*model.VersionLog, error)
	FindJoinIncrVersion(ctx context.Context, userID string, version uint, limit int) (*model.VersionLog, error)
}
