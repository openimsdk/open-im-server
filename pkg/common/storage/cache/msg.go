package cache

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

type MsgCache interface {
	SetSendMsgStatus(ctx context.Context, id string, status int32) error
	GetSendMsgStatus(ctx context.Context, id string) (int32, error)

	GetMessageBySeqs(ctx context.Context, conversationID string, seqs []int64) ([]*model.MsgInfoModel, error)
	DelMessageBySeqs(ctx context.Context, conversationID string, seqs []int64) error
	SetMessageBySeqs(ctx context.Context, conversationID string, msgs []*model.MsgInfoModel) error
}
