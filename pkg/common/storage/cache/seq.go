package cache

import (
	"context"
)

type SeqCache interface {
	SetMaxSeq(ctx context.Context, conversationID string, maxSeq int64) error
	GetMaxSeqs(ctx context.Context, conversationIDs []string) (map[string]int64, error)
	GetMaxSeq(ctx context.Context, conversationID string) (int64, error)
	SetMinSeq(ctx context.Context, conversationID string, minSeq int64) error
	SetMinSeqs(ctx context.Context, seqs map[string]int64) error
	GetMinSeqs(ctx context.Context, conversationIDs []string) (map[string]int64, error)
	GetMinSeq(ctx context.Context, conversationID string) (int64, error)
	GetConversationUserMinSeq(ctx context.Context, conversationID string, userID string) (int64, error)
	GetConversationUserMinSeqs(ctx context.Context, conversationID string, userIDs []string) (map[string]int64, error)
	SetConversationUserMinSeq(ctx context.Context, conversationID string, userID string, minSeq int64) error
	// seqs map: key userID value minSeq
	SetConversationUserMinSeqs(ctx context.Context, conversationID string, seqs map[string]int64) (err error)
	// seqs map: key conversationID value minSeq
	SetUserConversationsMinSeqs(ctx context.Context, userID string, seqs map[string]int64) error
	// has read seq
	SetHasReadSeq(ctx context.Context, userID string, conversationID string, hasReadSeq int64) error
	// k: user, v: seq
	SetHasReadSeqs(ctx context.Context, conversationID string, hasReadSeqs map[string]int64) error
	// k: conversation, v :seq
	UserSetHasReadSeqs(ctx context.Context, userID string, hasReadSeqs map[string]int64) error
	GetHasReadSeqs(ctx context.Context, userID string, conversationIDs []string) (map[string]int64, error)
	GetHasReadSeq(ctx context.Context, userID string, conversationID string) (int64, error)
}
