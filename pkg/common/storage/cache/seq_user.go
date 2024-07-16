package cache

import "context"

type SeqUser interface {
	GetMaxSeq(ctx context.Context, conversationID string, userID string) (int64, error)
	SetMaxSeq(ctx context.Context, conversationID string, userID string, seq int64) error
	GetMinSeq(ctx context.Context, conversationID string, userID string) (int64, error)
	SetMinSeq(ctx context.Context, conversationID string, userID string, seq int64) error
	GetReadSeq(ctx context.Context, conversationID string, userID string) (int64, error)
	SetReadSeq(ctx context.Context, conversationID string, userID string, seq int64) error
	SetMinSeqs(ctx context.Context, userID string, seqs map[string]int64) error
	SetReadSeqs(ctx context.Context, userID string, seqs map[string]int64) error
	GetReadSeqs(ctx context.Context, userID string, conversationIDs []string) (map[string]int64, error)
}
