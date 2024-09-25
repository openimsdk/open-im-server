// Copyright © 2024 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cachekey

const (
	ConversationKey                          = "CONVERSATION:"
	ConversationIDsKey                       = "CONVERSATION_IDS:"
	NotNotifyConversationIDsKey              = "NOT_NOTIFY_CONVERSATION_IDS:"
	PinnedConversationIDsKey                 = "PINNED_CONVERSATION_IDS:"
	ConversationIDsHashKey                   = "CONVERSATION_IDS_HASH:"
	ConversationHasReadSeqKey                = "CONVERSATION_HAS_READ_SEQ:"
	RecvMsgOptKey                            = "RECV_MSG_OPT:"
	SuperGroupRecvMsgNotNotifyUserIDsKey     = "SUPER_GROUP_RECV_MSG_NOT_NOTIFY_USER_IDS:"
	SuperGroupRecvMsgNotNotifyUserIDsHashKey = "SUPER_GROUP_RECV_MSG_NOT_NOTIFY_USER_IDS_HASH:"
	ConversationNotReceiveMessageUserIDsKey  = "CONVERSATION_NOT_RECEIVE_MESSAGE_USER_IDS:"
	ConversationUserMaxKey                   = "CONVERSATION_USER_MAX:"
)

func GetConversationKey(ownerUserID, conversationID string) string {
	return ConversationKey + ownerUserID + ":" + conversationID
}

func GetConversationIDsKey(ownerUserID string) string {
	return ConversationIDsKey + ownerUserID
}

func GetNotNotifyConversationIDsKey(ownerUserID string) string {
	return NotNotifyConversationIDsKey + ownerUserID
}

func GetPinnedConversationIDs(ownerUserID string) string {
	return PinnedConversationIDsKey + ownerUserID
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

func GetConversationUserMaxVersionKey(userID string) string {
	return ConversationUserMaxKey + userID
}
