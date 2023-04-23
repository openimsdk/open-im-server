package convert

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
	utils "github.com/OpenIMSDK/open_utils"
)

func FriendPb2DB(friend *sdkws.FriendInfo) *relation.FriendModel {
	dbFriend := &relation.FriendModel{}
	utils.CopyStructFields(dbFriend, friend)
	dbFriend.FriendUserID = friend.FriendUser.UserID
	dbFriend.CreateTime = utils.UnixSecondToTime(friend.CreateTime)
	return dbFriend
}

func FriendDB2Pb(ctx context.Context, friendDB *relation.FriendModel, getUser func(ctx context.Context, userIDs []string) ([]rpcclient.CommonUser, error)) (*sdkws.FriendInfo, error) {
	pbfriend := &sdkws.FriendInfo{FriendUser: &sdkws.UserInfo{}}
	utils.CopyStructFields(pbfriend, friendDB)
	users, err := getUser(ctx, []string{friendDB.FriendUserID})
	if err != nil {
		return nil, err
	}
	pbfriend.FriendUser.UserID = users[0].GetUserID()
	pbfriend.FriendUser.Nickname = users[0].GetNickname()
	pbfriend.FriendUser.FaceURL = users[0].GetFaceURL()
	pbfriend.FriendUser.Ex = users[0].GetEx()
	pbfriend.CreateTime = friendDB.CreateTime.Unix()
	return pbfriend, nil
}
