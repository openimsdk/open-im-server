package im_mysql_model

import "time"

type Friend struct {
	OwnerUserID    string    `gorm:"column:owner_user_id;primaryKey;"`
	FriendUserID   string    `gorm:"column:friend_user_id;primaryKey;"`
	Remark         string    `gorm:"column:remark"`
	CreateTime     time.Time `gorm:"column:create_time"`
	AddSource      int32     `gorm:"column:add_source"`
	OperatorUserID string    `gorm:"column:operator_user_id"`
	Ex             string    `gorm:"column:ex"`
}

type FriendRequest struct {
	FromUserID    string    `gorm:"column:from_user_id;primaryKey;"`
	ToUserID      string    `gorm:"column:to_user_id;primaryKey;"`
	HandleResult  int32     `gorm:"column:handle_result"`
	ReqMessage    string    `gorm:"column:req_message"`
	CreateTime    time.Time `gorm:"column:create_time"`
	HandlerUserID string    `gorm:"column:handler_user_id"`
	HandleMsg     string    `gorm:"column:handle_msg"`
	HandleTime    time.Time `gorm:"column:handle_time"`
	Ex            string    `gorm:"column:ex"`
}

type Group struct {
	GroupID       string    `gorm:"column:group_id;primaryKey;"`
	GroupName     string    `gorm:"column:name"`
	Introduction  string    `gorm:"column:introduction"`
	Notification  string    `gorm:"column:notification"`
	FaceUrl       string    `gorm:"column:face_url"`
	CreateTime    time.Time `gorm:"column:create_time"`
	Status        int32     `gorm:"column:status"`
	CreatorUserID string    `gorm:"column:creator_user_id"`
	GroupType     int32     `gorm:"column:group_type"`
	Ex            string    `gorm:"column:ex"`
}

type GroupMember struct {
	GroupID        string    `gorm:"column:group_id;primaryKey;"`
	UserID         string    `gorm:"column:user_id;primaryKey;"`
	Nickname       string    `gorm:"column:nickname"`
	FaceUrl        string    `gorm:"user_group_face_url"`
	RoleLevel      int32     `gorm:"column:role_level"`
	JoinTime       time.Time `gorm:"column:join_time"`
	JoinSource     int32     `gorm:"column:join_source"`
	OperatorUserID string    `gorm:"column:operator_user_id"`
	Ex             string    `gorm:"column:ex"`
}

type GroupRequest struct {
	UserID       string    `gorm:"column:user_id;primaryKey;"`
	GroupID      string    `gorm:"column:group_id;primaryKey;"`
	HandleResult int32     `gorm:"column:handle_result"`
	ReqMsg       string    `gorm:"column:req_msg"`
	HandledMsg   string    `gorm:"column:handled_msg"`
	ReqTime      time.Time `gorm:"column:req_time"`
	HandleUserID string    `gorm:"column:handle_user_id"`
	HandledTime  time.Time `gorm:"column:handle_time"`
	Ex           string    `gorm:"column:ex"`
}

type User struct {
	UserID      string    `gorm:"column:user_id;primaryKey;"`
	Nickname    string    `gorm:"column:name"`
	FaceUrl     string    `gorm:"column:icon"`
	Gender      int32     `gorm:"column:gender"`
	PhoneNumber string    `gorm:"column:phone_number"`
	Birth       string    `gorm:"column:birth"`
	Email       string    `gorm:"column:email"`
	Ex          string    `gorm:"column:ex"`
	CreateTime  time.Time `gorm:"column:create_time"`
}

type Black struct {
	OwnerUserID    string    `gorm:"column:owner_user_id;primaryKey;"`
	BlockUserID    string    `gorm:"column:block_user_id;primaryKey;"`
	CreateTime     time.Time `gorm:"column:create_time"`
	AddSource      int32     `gorm:"column:add_source"`
	OperatorUserID int32     `gorm:"column:operator_user_id"`
	Ex             string    `gorm:"column:ex"`
}
