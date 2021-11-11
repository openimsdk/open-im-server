package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"time"
)

func ReplaceIntoFriendReq(reqId, userId string, flag int32, reqMessage string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	err = dbConn.Exec("replace into friend_request(req_id,user_id,flag,req_message,create_time) values(?,?,?,?,?)", reqId, userId, flag, reqMessage, time.Now()).Error
	if err != nil {
		return err
	}
	return nil
}

func FindFriendsApplyFromFriendReq(userId string) ([]FriendRequest, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var usersInfo []FriendRequest
	//dbConn.LogMode(true)
	err = dbConn.Table("friend_request").Where("user_id=?", userId).Find(&usersInfo).Error
	if err != nil {
		return nil, err
	}
	return usersInfo, nil
}

func FindSelfApplyFromFriendReq(userId string) ([]FriendRequest, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var usersInfo []FriendRequest
	err = dbConn.Table("friend_request").Where("req_id=?", userId).Find(&usersInfo).Error
	if err != nil {
		return nil, err
	}
	return usersInfo, nil
}

func FindFriendApplyFromFriendReqByUid(reqId, userId string) (*FriendRequest, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var friendRequest FriendRequest
	err = dbConn.Table("friend_request").Where("req_id=? and user_id=?", reqId, userId).Find(&friendRequest).Error
	if err != nil {
		return nil, err
	}
	return &friendRequest, nil
}

func UpdateFriendRelationshipToFriendReq(reqId, userId string, flag int32) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	err = dbConn.Exec("update friend_request set flag=? where req_id=? and user_id=?", flag, reqId, userId).Error
	if err != nil {
		return err
	}
	return nil
}
