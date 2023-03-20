package call_back_struct

import (
	"Open_IM/pkg/proto/msg"
	sdk_ws "Open_IM/pkg/proto/sdk_ws"
)

type CallbackBeforeSendSingleMsgReq struct {
	CommonCallbackReq
	RecvID string `json:"recvID"`
}

type CallbackBeforeSendSingleMsgResp struct {
	*CommonCallbackResp
}

type CallbackAfterSendSingleMsgReq struct {
	CommonCallbackReq
	RecvID string `json:"recvID"`
}

type CallbackAfterSendSingleMsgResp struct {
	*CommonCallbackResp
}

type CallbackBeforeSendGroupMsgReq struct {
	CommonCallbackReq
	GroupID string `json:"groupID"`
}

type CallbackBeforeSendGroupMsgResp struct {
	*CommonCallbackResp
}

type CallbackAfterSendGroupMsgReq struct {
	CommonCallbackReq
	GroupID string `json:"groupID"`
}

type CallbackAfterSendGroupMsgResp struct {
	*CommonCallbackResp
}

type CallbackAfterConsumeGroupMsgReq struct {
	CommonCallbackReq
	GroupID string `json:"groupID"`
}

type CallbackAfterConsumeGroupMsgResp struct {
	*CommonCallbackResp
}
type CallbackMsgModifyCommandReq struct {
	CommonCallbackReq
	RecvID  string `json:"recvID,omitempty"`
	GroupID string `json:"groupID,omitempty"`
}

type CallbackMsgModifyCommandResp struct {
	*CommonCallbackResp
	Content          *string                 `json:"content"`
	RecvID           *string                 `json:"recvID"`
	GroupID          *string                 `json:"groupID"`
	ClientMsgID      *string                 `json:"clientMsgID"`
	ServerMsgID      *string                 `json:"serverMsgID"`
	SenderPlatformID *int32                  `json:"senderPlatformID"`
	SenderNickname   *string                 `json:"senderNickname"`
	SenderFaceURL    *string                 `json:"senderFaceURL"`
	SessionType      *int32                  `json:"sessionType"`
	MsgFrom          *int32                  `json:"msgFrom"`
	ContentType      *int32                  `json:"contentType"`
	Status           *int32                  `json:"status"`
	Options          *map[string]bool        `json:"options"`
	OfflinePushInfo  *sdk_ws.OfflinePushInfo `json:"offlinePushInfo"`
	AtUserIDList     *[]string               `json:"atUserIDList"`
	MsgDataList      *[]byte                 `json:"msgDataList"`
	AttachedInfo     *string                 `json:"attachedInfo"`
	Ex               *string                 `json:"ex"`
}
type CallbackBeforeSetMessageReactionExtReq struct {
	OperationID           string                      `json:"operationID"`
	CallbackCommand       string                      `json:"callbackCommand"`
	SourceID              string                      `json:"sourceID"`
	OpUserID              string                      `json:"opUserID"`
	SessionType           int32                       `json:"sessionType"`
	ReactionExtensionList map[string]*sdk_ws.KeyValue `json:"reactionExtensionList"`
	ClientMsgID           string                      `json:"clientMsgID"`
	IsReact               bool                        `json:"isReact"`
	IsExternalExtensions  bool                        `json:"isExternalExtensions"`
	MsgFirstModifyTime    int64                       `json:"msgFirstModifyTime"`
}
type CallbackBeforeSetMessageReactionExtResp struct {
	*CommonCallbackResp
	ResultReactionExtensionList []*msg.KeyValueResp `json:"resultReactionExtensionList"`
	MsgFirstModifyTime          int64               `json:"msgFirstModifyTime"`
}
type CallbackDeleteMessageReactionExtReq struct {
	OperationID           string             `json:"operationID"`
	CallbackCommand       string             `json:"callbackCommand"`
	SourceID              string             `json:"sourceID"`
	OpUserID              string             `json:"opUserID"`
	SessionType           int32              `json:"sessionType"`
	ReactionExtensionList []*sdk_ws.KeyValue `json:"reactionExtensionList"`
	ClientMsgID           string             `json:"clientMsgID"`
	IsExternalExtensions  bool               `json:"isExternalExtensions"`
	MsgFirstModifyTime    int64              `json:"msgFirstModifyTime"`
}
type CallbackDeleteMessageReactionExtResp struct {
	*CommonCallbackResp
	ResultReactionExtensionList []*msg.KeyValueResp `json:"resultReactionExtensionList"`
	MsgFirstModifyTime          int64               `json:"msgFirstModifyTime"`
}
type CallbackGetMessageListReactionExtReq struct {
	OperationID     string                                                        `json:"operationID"`
	CallbackCommand string                                                        `json:"callbackCommand"`
	SourceID        string                                                        `json:"sourceID"`
	OpUserID        string                                                        `json:"opUserID"`
	SessionType     int32                                                         `json:"sessionType"`
	TypeKeyList     []string                                                      `json:"typeKeyList"`
	MessageKeyList  []*msg.GetMessageListReactionExtensionsReq_MessageReactionKey `json:"messageKeyList"`
}
type CallbackGetMessageListReactionExtResp struct {
	*CommonCallbackResp
	MessageResultList []*msg.SingleMessageExtensionResult `json:"messageResultList"`
}

type CallbackAddMessageReactionExtReq struct {
	OperationID           string                      `json:"operationID"`
	CallbackCommand       string                      `json:"callbackCommand"`
	SourceID              string                      `json:"sourceID"`
	OpUserID              string                      `json:"opUserID"`
	SessionType           int32                       `json:"sessionType"`
	ReactionExtensionList map[string]*sdk_ws.KeyValue `json:"reactionExtensionList"`
	ClientMsgID           string                      `json:"clientMsgID"`
	IsReact               bool                        `json:"isReact"`
	IsExternalExtensions  bool                        `json:"isExternalExtensions"`
	MsgFirstModifyTime    int64                       `json:"msgFirstModifyTime"`
}
type CallbackAddMessageReactionExtResp struct {
	*CommonCallbackResp
	ResultReactionExtensionList []*msg.KeyValueResp `json:"resultReactionExtensionList"`
	IsReact                     bool                `json:"isReact"`
	MsgFirstModifyTime          int64               `json:"msgFirstModifyTime"`
}
