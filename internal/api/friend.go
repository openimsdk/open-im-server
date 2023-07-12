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
	"github.com/OpenIMSDK/Open-IM-Server/pkg/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/friend"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
	"github.com/gin-gonic/gin"
)

// type FriendAPI
type FriendAPI rpcclient.Friend

// NewFriendAPI creates a new friend
func NewFriendAPI(discov discoveryregistry.SvcDiscoveryRegistry) FriendAPI {
	return FriendAPI(*rpcclient.NewFriend(discov))
}

// apply to add a friend
func (o *FriendAPI) ApplyToAddFriend(c *gin.Context) {
	a2r.Call(friend.FriendClient.ApplyToAddFriend, o.Client, c)
}

// response Friend's apply
func (o *FriendAPI) RespondFriendApply(c *gin.Context) {
	a2r.Call(friend.FriendClient.RespondFriendApply, o.Client, c)
}

// delete a friend
func (o *FriendAPI) DeleteFriend(c *gin.Context) {
	a2r.Call(friend.FriendClient.DeleteFriend, o.Client, c)
}

// get friend list
func (o *FriendAPI) GetFriendApplyList(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetPaginationFriendsApplyTo, o.Client, c)
}

// get friend self list for apply
func (o *FriendAPI) GetSelfApplyList(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetPaginationFriendsApplyFrom, o.Client, c)
}

// get friend list
func (o *FriendAPI) GetFriendList(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetPaginationFriends, o.Client, c)
}

// set friend remark sign
func (o *FriendAPI) SetFriendRemark(c *gin.Context) {
	a2r.Call(friend.FriendClient.SetFriendRemark, o.Client, c)
}

// add friend to blacklist
func (o *FriendAPI) AddBlack(c *gin.Context) {
	a2r.Call(friend.FriendClient.AddBlack, o.Client, c)
}

// get balck list with pagenation
func (o *FriendAPI) GetPaginationBlacks(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetPaginationBlacks, o.Client, c)
}

// remove friend from black
func (o *FriendAPI) RemoveBlack(c *gin.Context) {
	a2r.Call(friend.FriendClient.RemoveBlack, o.Client, c)
}

// import friends
func (o *FriendAPI) ImportFriends(c *gin.Context) {
	a2r.Call(friend.FriendClient.ImportFriends, o.Client, c)
}

// judege friend is or not friend
func (o *FriendAPI) IsFriend(c *gin.Context) {
	a2r.Call(friend.FriendClient.IsFriend, o.Client, c)
}
