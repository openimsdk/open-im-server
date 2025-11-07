package model

import (
	"time"
)

type GroupMember struct {
	GroupID        string    `bson:"group_id"`
	UserID         string    `bson:"user_id"`
	Nickname       string    `bson:"nickname"`
	FaceURL        string    `bson:"face_url"`
	RoleLevel      int32     `bson:"role_level"`
	JoinTime       time.Time `bson:"join_time"`
	JoinSource     int32     `bson:"join_source"`
	InviterUserID  string    `bson:"inviter_user_id"`
	OperatorUserID string    `bson:"operator_user_id"`
	MuteEndTime    time.Time `bson:"mute_end_time"`
	Ex             string    `bson:"ex"`
}
