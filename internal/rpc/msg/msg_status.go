package msg

import (
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/tracelog"
	pbMsg "OpenIM/pkg/proto/msg"
	"context"
)

func (m *msgServer) SetSendMsgStatus(ctx context.Context, req *pbMsg.SetSendMsgStatusReq) (*pbMsg.SetSendMsgStatusResp, error) {
	resp := &pbMsg.SetSendMsgStatusResp{}
	if err := m.MsgInterface.SetSendMsgStatus(ctx, tracelog.GetOperationID(ctx), req.Status); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *msgServer) GetSendMsgStatus(ctx context.Context, req *pbMsg.GetSendMsgStatusReq) (*pbMsg.GetSendMsgStatusResp, error) {
	resp := &pbMsg.GetSendMsgStatusResp{}
	status, err := m.MsgInterface.GetSendMsgStatus(ctx, tracelog.GetOperationID(ctx))
	if IsNotFound(err) {
		resp.Status = constant.MsgStatusNotExist
		return resp, nil
	} else if err != nil {
		return nil, err
	}
	resp.Status = status
	return resp, nil
}
