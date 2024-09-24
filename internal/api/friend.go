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
	"github.com/gin-gonic/gin"

	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	"github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/tools/a2r"
)

type FriendApi rpcclient.Friend

func NewFriendApi(client rpcclient.Friend) FriendApi {
	return FriendApi(client)
}

func (o *FriendApi) ApplyToAddFriend(c *gin.Context) {
	a2r.Call(relation.FriendClient.ApplyToAddFriend, o.Client, c)
}

func (o *FriendApi) RespondFriendApply(c *gin.Context) {
	a2r.Call(relation.FriendClient.RespondFriendApply, o.Client, c)
}

func (o *FriendApi) DeleteFriend(c *gin.Context) {
	a2r.Call(relation.FriendClient.DeleteFriend, o.Client, c)
}

func (o *FriendApi) GetFriendApplyList(c *gin.Context) {
	a2r.Call(relation.FriendClient.GetPaginationFriendsApplyTo, o.Client, c)
}

func (o *FriendApi) GetDesignatedFriendsApply(c *gin.Context) {
	a2r.Call(relation.FriendClient.GetDesignatedFriendsApply, o.Client, c)
}

func (o *FriendApi) GetSelfApplyList(c *gin.Context) {
	a2r.Call(relation.FriendClient.GetPaginationFriendsApplyFrom, o.Client, c)
}

func (o *FriendApi) GetFriendList(c *gin.Context) {
	a2r.Call(relation.FriendClient.GetPaginationFriends, o.Client, c)
}

func (o *FriendApi) GetDesignatedFriends(c *gin.Context) {
	a2r.Call(relation.FriendClient.GetDesignatedFriends, o.Client, c)
}

func (o *FriendApi) SetFriendRemark(c *gin.Context) {
	a2r.Call(relation.FriendClient.SetFriendRemark, o.Client, c)
}

func (o *FriendApi) AddBlack(c *gin.Context) {
	a2r.Call(relation.FriendClient.AddBlack, o.Client, c)
}

func (o *FriendApi) GetPaginationBlacks(c *gin.Context) {
	a2r.Call(relation.FriendClient.GetPaginationBlacks, o.Client, c)
}

func (o *FriendApi) GetSpecifiedBlacks(c *gin.Context) {
	a2r.Call(relation.FriendClient.GetSpecifiedBlacks, o.Client, c)
}

func (o *FriendApi) RemoveBlack(c *gin.Context) {
	a2r.Call(relation.FriendClient.RemoveBlack, o.Client, c)
}

func (o *FriendApi) ImportFriends(c *gin.Context) {
	a2r.Call(relation.FriendClient.ImportFriends, o.Client, c)
}

func (o *FriendApi) IsFriend(c *gin.Context) {
	a2r.Call(relation.FriendClient.IsFriend, o.Client, c)
}

func (o *FriendApi) GetFriendIDs(c *gin.Context) {
	a2r.Call(relation.FriendClient.GetFriendIDs, o.Client, c)
}

func (o *FriendApi) GetSpecifiedFriendsInfo(c *gin.Context) {
	a2r.Call(relation.FriendClient.GetSpecifiedFriendsInfo, o.Client, c)
}

func (o *FriendApi) UpdateFriends(c *gin.Context) {
	a2r.Call(relation.FriendClient.UpdateFriends, o.Client, c)
}

func (o *FriendApi) GetIncrementalFriends(c *gin.Context) {
	a2r.Call(relation.FriendClient.GetIncrementalFriends, o.Client, c)
}

// GetIncrementalBlacks is temporarily unused.
// Deprecated: This function is currently unused and may be removed in future versions.
func (o *FriendApi) GetIncrementalBlacks(c *gin.Context) {
	a2r.Call(relation.FriendClient.GetIncrementalBlacks, o.Client, c)
}

func (o *FriendApi) GetFullFriendUserIDs(c *gin.Context) {
	a2r.Call(relation.FriendClient.GetFullFriendUserIDs, o.Client, c)
}
