package msg

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/tokenverify"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
)

func (m *msgServer) DelMsgs(ctx context.Context, req *msg.DelMsgsReq) (*msg.DelMsgsResp, error) {
	resp := &msg.DelMsgsResp{}
	if _, err := m.MsgDatabase.DelMsgBySeqs(ctx, req.UserID, req.Seqs); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *msgServer) DelSuperGroupMsg(ctx context.Context, req *msg.DelSuperGroupMsgReq) (*msg.DelSuperGroupMsgResp, error) {
	resp := &msg.DelSuperGroupMsgResp{}
	if err := tokenverify.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if err := m.MsgDatabase.DeleteConversationMsgsAndSetMinSeq(ctx, req.GroupID, 0); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *msgServer) ClearMsg(ctx context.Context, req *msg.ClearMsgReq) (*msg.ClearMsgResp, error) {
	resp := &msg.ClearMsgResp{}
	if err := tokenverify.CheckAccessV3(ctx, req.UserID); err != nil {
		return nil, err
	}
	if err := m.MsgDatabase.CleanUpConversationMsgs(ctx, req.UserID); err != nil {
		return nil, err
	}
	return resp, nil
}
