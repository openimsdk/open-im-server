package convert

import (
	"Open_IM/pkg/common/db/relation"
	"Open_IM/pkg/common/db/table"
	sdk "Open_IM/pkg/proto/sdk_ws"
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
	*table.FriendModel
}

func NewDBFriend(friend *table.FriendModel) *DBFriend {
	return &DBFriend{FriendModel: friend}
}

type PBFriend struct {
	*sdk.FriendInfo
}

func NewPBFriend(friendInfo *sdk.FriendInfo) *PBFriend {
	return &PBFriend{FriendInfo: friendInfo}
}

func (*PBFriend) PB2DB(friends []*sdk.FriendInfo) (DBFriends []*table.FriendModel, err error) {

}

func (*DBFriend) DB2PB(friends []*table.FriendModel) (PBFriends []*sdk.FriendInfo, err error) {

}

func (db *DBFriend) Convert() (*sdk.FriendInfo, error) {
	pbFriend := &sdk.FriendInfo{FriendUser: &sdk.UserInfo{}}
	utils.CopyStructFields(pbFriend, db)
	user, err := getUsersInfo([]string{db.FriendUserID})
	if err != nil {
		return nil, err
	}
	utils.CopyStructFields(pbFriend.FriendUser, user[0])
	pbFriend.CreateTime = db.CreateTime.Unix()

	pbFriend.FriendUser.CreateTime = db.CreateTime.Unix()
	return pbFriend, nil
}

func (pb *PBFriend) Convert() (*table.FriendModel, error) {
	dbFriend := &table.FriendModel{}
	utils.CopyStructFields(dbFriend, pb)
	dbFriend.FriendUserID = pb.FriendUser.UserID
	dbFriend.CreateTime = utils.UnixSecondToTime(pb.CreateTime)
	return dbFriend, nil
}

type DBFriendRequest struct {
	*table.FriendRequestModel
}

func NewDBFriendRequest(friendRequest *table.FriendRequestModel) *DBFriendRequest {
	return &DBFriendRequest{FriendRequestModel: friendRequest}
}

type PBFriendRequest struct {
	*sdk.FriendRequest
}

func NewPBFriendRequest(friendRequest *sdk.FriendRequest) *PBFriendRequest {
	return &PBFriendRequest{FriendRequest: friendRequest}
}

func (*PBFriendRequest) PB2DB(friendRequests []*sdk.FriendRequest) (DBFriendRequests []*table.FriendRequestModel, err error) {

}

func (*DBFriendRequest) DB2PB(friendRequests []*table.FriendRequestModel) (PBFriendRequests []*sdk.FriendRequest, err error) {

}

func (pb *PBFriendRequest) Convert() (*table.FriendRequestModel, error) {
	dbFriendRequest := &table.FriendRequestModel{}
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
	pbFriendRequest.CreateTime = db.CreateTime.Unix()
	pbFriendRequest.HandleTime = db.HandleTime.Unix()
	return pbFriendRequest, nil
}

type DBBlack struct {
	*relation.Black
}

func (*PBBlack) PB2DB(friendRequests []*sdk.BlackInfo) (DBFriendRequests []*table., err error) {

}

func (*DBBlack) DB2PB(friendRequests []*table.FriendRequestModel) (PBFriendRequests []*sdk.FriendRequest, err error) {

}


func NewDBBlack(black *relation.Black) *DBBlack {
	return &DBBlack{Black: black}
}

type PBBlack struct {
	*sdk.BlackInfo
}

func NewPBBlack(blackInfo *sdk.BlackInfo) *PBBlack {
	return &PBBlack{BlackInfo: blackInfo}
}

func (pb *PBBlack) Convert() (*relation.Black, error) {
	dbBlack := &relation.Black{}
	dbBlack.BlockUserID = pb.BlackUserInfo.UserID
	dbBlack.CreateTime = utils.UnixSecondToTime(int64(pb.CreateTime))
	return dbBlack, nil
}
func (db *DBBlack) Convert() (*sdk.BlackInfo, error) {
	pbBlack := &sdk.BlackInfo{}
	utils.CopyStructFields(pbBlack, db)
	pbBlack.CreateTime = db.CreateTime.Unix()
	user, err := getUsersInfo([]string{db.BlockUserID})
	if err != nil {
		return nil, err
	}
	utils.CopyStructFields(pbBlack.BlackUserInfo, user)
	return pbBlack, nil
}

type DBGroup struct {
	*table.GroupModel
}

func NewDBGroup(group *table.GroupModel) *DBGroup {
	return &DBGroup{GroupModel: group}
}

type PBGroup struct {
	*sdk.GroupInfo
}

