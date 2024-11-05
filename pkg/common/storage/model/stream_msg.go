package model

import (
	"time"
)

const (
	StreamMsgStatusWait = 0
	StreamMsgStatusDone = 1
	StreamMsgStatusFail = 2
)

type StreamMsg struct {
	ClientMsgID    string    `bson:"client_msg_id"`
	ConversationID string    `bson:"conversation_id"`
	UserID         string    `bson:"user_id"`
	Packets        []string  `bson:"packets"`
	End            bool      `bson:"end"`
	CreateTime     time.Time `bson:"create_time"`
	DeadlineTime   time.Time `bson:"deadline_time"`
}
