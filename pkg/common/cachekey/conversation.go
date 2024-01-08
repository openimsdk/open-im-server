package cachekey

import "time"

const (
	conversationKey                          = "CONVERSATION:"
	conversationIDsKey                       = "CONVERSATION_IDS:"
	conversationIDsHashKey                   = "CONVERSATION_IDS_HASH:"
	conversationHasReadSeqKey                = "CONVERSATION_HAS_READ_SEQ:"
	recvMsgOptKey                            = "RECV_MSG_OPT:"
	superGroupRecvMsgNotNotifyUserIDsKey     = "SUPER_GROUP_RECV_MSG_NOT_NOTIFY_USER_IDS:"
	superGroupRecvMsgNotNotifyUserIDsHashKey = "SUPER_GROUP_RECV_MSG_NOT_NOTIFY_USER_IDS_HASH:"
	conversationNotReceiveMessageUserIDsKey  = "CONVERSATION_NOT_RECEIVE_MESSAGE_USER_IDS:"

	conversationExpireTime = time.Second * 60 * 60 * 12
)

func GetConversationKey(ownerUserID, conversationID string) string {
	return conversationKey + ownerUserID + ":" + conversationID
}

func GetConversationIDsKey(ownerUserID string) string {
	return conversationIDsKey + ownerUserID
}

func GetSuperGroupRecvNotNotifyUserIDsKey(groupID string) string {
	return superGroupRecvMsgNotNotifyUserIDsKey + groupID
}

func GetRecvMsgOptKey(ownerUserID, conversationID string) string {
	return recvMsgOptKey + ownerUserID + ":" + conversationID
}

func GetSuperGroupRecvNotNotifyUserIDsHashKey(groupID string) string {
	return superGroupRecvMsgNotNotifyUserIDsHashKey + groupID
}

func GetConversationHasReadSeqKey(ownerUserID, conversationID string) string {
	return conversationHasReadSeqKey + ownerUserID + ":" + conversationID
}

func GetConversationNotReceiveMessageUserIDsKey(conversationID string) string {
	return conversationNotReceiveMessageUserIDsKey + conversationID
}

func GetUserConversationIDsHashKey(ownerUserID string) string {
	return conversationIDsHashKey + ownerUserID
}
