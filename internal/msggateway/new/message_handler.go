package new

import "context"

type Req struct {
	ReqIdentifier int32  `json:"reqIdentifier" validate:"required"`
	Token         string `json:"token" `
	SendID        string `json:"sendID" validate:"required"`
	OperationID   string `json:"operationID" validate:"required"`
	MsgIncr       string `json:"msgIncr" validate:"required"`
	Data          []byte `json:"data"`
}
type MessageHandler interface {
	GetSeq(context context.Context, data Req) ([]byte, error)
	SendMessage(context context.Context, data Req) ([]byte, error)
	SendSignalMessage(context context.Context, data Req) ([]byte, error)
	PullMessageBySeqList(context context.Context, data Req) ([]byte, error)
	UserLogout(context context.Context, data Req) ([]byte, error)
	SetUserDeviceBackground(context context.Context, data Req) ([]byte, error)
}

var _ MessageHandler = (*GrpcHandler)(nil)

type GrpcHandler struct {
}

func (g GrpcHandler) GetSeq(context context.Context, data Req) ([]byte, error) {
	panic("implement me")
}

func (g GrpcHandler) SendMessage(context context.Context, data Req) ([]byte, error) {
	panic("implement me")
}

func (g GrpcHandler) SendSignalMessage(context context.Context, data Req) ([]byte, error) {
	panic("implement me")
}

func (g GrpcHandler) PullMessageBySeqList(context context.Context, data Req) ([]byte, error) {
	panic("implement me")
}

func (g GrpcHandler) UserLogout(context context.Context, data Req) ([]byte, error) {
	panic("implement me")
}

func (g GrpcHandler) SetUserDeviceBackground(context context.Context, data Req) ([]byte, error) {
	panic("implement me")
}
