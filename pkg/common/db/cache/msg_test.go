// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cache

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/openimsdk/protocol/sdkws"
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
