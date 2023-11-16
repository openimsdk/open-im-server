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
	"github.com/openimsdk/open-im-server/v3/pkg/common/pagination"
	"time"
)

const (
	GroupModelTableName = "groups"
)

//type GroupModel struct {
//	GroupID                string    `gorm:"column:group_id;primary_key;size:64"                 json:"groupID"           binding:"required"`
//	GroupName              string    `gorm:"column:name;size:255"                                json:"groupName"`
//	Notification           string    `gorm:"column:notification;size:255"                        json:"notification"`
//	Introduction           string    `gorm:"column:introduction;size:255"                        json:"introduction"`
//	FaceURL                string    `gorm:"column:face_url;size:255"                            json:"faceURL"`
//	CreateTime             time.Time `gorm:"column:create_time;index:create_time;autoCreateTime"`
//	Ex                     string    `gorm:"column:ex"                                           json:"ex;size:1024"`
//	Status                 int32     `gorm:"column:status"`
//	CreatorUserID          string    `gorm:"column:creator_user_id;size:64"`
//	GroupType              int32     `gorm:"column:group_type"`
//	NeedVerification       int32     `gorm:"column:need_verification"`
//	LookMemberInfo         int32     `gorm:"column:look_member_info"                             json:"lookMemberInfo"`
//	ApplyMemberFriend      int32     `gorm:"column:apply_member_friend"                          json:"applyMemberFriend"`
//	NotificationUpdateTime time.Time `gorm:"column:notification_update_time"`
//	NotificationUserID     string    `gorm:"column:notification_user_id;size:64"`
//}

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

func (GroupModel) TableName() string {
	return GroupModelTableName
}

type GroupModelInterface interface {
	Create(ctx context.Context, groups []*GroupModel) (err error)
	UpdateMap(ctx context.Context, groupID string, args map[string]any) (err error)
	UpdateState(ctx context.Context, groupID string, state int32) (err error)
	Find(ctx context.Context, groupIDs []string) (groups []*GroupModel, err error)
	Take(ctx context.Context, groupID string) (group *GroupModel, err error)
	Search(ctx context.Context, keyword string, pagination pagination.Pagination) (total int64, groups []*GroupModel, err error)
	// 获取群总数
	CountTotal(ctx context.Context, before *time.Time) (count int64, err error)
	// 获取范围内群增量
	CountRangeEverydayTotal(ctx context.Context, start time.Time, end time.Time) (map[string]int64, error)
}
