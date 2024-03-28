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

type GroupRequestModel struct {
	UserID        string    `bson:"user_id"`
	GroupID       string    `bson:"group_id"`
	HandleResult  int32     `bson:"handle_result"`
	ReqMsg        string    `bson:"req_msg"`
	HandledMsg    string    `bson:"handled_msg"`
	ReqTime       time.Time `bson:"req_time"`
	HandleUserID  string    `bson:"handle_user_id"`
	HandledTime   time.Time `bson:"handled_time"`
	JoinSource    int32     `bson:"join_source"`
	InviterUserID string    `bson:"inviter_user_id"`
	Ex            string    `bson:"ex"`
}

type GroupRequestModelInterface interface {
	Create(ctx context.Context, groupRequests []*GroupRequestModel) (err error)
	Delete(ctx context.Context, groupID string, userID string) (err error)
	UpdateHandler(ctx context.Context, groupID string, userID string, handledMsg string, handleResult int32) (err error)
	Take(ctx context.Context, groupID string, userID string) (groupRequest *GroupRequestModel, err error)
	FindGroupRequests(ctx context.Context, groupID string, userIDs []string) ([]*GroupRequestModel, error)
	Page(ctx context.Context, userID string, pagination pagination.Pagination) (total int64, groups []*GroupRequestModel, err error)
	PageGroup(ctx context.Context, groupIDs []string, pagination pagination.Pagination) (total int64, groups []*GroupRequestModel, err error)
}
