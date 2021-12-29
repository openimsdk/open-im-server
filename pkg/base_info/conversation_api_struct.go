package base_info

import "Open_IM/pkg/proto/user"

type GetAllConversationMessageOptReq struct {
	OperationID string `json:"operationID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"`
}
type GetAllConversationMessageOptResp struct {
	CommResp
	ConversationOptResultList []*user.OptResult `json:"data"`
}
type GetReceiveMessageOptReq struct {
	ConversationIdList []string `json:"conversationIdList" binding:"required"`
	OperationID        string   `json:"operationID" binding:"required"`
	FromUserID         string   `json:"fromUserID" binding:"required"`
}
type GetReceiveMessageOptResp struct {
	CommResp
	ConversationOptResultList []*user.OptResult `json:"data"`
}
type SetReceiveMessageOptReq struct {
	OperationID        string   `json:"operationID" binding:"required"`
	Opt                *int32   `json:"opt" binding:"required"`
	ConversationIdList []string `json:"conversationIdList" binding:"required"`
}
type SetReceiveMessageOptResp struct {
	CommResp
	OptResultList []*user.OptResult `json:"data"`
}
