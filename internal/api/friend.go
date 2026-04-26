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

package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/a2r"
	"github.com/openimsdk/tools/errs"
)

type FriendApi struct {
	Client relation.FriendClient
}

func NewFriendApi(client relation.FriendClient) FriendApi {
	return FriendApi{client}
}

func (o *FriendApi) ApplyToAddFriend(c *gin.Context) {
	a2r.Call(c, relation.FriendClient.ApplyToAddFriend, o.Client)
}

func (o *FriendApi) RespondFriendApply(c *gin.Context) {
	a2r.Call(c, relation.FriendClient.RespondFriendApply, o.Client)
}

func (o *FriendApi) DeleteFriend(c *gin.Context) {
	a2r.Call(c, relation.FriendClient.DeleteFriend, o.Client)
}

func (o *FriendApi) GetFriendApplyList(c *gin.Context) {
	a2r.Call(c, relation.FriendClient.GetPaginationFriendsApplyTo, o.Client)
}

func (o *FriendApi) GetDesignatedFriendsApply(c *gin.Context) {
	a2r.Call(c, relation.FriendClient.GetDesignatedFriendsApply, o.Client)
}

func (o *FriendApi) GetSelfApplyList(c *gin.Context) {
	a2r.Call(c, relation.FriendClient.GetPaginationFriendsApplyFrom, o.Client)
}

func (o *FriendApi) GetFriendList(c *gin.Context) {
	a2r.Call(c, relation.FriendClient.GetPaginationFriends, o.Client)
}

func (o *FriendApi) GetDesignatedFriends(c *gin.Context) {
	a2r.Call(c, relation.FriendClient.GetDesignatedFriends, o.Client)
}

func (o *FriendApi) SetFriendRemark(c *gin.Context) {
	a2r.Call(c, relation.FriendClient.SetFriendRemark, o.Client)
}

func (o *FriendApi) AddBlack(c *gin.Context) {
	a2r.Call(c, relation.FriendClient.AddBlack, o.Client)
}

func (o *FriendApi) GetPaginationBlacks(c *gin.Context) {
	a2r.Call(c, relation.FriendClient.GetPaginationBlacks, o.Client)
}

func (o *FriendApi) GetSpecifiedBlacks(c *gin.Context) {
	a2r.Call(c, relation.FriendClient.GetSpecifiedBlacks, o.Client)
}

func (o *FriendApi) RemoveBlack(c *gin.Context) {
	a2r.Call(c, relation.FriendClient.RemoveBlack, o.Client)
}

func (o *FriendApi) ImportFriends(c *gin.Context) {
	a2r.Call(c, relation.FriendClient.ImportFriends, o.Client)
}

func (o *FriendApi) IsFriend(c *gin.Context) {
	a2r.Call(c, relation.FriendClient.IsFriend, o.Client)
}

func (o *FriendApi) GetFriendIDs(c *gin.Context) {
	a2r.Call(c, relation.FriendClient.GetFriendIDs, o.Client)
}

func (o *FriendApi) GetSpecifiedFriendsInfo(c *gin.Context) {
	a2r.Call(c, relation.FriendClient.GetSpecifiedFriendsInfo, o.Client)
}

func (o *FriendApi) UpdateFriends(c *gin.Context) {
	a2r.Call(c, relation.FriendClient.UpdateFriends, o.Client)
}

func (o *FriendApi) GetIncrementalFriends(c *gin.Context) {
	a2r.Call(c, relation.FriendClient.GetIncrementalFriends, o.Client)
}

// GetIncrementalBlacks is temporarily unused.
// Deprecated: This function is currently unused and may be removed in future versions.
func (o *FriendApi) GetIncrementalBlacks(c *gin.Context) {
	a2r.Call(c, relation.FriendClient.GetIncrementalBlacks, o.Client)
}

func (o *FriendApi) GetFullFriendUserIDs(c *gin.Context) {
	a2r.Call(c, relation.FriendClient.GetFullFriendUserIDs, o.Client)
}

func (o *FriendApi) GetSelfUnhandledApplyCount(c *gin.Context) {
	a2r.Call(c, relation.FriendClient.GetSelfUnhandledApplyCount, o.Client)
}

func (o *FriendApi) GetPinnedFriendIDs(c *gin.Context) {
	a2r.Call(c, relation.FriendClient.GetPinnedFriendIDs, o.Client)
}

func (o *FriendApi) AddOnewayFriend(c *gin.Context) {
	// Current generated relation grpc client may not include AddOnewayFriend yet.
	// Keep API route compile-safe and return a clear error instead of breaking build.
	client, ok := any(o.Client).(interface {
		AddOnewayFriend(ctx context.Context, in *relation.ApplyToAddFriendReq, opts ...grpc.CallOption) (*relation.ApplyToAddFriendResp, error)
	})
	if !ok {
		apiresp.GinError(c, errs.New("add_oneway_friend rpc is not generated in relation_grpc.pb.go"))
		return
	}
	var req relation.ApplyToAddFriendReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap())
		return
	}
	resp, err := client.AddOnewayFriend(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}
