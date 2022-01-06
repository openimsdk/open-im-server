package utils

import (
	db "Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/token_verify"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"math/rand"
	"strconv"
	"time"
)

func OperationIDGenerator() string {
	return strconv.FormatInt(time.Now().UnixNano()+int64(rand.Uint32()), 10)
}

func FriendOpenIMCopyDB(dst *db.Friend, src *open_im_sdk.FriendInfo) {
	utils.CopyStructFields(dst, src)
	dst.FriendUserID = src.FriendUser.UserID
}

func FriendDBCopyOpenIM(dst *open_im_sdk.FriendInfo, src *db.Friend) {
	utils.CopyStructFields(dst, src)
	user, _ := imdb.GetUserByUserID(src.FriendUserID)
	if user != nil {
		utils.CopyStructFields(dst.FriendUser, user)
	}
	dst.CreateTime = uint32(src.CreateTime.Unix())
	dst.FriendUser.CreateTime = user.CreateTime.Unix()
}

//
func FriendRequestOpenIMCopyDB(dst *db.FriendRequest, src *open_im_sdk.FriendRequest) {
	utils.CopyStructFields(dst, src)
}

func FriendRequestDBCopyOpenIM(dst *open_im_sdk.FriendRequest, src *db.FriendRequest) {
	utils.CopyStructFields(dst, src)
	dst.CreateTime = src.CreateTime.Unix()
	dst.HandleTime = src.HandleTime.Unix()
}

func GroupOpenIMCopyDB(dst *db.Group, src *open_im_sdk.GroupInfo) {
	utils.CopyStructFields(dst, src)
}

func GroupDBCopyOpenIM(dst *open_im_sdk.GroupInfo, src *db.Group) {
	utils.CopyStructFields(dst, src)
	user, _ := imdb.GetGroupOwnerInfoByGroupID(src.GroupID)
	if user != nil {
		dst.OwnerUserID = user.UserID
	}
	dst.MemberCount = imdb.GetGroupMemberNumByGroupID(src.GroupID)
	dst.CreateTime = src.CreateTime.Unix()
}

func GroupMemberOpenIMCopyDB(dst *db.GroupMember, src *open_im_sdk.GroupMemberFullInfo) {
	utils.CopyStructFields(dst, src)
}

func GroupMemberDBCopyOpenIM(dst *open_im_sdk.GroupMemberFullInfo, src *db.GroupMember) {
	utils.CopyStructFields(dst, src)
	if token_verify.IsMangerUserID(src.UserID) {
		u, _ := imdb.GetUserByUserID(src.UserID)
		if u != nil {
			utils.CopyStructFields(dst, u)
		}
		dst.AppMangerLevel = 1
	}
	dst.JoinTime = src.JoinTime.Unix()
}

func GroupRequestOpenIMCopyDB(dst *db.GroupRequest, src *open_im_sdk.GroupRequest) {
	utils.CopyStructFields(dst, src)
}

func GroupRequestDBCopyOpenIM(dst *open_im_sdk.GroupRequest, src *db.GroupRequest) {
	utils.CopyStructFields(dst, src)
	dst.ReqTime = src.ReqTime.Unix()
	dst.HandleTime = src.HandledTime.Unix()
}

func UserOpenIMCopyDB(dst *db.User, src *open_im_sdk.UserInfo) {
	utils.CopyStructFields(dst, src)
}

func UserDBCopyOpenIM(dst *open_im_sdk.UserInfo, src *db.User) {
	utils.CopyStructFields(dst, src)
	dst.CreateTime = src.CreateTime.Unix()
}

func BlackOpenIMCopyDB(dst *db.Black, src *open_im_sdk.BlackInfo) {
	utils.CopyStructFields(dst, src)
	dst.BlockUserID = src.BlackUserInfo.UserID
}

func BlackDBCopyOpenIM(dst *open_im_sdk.BlackInfo, src *db.Black) {
	utils.CopyStructFields(dst, src)
	dst.CreateTime = src.CreateTime.Unix()
	user, _ := imdb.GetUserByUserID(src.BlockUserID)
	if user != nil {
		utils.CopyStructFields(dst.BlackUserInfo, user)
	}
}
