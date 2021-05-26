package im_mysql_model

import "Open_IM/src/common/db"

func InsertIntoGroupMember(groupId, userId string, isAdmin int64) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	err = dbConn.Exec("insert into `group_member`(group_id,user_id,is_admin) values(?,?,?)", groupId, userId, isAdmin).Error
	if err != nil {
		return err
	}
	return nil
}

func FindGroupMemberListByUserId(userId string) ([]GroupMember, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var groupMemberList []GroupMember
	err = dbConn.Raw("select * from `group_member` where user_id=?", userId).Find(&groupMemberList).Error
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

func FindGroupMemberInfoByGroupIdAndUserId(groupId, userId string) (*GroupMember, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var groupMember GroupMember
	err = dbConn.Raw("select * from `group_member` where group_id=? and user_id=? limit 1", groupId, userId).Scan(&groupMember).Error
	if err != nil {
		return nil, err
	}
	return &groupMember, nil
}

func DeleteGroupMemberByGroupIdAndUserId(groupId, userId string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	err = dbConn.Exec("delete  from `group_member` where group_id=? and user_id=?", groupId, userId).Error
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
	err = dbConn.Exec("update `group_member` set nickname=? where group_id=? and user_id=?", groupNickName, groupId, userId).Error
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
