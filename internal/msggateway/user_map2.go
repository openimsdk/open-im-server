package msggateway

import (
	"sync"
	"time"
)

type UMap interface {
	GetAll(userID string) ([]*Client, bool)
	Get(userID string, platformID int) ([]*Client, bool, bool)
	Set(userID string, v *Client)
	Delete(userID string, connRemoteAddr string) (isDeleteUser bool)
	DeleteClients(userID string, clients []*Client) (isDeleteUser bool)
	UserState() <-chan UserState
	GetAllUserStatus(deadline, nowtime time.Time) []UserState
}

var _ UMap = (*UserMap2)(nil)

type UserPlatform struct {
	Time    time.Time
	Clients map[string]*Client
}

func (u *UserPlatform) PlatformIDs() []int32 {
	if len(u.Clients) == 0 {
		return nil
	}
	platformIDs := make([]int32, 0, len(u.Clients))
	for _, client := range u.Clients {
		platformIDs = append(platformIDs, int32(client.PlatformID))
	}
	return platformIDs
}

func newUserMap() UMap {
	return &UserMap2{
		data: make(map[string]*UserPlatform),
		ch:   make(chan UserState, 10000),
	}
}

type UserMap2 struct {
	lock sync.RWMutex
	data map[string]*UserPlatform
	ch   chan UserState
}

func (u *UserMap2) push(userID string, userPlatform *UserPlatform, offline []int32) bool {
	select {
	case u.ch <- UserState{UserID: userID, Online: userPlatform.PlatformIDs(), Offline: offline}:
		userPlatform.Time = time.Now()
		return true
	default:
		return false
	}
}

func (u *UserMap2) GetAll(userID string) ([]*Client, bool) {
	u.lock.RLock()
	defer u.lock.RUnlock()
	result, ok := u.data[userID]
	if !ok {
		return nil, false
	}
	clients := make([]*Client, 0, len(result.Clients))
	for _, client := range result.Clients {
		clients = append(clients, client)
	}
	return clients, true
}

func (u *UserMap2) Get(userID string, platformID int) ([]*Client, bool, bool) {
	u.lock.RLock()
	defer u.lock.RUnlock()
	result, ok := u.data[userID]
	if !ok {
		return nil, false, false
	}
	var clients []*Client
	for _, client := range result.Clients {
		if client.PlatformID == platformID {
			clients = append(clients, client)
		}
	}
	return clients, true, len(clients) > 0
}

func (u *UserMap2) Set(userID string, client *Client) {
	u.lock.Lock()
	defer u.lock.Unlock()
	result, ok := u.data[userID]
	if ok {
		result.Clients[client.ctx.GetRemoteAddr()] = client
	} else {
		result = &UserPlatform{
			Clients: map[string]*Client{
				client.ctx.GetRemoteAddr(): client,
			},
		}
	}
	u.push(client.UserID, result, nil)
}

func (u *UserMap2) Delete(userID string, connRemoteAddr string) (isDeleteUser bool) {
	u.lock.Lock()
	defer u.lock.Unlock()
	result, ok := u.data[userID]
	if !ok {
		return false
	}
	client, ok := result.Clients[connRemoteAddr]
	if !ok {
		return false
	}
	delete(result.Clients, connRemoteAddr)
	defer u.push(userID, result, []int32{int32(client.PlatformID)})
	if len(result.Clients) > 0 {
		return false
	}
	delete(u.data, userID)
	return true
}

func (u *UserMap2) DeleteClients(userID string, clients []*Client) (isDeleteUser bool) {
	if len(clients) == 0 {
		return false
	}
	u.lock.Lock()
	defer u.lock.Unlock()
	result, ok := u.data[userID]
	if !ok {
		return false
	}
	offline := make([]int32, 0, len(clients))
	for _, client := range clients {
		offline = append(offline, int32(client.PlatformID))
		delete(result.Clients, client.ctx.GetRemoteAddr())
	}
	defer u.push(userID, result, offline)
	if len(result.Clients) > 0 {
		return false
	}
	delete(u.data, userID)
	return true
}

func (u *UserMap2) GetAllUserStatus(deadline, nowtime time.Time) []UserState {
	u.lock.RLock()
	defer u.lock.RUnlock()
	if len(u.data) == 0 {
		return nil
	}
	result := make([]UserState, 0, len(u.data))
	for userID, p := range u.data {
		if len(result) == cap(result) {
			break
		}
		if p.Time.Before(deadline) {
			continue
		}
		p.Time = nowtime
		online := make([]int32, 0, len(p.Clients))
		for _, client := range p.Clients {
			online = append(online, int32(client.PlatformID))
		}
		result = append(result, UserState{UserID: userID, Online: online})
	}
	return result
}

func (u *UserMap2) UserState() <-chan UserState {
	return u.ch
}

type UserState struct {
	UserID  string
	Online  []int32
	Offline []int32
}
