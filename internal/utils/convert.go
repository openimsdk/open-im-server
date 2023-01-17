package utils

import (
	imdb "Open_IM/pkg/common/db/mysql"
	sdk "Open_IM/pkg/proto/sdk_ws"
	utils2 "Open_IM/pkg/utils"
	utils "github.com/OpenIMSDK/open_utils"
	"time"
)

func getUsersInfo(userIDs []string) ([]*sdk.UserInfo, error) {
	return nil, nil
}

func getGroupOwnerInfo(groupID string) (*sdk.GroupMemberFullInfo, error) {
	return nil, nil
}
func getNumberOfGroupMember(groupID string) (int32, error) {
	return 0, nil
}

type DBFriend struct {
	*imdb.Friend
}

type PBFriend struct {
	*sdk.FriendInfo
}

func (db *DBFriend) convert() (*sdk.FriendInfo, error) {
	pbFriend := &sdk.FriendInfo{FriendUser: &sdk.UserInfo{}}
	utils.CopyStructFields(pbFriend, db)
	user, err := getUsersInfo([]string{db.FriendUserID})
	if err != nil {
		return nil, err
	}
	utils2.CopyStructFields(pbFriend.FriendUser, user[0])
	pbFriend.CreateTime = uint32(db.CreateTime.Unix())

	pbFriend.FriendUser.CreateTime = uint32(db.CreateTime.Unix())
	return pbFriend, nil
}

func (pb *PBFriend) Convert() (*imdb.Friend, error) {
	dbFriend := &imdb.Friend{}
	utils2.CopyStructFields(dbFriend, pb)
	dbFriend.FriendUserID = pb.FriendUser.UserID
	dbFriend.CreateTime = utils2.UnixSecondToTime(int64(pb.CreateTime))
	return dbFriend, nil
}

type DBFriendRequest struct {
	*imdb.FriendRequest
}

type PBFriendRequest struct {
	*sdk.FriendRequest
}

func (pb *PBFriendRequest) Convert() (*imdb.FriendRequest, error) {
	dbFriendRequest := &imdb.FriendRequest{}
	utils.CopyStructFields(dbFriendRequest, pb)
	dbFriendRequest.CreateTime = utils.UnixSecondToTime(int64(pb.CreateTime))
	dbFriendRequest.HandleTime = utils.UnixSecondToTime(int64(pb.HandleTime))
	return dbFriendRequest, nil
}
func (db *DBFriendRequest) Convert() (*sdk.FriendRequest, error) {
	pbFriendRequest := &sdk.FriendRequest{}
	utils.CopyStructFields(pbFriendRequest, db)
	user, err := getUsersInfo([]string{db.FromUserID})
	if err != nil {
		return nil, err
	}
	pbFriendRequest.FromNickname = user[0].Nickname
	pbFriendRequest.FromFaceURL = user[0].FaceURL
	pbFriendRequest.FromGender = user[0].Gender
	user, err = getUsersInfo([]string{db.ToUserID})
	if err != nil {
		return nil, err
	}
	pbFriendRequest.ToNickname = user[0].Nickname
	pbFriendRequest.ToFaceURL = user[0].FaceURL
	pbFriendRequest.ToGender = user[0].Gender
	pbFriendRequest.CreateTime = uint32(db.CreateTime.Unix())
	pbFriendRequest.HandleTime = uint32(db.HandleTime.Unix())
	return pbFriendRequest, nil
}

type DBBlack struct {
	*imdb.Black
}

type PBBlack struct {
	*sdk.BlackInfo
}

func (pb *PBBlack) Convert() (*imdb.Black, error) {
	dbBlack := &imdb.Black{}
	dbBlack.BlockUserID = pb.BlackUserInfo.UserID
	dbBlack.CreateTime = utils.UnixSecondToTime(int64(pb.CreateTime))
	return dbBlack, nil
}
func (db *DBBlack) Convert() (*sdk.BlackInfo, error) {
	pbBlack := &sdk.BlackInfo{}
	utils.CopyStructFields(pbBlack, db)
	pbBlack.CreateTime = uint32(db.CreateTime.Unix())
	user, err := getUsersInfo([]string{db.BlockUserID})
	if err != nil {
		return nil, err
	}
	utils.CopyStructFields(pbBlack.BlackUserInfo, user)
	return pbBlack, nil
}

