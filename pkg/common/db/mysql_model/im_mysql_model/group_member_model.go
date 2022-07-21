package im_mysql_model

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/utils"
	"errors"
	"fmt"
	"time"
)

//type GroupMember struct {
//	GroupID            string    `gorm:"column:group_id;primaryKey;"`
//	UserID             string    `gorm:"column:user_id;primaryKey;"`
//	NickName           string    `gorm:"column:nickname"`
//	FaceUrl            string    `gorm:"user_group_face_url"`
//	RoleLevel int32     `gorm:"column:role_level"`
//	JoinTime           time.Time `gorm:"column:join_time"`
//	JoinSource int32 `gorm:"column:join_source"`
//	OperatorUserID  string `gorm:"column:operator_user_id"`
//	Ex string `gorm:"column:ex"`
//}

func InsertIntoGroupMember(toInsertInfo db.GroupMember) error {
	toInsertInfo.JoinTime = time.Now()
	if toInsertInfo.RoleLevel == 0 {
		toInsertInfo.RoleLevel = constant.GroupOrdinaryUsers
	}
	toInsertInfo.MuteEndTime = time.Unix(int64(time.Now().Second()), 0)
	err := db.DB.MysqlDB.DefaultGormDB().Table("group_members").Create(toInsertInfo).Error
	if err != nil {
		return err
	}
	return nil
}

func GetGroupMemberListByUserID(userID string) ([]db.GroupMember, error) {
	var groupMemberList []db.GroupMember
	err := db.DB.MysqlDB.DefaultGormDB().Table("group_members").Where("user_id=?", userID).Find(&groupMemberList).Error
	//err = dbConn.Table("group_members").Where("user_id=?", userID).Take(&groupMemberList).Error
	if err != nil {
		return nil, err
	}
	return groupMemberList, nil
}

func GetGroupMemberListByGroupID(groupID string) ([]db.GroupMember, error) {
	var groupMemberList []db.GroupMember
	err := db.DB.MysqlDB.DefaultGormDB().Table("group_members").Where("group_id=?", groupID).Find(&groupMemberList).Error
	if err != nil {
		return nil, err
	}
	return groupMemberList, nil
}

func GetGroupMemberIDListByGroupID(groupID string) ([]string, error) {
	var groupMemberIDList []string
	err := db.DB.MysqlDB.DefaultGormDB().Table("group_members").Where("group_id=?", groupID).Pluck("user_id", &groupMemberIDList).Error
	if err != nil {
		return nil, err
	}
	return groupMemberIDList, nil
}

func GetGroupMemberListByGroupIDAndRoleLevel(groupID string, roleLevel int32) ([]db.GroupMember, error) {
	var groupMemberList []db.GroupMember
	err := db.DB.MysqlDB.DefaultGormDB().Table("group_members").Where("group_id=? and role_level=?", groupID, roleLevel).Find(&groupMemberList).Error
	if err != nil {
		return nil, err
	}
	return groupMemberList, nil
}

func GetGroupMemberInfoByGroupIDAndUserID(groupID, userID string) (*db.GroupMember, error) {
	var groupMember db.GroupMember
	err := db.DB.MysqlDB.DefaultGormDB().Table("group_members").Where("group_id=? and user_id=? ", groupID, userID).Limit(1).Take(&groupMember).Error
	if err != nil {
		return nil, err
	}
	return &groupMember, nil
}

func DeleteGroupMemberByGroupIDAndUserID(groupID, userID string) error {
	return db.DB.MysqlDB.DefaultGormDB().Table("group_members").Where("group_id=? and user_id=? ", groupID, userID).Delete(db.GroupMember{}).Error
}

func DeleteGroupMemberByGroupID(groupID string) error {
	return db.DB.MysqlDB.DefaultGormDB().Table("group_members").Where("group_id=?  ", groupID).Delete(db.GroupMember{}).Error
}

func UpdateGroupMemberInfo(groupMemberInfo db.GroupMember) error {
	return db.DB.MysqlDB.DefaultGormDB().Table("group_members").Where("group_id=? and user_id=?", groupMemberInfo.GroupID, groupMemberInfo.UserID).Updates(&groupMemberInfo).Error
}

func UpdateGroupMemberInfoByMap(groupMemberInfo db.GroupMember, m map[string]interface{}) error {
	return db.DB.MysqlDB.DefaultGormDB().Table("group_members").Where("group_id=? and user_id=?", groupMemberInfo.GroupID, groupMemberInfo.UserID).Updates(m).Error
}

