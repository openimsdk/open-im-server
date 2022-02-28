package im_mysql_model

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/utils"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
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
	err = dbConn.Table("groups").Where("group_id=?", groupId).Take(&groupInfo).Error
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
		Status:  groupStatus,
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

func OperateGroupRole(userId, groupId string, roleLevel int32) (string, string, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return "", "", err
	}
	dbConn.LogMode(true)
	groupMember := db.GroupMember{
		UserID:  userId,
		GroupID: groupId,
	}
	updateInfo := db.GroupMember{
		RoleLevel: roleLevel,
	}
	groupMaster := db.GroupMember{}
	switch roleLevel {
	case constant.GroupOwner:
		err = dbConn.Transaction(func(tx *gorm.DB) error {
			result := dbConn.Table("group_members").Where("group_id = ? and role_level = ?", groupId, constant.GroupOwner).First(&groupMaster).Update(&db.GroupMember{
				RoleLevel: constant.GroupOrdinaryUsers,
			})
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return errors.New(fmt.Sprintf("user %s not exist in group %s or already operate", userId, groupId))
			}

			result = dbConn.Table("group_members").First(&groupMember).Update(updateInfo)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return errors.New(fmt.Sprintf("user %s not exist in group %s or already operate", userId, groupId))
			}
			return nil
		})

	case constant.GroupOrdinaryUsers:
		err = dbConn.Transaction(func(tx *gorm.DB) error {
			result := dbConn.Table("group_members").Where("group_id = ? and role_level = ?", groupId, constant.GroupOwner).First(&groupMaster)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return errors.New(fmt.Sprintf("user %s not exist in group %s or already operate", userId, groupId))
			}
			if groupMaster.UserID == userId {
				return errors.New(fmt.Sprintf("user %s is master of %s, cant set to ordinary user", userId, groupId))
			} else {
				result = dbConn.Table("group_members").Find(&groupMember).Update(updateInfo)
				if result.Error != nil {
					return result.Error
				}
				if result.RowsAffected == 0 {
					return errors.New(fmt.Sprintf("user %s not exist in group %s or already operate", userId, groupId))
				}
			}
			return nil
		})
	}
	return "", "", nil
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

func GetGroupById(groupId string) (db.Group, error) {
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
