package cache

import (
	"Open_IM/pkg/common/db/relation"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"
	"github.com/dtm-labs/rockscache"
	"time"
)

const (
	blackIDsKey     = "BLACK_IDS:"
	blackExpireTime = time.Second * 60 * 60 * 12
)

type BlackCache struct {
	blackDB    *relation.Black
	expireTime time.Duration
	rcClient   *rockscache.Client
}

func NewBlackCache(blackDB *relation.Black) *BlackCache {
	return &BlackCache{
		blackDB:    nil,
		expireTime: 0,
		rcClient:   nil,
	}
}

func (b *BlackCache) getBlackIDsKey(ownerUserID string) string {
	return blackIDsKey + ownerUserID
}

func (b *BlackCache) GetBlackIDs(ctx context.Context, userID string) (blackIDs []string, err error) {
	getBlackIDList := func() (string, error) {
		blackIDs, err := b.blackDB.GetBlackIDs(ctx, userID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(blackIDs)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userID", userID, "blackIDList", blackIDs)
	}()
	blackIDListStr, err := b.rcClient.Fetch(blackListCache+userID, time.Second*30*60, getBlackIDList)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	err = json.Unmarshal([]byte(blackIDListStr), &blackIDs)
	return blackIDs, utils.Wrap(err, "")
}

func (b *BlackCache) DelBlackIDListFromCache(ctx context.Context, userID string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ctx", ctx)
	}()
	return b.rcClient.TagAsDeleted(blackListCache + userID)
}
