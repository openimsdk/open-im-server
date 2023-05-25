package msg

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

type MessageInterceptorFunc func(ctx context.Context, req *msg.SendMsgReq) (*sdkws.MsgData, error)

func MessageHasReadEnabled(_ context.Context, req *msg.SendMsgReq) (*sdkws.MsgData, error) {
	switch req.MsgData.ContentType {
	case constant.HasReadReceipt:
		if config.Config.SingleMessageHasReadReceiptEnable {
			return req.MsgData, nil
		} else {
			return nil, errs.ErrMessageHasReadDisable.Wrap()
		}
	case constant.GroupHasReadReceipt:
		if config.Config.GroupMessageHasReadReceiptEnable {
			return req.MsgData, nil
		} else {
			return nil, errs.ErrMessageHasReadDisable.Wrap()
		}
	}
	return req.MsgData, nil
}
func MessageModifyCallback(ctx context.Context, req *msg.SendMsgReq) (*sdkws.MsgData, error) {
	if err := callbackMsgModify(ctx, req); err != nil && err != errs.ErrCallbackContinue {
		log.ZWarn(ctx, "CallbackMsgModify failed", err, "req", req.String())
		return nil, err
	}
	return req.MsgData, nil
}
func MessageBeforeSendCallback(ctx context.Context, req *msg.SendMsgReq) (*sdkws.MsgData, error) {
	switch req.MsgData.SessionType {
	case constant.SingleChatType:
		if err := callbackBeforeSendSingleMsg(ctx, req); err != nil && err != errs.ErrCallbackContinue {
			log.ZWarn(ctx, "CallbackBeforeSendSingleMsg failed", err, "req", req.String())
			return nil, err
		}
	case constant.NotificationChatType:
	case constant.SuperGroupChatType:
		if err := callbackBeforeSendGroupMsg(ctx, req); err != nil && err != errs.ErrCallbackContinue {
			log.ZWarn(ctx, "CallbackBeforeSendGroupMsg failed", err, "req", req.String())
			return nil, err
		}
	default:
		return nil, errs.ErrArgs.Wrap("unknown sessionType")
	}
	return req.MsgData, nil
}
