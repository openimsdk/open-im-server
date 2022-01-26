package im_mysql_model

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/utils"
	"time"
)

//type Group struct {
//	GroupID       string    `gorm:"column:group_id;primaryKey;"`
//	GroupName     string    `gorm:"column:name"`
//	Introduction  string    `gorm:"column:introduction"`
//	Notification  string    `gorm:"column:notification"`
//	FaceUrl       string    `gorm:"column:face_url"`
//	CreateTime    time.Time `gorm:"column:create_time"`
//	Status        int32     `gorm:"column:status"`
//	CreatorUserID string    `gorm:"column:creator_user_id"`
//	GroupType     int32     `gorm:"column:group_type"`
//	Ex            string    `gorm:"column:ex"`
//}

func InsertIntoGroup(groupInfo db.Group) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	if groupInfo.GroupName == "" {
		groupInfo.GroupName = "Group Chat"
	}
	groupInfo.CreateTime = time.Now()
	err = dbConn.Table("groups").Create(groupInfo).Error
	if err != nil {
		return err
	}
	return nil
}

func GetGroupInfoByGroupID(groupId string) (*db.Group, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var groupInfo db.Group
	err = dbConn.Table("groups").Where("group_id=?", groupId).Find(&groupInfo).Error
	if err != nil {
		return nil, err
	}
	return &groupInfo, nil
}

func SetGroupInfo(groupInfo db.Group) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	dbConn.LogMode(true)
	err = dbConn.Table("groups").Where("group_id=?", groupInfo.GroupID).Update(&groupInfo).Error
	return err
}

func GetGroupsByName(groupName string, pageNumber, showNumber int32) ([]db.Group, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	var groups []db.Group
	if err != nil {
		return groups, err
	}
	dbConn.LogMode(true)
	err = dbConn.Table("groups").Where("name=?", groupName).Limit(showNumber).Offset(showNumber * (pageNumber - 1)).Find(&groups).Error
	return groups, err
}

func GetGroups(pageNumber, showNumber int) ([]db.Group, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	var groups []db.Group
	if err != nil {
		return groups, err
	}
	dbConn.LogMode(true)
	if err = dbConn.Table("groups").Limit(showNumber).Offset(showNumber * (pageNumber - 1)).Find(&groups).Error; err != nil {
		return groups, err
	}
	return groups, nil
}

func BanGroupChat(groupId string) error {
	var group db.Group
	group.Status = constant.GroupBanChat
	if err := SetGroupInfo(group); err != nil {
		return err
	}
	return nil
}

func BanPrivateChat(groupId string) error {
	var group db.Group
	group.Status = constant.GroupBanPrivateChat
	if err := SetGroupInfo(group); err != nil {
		return err
	}
	return nil
}

func DeleteGroup(groupId string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	dbConn.LogMode(true)
	var group db.Group
	if err := dbConn.Table("groups").Where("").Delete(&group).Error; err != nil {
		return err
	}
	return nil
}

func SetGroupMaster(userId, groupId string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	dbConn.LogMode(true)
	groupMember := db.GroupMember{
		UserID:  userId,
		GroupID: groupId,
	}
	updateInfo := db.GroupMember{
		RoleLevel: constant.GroupOwner,
	}
	if err := dbConn.Find(&groupMember).Update(updateInfo).Error; err != nil {
		return err
	}
	return nil
}

func GetGroupsCountNum() (int, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return 0, err
	}
	dbConn.LogMode(true)
	var count int
	if err := dbConn.Model(&db.Group{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
