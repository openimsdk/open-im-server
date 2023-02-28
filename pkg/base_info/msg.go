package base_info

import (
	"Open_IM/pkg/proto/msg"
	sdk_ws "Open_IM/pkg/proto/sdk_ws"
)

type DelMsgReq struct {
	UserID      string   `json:"userID,omitempty" binding:"required"`
	SeqList     []uint32 `json:"seqList,omitempty" binding:"required"`
	OperationID string   `json:"operationID,omitempty" binding:"required"`
}

type DelMsgResp struct {
	CommResp
}

type CleanUpMsgReq struct {
	UserID      string `json:"userID"  binding:"required"`
	OperationID string `json:"operationID"  binding:"required"`
}

type CleanUpMsgResp struct {
	CommResp
}

type DelSuperGroupMsgReq struct {
	UserID      string   `json:"userID" binding:"required"`
	GroupID     string   `json:"groupID" binding:"required"`
	SeqList     []uint32 `json:"seqList,omitempty"`
	IsAllDelete bool     `json:"isAllDelete"`
	OperationID string   `json:"operationID" binding:"required"`
}

type DelSuperGroupMsgResp struct {
	CommResp
}

type MsgDeleteNotificationElem struct {
	GroupID     string   `json:"groupID"`
	IsAllDelete bool     `json:"isAllDelete"`
	SeqList     []uint32 `json:"seqList"`
}

type SetMsgMinSeqReq struct {
	UserID      string `json:"userID"  binding:"required"`
	GroupID     string `json:"groupID"`
	MinSeq      uint32 `json:"minSeq"  binding:"required"`
	OperationID string `json:"operationID"  binding:"required"`
}

type SetMsgMinSeqResp struct {
	CommResp
}

type ModifyMessageReactionExtensionsReq struct {
	OperationID           string                      `json:"operationID" binding:"required"`
	SourceID              string                      `json:"sourceID"  binding:"required"`
	SessionType           int32                       `json:"sessionType" binding:"required"`
	ReactionExtensionList map[string]*sdk_ws.KeyValue `json:"reactionExtensionList,omitempty" binding:"required"`
	ClientMsgID           string                      `json:"clientMsgID" binding:"required"`
	Ex                    *string                     `json:"ex"`
	AttachedInfo          *string                     `json:"attachedInfo"`
	IsReact               bool                        `json:"isReact"`
	IsExternalExtensions  bool                        `json:"isExternalExtensions"`
	MsgFirstModifyTime    int64                       `json:"msgFirstModifyTime"`
	Seq                   uint32                      `json:"seq"`
}

type ModifyMessageReactionExtensionsResp struct {
	CommResp
	Data struct {
		ResultKeyValue     []*msg.KeyValueResp `json:"result"`
		MsgFirstModifyTime int64               `json:"msgFirstModifyTime"`
		IsReact            bool                `json:"isReact"`
	} `json:"data"`
}

type OperateMessageListReactionExtensionsReq struct {
	OperationID            string                                                        `json:"operationID" binding:"required"`
	SourceID               string                                                        `json:"sourceID"  binding:"required"`
	SessionType            int32                                                         `json:"sessionType" binding:"required"`
	IsExternalExtensions   bool                                                          `json:"isExternalExtensions"`
	TypeKeyList            []string                                                      `json:"typeKeyList"`
	MessageReactionKeyList []*msg.GetMessageListReactionExtensionsReq_MessageReactionKey `json:"messageReactionKeyList" binding:"required"`
}

type OperateMessageListReactionExtensionsResp struct {
	CommResp
	Data struct {
		SuccessList []*msg.ExtendMsgResp `json:"successList"`
		FailedList  []*msg.ExtendMsgResp `json:"failedList"`
	} `json:"data"`
}

type SetMessageReactionExtensionsReq ModifyMessageReactionExtensionsReq

type SetMessageReactionExtensionsResp ModifyMessageReactionExtensionsResp

type AddMessageReactionExtensionsReq ModifyMessageReactionExtensionsReq

type AddMessageReactionExtensionsResp ModifyMessageReactionExtensionsResp
type GetMessageListReactionExtensionsReq OperateMessageListReactionExtensionsReq

type GetMessageListReactionExtensionsResp struct {
	CommResp
	Data []*msg.SingleMessageExtensionResult `json:"data"`
}

type DeleteMessageReactionExtensionsReq struct {
	OperationID           string             `json:"operationID" binding:"required"`
	SourceID              string             `json:"sourceID" binding:"required"`
	SessionType           int32              `json:"sessionType" binding:"required"`
	ClientMsgID           string             `json:"clientMsgID" binding:"required"`
	IsExternalExtensions  bool               `json:"isExternalExtensions"`
	MsgFirstModifyTime    int64              `json:"msgFirstModifyTime" binding:"required"`
	ReactionExtensionList []*sdk_ws.KeyValue `json:"reactionExtensionList" binding:"required"`
}

type DeleteMessageReactionExtensionsResp struct {
	CommResp
	Data []*msg.KeyValueResp
}

type ReactionMessageModifierNotification struct {
	Operation                    int                         `json:"operation" binding:"required"`
	SourceID                     string                      `json:"sourceID"  binding:"required"`
	OpUserID                     string                      `json:"opUserID"  binding:"required"`
	SessionType                  int32                       `json:"sessionType" binding:"required"`
	SuccessReactionExtensionList map[string]*sdk_ws.KeyValue `json:"reactionExtensionList,omitempty" binding:"required"`
	ClientMsgID                  string                      `json:"clientMsgID" binding:"required"`
	IsReact                      bool                        `json:"isReact"`
	IsExternalExtensions         bool                        `json:"isExternalExtensions"`
	MsgFirstModifyTime           int64                       `json:"msgFirstModifyTime"`
	Seq                          uint32                      `json:"seq"`
}

type ReactionMessageDeleteNotification struct {
	SourceID                     string                      `json:"sourceID"  binding:"required"`
	OpUserID                     string                      `json:"opUserID"  binding:"required"`
	SessionType                  int32                       `json:"sessionType" binding:"required"`
	SuccessReactionExtensionList map[string]*sdk_ws.KeyValue `json:"reactionExtensionList,omitempty" binding:"required"`
	ClientMsgID                  string                      `json:"clientMsgID" binding:"required"`
	MsgFirstModifyTime           int64                       `json:"msgFirstModifyTime"`
}
