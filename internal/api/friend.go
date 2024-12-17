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
	a2r.CallV2(c, relation.ApplyToAddFriendCaller.Invoke)
}

func (o *FriendApi) RespondFriendApply(c *gin.Context) {
	a2r.CallV2(c, relation.RespondFriendApplyCaller.Invoke)
}

func (o *FriendApi) DeleteFriend(c *gin.Context) {
	a2r.CallV2(c, relation.DeleteFriendCaller.Invoke)
}

func (o *FriendApi) GetFriendApplyList(c *gin.Context) {
	a2r.CallV2(c, relation.GetPaginationFriendsApplyToCaller.Invoke)
}

func (o *FriendApi) GetDesignatedFriendsApply(c *gin.Context) {
	a2r.CallV2(c, relation.GetDesignatedFriendsApplyCaller.Invoke)
}

func (o *FriendApi) GetSelfApplyList(c *gin.Context) {
	a2r.CallV2(c, relation.GetPaginationFriendsApplyFromCaller.Invoke)
}

func (o *FriendApi) GetFriendList(c *gin.Context) {
	a2r.CallV2(c, relation.GetPaginationFriendsCaller.Invoke)
}

func (o *FriendApi) GetDesignatedFriends(c *gin.Context) {
	a2r.CallV2(c, relation.GetDesignatedFriendsCaller.Invoke)
}

func (o *FriendApi) SetFriendRemark(c *gin.Context) {
	a2r.CallV2(c, relation.SetFriendRemarkCaller.Invoke)
}

func (o *FriendApi) AddBlack(c *gin.Context) {
	a2r.CallV2(c, relation.AddBlackCaller.Invoke)
}

func (o *FriendApi) GetPaginationBlacks(c *gin.Context) {
	a2r.CallV2(c, relation.GetPaginationBlacksCaller.Invoke)
}

func (o *FriendApi) GetSpecifiedBlacks(c *gin.Context) {
	a2r.CallV2(c, relation.GetSpecifiedBlacksCaller.Invoke)
}

func (o *FriendApi) RemoveBlack(c *gin.Context) {
	a2r.CallV2(c, relation.RemoveBlackCaller.Invoke)
}

func (o *FriendApi) ImportFriends(c *gin.Context) {
	a2r.CallV2(c, relation.ImportFriendsCaller.Invoke)
}

func (o *FriendApi) IsFriend(c *gin.Context) {
	a2r.CallV2(c, relation.IsFriendCaller.Invoke)
}

func (o *FriendApi) GetFriendIDs(c *gin.Context) {
	a2r.CallV2(c, relation.GetFriendIDsCaller.Invoke)
}

func (o *FriendApi) GetSpecifiedFriendsInfo(c *gin.Context) {
	a2r.CallV2(c, relation.GetSpecifiedFriendsInfoCaller.Invoke)
}

func (o *FriendApi) UpdateFriends(c *gin.Context) {
	a2r.CallV2(c, relation.UpdateFriendsCaller.Invoke)
}

func (o *FriendApi) GetIncrementalFriends(c *gin.Context) {
	a2r.CallV2(c, relation.GetIncrementalFriendsCaller.Invoke)
}

// GetIncrementalBlacks is temporarily unused.
// Deprecated: This function is currently unused and may be removed in future versions.
func (o *FriendApi) GetIncrementalBlacks(c *gin.Context) {
	a2r.CallV2(c, relation.GetIncrementalBlacksCaller.Invoke)
}

func (o *FriendApi) GetFullFriendUserIDs(c *gin.Context) {
	a2r.CallV2(c, relation.GetFullFriendUserIDsCaller.Invoke)
}
