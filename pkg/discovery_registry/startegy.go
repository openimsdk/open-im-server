package discoveryRegistry

import "google.golang.org/grpc"

type Robin struct {
	next int
}

func (r *Robin) Robin(slice []*grpc.ClientConn) int {
	index := r.next
	r.next += 1
	if r.next > len(slice)-1 {
		r.next = 0
	}
	return index
}

type Hash struct {
}

func (r *Hash) Hash(slice []*grpc.ClientConn) int {
	index := r.next
	r.next += 1
	if r.next > len(slice)-1 {
		r.next = 0
	}
	return index
}
