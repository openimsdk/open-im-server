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

package rpcclient

import (
	"context"

	"github.com/openimsdk/protocol/relation"
	sdkws "github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/system/program"
	"google.golang.org/grpc"
)

type Friend struct {
	conn   grpc.ClientConnInterface
	Client relation.FriendClient
	discov discovery.SvcDiscoveryRegistry
}

func NewFriend(discov discovery.SvcDiscoveryRegistry, rpcRegisterName string) *Friend {
	conn, err := discov.GetConn(context.Background(), rpcRegisterName)
	if err != nil {
		program.ExitWithError(err)
	}
	client := relation.NewFriendClient(conn)
	return &Friend{discov: discov, conn: conn, Client: client}
}

type FriendRpcClient Friend

func NewFriendRpcClient(discov discovery.SvcDiscoveryRegistry, rpcRegisterName string) FriendRpcClient {
	return FriendRpcClient(*NewFriend(discov, rpcRegisterName))
}

func (f *FriendRpcClient) GetFriendsInfo(
	ctx context.Context,
	ownerUserID, relationUserID string,
) (resp *sdkws.FriendInfo, err error) {
	r, err := f.Client.GetDesignatedFriends(
		ctx,
		&relation.GetDesignatedFriendsReq{OwnerUserID: ownerUserID, FriendUserIDs: []string{relationUserID}},
	)
	if err != nil {
		return nil, err
	}
	resp = r.FriendsInfo[0]
	return
}

// possibleFriendUserID Is PossibleFriendUserId's relations.
func (f *FriendRpcClient) IsFriend(ctx context.Context, possibleFriendUserID, userID string) (bool, error) {
	resp, err := f.Client.IsFriend(ctx, &relation.IsFriendReq{UserID1: userID, UserID2: possibleFriendUserID})
	if err != nil {
		return false, err
	}
	return resp.InUser1Friends, nil
}

func (f *FriendRpcClient) GetFriendIDs(ctx context.Context, ownerUserID string) (relationIDs []string, err error) {
	req := relation.GetFriendIDsReq{UserID: ownerUserID}
	resp, err := f.Client.GetFriendIDs(ctx, &req)
	if err != nil {
		return nil, err
	}
	return resp.FriendIDs, nil
}

func (f *FriendRpcClient) IsBlack(ctx context.Context, possibleBlackUserID, userID string) (bool, error) {
	r, err := f.Client.IsBlack(ctx, &relation.IsBlackReq{UserID1: possibleBlackUserID, UserID2: userID})
	if err != nil {
		return false, err
	}
	return r.InUser2Blacks, nil
}
