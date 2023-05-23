package msg

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
)

func (m *msgServer) RevokeMsg(ctx context.Context, req *msg.RevokeMsgReq) (*msg.RevokeMsgResp, error) {
	//msgs := []*sdkws.MsgData{}
	//if err := m.MsgDatabase.MsgToMongoMQ(ctx, "", req.ConversationID, msgs, 0); err != nil {
	//	return nil, err
	//}
	//msg := sdkws.MsgData{
	//	SendID: "",
	//}
	//
	//m.SendMsg(ctx, &msg.SendMsgReq{})

	return &msg.RevokeMsgResp{}, nil
}
