package filters

import (
	rpcChat "Open_IM/internal/rpc/msg"
	"Open_IM/pkg/common/constant"
	pbChat "Open_IM/pkg/proto/chat"
	"errors"
	"fmt"
)

func MockBeforeSendFilter1(ctx *rpcChat.SendContext, pb *pbChat.SendMsgReq) (*pbChat.SendMsgResp, bool, error) {
	ctxKey := "test_key"
	v := true
	fmt.Printf("MockBeforeSendFilter1:%s set value to ctx,value is :%v\n", ctxKey, v)
	ctx.WithValue(ctxKey, v)

	return nil, true, nil
}

// MockBeforeSendFilter is a mock handle that handles custom logic before send msg.
func MockBeforeSendFilter2(ctx *rpcChat.SendContext, pb *pbChat.SendMsgReq) (*pbChat.SendMsgResp, bool, error) {
	ctxKey := "test_key"
	v, ok := ctx.Value(ctxKey).(bool)
	if ok {
		fmt.Printf("MockBeforeSendFilter2:%s selected from ctx,value is :%v\n", ctxKey, v)
	}

	fmt.Printf("MockBeforeSendHandler trigger,contentType:%d\n", pb.MsgData.GetContentType())
	if pb.MsgData.ContentType == constant.Text {
		msg := string(pb.MsgData.Content)
		fmt.Printf("text msg:%s", msg)
		if msg == "this is a m..m..mock msg" {
			fmt.Println(".==>msg had banned")
			return nil, false, errors.New("BANG! This msg has been banned by MockBeforeSendHandler")
		}
	}

	return nil, true, nil
}
