package widget

import (
	"context"

	pbChat "Open_IM/pkg/proto/chat"
)

// BeforeSendHandler handles custom logic before send msg.
type BeforeSendHandler func(ctx context.Context, pb *pbChat.SendMsgReq) (*pbChat.SendMsgResp, bool, error)
