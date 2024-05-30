package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type VersionLogElem struct {
	EID        string    `bson:"e_id"`
	Deleted    bool      `bson:"deleted"`
	Version    uint      `bson:"version"`
	LastUpdate time.Time `bson:"last_update"`
}

type VersionLogTable struct {
	ID         primitive.ObjectID `bson:"_id"`
	DID        string             `bson:"d_id"`
	Logs       []VersionLogElem   `bson:"logs"`
	Version    uint               `bson:"version"`
	Deleted    uint               `bson:"deleted"`
	LastUpdate time.Time          `bson:"last_update"`
}

func (v *VersionLogTable) VersionLog() *VersionLog {
	return &VersionLog{
		ID:         v.ID,
		DID:        v.DID,
		Logs:       v.Logs,
		Version:    v.Version,
		Deleted:    v.Deleted,
		LastUpdate: v.LastUpdate,
		LogLen:     0,
		queryDoc:   true,
	}
}

type VersionLog struct {
	ID         primitive.ObjectID `bson:"_id"`
	DID        string             `bson:"d_id"`
	Logs       []VersionLogElem   `bson:"logs"`
	Version    uint               `bson:"version"`
	Deleted    uint               `bson:"deleted"`
	LastUpdate time.Time          `bson:"last_update"`
	LogLen     int                `bson:"log_len"`
	queryDoc   bool               `bson:"-"`
}

func (w *VersionLog) Full() bool {
	return w.queryDoc || w.Version == 0 || len(w.Logs) != w.LogLen
}

func (w *VersionLog) DeleteAndChangeIDs() (delIds []string, changeIds []string) {
	for _, l := range w.Logs {
		if l.Deleted {
			delIds = append(delIds, l.EID)
		} else {
			changeIds = append(changeIds, l.EID)
		}
	}
	return
}
