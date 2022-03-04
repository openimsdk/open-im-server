package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/utils"
	"time"
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
func GetReceivedFriendsApplicationListByUserID(ToUserID string) ([]db.FriendRequest, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var usersInfo []db.FriendRequest
	err = dbConn.Table("friend_requests").Where("to_user_id=?", ToUserID).Find(&usersInfo).Error
	if err != nil {
		return nil, err
	}
	return usersInfo, nil
}

//I apply to add somebody
func GetSendFriendApplicationListByUserID(FromUserID string) ([]db.FriendRequest, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var usersInfo []db.FriendRequest
	err = dbConn.Table("friend_requests").Where("from_user_id=?", FromUserID).Find(&usersInfo).Error
	if err != nil {
		return nil, err
	}
	return usersInfo, nil
}

//FromUserId apply to add ToUserID
func GetFriendApplicationByBothUserID(FromUserID, ToUserID string) (*db.FriendRequest, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var friendRequest db.FriendRequest
	err = dbConn.Table("friend_requests").Where("from_user_id=? and to_user_id=?", FromUserID, ToUserID).Take(&friendRequest).Error
	if err != nil {
		return nil, err
	}
	return &friendRequest, nil
}

func UpdateFriendApplication(friendRequest *db.FriendRequest) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	friendRequest.CreateTime = time.Now()

	return dbConn.Table("friend_requests").Where("from_user_id=? and to_user_id=?",
		friendRequest.FromUserID, friendRequest.ToUserID).Update(&friendRequest).Error
}

func InsertFriendApplication(friendRequest *db.FriendRequest, args map[string]interface{}) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}

	if err = dbConn.Table("friend_requests").Create(friendRequest).Error; err == nil {
		return nil
	}

	//t := dbConn.Debug().Table("friend_requests").Where("from_user_id = ? and to_user_id = ?", friendRequest.FromUserID, friendRequest.ToUserID).Select("*").Updates(*friendRequest)
	//if t.RowsAffected == 0 {
	//	return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	//}
	//return utils.Wrap(t.Error, "")

	friendRequest.CreateTime = time.Now()
	args["create_time"] = friendRequest.CreateTime
	u := dbConn.Model(friendRequest).Updates(args)
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
	err = dbConn.Table("friend_requests").Create(friendRequest).Error
	if err != nil {
		return err
	}
	return nil
}
