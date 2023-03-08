package apistruct

import (
	"OpenIM/pkg/proto/msg"
	sdkws "OpenIM/pkg/proto/sdkws"
)

type DelMsgReq struct {
	UserID      string   `json:"userID,omitempty" binding:"required"`
	SeqList     []uint32 `json:"seqList,omitempty" binding:"required"`
	OperationID string   `json:"operationID,omitempty" binding:"required"`
}

type DelMsgResp struct {
}

type CleanUpMsgReq struct {
	UserID      string `json:"userID"  binding:"required"`
	OperationID string `json:"operationID"  binding:"required"`
}

type CleanUpMsgResp struct {
}

type DelSuperGroupMsgReq struct {
	UserID      string   `json:"userID" binding:"required"`
	GroupID     string   `json:"groupID" binding:"required"`
	SeqList     []uint32 `json:"seqList,omitempty"`
	IsAllDelete bool     `json:"isAllDelete"`
	OperationID string   `json:"operationID" binding:"required"`
}

type DelSuperGroupMsgResp struct {
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
}

type ModifyMessageReactionExtensionsReq struct {
	OperationID           string                     `json:"operationID" binding:"required"`
	SourceID              string                     `json:"sourceID"  binding:"required"`
	SessionType           int32                      `json:"sessionType" binding:"required"`
	ReactionExtensionList map[string]*sdkws.KeyValue `json:"reactionExtensionList,omitempty" binding:"required"`
	ClientMsgID           string                     `json:"clientMsgID" binding:"required"`
	Ex                    *string                    `json:"ex"`
	AttachedInfo          *string                    `json:"attachedInfo"`
	IsReact               bool                       `json:"isReact"`
	IsExternalExtensions  bool                       `json:"isExternalExtensions"`
	MsgFirstModifyTime    int64                      `json:"msgFirstModifyTime"`
}

type ModifyMessageReactionExtensionsResp struct {
	Data struct {
		ResultKeyValue     []*msg.KeyValueResp `json:"result"`
		MsgFirstModifyTime int64               `json:"msgFirstModifyTime"`
		IsReact            bool                `json:"isReact"`
	} `json:"data"`
}

//type OperateMessageListReactionExtensionsReq struct {
//	OperationID            string                                                        `json:"operationID" binding:"required"`
//	SourceID               string                                                        `json:"sourceID"  binding:"required"`
//	SessionType            string                                                        `json:"sessionType" binding:"required"`
//	MessageReactionKeyList []*msg.GetMessageListReactionExtensionsReq_MessageReactionKey `json:"messageReactionKeyList" binding:"required"`
//}

type OperateMessageListReactionExtensionsResp struct {
	Data struct {
		SuccessList []*msg.ExtendMsgResp `json:"successList"`
		FailedList  []*msg.ExtendMsgResp `json:"failedList"`
	} `json:"data"`
}

type SetMessageReactionExtensionsCallbackReq ModifyMessageReactionExtensionsReq

type SetMessageReactionExtensionsCallbackResp ModifyMessageReactionExtensionsResp

//type GetMessageListReactionExtensionsReq OperateMessageListReactionExtensionsReq

type GetMessageListReactionExtensionsResp struct {
	Data []*msg.SingleMessageExtensionResult `json:"data"`
}

type AddMessageReactionExtensionsReq ModifyMessageReactionExtensionsReq

type AddMessageReactionExtensionsResp ModifyMessageReactionExtensionsResp

type DeleteMessageReactionExtensionsReq struct {
	OperationID           string            `json:"operationID" binding:"required"`
	SourceID              string            `json:"sourceID" binding:"required"`
	SessionType           int32             `json:"sessionType" binding:"required"`
	ClientMsgID           string            `json:"clientMsgID" binding:"required"`
	IsExternalExtensions  bool              `json:"isExternalExtensions"`
	MsgFirstModifyTime    int64             `json:"msgFirstModifyTime" binding:"required"`
	ReactionExtensionList []*sdkws.KeyValue `json:"reactionExtensionList" binding:"required"`
}

type DeleteMessageReactionExtensionsResp struct {
	Data []*msg.KeyValueResp
}

type ReactionMessageModifierNotification struct {
	SourceID                  string                     `json:"sourceID"  binding:"required"`
	OpUserID                  string                     `json:"opUserID"  binding:"required"`
	SessionType               int32                      `json:"sessionType" binding:"required"`
	SuccessReactionExtensions map[string]*sdkws.KeyValue `json:"reactionExtensionList,omitempty" binding:"required"`
	ClientMsgID               string                     `json:"clientMsgID" binding:"required"`
	IsReact                   bool                       `json:"isReact"`
	IsExternalExtensions      bool                       `json:"isExternalExtensions"`
	MsgFirstModifyTime        int64                      `json:"msgFirstModifyTime"`
}

type ReactionMessageDeleteNotification struct {
	SourceID                  string                     `json:"sourceID"  binding:"required"`
	OpUserID                  string                     `json:"opUserID"  binding:"required"`
	SessionType               int32                      `json:"sessionType" binding:"required"`
	SuccessReactionExtensions map[string]*sdkws.KeyValue `json:"reactionExtensionList,omitempty" binding:"required"`
	ClientMsgID               string                     `json:"clientMsgID" binding:"required"`
	MsgFirstModifyTime        int64                      `json:"msgFirstModifyTime"`
}
