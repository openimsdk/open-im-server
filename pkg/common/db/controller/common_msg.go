package controller

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

// sourceID 可以是通知也可以是conversation
type commonMsgDatabase interface {
	BatchInsertChat2DB(ctx context.Context, sourceID string, msgs []*sdkws.MsgData, currentMaxSeq int64) error
	DeleteMessageFromCache(ctx context.Context, sourceID string, msgs []*sdkws.MsgData) error

	GetMsgBySeqs(ctx context.Context, sourceID string, seqs []int64) (seqMsg []*sdkws.MsgData, err error)
	GetMsgBySeqsRange(ctx context.Context, sourceID string, begin, end, num int64) (seqMsg []*sdkws.MsgData, err error)
	CleanUpUserMsg(ctx context.Context, sourceID string) error
	DelMsgsBySeqs(ctx context.Context, sourceID string, seqs []int64) (totalUnExistSeqs []int64, err error)
	DelMsgsAndResetMinSeq(ctx context.Context, sourceID string, userIDs []string, remainTime int64) error

	GetMinMaxSeqInMongoAndCache(ctx context.Context, sourceID string) (minSeqMongo, maxSeqMongo, minSeqCache, maxSeqCache int64, err error)

	GetMaxSeq(ctx context.Context, sourceID string) (int64, error)
	GetMinSeq(ctx context.Context, sourceID string) (int64, error)
	SetMaxSeq(ctx context.Context, sourceID string, seq int64) error
	SetMinSeq(ctx context.Context, sourceID string, seq int64) error

	MsgToMQ(ctx context.Context, sourceID string, msg2mq *sdkws.MsgData) error
	MsgToModifyMQ(ctx context.Context, sourceID string, messages []*sdkws.MsgData) error
	MsgToPushMQ(ctx context.Context, sourceID string, msg2mq *sdkws.MsgData) (int32, int64, error)
	MsgToMongoMQ(ctx context.Context, sourceID string, messages []*sdkws.MsgData, lastSeq int64) error
}
