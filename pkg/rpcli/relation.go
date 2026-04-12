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
	if len(friendUserIDs) == 0 {
		return nil, nil
	}
	req := &relation.GetFriendInfoReq{OwnerUserID: ownerUserID, FriendUserIDs: friendUserIDs}
	return extractField(ctx, x.FriendClient.GetFriendInfo, req, (*relation.GetFriendInfoResp).GetFriendInfos)
}

// IsFriend checks whether userID2 is in userID1's friend list.
func (x *RelationClient) IsFriend(ctx context.Context, ownerUserID, friendUserID string) (bool, error) {
	resp, err := x.FriendClient.IsFriend(ctx, &relation.IsFriendReq{
		UserID1: ownerUserID,
		UserID2: friendUserID,
	})
	if err != nil {
		return false, err
	}
	return resp.InUser1Friends, nil
}
