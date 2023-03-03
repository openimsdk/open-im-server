package msg

import (
	"OpenIM/pkg/common/tokenverify"
	"OpenIM/pkg/proto/msg"
	"context"
)

func (m *msgServer) DelMsgs(ctx context.Context, req *msg.DelMsgsReq) (*msg.DelMsgsResp, error) {
	resp := &msg.DelMsgsResp{}
	if _, err := m.MsgDatabase.DelMsgBySeqs(ctx, req.UserID, req.Seqs); err != nil {
		return nil, err
	}
	//DeleteMessageNotification(ctx, req.UserID, req.Seqs)
	return resp, nil
}

func (m *msgServer) DelSuperGroupMsg(ctx context.Context, req *msg.DelSuperGroupMsgReq) (*msg.DelSuperGroupMsgResp, error) {
	resp := &msg.DelSuperGroupMsgResp{}
	if err := tokenverify.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if err := m.MsgDatabase.DeleteUserSuperGroupMsgsAndSetMinSeq(ctx, req.GroupID, []string{req.UserID}, 0); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *msgServer) ClearMsg(ctx context.Context, req *msg.ClearMsgReq) (*msg.ClearMsgResp, error) {
	resp := &msg.ClearMsgResp{}
	if err := tokenverify.CheckAccessV3(ctx, req.UserID); err != nil {
		return nil, err
	}
	if err := m.MsgDatabase.CleanUpUserMsg(ctx, req.UserID); err != nil {
		return nil, err
	}
	return resp, nil
}
