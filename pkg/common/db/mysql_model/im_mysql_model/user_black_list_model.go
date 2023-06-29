package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/utils"
	"time"
)

func InsertInToUserBlackList(black db.Black) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	black.CreateTime = time.Now()
	err = dbConn.Table("blacks").Create(black).Error
	return err
}

//type Black struct {
//	OwnerUserID    string    `gorm:"column:owner_user_id;primaryKey;"`
//	BlockUserID    string    `gorm:"column:block_user_id;primaryKey;"`
//	CreateTime     time.Time `gorm:"column:create_time"`
//	AddSource      int32     `gorm:"column:add_source"`
//	OperatorUserID int32     `gorm:"column:operator_user_id"`
//	Ex             string    `gorm:"column:ex"`
//}

func CheckBlack(ownerUserID, blockUserID string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	var black db.Black
	err = dbConn.Table("blacks").Where("owner_user_id=? and block_user_id=?", ownerUserID, blockUserID).Find(&black).Error
	return err
}

func RemoveBlackList(ownerUserID, blockUserID string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	err = dbConn.Table("blacks").Where("owner_user_id=? and block_user_id=?", ownerUserID, blockUserID).Delete(db.Black{}).Error
	return utils.Wrap(err, "RemoveBlackList failed")
}

func GetBlackListByUserID(ownerUserID string) ([]db.Black, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var blackListUsersInfo []db.Black
	err = dbConn.Table("blacks").Where("owner_user_id=?", ownerUserID).Find(&blackListUsersInfo).Error
	if err != nil {
		return nil, err
	}
	return blackListUsersInfo, nil
}
