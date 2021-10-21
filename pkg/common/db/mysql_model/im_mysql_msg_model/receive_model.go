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

//func InsertMessageToReceive(seq int64, userid, msgid string) error {
//	conn := db.NewDbConnection()
//	receive := Receive{
//		UID: userid,
//		Seq:    seq,
//		MsgId:  msgid,
//	}
//	err := conn.Table("receive").Create(&receive).Error
//	return err
//}
//func GetBiggestSeqFromReceive(userid string) (seq int64, err error) {
//	//得到数据库的连接（并非真连接，调用时才连接，由gorm自动维护数据库连接池）
//	conn := db.NewDbConnection()
//	err = conn.Raw("select max(seq) from receive where user_id = ?", userid).Row().Scan(&seq)
//	return seq, err
//}
