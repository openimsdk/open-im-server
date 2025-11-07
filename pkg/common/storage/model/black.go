package model

import (
	"time"
)

type Black struct {
	OwnerUserID    string    `bson:"owner_user_id"`
	BlockUserID    string    `bson:"block_user_id"`
	CreateTime     time.Time `bson:"create_time"`
	AddSource      int32     `bson:"add_source"`
	OperatorUserID string    `bson:"operator_user_id"`
	Ex             string    `bson:"ex"`
}
