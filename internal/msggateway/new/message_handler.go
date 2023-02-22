package new

import (
	"Open_IM/internal/common/notification"
	"Open_IM/pkg/proto/msg"
	pbRtc "Open_IM/pkg/proto/rtc"
	"Open_IM/pkg/proto/sdkws"
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/golang/protobuf/proto"
)

type Req struct {
	ReqIdentifier int32  `json:"reqIdentifier" validate:"required"`
	Token         string `json:"token" `
	SendID        string `json:"sendID" validate:"required"`
	OperationID   string `json:"operationID" validate:"required"`
	MsgIncr       string `json:"msgIncr" validate:"required"`
	Data          []byte `json:"data"`
}
type Resp struct {
	ReqIdentifier int32  `json:"reqIdentifier"`
	MsgIncr       string `json:"msgIncr"`
	OperationID   string `json:"operationID"`
	ErrCode       int32  `json:"errCode"`
	ErrMsg        string `json:"errMsg"`
	Data          []byte `json:"data"`
}
type MessageHandler interface {
	GetSeq(context context.Context, data Req) ([]byte, error)
	SendMessage(context context.Context, data Req) ([]byte, error)
	SendSignalMessage(context context.Context, data Req) ([]byte, error)
	PullMessageBySeqList(context context.Context, data Req) ([]byte, error)
	UserLogout(context context.Context, data Req) ([]byte, error)
	SetUserDeviceBackground(context context.Context, data Req) ([]byte, bool, error)
}

var _ MessageHandler = (*GrpcHandler)(nil)

type GrpcHandler struct {
	notification *notification.Check
	validate     *validator.Validate
}

func NewGrpcHandler(validate *validator.Validate, notification *notification.Check) *GrpcHandler {
	return &GrpcHandler{notification: notification, validate: validate}
}

func (g GrpcHandler) GetSeq(context context.Context, data Req) ([]byte, error) {
	req := sdkws.GetMaxAndMinSeqReq{}
	if err := proto.Unmarshal(data.Data, &req); err != nil {
		return nil, err
	}
	if err := g.validate.Struct(req); err != nil {
		return nil, err
	}
	resp, err := g.notification.Msg.GetMaxAndMinSeq(context, &req)
	if err != nil {
		return nil, err
	}
	c, err := proto.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (g GrpcHandler) SendMessage(context context.Context, data Req) ([]byte, error) {
	msgData := sdkws.MsgData{}
	if err := proto.Unmarshal(data.Data, &msgData); err != nil {
		return nil, err
	}
	if err := g.validate.Struct(msgData); err != nil {
		return nil, err
	}
	req := msg.SendMsgReq{MsgData: &msgData}
	resp, err := g.notification.Msg.SendMsg(context, &req)
	if err != nil {
		return nil, err
	}
	c, err := proto.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (g GrpcHandler) SendSignalMessage(context context.Context, data Req) ([]byte, error) {
	signalReq := pbRtc.SignalReq{}
	if err := proto.Unmarshal(data.Data, &signalReq); err != nil {
		return nil, err
	}
	if err := g.validate.Struct(signalReq); err != nil {
		return nil, err
	}
	//req := pbRtc.SignalMessageAssembleReq{SignalReq: &signalReq, OperationID: "111"}
	//todo rtc rpc call
	resp, err := g.notification.Msg.SendMsg(context, nil)
	if err != nil {
		return nil, err
	}
	c, err := proto.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (g GrpcHandler) PullMessageBySeqList(context context.Context, data Req) ([]byte, error) {
	req := sdkws.PullMessageBySeqListReq{}
	if err := proto.Unmarshal(data.Data, &req); err != nil {
		return nil, err
	}
	if err := g.validate.Struct(data); err != nil {
		return nil, err
	}
	resp, err := g.notification.Msg.PullMessageBySeqList(context, &req)
	if err != nil {
		return nil, err
	}
	c, err := proto.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (g GrpcHandler) UserLogout(context context.Context, data Req) ([]byte, error) {
	//todo
	resp, err := g.notification.Msg.PullMessageBySeqList(context, nil)
	if err != nil {
		return nil, err
	}
	c, err := proto.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return c, nil
}
func (g GrpcHandler) SetUserDeviceBackground(_ context.Context, data Req) ([]byte, bool, error) {
	req := sdkws.SetAppBackgroundStatusReq{}
	if err := proto.Unmarshal(data.Data, &req); err != nil {
		return nil, false, err
	}
	if err := g.validate.Struct(data); err != nil {
		return nil, false, err
	}
	return nil, req.IsBackground, nil
}
