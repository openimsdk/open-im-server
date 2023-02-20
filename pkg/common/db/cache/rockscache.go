package cache

import (
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"
	"github.com/dtm-labs/rockscache"
	"time"
)

const scanCount = 3000

//func (rc *RcClient) DelKeys() {
//	for _, key := range []string{"GROUP_CACHE:", "FRIEND_RELATION_CACHE", "BLACK_LIST_CACHE:", "USER_INFO_CACHE:", "GROUP_INFO_CACHE", groupOwnerIDCache, joinedGroupListCache,
//		groupMemberInfoCache, groupAllMemberInfoCache, "ALL_FRIEND_INFO_CACHE:"} {
//		fName := utils.GetSelfFuncName()
//		var cursor uint64
//		var n int
//		for {
//			var keys []string
//			var err error
//			keys, cursor, err = rc.rdb.Scan(context.Background(), cursor, key+"*", scanCount).Result()
//			if err != nil {
//				panic(err.Error())
//			}
//			n += len(keys)
//			// for each for redis cluster
//			for _, key := range keys {
//				if err = rc.rdb.Del(context.Background(), key).Err(); err != nil {
//					log.NewError("", fName, key, err.Error())
//					err = rc.rdb.Del(context.Background(), key).Err()
//					if err != nil {
//						panic(err.Error())
//					}
//				}
//			}
//			if cursor == 0 {
//				break
//			}
//		}
//	}
//}

func GetCache[T any](ctx context.Context, rcClient *rockscache.Client, key string, expire time.Duration, fn func(ctx context.Context) (T, error)) (T, error) {
	var t T
	var write bool
	v, err := rcClient.Fetch(key, expire, func() (s string, err error) {
		t, err = fn(ctx)
		if err != nil {
			return "", err
		}
		bs, err := json.Marshal(t)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		write = true
		return string(bs), nil
	})
	if err != nil {
		return t, err
	}
	if write {
		return t, nil
	}
	err = json.Unmarshal([]byte(v), &t)
	if err != nil {
		return t, utils.Wrap(err, "")
	}
	return t, nil
}

func GetCacheFor[E any, T any](ctx context.Context, list []E, fn func(ctx context.Context, item E) (T, error)) ([]T, error) {
	rs := make([]T, 0, len(list))
	for _, e := range list {
		r, err := fn(ctx, e)
		if err != nil {
			return nil, err
		}
		rs = append(rs, r)
	}
	return rs, nil
}
