package controller

import (
	"context"
	"testing"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/stretchr/testify/require"
)

type effectiveReadSeqCache struct{ cache.SeqUser }

func (*effectiveReadSeqCache) SetUserReadSeq(context.Context, string, string, int64) (int64, error) {
	return 346, nil
}

func TestReadSeqSenderQueuesEffectiveProgress(t *testing.T) {
	db := &msgTransferDatabase{seqUser: &effectiveReadSeqCache{}}
	pending := map[string]int64{"user": 200}
	require.NoError(t, db.SetHasReadSeqs(context.Background(), "group", pending))
	require.EqualValues(t, 346, pending["user"])
}
