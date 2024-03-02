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
	FriendModelTableName = "friends"
)

type FriendModel struct {
	OwnerUserID    string    `gorm:"column:owner_user_id;primary_key;size:64"`
	FriendUserID   string    `gorm:"column:friend_user_id;primary_key;size:64"`
	Remark         string    `gorm:"column:remark;size:255"`
	CreateTime     time.Time `gorm:"column:create_time;autoCreateTime"`
	AddSource      int32     `gorm:"column:add_source"`
	OperatorUserID string    `gorm:"column:operator_user_id;size:64"`
	Ex             string    `gorm:"column:ex;size:1024"`
}

func (FriendModel) TableName() string {
	return FriendModelTableName
}

type FriendModelInterface interface {
	// Create inserts multiple friend records.
	Create(ctx context.Context, friends []*FriendModel) error
	// Delete removes specified friends for an owner user.
	Delete(ctx context.Context, ownerUserID string, friendUserIDs []string) error
	// UpdateByMap updates a single friend's information for an owner user based on a map of arguments. Zero values are updated.
	UpdateByMap(ctx context.Context, ownerUserID string, friendUserID string, args map[string]interface{}) error
	// Update modifies the information of friends, excluding zero values.
	Update(ctx context.Context, friends []*FriendModel) error
	// UpdateRemark updates the remark for a friend, supporting zero values.
	UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) error
	// Take retrieves a single friend's information. Returns an error if not found.
	Take(ctx context.Context, ownerUserID, friendUserID string) (*FriendModel, error)
	// FindUserState finds the friendship status between two users, returning both if a mutual friendship exists.
	FindUserState(ctx context.Context, userID1, userID2 string) ([]*FriendModel, error)
	// FindFriends retrieves a list of friends for an owner, not returning an error for non-existent friendUserIDs.
	FindFriends(ctx context.Context, ownerUserID string, friendUserIDs []string) ([]*FriendModel, error)
	// FindReversalFriends finds who has added the specified user as a friend, not returning an error for non-existent ownerUserIDs.
	FindReversalFriends(ctx context.Context, friendUserID string, ownerUserIDs []string) ([]*FriendModel, error)
	// FindOwnerFriends paginates through the friends list of an owner user.
	FindOwnerFriends(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) ([]*FriendModel, int64, error)
	// FindInWhoseFriends paginates through users who have added the specified user as a friend.
	FindInWhoseFriends(ctx context.Context, friendUserID string, pageNumber, showNumber int32) ([]*FriendModel, int64, error)
	// FindFriendUserIDs retrieves a list of friend user IDs for an owner user.
	FindFriendUserIDs(ctx context.Context, ownerUserID string) ([]string, error)
	// NewTx creates a new transactional instance of the FriendModelInterface.
	NewTx(tx any) FriendModelInterface
}
