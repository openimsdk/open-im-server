package msg

import (
	"Open_IM/pkg/common/tracelog"
	pbMsg "Open_IM/pkg/proto/msg"
	"context"
)

func (s *msgServer) SetSendMsgStatus(ctx context.Context, req *pbMsg.SetSendMsgStatusReq) (*pbMsg.SetSendMsgStatusResp, error) {
	resp := &pbMsg.SetSendMsgStatusResp{}
	if err := s.MsgInterface.SetSendMsgStatus(ctx, tracelog.GetOperationID(ctx), req.Status); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *msgServer) GetSendMsgStatus(ctx context.Context, req *pbMsg.GetSendMsgStatusReq) (resp *pbMsg.GetSendMsgStatusResp, err error) {
	resp = &pbMsg.GetSendMsgStatusResp{}
	resp.Status, err = s.MsgInterface.GetSendMsgStatus(ctx, tracelog.GetOperationID(ctx))
	if err != nil {
		return nil, err
	}
	return resp, nil
}
