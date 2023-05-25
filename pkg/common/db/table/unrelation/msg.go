package unrelation

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
	"strings"
)

const (
	singleGocMsgNum = 5000
	Msg             = "msg"
	OldestList      = 0
	NewestList      = -1
)

type MsgDocModel struct {
	DocID string          `bson:"doc_id"`
	Msg   []*MsgInfoModel `bson:"msgs"`
}

type RevokeModel struct {
	UserID   string `bson:"user_id"`
	Nickname string `bson:"nickname"`
	Time     int64  `bson:"time"`
}

type OfflinePushModel struct {
	Title         string `bson:"title"`
	Desc          string `bson:"desc"`
	Ex            string `bson:"ex"`
	IOSPushSound  string `bson:"ios_push_sound"`
	IOSBadgeCount bool   `bson:"ios_badge_count"`
}

type MsgDataModel struct {
	SendID           string            `bson:"send_id"`
	RecvID           string            `bson:"recv_id"`
	GroupID          string            `bson:"group_id"`
	ClientMsgID      string            `bson:"client_msg_id"`
	ServerMsgID      string            `bson:"server_msg_id"`
	SenderPlatformID int32             `bson:"sender_platform_id"`
	SenderNickname   string            `bson:"sender_nickname"`
	SenderFaceURL    string            `bson:"sender_face_url"`
	SessionType      int32             `bson:"session_type"`
	MsgFrom          int32             `bson:"msg_from"`
	ContentType      int32             `bson:"content_type"`
	Content          string            `bson:"content"`
	Seq              int64             `bson:"seq"`
	SendTime         int64             `bson:"send_time"`
	CreateTime       int64             `bson:"create_time"`
	Status           int32             `bson:"status"`
	Options          map[string]bool   `bson:"options"`
	OfflinePush      *OfflinePushModel `bson:"offline_push"`
	AtUserIDList     []string          `bson:"at_user_id_list"`
	AttachedInfo     string            `bson:"attached_info"`
	Ex               string            `bson:"ex"`
}

type MsgInfoModel struct {
	Msg      *MsgDataModel `bson:"msg"`
	Revoke   *RevokeModel  `bson:"revoke"`
	DelList  []string      `bson:"del_list"`
	ReadList []string      `bson:"read_list"`
}

type MsgDocModelInterface interface {
	PushMsgsToDoc(ctx context.Context, docID string, msgsToMongo []MsgInfoModel) error
	Create(ctx context.Context, model *MsgDocModel) error
	UpdateMsg(ctx context.Context, docID string, index int64, key string, value any) (*mongo.UpdateResult, error)
	PushUnique(ctx context.Context, docID string, index int64, key string, value any) (*mongo.UpdateResult, error)
	UpdateMsgContent(ctx context.Context, docID string, index int64, msg []byte) error
	IsExistDocID(ctx context.Context, docID string) (bool, error)
	UpdateMsgStatusByIndexInOneDoc(ctx context.Context, docID string, msg *sdkws.MsgData, seqIndex int, status int32) error
	FindOneByDocID(ctx context.Context, docID string) (*MsgDocModel, error)
	GetMsgBySeqIndexIn1Doc(ctx context.Context, docID string, seqs []int64) ([]*sdkws.MsgData, error)
	GetMsgAndIndexBySeqsInOneDoc(ctx context.Context, docID string, seqs []int64) (seqMsgs []*sdkws.MsgData, indexes []int, unExistSeqs []int64, err error)
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
	//return m.Msg[len(m.Msg)-1].SendTime != 0
	return false
}

func (m MsgDocModel) GetDocID(conversationID string, seq int64) string {
	seqSuffix := seq / singleGocMsgNum
	return m.indexGen(conversationID, seqSuffix)
}

func (m MsgDocModel) IndexDocID(conversationID string, index int64) string {
	return m.indexGen(conversationID, index)
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

func (m MsgDocModel) ToNextDoc(docID string) string {
	l := strings.Split(docID, ":")
	index, _ := strconv.Atoi(l[len(l)-1])
	index++
	return strings.Split(docID, ":")[0] + ":" + strconv.Itoa(index)
}

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
