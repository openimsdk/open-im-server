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
	"github.com/openimsdk/protocol/friend"
	"github.com/openimsdk/tools/a2r"
)

type FriendApi rpcclient.Friend

func NewFriendApi(client rpcclient.Friend) FriendApi {
	return FriendApi(client)
}

func (o *FriendApi) ApplyToAddFriend(c *gin.Context) {
	a2r.Call(friend.FriendClient.ApplyToAddFriend, o.Client, c)
}

func (o *FriendApi) RespondFriendApply(c *gin.Context) {
	a2r.Call(friend.FriendClient.RespondFriendApply, o.Client, c)
}

func (o *FriendApi) DeleteFriend(c *gin.Context) {
	a2r.Call(friend.FriendClient.DeleteFriend, o.Client, c)
}

func (o *FriendApi) GetFriendApplyList(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetPaginationFriendsApplyTo, o.Client, c)
}

func (o *FriendApi) GetDesignatedFriendsApply(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetDesignatedFriendsApply, o.Client, c)
}

func (o *FriendApi) GetSelfApplyList(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetPaginationFriendsApplyFrom, o.Client, c)
}

func (o *FriendApi) GetFriendList(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetPaginationFriends, o.Client, c)
}

func (o *FriendApi) GetDesignatedFriends(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetDesignatedFriends, o.Client, c)
}

func (o *FriendApi) SetFriendRemark(c *gin.Context) {
	a2r.Call(friend.FriendClient.SetFriendRemark, o.Client, c)
}

func (o *FriendApi) AddBlack(c *gin.Context) {
	a2r.Call(friend.FriendClient.AddBlack, o.Client, c)
}

func (o *FriendApi) GetPaginationBlacks(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetPaginationBlacks, o.Client, c)
}

func (o *FriendApi) RemoveBlack(c *gin.Context) {
	a2r.Call(friend.FriendClient.RemoveBlack, o.Client, c)
}

func (o *FriendApi) ImportFriends(c *gin.Context) {
	a2r.Call(friend.FriendClient.ImportFriends, o.Client, c)
}

func (o *FriendApi) IsFriend(c *gin.Context) {
	a2r.Call(friend.FriendClient.IsFriend, o.Client, c)
}

func (o *FriendApi) GetFriendIDs(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetFriendIDs, o.Client, c)
}

func (o *FriendApi) GetSpecifiedFriendsInfo(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetSpecifiedFriendsInfo, o.Client, c)
}
func (o *FriendApi) UpdateFriends(c *gin.Context) {
	a2r.Call(friend.FriendClient.UpdateFriends, o.Client, c)
}
