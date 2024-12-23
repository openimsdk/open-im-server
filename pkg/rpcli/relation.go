package rpcli

import (
	"context"
	"github.com/openimsdk/protocol/relation"
	"google.golang.org/grpc"
)

func NewRelationClient(cc grpc.ClientConnInterface) *RelationClient {
	return &RelationClient{relation.NewFriendClient(cc)}
}

type RelationClient struct {
	relation.FriendClient
}

func (x *RelationClient) GetFriendsInfo(ctx context.Context, ownerUserID string, friendUserIDs []string) ([]*relation.FriendInfoOnly, error) {
	req := &relation.GetFriendInfoReq{OwnerUserID: ownerUserID, FriendUserIDs: friendUserIDs}
	return extractField(ctx, x.FriendClient.GetFriendInfo, req, (*relation.GetFriendInfoResp).GetFriendInfos)
}
