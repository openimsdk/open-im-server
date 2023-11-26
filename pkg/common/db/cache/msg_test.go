package cache

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestParallelSetMessageToCache(t *testing.T) {
	var (
		cid      = fmt.Sprintf("cid-%v", rand.Int63())
		seqFirst = rand.Int63()
		msgs     = []*sdkws.MsgData{}
	)

	for i := 0; i < 100; i++ {
		msgs = append(msgs, &sdkws.MsgData{
			Seq: seqFirst + int64(i),
		})
	}

	testParallelSetMessageToCache(t, cid, msgs)
}

func testParallelSetMessageToCache(t *testing.T, cid string, msgs []*sdkws.MsgData) {
	rdb := redis.NewClient(&redis.Options{})
	defer rdb.Close()

	cacher := msgCache{rdb: rdb}

	ret, err := cacher.ParallelSetMessageToCache(context.Background(), cid, msgs)
	assert.Nil(t, err)
	assert.Equal(t, len(msgs), ret)

	// validate
	for _, msg := range msgs {
		key := cacher.getMessageCacheKey(cid, msg.Seq)
		val, err := rdb.Exists(context.Background(), key).Result()
		assert.Nil(t, err)
		assert.EqualValues(t, 1, val)
	}
}

func TestPipeSetMessageToCache(t *testing.T) {
	var (
		cid      = fmt.Sprintf("cid-%v", rand.Int63())
		seqFirst = rand.Int63()
		msgs     = []*sdkws.MsgData{}
	)

	for i := 0; i < 100; i++ {
		msgs = append(msgs, &sdkws.MsgData{
			Seq: seqFirst + int64(i),
		})
	}

	testPipeSetMessageToCache(t, cid, msgs)
}

func testPipeSetMessageToCache(t *testing.T, cid string, msgs []*sdkws.MsgData) {
	rdb := redis.NewClient(&redis.Options{})
	defer rdb.Close()

	cacher := msgCache{rdb: rdb}

	ret, err := cacher.PipeSetMessageToCache(context.Background(), cid, msgs)
	assert.Nil(t, err)
	assert.Equal(t, len(msgs), ret)

	// validate
	for _, msg := range msgs {
		key := cacher.getMessageCacheKey(cid, msg.Seq)
		val, err := rdb.Exists(context.Background(), key).Result()
		assert.Nil(t, err)
		assert.EqualValues(t, 1, val)
	}
}

func TestGetMessagesBySeq(t *testing.T) {
	var (
		cid      = fmt.Sprintf("cid-%v", rand.Int63())
		seqFirst = rand.Int63()
		msgs     = []*sdkws.MsgData{}
	)

	seqs := []int64{}
	for i := 0; i < 100; i++ {
		msgs = append(msgs, &sdkws.MsgData{
			Seq:    seqFirst + int64(i),
			SendID: fmt.Sprintf("fake-sendid-%v", i),
		})
		seqs = append(seqs, seqFirst+int64(i))
	}

	// set data to cache
	testPipeSetMessageToCache(t, cid, msgs)

	// get data from cache with parallet mode
	testParallelGetMessagesBySeq(t, cid, seqs, msgs)

	// get data from cache with pipeline mode
	testPipeGetMessagesBySeq(t, cid, seqs, msgs)
}

func testParallelGetMessagesBySeq(t *testing.T, cid string, seqs []int64, inputMsgs []*sdkws.MsgData) {
	rdb := redis.NewClient(&redis.Options{})
	defer rdb.Close()

	cacher := msgCache{rdb: rdb}

	respMsgs, failedSeqs, err := cacher.ParallelGetMessagesBySeq(context.Background(), cid, seqs)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(failedSeqs))
	assert.Equal(t, len(respMsgs), len(seqs))

	// validate
	for idx, msg := range respMsgs {
		assert.Equal(t, msg.Seq, inputMsgs[idx].Seq)
		assert.Equal(t, msg.SendID, inputMsgs[idx].SendID)
	}
}

