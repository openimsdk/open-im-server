package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

func InsertToFriend(toInsertFollow *db.Friend) error {
	toInsertFollow.CreateTime = time.Now()
	err := db.DB.MysqlDB.DefaultGormDB().Table("friends").Create(toInsertFollow).Error
	if err != nil {
		return err
	}
	return nil
}

func GetFriendRelationshipFromFriend(OwnerUserID, FriendUserID string) (*db.Friend, error) {
	var friend db.Friend
	err := db.DB.MysqlDB.DefaultGormDB().Table("friends").Where("owner_user_id=? and friend_user_id=?", OwnerUserID, FriendUserID).Take(&friend).Error
	if err != nil {
		return nil, err
	}
	return &friend, err
}

func GetFriendListByUserID(OwnerUserID string) ([]db.Friend, error) {
	var friends []db.Friend
	var x db.Friend
	x.OwnerUserID = OwnerUserID
	err := db.DB.MysqlDB.DefaultGormDB().Table("friends").Where("owner_user_id=?", OwnerUserID).Find(&friends).Error
	if err != nil {
		return nil, err
	}
	return friends, nil
}

func GetFriendIDListByUserID(OwnerUserID string) ([]string, error) {
	var friendIDList []string
	err := db.DB.MysqlDB.DefaultGormDB().Table("friends").Where("owner_user_id=?", OwnerUserID).Pluck("friend_user_id", &friendIDList).Error
	if err != nil {
		return nil, err
	}
	return friendIDList, nil
}

func UpdateFriendComment(OwnerUserID, FriendUserID, Remark string) error {
	return db.DB.MysqlDB.DefaultGormDB().Exec("update friends set remark=? where owner_user_id=? and friend_user_id=?", Remark, OwnerUserID, FriendUserID).Error
}

func DeleteSingleFriendInfo(OwnerUserID, FriendUserID string) error {
	return db.DB.MysqlDB.DefaultGormDB().Table("friends").Where("owner_user_id=? and friend_user_id=?", OwnerUserID, FriendUserID).Delete(db.Friend{}).Error
}
