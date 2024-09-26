package cache

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
)

type SeqConversationCache interface {
	Malloc(ctx context.Context, conversationID string, size int64) (int64, error)
	GetMaxSeq(ctx context.Context, conversationID string) (int64, error)
	SetMinSeq(ctx context.Context, conversationID string, seq int64) error
	GetMinSeq(ctx context.Context, conversationID string) (int64, error)
	GetMaxSeqs(ctx context.Context, conversationIDs []string) (map[string]int64, error)
	SetMinSeqs(ctx context.Context, seqs map[string]int64) error
	GetCacheMaxSeqWithTime(ctx context.Context, conversationIDs []string) (map[string]database.SeqTime, error)
	GetMaxSeqsWithTime(ctx context.Context, conversationIDs []string) (map[string]database.SeqTime, error)
	GetMaxSeqWithTime(ctx context.Context, conversationID string) (database.SeqTime, error)
}
