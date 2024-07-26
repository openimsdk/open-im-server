package database

import "context"

type SeqUser interface {
	GetUserMaxSeq(ctx context.Context, conversationID string, userID string) (int64, error)
	SetUserMaxSeq(ctx context.Context, conversationID string, userID string, seq int64) error
	GetUserMinSeq(ctx context.Context, conversationID string, userID string) (int64, error)
	SetUserMinSeq(ctx context.Context, conversationID string, userID string, seq int64) error
	GetUserReadSeq(ctx context.Context, conversationID string, userID string) (int64, error)
	SetUserReadSeq(ctx context.Context, conversationID string, userID string, seq int64) error
	GetUserReadSeqs(ctx context.Context, userID string, conversationID []string) (map[string]int64, error)
}
