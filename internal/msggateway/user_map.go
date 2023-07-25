// Copyright © 2023 OpenIM. All rights reserved.
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

package msggateway

import (
	"context"
	"sync"

	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/utils"
)

type UserMap struct {
	m sync.Map
}

func newUserMap() *UserMap {
	return &UserMap{}
}

func (u *UserMap) GetAll(key string) ([]*Client, bool) {
	allClients, ok := u.m.Load(key)
	if ok {
		return allClients.([]*Client), ok
	}
	return nil, ok
}

func (u *UserMap) Get(key string, platformID int) ([]*Client, bool, bool) {
	allClients, userExisted := u.m.Load(key)
	if userExisted {
		var clients []*Client
		for _, client := range allClients.([]*Client) {
			if client.PlatformID == platformID {
				clients = append(clients, client)
			}
		}
		if len(clients) > 0 {
			return clients, userExisted, true
		}
		return clients, userExisted, false
	}
	return nil, userExisted, false
}

func (u *UserMap) Set(key string, v *Client) {
	allClients, existed := u.m.Load(key)
	if existed {
		log.ZDebug(context.Background(), "Set existed", "user_id", key, "client", *v)
		oldClients := allClients.([]*Client)
		oldClients = append(oldClients, v)
		u.m.Store(key, oldClients)
	} else {
		log.ZDebug(context.Background(), "Set not existed", "user_id", key, "client", *v)
		var clients []*Client
		clients = append(clients, v)
		u.m.Store(key, clients)
	}
}

func (u *UserMap) delete(key string, connRemoteAddr string) (isDeleteUser bool) {
	allClients, existed := u.m.Load(key)
	if existed {
		oldClients := allClients.([]*Client)
		var a []*Client
		for _, client := range oldClients {
			if client.ctx.GetRemoteAddr() != connRemoteAddr {
				a = append(a, client)
			}
		}
		if len(a) == 0 {
			u.m.Delete(key)
			return true
		} else {
			u.m.Store(key, a)
			return false
		}
	}
	return existed
}

func (u *UserMap) deleteClients(key string, clients []*Client) (isDeleteUser bool) {
	m := utils.SliceToMapAny(clients, func(c *Client) (string, struct{}) {
		return c.ctx.GetRemoteAddr(), struct{}{}
	})
	allClients, existed := u.m.Load(key)
	if existed {
		oldClients := allClients.([]*Client)
		var a []*Client
		for _, client := range oldClients {
			if _, ok := m[client.ctx.GetRemoteAddr()]; !ok {
				a = append(a, client)
			}
		}
		if len(a) == 0 {
			u.m.Delete(key)
			return true
		} else {
			u.m.Store(key, a)
			return false
		}
	}
	return existed
}

func (u *UserMap) DeleteAll(key string) {
	u.m.Delete(key)
}
