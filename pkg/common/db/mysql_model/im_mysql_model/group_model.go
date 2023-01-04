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

func TakeGroupInfoByGroupID(groupID string) (*db.Group, error) {
	var groupInfo db.Group
	err := db.DB.MysqlDB.DefaultGormDB().Table("groups").Where("group_id=?", groupID).Take(&groupInfo).Error
	return &groupInfo, err
}

func GetGroupInfoByGroupID(groupID string) (*db.Group, error) {
	var groupInfo db.Group
	err := db.DB.MysqlDB.DefaultGormDB().Table("groups").Where("group_id=?", groupID).Take(&groupInfo).Error
	return &groupInfo, err
}

func SetGroupInfo(groupInfo db.Group) error {
	return db.DB.MysqlDB.DefaultGormDB().Table("groups").Where("group_id=?", groupInfo.GroupID).Updates(&groupInfo).Error
}

type GroupWithNum struct {
	db.Group
	MemberCount int `gorm:"column:num"`
}

func GetGroupsByName(groupName string, pageNumber, showNumber int32) ([]GroupWithNum, int64, error) {
	var groups []GroupWithNum
	var count int64
	sql := db.DB.MysqlDB.DefaultGormDB().Table("groups").Select("groups.*, (select count(*) from group_members where group_members.group_id=groups.group_id) as num").
		Where(" name like ? and status != ?", fmt.Sprintf("%%%s%%", groupName), constant.GroupStatusDismissed)
	if err := sql.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	err := sql.Limit(int(showNumber)).Offset(int(showNumber * (pageNumber - 1))).Find(&groups).Error
	return groups, count, err
}

func GetGroups(pageNumber, showNumber int) ([]GroupWithNum, error) {
	var groups []GroupWithNum
	if err := db.DB.MysqlDB.DefaultGormDB().Table("groups").Select("groups.*, (select count(*) from group_members where group_members.group_id=groups.group_id) as num").
		Limit(showNumber).Offset(showNumber * (pageNumber - 1)).Find(&groups).Error; err != nil {
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

func GetGroupsCountNum(group db.Group) (int32, error) {
	var count int64
	if err := db.DB.MysqlDB.DefaultGormDB().Table("groups").Where(" name like ? ", fmt.Sprintf("%%%s%%", group.GroupName)).Count(&count).Error; err != nil {
		return 0, err
	}
	return int32(count), nil
}

func UpdateGroupInfoDefaultZero(groupID string, args map[string]interface{}) error {
	return db.DB.MysqlDB.DefaultGormDB().Table("groups").Where("group_id = ? ", groupID).Updates(args).Error
}

func GetGroupIDListByGroupType(groupType int) ([]string, error) {
	var groupIDList []string
	if err := db.DB.MysqlDB.DefaultGormDB().Table("groups").Where("group_type = ? ", groupType).Pluck("group_id", &groupIDList).Error; err != nil {
		return nil, err
	}
	return groupIDList, nil
}
