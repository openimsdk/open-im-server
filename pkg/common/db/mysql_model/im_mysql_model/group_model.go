package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/proto/group"
	"errors"
	"time"
)

func InsertIntoGroup(groupId, name, introduction, notification, faceUrl, ex string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	//Default group name
	if name == "" {
		name = "groupChat"
	}
	toInsertInfo := Group{GroupId: groupId, Name: name, Introduction: introduction, Notification: notification, FaceUrl: faceUrl, CreateTime: time.Now(), Ex: ex}
	err = dbConn.Table("group").Create(toInsertInfo).Error
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

func SetGroupInfo(groupId, groupName, introduction, notification, groupFaceUrl, ex string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	dbConn.LogMode(true)
	if err != nil {
		return err
	}
	if groupName != "" {
		if err = dbConn.Exec("update `group` set name=? where group_id=?", groupName, groupId).Error; err != nil {
			return err
		}
	}
	if introduction != "" {
		if err = dbConn.Exec("update `group` set introduction=? where group_id=?", introduction, groupId).Error; err != nil {
			return err
		}
	}
	if notification != "" {
		if err = dbConn.Exec("update `group` set notification=? where group_id=?", notification, groupId).Error; err != nil {
			return err
		}
	}
	if groupFaceUrl != "" {
		if err = dbConn.Exec("update `group` set face_url=? where group_id=?", groupFaceUrl, groupId).Error; err != nil {
			return err
		}
	}
	if ex != "" {
		if err = dbConn.Exec("update `group` set ex=? where group_id=?", ex, groupId).Error; err != nil {
			return err
		}
	}
	return nil
}

func GetGroupApplicationList(uid string) (*group.GetGroupApplicationListResp, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}

	var gID string
	var gIDs []string
	rows, err := dbConn.Raw("select group_id from `group_member` where uid = ? and administrator_level > 0", uid).Rows()
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		rows.Scan(&gID)
		gIDs = append(gIDs, gID)
	}

	if len(gIDs) == 0 {
		return &group.GetGroupApplicationListResp{}, nil
	}

	sql := "select id, group_id, from_user_id, to_user_id, flag, req_msg, handled_msg, create_time, " +
		"from_user_nickname, to_user_nickname, from_user_face_url, to_user_face_url, handled_user  from `group_request` where  group_id in ( "
	for i := 0; i < len(gIDs); i++ {
		if i == len(gIDs)-1 {
			sql = sql + "\"" + gIDs[i] + "\"" + " )"
		} else {
			sql = sql + "\"" + gIDs[i] + "\"" + ", "
		}
	}

	var groupRequest GroupRequest
	var groupRequests []GroupRequest
	log.Info("", "", sql)
	rows, err = dbConn.Raw(sql).Rows()
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		rows.Scan(&groupRequest.ID, &groupRequest.GroupID, &groupRequest.FromUserID, &groupRequest.ToUserID, &groupRequest.Flag, &groupRequest.ReqMsg,
			&groupRequest.HandledMsg, &groupRequest.CreateTime, &groupRequest.FromUserNickname, &groupRequest.ToUserNickname,
			&groupRequest.FromUserFaceUrl, &groupRequest.ToUserFaceUrl, &groupRequest.HandledUser)
		groupRequests = append(groupRequests, groupRequest)
	}

	reply := &group.GetGroupApplicationListResp{}
	reply.Data = &group.GetGroupApplicationListData{}
	reply.Data.Count = int32(len(groupRequests))
	for i := 0; i < int(reply.Data.Count); i++ {
		addUser := group.GetGroupApplicationList_Data_User{
			ID:               groupRequests[i].ID,
			GroupID:          groupRequests[i].GroupID,
			FromUserID:       groupRequests[i].FromUserID,
			FromUserNickname: groupRequests[i].FromUserNickname,
			FromUserFaceUrl:  groupRequests[i].FromUserFaceUrl,
			ToUserID:         groupRequests[i].ToUserID,
			AddTime:          groupRequests[i].CreateTime.Unix(),
			RequestMsg:       groupRequests[i].ReqMsg,
			HandledMsg:       groupRequests[i].HandledMsg,
			Flag:             groupRequests[i].Flag,
			ToUserNickname:   groupRequests[i].ToUserNickname,
			ToUserFaceUrl:    groupRequests[i].ToUserFaceUrl,
			HandledUser:      groupRequests[i].HandledUser,
			Type:             0,
			HandleStatus:     0,
			HandleResult:     0,
		}

		if addUser.ToUserID != "0" {
			addUser.Type = 1
		}

		if len(groupRequests[i].HandledUser) > 0 {
			if groupRequests[i].HandledUser == uid {
				addUser.HandleStatus = 2
			} else {
				addUser.HandleStatus = 1
			}
		}

		if groupRequests[i].Flag == 1 {
			addUser.HandleResult = 1
		}

		reply.Data.User = append(reply.Data.User, &addUser)
	}
	return reply, nil
}

func TransferGroupOwner(pb *group.TransferGroupOwnerReq) (*group.TransferGroupOwnerResp, error) {
	oldOwner, err := FindGroupMemberInfoByGroupIdAndUserId(pb.GroupID, pb.OldOwner)
	if err != nil {
		return nil, err
	}
	newOwner, err := FindGroupMemberInfoByGroupIdAndUserId(pb.GroupID, pb.NewOwner)
	if err != nil {
		return nil, err
	}

	if oldOwner.Uid == newOwner.Uid {
		return nil, errors.New("the self")
	}

	if err = UpdateTheUserAdministratorLevel(pb.GroupID, pb.OldOwner, 0); err != nil {
		return nil, err
	}

	if err = UpdateTheUserAdministratorLevel(pb.GroupID, pb.NewOwner, 1); err != nil {
		UpdateTheUserAdministratorLevel(pb.GroupID, pb.OldOwner, 1)
		return nil, err
	}

	return &group.TransferGroupOwnerResp{}, nil
}

