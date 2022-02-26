package msg

import pbChat "Open_IM/pkg/proto/chat"

// BeforeSendFilter handles custom logic before send msg.
type BeforeSendFilter func(ctx *SendContext, pb *pbChat.SendMsgReq) (*pbChat.SendMsgResp, bool, error)
