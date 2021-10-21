package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"time"
)

func InsertInToUserBlackList(ownerID, blockID string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	toInsertInfo := BlackList{OwnerId: ownerID, BlockId: blockID, CreateTime: time.Now()}
	err = dbConn.Table("user_black_list").Create(toInsertInfo).Error
	return err
}

func FindRelationshipFromBlackList(ownerID, blockID string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	var blackList BlackList
	err = dbConn.Table("user_black_list").Where("owner_id=? and block_id=?", ownerID, blockID).Find(&blackList).Error
	return err
}

func RemoveBlackList(ownerID, blockID string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	err = dbConn.Exec("delete from user_black_list where owner_id=? and block_id=?", ownerID, blockID).Error
	return err
}

func GetBlackListByUID(ownerID string) ([]BlackList, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var blackListUsersInfo []BlackList
	err = dbConn.Table("user_black_list").Where("owner_id=?", ownerID).Find(&blackListUsersInfo).Error
	if err != nil {
		return nil, err
	}
	return blackListUsersInfo, nil
}
