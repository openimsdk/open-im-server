package model

import (
	"time"
)

type GroupRequest struct {
	UserID        string    `bson:"user_id"`
	GroupID       string    `bson:"group_id"`
	HandleResult  int32     `bson:"handle_result"`
	ReqMsg        string    `bson:"req_msg"`
	HandledMsg    string    `bson:"handled_msg"`
	ReqTime       time.Time `bson:"req_time"`
	HandleUserID  string    `bson:"handle_user_id"`
	HandledTime   time.Time `bson:"handled_time"`
	JoinSource    int32     `bson:"join_source"`
	InviterUserID string    `bson:"inviter_user_id"`
	Ex            string    `bson:"ex"`
}