func testPipeGetMessagesBySeq(t *testing.T, cid string, seqs []int64, inputMsgs []*sdkws.MsgData) {
	rdb := redis.NewClient(&redis.Options{})
	defer rdb.Close()

	cacher := msgCache{rdb: rdb}

	respMsgs, failedSeqs, err := cacher.PipeGetMessagesBySeq(context.Background(), cid, seqs)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(failedSeqs))
	assert.Equal(t, len(respMsgs), len(seqs))

	// validate
	for idx, msg := range respMsgs {
		assert.Equal(t, msg.Seq, inputMsgs[idx].Seq)
		assert.Equal(t, msg.SendID, inputMsgs[idx].SendID)
	}
}

func TestGetMessagesBySeqWithEmptySeqs(t *testing.T) {
	var (
		cid            = fmt.Sprintf("cid-%v", rand.Int63())
		seqFirst int64 = 0
		msgs           = []*sdkws.MsgData{}
	)

	seqs := []int64{}
	for i := 0; i < 100; i++ {
		msgs = append(msgs, &sdkws.MsgData{
			Seq:    seqFirst + int64(i),
			SendID: fmt.Sprintf("fake-sendid-%v", i),
		})
		seqs = append(seqs, seqFirst+int64(i))
	}

	// don't set cache, only get data from cache.

	// get data from cache with parallet mode
	testParallelGetMessagesBySeqWithEmptry(t, cid, seqs, msgs)

	// get data from cache with pipeline mode
	testPipeGetMessagesBySeqWithEmptry(t, cid, seqs, msgs)
}

func testParallelGetMessagesBySeqWithEmptry(t *testing.T, cid string, seqs []int64, inputMsgs []*sdkws.MsgData) {
	rdb := redis.NewClient(&redis.Options{})
	defer rdb.Close()

	cacher := msgCache{rdb: rdb}

	respMsgs, failedSeqs, err := cacher.ParallelGetMessagesBySeq(context.Background(), cid, seqs)
	assert.Nil(t, err)
	assert.Equal(t, len(seqs), len(failedSeqs))
	assert.Equal(t, 0, len(respMsgs))
}

func testPipeGetMessagesBySeqWithEmptry(t *testing.T, cid string, seqs []int64, inputMsgs []*sdkws.MsgData) {
	rdb := redis.NewClient(&redis.Options{})
	defer rdb.Close()

	cacher := msgCache{rdb: rdb}

	respMsgs, failedSeqs, err := cacher.PipeGetMessagesBySeq(context.Background(), cid, seqs)
	assert.Equal(t, err, redis.Nil)
	assert.Equal(t, len(seqs), len(failedSeqs))
	assert.Equal(t, 0, len(respMsgs))
}

func TestGetMessagesBySeqWithLostHalfSeqs(t *testing.T) {
	var (
		cid            = fmt.Sprintf("cid-%v", rand.Int63())
		seqFirst int64 = 0
		msgs           = []*sdkws.MsgData{}
	)

	seqs := []int64{}
	for i := 0; i < 100; i++ {
		msgs = append(msgs, &sdkws.MsgData{
			Seq:    seqFirst + int64(i),
			SendID: fmt.Sprintf("fake-sendid-%v", i),
		})
		seqs = append(seqs, seqFirst+int64(i))
	}

	// Only set half the number of messages.
	testParallelSetMessageToCache(t, cid, msgs[:50])

	// get data from cache with parallet mode
	testParallelGetMessagesBySeqWithLostHalfSeqs(t, cid, seqs, msgs)

	// get data from cache with pipeline mode
	testPipeGetMessagesBySeqWithLostHalfSeqs(t, cid, seqs, msgs)
}

