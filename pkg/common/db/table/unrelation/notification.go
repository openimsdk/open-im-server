package unrelation

import (
	"context"
	"strconv"
	"strings"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

const (
	singleGocNotificationNum = 5000
	Notification             = "notification"
	//OldestList      = 0
	//NewestList      = -1
)

type NotificationDocModel struct {
	DocID string                  `bson:"uid"`
	Msg   []NotificationInfoModel `bson:"msg"`
}

type NotificationInfoModel struct {
	SendTime int64  `bson:"sendtime"`
	Msg      []byte `bson:"msg"`
}

type NotificationDocModelInterface interface {
	PushMsgsToDoc(ctx context.Context, docID string, msgsToMongo []NotificationInfoModel) error
	Create(ctx context.Context, model *NotificationDocModel) error
	UpdateMsgStatusByIndexInOneDoc(ctx context.Context, docID string, msg *sdkws.MsgData, seqIndex int, status int32) error
	FindOneByDocID(ctx context.Context, docID string) (*NotificationDocModel, error)
	GetNewestMsg(ctx context.Context, conversationID string) (*NotificationInfoModel, error)
	GetOldestMsg(ctx context.Context, conversationID string) (*NotificationInfoModel, error)
	Delete(ctx context.Context, docIDs []string) error
	GetMsgsByIndex(ctx context.Context, conversationID string, index int64) (*NotificationDocModel, error)
	UpdateOneDoc(ctx context.Context, msg *NotificationDocModel) error
}

func (NotificationDocModel) TableName() string {
	return Notification
}

func (NotificationDocModel) GetsingleGocNotificationNum() int64 {
	return singleGocNotificationNum
}

func (m *NotificationDocModel) IsFull() bool {
	index, _ := strconv.Atoi(strings.Split(m.DocID, ":")[1])
	if index == 0 {
		if len(m.Msg) >= singleGocNotificationNum-1 {
			return true
		}
	}
	if len(m.Msg) >= singleGocNotificationNum {
		return true
	}

	return false
}

func (m NotificationDocModel) GetDocID(conversationID string, seq int64) string {
	seqSuffix := seq / singleGocNotificationNum
	return m.indexGen(conversationID, seqSuffix)
}

func (m NotificationDocModel) GetSeqDocIDList(userID string, maxSeq int64) []string {
	seqMaxSuffix := maxSeq / singleGocNotificationNum
	var seqUserIDs []string
	for i := 0; i <= int(seqMaxSuffix); i++ {
		seqUserID := m.indexGen(userID, int64(i))
		seqUserIDs = append(seqUserIDs, seqUserID)
	}
	return seqUserIDs
}

func (m NotificationDocModel) getSeqSuperGroupID(groupID string, seq int64) string {
	seqSuffix := seq / singleGocNotificationNum
	return m.superGroupIndexGen(groupID, seqSuffix)
}

func (m NotificationDocModel) superGroupIndexGen(groupID string, seqSuffix int64) string {
	return "super_group_" + groupID + ":" + strconv.FormatInt(int64(seqSuffix), 10)
}

func (m NotificationDocModel) GetDocIDSeqsMap(conversationID string, seqs []int64) map[string][]int64 {
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

func (m NotificationDocModel) getMsgIndex(seq uint32) int {
	seqSuffix := seq / singleGocNotificationNum
	var index uint32
	if seqSuffix == 0 {
		index = (seq - seqSuffix*singleGocNotificationNum) - 1
	} else {
		index = seq - seqSuffix*singleGocNotificationNum
	}
	return int(index)
}

func (m NotificationDocModel) indexGen(conversationID string, seqSuffix int64) string {
	return conversationID + ":" + strconv.FormatInt(seqSuffix, 10)
}

func (NotificationDocModel) GenExceptionMessageBySeqs(seqs []int64) (exceptionMsg []*sdkws.MsgData) {
	for _, v := range seqs {
		msg := new(sdkws.MsgData)
		msg.Seq = v
		exceptionMsg = append(exceptionMsg, msg)
	}
	return exceptionMsg
}

func (NotificationDocModel) GenExceptionSuperGroupMessageBySeqs(seqs []int64, groupID string) (exceptionMsg []*sdkws.MsgData) {
	for _, v := range seqs {
		msg := new(sdkws.MsgData)
		msg.Seq = v
		msg.GroupID = groupID
		msg.SessionType = constant.SuperGroupChatType
		exceptionMsg = append(exceptionMsg, msg)
	}
	return exceptionMsg
}
