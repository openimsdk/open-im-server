/*
** description("").
** copyright('tuoyun,www.tuoyun.net').
** author("fg,Gordon@tuoyun.net").
** time(2021/3/4 11:18).
 */
package im_mysql_msg_model

import (
	"time"
)

// Receive Inbox table structure
type Receive struct {
	UserId     string `gorm:"primary_key"` // 收件箱主键ID
	Seq        int64  `gorm:"primary_key"` // 收件箱主键ID
	MsgId      string
	CreateTime *time.Time
}
