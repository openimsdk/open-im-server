package model

import (
	"time"
)

type Log struct {
	LogID        string    `bson:"log_id"`
	Platform     string    `bson:"platform"`
	UserID       string    `bson:"user_id"`
	CreateTime   time.Time `bson:"create_time"`
	Url          string    `bson:"url"`
	FileName     string    `bson:"file_name"`
	SystemType   string    `bson:"system_type"`
	AppFramework string    `bson:"app_framework"`
	Version      string    `bson:"version"`
	Ex           string    `bson:"ex"`
}
