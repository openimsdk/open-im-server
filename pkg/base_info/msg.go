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
}

type ModifyMessageReactionExtensionsResp struct {
	CommResp
	Data struct {
		SuccessList []*msg.ExtendMsgResp `json:"successList,omitempty"`
		FailedList  []*msg.ExtendMsgResp `json:"failedList,omitempty"`
	} `json:"data"`
}

type OperateMessageListReactionExtensionsReq struct {
	OperationID            string                                                            `json:"operationID" binding:"required"`
	SourceID               string                                                            `json:"sourceID"  binding:"required"`
	SessionType            string                                                            `json:"sessionType" binding:"required"`
	MessageReactionKeyList []*msg.OperateMessageListReactionExtensionsReq_MessageReactionKey `json:"messageReactionKeyList" binding:"required"`
}

type OperateMessageListReactionExtensionsResp struct {
	CommResp
	Data struct {
		SuccessList []*msg.ExtendMsgResp `json:"successList"`
		FailedList  []*msg.ExtendMsgResp `json:"failedList"`
	} `json:"data"`
}

type SetMessageReactionExtensionsCallbackReq ModifyMessageReactionExtensionsReq

type SetMessageReactionExtensionsCallbackResp ModifyMessageReactionExtensionsResp

type GetMessageListReactionExtensionsReq OperateMessageListReactionExtensionsReq

type GetMessageListReactionExtensionsResp OperateMessageListReactionExtensionsResp

type AddMessageReactionExtensionsReq ModifyMessageReactionExtensionsReq

type AddMessageReactionExtensionsResp ModifyMessageReactionExtensionsResp

type DeleteMessageReactionExtensionsReq OperateMessageListReactionExtensionsReq

type DeleteMessageReactionExtensionsResp OperateMessageListReactionExtensionsResp
