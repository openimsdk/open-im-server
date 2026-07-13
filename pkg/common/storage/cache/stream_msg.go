package cache

import "context"

type StreamMsg struct {
	SendUserID    string
	RecvID        string
	SessionType   int32
	StreamType    string
	StreamContent string
	Packets       map[int64]string
	End           bool
	UpdateTime    int64
}

type StreamMsgCache interface {
	CreateStreamMsg(ctx context.Context, conversationID string, clientMsgID string, msg *StreamMsg) error
	AppendStreamMsg(ctx context.Context, conversationID string, clientMsgID string, startIndex int, packets []string, end bool, retPacket bool) (*StreamMsg, error)
	GetStreamMsg(ctx context.Context, conversationID string, clientMsgID string) (*StreamMsg, error)
	GetStreamMsgEnd(ctx context.Context, conversationID string, clientMsgID string) (bool, error)
}
