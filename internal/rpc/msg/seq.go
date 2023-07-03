package msg

import (
	"context"

	pbMsg "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
)

func (m *msgServer) GetConversationMaxSeq(
	ctx context.Context,
	req *pbMsg.GetConversationMaxSeqReq,
) (resp *pbMsg.GetConversationMaxSeqResp, err error) {
	maxSeq, err := m.MsgDatabase.GetMaxSeq(ctx, req.ConversationID)
	if err != nil {
		return nil, err
	}
	return &pbMsg.GetConversationMaxSeqResp{MaxSeq: maxSeq}, nil
}
