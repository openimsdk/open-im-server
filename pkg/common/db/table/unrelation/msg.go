package unrelation

import "strconv"

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
