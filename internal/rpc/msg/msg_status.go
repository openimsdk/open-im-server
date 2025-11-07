package msg

import (
	"context"

	"github.com/openimsdk/protocol/constant"
	pbmsg "github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/tools/mcontext"
)

func (m *msgServer) SetSendMsgStatus(ctx context.Context, req *pbmsg.SetSendMsgStatusReq) (*pbmsg.SetSendMsgStatusResp, error) {
	resp := &pbmsg.SetSendMsgStatusResp{}
	if err := m.MsgDatabase.SetSendMsgStatus(ctx, mcontext.GetOperationID(ctx), req.Status); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *msgServer) GetSendMsgStatus(ctx context.Context, req *pbmsg.GetSendMsgStatusReq) (*pbmsg.GetSendMsgStatusResp, error) {
	resp := &pbmsg.GetSendMsgStatusResp{}
	status, err := m.MsgDatabase.GetSendMsgStatus(ctx, mcontext.GetOperationID(ctx))
	if IsNotFound(err) {
		resp.Status = constant.MsgStatusNotExist
		return resp, nil
	} else if err != nil {
		return nil, err
	}
	resp.Status = status
	return resp, nil
}
