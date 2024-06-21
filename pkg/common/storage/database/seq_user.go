package database

import "context"

type SeqUser interface {
	GetMaxSeq(ctx context.Context, userID string, conversationID string) (int64, error)
	SetMaxSeq(ctx context.Context, userID string, conversationID string, seq int64) error
	GetMinSeq(ctx context.Context, userID string, conversationID string) (int64, error)
	SetMinSeq(ctx context.Context, userID string, conversationID string, seq int64) error
	GetReadSeq(ctx context.Context, userID string, conversationID string) (int64, error)
	SetReadSeq(ctx context.Context, userID string, conversationID string, seq int64) error
}
