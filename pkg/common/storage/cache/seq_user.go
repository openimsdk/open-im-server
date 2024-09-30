package cache

import "context"

type SeqUser interface {
	GetUserMaxSeq(ctx context.Context, conversationID string, userID string) (int64, error)
	SetUserMaxSeq(ctx context.Context, conversationID string, userID string, seq int64) error
	GetUserMinSeq(ctx context.Context, conversationID string, userID string) (int64, error)
	SetUserMinSeq(ctx context.Context, conversationID string, userID string, seq int64) error
	GetUserReadSeq(ctx context.Context, conversationID string, userID string) (int64, error)
	SetUserReadSeq(ctx context.Context, conversationID string, userID string, seq int64) error
	SetUserReadSeqToDB(ctx context.Context, conversationID string, userID string, seq int64) error
	SetUserMinSeqs(ctx context.Context, userID string, seqs map[string]int64) error
	SetUserReadSeqs(ctx context.Context, userID string, seqs map[string]int64) error
	GetUserReadSeqs(ctx context.Context, userID string, conversationIDs []string) (map[string]int64, error)
}