type DBGroup struct {
	*imdb.Group
}

type PBGroup struct {
	*sdk.GroupInfo
}

func (pb *PBGroup) Convert() (*imdb.Group, error) {
	dst := &imdb.Group{}
	utils.CopyStructFields(dst, pb)
	return dst, nil
}
func (db *DBGroup) Convert() (*sdk.GroupInfo, error) {
	dst := &sdk.GroupInfo{}
	utils.CopyStructFields(dst, db)
	user, err := getGroupOwnerInfo(db.GroupID)
	if err != nil {
		return nil, err
	}
	dst.OwnerUserID = user.UserID

	memberCount, err := getNumberOfGroupMember(db.GroupID)
	if err != nil {
		return nil, err
	}
	dst.MemberCount = uint32(memberCount)
	dst.CreateTime = uint32(db.CreateTime.Unix())
	dst.NotificationUpdateTime = uint32(db.NotificationUpdateTime.Unix())
	if db.NotificationUpdateTime.Unix() < 0 {
		dst.NotificationUpdateTime = 0
	}
	return dst, nil
}

type DBGroupMember struct {
	*imdb.GroupMember
}

type PBGroupMember struct {
	*sdk.GroupMemberFullInfo
}

func (pb *PBGroupMember) Convert() (*imdb.GroupMember, error) {
	dst := &imdb.GroupMember{}
	utils.CopyStructFields(dst, pb)
	dst.JoinTime = utils.UnixSecondToTime(int64(pb.JoinTime))
	dst.MuteEndTime = utils.UnixSecondToTime(int64(pb.MuteEndTime))
	return dst, nil
}
func (db *DBGroupMember) Convert() (*sdk.GroupMemberFullInfo, error) {
	dst := &sdk.GroupMemberFullInfo{}
	utils.CopyStructFields(dst, db)

	user, err := getUsersInfo([]string{db.UserID})
	if err != nil {
		return nil, err
	}
	dst.AppMangerLevel = user[0].AppMangerLevel

	dst.JoinTime = int32(db.JoinTime.Unix())
	if db.JoinTime.Unix() < 0 {
		dst.JoinTime = 0
	}
	dst.MuteEndTime = uint32(db.MuteEndTime.Unix())
	if dst.MuteEndTime < uint32(time.Now().Unix()) {
		dst.MuteEndTime = 0
	}
	return dst, nil
}

type DBGroupRequest struct {
	*imdb.GroupRequest
}

type PBGroupRequest struct {
	*sdk.GroupRequest
}

func (pb *PBGroupRequest) Convert() (*imdb.GroupRequest, error) {
	dst := &imdb.GroupRequest{}
	utils.CopyStructFields(dst, pb)
	dst.ReqTime = utils.UnixSecondToTime(int64(pb.ReqTime))
	dst.HandledTime = utils.UnixSecondToTime(int64(pb.HandleTime))
	return dst, nil
}
func (db *DBGroupRequest) Convert() (*sdk.GroupRequest, error) {
	dst := &sdk.GroupRequest{}
	utils.CopyStructFields(dst, db)
	dst.ReqTime = uint32(db.ReqTime.Unix())
	dst.HandleTime = uint32(db.HandledTime.Unix())
	return dst, nil
}

type DBUser struct {
	*imdb.User
}

type PBUser struct {
	*sdk.UserInfo
}

func (pb *PBUser) Convert() (*imdb.User, error) {
	dst := &imdb.User{}
	utils.CopyStructFields(dst, pb)
	dst.Birth = utils.UnixSecondToTime(pb.Birthday)
	dst.CreateTime = utils.UnixSecondToTime(int64(pb.CreateTime))
	return dst, nil
}

func (db *DBUser) Convert() (*sdk.UserInfo, error) {
	dst := &sdk.UserInfo{}
	utils.CopyStructFields(dst, db)
	dst.CreateTime = uint32(db.CreateTime.Unix())
	dst.Birthday = db.Birth.Unix()
	return dst, nil
}

func (db *DBUser) ConvertPublic() (*sdk.PublicUserInfo, error) {
	dst := &sdk.PublicUserInfo{}
	utils.CopyStructFields(dst, db)
	return dst, nil
}
