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

package fcm

import (
	"context"
	"path/filepath"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/OpenIMSDK/protocol/constant"
	"github.com/openimsdk/open-im-server/v3/internal/push/offlinepush"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	"github.com/redis/go-redis/v9"
	"google.golang.org/api/option"
)

const SinglePushCountLimit = 400

var Terminal = []int{constant.IOSPlatformID, constant.AndroidPlatformID, constant.WebPlatformID}

type Fcm struct {
	fcmMsgCli *messaging.Client
	cache     cache.MsgModel
}

// NewClient initializes a new FCM client using the Firebase Admin SDK.
// It requires the FCM service account credentials file located within the project's configuration directory.
func NewClient(cache cache.MsgModel) *Fcm {
	projectRoot, _ := config.GetProjectRoot()
	credentialsFilePath := filepath.Join(projectRoot, "config", config.Config.Push.Fcm.ServiceAccount)
	opt := option.WithCredentialsFile(credentialsFilePath)
	fcmApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil
	}
	ctx := context.Background()
	fcmMsgClient, err := fcmApp.Messaging(ctx)
	if err != nil {
		return nil
	}

	return &Fcm{fcmMsgCli: fcmMsgClient, cache: cache}
}

func (f *Fcm) Push(ctx context.Context, userIDs []string, title, content string, opts *offlinepush.Opts) error {
	// accounts->registrationToken
	allTokens := make(map[string][]string, 0)
	for _, account := range userIDs {
		var personTokens []string
		for _, v := range Terminal {
			Token, err := f.cache.GetFcmToken(ctx, account, v)
			if err == nil {
				personTokens = append(personTokens, Token)
			}
		}
		allTokens[account] = personTokens
	}
	Success := 0
	Fail := 0
	notification := &messaging.Notification{}
	notification.Body = content
	notification.Title = title
	var messages []*messaging.Message
	for userID, personTokens := range allTokens {
		apns := &messaging.APNSConfig{Payload: &messaging.APNSPayload{Aps: &messaging.Aps{Sound: opts.IOSPushSound}}}
		messageCount := len(messages)
		if messageCount >= SinglePushCountLimit {
			response, err := f.fcmMsgCli.SendAll(ctx, messages)
			if err != nil {
				Fail = Fail + messageCount
			} else {
				Success = Success + response.SuccessCount
				Fail = Fail + response.FailureCount
			}
			messages = messages[0:0]
		}
		if opts.IOSBadgeCount {
			unreadCountSum, err := f.cache.IncrUserBadgeUnreadCountSum(ctx, userID)
			if err == nil {
				apns.Payload.Aps.Badge = &unreadCountSum
			} else {
				// log.Error(operationID, "IncrUserBadgeUnreadCountSum redis err", err.Error(), uid)
				Fail++
				continue
			}
		} else {
			unreadCountSum, err := f.cache.GetUserBadgeUnreadCountSum(ctx, userID)
			if err == nil && unreadCountSum != 0 {
				apns.Payload.Aps.Badge = &unreadCountSum
			} else if err == redis.Nil || unreadCountSum == 0 {
				zero := 1
				apns.Payload.Aps.Badge = &zero
			} else {
				// log.Error(operationID, "GetUserBadgeUnreadCountSum redis err", err.Error(), uid)
				Fail++
				continue
			}
		}
		for _, token := range personTokens {
			temp := &messaging.Message{
				Data:         map[string]string{"ex": opts.Ex},
				Token:        token,
				Notification: notification,
				APNS:         apns,
			}
			messages = append(messages, temp)
		}
	}
	messageCount := len(messages)
	if messageCount > 0 {
		response, err := f.fcmMsgCli.SendAll(ctx, messages)
		if err != nil {
			Fail = Fail + messageCount
		} else {
			Success = Success + response.SuccessCount
			Fail = Fail + response.FailureCount
		}
	}
	return nil
}
