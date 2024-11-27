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

	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/mcontext"

	cbapi "github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
)

func (ws *WsServer) webhookAfterUserOnline(ctx context.Context, after *config.AfterConfig, userID string, platformID int, isAppBackground bool, connID string) {
	req := cbapi.CallbackUserOnlineReq{
		UserStatusCallbackReq: cbapi.UserStatusCallbackReq{
			UserStatusBaseCallback: cbapi.UserStatusBaseCallback{
				CallbackCommand: cbapi.CallbackAfterUserOnlineCommand,
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
	ws.webhookClient.AsyncPost(ctx, req.GetCallbackCommand(), req, &cbapi.CommonCallbackResp{}, after)
}

func (ws *WsServer) webhookAfterUserOffline(ctx context.Context, after *config.AfterConfig, userID string, platformID int, connID string) {
	req := &cbapi.CallbackUserOfflineReq{
		UserStatusCallbackReq: cbapi.UserStatusCallbackReq{
			UserStatusBaseCallback: cbapi.UserStatusBaseCallback{
				CallbackCommand: cbapi.CallbackAfterUserOfflineCommand,
				OperationID:     mcontext.GetOperationID(ctx),
				PlatformID:      platformID,
				Platform:        constant.PlatformIDToName(platformID),
			},
			UserID: userID,
		},
		Seq:    time.Now().UnixMilli(),
		ConnID: connID,
	}
	ws.webhookClient.AsyncPost(ctx, req.GetCallbackCommand(), req, &cbapi.CallbackUserOfflineResp{}, after)
}

func (ws *WsServer) webhookAfterUserKickOff(ctx context.Context, after *config.AfterConfig, userID string, platformID int) {
	req := &cbapi.CallbackUserKickOffReq{
		UserStatusCallbackReq: cbapi.UserStatusCallbackReq{
			UserStatusBaseCallback: cbapi.UserStatusBaseCallback{
				CallbackCommand: cbapi.CallbackAfterUserKickOffCommand,
				OperationID:     mcontext.GetOperationID(ctx),
				PlatformID:      platformID,
				Platform:        constant.PlatformIDToName(platformID),
			},
			UserID: userID,
		},
		Seq: time.Now().UnixMilli(),
	}
	ws.webhookClient.AsyncPost(ctx, req.GetCallbackCommand(), req, &cbapi.CommonCallbackResp{}, after)
}
