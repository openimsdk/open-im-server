package callbackstruct

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	sdkws "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

type CallbackBeforeSendSingleMsgReq struct {
	CommonCallbackReq
	RecvID string `json:"recvID"`
}

type CallbackBeforeSendSingleMsgResp struct {
	CommonCallbackResp
}

type CallbackAfterSendSingleMsgReq struct {
	CommonCallbackReq
	RecvID string `json:"recvID"`
}

type CallbackAfterSendSingleMsgResp struct {
	CommonCallbackResp
}

type CallbackBeforeSendGroupMsgReq struct {
	CommonCallbackReq
	GroupID string `json:"groupID"`
}

type CallbackBeforeSendGroupMsgResp struct {
	CommonCallbackResp
}

type CallbackAfterSendGroupMsgReq struct {
	CommonCallbackReq
	GroupID string `json:"groupID"`
}

type CallbackAfterSendGroupMsgResp struct {
	CommonCallbackResp
}

type CallbackMsgModifyCommandReq struct {
	CommonCallbackReq
}

type CallbackMsgModifyCommandResp struct {
	CommonCallbackResp
	Content          *string                `json:"content"`
	RecvID           *string                `json:"recvID"`
	GroupID          *string                `json:"groupID"`
	ClientMsgID      *string                `json:"clientMsgID"`
	ServerMsgID      *string                `json:"serverMsgID"`
	SenderPlatformID *int32                 `json:"senderPlatformID"`
	SenderNickname   *string                `json:"senderNickname"`
	SenderFaceURL    *string                `json:"senderFaceURL"`
	SessionType      *int32                 `json:"sessionType"`
	MsgFrom          *int32                 `json:"msgFrom"`
	ContentType      *int32                 `json:"contentType"`
	Status           *int32                 `json:"status"`
	Options          *map[string]bool       `json:"options"`
	OfflinePushInfo  *sdkws.OfflinePushInfo `json:"offlinePushInfo"`
	AtUserIDList     *[]string              `json:"atUserIDList"`
	MsgDataList      *[]byte                `json:"msgDataList"`
	AttachedInfo     *string                `json:"attachedInfo"`
	Ex               *string                `json:"ex"`
}
type CallbackBeforeSetMessageReactionExtReq struct {
	OperationID           string `json:"operationID"`
	CallbackCommand       `json:"callbackCommand"`
	ConversationID        string                     `json:"conversationID"`
	OpUserID              string                     `json:"opUserID"`
	SessionType           int32                      `json:"sessionType"`
	ReactionExtensionList map[string]*sdkws.KeyValue `json:"reactionExtensionList"`
	ClientMsgID           string                     `json:"clientMsgID"`
	IsReact               bool                       `json:"isReact"`
	IsExternalExtensions  bool                       `json:"isExternalExtensions"`
	MsgFirstModifyTime    int64                      `json:"msgFirstModifyTime"`
}
type CallbackBeforeSetMessageReactionExtResp struct {
	CommonCallbackResp
	ResultReactionExtensionList []*msg.KeyValueResp `json:"resultReactionExtensionList"`
	MsgFirstModifyTime          int64               `json:"msgFirstModifyTime"`
}
type CallbackDeleteMessageReactionExtReq struct {
	CallbackCommand       `json:"callbackCommand"`
	OperationID           string            `json:"operationID"`
	ConversationID        string            `json:"conversationID"`
	OpUserID              string            `json:"opUserID"`
	SessionType           int32             `json:"sessionType"`
	ReactionExtensionList []*sdkws.KeyValue `json:"reactionExtensionList"`
	ClientMsgID           string            `json:"clientMsgID"`
	IsExternalExtensions  bool              `json:"isExternalExtensions"`
	MsgFirstModifyTime    int64             `json:"msgFirstModifyTime"`
}
type CallbackDeleteMessageReactionExtResp struct {
	CommonCallbackResp
	ResultReactionExtensionList []*msg.KeyValueResp `json:"resultReactionExtensionList"`
	MsgFirstModifyTime          int64               `json:"msgFirstModifyTime"`
}

type CallbackGetMessageListReactionExtReq struct {
	OperationID     string `json:"operationID"`
	CallbackCommand `json:"callbackCommand"`
	ConversationID  string   `json:"conversationID"`
	OpUserID        string   `json:"opUserID"`
	SessionType     int32    `json:"sessionType"`
	TypeKeyList     []string `json:"typeKeyList"`
	//MessageKeyList  []*msg.GetMessageListReactionExtensionsReq_MessageReactionKey `json:"messageKeyList"`
}

type CallbackGetMessageListReactionExtResp struct {
	CommonCallbackResp
	MessageResultList []*msg.SingleMessageExtensionResult `json:"messageResultList"`
}

type CallbackAddMessageReactionExtReq struct {
	OperationID           string `json:"operationID"`
	CallbackCommand       `json:"callbackCommand"`
	ConversationID        string                     `json:"conversationID"`
	OpUserID              string                     `json:"opUserID"`
	SessionType           int32                      `json:"sessionType"`
	ReactionExtensionList map[string]*sdkws.KeyValue `json:"reactionExtensionList"`
	ClientMsgID           string                     `json:"clientMsgID"`
	IsReact               bool                       `json:"isReact"`
	IsExternalExtensions  bool                       `json:"isExternalExtensions"`
	MsgFirstModifyTime    int64                      `json:"msgFirstModifyTime"`
}

type CallbackAddMessageReactionExtResp struct {
	CommonCallbackResp
	ResultReactionExtensionList []*msg.KeyValueResp `json:"resultReactionExtensionList"`
	IsReact                     bool                `json:"isReact"`
	MsgFirstModifyTime          int64               `json:"msgFirstModifyTime"`
}
