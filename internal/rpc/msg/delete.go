package msg

import (
	"Open_IM/pkg/common/tokenverify"
	"Open_IM/pkg/proto/msg"
	common "Open_IM/pkg/proto/sdkws"
	"context"
)

func (s *msgServer) DelMsgList(ctx context.Context, req *common.DelMsgListReq) (*common.DelMsgListResp, error) {
	resp := &common.DelMsgListResp{}
	if err := s.MsgInterface.DelMsgFromCache(ctx, req.UserID, req.SeqList); err != nil {
		return nil, err
	}
	DeleteMessageNotification(ctx, req.UserID, req.SeqList)
	return resp, nil
}

func (s *msgServer) DelSuperGroupMsg(ctx context.Context, req *msg.DelSuperGroupMsgReq) (*msg.DelSuperGroupMsgResp, error) {
	resp := &msg.DelSuperGroupMsgResp{}
	if err := tokenverify.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	maxSeq, err := s.MsgInterface.GetGroupMaxSeq(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if err := s.MsgInterface.SetGroupUserMinSeq(ctx, req.GroupID, maxSeq); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *msgServer) ClearMsg(ctx context.Context, req *msg.ClearMsgReq) (*msg.ClearMsgResp, error) {
	resp := &msg.ClearMsgResp{}
	if err := tokenverify.CheckAccessV3(ctx, req.UserID); err != nil {
		return nil, err
	}
	if err := s.MsgInterface.DelUserAllSeq(ctx, req.UserID); err != nil {
		return nil, err
	}
	return resp, nil
}