func NewPBGroup(groupInfo *sdk.GroupInfo) *PBGroup {
	return &PBGroup{GroupInfo: groupInfo}
}

func (pb *PBGroup) Convert() *table.GroupModel {
	dst := &table.GroupModel{}
	_ = utils.CopyStructFields(dst, pb)
	return dst
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
	dst.CreateTime = db.CreateTime.Unix()
	dst.NotificationUpdateTime = db.NotificationUpdateTime.Unix()
	if db.NotificationUpdateTime.Unix() < 0 {
		dst.NotificationUpdateTime = 0
	}
	return dst, nil
}

type DBGroupMember struct {
	*table.GroupMemberModel
}

func NewDBGroupMember(groupMember *table.GroupMemberModel) *DBGroupMember {
	return &DBGroupMember{GroupMemberModel: groupMember}
}

type PBGroupMember struct {
	*sdk.GroupMemberFullInfo
}

func NewPBGroupMember(groupMemberFullInfo *sdk.GroupMemberFullInfo) *PBGroupMember {
	return &PBGroupMember{GroupMemberFullInfo: groupMemberFullInfo}
}

func (pb *PBGroupMember) Convert() (*table.GroupMemberModel, error) {
	dst := &table.GroupMemberModel{}
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

	dst.JoinTime = db.JoinTime.Unix()
	if db.JoinTime.Unix() < 0 {
		dst.JoinTime = 0
	}
	dst.MuteEndTime = db.MuteEndTime.Unix()
	if dst.MuteEndTime < time.Now().Unix() {
		dst.MuteEndTime = 0
	}
	return dst, nil
}

type DBGroupRequest struct {
	*table.GroupRequestModel
}

func NewDBGroupRequest(groupRequest *table.GroupRequestModel) *DBGroupRequest {
	return &DBGroupRequest{GroupRequestModel: groupRequest}
}

type PBGroupRequest struct {
	*sdk.GroupRequest
}

func NewPBGroupRequest(groupRequest *sdk.GroupRequest) *PBGroupRequest {
	return &PBGroupRequest{GroupRequest: groupRequest}
}

func (pb *PBGroupRequest) Convert() (*table.GroupRequestModel, error) {
	dst := &table.GroupRequestModel{}
	utils.CopyStructFields(dst, pb)
	dst.ReqTime = utils.UnixSecondToTime(int64(pb.ReqTime))
	dst.HandledTime = utils.UnixSecondToTime(int64(pb.HandleTime))
	return dst, nil
}
func (db *DBGroupRequest) Convert() (*sdk.GroupRequest, error) {
	dst := &sdk.GroupRequest{}
	utils.CopyStructFields(dst, db)
	dst.ReqTime = db.ReqTime.Unix()
	dst.HandleTime = db.HandledTime.Unix()
	return dst, nil
}

type DBUser struct {
	*table.UserModel
}

func NewDBUser(user *table.UserModel) *DBUser {
	return &DBUser{UserModel: user}
}

type PBUser struct {
	*sdk.UserInfo
}

func NewPBUser(userInfo *sdk.UserInfo) *PBUser {
	return &PBUser{UserInfo: userInfo}
}

func (*PBUser) PB2DB(users []*sdk.UserInfo) (DBUsers []*table.UserModel, err error) {
	for _, v := range users {
		u, err := NewPBUser(v).Convert()
		if err != nil {
			return nil, err
		}
		DBUsers = append(DBUsers, u)
	}
	return
}

func (*DBUser) DB2PB(users []*table.UserModel) (PBUsers []*sdk.UserInfo, err error) {
	for _, v := range users {
		u, err := NewDBUser(v).Convert()
		if err != nil {
			return nil, err
		}
		PBUsers = append(PBUsers, u)
	}
	return
}

func (pb *PBUser) Convert() (*table.UserModel, error) {
	dst := &table.UserModel{}
	utils.CopyStructFields(dst, pb)
	dst.Birth = utils.UnixSecondToTime(pb.Birthday)
	dst.CreateTime = utils.UnixSecondToTime(int64(pb.CreateTime))
	return dst, nil
}

func (db *DBUser) Convert() (*sdk.UserInfo, error) {
	dst := &sdk.UserInfo{}
	utils.CopyStructFields(dst, db)
	dst.CreateTime = db.CreateTime.Unix()
	dst.Birthday = db.Birth.Unix()
	return dst, nil
}

func (db *DBUser) ConvertPublic() (*sdk.PublicUserInfo, error) {
	dst := &sdk.PublicUserInfo{}
	utils.CopyStructFields(dst, db)
	return dst, nil
}
