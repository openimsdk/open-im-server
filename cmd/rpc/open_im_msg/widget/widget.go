package widget

import (
	"Open_IM/pkg/common/constant"
	pbChat "Open_IM/pkg/proto/chat"
	"context"
	"errors"
)

// MockBeforeSendHandler is a mock handle that handles custom logic before send msg.
func MockBeforeSendHandler(ctx context.Context, pb *pbChat.SendMsgReq) (*pbChat.SendMsgResp, bool, error) {
	if pb.MsgData.ContentType == constant.Text {
		msg := string(pb.MsgData.Content)
		if msg == "this is a m..m..mock msg" {
			return nil, false, errors.New("BANG! This msg has been banned by MockBeforeSendHandler")
		}
	}

	return nil, true, nil
}
