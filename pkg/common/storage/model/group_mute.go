package model

import "time"

// GroupMute is per-user mute of group message notifications (e.g. offline push).
// OwnerUserID must be a group member; MuteEndTime 0 means permanent.
// MuteDuration is the configured interval at set time: -1 permanent, >0 seconds.
type GroupMute struct {
	OwnerUserID  string    `bson:"owner_user_id"`
	GroupID      string    `bson:"group_id"`
	MuteEndTime  int64     `bson:"mute_end_time"`  // Unix seconds; 0 = permanent
	MuteDuration int64     `bson:"mute_duration"`  // -1 permanent, >0 seconds
	CreateTime   time.Time `bson:"create_time"`
}
