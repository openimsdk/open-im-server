package base_info

type OptResult struct {
	ConversationID string `json:"conversationID"`
	Result         *int32 `json:"result"`
}
type GetAllConversationMessageOptReq struct {
	OperationID string `json:"operationID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"`
}
type GetAllConversationMessageOptResp struct {
	CommResp
	ConversationOptResultList []*OptResult `json:"data"`
}
type GetReceiveMessageOptReq struct {
	ConversationIDList []string `json:"conversationIDList" binding:"required"`
	OperationID        string   `json:"operationID" binding:"required"`
	FromUserID         string   `json:"fromUserID" binding:"required"`
}
type GetReceiveMessageOptResp struct {
	CommResp
	ConversationOptResultList []*OptResult `json:"data"`
}
type SetReceiveMessageOptReq struct {
	OperationID        string   `json:"operationID" binding:"required"`
	Opt                *int32   `json:"opt" binding:"required"`
	ConversationIDList []string `json:"conversationIDList" binding:"required"`
}
type SetReceiveMessageOptResp struct {
	CommResp
	ConversationOptResultList []*OptResult `json:"data"`
}
