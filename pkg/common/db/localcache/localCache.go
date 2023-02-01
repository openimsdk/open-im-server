package localcache

import "google.golang.org/grpc"

type GroupLocalCache struct {
	cache map[string]GroupMemberIDsHash
	rpc   *grpc.ClientConn
}

type GroupMemberIDsHash struct {
	MemberListHash uint64
	UserIDs        []string
}

func NewGroupMemberIDsLocalCache(rpc *grpc.ClientConn) GroupLocalCache {
	return GroupLocalCache{
		cache: make(map[string]GroupMemberIDsHash, 0),
		rpc:   rpc,
	}
}

func (g *GroupMemberIDsHash) GetGroupMemberIDs(groupID string) []string {
	return []string{}
}
