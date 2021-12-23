package im_mysql_model

import "time"

type User struct {
	UserID     string    `gorm:"column:uid;primaryKey;"`
	Nickname   string    `gorm:"column:name"`
	FaceUrl    string    `gorm:"column:icon"`
	Gender     int32     `gorm:"column:gender"`
	Mobile     string    `gorm:"column:mobile"`
	Birth      string    `gorm:"column:birth"`
	Email      string    `gorm:"column:email"`
	Ex         string    `gorm:"column:ex"`
	CreateTime time.Time `gorm:"column:create_time"`
}

type Friend struct {
	OwnerUserID  string    `gorm:"column:owner_id"`
	FriendUserID string    `gorm:"column:friend_id"`
	Remark       string    `gorm:"column:comment"`
	FriendFlag   int32     `gorm:"column:friend_flag"`
	CreateTime   time.Time `gorm:"column:create_time"`
}
type FriendRequest struct {
	ReqID      string    `gorm:"column:req_id"`
	UserID     string    `gorm:"column:user_id"`
	Flag       int32     `gorm:"column:flag"`
	ReqMessage string    `gorm:"column:req_message"`
	CreateTime time.Time `gorm:"column:create_time"`
}
type BlackList struct {
	OwnerUserID string    `gorm:"column:owner_id"`
	BlockUserID string    `gorm:"column:block_id"`
	CreateTime  time.Time `gorm:"column:create_time"`
}

type Group struct {
	GroupID      string    `gorm:"column:group_id"`
	GroupName    string    `gorm:"column:name"`
	Introduction string    `gorm:"column:introduction"`
	Notification string    `gorm:"column:notification"`
	FaceUrl      string    `gorm:"column:face_url"`
	CreateTime   time.Time `gorm:"column:create_time"`
	Ext          string    `gorm:"column:ex"`
}

type GroupMember struct {
	GroupID            string    `gorm:"column:group_id"`
	UserID             string    `gorm:"column:uid"`
	NickName           string    `gorm:"column:nickname"`
	AdministratorLevel int32     `gorm:"column:administrator_level"`
	JoinTime           time.Time `gorm:"column:join_time"`
	FaceUrl            string    `gorm:"user_group_face_url"`
}

type GroupRequest struct {
	ID               string    `gorm:"column:id"`
	GroupID          string    `gorm:"column:group_id"`
	FromUserID       string    `gorm:"column:from_user_id"`
	ToUserID         string    `gorm:"column:to_user_id"`
	Flag             int32     `gorm:"column:flag"`
	ReqMsg           string    `gorm:"column:req_msg"`
	HandledMsg       string    `gorm:"column:handled_msg"`
	CreateTime       time.Time `gorm:"column:create_time"`
	FromUserNickname string    `gorm:"from_user_nickname"`
	ToUserNickname   string    `gorm:"to_user_nickname"`
	FromUserFaceUrl  string    `gorm:"from_user_face_url"`
	ToUserFaceUrl    string    `gorm:"to_user_face_url"`
	HandledUser      string    `gorm:"handled_user"`
}
