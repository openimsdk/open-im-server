package convert

import (
	relationtb "github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/utils/datautil"
	"time"

	"github.com/openimsdk/protocol/sdkws"
)

func UserDB2Pb(user *relationtb.User) *sdkws.UserInfo {
	return &sdkws.UserInfo{
		UserID:           user.UserID,
		Nickname:         user.Nickname,
		FaceURL:          user.FaceURL,
		Ex:               user.Ex,
		CreateTime:       user.CreateTime.UnixMilli(),
		AppMangerLevel:   user.AppMangerLevel,
		GlobalRecvMsgOpt: user.GlobalRecvMsgOpt,
	}
}

func UsersDB2Pb(users []*relationtb.User) []*sdkws.UserInfo {
	return datautil.Slice(users, UserDB2Pb)
}

func UserPb2DB(user *sdkws.UserInfo) *relationtb.User {
	return &relationtb.User{
		UserID:           user.UserID,
		Nickname:         user.Nickname,
		FaceURL:          user.FaceURL,
		Ex:               user.Ex,
		CreateTime:       time.UnixMilli(user.CreateTime),
		AppMangerLevel:   user.AppMangerLevel,
		GlobalRecvMsgOpt: user.GlobalRecvMsgOpt,
	}
}

func UserPb2DBMap(user *sdkws.UserInfo) map[string]any {
	if user == nil {
		return nil
	}
	val := make(map[string]any)
	fields := map[string]any{
		"nickname":            user.Nickname,
		"face_url":            user.FaceURL,
		"ex":                  user.Ex,
		"app_manager_level":   user.AppMangerLevel,
		"global_recv_msg_opt": user.GlobalRecvMsgOpt,
	}
	for key, value := range fields {
		if v, ok := value.(string); ok && v != "" {
			val[key] = v
		} else if v, ok := value.(int32); ok && v != 0 {
			val[key] = v
		}
	}
	return val
}
func UserPb2DBMapEx(user *sdkws.UserInfoWithEx) map[string]any {
	if user == nil {
		return nil
	}
	val := make(map[string]any)

	// Map fields from UserInfoWithEx to val
	if user.Nickname != nil {
		val["nickname"] = user.Nickname.Value
	}
	if user.FaceURL != nil {
		val["face_url"] = user.FaceURL.Value
	}
	if user.Ex != nil {
		val["ex"] = user.Ex.Value
	}
	if user.GlobalRecvMsgOpt != nil {
		val["global_recv_msg_opt"] = user.GlobalRecvMsgOpt.Value
	}

	return val
}
