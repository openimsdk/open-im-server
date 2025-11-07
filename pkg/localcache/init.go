package localcache

import (
	"strings"
	"sync"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
)

var (
	once      sync.Once
	subscribe map[string][]string
)

func InitLocalCache(localCache *config.LocalCache) {
	once.Do(func() {
		list := []struct {
			Local config.CacheConfig
			Keys  []string
		}{
			{
				Local: localCache.Auth,
				Keys:  []string{cachekey.UidPidToken},
			},
			{
				Local: localCache.User,
				Keys:  []string{cachekey.UserInfoKey, cachekey.UserGlobalRecvMsgOptKey},
			},
			{
				Local: localCache.Group,
				Keys:  []string{cachekey.GroupMemberIDsKey, cachekey.GroupInfoKey, cachekey.GroupMemberInfoKey},
			},
			{
				Local: localCache.Friend,
				Keys:  []string{cachekey.FriendIDsKey, cachekey.BlackIDsKey},
			},
			{
				Local: localCache.Conversation,
				Keys:  []string{cachekey.ConversationKey, cachekey.ConversationIDsKey, cachekey.ConversationNotReceiveMessageUserIDsKey},
			},
		}
		subscribe = make(map[string][]string)
		for _, v := range list {
			if v.Local.Enable() {
				subscribe[v.Local.Topic] = v.Keys
			}
		}
	})
}

func GetPublishKeysByTopic(topics []string, keys []string) map[string][]string {
	keysByTopic := make(map[string][]string)
	for _, topic := range topics {
		keysByTopic[topic] = []string{}
	}

	for _, key := range keys {
		for _, topic := range topics {
			prefixes, ok := subscribe[topic]
			if !ok {
				continue
			}
			for _, prefix := range prefixes {
				if strings.HasPrefix(key, prefix) {
					keysByTopic[topic] = append(keysByTopic[topic], key)
					break
				}
			}
		}
	}

	return keysByTopic
}
