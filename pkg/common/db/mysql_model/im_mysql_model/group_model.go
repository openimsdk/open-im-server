package im_mysql_model

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
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
	if groupInfo.GroupName == "" {
		groupInfo.GroupName = "Group Chat"
	}
	groupInfo.CreateTime = time.Now()

	if groupInfo.NotificationUpdateTime.Unix() < 0 {
		groupInfo.NotificationUpdateTime = utils.UnixSecondToTime(0)
	}
	err := db.DB.MysqlDB.DefaultGormDB().Table("groups").Create(groupInfo).Error
	if err != nil {
		return err
	}

	return nil
}

func GetGroupInfoByGroupID(groupId string) (*db.Group, error) {
	var groupInfo db.Group
	err := db.DB.MysqlDB.DefaultGormDB().Table("groups").Where("group_id=?", groupId).Take(&groupInfo).Error
	if err != nil {
		return nil, err
	}
	return &groupInfo, nil
}

func SetGroupInfo(groupInfo db.Group) error {
	return db.DB.MysqlDB.DefaultGormDB().Table("groups").Where("group_id=?", groupInfo.GroupID).Updates(&groupInfo).Error
}

func GetGroupsByName(groupName string, pageNumber, showNumber int32) ([]db.Group, error) {
	var groups []db.Group
	err := db.DB.MysqlDB.DefaultGormDB().Table("groups").Where(fmt.Sprintf(" name like '%%%s%%' ", groupName)).Limit(int(showNumber)).Offset(int(showNumber * (pageNumber - 1))).Find(&groups).Error
	return groups, err
}

func GetGroups(pageNumber, showNumber int) ([]db.Group, error) {
	var groups []db.Group
	if err := db.DB.MysqlDB.DefaultGormDB().Table("groups").Limit(showNumber).Offset(showNumber * (pageNumber - 1)).Find(&groups).Error; err != nil {
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
	var group db.Group
	var groupMembers []db.GroupMember
	if err := db.DB.MysqlDB.DefaultGormDB().Table("groups").Where("group_id=?", groupId).Delete(&group).Error; err != nil {
		return err
	}
	if err := db.DB.MysqlDB.DefaultGormDB().Table("group_members").Where("group_id=?", groupId).Delete(groupMembers).Error; err != nil {
		return err
	}
	return nil
}

func OperateGroupRole(userId, groupId string, roleLevel int32) (string, string, error) {
	groupMember := db.GroupMember{
		UserID:  userId,
		GroupID: groupId,
	}
	updateInfo := db.GroupMember{
		RoleLevel: roleLevel,
	}
	groupMaster := db.GroupMember{}
	var err error
	switch roleLevel {
	case constant.GroupOwner:
		err = db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
			result := db.DB.MysqlDB.DefaultGormDB().Table("group_members").Where("group_id = ? and role_level = ?", groupId, constant.GroupOwner).First(&groupMaster).Updates(&db.GroupMember{
				RoleLevel: constant.GroupOrdinaryUsers,
			})
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return errors.New(fmt.Sprintf("user %s not exist in group %s or already operate", userId, groupId))
			}

			result = db.DB.MysqlDB.DefaultGormDB().Table("group_members").First(&groupMember).Updates(updateInfo)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return errors.New(fmt.Sprintf("user %s not exist in group %s or already operate", userId, groupId))
			}
			return nil
		})

	case constant.GroupOrdinaryUsers:
		err = db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
			result := db.DB.MysqlDB.DefaultGormDB().Table("group_members").Where("group_id = ? and role_level = ?", groupId, constant.GroupOwner).First(&groupMaster)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return errors.New(fmt.Sprintf("user %s not exist in group %s or already operate", userId, groupId))
			}
			if groupMaster.UserID == userId {
				return errors.New(fmt.Sprintf("user %s is master of %s, cant set to ordinary user", userId, groupId))
			} else {
				result = db.DB.MysqlDB.DefaultGormDB().Table("group_members").Find(&groupMember).Updates(updateInfo)
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
	return "", "", err
}

func GetGroupsCountNum(group db.Group) (int32, error) {
	var count int64
	if err := db.DB.MysqlDB.DefaultGormDB().Table("groups").Where(fmt.Sprintf(" name like '%%%s%%' ", group.GroupName)).Count(&count).Error; err != nil {
		return 0, err
	}
	return int32(count), nil
}

func GetGroupById(groupId string) (db.Group, error) {
	group := db.Group{
		GroupID: groupId,
	}
	if err := db.DB.MysqlDB.DefaultGormDB().Table("groups").Find(&group).Error; err != nil {
		return group, err
	}
	return group, nil
}

func GetGroupMaster(groupId string) (db.GroupMember, error) {
	groupMember := db.GroupMember{}
	if err := db.DB.MysqlDB.DefaultGormDB().Table("group_members").Where("role_level=? and group_id=?", constant.GroupOwner, groupId).Find(&groupMember).Error; err != nil {
		return groupMember, err
	}
	return groupMember, nil
}

func UpdateGroupInfoDefaultZero(groupID string, args map[string]interface{}) error {
	return db.DB.MysqlDB.DefaultGormDB().Table("groups").Where("group_id = ? ", groupID).Updates(args).Error
}

func GetAllGroupIDList() ([]string, error) {
	var groupIDList []string
	err := db.DB.MysqlDB.DefaultGormDB().Table("groups").Pluck("group_id", &groupIDList).Error
	return groupIDList, err
}
