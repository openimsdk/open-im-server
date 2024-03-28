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

	"github.com/openimsdk/tools/db/pagination"
)

type GroupModel struct {
	GroupID                string    `bson:"group_id"`
	GroupName              string    `bson:"group_name"`
	Notification           string    `bson:"notification"`
	Introduction           string    `bson:"introduction"`
	FaceURL                string    `bson:"face_url"`
	CreateTime             time.Time `bson:"create_time"`
	Ex                     string    `bson:"ex"`
	Status                 int32     `bson:"status"`
	CreatorUserID          string    `bson:"creator_user_id"`
	GroupType              int32     `bson:"group_type"`
	NeedVerification       int32     `bson:"need_verification"`
	LookMemberInfo         int32     `bson:"look_member_info"`
	ApplyMemberFriend      int32     `bson:"apply_member_friend"`
	NotificationUpdateTime time.Time `bson:"notification_update_time"`
	NotificationUserID     string    `bson:"notification_user_id"`
}

type GroupModelInterface interface {
	Create(ctx context.Context, groups []*GroupModel) (err error)
	UpdateMap(ctx context.Context, groupID string, args map[string]any) (err error)
	UpdateStatus(ctx context.Context, groupID string, status int32) (err error)
	Find(ctx context.Context, groupIDs []string) (groups []*GroupModel, err error)
	Take(ctx context.Context, groupID string) (group *GroupModel, err error)
	Search(ctx context.Context, keyword string, pagination pagination.Pagination) (total int64, groups []*GroupModel, err error)
	// Get Group total quantity
	CountTotal(ctx context.Context, before *time.Time) (count int64, err error)
	// Get Group total quantity every day
	CountRangeEverydayTotal(ctx context.Context, start time.Time, end time.Time) (map[string]int64, error)
}
