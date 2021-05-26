package im_mysql_model

import "time"

type User struct {
	UID        string    `gorm:"column:uid"`
	Name       string    `gorm:"column:name"`
	Icon       string    `gorm:"column:icon"`
	Gender     int32     `gorm:"column:gender"`
	Mobile     string    `gorm:"column:mobile"`
	Birth      string    `gorm:"column:birth"`
	Email      string    `gorm:"column:email"`
	Ex         string    `gorm:"column:ex"`
	CreateTime time.Time `gorm:"column:create_time"`
}

type Friend struct {
	OwnerId    string    `gorm:"column:owner_id"`
	FriendId   string    `gorm:"column:friend_id"`
	Comment    string    `gorm:"column:comment"`
	FriendFlag int32     `gorm:"column:friend_flag"`
	CreateTime time.Time `gorm:"column:create_time"`
}
type FriendRequest struct {
	ReqId      string    `gorm:"column:req_id"`
	UserId     string    `gorm:"column:user_id"`
	Flag       int32     `gorm:"column:flag"`
	ReqMessage string    `gorm:"column:req_message"`
	CreateTime time.Time `gorm:"column:create_time"`
}
type BlackList struct {
	OwnerId    string    `gorm:"column:owner_id"`
	BlockId    string    `gorm:"column:block_id"`
	CreateTime time.Time `gorm:"column:create_time"`
}
type Group struct {
	GroupId  string `gorm:"column:group_id"`
	Name     string `gorm:"column:name"`
	HeadURL  string `gorm:"column:head_url"`
	Bulletin string `gorm:"column:bulletin"`
}

type GroupMember struct {
	GroupId  string `gorm:"column:group_id"`
	UserId   string `gorm:"column:user_id"`
	NickName string `gorm:"column:nickname"`
	IsAdmin  int32  `gorm:"column:is_admin"`
}
