package im_mysql_model

import (
	"Open_IM/pkg/common/db"
)

//type FriendRequest struct {
//	FromUserID    string    `gorm:"column:from_user_id;primaryKey;"`
//	ToUserID      string    `gorm:"column:to_user_id;primaryKey;"`
//	HandleResult  int32     `gorm:"column:handle_result"`
//	ReqMessage    string    `gorm:"column:req_message"`
//	CreateTime    time.Time `gorm:"column:create_time"`
//	HandlerUserID string    `gorm:"column:handler_user_id"`
//	HandleMsg     string    `gorm:"column:handle_msg"`
//	HandleTime    time.Time `gorm:"column:handle_time"`
//	Ex            string    `gorm:"column:ex"`
//}

// who apply to add me
func GetReceivedFriendsApplicationListByUserID(ToUserID string) ([]FriendRequest, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var usersInfo []FriendRequest
	err = dbConn.Table("friend_request").Where("to_user_id=?", ToUserID).Find(&usersInfo).Error
	if err != nil {
		return nil, err
	}
	return usersInfo, nil
}

//I apply to add somebody
func GetSendFriendApplicationListByUserID(FromUserID string) ([]FriendRequest, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var usersInfo []FriendRequest
	err = dbConn.Table("friend_request").Where("from_user_id=?", FromUserID).Find(&usersInfo).Error
	if err != nil {
		return nil, err
	}
	return usersInfo, nil
}

//reqId apply to add userId already
func FindFriendApplicationByBothUserID(FromUserId, ToUserID string) (*FriendRequest, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var friendRequest FriendRequest
	err = dbConn.Table("friend_request").Where("from_user_id=? and to_user_id=?", FromUserId, ToUserID).Find(&friendRequest).Error
	if err != nil {
		return nil, err
	}
	return &friendRequest, nil
}

func UpdateFriendApplication(friendRequest FriendRequest) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	err = dbConn.Table("friend_request").Where("from_user_id=? and to_user_id=?", friendRequest.FromUserID, friendRequest.ToUserID).Update(&friendRequest).Error
	if err != nil {
		return err
	}
	return nil
}
