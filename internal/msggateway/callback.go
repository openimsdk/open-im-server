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

package msggateway

import (
	"context"
	"time"

	"github.com/OpenIMSDK/protocol/constant"
	cbapi "github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/http"
	"github.com/openimsdk/tools/mcontext"
)

func CallbackUserOnline(ctx context.Context, callback *config.Callback, userID string, platformID int, isAppBackground bool, connID string) error {
	if !callback.CallbackUserOnline.Enable {
		return nil
	}
	req := cbapi.CallbackUserOnlineReq{
		UserStatusCallbackReq: cbapi.UserStatusCallbackReq{
			UserStatusBaseCallback: cbapi.UserStatusBaseCallback{
				CallbackCommand: cbapi.CallbackUserOnlineCommand,
				OperationID:     mcontext.GetOperationID(ctx),
				PlatformID:      platformID,
				Platform:        constant.PlatformIDToName(platformID),
			},
			UserID: userID,
		},
		Seq:             time.Now().UnixMilli(),
		IsAppBackground: isAppBackground,
		ConnID:          connID,
	}
	resp := cbapi.CommonCallbackResp{}
	if err := http.CallBackPostReturn(ctx, callback.CallbackUrl, &req, &resp, callback.CallbackUserOnline); err != nil {
		return err
	}
	return nil
}

func CallbackUserOffline(ctx context.Context, callback *config.Callback, userID string, platformID int, connID string) error {
	if !callback.CallbackUserOffline.Enable {
		return nil
	}
	req := &cbapi.CallbackUserOfflineReq{
		UserStatusCallbackReq: cbapi.UserStatusCallbackReq{
			UserStatusBaseCallback: cbapi.UserStatusBaseCallback{
				CallbackCommand: cbapi.CallbackUserOfflineCommand,
				OperationID:     mcontext.GetOperationID(ctx),
				PlatformID:      platformID,
				Platform:        constant.PlatformIDToName(platformID),
			},
			UserID: userID,
		},
		Seq:    time.Now().UnixMilli(),
		ConnID: connID,
	}
	resp := &cbapi.CallbackUserOfflineResp{}
	if err := http.CallBackPostReturn(ctx, callback.CallbackUrl, req, resp, callback.CallbackUserOffline); err != nil {
		return err
	}
	return nil
}

func CallbackUserKickOff(ctx context.Context, callback *config.Callback, userID string, platformID int) error {
	if !callback.CallbackUserKickOff.Enable {
		return nil
	}
	req := &cbapi.CallbackUserKickOffReq{
		UserStatusCallbackReq: cbapi.UserStatusCallbackReq{
			UserStatusBaseCallback: cbapi.UserStatusBaseCallback{
				CallbackCommand: cbapi.CallbackUserKickOffCommand,
				OperationID:     mcontext.GetOperationID(ctx),
				PlatformID:      platformID,
				Platform:        constant.PlatformIDToName(platformID),
			},
			UserID: userID,
		},
		Seq: time.Now().UnixMilli(),
	}
	resp := &cbapi.CommonCallbackResp{}
	if err := http.CallBackPostReturn(ctx, callback.CallbackUrl, req, resp, callback.CallbackUserOffline); err != nil {
		return err
	}
	return nil
}