func GroupApplicationResponse(pb *group.GroupApplicationResponseReq) (*group.GroupApplicationResponseResp, error) {

	ownerUser, err := FindGroupMemberInfoByGroupIdAndUserId(pb.GroupID, pb.OwnerID)
	if err != nil {
		log.ErrorByKv("FindGroupMemberInfoByGroupIdAndUserId failed", pb.OperationID, "groupId", pb.GroupID, "ownerID", pb.OwnerID)
		return nil, err
	}
	if ownerUser.AdministratorLevel <= 0 {
		return nil, errors.New("insufficient permissions")
	}

	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var groupRequest GroupRequest
	err = dbConn.Raw("select * from `group_request` where handled_user = ? and group_id = ? and from_user_id = ? and to_user_id = ?",
		"", pb.GroupID, pb.FromUserID, pb.ToUserID).Scan(&groupRequest).Error
	if err != nil {
		log.ErrorByKv("find group_request info failed", pb.OperationID, "groupId", pb.GroupID, "fromUserId", pb.FromUserID, "toUserId", pb.OwnerID)
		return nil, err
	}

	if groupRequest.Flag != 0 {
		return nil, errors.New("application has already handle")
	}

	var saveFlag int
	if pb.HandleResult == 0 {
		saveFlag = -1
	} else if pb.HandleResult == 1 {
		saveFlag = 1
	} else {
		return nil, errors.New("parma HandleResult error")
	}
	err = dbConn.Exec("update `group_request` set flag = ?, handled_msg = ?, handled_user = ? where group_id = ? and from_user_id = ? and to_user_id = ?",
		saveFlag, pb.HandledMsg, pb.OwnerID, groupRequest.GroupID, groupRequest.FromUserID, groupRequest.ToUserID).Error
	if err != nil {
		log.ErrorByKv("update group request failed", pb.OperationID, "groupID", pb.GroupID, "flag", saveFlag, "ownerId", pb.OwnerID, "fromUserId", pb.FromUserID, "toUserID", pb.ToUserID)
		return nil, err
	}

	if saveFlag == 1 {
		if groupRequest.ToUserID == "0" {
			err = InsertIntoGroupMember(pb.GroupID, pb.FromUserID, groupRequest.FromUserNickname, groupRequest.FromUserFaceUrl, 0)
			if err != nil {
				log.ErrorByKv("InsertIntoGroupMember failed", pb.OperationID, "groupID", pb.GroupID, "fromUserId", pb.FromUserID)
				return nil, err
			}
		} else {
			err = InsertIntoGroupMember(pb.GroupID, pb.ToUserID, groupRequest.ToUserNickname, groupRequest.ToUserFaceUrl, 0)
			if err != nil {
				log.ErrorByKv("InsertIntoGroupMember failed", pb.OperationID, "groupID", pb.GroupID, "fromUserId", pb.FromUserID)
				return nil, err
			}
		}
	}

	//if err != nil {
	//	err = dbConn.Raw("select * from `group_request` where handled_user = ? and group_id = ? and to_user_id = ? and from_user_id = ?", "", pb.GroupID, "0", pb.UID).Scan(&groupRequest).Error
	//	if err != nil {
	//		return nil, err
	//	}
	//	if pb.Flag == 1 {
	//		err = dbConn.Exec("update `group_request` set flag = ?, handled_msg = ?, handled_user = ? where group_id = ? and to_user_id = ? and from_user_id = ?",
	//			pb.Flag, pb.RespMsg, pb.OwnerID, pb.GroupID, "0", pb.UID).Error
	//		if err != nil {
	//			return nil, err
	//		}
	//
	//		// add to group member
	//		err = InsertIntoGroupMember(pb.GroupID, pb.UID, groupRequest.FromUserNickname, groupRequest.FromUserFaceUrl, 0)
	//		if err != nil {
	//			return nil, err
	//		}
	//	} else if pb.Flag == -1 {
	//		err = dbConn.Exec("update `group_request` set flag = ?, handled_msg = ?, handled_user = ? where group_id = ? and to_user_id = ? and from_user_id = ?",
	//			pb.Flag, pb.RespMsg, pb.OwnerID, pb.GroupID, "0", pb.UID).Error
	//		if err != nil {
	//			return nil, err
	//		}
	//	} else {
	//		return nil, errors.New("flag error")
	//	}
	//} else {
	//	if pb.Flag == 1 {
	//		err = dbConn.Exec("update `group_request` set flag = ?, handled_msg = ?, handled_user = ? where group_id = ? and to_user_id = ?",
	//			pb.Flag, pb.RespMsg, pb.OwnerID, pb.GroupID, pb.UID).Error
	//		if err != nil {
	//			return nil, err
	//		}
	//
	//		// add to group member
	//		err = InsertIntoGroupMember(pb.GroupID, pb.UID, groupRequest.ToUserNickname, groupRequest.ToUserFaceUrl, 0)
	//		if err != nil {
	//			return nil, err
	//		}
	//	} else if pb.Flag == -1 {
	//		err = dbConn.Exec("update `group_request` set flag = ?, handled_msg = ?, handled_user = ? where group_id = ? and to_user_id = ?",
	//			pb.Flag, pb.RespMsg, pb.OwnerID, pb.GroupID, pb.UID).Error
	//		if err != nil {
	//			return nil, err
	//		}
	//	} else {
	//		return nil, errors.New("flag error")
	//	}
	//}

	return &group.GroupApplicationResponseResp{}, nil
}
