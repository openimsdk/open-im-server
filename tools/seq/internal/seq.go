package internal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database/mgo"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/redisutil"
	"github.com/openimsdk/tools/utils/runtimeenv"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const StructTagName = "yaml"

const (
	MaxSeq                 = "MAX_SEQ:"
	MinSeq                 = "MIN_SEQ:"
	ConversationUserMinSeq = "CON_USER_MIN_SEQ:"
	HasReadSeq             = "HAS_READ_SEQ:"
)

const (
	batchSize             = 100
	dataVersionCollection = "data_version"
	seqKey                = "seq"
	seqVersion            = 38
)

func readConfig[T any](dir string, name string) (*T, error) {
	if runtimeenv.RuntimeEnvironment() == config.KUBERNETES {
		dir = os.Getenv(config.MountConfigFilePath)
	}
	v := viper.New()
	v.SetEnvPrefix(config.EnvPrefixMap[name])
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetConfigFile(filepath.Join(dir, name))
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var conf T
	if err := v.Unmarshal(&conf, func(config *mapstructure.DecoderConfig) {
		config.TagName = StructTagName
	}); err != nil {
		return nil, err
	}

	return &conf, nil
}

func Main(conf string, del time.Duration) error {
	redisConfig, err := readConfig[config.Redis](conf, config.RedisConfigFileName)
	if err != nil {
		return err
	}
	
	mongodbConfig, err := readConfig[config.Mongo](conf, config.MongodbConfigFileName)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	rdb, err := redisutil.NewRedisClient(ctx, redisConfig.Build())
	if err != nil {
		return err
	}
	mgocli, err := mongoutil.NewMongoDB(ctx, mongodbConfig.Build())
	if err != nil {
		return err
	}
	versionColl := mgocli.GetDB().Collection(dataVersionCollection)
	converted, err := CheckVersion(versionColl, seqKey, seqVersion)
	if err != nil {
		return err
	}
	if converted {
		fmt.Println("[seq] seq data has been converted")
		return nil
	}
	if _, err := mgo.NewSeqConversationMongo(mgocli.GetDB()); err != nil {
		return err
	}
	cSeq, err := mgo.NewSeqConversationMongo(mgocli.GetDB())
	if err != nil {
		return err
	}
	uSeq, err := mgo.NewSeqUserMongo(mgocli.GetDB())
	if err != nil {
		return err
	}
	uSpitHasReadSeq := func(id string) (conversationID string, userID string, err error) {
		// HasReadSeq + userID + ":" + conversationID
		arr := strings.Split(id, ":")
		if len(arr) != 2 || arr[0] == "" || arr[1] == "" {
			return "", "", fmt.Errorf("invalid has read seq id %s", id)
		}
		userID = arr[0]
		conversationID = arr[1]
		return
	}
	uSpitConversationUserMinSeq := func(id string) (conversationID string, userID string, err error) {
		// ConversationUserMinSeq + conversationID + "u:" + userID
		arr := strings.Split(id, "u:")
		if len(arr) != 2 || arr[0] == "" || arr[1] == "" {
			return "", "", fmt.Errorf("invalid has read seq id %s", id)
		}
		conversationID = arr[0]
		userID = arr[1]
		return
	}

	ts := []*taskSeq{
		{
			Prefix: MaxSeq,
			GetSeq: cSeq.GetMaxSeq,
			SetSeq: cSeq.SetMaxSeq,
		},
		{
			Prefix: MinSeq,
			GetSeq: cSeq.GetMinSeq,
			SetSeq: cSeq.SetMinSeq,
		},
		{
			Prefix: HasReadSeq,
			GetSeq: func(ctx context.Context, id string) (int64, error) {
				conversationID, userID, err := uSpitHasReadSeq(id)
				if err != nil {
					return 0, err
				}
				return uSeq.GetUserReadSeq(ctx, conversationID, userID)
			},
			SetSeq: func(ctx context.Context, id string, seq int64) error {
				conversationID, userID, err := uSpitHasReadSeq(id)
				if err != nil {
					return err
				}
				return uSeq.SetUserReadSeq(ctx, conversationID, userID, seq)
			},
		},
		{
			Prefix: ConversationUserMinSeq,
			GetSeq: func(ctx context.Context, id string) (int64, error) {
				conversationID, userID, err := uSpitConversationUserMinSeq(id)
				if err != nil {
					return 0, err
				}
				return uSeq.GetUserMinSeq(ctx, conversationID, userID)
			},
			SetSeq: func(ctx context.Context, id string, seq int64) error {
				conversationID, userID, err := uSpitConversationUserMinSeq(id)
				if err != nil {
					return err
				}
				return uSeq.SetUserMinSeq(ctx, conversationID, userID, seq)
			},
		},
	}

	cancel()
	ctx = context.Background()

	var wg sync.WaitGroup
	wg.Add(len(ts))

	for i := range ts {
		go func(task *taskSeq) {
			defer wg.Done()
			err := seqRedisToMongo(ctx, rdb, task.GetSeq, task.SetSeq, task.Prefix, del, &task.Count)
			task.End = time.Now()
			task.Error = err
		}(ts[i])
	}
	start := time.Now()
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	var buf bytes.Buffer

	printTaskInfo := func(now time.Time) {
		buf.Reset()
		buf.WriteString(now.Format(time.DateTime))
		buf.WriteString(" \n")
		for i := range ts {
			task := ts[i]
			if task.Error == nil {
				if task.End.IsZero() {
					buf.WriteString(fmt.Sprintf("[%s] converting %s* count %d", now.Sub(start), task.Prefix, atomic.LoadInt64(&task.Count)))
				} else {
					buf.WriteString(fmt.Sprintf("[%s] success %s* count %d", task.End.Sub(start), task.Prefix, atomic.LoadInt64(&task.Count)))
				}
			} else {
				buf.WriteString(fmt.Sprintf("[%s] failed %s* count %d error %s", task.End.Sub(start), task.Prefix, atomic.LoadInt64(&task.Count), task.Error))
			}
			buf.WriteString("\n")
		}
		fmt.Println(buf.String())
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case s := <-sigs:
			return fmt.Errorf("exit by signal %s", s)
		case <-done:
			errs := make([]error, 0, len(ts))
			for i := range ts {
				task := ts[i]
				if task.Error != nil {
					errs = append(errs, fmt.Errorf("seq %s failed %w", task.Prefix, task.Error))
				}
			}
			if len(errs) > 0 {
				return errors.Join(errs...)
			}
			printTaskInfo(time.Now())
			if err := SetVersion(versionColl, seqKey, seqVersion); err != nil {
				return fmt.Errorf("set mongodb seq version %w", err)
			}
			return nil
		case now := <-ticker.C:
			printTaskInfo(now)
		}
	}
}

