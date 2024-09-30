package database

import "context"

type SeqTime struct {
	Seq  int64
	Time int64
}

type SeqConversation interface {
	Malloc(ctx context.Context, conversationID string, size int64) (int64, error)
	GetMaxSeq(ctx context.Context, conversationID string) (int64, error)
	SetMaxSeq(ctx context.Context, conversationID string, seq int64) error
	GetMinSeq(ctx context.Context, conversationID string) (int64, error)
	SetMinSeq(ctx context.Context, conversationID string, seq int64) error
}
