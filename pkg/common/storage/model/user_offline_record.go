package model

import "time"

// UserOfflineRecord 记录用户全平台离线的时刻及账号自动删除截止时间。
// 用户上线时删除记录；用户全部平台离线时 upsert 记录。
// crontask 每小时扫描此集合，删除 DeleteUserDeadline <= now 的账号。
type UserOfflineRecord struct {
	UserID             string    `bson:"user_id"`
	OfflineTime        time.Time `bson:"offline_time"`
	// DeleteUserDeadline = OfflineTime + delete_account_interval（秒）
	// 用户修改 delete_account_interval 时同步刷新此字段。
	DeleteUserDeadline time.Time `bson:"delete_user_deadline"`
}
