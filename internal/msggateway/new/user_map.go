package new

import "sync"

type UserMap struct {
	m sync.Map
}

func newUserMap() *UserMap {
	return &UserMap{}
}
func (u *UserMap) GetAll(key string) []*Client {
	allClients, ok := u.m.Load(key)
	if ok {
		return allClients.([]*Client)
	}
	return nil
}
func (u *UserMap) Get(key string, platformID int32) (*Client, bool) {
	allClients, existed := u.m.Load(key)
	if existed {
		for _, client := range allClients.([]*Client) {
			if client.PlatformID == platformID {
				return client, existed
			}
		}
		return nil, false
	}
	return nil, existed
}
func (u *UserMap) Set(key string, v *Client) {
	allClients, existed := u.m.Load(key)
	if existed {
		oldClients := allClients.([]*Client)
		oldClients = append(oldClients, v)
		u.m.Store(key, oldClients)
	} else {
		clients := make([]*Client, 3)
		clients = append(clients, v)
		u.m.Store(key, clients)
	}
}
func (u *UserMap) delete(key string, platformID int32) {
	allClients, existed := u.m.Load(key)
	if existed {
		oldClients := allClients.([]*Client)

		a := make([]*Client, len(oldClients))
		for _, client := range oldClients {
			if client.PlatformID != platformID {
				a = append(a, client)
			}
		}
		if len(a) == 0 {
			u.m.Delete(key)
		} else {
			u.m.Store(key, a)

		}
	}
}
func (u *UserMap) DeleteAll(key string) {
	u.m.Delete(key)
}
