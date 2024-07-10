package database

import "context"

type SeqUser interface {
	GetMaxSeq(ctx context.Context, conversationID string, userID string) (int64, error)
	SetMaxSeq(ctx context.Context, conversationID string, userID string, seq int64) error
	GetMinSeq(ctx context.Context, conversationID string, userID string) (int64, error)
	SetMinSeq(ctx context.Context, conversationID string, userID string, seq int64) error
	GetReadSeq(ctx context.Context, conversationID string, userID string) (int64, error)
	SetReadSeq(ctx context.Context, conversationID string, userID string, seq int64) error
}
