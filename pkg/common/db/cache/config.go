// Copyright © 2024 OpenIM. All rights reserved.
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
	"strings"
	"sync"

	"github.com/openimsdk/open-im-server/v3/pkg/common/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
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
