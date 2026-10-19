package redis

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

type readSeqStore struct {
	database.SeqUser
	read      func(context.Context, string, string) (int64, error)
	readBatch func(context.Context, string, []string) (map[string]int64, error)
	write     func(context.Context, string, string, int64) error
}

func (s *readSeqStore) SetUserReadSeq(ctx context.Context, conv, user string, seq int64) error {
	return s.write(ctx, conv, user, seq)
}

func (s *readSeqStore) GetUserReadSeq(ctx context.Context, conv, user string) (int64, error) {
	if s.read != nil {
		return s.read(ctx, conv, user)
	}
	return 100, nil
}

func (s *readSeqStore) GetUserReadSeqs(ctx context.Context, user string, convs []string) (map[string]int64, error) {
	if s.readBatch != nil {
		return s.readBatch(ctx, user, convs)
	}
	result := make(map[string]int64, len(convs))
	for _, conv := range convs {
		result[conv] = 100
	}
	return result, nil
}

func readSeqTestCache(t *testing.T, store database.SeqUser) (*seqUserCacheRedis, *redis.Client, *miniredis.Miniredis) {
	t.Helper()
	server := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: server.Addr()})
	t.Cleanup(func() { _ = rdb.Close() })
	return NewSeqUserCacheRedis(rdb, store).(*seqUserCacheRedis), rdb, server
}

// Delay the lower writer at the Redis command boundary, after its cache read.
// This covers both the old HSET and the atomic-script implementation.
type pauseReadSeqWrite struct {
	started chan struct{}
	resume  chan struct{}
	once    sync.Once
}

func (h *pauseReadSeqWrite) DialHook(next redis.DialHook) redis.DialHook { return next }
func (h *pauseReadSeqWrite) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return next
}
func (h *pauseReadSeqWrite) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		if cmd.Name() == "hset" || cmd.Name() == "evalsha" || cmd.Name() == "eval" {
			for _, arg := range cmd.Args() {
				if fmt.Sprint(arg) == "200" {
					h.once.Do(func() {
						close(h.started)
						select {
						case <-h.resume:
						case <-ctx.Done():
						}
					})
					break
				}
			}
		}
		return next(ctx, cmd)
	}
}

// Keep a compile-time check that this hook remains usable by go-redis.
var _ redis.Hook = (*pauseReadSeqWrite)(nil)

func TestReadSeqConcurrentWritersNeverRegress(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s, rdb, _ := readSeqTestCache(t, &readSeqStore{})
	key := s.getSeqUserReadSeqKey("group", "user")
	require.NoError(t, rdb.HSet(ctx, key, "value", "100").Err())
	hook := &pauseReadSeqWrite{started: make(chan struct{}), resume: make(chan struct{})}
	rdb.AddHook(hook)
	done := make(chan error, 1)
	go func() { _, err := s.SetUserReadSeq(ctx, "group", "user", 200); done <- err }()
	select {
	case <-hook.started:
	case <-ctx.Done():
		t.Fatal("lower writer never reached Redis")
	}
	highSeq, highErr := s.SetUserReadSeq(ctx, "group", "user", 346)
	close(hook.resume)
	require.NoError(t, highErr)
	require.EqualValues(t, 346, highSeq)
	require.NoError(t, <-done)
	require.Equal(t, "346", rdb.HGet(ctx, key, "value").Val())
}

func TestReadSeqRefillDoesNotOverwriteNewerProgress(t *testing.T) {
	for _, batch := range []bool{false, true} {
		t.Run(fmt.Sprintf("batch=%v", batch), func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			started, resume := make(chan struct{}), make(chan struct{})
			var once sync.Once
			wait := func() {
				once.Do(func() {
					close(started)
					select {
					case <-resume:
					case <-ctx.Done():
					}
				})
			}
			store := &readSeqStore{}
			if batch {
				store.readBatch = func(context.Context, string, []string) (map[string]int64, error) {
					wait()
					return map[string]int64{"group": 100}, nil
				}
			} else {
				store.read = func(context.Context, string, string) (int64, error) { wait(); return 100, nil }
			}
			s, rdb, _ := readSeqTestCache(t, store)
			type result struct {
				seq int64
				err error
			}
			done := make(chan result, 1)
			go func() {
				if batch {
					values, err := s.GetUserReadSeqs(ctx, "user", []string{"group"})
					done <- result{values["group"], err}
				} else {
					seq, err := s.GetUserReadSeq(ctx, "group", "user")
					done <- result{seq, err}
				}
			}()
			select {
			case <-started:
			case <-ctx.Done():
				t.Fatal("refill never reached MongoDB")
			}
			// Seed the new baseline as another request completes its refill, then
			// advance it through the real writer while the older query is pending.
			key := s.getSeqUserReadSeqKey("group", "user")
			require.NoError(t, rdb.HSet(ctx, key, "value", "100").Err())
			require.NoError(t, rdb.HDel(ctx, key, "lockUntil", "lockOwner").Err())
			writeErr := s.SetUserReadSeqs(ctx, "user", map[string]int64{"group": 346})
			close(resume)
			require.NoError(t, writeErr)
			got := <-done
			require.NoError(t, got.err)
			require.EqualValues(t, 346, got.seq, "the response must include the concurrent advancement")
			require.Equal(t, "346", rdb.HGet(ctx, key, "value").Val())
		})
	}
}

