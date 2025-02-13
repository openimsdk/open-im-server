package mcache

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
)

func NewSeqConversationCache(sc database.SeqConversation) cache.SeqConversationCache {
	return &seqConversationCache{
		sc: sc,
	}
}

type seqConversationCache struct {
	sc database.SeqConversation
}

func (x *seqConversationCache) Malloc(ctx context.Context, conversationID string, size int64) (int64, error) {
	return x.sc.Malloc(ctx, conversationID, size)
}

func (x *seqConversationCache) SetMinSeq(ctx context.Context, conversationID string, seq int64) error {
	return x.sc.SetMinSeq(ctx, conversationID, seq)
}

func (x *seqConversationCache) GetMinSeq(ctx context.Context, conversationID string) (int64, error) {
	return x.sc.GetMinSeq(ctx, conversationID)
}

func (x *seqConversationCache) GetMaxSeqs(ctx context.Context, conversationIDs []string) (map[string]int64, error) {
	res := make(map[string]int64)
	for _, conversationID := range conversationIDs {
		seq, err := x.GetMinSeq(ctx, conversationID)
		if err != nil {
			return nil, err
		}
		res[conversationID] = seq
	}
	return res, nil
}

func (x *seqConversationCache) GetMaxSeqsWithTime(ctx context.Context, conversationIDs []string) (map[string]database.SeqTime, error) {
	res := make(map[string]database.SeqTime)
	for _, conversationID := range conversationIDs {
		seq, err := x.GetMinSeq(ctx, conversationID)
		if err != nil {
			return nil, err
		}
		res[conversationID] = database.SeqTime{Seq: seq}
	}
	return res, nil
}

func (x *seqConversationCache) GetMaxSeq(ctx context.Context, conversationID string) (int64, error) {
	return x.sc.GetMaxSeq(ctx, conversationID)
}

func (x *seqConversationCache) GetMaxSeqWithTime(ctx context.Context, conversationID string) (database.SeqTime, error) {
	seq, err := x.GetMinSeq(ctx, conversationID)
	if err != nil {
		return database.SeqTime{}, err
	}
	return database.SeqTime{Seq: seq}, nil
}

func (x *seqConversationCache) SetMinSeqs(ctx context.Context, seqs map[string]int64) error {
	for conversationID, seq := range seqs {
		if err := x.sc.SetMinSeq(ctx, conversationID, seq); err != nil {
			return err
		}
	}
	return nil
}

func (x *seqConversationCache) GetCacheMaxSeqWithTime(ctx context.Context, conversationIDs []string) (map[string]database.SeqTime, error) {
	return x.GetMaxSeqsWithTime(ctx, conversationIDs)
}
