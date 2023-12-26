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

package convert

import (
	"context"
	"fmt"
	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/utils"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
)

func FriendPb2DB(friend *sdkws.FriendInfo) *relation.FriendModel {
	dbFriend := &relation.FriendModel{}
	utils.CopyStructFields(dbFriend, friend)
	dbFriend.FriendUserID = friend.FriendUser.UserID
	dbFriend.CreateTime = utils.UnixSecondToTime(friend.CreateTime)
	return dbFriend
}

func FriendDB2Pb(ctx context.Context, friendDB *relation.FriendModel,
	getUsers func(ctx context.Context, userIDs []string) (map[string]*sdkws.UserInfo, error),
) (*sdkws.FriendInfo, error) {
	users, err := getUsers(ctx, []string{friendDB.FriendUserID})
	if err != nil {
		return nil, err
	}
	user, ok := users[friendDB.FriendUserID]
	if !ok {
		return nil, fmt.Errorf("user not found: %s", friendDB.FriendUserID)
	}

	return &sdkws.FriendInfo{
		FriendUser: user,
		CreateTime: friendDB.CreateTime.Unix(),
	}, nil
}

func FriendsDB2Pb(
	ctx context.Context,
	friendsDB []*relation.FriendModel,
	getUsers func(ctx context.Context, userIDs []string) (map[string]*sdkws.UserInfo, error),
) (friendsPb []*sdkws.FriendInfo, err error) {
	if len(friendsDB) == 0 {
		return nil, nil
	}
	var userID []string
	for _, friendDB := range friendsDB {
		userID = append(userID, friendDB.FriendUserID)
	}

	users, err := getUsers(ctx, userID)
	if err != nil {
		return nil, err
	}
	for _, friend := range friendsDB {
		friendPb := &sdkws.FriendInfo{FriendUser: &sdkws.UserInfo{}}
		utils.CopyStructFields(friendPb, friend)
		friendPb.FriendUser.UserID = users[friend.FriendUserID].UserID
		friendPb.FriendUser.Nickname = users[friend.FriendUserID].Nickname
		friendPb.FriendUser.FaceURL = users[friend.FriendUserID].FaceURL
		friendPb.FriendUser.Ex = users[friend.FriendUserID].Ex
		friendPb.CreateTime = friend.CreateTime.Unix()
		friendPb.IsPinned = friend.IsPinned
		friendsPb = append(friendsPb, friendPb)
	}
	return friendsPb, nil
}

func FriendRequestDB2Pb(
	ctx context.Context,
	friendRequests []*relation.FriendRequestModel,
	getUsers func(ctx context.Context, userIDs []string) (map[string]*sdkws.UserInfo, error),
) ([]*sdkws.FriendRequest, error) {
	if len(friendRequests) == 0 {
		return nil, nil
	}
	userIDMap := make(map[string]struct{})
	for _, friendRequest := range friendRequests {
		userIDMap[friendRequest.ToUserID] = struct{}{}
		userIDMap[friendRequest.FromUserID] = struct{}{}
	}
	users, err := getUsers(ctx, utils.Keys(userIDMap))
	if err != nil {
		return nil, err
	}
	res := make([]*sdkws.FriendRequest, 0, len(friendRequests))
	for _, friendRequest := range friendRequests {
		toUser := users[friendRequest.ToUserID]
		fromUser := users[friendRequest.FromUserID]
		res = append(res, &sdkws.FriendRequest{
			FromUserID:    friendRequest.FromUserID,
			FromNickname:  fromUser.Nickname,
			FromFaceURL:   fromUser.FaceURL,
			ToUserID:      friendRequest.ToUserID,
			ToNickname:    toUser.Nickname,
			ToFaceURL:     toUser.FaceURL,
			HandleResult:  friendRequest.HandleResult,
			ReqMsg:        friendRequest.ReqMsg,
			CreateTime:    friendRequest.CreateTime.UnixMilli(),
			HandlerUserID: friendRequest.HandlerUserID,
			HandleMsg:     friendRequest.HandleMsg,
			HandleTime:    friendRequest.HandleTime.UnixMilli(),
			Ex:            friendRequest.Ex,
		})
	}
	return res, nil
}

// FriendPb2DBMap converts a FriendInfo protobuf object to a map suitable for database operations.
// It only includes non-zero or non-empty fields in the map.
func FriendPb2DBMap(friend *sdkws.FriendInfo) map[string]any {
	if friend == nil {
		return nil
	}

	val := make(map[string]any)

	// Assuming FriendInfo has similar fields to those in FriendModel.
	// Add or remove fields based on your actual FriendInfo and FriendModel structures.
	if friend.FriendUser != nil {
		if friend.FriendUser.UserID != "" {
			val["friend_user_id"] = friend.FriendUser.UserID
		}
		if friend.FriendUser.Nickname != "" {
			val["nickname"] = friend.FriendUser.Nickname
		}
		if friend.FriendUser.FaceURL != "" {
			val["face_url"] = friend.FriendUser.FaceURL
		}
		if friend.FriendUser.Ex != "" {
			val["ex"] = friend.FriendUser.Ex
		}
	}
	if friend.CreateTime != 0 {
		val["create_time"] = friend.CreateTime // You might need to convert this to a proper time format.
	}

	// Include other fields from FriendInfo as needed, similar to the above pattern.

	return val
}