func TestReadSeqBatchPreservesMongoBaseline(t *testing.T) {
	store := &readSeqStore{readBatch: func(context.Context, string, []string) (map[string]int64, error) {
		return map[string]int64{"group": 346}, nil
	}}
	s, rdb, _ := readSeqTestCache(t, store)
	ctx := context.Background()
	require.NoError(t, s.SetUserReadSeqs(ctx, "user", map[string]int64{"group": 200}))
	require.Equal(t, "346", rdb.HGet(ctx, s.getSeqUserReadSeqKey("group", "user"), "value").Val())
}

func TestReadSeqBatchNeverRegresses(t *testing.T) {
	s, rdb, _ := readSeqTestCache(t, &readSeqStore{})
	ctx := context.Background()
	key := s.getSeqUserReadSeqKey("group", "user")
	require.NoError(t, rdb.HSet(ctx, key, "value", "346").Err())
	require.NoError(t, s.SetUserReadSeqs(ctx, "user", map[string]int64{"group": 200}))
	require.Equal(t, "346", rdb.HGet(ctx, key, "value").Val())
}

func TestReadSeqInt64Precision(t *testing.T) {
	s, rdb, server := readSeqTestCache(t, &readSeqStore{})
	ctx := context.Background()
	const high int64 = 9223372036854775807
	key := s.getSeqUserReadSeqKey("group", "user")
	require.NoError(t, rdb.HSet(ctx, key, "value", strconv.FormatInt(high, 10)).Err())
	require.NoError(t, s.SetUserReadSeqs(ctx, "user", map[string]int64{"group": high - 1}))
	require.Equal(t, strconv.FormatInt(high, 10), rdb.HGet(ctx, key, "value").Val())
	require.Equal(t, s.readExpireTime, server.TTL(key))
}

func TestReadSeqFencesLegacyRocksCacheRefill(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s, rdb, _ := readSeqTestCache(t, &readSeqStore{})
	key := s.getSeqUserReadSeqKey("group", "user")
	started, resume := make(chan struct{}), make(chan struct{})
	done := make(chan error, 1)
	go func() {
		_, err := s.rocks.GetClient().Fetch2(ctx, key, time.Hour, func() (string, error) {
			close(started)
			select {
			case <-resume:
			case <-ctx.Done():
			}
			return "100", nil
		})
		done <- err
	}()
	select {
	case <-started:
	case <-ctx.Done():
		t.Fatal("legacy refill did not start")
	}
	effective, writeErr := s.SetUserReadSeq(ctx, "group", "user", 346)
	close(resume)
	require.NoError(t, writeErr)
	require.EqualValues(t, 346, effective)
	require.NoError(t, <-done)
	require.Equal(t, "346", rdb.HGet(ctx, key, "value").Val())
}

func TestReadSeqWithoutRedis(t *testing.T) {
	var stored int64 = 346
	store := &readSeqStore{
		read:  func(context.Context, string, string) (int64, error) { return stored, nil },
		write: func(_ context.Context, _, _ string, seq int64) error { stored = max(stored, seq); return nil },
	}
	s := NewSeqUserCacheRedis(nil, store)
	ctx := context.Background()
	effective, err := s.SetUserReadSeq(ctx, "group", "user", 200)
	require.NoError(t, err)
	require.EqualValues(t, 346, effective)
	require.NoError(t, s.SetUserReadSeqs(ctx, "user", map[string]int64{"group": 400}))
	require.EqualValues(t, 400, stored)
}

func TestReadSeqMissingAndErrorCases(t *testing.T) {
	t.Run("missing rows fill zero and duplicate IDs are deduplicated", func(t *testing.T) {
		store := &readSeqStore{readBatch: func(_ context.Context, _ string, convs []string) (map[string]int64, error) {
			require.Equal(t, []string{"group"}, convs)
			return nil, nil
		}}
		s, _, _ := readSeqTestCache(t, store)
		got, err := s.GetUserReadSeqs(context.Background(), "user", []string{"group", "group"})
		require.NoError(t, err)
		require.Equal(t, map[string]int64{"group": 0}, got)
	})
	t.Run("MongoDB errors do not populate the cache", func(t *testing.T) {
		failure := fmt.Errorf("database unavailable")
		store := &readSeqStore{read: func(context.Context, string, string) (int64, error) { return 0, failure }}
		s, rdb, _ := readSeqTestCache(t, store)
		_, err := s.SetUserReadSeq(context.Background(), "group", "user", 346)
		require.ErrorIs(t, err, failure)
		require.Zero(t, rdb.Exists(context.Background(), s.getSeqUserReadSeqKey("group", "user")).Val())
	})
	t.Run("negative writes are rejected", func(t *testing.T) {
		s, _, _ := readSeqTestCache(t, &readSeqStore{})
		_, err := s.SetUserReadSeq(context.Background(), "group", "user", -1)
		require.Error(t, err)
		require.Error(t, s.SetUserReadSeqs(context.Background(), "user", map[string]int64{"group": -1}))
	})
}
