package msg

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

type MessageInterceptorFunc func(ctx context.Context, req *msg.SendMsgReq) (*sdkws.MsgData, error)

func MessageHasReadEnabled(_ context.Context, req *msg.SendMsgReq) (*sdkws.MsgData, error) {
	switch {
	case req.MsgData.ContentType == constant.HasReadReceipt && req.MsgData.SessionType == constant.SingleChatType:
		if config.Config.SingleMessageHasReadReceiptEnable {
			return req.MsgData, nil
		} else {
			return nil, errs.ErrMessageHasReadDisable.Wrap()
		}
	case req.MsgData.ContentType == constant.HasReadReceipt && req.MsgData.SessionType == constant.SuperGroupChatType:
		if config.Config.GroupMessageHasReadReceiptEnable {
			return req.MsgData, nil
		} else {
			return nil, errs.ErrMessageHasReadDisable.Wrap()
		}
	}
	return req.MsgData, nil
}
