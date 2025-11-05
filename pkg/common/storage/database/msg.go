package database

import (
	"context"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/protocol/msg"
	"go.mongodb.org/mongo-driver/mongo"
)

type Msg interface {
	Create(ctx context.Context, model *model.MsgDocModel) error
	UpdateMsg(ctx context.Context, docID string, index int64, key string, value any) (*mongo.UpdateResult, error)
	PushUnique(ctx context.Context, docID string, index int64, key string, value any) (*mongo.UpdateResult, error)
	FindOneByDocID(ctx context.Context, docID string) (*model.MsgDocModel, error)
	GetMsgBySeqIndexIn1Doc(ctx context.Context, userID, docID string, seqs []int64) ([]*model.MsgInfoModel, error)
	GetNewestMsg(ctx context.Context, conversationID string) (*model.MsgInfoModel, error)
	GetOldestMsg(ctx context.Context, conversationID string) (*model.MsgInfoModel, error)
	DeleteMsgsInOneDocByIndex(ctx context.Context, docID string, indexes []int) error
	MarkSingleChatMsgsAsRead(ctx context.Context, userID string, docID string, indexes []int64) error
	SearchMessage(ctx context.Context, req *msg.SearchMessageReq) (int64, []*model.MsgInfoModel, error)
	RangeUserSendCount(ctx context.Context, start time.Time, end time.Time, group bool, ase bool, pageNumber int32, showNumber int32) (msgCount int64, userCount int64, users []*model.UserCount, dateCount map[string]int64, err error)
	RangeGroupSendCount(ctx context.Context, start time.Time, end time.Time, ase bool, pageNumber int32, showNumber int32) (msgCount int64, userCount int64, groups []*model.GroupCount, dateCount map[string]int64, err error)
	DeleteDoc(ctx context.Context, docID string) error
	GetRandBeforeMsg(ctx context.Context, ts int64, limit int) ([]*model.MsgDocModel, error)
	GetLastMessageSeqByTime(ctx context.Context, conversationID string, time int64) (int64, error)
	GetLastMessage(ctx context.Context, conversationID string) (*model.MsgInfoModel, error)
	FindSeqs(ctx context.Context, conversationID string, seqs []int64) ([]*model.MsgInfoModel, error)
}
