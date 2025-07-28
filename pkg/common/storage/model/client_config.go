package model

type ClientConfig struct {
	Key    string `bson:"key"`
	UserID string `bson:"user_id"`
	Value  string `bson:"value"`
}
