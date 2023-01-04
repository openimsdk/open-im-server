package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"fmt"
	"time"
)

func InsertToFriend(toInsertFollow *Friend) error {
	toInsertFollow.CreateTime = time.Now()
	err := db.DB.MysqlDB.DefaultGormDB().Table("friends").Create(toInsertFollow).Error
	if err != nil {
		return err
	}
	return nil
}

func GetFriendRelationshipFromFriend(OwnerUserID, FriendUserID string) (*Friend, error) {
	var friend Friend
	err := db.DB.MysqlDB.DefaultGormDB().Table("friends").Where("owner_user_id=? and friend_user_id=?", OwnerUserID, FriendUserID).Take(&friend).Error
	if err != nil {
		return nil, err
	}
	return &friend, err
}

func GetFriendListByUserID(OwnerUserID string) ([]Friend, error) {
	var friends []Friend
	var x Friend
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
	return db.DB.MysqlDB.DefaultGormDB().Table("friends").Where("owner_user_id=? and friend_user_id=?", OwnerUserID, FriendUserID).Delete(Friend{}).Error
}

type FriendUser struct {
	Friend
	Nickname string `gorm:"column:name;size:255"`
}

func GetUserFriendsCMS(ownerUserID, friendUserName string, pageNumber, showNumber int32) (friendUserList []*FriendUser, count int64, err error) {
	db := db.DB.MysqlDB.DefaultGormDB().Table("friends").
		Select("friends.*, users.name").
		Where("friends.owner_user_id=?", ownerUserID).Limit(int(showNumber)).
		Joins("left join users on friends.friend_user_id = users.user_id").
		Offset(int(showNumber * (pageNumber - 1)))
	if friendUserName != "" {
		db = db.Where("users.name like ?", fmt.Sprintf("%%%s%%", friendUserName))
	}
	if err = db.Count(&count).Error; err != nil {
		return
	}
	err = db.Find(&friendUserList).Error
	return
}

func GetFriendByIDCMS(ownerUserID, friendUserID string) (friendUser *FriendUser, err error) {
	friendUser = &FriendUser{}
	err = db.DB.MysqlDB.DefaultGormDB().Table("friends").
		Select("friends.*, users.name").
		Where("friends.owner_user_id=? and friends.friend_user_id=?", ownerUserID, friendUserID).
		Joins("left join users on friends.friend_user_id = users.user_id").
		Take(friendUser).Error
	return friendUser, err
}
