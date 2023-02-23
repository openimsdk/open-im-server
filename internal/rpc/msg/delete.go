package msg

import (
	"Open_IM/pkg/common/tokenverify"
	"Open_IM/pkg/proto/msg"
	"Open_IM/pkg/proto/sdkws"
	"context"
)

func (m *msgServer) DelMsgList(ctx context.Context, req *sdkws.DelMsgListReq) (*sdkws.DelMsgListResp, error) {
	resp := &common.DelMsgListResp{}
	if _, err := m.MsgInterface.DelMsgBySeqs(ctx, req.UserID, req.SeqList); err != nil {
		return nil, err
	}
	DeleteMessageNotification(ctx, req.UserID, req.SeqList)
	return resp, nil
}

func (m *msgServer) DelSuperGroupMsg(ctx context.Context, req *msg.DelSuperGroupMsgReq) (*msg.DelSuperGroupMsgResp, error) {
	resp := &msg.DelSuperGroupMsgResp{}
	if err := tokenverify.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	//maxSeq, err := m.MsgInterface.GetGroupMaxSeq(ctx, req.GroupID)
	//if err != nil {
	//	return nil, err
	//}
	//if err := m.MsgInterface.SetGroupUserMinSeq(ctx, req.GroupID, maxSeq); err != nil {
	//	return nil, err
	//}
	if err := m.MsgInterface.DeleteUserSuperGroupMsgsAndSetMinSeq(ctx, req.GroupID, []string{req.UserID}, 0); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *msgServer) ClearMsg(ctx context.Context, req *msg.ClearMsgReq) (*msg.ClearMsgResp, error) {
	resp := &msg.ClearMsgResp{}
	if err := tokenverify.CheckAccessV3(ctx, req.UserID); err != nil {
		return nil, err
	}
	if err := m.MsgInterface.CleanUpUserMsg(ctx, req.UserID); err != nil {
		return nil, err
	}
	//if err := m.MsgInterface.DelUserAllSeq(ctx, req.UserID); err != nil {
	//	return nil, err
	//}
	return resp, nil
}
