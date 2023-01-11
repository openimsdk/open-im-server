package im_mysql_model

import (
	"Open_IM/pkg/utils"
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
}

func (FriendRequest) TableName() string {
	return "friend_requests"
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
