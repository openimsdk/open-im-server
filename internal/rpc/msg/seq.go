package msg

import (
	"context"

	pbMsg "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
)

func (m *msgServer) GetConversationMaxSeq(ctx context.Context, req *pbMsg.GetConversationMaxSeqReq) (resp *pbMsg.GetConversationMaxSeqResp, err error) {
	resp = &pbMsg.GetConversationMaxSeqResp{}
	resp.MaxSeq, err = m.MsgDatabase.GetMaxSeq(ctx, req.ConversationID)
	return resp, err
}
