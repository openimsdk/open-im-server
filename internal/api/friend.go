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

	"github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/tools/a2r"
)

type FriendApi struct{}

func NewFriendApi() FriendApi {
	return FriendApi{}
}

func (o *FriendApi) ApplyToAddFriend(c *gin.Context) {
	a2r.CallV2(relation.ApplyToAddFriendCaller.Invoke, c)
}

func (o *FriendApi) RespondFriendApply(c *gin.Context) {
	a2r.CallV2(relation.RespondFriendApplyCaller.Invoke, c)
}

func (o *FriendApi) DeleteFriend(c *gin.Context) {
	a2r.CallV2(relation.DeleteFriendCaller.Invoke, c)
}

func (o *FriendApi) GetFriendApplyList(c *gin.Context) {
	a2r.CallV2(relation.GetPaginationFriendsApplyToCaller.Invoke, c)
}

func (o *FriendApi) GetDesignatedFriendsApply(c *gin.Context) {
	a2r.CallV2(relation.GetDesignatedFriendsApplyCaller.Invoke, c)
}

func (o *FriendApi) GetSelfApplyList(c *gin.Context) {
	a2r.CallV2(relation.GetPaginationFriendsApplyFromCaller.Invoke, c)
}

func (o *FriendApi) GetFriendList(c *gin.Context) {
	a2r.CallV2(relation.GetPaginationFriendsCaller.Invoke, c)
}

func (o *FriendApi) GetDesignatedFriends(c *gin.Context) {
	a2r.CallV2(relation.GetDesignatedFriendsCaller.Invoke, c)
}

func (o *FriendApi) SetFriendRemark(c *gin.Context) {
	a2r.CallV2(relation.SetFriendRemarkCaller.Invoke, c)
}

func (o *FriendApi) AddBlack(c *gin.Context) {
	a2r.CallV2(relation.AddBlackCaller.Invoke, c)
}

func (o *FriendApi) GetPaginationBlacks(c *gin.Context) {
	a2r.CallV2(relation.GetPaginationBlacksCaller.Invoke, c)
}

func (o *FriendApi) GetSpecifiedBlacks(c *gin.Context) {
	a2r.CallV2(relation.GetSpecifiedBlacksCaller.Invoke, c)
}

func (o *FriendApi) RemoveBlack(c *gin.Context) {
	a2r.CallV2(relation.RemoveBlackCaller.Invoke, c)
}

func (o *FriendApi) ImportFriends(c *gin.Context) {
	a2r.CallV2(relation.ImportFriendsCaller.Invoke, c)
}

func (o *FriendApi) IsFriend(c *gin.Context) {
	a2r.CallV2(relation.IsFriendCaller.Invoke, c)
}

func (o *FriendApi) GetFriendIDs(c *gin.Context) {
	a2r.CallV2(relation.GetFriendIDsCaller.Invoke, c)
}

func (o *FriendApi) GetSpecifiedFriendsInfo(c *gin.Context) {
	a2r.CallV2(relation.GetSpecifiedFriendsInfoCaller.Invoke, c)
}

func (o *FriendApi) UpdateFriends(c *gin.Context) {
	a2r.CallV2(relation.UpdateFriendsCaller.Invoke, c)
}

func (o *FriendApi) GetIncrementalFriends(c *gin.Context) {
	a2r.CallV2(relation.GetIncrementalFriendsCaller.Invoke, c)
}

// GetIncrementalBlacks is temporarily unused.
// Deprecated: This function is currently unused and may be removed in future versions.
func (o *FriendApi) GetIncrementalBlacks(c *gin.Context) {
	a2r.CallV2(relation.GetIncrementalBlacksCaller.Invoke, c)
}

func (o *FriendApi) GetFullFriendUserIDs(c *gin.Context) {
	a2r.CallV2(relation.GetFullFriendUserIDsCaller.Invoke, c)
}
