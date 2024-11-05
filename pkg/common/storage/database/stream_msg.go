package database

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"time"
)

type StreamMsg interface {
	CreateStreamMsg(ctx context.Context, model *model.StreamMsg) error
	AppendStreamMsg(ctx context.Context, clientMsgID string, startIndex int, packets []string, end bool, deadlineTime time.Time) error
	GetStreamMsg(ctx context.Context, clientMsgID string) (*model.StreamMsg, error)
}
