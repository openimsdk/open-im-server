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

	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
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

// Set adds a client to the map.
func (u *UserMap) Set(key string, v *Client) {
	allClients, existed := u.m.Load(key)
	if existed {
		log.ZDebug(context.Background(), "Set existed", "user_id", key, "client_user_id", v.UserID)
		oldClients := allClients.([]*Client)
		oldClients = append(oldClients, v)
		u.m.Store(key, oldClients)
	} else {
		log.ZDebug(context.Background(), "Set not existed", "user_id", key, "client_user_id", v.UserID)

		var clients []*Client
		clients = append(clients, v)
		u.m.Store(key, clients)
	}
}

func (u *UserMap) delete(key string, connRemoteAddr string) (isDeleteUser bool) {
	// Attempt to load the clients associated with the key.
	allClients, existed := u.m.Load(key)
	if !existed {
		// Return false immediately if the key does not exist.
		return false
	}

	// Convert allClients to a slice of *Client.
	oldClients := allClients.([]*Client)
	var remainingClients []*Client
	for _, client := range oldClients {
		// Keep clients that do not match the connRemoteAddr.
		if client.ctx.GetRemoteAddr() != connRemoteAddr {
			remainingClients = append(remainingClients, client)
		}
	}

	// If no clients remain after filtering, delete the key from the map.
	if len(remainingClients) == 0 {
		u.m.Delete(key)
		return true
	}

	// Otherwise, update the key with the remaining clients.
	u.m.Store(key, remainingClients)
	return false
}

func (u *UserMap) deleteClients(key string, clients []*Client) (isDeleteUser bool) {
	m := datautil.SliceToMapAny(clients, func(c *Client) (string, struct{}) {
		return c.ctx.GetRemoteAddr(), struct{}{}
	})
	allClients, existed := u.m.Load(key)
	if !existed {
		// If the key doesn't exist, return false.
		return false
	}

	// Filter out clients that are in the deleteMap.
	oldClients := allClients.([]*Client)
	var remainingClients []*Client
	for _, client := range oldClients {
		if _, shouldBeDeleted := m[client.ctx.GetRemoteAddr()]; !shouldBeDeleted {
			remainingClients = append(remainingClients, client)
		}
	}

	// Update or delete the key based on the remaining clients.
	if len(remainingClients) == 0 {
		u.m.Delete(key)
		return true
	}

	u.m.Store(key, remainingClients)
	return false
}

func (u *UserMap) DeleteAll(key string) {
	u.m.Delete(key)
}
