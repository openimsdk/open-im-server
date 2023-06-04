package zookeeper

import (
	"sync"

	"google.golang.org/grpc"
)

type RoundRobin struct {
	index int
	lock  sync.Mutex
}

func (r *RoundRobin) getConnBalance(conns []*grpc.ClientConn) (conn *grpc.ClientConn) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.index < len(conns)-1 {
		r.index++
	} else {
		r.index = 0
	}
	return conns[r.index]
}
