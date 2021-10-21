package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

func InsertToFriend(ownerId, friendId string, flag int32) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	toInsertFollow := Friend{
		OwnerId:    ownerId,
		FriendId:   friendId,
		FriendFlag: flag,
		CreateTime: time.Now(),
	}
	err = dbConn.Table("friend").Create(toInsertFollow).Error
	if err != nil {
		return err
	}
	return nil
}

func FindFriendRelationshipFromFriend(ownerId, friendId string) (*Friend, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var friend Friend
	err = dbConn.Table("friend").Where("owner_id=? and friend_id=?", ownerId, friendId).Find(&friend).Error
	if err != nil {
		return nil, err
	}
	return &friend, err
}

func FindUserInfoFromFriend(ownerId string) ([]Friend, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var friends []Friend
	err = dbConn.Table("friend").Where("owner_id=?", ownerId).Find(&friends).Error
	if err != nil {
		return nil, err
	}
	return friends, nil
}

func UpdateFriendComment(ownerId, friendId, comment string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	err = dbConn.Exec("update friend set comment=? where owner_id=? and friend_id=?", comment, ownerId, friendId).Error
	return err
}

func DeleteSingleFriendInfo(ownerId, friendId string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	err = dbConn.Table("friend").Where("owner_id=? and friend_id=?", ownerId, friendId).Delete(Friend{}).Error
	return err
}
