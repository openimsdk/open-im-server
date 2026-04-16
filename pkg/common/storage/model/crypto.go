package model

import "time"

type CryptoDevice struct {
	DeviceID       string    `bson:"device_id"`
	UserID         string    `bson:"user_id"`
	Platform       string    `bson:"platform"`
	DeviceModel    string    `bson:"device_model"`
	AppVersion     string    `bson:"app_version"`
	VirgilIdentity string    `bson:"virgil_identity"`
	Status         string    `bson:"status"`
	LastSeenAt     time.Time `bson:"last_seen_at"`
	CreateTime     time.Time `bson:"create_time"`
}

type GroupKeyVersion struct {
	GroupID         string `bson:"group_id"`
	GroupKeyVersion int64  `bson:"group_key_version"`
}

type GroupKeyEvent struct {
	EventID         string    `bson:"event_id"`
	GroupID         string    `bson:"group_id"`
	GroupKeyVersion int64     `bson:"group_key_version"`
	EventType       string    `bson:"event_type"`
	OperatorUserID  string    `bson:"operator_user_id"`
	CreateTime      time.Time `bson:"create_time"`
}
