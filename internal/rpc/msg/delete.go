package msg

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
)

func (m *msgServer) ClearConversationsMsg(ctx context.Context, req *msg.ClearConversationsMsgReq) (*msg.ClearConversationsMsgResp, error) {
	return &msg.ClearConversationsMsgResp{}, nil
}

func (m *msgServer) UserClearAllMsg(ctx context.Context, req *msg.UserClearAllMsgReq) (*msg.UserClearAllMsgResp, error) {
	return &msg.UserClearAllMsgResp{}, nil
}

func (m *msgServer) DeleteMsgs(ctx context.Context, req *msg.DeleteMsgsReq) (*msg.DeleteMsgsResp, error) {
	return &msg.DeleteMsgsResp{}, nil
}

func (m *msgServer) DeleteMsgPhysicalBySeq(ctx context.Context, req *msg.DeleteMsgPhysicalBySeqReq) (*msg.DeleteMsgPhysicalBySeqResp, error) {
	return &msg.DeleteMsgPhysicalBySeqResp{}, nil
}
