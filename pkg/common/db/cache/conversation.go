package cache

import (
	"Open_IM/pkg/common/db/table"
	"Open_IM/pkg/utils"
	"encoding/json"
	"github.com/dtm-labs/rockscache"
	"time"
)

type DBFun func() (string, error)

type ConversationCache interface {
	GetUserConversationIDListFromCache(userID string, fn DBFun) ([]string, error)
	DelUserConversationIDListFromCache(userID string) error
	GetConversationFromCache(ownerUserID, conversationID string, fn DBFun) (*table.ConversationModel, error)
	GetConversationsFromCache(ownerUserID string, conversationIDList []string, fn DBFun) ([]*table.ConversationModel, error)
	GetUserAllConversationList(ownerUserID string, fn DBFun) ([]*table.ConversationModel, error)
	DelConversationFromCache(ownerUserID, conversationID string) error
}
type ConversationRedis struct {
	rcClient *rockscache.Client
}

func NewConversationRedis(rcClient *rockscache.Client) *ConversationRedis {
	return &ConversationRedis{rcClient: rcClient}
}

func (c *ConversationRedis) GetUserConversationIDListFromCache(userID string, fn DBFun) ([]string, error) {
	conversationIDListStr, err := c.rcClient.Fetch(conversationIDListCache+userID, time.Second*30*60, fn)
	var conversationIDList []string
	err = json.Unmarshal([]byte(conversationIDListStr), &conversationIDList)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	return conversationIDList, nil
}

func (c *ConversationRedis) DelUserConversationIDListFromCache(userID string) error {
	return utils.Wrap(c.rcClient.TagAsDeleted(conversationIDListCache+userID), "DelUserConversationIDListFromCache err")
}

func (c *ConversationRedis) GetConversationFromCache(ownerUserID, conversationID string, fn DBFun) (*table.ConversationModel, error) {
	conversationStr, err := c.rcClient.Fetch(conversationCache+ownerUserID+":"+conversationID, time.Second*30*60, fn)
	if err != nil {
		return nil, utils.Wrap(err, "Fetch failed")
	}
	conversation := table.ConversationModel{}
	err = json.Unmarshal([]byte(conversationStr), &conversation)
	if err != nil {
		return nil, utils.Wrap(err, "Unmarshal failed")
	}
	return &conversation, nil
}

func (c *ConversationRedis) GetConversationsFromCache(ownerUserID string, conversationIDList []string, fn DBFun) ([]*table.ConversationModel, error) {
	var conversationList []*table.ConversationModel
	for _, conversationID := range conversationIDList {
		conversation, err := c.GetConversationFromCache(ownerUserID, conversationID, fn)
		if err != nil {
			return nil, utils.Wrap(err, "GetConversationFromCache failed")
		}
		conversationList = append(conversationList, conversation)
	}
	return conversationList, nil
}

func (c *ConversationRedis) GetUserAllConversationList(ownerUserID string, fn DBFun) ([]*table.ConversationModel, error) {
	IDList, err := c.GetUserConversationIDListFromCache(ownerUserID, fn)
	if err != nil {
		return nil, err
	}
	var conversationList []*table.ConversationModel
	for _, conversationID := range IDList {
		conversation, err := c.GetConversationFromCache(ownerUserID, conversationID, fn)
		if err != nil {
			return nil, utils.Wrap(err, "GetConversationFromCache failed")
		}
		conversationList = append(conversationList, conversation)
	}
	return conversationList, nil
}

func (c *ConversationRedis) DelConversationFromCache(ownerUserID, conversationID string) error {
	return utils.Wrap(c.rcClient.TagAsDeleted(conversationCache+ownerUserID+":"+conversationID), "DelConversationFromCache err")
}
