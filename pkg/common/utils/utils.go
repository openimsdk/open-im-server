package utils

import (
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

func FriendOpenIMCopyDB(dst *imdb.Friend, src open_im_sdk.FriendInfo) {
	utils.CopyStructFields(dst, src)
	dst.FriendUserID = src.FriendUser.UserID
}

func FriendDBCopyOpenIM(dst *open_im_sdk.FriendInfo, src imdb.Friend) {
	utils.CopyStructFields(dst, src)
	user, _ := imdb.GetUserByUserID(src.FriendUserID)
	if user != nil {
		utils.CopyStructFields(dst.FriendUser, user)
	}
	dst.CreateTime = src.CreateTime.Unix()
	dst.FriendUser.CreateTime = user.CreateTime.Unix()
}

//
func FriendRequestOpenIMCopyDB(dst *imdb.FriendRequest, src open_im_sdk.FriendRequest) {
	utils.CopyStructFields(dst, src)
}

func FriendRequestDBCopyOpenIM(dst *open_im_sdk.FriendRequest, src imdb.FriendRequest) {
	utils.CopyStructFields(dst, src)
	dst.CreateTime = src.CreateTime.Unix()
	dst.HandleTime = src.HandleTime.Unix()
}

func GroupOpenIMCopyDB(dst *imdb.Group, src open_im_sdk.GroupInfo) {
	utils.CopyStructFields(dst, src)
}

func GroupDBCopyOpenIM(dst *open_im_sdk.GroupInfo, src imdb.Group) {
	utils.CopyStructFields(dst, src)
	user, _ := imdb.GetGroupOwnerInfoByGroupID(src.GroupID)
	if user != nil {
		dst.OwnerUserID = user.UserID
	}
	dst.MemberCount = imdb.GetGroupMemberNumByGroupID(src.GroupID)
	dst.CreateTime = src.CreateTime.Unix()
}

func GroupMemberOpenIMCopyDB(dst *imdb.GroupMember, src open_im_sdk.GroupMemberFullInfo) {
	utils.CopyStructFields(dst, src)
}

func GroupMemberDBCopyOpenIM(dst *open_im_sdk.GroupMemberFullInfo, src imdb.GroupMember) {
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

func GroupRequestOpenIMCopyDB(dst *imdb.GroupRequest, src open_im_sdk.GroupRequest) {
	utils.CopyStructFields(dst, src)
}

func GroupRequestDBCopyOpenIM(dst *open_im_sdk.GroupRequest, src imdb.GroupRequest) {
	utils.CopyStructFields(dst, src)
	dst.ReqTime = src.ReqTime.Unix()
	dst.HandleTime = src.HandledTime.Unix()
}

func UserOpenIMCopyDB(dst *imdb.User, src open_im_sdk.UserInfo) {
	utils.CopyStructFields(dst, src)
}

func UserDBCopyOpenIM(dst *open_im_sdk.UserInfo, src imdb.User) {
	utils.CopyStructFields(dst, src)
	dst.CreateTime = src.CreateTime.Unix()
}

func BlackOpenIMCopyDB(dst *imdb.Black, src open_im_sdk.BlackInfo) {
	utils.CopyStructFields(dst, src)
	dst.BlockUserID = src.BlackUserInfo.UserID
}

func BlackDBCopyOpenIM(dst *open_im_sdk.BlackInfo, src imdb.Black) {
	utils.CopyStructFields(dst, src)
	dst.CreateTime = src.CreateTime.Unix()
	user, _ := imdb.GetUserByUserID(src.BlockUserID)
	if user != nil {
		utils.CopyStructFields(dst.BlackUserInfo, user)
	}
}
