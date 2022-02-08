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
	dst.CreateTime = utils.UnixSecondToTime(int64(src.CreateTime))
}

func FriendDBCopyOpenIM(dst *open_im_sdk.FriendInfo, src *db.Friend) error {
	utils.CopyStructFields(dst, src)
	user, err := imdb.GetUserByUserID(src.FriendUserID)
	if err != nil {
		return utils.Wrap(err, "")
	}
	utils.CopyStructFields(dst.FriendUser, user)
	dst.CreateTime = uint32(src.CreateTime.Unix())
	if dst.FriendUser == nil {
		dst.FriendUser = &open_im_sdk.UserInfo{}
	}
	dst.FriendUser.CreateTime = uint32(user.CreateTime.Unix())
	return nil
}

//
func FriendRequestOpenIMCopyDB(dst *db.FriendRequest, src *open_im_sdk.FriendRequest) {
	utils.CopyStructFields(dst, src)
	dst.CreateTime = utils.UnixSecondToTime(int64(src.CreateTime))
	dst.HandleTime = utils.UnixSecondToTime(int64(src.HandleTime))
}

func FriendRequestDBCopyOpenIM(dst *open_im_sdk.FriendRequest, src *db.FriendRequest) error {
	utils.CopyStructFields(dst, src)
	user, err := imdb.GetUserByUserID(src.FromUserID)
	if err != nil {
		return utils.Wrap(err, "")
	}
	dst.FromNickname = user.Nickname
	dst.FromFaceURL = user.FaceURL
	dst.FromGender = user.Gender
	user, err = imdb.GetUserByUserID(src.ToUserID)
	if err != nil {
		return utils.Wrap(err, "")
	}
	dst.ToNickname = user.Nickname
	dst.ToFaceURL = user.FaceURL
	dst.ToGender = user.Gender
	dst.CreateTime = uint32(src.CreateTime.Unix())
	dst.HandleTime = uint32(src.HandleTime.Unix())
	return nil
}

func BlackOpenIMCopyDB(dst *db.Black, src *open_im_sdk.BlackInfo) {
	utils.CopyStructFields(dst, src)
	dst.BlockUserID = src.BlackUserInfo.UserID
	dst.CreateTime = utils.UnixSecondToTime(int64(src.CreateTime))
}

func BlackDBCopyOpenIM(dst *open_im_sdk.BlackInfo, src *db.Black) error {
	utils.CopyStructFields(dst, src)
	dst.CreateTime = uint32(src.CreateTime.Unix())
	user, err := imdb.GetUserByUserID(src.BlockUserID)
	if err != nil {
		return utils.Wrap(err, "")
	}
	utils.CopyStructFields(dst.BlackUserInfo, user)
	return nil
}

func GroupOpenIMCopyDB(dst *db.Group, src *open_im_sdk.GroupInfo) {
	utils.CopyStructFields(dst, src)
}

func GroupDBCopyOpenIM(dst *open_im_sdk.GroupInfo, src *db.Group) error {
	utils.CopyStructFields(dst, src)
	user, err := imdb.GetGroupOwnerInfoByGroupID(src.GroupID)
	if err != nil {
		return utils.Wrap(err, "")
	}
	dst.OwnerUserID = user.UserID

	dst.MemberCount, err = imdb.GetGroupMemberNumByGroupID(src.GroupID)
	if err != nil {
		return utils.Wrap(err, "")
	}
	dst.CreateTime = uint32(src.CreateTime.Unix())
	return nil
}

func GroupMemberOpenIMCopyDB(dst *db.GroupMember, src *open_im_sdk.GroupMemberFullInfo) {
	utils.CopyStructFields(dst, src)
}

func GroupMemberDBCopyOpenIM(dst *open_im_sdk.GroupMemberFullInfo, src *db.GroupMember) error {
	utils.CopyStructFields(dst, src)
	if token_verify.IsMangerUserID(src.UserID) {
		u, err := imdb.GetUserByUserID(src.UserID)
		if err != nil {
			return utils.Wrap(err, "")
		}

		utils.CopyStructFields(dst, u)

		dst.AppMangerLevel = 1
	}
	dst.JoinTime = src.JoinTime.Unix()
	return nil
}

func GroupRequestOpenIMCopyDB(dst *db.GroupRequest, src *open_im_sdk.GroupRequest) {
	utils.CopyStructFields(dst, src)
}

func GroupRequestDBCopyOpenIM(dst *open_im_sdk.GroupRequest, src *db.GroupRequest) {
	utils.CopyStructFields(dst, src)
	dst.ReqTime = uint32(src.ReqTime.Unix())
	dst.HandleTime = uint32(src.HandledTime.Unix())
}

func UserOpenIMCopyDB(dst *db.User, src *open_im_sdk.UserInfo) {
	utils.CopyStructFields(dst, src)
	dst.Birth = utils.UnixSecondToTime(int64(src.Birth))
	dst.CreateTime = utils.UnixSecondToTime(int64(src.CreateTime))
}

func UserDBCopyOpenIM(dst *open_im_sdk.UserInfo, src *db.User) {
	utils.CopyStructFields(dst, src)
	dst.CreateTime = uint32(src.CreateTime.Unix())
	dst.Birth = uint32(src.Birth.Unix())
}

func UserDBCopyOpenIMPublicUser(dst *open_im_sdk.PublicUserInfo, src *db.User) {
	utils.CopyStructFields(dst, src)
}

//
//func PublicUserDBCopyOpenIM(dst *open_im_sdk.PublicUserInfo, src *db.User){
//
//}
