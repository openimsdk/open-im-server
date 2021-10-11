package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"time"
)

func InsertIntoGroupRequest(groupId, fromUserId, toUserId, reqMsg, fromUserNickName, fromUserFaceUrl string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	toInsertInfo := GroupRequest{GroupID: groupId, FromUserID: fromUserId, ToUserID: toUserId, ReqMsg: reqMsg, FromUserNickname: fromUserNickName, FromUserFaceUrl: fromUserFaceUrl, CreateTime: time.Now()}
	err = dbConn.Table("group_request").Create(&toInsertInfo).Error
	if err != nil {
		return err
	}
	return nil
}

func FindGroupRequestUserInfoByGroupIDAndUid(groupId, uid string) (*GroupRequest, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var requestUserInfo GroupRequest
	err = dbConn.Table("group_request").Where("from_user_id=? and group_id=?", uid, groupId).Find(&requestUserInfo).Error
	if err != nil {
		return nil, err
	}
	return &requestUserInfo, nil
}

func DelGroupRequest(groupId, fromUserId, toUserId string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	err = dbConn.Exec("delete from group_request where group_id=? and from_user_id=? and to_user_id=?", groupId, fromUserId, toUserId).Error
	if err != nil {
		return err
	}
	return nil
}

func FindGroupBeInvitedRequestInfoByUidAndGroupID(groupId, uid string) (*GroupRequest, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var beInvitedRequestUserInfo GroupRequest
	err = dbConn.Table("group_request").Where("to_user_id=? and group_id=?", uid, groupId).Find(&beInvitedRequestUserInfo).Error
	if err != nil {
		return nil, err
	}
	return &beInvitedRequestUserInfo, nil

}

func InsertGroupRequest(groupId, fromUser, fromUserNickName, fromUserFaceUrl, toUser, requestMsg, handledMsg string, handleStatus int) error {
	return nil
}
