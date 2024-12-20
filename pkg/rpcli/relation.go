package rpcli

import "github.com/openimsdk/protocol/relation"

func NewRelationClient(cli relation.FriendClient) *RelationClient {
	return &RelationClient{cli}
}

type RelationClient struct {
	relation.FriendClient
}
