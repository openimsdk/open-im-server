package controller

import (
	"Open_IM/pkg/proto/msg"
	pbMsg "Open_IM/pkg/proto/msg"
	"context"
)

type MsgInterface interface {
	//消息写入队列
	MsgToMQ(ctx context.Context, key string, m *msg.MsgDataToMQ) error

	BatchInsertChat2DB(ctx context.Context, userID string, msgList []*pbMsg.MsgDataToMQ, currentMaxSeq uint64) error
	BatchInsertChat2Cache(ctx context.Context, insertID string, msgList []*pbMsg.MsgDataToMQ) (error, uint64)

	DelMsgBySeqList(ctx context.Context, userID string, seqList []uint32) (totalUnExistSeqList []uint32, err error)
	// logic delete
	DelMsgLogic(ctx context.Context, userID string, seqList []uint32) error
	DelMsgBySeqListInOneDoc(ctx context.Context, docID string, seqList []uint32) (unExistSeqList []uint32, err error)
	ReplaceMsgToBlankByIndex(docID string, index int) (replaceMaxSeq uint32, err error)
}

type MsgDatabaseInterface interface {
	BatchInsertChat2DB(ctx context.Context, userID string, msgList []*pbMsg.MsgDataToMQ, currentMaxSeq uint64) error
	BatchInsertChat2Cache(ctx context.Context, insertID string, msgList []*pbMsg.MsgDataToMQ) (error, uint64)
	DelMsgBySeqList(ctx context.Context, userID string, seqList []uint32) (totalUnExistSeqList []uint32, err error)
}
