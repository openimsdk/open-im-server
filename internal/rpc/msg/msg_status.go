package msg

import (
	"Open_IM/pkg/common/constant"
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

func (s *msgServer) GetSendMsgStatus(ctx context.Context, req *pbMsg.GetSendMsgStatusReq) (*pbMsg.GetSendMsgStatusResp, error) {
	resp := &pbMsg.GetSendMsgStatusResp{}
	status, err := s.MsgInterface.GetSendMsgStatus(ctx, tracelog.GetOperationID(ctx))
	if IsNotFound(err) {
		resp.Status = constant.MsgStatusNotExist
		return resp, nil
	} else if err != nil {
		return nil, err
	}
	resp.Status = status
	return resp, nil
}
