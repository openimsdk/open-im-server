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
)

const (
	GroupRequestModelTableName = "group_requests"
)

type GroupRequestModel struct {
	UserID        string    `gorm:"column:user_id;primary_key;size:64"`
	GroupID       string    `gorm:"column:group_id;primary_key;size:64"`
	HandleResult  int32     `gorm:"column:handle_result"`
	ReqMsg        string    `gorm:"column:req_msg;size:1024"`
	HandledMsg    string    `gorm:"column:handle_msg;size:1024"`
	ReqTime       time.Time `gorm:"column:req_time"`
	HandleUserID  string    `gorm:"column:handle_user_id;size:64"`
	HandledTime   time.Time `gorm:"column:handle_time"`
	JoinSource    int32     `gorm:"column:join_source"`
	InviterUserID string    `gorm:"column:inviter_user_id;size:64"`
	Ex            string    `gorm:"column:ex;size:1024"`
}

func (GroupRequestModel) TableName() string {
	return GroupRequestModelTableName
}

type GroupRequestModelInterface interface {
	NewTx(tx any) GroupRequestModelInterface
	Create(ctx context.Context, groupRequests []*GroupRequestModel) (err error)
	Delete(ctx context.Context, groupID string, userID string) (err error)
	UpdateHandler(ctx context.Context, groupID string, userID string, handledMsg string, handleResult int32) (err error)
	Take(ctx context.Context, groupID string, userID string) (groupRequest *GroupRequestModel, err error)
	FindGroupRequests(ctx context.Context, groupID string, userIDs []string) (int64, []*GroupRequestModel, error)
	Page(
		ctx context.Context,
		userID string,
		pageNumber, showNumber int32,
	) (total uint32, groups []*GroupRequestModel, err error)
	PageGroup(
		ctx context.Context,
		groupIDs []string,
		pageNumber, showNumber int32,
	) (total uint32, groups []*GroupRequestModel, err error)
}
