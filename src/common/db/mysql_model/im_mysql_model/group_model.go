package im_mysql_model

import "Open_IM/src/common/db"

func InsertIntoGroup(groupId, name, groupHeadUrl string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	//Default group name
	if name == "" {
		name = "groupChat"
	}
	err = dbConn.Exec("insert into `group`(group_id,name,head_url) values(?,?,?)", groupId, name, groupHeadUrl).Error
	if err != nil {
		return err
	}
	return nil
}

func FindGroupInfoByGroupId(groupId string) (*Group, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var groupInfo Group
	err = dbConn.Raw("select * from `group` where group_id=?", groupId).Scan(&groupInfo).Error
	if err != nil {
		return nil, err
	}
	return &groupInfo, nil
}

func UpdateGroupName(groupId, groupName string) (err error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	err = dbConn.Exec("update `group` set name=? where group_id=?", groupName, groupId).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateGroupBulletin(groupId, bulletinContent string) (err error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	err = dbConn.Exec("update `group` set bulletin=? where group_id=?", bulletinContent, groupId).Error
	if err != nil {
		return err
	}
	return nil
}
func UpdateGroupHeadImage(groupId, headImageUrl string) (err error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	err = dbConn.Exec("update `group` set head_url=? where group_id=?", headImageUrl, groupId).Error
	if err != nil {
		return err
	}
	return nil
}
