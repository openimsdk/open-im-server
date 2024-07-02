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
	"time"

	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
)

func newUserMap1() UMap {
	return &UserMap{
		ch: make(chan UserState, 1024),
	}
}

type UserPlatform1 struct {
	Time    time.Time
	Clients []*Client
}

func (u *UserPlatform1) PlatformIDs() []int32 {
	if len(u.Clients) == 0 {
		return nil
	}
	platformIDs := make([]int32, 0, len(u.Clients))
	for _, client := range u.Clients {
		platformIDs = append(platformIDs, int32(client.PlatformID))
	}
	return platformIDs
}

type UserMap struct {
	m  sync.Map
	ch chan UserState
}

func (u *UserMap) UserState() <-chan UserState {
	return u.ch
}

func (u *UserMap) GetAllUserStatus(deadline time.Time) []UserState {
	var result []UserState
	u.m.Range(func(key, value any) bool {
		client := value.(*UserPlatform1)
		if client.Time.Before(deadline) {
			return true
		}
		client.Time = time.Now()
		us := UserState{
			UserID: key.(string),
			Online: make([]int32, 0, len(client.Clients)),
		}
		for _, c := range client.Clients {
			us.Online = append(us.Online, int32(c.PlatformID))
		}
		return true
	})
	return result
}

func (u *UserMap) push(userID string, userPlatform *UserPlatform1, offline []int32) bool {
	select {
	case u.ch <- UserState{UserID: userID, Online: userPlatform.PlatformIDs(), Offline: offline}:
		userPlatform.Time = time.Now()
		return true
	default:
		return false
	}
}

func (u *UserMap) GetAll(key string) ([]*Client, bool) {
	allClients, ok := u.m.Load(key)
	if ok {
		return allClients.(*UserPlatform1).Clients, ok
	}
	return nil, ok
}

func (u *UserMap) Get(key string, platformID int) ([]*Client, bool, bool) {
	allClients, userExisted := u.m.Load(key)
	if userExisted {
		var clients []*Client
		for _, client := range allClients.(*UserPlatform1).Clients {
			if client.PlatformID == platformID {
				clients = append(clients, client)
			}
		}
		if len(clients) > 0 {
			return clients, true, true
		}
		return clients, true, false
	}
	return nil, false, false
}

// Set adds a client to the map.
func (u *UserMap) Set(key string, v *Client) {
	allClients, existed := u.m.Load(key)
	if existed {
		log.ZDebug(context.Background(), "Set existed", "user_id", key, "client_user_id", v.UserID)
		oldClients := allClients.(*UserPlatform1)
		oldClients.Time = time.Now()
		oldClients.Clients = append(oldClients.Clients, v)
		u.push(key, oldClients, nil)
	} else {
		log.ZDebug(context.Background(), "Set not existed", "user_id", key, "client_user_id", v.UserID)
		cli := &UserPlatform1{
			Time:    time.Now(),
			Clients: []*Client{v},
		}
		u.m.Store(key, cli)
		u.push(key, cli, nil)
	}

}

func (u *UserMap) DeleteClients(key string, clients []*Client) (isDeleteUser bool) {
	m := datautil.SliceToMapAny(clients, func(c *Client) (string, struct{}) {
		return c.ctx.GetRemoteAddr(), struct{}{}
	})
	allClients, existed := u.m.Load(key)
	if !existed {
		// If the key doesn't exist, return false.
		return false
	}

	// Filter out clients that are in the deleteMap.
	oldClients := allClients.(*UserPlatform1)
	var (
		remainingClients []*Client
		offline          []int32
	)
	for _, client := range oldClients.Clients {
		if _, shouldBeDeleted := m[client.ctx.GetRemoteAddr()]; !shouldBeDeleted {
			remainingClients = append(remainingClients, client)
		} else {
			offline = append(offline, int32(client.PlatformID))
		}
	}

	oldClients.Clients = remainingClients
	defer u.push(key, oldClients, offline)
	// Update or delete the key based on the remaining clients.
	if len(remainingClients) == 0 {
		u.m.Delete(key)
		return true
	}

	return false
}
