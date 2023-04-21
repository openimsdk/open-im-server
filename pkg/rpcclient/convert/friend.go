package convert

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	sdk "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	utils "github.com/OpenIMSDK/open_utils"
)

func FriendPb2DB(friend *sdkws.FriendInfo) (*relation.FriendModel, error) {
	dbFriend := &relation.FriendModel{}
	utils.CopyStructFields(dbFriend, friend)
	dbFriend.FriendUserID = friend.FriendUser.UserID
	dbFriend.CreateTime = utils.UnixSecondToTime(friend.CreateTime)
	return dbFriend, nil
}

func FriendDB2Pb(ctx context.Context, friendDB *relation.FriendModel, fn func(ctx context.Context, userID string) (*sdkws.UserInfo, error)) (*sdk.FriendInfo, error) {
	pbfriend := &sdk.FriendInfo{FriendUser: &sdk.UserInfo{}}
	utils.CopyStructFields(pbfriend, friendDB)
	user, err := fn(ctx, friendDB.FriendUserID)
	if err != nil {
		return nil, err
	}
	utils.CopyStructFields(pbfriend.FriendUser, user)
	pbfriend.CreateTime = friendDB.CreateTime.Unix()
	pbfriend.FriendUser.CreateTime = friendDB.CreateTime.Unix()
	return pbfriend, nil
}
