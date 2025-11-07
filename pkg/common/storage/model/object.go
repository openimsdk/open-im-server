package model

import (
	"time"
)

type Object struct {
	Name        string    `bson:"name"`
	UserID      string    `bson:"user_id"`
	Hash        string    `bson:"hash"`
	Engine      string    `bson:"engine"`
	Key         string    `bson:"key"`
	Size        int64     `bson:"size"`
	ContentType string    `bson:"content_type"`
	Group       string    `bson:"group"`
	CreateTime  time.Time `bson:"create_time"`
}
