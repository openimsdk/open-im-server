package controller

import (
	"context"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

type StreamMsgDatabase interface {
	CreateStreamMsg(ctx context.Context, model *model.StreamMsg) error
	AppendStreamMsg(ctx context.Context, clientMsgID string, startIndex int, packets []string, end bool, deadlineTime time.Time) error
	GetStreamMsg(ctx context.Context, clientMsgID string) (*model.StreamMsg, error)
}

func NewStreamMsgDatabase(db database.StreamMsg) StreamMsgDatabase {
	return &streamMsgDatabase{db: db}
}

type streamMsgDatabase struct {
	db database.StreamMsg
}

func (m *streamMsgDatabase) CreateStreamMsg(ctx context.Context, model *model.StreamMsg) error {
	return m.db.CreateStreamMsg(ctx, model)
}

func (m *streamMsgDatabase) AppendStreamMsg(ctx context.Context, clientMsgID string, startIndex int, packets []string, end bool, deadlineTime time.Time) error {
	return m.db.AppendStreamMsg(ctx, clientMsgID, startIndex, packets, end, deadlineTime)
}

func (m *streamMsgDatabase) GetStreamMsg(ctx context.Context, clientMsgID string) (*model.StreamMsg, error) {
	return m.db.GetStreamMsg(ctx, clientMsgID)
}
