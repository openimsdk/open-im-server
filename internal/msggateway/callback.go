package msggateway

import (
	"context"
	"time"

	cbapi "github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/mcontext"
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
