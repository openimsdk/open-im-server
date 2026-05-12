package model

import "time"

// UserMute records a mute relationship: OwnerUserID has muted MutedUserID.
// Works for both friends and strangers. MuteEndTime == 0 means permanent mute.
type UserMute struct {
	OwnerUserID string    `bson:"owner_user_id"` // who set the mute
	MutedUserID string    `bson:"muted_user_id"` // who is muted
	MuteEndTime int64     `bson:"mute_end_time"` // Unix seconds; 0 = permanent
	CreateTime  time.Time `bson:"create_time"`
}
