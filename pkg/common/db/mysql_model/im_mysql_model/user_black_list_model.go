package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/utils"
	"time"
)

func InsertInToUserBlackList(black db.Black) error {
	black.CreateTime = time.Now()
	return db.DB.MysqlDB.DefaultGormDB().Table("blacks").Create(black).Error
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
	var black db.Black
	return db.DB.MysqlDB.DefaultGormDB().Table("blacks").Where("owner_user_id=? and block_user_id=?", ownerUserID, blockUserID).Find(&black).Error
}

func RemoveBlackList(ownerUserID, blockUserID string) error {
	err := db.DB.MysqlDB.DefaultGormDB().Table("blacks").Where("owner_user_id=? and block_user_id=?", ownerUserID, blockUserID).Delete(db.Black{}).Error
	return utils.Wrap(err, "RemoveBlackList failed")
}

func GetBlackListByUserID(ownerUserID string) ([]db.Black, error) {
	var blackListUsersInfo []db.Black
	err := db.DB.MysqlDB.DefaultGormDB().Table("blacks").Where("owner_user_id=?", ownerUserID).Find(&blackListUsersInfo).Error
	if err != nil {
		return nil, err
	}
	return blackListUsersInfo, nil
}

func GetBlackIDListByUserID(ownerUserID string) ([]string, error) {
	var blackIDList []string
	err := db.DB.MysqlDB.DefaultGormDB().Table("blacks").Where("owner_user_id=?", ownerUserID).Pluck("block_user_id", &blackIDList).Error
	if err != nil {
		return nil, err
	}
	return blackIDList, nil
}
