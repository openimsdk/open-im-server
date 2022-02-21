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
	FromUserID         string   `json:"fromUserID" binding:"required"`
	OperationID        string   `json:"operationID" binding:"required"`
	Opt                *int32   `json:"opt" binding:"required"`
	ConversationIDList []string `json:"conversationIDList" binding:"required"`
}
type SetReceiveMessageOptResp struct {
	CommResp
	ConversationOptResultList []*OptResult `json:"data"`
}

//type Conversation struct {
//	OwnerUserID      string `gorm:"column:owner_user_id;primary_key;type:char(128)" json:"OwnerUserID"`
//	ConversationID   string `gorm:"column:conversation_id;primary_key;type:char(128)" json:"conversationID"`
//	ConversationType int32  `gorm:"column:conversation_type" json:"conversationType"`
//	UserID           string `gorm:"column:user_id;type:char(64)" json:"userID"`
//	GroupID          string `gorm:"column:group_id;type:char(128)" json:"groupID"`
//	RecvMsgOpt       int32  `gorm:"column:recv_msg_opt" json:"recvMsgOpt"`
//	UnreadCount      int32  `gorm:"column:unread_count" json:"unreadCount"`
//	DraftTextTime    int64  `gorm:"column:draft_text_time" json:"draftTextTime"`
//	IsPinned         bool   `gorm:"column:is_pinned" json:"isPinned"`
//	AttachedInfo     string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
//	Ex               string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
//}

