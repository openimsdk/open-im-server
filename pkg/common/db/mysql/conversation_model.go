package mysql

import (
	"gorm.io/gorm"
)

var ConversationDB *gorm.DB

type Conversation struct {
	OwnerUserID           string `gorm:"column:owner_user_id;primary_key;type:char(128)" json:"OwnerUserID"`
	ConversationID        string `gorm:"column:conversation_id;primary_key;type:char(128)" json:"conversationID"`
	ConversationType      int32  `gorm:"column:conversation_type" json:"conversationType"`
	UserID                string `gorm:"column:user_id;type:char(64)" json:"userID"`
	GroupID               string `gorm:"column:group_id;type:char(128)" json:"groupID"`
	RecvMsgOpt            int32  `gorm:"column:recv_msg_opt" json:"recvMsgOpt"`
	UnreadCount           int32  `gorm:"column:unread_count" json:"unreadCount"`
	DraftTextTime         int64  `gorm:"column:draft_text_time" json:"draftTextTime"`
	IsPinned              bool   `gorm:"column:is_pinned" json:"isPinned"`
	IsPrivateChat         bool   `gorm:"column:is_private_chat" json:"isPrivateChat"`
	BurnDuration          int32  `gorm:"column:burn_duration;default:30" json:"burnDuration"`
	GroupAtType           int32  `gorm:"column:group_at_type" json:"groupAtType"`
	IsNotInGroup          bool   `gorm:"column:is_not_in_group" json:"isNotInGroup"`
	UpdateUnreadCountTime int64  `gorm:"column:update_unread_count_time" json:"updateUnreadCountTime"`
	AttachedInfo          string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
	Ex                    string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
}

func (Conversation) TableName() string {
	return "conversations"
}

func SetConversation(conversation Conversation) (bool, error) {
	var isUpdate bool
	newConversation := conversation
	if ConversationDB.Model(&Conversation{}).Find(&newConversation).RowsAffected == 0 {
		return isUpdate, ConversationDB.Model(&Conversation{}).Create(&conversation).Error
		// if exist, then update record
	} else {
		//force update
		isUpdate = true
		return isUpdate, ConversationDB.Model(conversation).Where("owner_user_id = ? and conversation_id = ?", conversation.OwnerUserID, conversation.ConversationID).
			Updates(map[string]interface{}{"recv_msg_opt": conversation.RecvMsgOpt, "is_pinned": conversation.IsPinned, "is_private_chat": conversation.IsPrivateChat,
				"group_at_type": conversation.GroupAtType, "is_not_in_group": conversation.IsNotInGroup}).Error
	}
}
func SetOneConversation(conversation Conversation) error {
	return ConversationDB.Model(&Conversation{}).Create(&conversation).Error

}

func PeerUserSetConversation(conversation Conversation) error {
	newConversation := conversation
	if ConversationDB.Model(&Conversation{}).Find(&newConversation).RowsAffected == 0 {
		return ConversationDB.Model(&Conversation{}).Create(&conversation).Error
		// if exist, then update record
	}
	//force update
	return ConversationDB.Model(conversation).Where("owner_user_id = ? and conversation_id = ?", conversation.OwnerUserID, conversation.ConversationID).
		Updates(map[string]interface{}{"is_private_chat": conversation.IsPrivateChat}).Error

}

func SetRecvMsgOpt(conversation Conversation) (bool, error) {
	var isUpdate bool
	newConversation := conversation
	if ConversationDB.Model(&Conversation{}).Find(&newConversation).RowsAffected == 0 {
		return isUpdate, ConversationDB.Model(&Conversation{}).Create(&conversation).Error
		// if exist, then update record
	} else {
		//force update
		isUpdate = true
		return isUpdate, ConversationDB.Model(conversation).Where("owner_user_id = ? and conversation_id = ?", conversation.OwnerUserID, conversation.ConversationID).
			Updates(map[string]interface{}{"recv_msg_opt": conversation.RecvMsgOpt}).Error
	}
}

func GetUserAllConversations(ownerUserID string) ([]Conversation, error) {
	var conversations []Conversation
	err := ConversationDB.Where("owner_user_id=?", ownerUserID).Find(&conversations).Error
	return conversations, err
}
func GetMultipleUserConversationByConversationID(ownerUserIDList []string, conversationID string) ([]Conversation, error) {
	var conversations []Conversation
	err := ConversationDB.Where("owner_user_id IN ? and  conversation_id=?", ownerUserIDList, conversationID).Find(&conversations).Error
	return conversations, err
}
func GetExistConversationUserIDList(ownerUserIDList []string, conversationID string) ([]string, error) {
	var resultArr []string
	err := ConversationDB.Table("conversations").Where(" owner_user_id IN (?) and conversation_id=?", ownerUserIDList, conversationID).Pluck("owner_user_id", &resultArr).Error
	if err != nil {
		return nil, err
	}
	return resultArr, nil
}

func GetConversation(OwnerUserID, conversationID string) (Conversation, error) {
	var conversation Conversation
	err := ConversationDB.Table("conversations").Where("owner_user_id=? and conversation_id=?", OwnerUserID, conversationID).Take(&conversation).Error
	return conversation, err
}

func GetConversations(OwnerUserID string, conversationIDs []string) ([]Conversation, error) {
	var conversations []Conversation
	err := ConversationDB.Model(&Conversation{}).Where("conversation_id IN (?) and  owner_user_id=?", conversationIDs, OwnerUserID).Find(&conversations).Error
	return conversations, err
}

func GetConversationsByConversationIDMultipleOwner(OwnerUserIDList []string, conversationID string) ([]Conversation, error) {
	var conversations []Conversation
	err := ConversationDB.Model(&Conversation{}).Where("owner_user_id IN (?) and  conversation_id=?", OwnerUserIDList, conversationID).Find(&conversations).Error
	return conversations, err
}

func UpdateColumnsConversations(ownerUserIDList []string, conversationID string, args map[string]interface{}) error {
	return ConversationDB.Model(&Conversation{}).Where("owner_user_id IN (?) and  conversation_id=?", ownerUserIDList, conversationID).Updates(args).Error

}

func GetConversationIDListByUserID(userID string) ([]string, error) {
	var IDList []string
	err := ConversationDB.Model(&Conversation{}).Where("owner_user_id=?", userID).Pluck("conversation_id", &IDList).Error
	return IDList, err
}
