package model

import (
	"time"
)

type FriendRequest struct {
	FromUserID    string    `bson:"from_user_id"`
	ToUserID      string    `bson:"to_user_id"`
	HandleResult  int32     `bson:"handle_result"`
	ReqMsg        string    `bson:"req_msg"`
	CreateTime    time.Time `bson:"create_time"`
	HandlerUserID string    `bson:"handler_user_id"`
	HandleMsg     string    `bson:"handle_msg"`
	HandleTime    time.Time `bson:"handle_time"`
	Ex            string    `bson:"ex"`
}
