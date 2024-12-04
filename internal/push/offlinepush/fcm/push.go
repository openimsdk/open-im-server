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
	"fmt"
	"path/filepath"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/redis/go-redis/v9"
	"google.golang.org/api/option"

	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/httputil"

	"github.com/openimsdk/open-im-server/v3/internal/push/offlinepush/options"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
)

const SinglePushCountLimit = 400

var Terminal = []int{constant.IOSPlatformID, constant.AndroidPlatformID, constant.WebPlatformID}

type Fcm struct {
	fcmMsgCli *messaging.Client
	cache     cache.ThirdCache
}

// NewClient initializes a new FCM client using the Firebase Admin SDK.
// It requires the FCM service account credentials file located within the project's configuration directory.
func NewClient(pushConf *config.Push, cache cache.ThirdCache, fcmConfigPath string) (*Fcm, error) {
	var opt option.ClientOption
	switch {
	case len(pushConf.FCM.FilePath) != 0:
		// with file path
		credentialsFilePath := filepath.Join(fcmConfigPath, pushConf.FCM.FilePath)
		opt = option.WithCredentialsFile(credentialsFilePath)
	case len(pushConf.FCM.AuthURL) != 0:
		// with authentication URL
		client := httputil.NewHTTPClient(httputil.NewClientConfig())
		resp, err := client.Get(pushConf.FCM.AuthURL)
		if err != nil {
			return nil, err
		}
		opt = option.WithCredentialsJSON(resp)
	default:
		return nil, errs.New("no FCM config").Wrap()
	}

	fcmApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	ctx := context.Background()
	fcmMsgClient, err := fcmApp.Messaging(ctx)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &Fcm{fcmMsgCli: fcmMsgClient, cache: cache}, nil
}

func (f *Fcm) Push(ctx context.Context, userIDs []string, title, content string, opts *options.Opts) error {
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
	var sendErrBuilder strings.Builder
	var msgErrBuilder strings.Builder
	for userID, personTokens := range allTokens {
		apns := &messaging.APNSConfig{Payload: &messaging.APNSPayload{Aps: &messaging.Aps{Sound: opts.IOSPushSound}}}
		messageCount := len(messages)
		if messageCount >= SinglePushCountLimit {
			response, err := f.fcmMsgCli.SendEach(ctx, messages)
			if err != nil {
				Fail = Fail + messageCount
				// Record push error
				sendErrBuilder.WriteString(err.Error())
				sendErrBuilder.WriteByte('.')
			} else {
				Success = Success + response.SuccessCount
				Fail = Fail + response.FailureCount
				if response.FailureCount != 0 {
					// Record message error
					for i := range response.Responses {
						if !response.Responses[i].Success {
							msgErrBuilder.WriteString(response.Responses[i].Error.Error())
							msgErrBuilder.WriteByte('.')
						}
					}
				}
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
		response, err := f.fcmMsgCli.SendEach(ctx, messages)
		if err != nil {
			Fail = Fail + messageCount
		} else {
			Success = Success + response.SuccessCount
			Fail = Fail + response.FailureCount
		}
	}
	if Fail != 0 {
		return errs.New(fmt.Sprintf("%d message send failed;send err:%s;message err:%s",
			Fail, sendErrBuilder.String(), msgErrBuilder.String())).Wrap()
	}
	return nil
}
