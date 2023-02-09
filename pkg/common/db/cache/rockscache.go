package cache

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/relation"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"
	"math/big"
	"sort"
	"strconv"
	"time"
)

const (
	//userInfoCache       = "USER_INFO_CACHE:"
	//friendRelationCache = "FRIEND_RELATION_CACHE:"
	blackListCache = "BLACK_LIST_CACHE:"
	//groupCache          = "GROUP_CACHE:"
	//groupInfoCache            = "GROUP_INFO_CACHE:"
	//groupOwnerIDCache         = "GROUP_OWNER_ID:"
	//joinedGroupListCache      = "JOINED_GROUP_LIST_CACHE:"
	//groupMemberInfoCache      = "GROUP_MEMBER_INFO_CACHE:"
	//groupAllMemberInfoCache   = "GROUP_ALL_MEMBER_INFO_CACHE:"
	//allFriendInfoCache = "ALL_FRIEND_INFO_CACHE:"
	//joinedSuperGroupListCache = "JOINED_SUPER_GROUP_LIST_CACHE:"
	//groupMemberListHashCache  = "GROUP_MEMBER_LIST_HASH_CACHE:"
	//groupMemberNumCache       = "GROUP_MEMBER_NUM_CACHE:"
	conversationCache       = "CONVERSATION_CACHE:"
	conversationIDListCache = "CONVERSATION_ID_LIST_CACHE:"

	extendMsgSetCache = "EXTEND_MSG_SET_CACHE:"
	extendMsgCache    = "EXTEND_MSG_CACHE:"
)

const scanCount = 3000
const RandomExpireAdjustment = 0.2

func (rc *RcClient) DelKeys() {
	for _, key := range []string{"GROUP_CACHE:", "FRIEND_RELATION_CACHE", "BLACK_LIST_CACHE:", "USER_INFO_CACHE:", "GROUP_INFO_CACHE", groupOwnerIDCache, joinedGroupListCache,
		groupMemberInfoCache, groupAllMemberInfoCache, "ALL_FRIEND_INFO_CACHE:"} {
		fName := utils.GetSelfFuncName()
		var cursor uint64
		var n int
		for {
			var keys []string
			var err error
			keys, cursor, err = rc.rdb.Scan(context.Background(), cursor, key+"*", scanCount).Result()
			if err != nil {
				panic(err.Error())
			}
			n += len(keys)
			// for each for redis cluster
			for _, key := range keys {
				if err = rc.rdb.Del(context.Background(), key).Err(); err != nil {
					log.NewError("", fName, key, err.Error())
					err = rc.rdb.Del(context.Background(), key).Err()
					if err != nil {
						panic(err.Error())
					}
				}
			}
			if cursor == 0 {
				break
			}
		}
	}
}
