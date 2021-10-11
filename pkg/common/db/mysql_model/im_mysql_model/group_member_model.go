package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"time"
)

func InsertIntoGroupMember(groupId, uid, nickName, userGroupFaceUrl string, administratorLevel int32) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	toInsertInfo := GroupMember{GroupId: groupId, Uid: uid, NickName: nickName, AdministratorLevel: administratorLevel, JoinTime: time.Now(), UserGroupFaceUrl: userGroupFaceUrl}
	err = dbConn.Table("group_member").Create(toInsertInfo).Error
	if err != nil {
		return err
	}
	return nil
}

func FindGroupMemberListByUserId(uid string) ([]GroupMember, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var groupMemberList []GroupMember
	err = dbConn.Raw("select * from `group_member` where uid=?", uid).Find(&groupMemberList).Error
	if err != nil {
		return nil, err
	}
	return groupMemberList, nil
}

func FindGroupMemberListByGroupId(groupId string) ([]GroupMember, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var groupMemberList []GroupMember
	err = dbConn.Raw("select * from `group_member` where group_id=?", groupId).Find(&groupMemberList).Error
	if err != nil {
		return nil, err
	}
	return groupMemberList, nil
}

func FindGroupMemberListByGroupIdAndFilterInfo(groupId string, filter int32) ([]GroupMember, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	dbConn.LogMode(true)
	if err != nil {
		return nil, err
	}
	var groupMemberList []GroupMember
	err = dbConn.Raw("select * from `group_member` where group_id=? and administrator_level=?", groupId, filter).Find(&groupMemberList).Error
	if err != nil {
		return nil, err
	}
	return groupMemberList, nil
}
func FindGroupMemberInfoByGroupIdAndUserId(groupId, uid string) (*GroupMember, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var groupMember GroupMember
	err = dbConn.Raw("select * from `group_member` where group_id=? and uid=? limit 1", groupId, uid).Scan(&groupMember).Error
	if err != nil {
		return nil, err
	}
	return &groupMember, nil
}

func DeleteGroupMemberByGroupIdAndUserId(groupId, uid string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	err = dbConn.Exec("delete  from `group_member` where group_id=? and uid=?", groupId, uid).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateOwnerGroupNickName(groupId, userId, groupNickName string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	err = dbConn.Exec("update `group_member` set nickname=? where group_id=? and uid=?", groupNickName, groupId, userId).Error
	if err != nil {
		return err
	}
	return nil
}

func SelectGroupList(groupID string) ([]string, error) {
	var groupUserID string
	var groupList []string
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return groupList, err
	}

	rows, err := dbConn.Model(&GroupMember{}).Where("group_id = ?", groupID).Select("user_id").Rows()
	if err != nil {
		return groupList, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&groupUserID)
		groupList = append(groupList, groupUserID)
	}
	return groupList, nil
}

func UpdateTheUserAdministratorLevel(groupId, uid string, administratorLevel int64) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	err = dbConn.Exec("update `group_member` set administrator_level=? where group_id=? and uid=?", administratorLevel, groupId, uid).Error
	if err != nil {
		return err
	}
	return nil
}

func GetOwnerManagerByGroupId(groupId string) ([]GroupMember, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var groupMemberList []GroupMember
	err = dbConn.Raw("select * from `group_member` where group_id=? and administrator_level > 0", groupId).Find(&groupMemberList).Error
	if err != nil {
		return nil, err
	}
	return groupMemberList, nil
}

func IsExistGroupMember(groupId, uid string) bool {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return false
	}
	var number int32
	err = dbConn.Raw("select count(*) from `group_member` where group_id = ? and uid = ?", groupId, uid).Count(&number).Error
	if err != nil {
		return false
	}

	if number != 1 {
		return false
	}
	return true
}

func RemoveGroupMember(groupId string, memberId string) error {
	return DeleteGroupMemberByGroupIdAndUserId(groupId, memberId)
}

func GetMemberInfoById(groupId string, memberId string) (*GroupMember, error) {
	return FindGroupMemberInfoByGroupIdAndUserId(groupId, memberId)
}

func GetGroupMemberByGroupId(groupId string, filter int32, begin int32, maxNumber int32) ([]GroupMember, error) {
	memberList, err := FindGroupMemberListByGroupId(groupId) //sorted by join time
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

func GetJoinedGroupIdListByMemberId(memberId string) ([]GroupMember, error) {
	return FindGroupMemberListByUserId(memberId)
}

func GetGroupMemberNumByGroupId(groupId string) int32 {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return 0
	}
	var number int32
	err = dbConn.Raw("select count(*) from `group_member` where group_id=? ", groupId).Count(&number).Error
	if err != nil {
		return 0
	}
	return number
}

func GetGroupOwnerByGroupId(groupId string) string {
	omList, err := GetOwnerManagerByGroupId(groupId)
	if err != nil {
		return ""
	}
	for _, v := range omList {
		if v.AdministratorLevel == 1 {
			return v.Uid
		}
	}
	return ""
}

func InsertGroupMember(groupId, userId, nickName, userFaceUrl string, role int32) error {
	return InsertIntoGroupMember(groupId, userId, nickName, userFaceUrl, role)
}
