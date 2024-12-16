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

package user

import (
	"context"

	relationtb "github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient/notification"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/sdkws"
)

type UserNotificationSender struct {
	*rpcclient.NotificationSender
	getUsersInfo func(ctx context.Context, userIDs []string) ([]notification.CommonUser, error)
	// db controller
	db controller.UserDatabase
}

type userNotificationSenderOptions func(*UserNotificationSender)

func WithUserDB(db controller.UserDatabase) userNotificationSenderOptions {
	return func(u *UserNotificationSender) {
		u.db = db
	}
}

func WithUserFunc(
	fn func(ctx context.Context, userIDs []string) (users []*relationtb.User, err error),
) userNotificationSenderOptions {
	return func(u *UserNotificationSender) {
		f := func(ctx context.Context, userIDs []string) (result []notification.CommonUser, err error) {
			users, err := fn(ctx, userIDs)
			if err != nil {
				return nil, err
			}
			for _, user := range users {
				result = append(result, user)
			}
			return result, nil
		}
		u.getUsersInfo = f
	}
}

func NewUserNotificationSender(config *Config, opts ...userNotificationSenderOptions) *UserNotificationSender {
	f := &UserNotificationSender{
		NotificationSender: rpcclient.NewNotificationSender(&config.NotificationConfig, rpcclient.WithRpcClient()),
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

/* func (u *UserNotificationSender) getUsersInfoMap(
	ctx context.Context,
	userIDs []string,
) (map[string]*sdkws.UserInfo, error) {
	users, err := u.getUsersInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	result := make(map[string]*sdkws.UserInfo)
	for _, user := range users {
		result[user.GetUserID()] = user.(*sdkws.UserInfo)
	}
	return result, nil
} */

/* func (u *UserNotificationSender) getFromToUserNickname(
	ctx context.Context,
	fromUserID, toUserID string,
) (string, string, error) {
	users, err := u.getUsersInfoMap(ctx, []string{fromUserID, toUserID})
	if err != nil {
		return "", "", nil
	}
	return users[fromUserID].Nickname, users[toUserID].Nickname, nil
} */

func (u *UserNotificationSender) UserStatusChangeNotification(
	ctx context.Context,
	tips *sdkws.UserStatusChangeTips,
) {
	u.Notification(ctx, tips.FromUserID, tips.ToUserID, constant.UserStatusChangeNotification, tips)
}
func (u *UserNotificationSender) UserCommandUpdateNotification(
	ctx context.Context,
	tips *sdkws.UserCommandUpdateTips,
) {
	u.Notification(ctx, tips.FromUserID, tips.ToUserID, constant.UserCommandUpdateNotification, tips)
}
func (u *UserNotificationSender) UserCommandAddNotification(
	ctx context.Context,
	tips *sdkws.UserCommandAddTips,
) {
	u.Notification(ctx, tips.FromUserID, tips.ToUserID, constant.UserCommandAddNotification, tips)
}
func (u *UserNotificationSender) UserCommandDeleteNotification(
	ctx context.Context,
	tips *sdkws.UserCommandDeleteTips,
) {
	u.Notification(ctx, tips.FromUserID, tips.ToUserID, constant.UserCommandDeleteNotification, tips)
}
