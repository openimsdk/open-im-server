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

// FriendModel represents the data structure for a friend relationship in MongoDB.
type FriendModel struct {
	OwnerUserID    string    `bson:"owner_user_id"`
	FriendUserID   string    `bson:"friend_user_id"`
	Remark         string    `bson:"remark"`
	CreateTime     time.Time `bson:"create_time"`
	AddSource      int32     `bson:"add_source"`
	OperatorUserID string    `bson:"operator_user_id"`
	Ex             string    `bson:"ex"`
}

// FriendModelInterface defines the operations for managing friends in MongoDB.
type FriendModelInterface interface {
	// Create inserts multiple friend records.
	Create(ctx context.Context, friends []*FriendModel) (err error)
	// Delete removes specified friends of the owner user.
	Delete(ctx context.Context, ownerUserID string, friendUserIDs []string) (err error)
	// UpdateByMap updates specific fields of a friend document using a map.
	UpdateByMap(ctx context.Context, ownerUserID string, friendUserID string, args map[string]any) (err error)
	// UpdateRemark modify remarks.
	UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) (err error)
	// Take retrieves a single friend document. Returns an error if not found.
	Take(ctx context.Context, ownerUserID, friendUserID string) (friend *FriendModel, err error)
	// FindUserState finds the friendship status between two users.
	FindUserState(ctx context.Context, userID1, userID2 string) (friends []*FriendModel, err error)
	// FindFriends retrieves a list of friends for a given owner. Missing friends do not cause an error.
	FindFriends(ctx context.Context, ownerUserID string, friendUserIDs []string) (friends []*FriendModel, err error)
	// FindReversalFriends finds users who have added the specified user as a friend.
	FindReversalFriends(ctx context.Context, friendUserID string, ownerUserIDs []string) (friends []*FriendModel, err error)
	// FindOwnerFriends retrieves a paginated list of friends for a given owner.
	FindOwnerFriends(ctx context.Context, ownerUserID string, pagination pagination.Pagination) (total int64, friends []*FriendModel, err error)
	// FindInWhoseFriends finds users who have added the specified user as a friend, with pagination.
	FindInWhoseFriends(ctx context.Context, friendUserID string, pagination pagination.Pagination) (total int64, friends []*FriendModel, err error)
	// FindFriendUserIDs retrieves a list of friend user IDs for a given owner.
	FindFriendUserIDs(ctx context.Context, ownerUserID string) (friendUserIDs []string, err error)
}
