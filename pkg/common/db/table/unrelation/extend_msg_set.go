package unrelation

import (
	"context"
	"strconv"
	"strings"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

const (
	CExtendMsgSet = "extend_msgs"

	ExtendMsgMaxNum = 100
)

type ExtendMsgSetModel struct {
	ConversationID   string                    `bson:"source_id" json:"conversationID"`
	SessionType      int32                     `bson:"session_type" json:"sessionType"`
	ExtendMsgs       map[string]ExtendMsgModel `bson:"extend_msgs" json:"extendMsgs"`
	ExtendMsgNum     int32                     `bson:"extend_msg_num" json:"extendMsgNum"`
	CreateTime       int64                     `bson:"create_time" json:"createTime"`               // this block's create time
	MaxMsgUpdateTime int64                     `bson:"max_msg_update_time" json:"maxMsgUpdateTime"` // index find msg
}

type KeyValueModel struct {
	TypeKey          string `bson:"type_key" json:"typeKey"`
	Value            string `bson:"value" json:"value"`
	LatestUpdateTime int64  `bson:"latest_update_time" json:"latestUpdateTime"`
}

type ExtendMsgModel struct {
	ReactionExtensionList map[string]KeyValueModel `bson:"reaction_extension_list" json:"reactionExtensionList"`
	ClientMsgID           string                   `bson:"client_msg_id" json:"clientMsgID"`
	MsgFirstModifyTime    int64                    `bson:"msg_first_modify_time" json:"msgFirstModifyTime"` // this extendMsg create time
	AttachedInfo          string                   `bson:"attached_info" json:"attachedInfo"`
	Ex                    string                   `bson:"ex" json:"ex"`
}

type ExtendMsgSetModelInterface interface {
	CreateExtendMsgSet(ctx context.Context, set *ExtendMsgSetModel) error
	GetAllExtendMsgSet(ctx context.Context, conversationID string, opts *GetAllExtendMsgSetOpts) (sets []*ExtendMsgSetModel, err error)
	GetExtendMsgSet(ctx context.Context, conversationID string, sessionType int32, maxMsgUpdateTime int64) (*ExtendMsgSetModel, error)
	InsertExtendMsg(ctx context.Context, conversationID string, sessionType int32, msg *ExtendMsgModel) error
	InsertOrUpdateReactionExtendMsgSet(ctx context.Context, conversationID string, sessionType int32, clientMsgID string, msgFirstModifyTime int64, reactionExtensionList map[string]*KeyValueModel) error
	DeleteReactionExtendMsgSet(ctx context.Context, conversationID string, sessionType int32, clientMsgID string, msgFirstModifyTime int64, reactionExtensionList map[string]*KeyValueModel) error
	TakeExtendMsg(ctx context.Context, conversationID string, sessionType int32, clientMsgID string, maxMsgUpdateTime int64) (extendMsg *ExtendMsgModel, err error)
}

func (ExtendMsgSetModel) TableName() string {
	return CExtendMsgSet
}

func (ExtendMsgSetModel) GetExtendMsgMaxNum() int32 {
	return ExtendMsgMaxNum
}

func (ExtendMsgSetModel) GetConversationID(ID string, index int32) string {
	return ID + ":" + strconv.Itoa(int(index))
}

func (e *ExtendMsgSetModel) SplitConversationIDAndGetIndex() int32 {
	l := strings.Split(e.ConversationID, ":")
	index, _ := strconv.Atoi(l[len(l)-1])
	return int32(index)
}

type GetAllExtendMsgSetOpts struct {
	ExcludeExtendMsgs bool
}

func (ExtendMsgSetModel) Pb2Model(reactionExtensionList map[string]*sdkws.KeyValue) map[string]*KeyValueModel {
	r := make(map[string]*KeyValueModel)
	for key, value := range reactionExtensionList {
		r[key] = &KeyValueModel{
			TypeKey:          value.TypeKey,
			Value:            value.Value,
			LatestUpdateTime: value.LatestUpdateTime,
		}
	}
	return r
}
