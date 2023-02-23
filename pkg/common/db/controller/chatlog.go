package controller

import (
	"OpenIM/pkg/common/db/relation"
	relationTb "OpenIM/pkg/common/db/table/relation"
	pbMsg "OpenIM/pkg/proto/msg"
	"gorm.io/gorm"
)

type ChatLogInterface interface {
	CreateChatLog(msg pbMsg.MsgDataToMQ) error
	GetChatLog(chatLog *relationTb.ChatLogModel, pageNumber, showNumber int32, contentTypeList []int32) (int64, []relationTb.ChatLogModel, error)
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

func (c *ChatLogController) GetChatLog(chatLog *relationTb.ChatLogModel, pageNumber, showNumber int32, contentTypeList []int32) (int64, []relationTb.ChatLogModel, error) {
	return c.database.GetChatLog(chatLog, pageNumber, showNumber, contentTypeList)
}

type ChatLogDataBaseInterface interface {
	CreateChatLog(msg pbMsg.MsgDataToMQ) error
	GetChatLog(chatLog *relationTb.ChatLogModel, pageNumber, showNumber int32, contentTypeList []int32) (int64, []relationTb.ChatLogModel, error)
}

type ChatLogDataBase struct {
	chatLogDB relationTb.ChatLogModelInterface
}

func NewChatLogDataBase(db *gorm.DB) ChatLogDataBaseInterface {
	return &ChatLogDataBase{chatLogDB: relation.NewChatLog(db)}
}

func (c *ChatLogDataBase) CreateChatLog(msg pbMsg.MsgDataToMQ) error {
	return c.chatLogDB.Create(msg)
}

func (c *ChatLogDataBase) GetChatLog(chatLog *relationTb.ChatLogModel, pageNumber, showNumber int32, contentTypeList []int32) (int64, []relationTb.ChatLogModel, error) {
	return c.chatLogDB.GetChatLog(chatLog, pageNumber, showNumber, contentTypeList)
}
