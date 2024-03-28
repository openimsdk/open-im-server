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

type FriendRequestModel struct {
	FromUserID    string    `bson:"from_user_id"`
	ToUserID      string    `bson:"to_user_id"`
	HandleResult  int32     `bson:"handle_result"`
	ReqMsg        string    `bson:"req_msg"`
	CreateTime    time.Time `bson:"create_time"`
	HandlerUserID string    `bson:"handler_user_id"`
	HandleMsg     string    `bson:"handle_msg"`
	HandleTime    time.Time `bson:"handle_time"`
	Ex            string    `bson:"ex"`
}

type FriendRequestModelInterface interface {
	// Insert multiple records
	Create(ctx context.Context, friendRequests []*FriendRequestModel) (err error)
	// Delete record
	Delete(ctx context.Context, fromUserID, toUserID string) (err error)
	// Update with zero values
	UpdateByMap(ctx context.Context, formUserID string, toUserID string, args map[string]any) (err error)
	// Update multiple records (non-zero values)
	Update(ctx context.Context, friendRequest *FriendRequestModel) (err error)
	// Get friend requests sent to a specific user, no error returned if not found
	Find(ctx context.Context, fromUserID, toUserID string) (friendRequest *FriendRequestModel, err error)
	Take(ctx context.Context, fromUserID, toUserID string) (friendRequest *FriendRequestModel, err error)
	// Get list of friend requests received by toUserID
	FindToUserID(ctx context.Context, toUserID string, pagination pagination.Pagination) (total int64, friendRequests []*FriendRequestModel, err error)
	// Get list of friend requests sent by fromUserID
	FindFromUserID(ctx context.Context, fromUserID string, pagination pagination.Pagination) (total int64, friendRequests []*FriendRequestModel, err error)
	FindBothFriendRequests(ctx context.Context, fromUserID, toUserID string) (friends []*FriendRequestModel, err error)
}
