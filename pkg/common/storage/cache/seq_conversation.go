package cache

import "context"

type SeqConversationCache interface {
	Malloc(ctx context.Context, conversationID string, size int64) (int64, error)
	GetMaxSeq(ctx context.Context, conversationID string) (int64, error)
	SetMinSeq(ctx context.Context, conversationID string, seq int64) error
	GetMinSeq(ctx context.Context, conversationID string) (int64, error)
	GetMaxSeqs(ctx context.Context, conversationIDs []string) (map[string]int64, error)
	SetMinSeqs(ctx context.Context, seqs map[string]int64) error
}
