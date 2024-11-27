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

package cache

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

type UserCache interface {
	BatchDeleter
	CloneUserCache() UserCache
	GetUserInfo(ctx context.Context, userID string) (userInfo *model.User, err error)
	GetUsersInfo(ctx context.Context, userIDs []string) ([]*model.User, error)
	DelUsersInfo(userIDs ...string) UserCache
	GetUserGlobalRecvMsgOpt(ctx context.Context, userID string) (opt int, err error)
	DelUsersGlobalRecvMsgOpt(userIDs ...string) UserCache
	//GetUserStatus(ctx context.Context, userIDs []string) ([]*user.OnlineStatus, error)
	//SetUserStatus(ctx context.Context, userID string, status, platformID int32) error
}