type taskSeq struct {
	Prefix string
	Count  int64
	Error  error
	End    time.Time
	GetSeq func(ctx context.Context, id string) (int64, error)
	SetSeq func(ctx context.Context, id string, seq int64) error
}

func seqRedisToMongo(ctx context.Context, rdb redis.UniversalClient, getSeq func(ctx context.Context, id string) (int64, error), setSeq func(ctx context.Context, id string, seq int64) error, prefix string, delAfter time.Duration, count *int64) error {
	var (
		cursor uint64
		keys   []string
		err    error
	)
	for {
		keys, cursor, err = rdb.Scan(ctx, cursor, prefix+"*", batchSize).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			for _, key := range keys {
				seqStr, err := rdb.Get(ctx, key).Result()
				if err != nil {
					return fmt.Errorf("redis get %s failed %w", key, err)
				}
				seq, err := strconv.Atoi(seqStr)
				if err != nil {
					return fmt.Errorf("invalid %s seq %s", key, seqStr)
				}
				if seq < 0 {
					return fmt.Errorf("invalid %s seq %s", key, seqStr)
				}
				id := strings.TrimPrefix(key, prefix)
				redisSeq := int64(seq)
				mongoSeq, err := getSeq(ctx, id)
				if err != nil {
					return fmt.Errorf("get mongo seq %s failed %w", key, err)
				}
				if mongoSeq < redisSeq {
					if err := setSeq(ctx, id, redisSeq); err != nil {
						return fmt.Errorf("set mongo seq %s failed %w", key, err)
					}
				}
				if delAfter > 0 {
					if err := rdb.Expire(ctx, key, delAfter).Err(); err != nil {
						return fmt.Errorf("redis expire key %s failed %w", key, err)
					}
				} else {
					if err := rdb.Del(ctx, key).Err(); err != nil {
						return fmt.Errorf("redis del key %s failed %w", key, err)
					}
				}
				atomic.AddInt64(count, 1)
			}
		}
		if cursor == 0 {
			return nil
		}
	}
}

func CheckVersion(coll *mongo.Collection, key string, currentVersion int) (converted bool, err error) {
	type VersionTable struct {
		Key   string `bson:"key"`
		Value string `bson:"value"`
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	res, err := mongoutil.FindOne[VersionTable](ctx, coll, bson.M{"key": key})
	if err == nil {
		ver, err := strconv.Atoi(res.Value)
		if err != nil {
			return false, fmt.Errorf("version %s parse error %w", res.Value, err)
		}
		if ver >= currentVersion {
			return true, nil
		}
		return false, nil
	} else if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	} else {
		return false, err
	}
}

func SetVersion(coll *mongo.Collection, key string, version int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	option := options.Update().SetUpsert(true)
	filter := bson.M{"key": key}
	update := bson.M{"$set": bson.M{"key": key, "value": strconv.Itoa(version)}}
	return mongoutil.UpdateOne(ctx, coll, filter, update, false, option)
}
