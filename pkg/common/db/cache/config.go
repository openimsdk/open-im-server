package cache

import (
	"github.com/openimsdk/open-im-server/v3/pkg/common/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"strings"
	"sync"
)

var (
	once      sync.Once
	subscribe map[string][]string
)

func getPublishKey(topic string, key []string) []string {
	if topic == "" || len(key) == 0 {
		return nil
	}
	once.Do(func() {
		list := []struct {
			Local config.LocalCache
			Keys  []string
		}{
			{
				Local: config.Config.LocalCache.User,
				Keys:  []string{cachekey.UserInfoKey, cachekey.UserGlobalRecvMsgOptKey},
			},
			{
				Local: config.Config.LocalCache.Group,
				Keys:  []string{cachekey.GroupMemberIDsKey, cachekey.GroupInfoKey, cachekey.GroupMemberInfoKey},
			},
			{
				Local: config.Config.LocalCache.Friend,
				Keys:  []string{cachekey.FriendIDsKey, cachekey.BlackIDsKey},
			},
			{
				Local: config.Config.LocalCache.Conversation,
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
	prefix, ok := subscribe[topic]
	if !ok {
		return nil
	}
	res := make([]string, 0, len(key))
	for _, k := range key {
		var exist bool
		for _, p := range prefix {
			if strings.HasPrefix(k, p) {
				exist = true
				break
			}
		}
		if exist {
			res = append(res, k)
		}
	}
	return res
}
