package model

import "time"

// UserGlobalBlack 全局黑名单记录，被加入黑名单的用户无法登录
type UserGlobalBlack struct {
	UserID     string    `bson:"user_id"`
	Nickname   string    `bson:"nickname"`
	OperatorID string    `bson:"operator_id"`
	Reason     string    `bson:"reason"`
	CreateTime time.Time `bson:"create_time"`
}
