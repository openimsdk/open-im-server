package controller

import (
	relationTb "OpenIM/pkg/common/db/table/relation"
	pbMsg "OpenIM/pkg/proto/msg"
)

type ChatLogDatabase interface {
	CreateChatLog(msg pbMsg.MsgDataToMQ) error
	GetChatLog(chatLog *relationTb.ChatLogModel, pageNumber, showNumber int32, contentTypeList []int32) (int64, []relationTb.ChatLogModel, error)
}

func NewChatLogDatabase(chatLogModelInterface relationTb.ChatLogModelInterface) ChatLogDatabase {
	return &ChatLogDataBase{chatLogModel: chatLogModelInterface}
}

type ChatLogDataBase struct {
	chatLogModel relationTb.ChatLogModelInterface
}

func (c *ChatLogDataBase) CreateChatLog(msg pbMsg.MsgDataToMQ) error {
	return c.chatLogModel.Create(msg)
}

func (c *ChatLogDataBase) GetChatLog(chatLog *relationTb.ChatLogModel, pageNumber, showNumber int32, contentTypeList []int32) (int64, []relationTb.ChatLogModel, error) {
	return c.chatLogModel.GetChatLog(chatLog, pageNumber, showNumber, contentTypeList)
}
