package im_mysql_model

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/utils"
	"fmt"

	"time"
)

func InsertIntoGroup(groupInfo Group) error {
	if groupInfo.GroupName == "" {
		groupInfo.GroupName = "Group Chat"
	}
	groupInfo.CreateTime = time.Now()

	if groupInfo.NotificationUpdateTime.Unix() < 0 {
		groupInfo.NotificationUpdateTime = utils.UnixSecondToTime(0)
	}
	err := GroupDB.Create(groupInfo).Error
	if err != nil {
		return err
	}
	return nil
}

func TakeGroupInfoByGroupID(groupID string) (*Group, error) {
	var groupInfo Group
	err := GroupDB.Where("group_id=?", groupID).Take(&groupInfo).Error
	return &groupInfo, err
}

func GetGroupInfoByGroupID(groupID string) (*Group, error) {
	var groupInfo Group
	err := GroupDB.Where("group_id=?", groupID).Take(&groupInfo).Error
	return &groupInfo, err
}

func SetGroupInfo(groupInfo Group) error {
	return GroupDB.Where("group_id=?", groupInfo.GroupID).Updates(&groupInfo).Error
}

type GroupWithNum struct {
	Group
	MemberCount int `gorm:"column:num"`
}

func GetGroupsByName(groupName string, pageNumber, showNumber int32) ([]GroupWithNum, int64, error) {
	var groups []GroupWithNum
	var count int64
	sql := GroupDB.Select("groups.*, (select count(*) from group_members where group_members.group_id=groups.group_id) as num").
		Where(" name like ? and status != ?", fmt.Sprintf("%%%s%%", groupName), constant.GroupStatusDismissed)
	if err := sql.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	err := sql.Limit(int(showNumber)).Offset(int(showNumber * (pageNumber - 1))).Find(&groups).Error
	return groups, count, err
}

func GetGroups(pageNumber, showNumber int) ([]GroupWithNum, error) {
	var groups []GroupWithNum
	if err := GroupDB.Select("groups.*, (select count(*) from group_members where group_members.group_id=groups.group_id) as num").
		Limit(showNumber).Offset(showNumber * (pageNumber - 1)).Find(&groups).Error; err != nil {
		return groups, err
	}
	return groups, nil
}

func OperateGroupStatus(groupId string, groupStatus int32) error {
	group := Group{
		GroupID: groupId,
		Status:  groupStatus,
	}
	if err := SetGroupInfo(group); err != nil {
		return err
	}
	return nil
}

func GetGroupsCountNum(group Group) (int32, error) {
	var count int64
	if err := GroupDB.Where(" name like ? ", fmt.Sprintf("%%%s%%", group.GroupName)).Count(&count).Error; err != nil {
		return 0, err
	}
	return int32(count), nil
}

func UpdateGroupInfoDefaultZero(groupID string, args map[string]interface{}) error {
	return GroupDB.Where("group_id = ? ", groupID).Updates(args).Error
}

func GetGroupIDListByGroupType(groupType int) ([]string, error) {
	var groupIDList []string
	if err := GroupDB.Where("group_type = ? ", groupType).Pluck("group_id", &groupIDList).Error; err != nil {
		return nil, err
	}
	return groupIDList, nil
}
