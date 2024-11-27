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

type FriendRequest interface {
	// Insert multiple records
	Create(ctx context.Context, friendRequests []*model.FriendRequest) (err error)
	// Delete record
	Delete(ctx context.Context, fromUserID, toUserID string) (err error)
	// Update with zero values
	UpdateByMap(ctx context.Context, formUserID string, toUserID string, args map[string]any) (err error)
	// Update multiple records (non-zero values)
	Update(ctx context.Context, friendRequest *model.FriendRequest) (err error)
	// Get friend requests sent to a specific user, no error returned if not found
	Find(ctx context.Context, fromUserID, toUserID string) (friendRequest *model.FriendRequest, err error)
	Take(ctx context.Context, fromUserID, toUserID string) (friendRequest *model.FriendRequest, err error)
	// Get list of friend requests received by toUserID
	FindToUserID(ctx context.Context, toUserID string, pagination pagination.Pagination) (total int64, friendRequests []*model.FriendRequest, err error)
	// Get list of friend requests sent by fromUserID
	FindFromUserID(ctx context.Context, fromUserID string, pagination pagination.Pagination) (total int64, friendRequests []*model.FriendRequest, err error)
	FindBothFriendRequests(ctx context.Context, fromUserID, toUserID string) (friends []*model.FriendRequest, err error)
}
