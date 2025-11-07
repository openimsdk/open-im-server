package model

import (
	"time"
)

type Group struct {
	GroupID                string    `bson:"group_id"`
	GroupName              string    `bson:"group_name"`
	Notification           string    `bson:"notification"`
	Introduction           string    `bson:"introduction"`
	FaceURL                string    `bson:"face_url"`
	CreateTime             time.Time `bson:"create_time"`
	Ex                     string    `bson:"ex"`
	Status                 int32     `bson:"status"`
	CreatorUserID          string    `bson:"creator_user_id"`
	GroupType              int32     `bson:"group_type"`
	NeedVerification       int32     `bson:"need_verification"`
	LookMemberInfo         int32     `bson:"look_member_info"`
	ApplyMemberFriend      int32     `bson:"apply_member_friend"`
	NotificationUpdateTime time.Time `bson:"notification_update_time"`
	NotificationUserID     string    `bson:"notification_user_id"`
}