func testParallelGetMessagesBySeqWithLostHalfSeqs(t *testing.T, cid string, seqs []int64, inputMsgs []*sdkws.MsgData) {
	rdb := redis.NewClient(&redis.Options{})
	defer rdb.Close()

	cacher := msgCache{rdb: rdb}

	respMsgs, failedSeqs, err := cacher.ParallelGetMessagesBySeq(context.Background(), cid, seqs)
	assert.Nil(t, err)
	assert.Equal(t, len(seqs)/2, len(failedSeqs))
	assert.Equal(t, len(seqs)/2, len(respMsgs))

	for idx, msg := range respMsgs {
		assert.Equal(t, msg.Seq, seqs[idx])
	}
}

func testPipeGetMessagesBySeqWithLostHalfSeqs(t *testing.T, cid string, seqs []int64, inputMsgs []*sdkws.MsgData) {
	rdb := redis.NewClient(&redis.Options{})
	defer rdb.Close()

	cacher := msgCache{rdb: rdb}

	respMsgs, failedSeqs, err := cacher.PipeGetMessagesBySeq(context.Background(), cid, seqs)
	assert.Nil(t, err)
	assert.Equal(t, len(seqs)/2, len(failedSeqs))
	assert.Equal(t, len(seqs)/2, len(respMsgs))

	for idx, msg := range respMsgs {
		assert.Equal(t, msg.Seq, seqs[idx])
	}
}

func TestPipeDeleteMessages(t *testing.T) {
	var (
		cid      = fmt.Sprintf("cid-%v", rand.Int63())
		seqFirst = rand.Int63()
		msgs     = []*sdkws.MsgData{}
	)

	var seqs []int64
	for i := 0; i < 100; i++ {
		msgs = append(msgs, &sdkws.MsgData{
			Seq: seqFirst + int64(i),
		})
		seqs = append(seqs, msgs[i].Seq)
	}

	testPipeSetMessageToCache(t, cid, msgs)
	testPipeDeleteMessagesOK(t, cid, seqs, msgs)

	// set again
	testPipeSetMessageToCache(t, cid, msgs)
	testPipeDeleteMessagesMix(t, cid, seqs[:90], msgs)
}

func testPipeDeleteMessagesOK(t *testing.T, cid string, seqs []int64, inputMsgs []*sdkws.MsgData) {
	rdb := redis.NewClient(&redis.Options{})
	defer rdb.Close()

	cacher := msgCache{rdb: rdb}

	err := cacher.PipeDeleteMessages(context.Background(), cid, seqs)
	assert.Nil(t, err)

	// validate
	for _, msg := range inputMsgs {
		key := cacher.getMessageCacheKey(cid, msg.Seq)
		val := rdb.Exists(context.Background(), key).Val()
		assert.EqualValues(t, 0, val)
	}
}

func testPipeDeleteMessagesMix(t *testing.T, cid string, seqs []int64, inputMsgs []*sdkws.MsgData) {
	rdb := redis.NewClient(&redis.Options{})
	defer rdb.Close()

	cacher := msgCache{rdb: rdb}

	err := cacher.PipeDeleteMessages(context.Background(), cid, seqs)
	assert.Nil(t, err)

	// validate
	for idx, msg := range inputMsgs {
		key := cacher.getMessageCacheKey(cid, msg.Seq)
		val, err := rdb.Exists(context.Background(), key).Result()
		assert.Nil(t, err)
		if idx < 90 {
			assert.EqualValues(t, 0, val) // not exists
			continue
		}

		assert.EqualValues(t, 1, val) // exists
	}
}

