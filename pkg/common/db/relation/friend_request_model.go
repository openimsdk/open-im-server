package relation

import (
	"Open_IM/pkg/common/trace_log"
	"Open_IM/pkg/utils"
	"context"
	"gorm.io/gorm"
	"time"
)

var FriendRequestDB *gorm.DB

type FriendRequest struct {
	FromUserID    string    `gorm:"column:from_user_id;primary_key;size:64"`
	ToUserID      string    `gorm:"column:to_user_id;primary_key;size:64"`
	HandleResult  int32     `gorm:"column:handle_result"`
	ReqMsg        string    `gorm:"column:req_msg;size:255"`
	CreateTime    time.Time `gorm:"column:create_time"`
	HandlerUserID string    `gorm:"column:handler_user_id;size:64"`
	HandleMsg     string    `gorm:"column:handle_msg;size:255"`
	HandleTime    time.Time `gorm:"column:handle_time"`
	Ex            string    `gorm:"column:ex;size:1024"`
	DB            *gorm.DB  `gorm:"-"`
}

func NewFriendRequest(db *gorm.DB) *FriendRequest {
	var fr FriendRequest
	fr.DB = initModel(db, &fr)
	return &fr
}

func (f *FriendRequest) Create(ctx context.Context, friends []*FriendRequest) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "friends", friends)
	}()
	return utils.Wrap(f.DB.Create(&friends).Error, "")
}

func (f *FriendRequest) Delete(ctx context.Context, fromUserID, toUserID string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "fromUserID", fromUserID, "toUserID", toUserID)
	}()
	return utils.Wrap(f.DB.Where("from_user_id = ? and to_user_id = ?", fromUserID, toUserID).Delete(&FriendRequest{}).Error, "")
}

func (f *FriendRequest) UpdateByMap(ctx context.Context, ownerUserID string, args map[string]interface{}) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "ownerUserID", ownerUserID, "args", args)
	}()
	return utils.Wrap(f.DB.Where("owner_user_id = ?", ownerUserID).Updates(args).Error, "")
}

func (f *FriendRequest) Update(ctx context.Context, friends []*FriendRequest) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "friends", friends)
	}()
	return utils.Wrap(f.DB.Updates(&friends).Error, "")
}

func (f *FriendRequest) Find(ctx context.Context, ownerUserID string) (friends []*FriendRequest, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "ownerUserID", ownerUserID, "friends", friends)
	}()
	return friends, utils.Wrap(f.DB.Where("owner_user_id = ?", ownerUserID).Find(&friends).Error, "")
}

func (f *FriendRequest) Take(ctx context.Context, fromUserID, toUserID string) (friend *FriendRequest, err error) {
	friend = &FriendRequest{}
	defer tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "fromUserID", fromUserID, "toUserID", toUserID, "friend", friend)
	return friend, utils.Wrap(f.DB.Where("from_user_id = ? and to_user_id", fromUserID, toUserID).Take(friend).Error, "")
}

func (f *FriendRequest) FindToUserID(ctx context.Context, toUserID string) (friends []*FriendRequest, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "toUserID", toUserID, "friends", friends)
	}()
	return friends, utils.Wrap(f.DB.Where("to_user_id = ?", toUserID).Find(&friends).Error, "")
}

func (f *FriendRequest) FindFromUserID(ctx context.Context, fromUserID string) (friends []*FriendRequest, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "fromUserID", fromUserID, "friends", friends)
	}()
	return friends, utils.Wrap(f.DB.Where("from_user_id = ?", fromUserID).Find(&friends).Error, "")
}

// who apply to add me
func GetReceivedFriendsApplicationListByUserID(ToUserID string) ([]FriendRequest, error) {
	var usersInfo []FriendRequest
	err := FriendRequestDB.Table("friend_requests").Where("to_user_id=?", ToUserID).Find(&usersInfo).Error
	if err != nil {
		return nil, err
	}
	return usersInfo, nil
}

// I apply to add somebody
func GetSendFriendApplicationListByUserID(FromUserID string) ([]FriendRequest, error) {
	var usersInfo []FriendRequest
	err := FriendRequestDB.Table("friend_requests").Where("from_user_id=?", FromUserID).Find(&usersInfo).Error
	if err != nil {
		return nil, err
	}
	return usersInfo, nil
}

// FromUserId apply to add ToUserID
func GetFriendApplicationByBothUserID(FromUserID, ToUserID string) (*FriendRequest, error) {
	var friendRequest FriendRequest
	err := FriendRequestDB.Table("friend_requests").Where("from_user_id=? and to_user_id=?", FromUserID, ToUserID).Take(&friendRequest).Error
	if err != nil {
		return nil, err
	}
	return &friendRequest, nil
}

func UpdateFriendApplication(friendRequest *FriendRequest) error {
	friendRequest.CreateTime = time.Now()
	return FriendRequestDB.Table("friend_requests").Where("from_user_id=? and to_user_id=?",
		friendRequest.FromUserID, friendRequest.ToUserID).Updates(&friendRequest).Error
}

func InsertFriendApplication(friendRequest *FriendRequest, args map[string]interface{}) error {
	if err := FriendRequestDB.Table("friend_requests").Create(friendRequest).Error; err == nil {
		return nil
	}

	//t := dbConn.Debug().Table("friend_requests").Where("from_user_id = ? and to_user_id = ?", friendRequest.FromUserID, friendRequest.ToUserID).Select("*").Updates(*friendRequest)
	//if t.RowsAffected == 0 {
	//	return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	//}
	//return utils.Wrap(t.Error, "")

	friendRequest.CreateTime = time.Now()
	args["create_time"] = friendRequest.CreateTime
	u := FriendRequestDB.Model(friendRequest).Updates(args)
	//u := dbConn.Table("friend_requests").Where("from_user_id=? and to_user_id=?",
	// friendRequest.FromUserID, friendRequest.ToUserID).Update(&friendRequest)
	//u := dbConn.Table("friend_requests").Where("from_user_id=? and to_user_id=?",
	//	friendRequest.FromUserID, friendRequest.ToUserID).Update(&friendRequest)
	if u.RowsAffected != 0 {
		return nil
	}

	if friendRequest.CreateTime.Unix() < 0 {
		friendRequest.CreateTime = time.Now()
	}
	if friendRequest.HandleTime.Unix() < 0 {
		friendRequest.HandleTime = utils.UnixSecondToTime(0)
	}
	err := FriendRequestDB.Table("friend_requests").Create(friendRequest).Error
	if err != nil {
		return err
	}
	return nil
}
