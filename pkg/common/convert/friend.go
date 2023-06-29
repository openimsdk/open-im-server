package convert

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
)

func FriendPb2DB(friend *sdkws.FriendInfo) *relation.FriendModel {
	dbFriend := &relation.FriendModel{}
	utils.CopyStructFields(dbFriend, friend)
	dbFriend.FriendUserID = friend.FriendUser.UserID
	dbFriend.CreateTime = utils.UnixSecondToTime(friend.CreateTime)
	return dbFriend
}

func FriendDB2Pb(ctx context.Context, friendDB *relation.FriendModel, getUsers func(ctx context.Context, userIDs []string) (map[string]*sdkws.UserInfo, error)) (*sdkws.FriendInfo, error) {
	pbfriend := &sdkws.FriendInfo{FriendUser: &sdkws.UserInfo{}}
	utils.CopyStructFields(pbfriend, friendDB)
	users, err := getUsers(ctx, []string{friendDB.FriendUserID})
	if err != nil {
		return nil, err
	}
	pbfriend.FriendUser.UserID = users[friendDB.FriendUserID].UserID
	pbfriend.FriendUser.Nickname = users[friendDB.FriendUserID].Nickname
	pbfriend.FriendUser.FaceURL = users[friendDB.FriendUserID].FaceURL
	pbfriend.FriendUser.Ex = users[friendDB.FriendUserID].Ex
	pbfriend.CreateTime = friendDB.CreateTime.Unix()
	return pbfriend, nil
}

func FriendsDB2Pb(ctx context.Context, friendsDB []*relation.FriendModel, getUsers func(ctx context.Context, userIDs []string) (map[string]*sdkws.UserInfo, error)) (friendsPb []*sdkws.FriendInfo, err error) {
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
		friendsPb = append(friendsPb, friendPb)
	}
	return friendsPb, nil
}

func FriendRequestDB2Pb(ctx context.Context, friendRequests []*relation.FriendRequestModel, getUsers func(ctx context.Context, userIDs []string) (map[string]*sdkws.UserInfo, error)) ([]*sdkws.FriendRequest, error) {
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
