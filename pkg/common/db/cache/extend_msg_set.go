package cache

import (
	"Open_IM/pkg/common/db/relation"
	"github.com/dtm-labs/rockscache"
	"time"
)

type ExtendMsgSetCache struct {
	friendDB   *relation.FriendGorm
	expireTime time.Duration
	rcClient   *rockscache.Client
}
