package model

import "time"

// UserGlobalBlack 全局黑名单/冻结记录。
// Status: 1=冻结（可登录，不能收发消息）；2=黑名单（不可登录，自动踢下线，不能收发消息）
type UserGlobalBlack struct {
	UserID     string    `bson:"user_id"`
	Nickname   string    `bson:"nickname"`
	OperatorID string    `bson:"operator_id"`
	Reason     string    `bson:"reason"`
	CreateTime time.Time `bson:"create_time"`
	// Status 限制类型：1=冻结，2=黑名单
	Status int32 `bson:"status"`
}
