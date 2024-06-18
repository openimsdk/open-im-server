package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/openimsdk/open-im-server/v3/pkg/common/cmd"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database/mgo"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/redisutil"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const batchSize = 5

func readConfig[T any](dir string, name string) (*T, error) {
	data, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		return nil, err
	}
	var conf T
	if err := yaml.Unmarshal(data, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

func redisKey(rdb redis.UniversalClient, prefix string, fn func(ctx context.Context, key string, delKey map[string]struct{}) error) error {
	var (
		cursor uint64
		keys   []string
		err    error
	)
	ctx := context.Background()
	for {
		keys, cursor, err = rdb.Scan(ctx, cursor, prefix+"*", batchSize).Result()
		if err != nil {
			return err
		}
		delKey := make(map[string]struct{})
		if len(keys) > 0 {
			for _, key := range keys {
				if err := fn(ctx, key, delKey); err != nil {
					return err
				}
			}
		}
		if len(delKey) > 0 {
			//if err := rdb.Del(ctx, datautil.Keys(delKey)...).Err(); err != nil {
			//	return err
			//}
		}
		if cursor == 0 {
			return nil
		}
	}
}

func Main(conf string) error {
	redisConfig, err := readConfig[config.Redis](conf, cmd.RedisConfigFileName)
	if err != nil {
		return err
	}
	mongodbConfig, err := readConfig[config.Mongo](conf, cmd.MongodbConfigFileName)
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
	if _, err := mgo.NewSeqConversationMongo(mgocli.GetDB()); err != nil {
		return err
	}
	coll := mgocli.GetDB().Collection(database.SeqConversationName)
	const prefix = cachekey.MaxSeq
	return redisKey(rdb, prefix, func(ctx context.Context, key string, delKey map[string]struct{}) error {
		conversationId := strings.TrimPrefix(key, prefix)
		delKey[key] = struct{}{}
		maxValue, err := rdb.Get(ctx, key).Result()
		if err != nil {
			return err
		}
		seq, err := strconv.Atoi(maxValue)
		if err != nil {
			return fmt.Errorf("invalid max seq %s", maxValue)
		}
		if seq == 0 {
			return nil
		}
		if seq < 0 {
			return fmt.Errorf("invalid max seq %s", maxValue)
		}
		var (
			minSeq int64
			maxSeq = int64(seq)
		)
		minKey := cachekey.MinSeq + conversationId
		delKey[minKey] = struct{}{}
		minValue, err := rdb.Get(ctx, minKey).Result()
		if err == nil {
			seq, err := strconv.Atoi(minValue)
			if err != nil {
				return fmt.Errorf("invalid min seq %s", minValue)
			}
			if seq < 0 {
				return fmt.Errorf("invalid min seq %s", minValue)
			}
			minSeq = int64(seq)
		} else if !errors.Is(err, redis.Nil) {
			return err
		}
		if maxSeq < minSeq {
			return fmt.Errorf("invalid max seq %d < min seq %d", maxSeq, minSeq)
		}
		res, err := mongoutil.FindOne[*model.SeqConversation](ctx, coll, bson.M{"conversation_id": conversationId}, nil)
		if err == nil {
			if res.MaxSeq < int64(seq) {
				_, err = coll.UpdateOne(ctx, bson.M{"conversation_id": conversationId}, bson.M{"$set": bson.M{"max_seq": maxSeq, "min_seq": minSeq}})
			}
			return err
		} else if errors.Is(err, mongo.ErrNoDocuments) {
			res = &model.SeqConversation{
				ConversationID: conversationId,
				MaxSeq:         maxSeq,
				MinSeq:         minSeq,
			}
			_, err := coll.InsertOne(ctx, res)
			return err
		} else {
			return err
		}
	})
}
