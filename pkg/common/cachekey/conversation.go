package cachekey

const (
	ConversationKey                          = "CONVERSATION:"
	ConversationIDsKey                       = "CONVERSATION_IDS:"
	ConversationIDsHashKey                   = "CONVERSATION_IDS_HASH:"
	ConversationHasReadSeqKey                = "CONVERSATION_HAS_READ_SEQ:"
	RecvMsgOptKey                            = "RECV_MSG_OPT:"
	SuperGroupRecvMsgNotNotifyUserIDsKey     = "SUPER_GROUP_RECV_MSG_NOT_NOTIFY_USER_IDS:"
	SuperGroupRecvMsgNotNotifyUserIDsHashKey = "SUPER_GROUP_RECV_MSG_NOT_NOTIFY_USER_IDS_HASH:"
	ConversationNotReceiveMessageUserIDsKey  = "CONVERSATION_NOT_RECEIVE_MESSAGE_USER_IDS:"
)

func GetConversationKey(ownerUserID, conversationID string) string {
	return ConversationKey + ownerUserID + ":" + conversationID
}

func GetConversationIDsKey(ownerUserID string) string {
	return ConversationIDsKey + ownerUserID
}

func GetSuperGroupRecvNotNotifyUserIDsKey(groupID string) string {
	return SuperGroupRecvMsgNotNotifyUserIDsKey + groupID
}

func GetRecvMsgOptKey(ownerUserID, conversationID string) string {
	return RecvMsgOptKey + ownerUserID + ":" + conversationID
}

func GetSuperGroupRecvNotNotifyUserIDsHashKey(groupID string) string {
	return SuperGroupRecvMsgNotNotifyUserIDsHashKey + groupID
}

func GetConversationHasReadSeqKey(ownerUserID, conversationID string) string {
	return ConversationHasReadSeqKey + ownerUserID + ":" + conversationID
}

func GetConversationNotReceiveMessageUserIDsKey(conversationID string) string {
	return ConversationNotReceiveMessageUserIDsKey + conversationID
}

func GetUserConversationIDsHashKey(ownerUserID string) string {
	return ConversationIDsHashKey + ownerUserID
}
