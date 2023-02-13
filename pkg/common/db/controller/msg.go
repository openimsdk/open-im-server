package controller

import (
	"Open_IM/pkg/proto/msg"
	pbMsg "Open_IM/pkg/proto/msg"
	"Open_IM/pkg/proto/sdkws"
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

	// status
	SetSendMsgStatus(ctx context.Context, id string, status int32) error
	GetSendMsgStatus(ctx context.Context, id string) (int32, error) // 不存在返回 constant.MsgStatusNotExist
	// delete
	DelMsgFromCache(ctx context.Context, userID string, seqs []uint32) error
	GetGroupMaxSeq(ctx context.Context, groupID string) (uint32, error)
	GetGroupMinSeq(ctx context.Context, groupID string) (uint32, error)
	SetGroupUserMinSeq(ctx context.Context, groupID string, seq uint32) error
	DelUserAllSeq(ctx context.Context, userID string) error // redis and mongodb
	GetUserMaxSeq(ctx context.Context, userID string) (uint32, error)
	GetUserMinSeq(ctx context.Context, userID string) (uint32, error)

	GetMessageListBySeq(ctx context.Context, userID string, seqs []uint32) ([]*sdkws.MsgData, error)
	GetSuperGroupMsg(ctx context.Context, groupID string, seq uint32) (*sdkws.MsgData, error)
}

type MsgDatabaseInterface interface {
	BatchInsertChat2DB(ctx context.Context, userID string, msgList []*pbMsg.MsgDataToMQ, currentMaxSeq uint64) error
	BatchInsertChat2Cache(ctx context.Context, insertID string, msgList []*pbMsg.MsgDataToMQ) (error, uint64)
	DelMsgBySeqList(ctx context.Context, userID string, seqList []uint32) (totalUnExistSeqList []uint32, err error)
}
