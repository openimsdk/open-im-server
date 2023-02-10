package controller

import (
	pbMsg "Open_IM/pkg/proto/msg"
	"Open_IM/pkg/proto/sdkws"
	"context"
	"encoding/json"
)

type MsgInterface interface {
	BatchInsertChat2DB(ctx context.Context, userID string, msgList []*pbMsg.MsgDataToMQ, currentMaxSeq uint64) error
	BatchInsertChat2Cache(ctx context.Context, insertID string, msgList []*pbMsg.MsgDataToMQ) (error, uint64)

	DelMsgBySeqList(ctx context.Context, userID string, seqList []uint32) (totalUnExistSeqList []uint32, err error)
	// logic delete
	DelMsgLogic(ctx context.Context, userID string, seqList []uint32) error
	DelMsgBySeqListInOneDoc(ctx context.Context, docID string, seqList []uint32) (unExistSeqList []uint32, err error)
	ReplaceMsgToBlankByIndex(docID string, index int) (replaceMaxSeq uint32, err error)
	ReplaceMsgByIndex(ctx context.Context, suffixUserID string, msg *sdkws.MsgData, seqIndex int) error
	// 获取群ID或者UserID最新一条在mongo里面的消息
	GetNewestMsg(ID string) (msg *sdkws.MsgData, err error)
	// 获取群ID或者UserID最老一条在mongo里面的消息
	GetOldestMsg(ID string) (msg *sdkws.MsgData, err error)

	GetMsgBySeqListMongo2(uid string, seqList []uint32, operationID string) (seqMsg []*sdkws.MsgData, err error)
	GetSuperGroupMsgBySeqListMongo(groupID string, seqList []uint32, operationID string) (seqMsg []*sdkws.MsgData, err error)
	GetMsgAndIndexBySeqListInOneMongo2(suffixUserID string, seqList []uint32, operationID string) (seqMsg []*sdkws.MsgData, indexList []int, unExistSeqList []uint32, err error)
	SaveUserChatMongo2(uid string, sendTime int64, m *pbMsg.MsgDataToDB) error

	CleanUpUserMsgFromMongo(userID string, operationID string) error
}

func NewMsgController() MsgDatabaseInterface {
	return MsgController
}

type MsgController struct {
}

type MsgDatabaseInterface interface {
	BatchInsertChat2DB(ctx context.Context, userID string, msgList []*pbMsg.MsgDataToMQ, currentMaxSeq uint64) error
	BatchInsertChat2Cache(ctx context.Context, insertID string, msgList []*pbMsg.MsgDataToMQ) (error, uint64)

	DelMsgBySeqList(ctx context.Context, userID string, seqList []uint32) (totalUnExistSeqList []uint32, err error)
	// logic delete
	DelMsgLogic(ctx context.Context, userID string, seqList []uint32) error
	DelMsgBySeqListInOneDoc(ctx context.Context, docID string, seqList []uint32) (unExistSeqList []uint32, err error)
	ReplaceMsgToBlankByIndex(docID string, index int) (replaceMaxSeq uint32, err error)
	ReplaceMsgByIndex(ctx context.Context, suffixUserID string, msg *sdkws.MsgData, seqIndex int) error
	// 获取群ID或者UserID最新一条在mongo里面的消息
	GetNewestMsg(ID string) (msg *sdkws.MsgData, err error)
	// 获取群ID或者UserID最老一条在mongo里面的消息
	GetOldestMsg(ID string) (msg *sdkws.MsgData, err error)

	GetMsgBySeqListMongo2(uid string, seqList []uint32, operationID string) (seqMsg []*sdkws.MsgData, err error)
	GetSuperGroupMsgBySeqListMongo(groupID string, seqList []uint32, operationID string) (seqMsg []*sdkws.MsgData, err error)
	GetMsgAndIndexBySeqListInOneMongo2(suffixUserID string, seqList []uint32, operationID string) (seqMsg []*sdkws.MsgData, indexList []int, unExistSeqList []uint32, err error)
	SaveUserChatMongo2(uid string, sendTime int64, m *pbMsg.MsgDataToDB) error
	// 删除用户所有消息/redis/mongo然后重置seq
	CleanUpUserMsgFromMongo(userID string, operationID string) error
}

func NewMsgDatabase() MsgDatabaseInterface {
	return MsgDatabase
}

type MsgDatabase struct {
}

func (m *MsgDatabase) BatchInsertChat2DB(ctx context.Context, userID string, msgList []*pbMsg.MsgDataToMQ, currentMaxSeq uint64) error {

}

func (m *MsgDatabase) CleanUpUserMsgFromMongo(userID string, operationID string) error {

}
