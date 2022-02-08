package im_mysql_model

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/utils"
	"fmt"
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
	err = dbConn.Table("groups").Where(fmt.Sprintf(" name like '%%%s%%' ", groupName)).Limit(showNumber).Offset(showNumber * (pageNumber - 1)).Find(&groups).Error
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


func OperateGroupStatus(groupId string, groupStatus int32) error {
	group := db.Group{
		GroupID: groupId,
		Status: groupStatus,
	}
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
	var groupMembers []db.GroupMember
	if err := dbConn.Table("groups").Where("group_id=?", groupId).Delete(&group).Error; err != nil {
		return err
	}
	if err := dbConn.Table("group_members").Where("group_id=?", groupId).Delete(groupMembers).Error; err != nil {
		return err
	}
	return nil
}

func OperateGroupRole(userId, groupId string, roleLevel int32) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	dbConn.LogMode(true)
	groupMember := db.GroupMember{
		UserID:  userId,
		GroupID: groupId,
		RoleLevel: roleLevel,
	}
	updateInfo := db.GroupMember{
		RoleLevel: constant.GroupOwner,
	}
	if err := dbConn.Find(&groupMember).Update(updateInfo).Error; err != nil {
		return err
	}
	return nil
}

func GetGroupsCountNum(group db.Group) (int32, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return 0, err
	}
	dbConn.LogMode(true)
	var count int32
	if err := dbConn.Table("groups").Where(fmt.Sprintf(" name like '%%%s%%' ", group.GroupName)).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func GetGroupsById(groupId string) (db.Group, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	group := db.Group{
		GroupID: groupId,
	}
	if err != nil {
		return group, err
	}
	dbConn.LogMode(true)
	if err := dbConn.Table("groups").Find(&group).Error; err != nil {
		return group, err
	}
	return group, nil
}

func GetGroupMaster(groupId string) (db.GroupMember, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	groupMember := db.GroupMember{}
	if err != nil {
		return groupMember, err
	}
	dbConn.LogMode(true)
	if err := dbConn.Table("group_members").Where("role_level=? and group_id=?", constant.GroupOwner, groupId).Find(&groupMember).Error; err != nil {
		return groupMember, err
	}
	return groupMember, nil
}