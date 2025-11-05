package cache

import (
	"context"
	relationtb "github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

// arg fn will exec when no data in msgCache.
type ConversationCache interface {
	BatchDeleter
	CloneConversationCache() ConversationCache
	// get user's conversationIDs from msgCache
	GetUserConversationIDs(ctx context.Context, ownerUserID string) ([]string, error)
	GetUserNotNotifyConversationIDs(ctx context.Context, userID string) ([]string, error)
	GetPinnedConversationIDs(ctx context.Context, userID string) ([]string, error)
	DelConversationIDs(userIDs ...string) ConversationCache

	GetUserConversationIDsHash(ctx context.Context, ownerUserID string) (hash uint64, err error)
	DelUserConversationIDsHash(ownerUserIDs ...string) ConversationCache

	// get one conversation from msgCache
	GetConversation(ctx context.Context, ownerUserID, conversationID string) (*relationtb.Conversation, error)
	DelConversations(ownerUserID string, conversationIDs ...string) ConversationCache
	DelUsersConversation(conversationID string, ownerUserIDs ...string) ConversationCache
	// get one conversation from msgCache
	GetConversations(ctx context.Context, ownerUserID string,
		conversationIDs []string) ([]*relationtb.Conversation, error)
	// get one user's all conversations from msgCache
	GetUserAllConversations(ctx context.Context, ownerUserID string) ([]*relationtb.Conversation, error)
	// get user conversation recv msg from msgCache
	GetUserRecvMsgOpt(ctx context.Context, ownerUserID, conversationID string) (opt int, err error)
	DelUserRecvMsgOpt(ownerUserID, conversationID string) ConversationCache
	// get one super group recv msg but do not notification userID list
	// GetSuperGroupRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string) (userIDs []string, err error)
	DelSuperGroupRecvMsgNotNotifyUserIDs(groupID string) ConversationCache
	// get one super group recv msg but do not notification userID list hash
	// GetSuperGroupRecvMsgNotNotifyUserIDsHash(ctx context.Context, groupID string) (hash uint64, err error)
	DelSuperGroupRecvMsgNotNotifyUserIDsHash(groupID string) ConversationCache

	// GetUserAllHasReadSeqs(ctx context.Context, ownerUserID string) (map[string]int64, error)
	DelUserAllHasReadSeqs(ownerUserID string, conversationIDs ...string) ConversationCache

	GetConversationNotReceiveMessageUserIDs(ctx context.Context, conversationID string) ([]string, error)
	DelConversationNotReceiveMessageUserIDs(conversationIDs ...string) ConversationCache
	DelConversationNotNotifyMessageUserIDs(userIDs ...string) ConversationCache
	DelConversationPinnedMessageUserIDs(userIDs ...string) ConversationCache
	DelConversationVersionUserIDs(userIDs ...string) ConversationCache

	FindMaxConversationUserVersion(ctx context.Context, userID string) (*relationtb.VersionLog, error)
}
