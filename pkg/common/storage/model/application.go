package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Application struct {
	ID         primitive.ObjectID `bson:"_id"`
	Platform   string             `bson:"platform"`
	Hot        bool               `bson:"hot"`
	Version    string             `bson:"version"`
	Url        string             `bson:"url"`
	Text       string             `bson:"text"`
	Force      bool               `bson:"force"`
	Latest     bool               `bson:"latest"`
	CreateTime time.Time          `bson:"create_time"`
}
