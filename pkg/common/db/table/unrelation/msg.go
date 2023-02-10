package unrelation

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/proto/sdkws"
	"strconv"
)

const (
	singleGocMsgNum = 5000
	CChat           = "msg"
)

type UserMsgDocModel struct {
	DocID string         `bson:"uid"`
	Msg   []MsgInfoModel `bson:"msg"`
}

type MsgInfoModel struct {
	SendTime int64  `bson:"sendtime"`
	Msg      []byte `bson:"msg"`
}

func (UserMsgDocModel) TableName() string {
	return CChat
}

func (UserMsgDocModel) GetSingleDocMsgNum() int {
	return singleGocMsgNum
}

func (u UserMsgDocModel) getSeqUid(uid string, seq uint32) string {
	seqSuffix := seq / singleGocMsgNum
	return u.indexGen(uid, seqSuffix)
}

func (u UserMsgDocModel) getSeqUserIDList(userID string, maxSeq uint32) []string {
	seqMaxSuffix := maxSeq / singleGocMsgNum
	var seqUserIDList []string
	for i := 0; i <= int(seqMaxSuffix); i++ {
		seqUserID := u.indexGen(userID, uint32(i))
		seqUserIDList = append(seqUserIDList, seqUserID)
	}
	return seqUserIDList
}

func (UserMsgDocModel) getSeqSuperGroupID(groupID string, seq uint32) string {
	seqSuffix := seq / singleGocMsgNum
	return superGroupIndexGen(groupID, seqSuffix)
}

func (u UserMsgDocModel) GetSeqUid(uid string, seq uint32) string {
	return u.getSeqUid(uid, seq)
}

func (u UserMsgDocModel) GetDocIDSeqsMap(uid string, seqs []uint32) map[string][]uint32 {
	t := make(map[string][]uint32)
	for i := 0; i < len(seqs); i++ {
		seqUid := u.getSeqUid(uid, seqs[i])
		if value, ok := t[seqUid]; !ok {
			var temp []uint32
			t[seqUid] = append(temp, seqs[i])
		} else {
			t[seqUid] = append(value, seqs[i])
		}
	}
	return t
}

func (UserMsgDocModel) getMsgIndex(seq uint32) int {
	seqSuffix := seq / singleGocMsgNum
	var index uint32
	if seqSuffix == 0 {
		index = (seq - seqSuffix*singleGocMsgNum) - 1
	} else {
		index = seq - seqSuffix*singleGocMsgNum
	}
	return int(index)
}

func (UserMsgDocModel) indexGen(uid string, seqSuffix uint32) string {
	return uid + ":" + strconv.FormatInt(int64(seqSuffix), 10)
}

func (UserMsgDocModel) genExceptionMessageBySeqList(seqList []uint32) (exceptionMsg []*sdkws.MsgData) {
	for _, v := range seqList {
		msg := new(sdkws.MsgData)
		msg.Seq = v
		exceptionMsg = append(exceptionMsg, msg)
	}
	return exceptionMsg
}

func (UserMsgDocModel) genExceptionSuperGroupMessageBySeqList(seqList []uint32, groupID string) (exceptionMsg []*sdkws.MsgData) {
	for _, v := range seqList {
		msg := new(sdkws.MsgData)
		msg.Seq = v
		msg.GroupID = groupID
		msg.SessionType = constant.SuperGroupChatType
		exceptionMsg = append(exceptionMsg, msg)
	}
	return exceptionMsg
}
