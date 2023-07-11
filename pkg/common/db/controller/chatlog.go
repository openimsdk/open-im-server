package controller

import (
	"context"
	relationTb "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	pbMsg "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
)

type ChatLogDatabase interface {
	CreateChatLog(msg *pbMsg.MsgDataToMQ) error
	GetChatLog(ctx context.Context, chatLog *relationTb.ChatLogModel, pageNumber, showNumber int32, contentTypes []int32) (int64, []relationTb.ChatLogModel, error)
}

func NewChatLogDatabase(chatLogModelInterface relationTb.ChatLogModelInterface) ChatLogDatabase {
	return &chatLogDatabase{chatLogModel: chatLogModelInterface}
}

type chatLogDatabase struct {
	chatLogModel relationTb.ChatLogModelInterface
}

func (c *chatLogDatabase) CreateChatLog(msg *pbMsg.MsgDataToMQ) error {
	return c.chatLogModel.Create(msg)
}

func (c *chatLogDatabase) GetChatLog(ctx context.Context, chatLog *relationTb.ChatLogModel, pageNumber, showNumber int32, contentTypes []int32) (int64, []relationTb.ChatLogModel, error) {
	return c.chatLogModel.GetChatLog(ctx, chatLog, pageNumber, showNumber, contentTypes)
}
