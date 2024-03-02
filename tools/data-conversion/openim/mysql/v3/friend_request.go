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

const FriendRequestModelTableName = "friend_requests"

type FriendRequestModel struct {
	FromUserID    string    `gorm:"column:from_user_id;primary_key;size:64"`
	ToUserID      string    `gorm:"column:to_user_id;primary_key;size:64"`
	HandleResult  int32     `gorm:"column:handle_result"`
	ReqMsg        string    `gorm:"column:req_msg;size:255"`
	CreateTime    time.Time `gorm:"column:create_time; autoCreateTime"`
	HandlerUserID string    `gorm:"column:handler_user_id;size:64"`
	HandleMsg     string    `gorm:"column:handle_msg;size:255"`
	HandleTime    time.Time `gorm:"column:handle_time"`
	Ex            string    `gorm:"column:ex;size:1024"`
}

func (FriendRequestModel) TableName() string {
	return FriendRequestModelTableName
}

type FriendRequestModelInterface interface {
	// Insert multiple records
	Create(ctx context.Context, friendRequests []*FriendRequestModel) (err error)

	// Delete a record
	Delete(ctx context.Context, fromUserID, toUserID string) (err error)

	// Update records with zero values based on a map of changes
	UpdateByMap(ctx context.Context, formUserID, toUserID string, args map[string]interface{}) (err error)

	// Update multiple records (non-zero values)
	Update(ctx context.Context, friendRequest *FriendRequestModel) (err error)

	// Find a friend request sent to a specific user; does not return an error if not found
	Find(ctx context.Context, fromUserID, toUserID string) (friendRequest *FriendRequestModel, err error)

	// Alias for Find (retrieves a friend request between two users)
	Take(ctx context.Context, fromUserID, toUserID string) (friendRequest *FriendRequestModel, err error)

	// Get a list of friend requests received by `toUserID`
	FindToUserID(ctx context.Context, toUserID string, pageNumber, showNumber int32) (friendRequests []*FriendRequestModel, total int64, err error)

	// Get a list of friend requests sent by `fromUserID`
	FindFromUserID(ctx context.Context, fromUserID string, pageNumber, showNumber int32) (friendRequests []*FriendRequestModel, total int64, err error)

	// Find all friend requests between two users (both directions)
	FindBothFriendRequests(ctx context.Context, fromUserID, toUserID string) (friends []*FriendRequestModel, err error)

	// Create a new transaction
	NewTx(tx any) FriendRequestModelInterface
}
