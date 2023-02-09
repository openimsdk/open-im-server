package unrelation

import (
	common "Open_IM/pkg/proto/sdk_ws"
	"context"
	"strconv"
	"strings"
)

const (
	CExtendMsgSet = "extend_msgs"

	ExtendMsgMaxNum = 100
)

type ExtendMsgSetModel struct {
	SourceID         string               `bson:"source_id" json:"sourceID"`
	SessionType      int32                `bson:"session_type" json:"sessionType"`
	ExtendMsgs       map[string]ExtendMsg `bson:"extend_msgs" json:"extendMsgs"`
	ExtendMsgNum     int32                `bson:"extend_msg_num" json:"extendMsgNum"`
	CreateTime       int64                `bson:"create_time" json:"createTime"`               // this block's create time
	MaxMsgUpdateTime int64                `bson:"max_msg_update_time" json:"maxMsgUpdateTime"` // index find msg
}

type KeyValueModel struct {
	TypeKey          string `bson:"type_key" json:"typeKey"`
	Value            string `bson:"value" json:"value"`
	LatestUpdateTime int64  `bson:"latest_update_time" json:"latestUpdateTime"`
}

type ExtendMsg struct {
	ReactionExtensionList map[string]KeyValueModel `bson:"reaction_extension_list" json:"reactionExtensionList"`
	ClientMsgID           string                   `bson:"client_msg_id" json:"clientMsgID"`
	MsgFirstModifyTime    int64                    `bson:"msg_first_modify_time" json:"msgFirstModifyTime"` // this extendMsg create time
	AttachedInfo          string                   `bson:"attached_info" json:"attachedInfo"`
	Ex                    string                   `bson:"ex" json:"ex"`
}

func (ExtendMsgSetModel) TableName() string {
	return CExtendMsgSet
}

func (ExtendMsgSetModel) GetExtendMsgMaxNum() int32 {
	return ExtendMsgMaxNum
}

func (ExtendMsgSetModel) GetSourceID(ID string, index int32) string {
	return ID + ":" + strconv.Itoa(int(index))
}

func (e *ExtendMsgSetModel) SplitSourceIDAndGetIndex() int32 {
	l := strings.Split(e.SourceID, ":")
	index, _ := strconv.Atoi(l[len(l)-1])
	return int32(index)
}

type GetAllExtendMsgSetOpts struct {
	ExcludeExtendMsgs bool
}

type ExtendMsgSetInterface interface {
	CreateExtendMsgSet(ctx context.Context, set *ExtendMsgSetModel) error
	GetAllExtendMsgSet(ctx context.Context, ID string, opts *GetAllExtendMsgSetOpts) (sets []*ExtendMsgSetModel, err error)
	GetExtendMsgSet(ctx context.Context, sourceID string, sessionType int32, maxMsgUpdateTime int64) (*ExtendMsgSetModel, error)
	InsertExtendMsg(ctx context.Context, sourceID string, sessionType int32, msg *ExtendMsg) error
	InsertOrUpdateReactionExtendMsgSet(ctx context.Context, sourceID string, sessionType int32, clientMsgID string, msgFirstModifyTime int64, reactionExtensionList map[string]*common.KeyValue) error
	DeleteReactionExtendMsgSet(ctx context.Context, sourceID string, sessionType int32, clientMsgID string, msgFirstModifyTime int64, reactionExtensionList map[string]*common.KeyValue) error
	GetExtendMsg(ctx context.Context, sourceID string, sessionType int32, clientMsgID string, maxMsgUpdateTime int64) (extendMsg *ExtendMsg, err error)
}
