package unrelation

import (
	"context"
	"strconv"
	"strings"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

const (
	singleGocMsgNum = 5000
	Msg             = "msg"
	OldestList      = 0
	NewestList      = -1
)

type MsgDocModel struct {
	DocID string         `bson:"doc_id"`
	Msg   []MsgInfoModel `bson:"msgs"`
}

type MsgInfoModel struct {
	SendTime int64  `bson:"sendtime"`
	Msg      []byte `bson:"msg"`
}

type MsgDocModelInterface interface {
	PushMsgsToDoc(ctx context.Context, docID string, msgsToMongo []MsgInfoModel) error
	Create(ctx context.Context, model *MsgDocModel) error
	UpdateMsgStatusByIndexInOneDoc(ctx context.Context, docID string, msg *sdkws.MsgData, seqIndex int, status int32) error
	FindOneByDocID(ctx context.Context, docID string) (*MsgDocModel, error)
	GetMsgBySeqIndexIn1Doc(ctx context.Context, docID string, beginSeq, endSeq int64) ([]*sdkws.MsgData, []int64, error)
	GetNewestMsg(ctx context.Context, conversationID string) (*MsgInfoModel, error)
	GetOldestMsg(ctx context.Context, conversationID string) (*MsgInfoModel, error)
	Delete(ctx context.Context, docIDs []string) error
	GetMsgsByIndex(ctx context.Context, conversationID string, index int64) (*MsgDocModel, error)
	UpdateOneDoc(ctx context.Context, msg *MsgDocModel) error
}

func (MsgDocModel) TableName() string {
	return Msg
}

func (MsgDocModel) GetSingleGocMsgNum() int64 {
	return singleGocMsgNum
}

func (m *MsgDocModel) IsFull() bool {
	index, _ := strconv.Atoi(strings.Split(m.DocID, ":")[1])
	if index == 0 {
		if len(m.Msg) >= singleGocMsgNum-1 {
			return true
		}
	}
	if len(m.Msg) >= singleGocMsgNum {
		return true
	}

	return false
}

func (m MsgDocModel) GetDocID(conversationID string, seq int64) string {
	seqSuffix := seq / singleGocMsgNum
	return m.indexGen(conversationID, seqSuffix)
}

func (m MsgDocModel) GetSeqDocIDList(userID string, maxSeq int64) []string {
	seqMaxSuffix := maxSeq / singleGocMsgNum
	var seqUserIDs []string
	for i := 0; i <= int(seqMaxSuffix); i++ {
		seqUserID := m.indexGen(userID, int64(i))
		seqUserIDs = append(seqUserIDs, seqUserID)
	}
	return seqUserIDs
}

// func (m MsgDocModel) getSeqSuperGroupID(groupID string, seq int64) string {
// 	seqSuffix := seq / singleGocMsgNum
// 	return m.superGroupIndexGen(groupID, seqSuffix)
// }

// func (m MsgDocModel) superGroupIndexGen(groupID string, seqSuffix int64) string {
// 	return "super_group_" + groupID + ":" + strconv.FormatInt(int64(seqSuffix), 10)
// }

func (m MsgDocModel) GetDocIDSeqsMap(conversationID string, seqs []int64) map[string][]int64 {
	t := make(map[string][]int64)
	for i := 0; i < len(seqs); i++ {
		docID := m.GetDocID(conversationID, seqs[i])
		if value, ok := t[docID]; !ok {
			var temp []int64
			t[docID] = append(temp, seqs[i])
		} else {
			t[docID] = append(value, seqs[i])
		}
	}
	return t
}

func (m MsgDocModel) GetSeqsBeginEnd(seqs []int64) (int64, int64) {
	if len(seqs) == 0 {
		return 0, 0
	}
	return seqs[0], seqs[len(seqs)-1]
}

func (m MsgDocModel) GetMsgIndex(seq int64) int64 {
	seqSuffix := seq / singleGocMsgNum
	var index int64
	if seqSuffix == 0 {
		index = (seq - seqSuffix*singleGocMsgNum) - 1
	} else {
		index = seq - seqSuffix*singleGocMsgNum
	}
	return index
}

func (m MsgDocModel) indexGen(conversationID string, seqSuffix int64) string {
	return conversationID + ":" + strconv.FormatInt(seqSuffix, 10)
}

func (MsgDocModel) GenExceptionMessageBySeqs(seqs []int64) (exceptionMsg []*sdkws.MsgData) {
	for _, v := range seqs {
		msg := new(sdkws.MsgData)
		msg.Seq = v
		exceptionMsg = append(exceptionMsg, msg)
	}
	return exceptionMsg
}

func (MsgDocModel) GenExceptionSuperGroupMessageBySeqs(seqs []int64, groupID string) (exceptionMsg []*sdkws.MsgData) {
	for _, v := range seqs {
		msg := new(sdkws.MsgData)
		msg.Seq = v
		msg.GroupID = groupID
		msg.SessionType = constant.SuperGroupChatType
		exceptionMsg = append(exceptionMsg, msg)
	}
	return exceptionMsg
}
