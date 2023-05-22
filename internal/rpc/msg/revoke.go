package msg

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
)

func (m *msgServer) RevokeMsg(ctx context.Context, req *msg.RevokeMsgReq) (*msg.RevokeMsgResp, error) {
	// if _, err := m.MsgDatabase.DelMsgBySeqs(ctx, req.UserID, req.Seqs); err != nil {
	// 	return nil, err
	// }
	return &msg.RevokeMsgResp{}, nil
}
