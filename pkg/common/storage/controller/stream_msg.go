package controller

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

type StreamMsgDatabase interface {
	CreateStreamMsg(ctx context.Context, model *model.StreamMsg) error
	AppendStreamMsg(ctx context.Context, clientMsgID string, startIndex int, packets []string, end bool) error
	GetStreamMsg(ctx context.Context, clientMsgID string) (*model.StreamMsg, error)
}
