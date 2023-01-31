package controller

import (
	"Open_IM/pkg/common/db/relation"
	pbMsg "Open_IM/pkg/proto/msg"
	"gorm.io/gorm"
)

type ChatLogInterface interface {
	CreateChatLog(msg pbMsg.MsgDataToMQ) error
	GetChatLog(chatLog *relation.ChatLog, pageNumber, showNumber int32, contentTypeList []int32) (int64, []relation.ChatLog, error)
}

func NewChatLogController(db *gorm.DB) ChatLogInterface {
	return &ChatLogController{database: NewChatLogDataBase(db)}
}

type ChatLogController struct {
	database ChatLogDataBaseInterface
}

func (c *ChatLogController) CreateChatLog(msg pbMsg.MsgDataToMQ) error {
	return c.database.CreateChatLog(msg)
}

func (c *ChatLogController) GetChatLog(chatLog *relation.ChatLog, pageNumber, showNumber int32, contentTypeList []int32) (int64, []relation.ChatLog, error) {
	return c.database.GetChatLog(chatLog, pageNumber, showNumber, contentTypeList)
}

type ChatLogDataBaseInterface interface {
	CreateChatLog(msg pbMsg.MsgDataToMQ) error
	GetChatLog(chatLog *relation.ChatLog, pageNumber, showNumber int32, contentTypeList []int32) (int64, []relation.ChatLog, error)
}

type ChatLogDataBase struct {
	chatLogDB *relation.ChatLog
}

func NewChatLogDataBase(db *gorm.DB) ChatLogDataBaseInterface {
	return &ChatLogDataBase{chatLogDB: relation.NewChatLog(db)}
}

func (c *ChatLogDataBase) CreateChatLog(msg pbMsg.MsgDataToMQ) error {
	return c.chatLogDB.Create(msg)
}

func (c *ChatLogDataBase) GetChatLog(chatLog *relation.ChatLog, pageNumber, showNumber int32, contentTypeList []int32) (int64, []relation.ChatLog, error) {
	return c.chatLogDB.GetChatLog(chatLog, pageNumber, showNumber, contentTypeList)
}
