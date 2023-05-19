package convert

import (
	"time"

	relationTb "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

func UsersDB2Pb(users []*relationTb.UserModel) (result []*sdkws.UserInfo) {
	for _, user := range users {
		var userPb sdkws.UserInfo
		userPb.UserID = user.UserID
		userPb.Nickname = user.Nickname
		userPb.FaceURL = user.FaceURL
		userPb.Ex = user.Ex
		userPb.CreateTime = user.CreateTime.UnixMilli()
		userPb.AppMangerLevel = user.AppMangerLevel
		userPb.GlobalRecvMsgOpt = user.GlobalRecvMsgOpt
		result = append(result, &userPb)
	}
	return result
}

func UserPb2DB(user *sdkws.UserInfo) *relationTb.UserModel {
	var userDB relationTb.UserModel
	userDB.UserID = user.UserID
	userDB.Nickname = user.Nickname
	userDB.FaceURL = user.FaceURL
	userDB.Ex = user.Ex
	userDB.CreateTime = time.UnixMilli(user.CreateTime)
	userDB.AppMangerLevel = user.AppMangerLevel
	userDB.GlobalRecvMsgOpt = user.GlobalRecvMsgOpt
	return &userDB
}
