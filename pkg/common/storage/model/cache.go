package model

import "time"

type Cache struct {
	Key      string     `bson:"key"`
	Value    string     `bson:"value"`
	ExpireAt *time.Time `bson:"expire_at"`
}
