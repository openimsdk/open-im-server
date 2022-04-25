package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
)

func SetConversation(conversation db.Conversation) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	dbConn.LogMode(false)
	newConversation := conversation
	if dbConn.Model(&db.Conversation{}).Find(&newConversation).RowsAffected == 0 {
		log.NewDebug("", utils.GetSelfFuncName(), "conversation", conversation, "not exist in db, create")
		return dbConn.Model(&db.Conversation{}).Create(conversation).Error
		// if exist, then update record
	} else {
		log.NewDebug("", utils.GetSelfFuncName(), "conversation", conversation, "exist in db, update")
		//force update
		return dbConn.Model(conversation).Where("owner_user_id = ? and conversation_id = ?", conversation.OwnerUserID, conversation.ConversationID).
			Update(map[string]interface{}{"recv_msg_opt": conversation.RecvMsgOpt, "is_pinned": conversation.IsPinned, "is_private_chat": conversation.IsPrivateChat,
				"group_at_type": conversation.GroupAtType, "is_not_in_group": conversation.IsNotInGroup}).Error
	}
}
func SetOneConversation(conversation db.Conversation) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	dbConn.LogMode(false)
	return dbConn.Model(&db.Conversation{}).Create(conversation).Error

}

func PeerUserSetConversation(conversation db.Conversation) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	dbConn.LogMode(false)
	newConversation := conversation
	if dbConn.Model(&db.Conversation{}).Find(&newConversation).RowsAffected == 0 {
		log.NewDebug("", utils.GetSelfFuncName(), "conversation", conversation, "not exist in db, create")
		return dbConn.Model(&db.Conversation{}).Create(conversation).Error
		// if exist, then update record
	}
	log.NewDebug("", utils.GetSelfFuncName(), "conversation", conversation, "exist in db, update")
	//force update
	return dbConn.Model(conversation).Where("owner_user_id = ? and conversation_id = ?", conversation.OwnerUserID, conversation.ConversationID).
		Update(map[string]interface{}{"is_private_chat": conversation.IsPrivateChat}).Error

}

func SetRecvMsgOpt(conversation db.Conversation) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	dbConn.LogMode(false)
	newConversation := conversation
	if dbConn.Model(&db.Conversation{}).Find(&newConversation).RowsAffected == 0 {
		log.NewDebug("", utils.GetSelfFuncName(), "conversation", conversation, "not exist in db, create")
		return dbConn.Model(&db.Conversation{}).Create(conversation).Error
		// if exist, then update record
	} else {
		log.NewDebug("", utils.GetSelfFuncName(), "conversation", conversation, "exist in db, update")
		//force update
		return dbConn.Model(conversation).Where("owner_user_id = ? and conversation_id = ?", conversation.OwnerUserID, conversation.ConversationID).
			Update(map[string]interface{}{"recv_msg_opt": conversation.RecvMsgOpt}).Error
	}
}

func GetUserAllConversations(ownerUserID string) ([]db.Conversation, error) {
	var conversations []db.Conversation
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return conversations, err
	}
	dbConn.LogMode(false)
	err = dbConn.Model(&db.Conversation{}).Where("owner_user_id=?", ownerUserID).Find(&conversations).Error
	return conversations, err
}
func GetMultipleUserConversationByConversationID(ownerUserIDList []string, conversationID string) ([]db.Conversation, error) {
	var conversations []db.Conversation
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return conversations, err
	}
	dbConn.LogMode(false)
	err = dbConn.Model(&db.Conversation{}).Where("owner_user_id IN ? and  conversation_id=?", ownerUserIDList, conversationID).Find(&conversations).Error
	return conversations, err
}
func GetExistConversationUserIDList(ownerUserIDList []string, conversationID string) ([]string, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var resultArr []string
	err = dbConn.Table("conversations").Where(" owner_user_id IN (?) and conversation_id=?", ownerUserIDList, conversationID).Pluck("owner_user_id", &resultArr).Error
	if err != nil {
		return nil, err
	}
	return resultArr, nil
}

func GetConversation(OwnerUserID, conversationID string) (db.Conversation, error) {
	var conversation db.Conversation
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return conversation, err
	}
	err = dbConn.Table("conversations").Where("owner_user_id=? and conversation_id=?", OwnerUserID, conversationID).Take(&conversation).Error
	return conversation, err
}

func GetConversations(OwnerUserID string, conversationIDs []string) ([]db.Conversation, error) {
	var conversations []db.Conversation
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return conversations, err
	}
	err = dbConn.Model(&db.Conversation{}).Where("conversation_id IN (?) and  owner_user_id=?", conversationIDs, OwnerUserID).Find(&conversations).Error
	return conversations, err
}
func GetConversationsByConversationIDMultipleOwner(OwnerUserIDList []string, conversationID string) ([]db.Conversation, error) {
	var conversations []db.Conversation
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return conversations, err
	}
	err = dbConn.Model(&db.Conversation{}).Where("owner_user_id IN (?) and  conversation_id=?", OwnerUserIDList, conversationID).Find(&conversations).Error
	return conversations, err
}
func UpdateColumnsConversations(ownerUserIDList []string, conversationID string, args map[string]interface{}) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	dbConn.LogMode(false)
	return dbConn.Debug().Model(&db.Conversation{}).Where("owner_user_id IN (?) and  conversation_id=?", ownerUserIDList, conversationID).Updates(args).Error

}