func TestParallelDeleteMessages(t *testing.T) {
	var (
		cid      = fmt.Sprintf("cid-%v", rand.Int63())
		seqFirst = rand.Int63()
		msgs     = []*sdkws.MsgData{}
	)

	var seqs []int64
	for i := 0; i < 100; i++ {
		msgs = append(msgs, &sdkws.MsgData{
			Seq: seqFirst + int64(i),
		})
		seqs = append(seqs, msgs[i].Seq)
	}

	randSeqs := []int64{}
	for i := seqFirst + 100; i < seqFirst+200; i++ {
		randSeqs = append(randSeqs, i)
	}

	testParallelSetMessageToCache(t, cid, msgs)
	testParallelDeleteMessagesOK(t, cid, seqs, msgs)

	// set again
	testParallelSetMessageToCache(t, cid, msgs)
	testParallelDeleteMessagesMix(t, cid, seqs[:90], msgs, 90)
	testParallelDeleteMessagesOK(t, cid, seqs[90:], msgs[:90])

	// set again
	testParallelSetMessageToCache(t, cid, msgs)
	testParallelDeleteMessagesMix(t, cid, randSeqs, msgs, 0)
}

func testParallelDeleteMessagesOK(t *testing.T, cid string, seqs []int64, inputMsgs []*sdkws.MsgData) {
	rdb := redis.NewClient(&redis.Options{})
	defer rdb.Close()

	cacher := msgCache{rdb: rdb}

	err := cacher.PipeDeleteMessages(context.Background(), cid, seqs)
	assert.Nil(t, err)

	// validate
	for _, msg := range inputMsgs {
		key := cacher.getMessageCacheKey(cid, msg.Seq)
		val := rdb.Exists(context.Background(), key).Val()
		assert.EqualValues(t, 0, val)
	}
}

func testParallelDeleteMessagesMix(t *testing.T, cid string, seqs []int64, inputMsgs []*sdkws.MsgData, lessValNonExists int) {
	rdb := redis.NewClient(&redis.Options{})
	defer rdb.Close()

	cacher := msgCache{rdb: rdb}

	err := cacher.PipeDeleteMessages(context.Background(), cid, seqs)
	assert.Nil(t, err)

	// validate
	for idx, msg := range inputMsgs {
		key := cacher.getMessageCacheKey(cid, msg.Seq)
		val, err := rdb.Exists(context.Background(), key).Result()
		assert.Nil(t, err)
		if idx < lessValNonExists {
			assert.EqualValues(t, 0, val) // not exists
			continue
		}

		assert.EqualValues(t, 1, val) // exists
	}
}

func TestCleanUpOneConversationAllMsg(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{})
	defer rdb.Close()

	cacher := msgCache{rdb: rdb}
	count := 1000
	prefix := fmt.Sprintf("%v", rand.Int63())

	ids := []string{}
	for i := 0; i < count; i++ {
		id := fmt.Sprintf("%v-cid-%v", prefix, rand.Int63())
		ids = append(ids, id)

		key := cacher.allMessageCacheKey(id)
		rdb.Set(context.Background(), key, "openim", 0)
	}

	// delete 100 keys with scan.
	for i := 0; i < 100; i++ {
		pickedKey := ids[i]
		err := cacher.CleanUpOneConversationAllMsg(context.Background(), pickedKey)
		assert.Nil(t, err)

		ls, err := rdb.Keys(context.Background(), pickedKey).Result()
		assert.Nil(t, err)
		assert.Equal(t, 0, len(ls))

		rcode, err := rdb.Exists(context.Background(), pickedKey).Result()
		assert.Nil(t, err)
		assert.EqualValues(t, 0, rcode) // non-exists
	}

	sid := fmt.Sprintf("%v-cid-*", prefix)
	ls, err := rdb.Keys(context.Background(), cacher.allMessageCacheKey(sid)).Result()
	assert.Nil(t, err)
	assert.Equal(t, count-100, len(ls))

	// delete fuzzy matching keys.
	err = cacher.CleanUpOneConversationAllMsg(context.Background(), sid)
	assert.Nil(t, err)

	// don't contains keys matched `{prefix}-cid-{random}` on redis
	ls, err = rdb.Keys(context.Background(), cacher.allMessageCacheKey(sid)).Result()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(ls))
}