func GetOwnerManagerByGroupID(groupID string) ([]db.GroupMember, error) {
	var groupMemberList []db.GroupMember
	err := db.DB.MysqlDB.DefaultGormDB().Table("group_members").Where("group_id=? and role_level>?", groupID, constant.GroupOrdinaryUsers).Find(&groupMemberList).Error
	if err != nil {
		return nil, err
	}
	return groupMemberList, nil
}

func GetGroupMemberNumByGroupID(groupID string) (int64, error) {
	var number int64
	err := db.DB.MysqlDB.DefaultGormDB().Table("group_members").Where("group_id=?", groupID).Count(&number).Error
	if err != nil {
		return 0, utils.Wrap(err, "")
	}
	return number, nil
}

func GetGroupOwnerInfoByGroupID(groupID string) (*db.GroupMember, error) {
	omList, err := GetOwnerManagerByGroupID(groupID)
	if err != nil {
		return nil, err
	}
	for _, v := range omList {
		if v.RoleLevel == constant.GroupOwner {
			return &v, nil
		}
	}

	return nil, utils.Wrap(errors.New("no owner"), "")
}

func IsExistGroupMember(groupID, userID string) bool {
	var number int64
	err := db.DB.MysqlDB.DefaultGormDB().Table("group_members").Where("group_id = ? and user_id = ?", groupID, userID).Count(&number).Error
	if err != nil {
		return false
	}
	if number != 1 {
		return false
	}
	return true
}

func RemoveGroupMember(groupID string, UserID string) error {
	return DeleteGroupMemberByGroupIDAndUserID(groupID, UserID)
}

func GetMemberInfoByID(groupID string, userID string) (*db.GroupMember, error) {
	return GetGroupMemberInfoByGroupIDAndUserID(groupID, userID)
}

func GetGroupMemberByGroupID(groupID string, filter int32, begin int32, maxNumber int32) ([]db.GroupMember, error) {
	var memberList []db.GroupMember
	var err error
	if filter >= 0 {
		memberList, err = GetGroupMemberListByGroupIDAndRoleLevel(groupID, filter) //sorted by join time
	} else {
		memberList, err = GetGroupMemberListByGroupID(groupID)
	}

	if err != nil {
		return nil, err
	}
	if begin >= int32(len(memberList)) {
		return nil, nil
	}

	var end int32
	if begin+int32(maxNumber) < int32(len(memberList)) {
		end = begin + maxNumber
	} else {
		end = int32(len(memberList))
	}
	return memberList[begin:end], nil
}

func GetJoinedGroupIDListByUserID(userID string) ([]string, error) {
	memberList, err := GetGroupMemberListByUserID(userID)
	if err != nil {
		return nil, err
	}
	var groupIDList []string
	for _, v := range memberList {
		groupIDList = append(groupIDList, v.GroupID)
	}
	return groupIDList, nil
}

func IsGroupOwnerAdmin(groupID, UserID string) bool {
	groupMemberList, err := GetOwnerManagerByGroupID(groupID)
	if err != nil {
		return false
	}
	for _, v := range groupMemberList {
		if v.UserID == UserID && v.RoleLevel > constant.GroupOrdinaryUsers {
			return true
		}
	}
	return false
}

func GetGroupMembersByGroupIdCMS(groupId string, userName string, showNumber, pageNumber int32) ([]db.GroupMember, error) {
	var groupMembers []db.GroupMember
	err := db.DB.MysqlDB.DefaultGormDB().Table("group_members").Where("group_id=?", groupId).Where(fmt.Sprintf(" nickname like '%%%s%%' ", userName)).Limit(int(showNumber)).Offset(int(showNumber * (pageNumber - 1))).Find(&groupMembers).Error
	if err != nil {
		return nil, err
	}
	return groupMembers, nil
}

func GetGroupMembersCount(groupId, userName string) (int64, error) {
	var count int64
	if err := db.DB.MysqlDB.DefaultGormDB().Table("group_members").Where("group_id=?", groupId).Where(fmt.Sprintf(" nickname like '%%%s%%' ", userName)).Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}

func UpdateGroupMemberInfoDefaultZero(groupMemberInfo db.GroupMember, args map[string]interface{}) error {
	return db.DB.MysqlDB.DefaultGormDB().Model(groupMemberInfo).Updates(args).Error
}

//
//func SelectGroupList(groupID string) ([]string, error) {
//	var groupUserID string
//	var groupList []string
//	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
//	if err != nil {
//		return groupList, err
//	}
//
//	rows, err := dbConn.Model(&GroupMember{}).Where("group_id = ?", groupID).Select("user_id").Rows()
//	if err != nil {
//		return groupList, err
//	}
//	defer rows.Close()
//	for rows.Next() {
//		rows.Scan(&groupUserID)
//		groupList = append(groupList, groupUserID)
//	}
//	return groupList, nil
//}
