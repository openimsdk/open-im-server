package unrelation

import (
	"strconv"
	"strings"
)

const (
	CExtendMsgSet = "extend_msgs"

	ExtendMsgMaxNum = 100
)

type ExtendMsgSet struct {
	SourceID         string               `bson:"source_id" json:"sourceID"`
	SessionType      int32                `bson:"session_type" json:"sessionType"`
	ExtendMsgs       map[string]ExtendMsg `bson:"extend_msgs" json:"extendMsgs"`
	ExtendMsgNum     int32                `bson:"extend_msg_num" json:"extendMsgNum"`
	CreateTime       int64                `bson:"create_time" json:"createTime"`               // this block's create time
	MaxMsgUpdateTime int64                `bson:"max_msg_update_time" json:"maxMsgUpdateTime"` // index find msg
}

type KeyValue struct {
	TypeKey          string `bson:"type_key" json:"typeKey"`
	Value            string `bson:"value" json:"value"`
	LatestUpdateTime int64  `bson:"latest_update_time" json:"latestUpdateTime"`
}

type ExtendMsg struct {
	ReactionExtensionList map[string]KeyValue `bson:"reaction_extension_list" json:"reactionExtensionList"`
	ClientMsgID           string              `bson:"client_msg_id" json:"clientMsgID"`
	MsgFirstModifyTime    int64               `bson:"msg_first_modify_time" json:"msgFirstModifyTime"` // this extendMsg create time
	AttachedInfo          string              `bson:"attached_info" json:"attachedInfo"`
	Ex                    string              `bson:"ex" json:"ex"`
}

func (ExtendMsgSet) TableName() string {
	return CExtendMsgSet
}

func (ExtendMsgSet) GetExtendMsgMaxNum() int32 {
	return ExtendMsgMaxNum
}

func (ExtendMsgSet) GetSourceID(ID string, index int32) string {
	return ID + ":" + strconv.Itoa(int(index))
}

func (e *ExtendMsgSet) SplitSourceIDAndGetIndex() int32 {
	l := strings.Split(e.SourceID, ":")
	index, _ := strconv.Atoi(l[len(l)-1])
	return int32(index)
}
