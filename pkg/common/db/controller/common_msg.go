package controller

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

// conversationID 可以是通知也可以是conversation
type commonMsgDatabase interface {
	BatchInsertChat2DB(ctx context.Context, conversationID string, msgs []*sdkws.MsgData, currentMaxSeq int64) error
	DeleteMessageFromCache(ctx context.Context, conversationID string, msgs []*sdkws.MsgData) error

	GetMsgBySeqs(ctx context.Context, conversationID string, seqs []int64) (seqMsg []*sdkws.MsgData, err error)
	GetMsgBySeqsRange(ctx context.Context, conversationID string, begin, end, num int64) (seqMsg []*sdkws.MsgData, err error)
	CleanUpUserMsg(ctx context.Context, conversationID string) error
	DelMsgsBySeqs(ctx context.Context, conversationID string, seqs []int64) (totalUnExistSeqs []int64, err error)
	DelMsgsAndResetMinSeq(ctx context.Context, conversationID string, userIDs []string, remainTime int64) error

	GetMinMaxSeqInMongoAndCache(ctx context.Context, conversationID string) (minSeqMongo, maxSeqMongo, minSeqCache, maxSeqCache int64, err error)

	GetMaxSeq(ctx context.Context, conversationID string) (int64, error)
	GetMinSeq(ctx context.Context, conversationID string) (int64, error)
	SetMaxSeq(ctx context.Context, conversationID string, seq int64) error
	SetMinSeq(ctx context.Context, conversationID string, seq int64) error
	GetMaxSeqs(ctx context.Context, conversationIDs []string) (map[string]int64, error)
	GetMinSeqs(ctx context.Context, conversationIDs []string) (map[string]int64, error)

	MsgToMQ(ctx context.Context, conversationID string, msg2mq *sdkws.MsgData) error
	MsgToModifyMQ(ctx context.Context, conversationID string, messages []*sdkws.MsgData) error
	MsgToPushMQ(ctx context.Context, conversationID string, msg2mq *sdkws.MsgData) (int32, int64, error)
	MsgToMongoMQ(ctx context.Context, conversationID string, messages []*sdkws.MsgData, lastSeq int64) error
}
