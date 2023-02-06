package group

import (
	"Open_IM/pkg/common/tracelog"
	sdk_ws "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"math/big"
	"strconv"
	"time"
)

//
//func getDBGroupMember(ctx context.Context, groupID, userID string) (dbGroupMember *relation.GroupMember, err error) {
//	dbGroupMember = &relation.GroupMember{}
//
//	member := relation.GroupMember{}
//	member.GroupID = groupID
//	member.UserID = userID
//	member.RoleLevel = constant.GroupOrdinaryUsers
//	member.OperatorUserID = utils.OpUserID(ctx)
//
//	member.FaceURL = user.FaceURL
//	member.Nickname = user.Nickname
//	member.JoinSource = request.JoinSource
//	member.InviterUserID = request.InviterUserID
//	member.MuteEndTime = time.Unix(int64(time.Now().Second()), 0)
//
//	return dbGroupMember, nil
//}

func GetPublicUserInfoOne(ctx context.Context, userID string) (*sdk_ws.PublicUserInfo, error) {
	return nil, errors.New("todo")
}

func GetUsersInfo(ctx context.Context, userIDs []string) ([]*sdk_ws.UserInfo, error) {
	return nil, errors.New("todo")
}

func GetUserInfoMap(ctx context.Context, userIDs []string) (map[string]*sdk_ws.UserInfo, error) {
	users, err := GetUsersInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	return utils.SliceToMap(users, func(e *sdk_ws.UserInfo) string {
		return e.UserID
	}), nil
}

func GetPublicUserInfo(ctx context.Context, userIDs []string) ([]*sdk_ws.PublicUserInfo, error) {
	return nil, errors.New("todo")
}

func GetPublicUserInfoMap(ctx context.Context, userIDs []string) (map[string]*sdk_ws.PublicUserInfo, error) {
	users, err := GetPublicUserInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	return utils.SliceToMap(users, func(e *sdk_ws.PublicUserInfo) string {
		return e.UserID
	}), nil
}

func genGroupID(ctx context.Context, groupID string) string {
	if groupID != "" {
		return groupID
	}
	groupID = utils.Md5(tracelog.GetOperationID(ctx) + strconv.FormatInt(time.Now().UnixNano(), 10))
	bi := big.NewInt(0)
	bi.SetString(groupID[0:8], 16)
	groupID = bi.String()
	return groupID
}
