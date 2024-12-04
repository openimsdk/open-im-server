package model

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/openimsdk/tools/log"
)

const (
	VersionStateInsert = iota + 1
	VersionStateDelete
	VersionStateUpdate
)

const (
	VersionGroupChangeID = ""
	VersionSortChangeID  = "____S_O_R_T_I_D____"
)

type VersionLogElem struct {
	EID        string    `bson:"e_id"`
	State      int32     `bson:"state"`
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
		LogLen:     len(v.Logs),
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
}

func (v *VersionLog) DeleteAndChangeIDs() (insertIds, deleteIds, updateIds []string) {
	for _, l := range v.Logs {
		switch l.State {
		case VersionStateInsert:
			insertIds = append(insertIds, l.EID)
		case VersionStateDelete:
			deleteIds = append(deleteIds, l.EID)
		case VersionStateUpdate:
			updateIds = append(updateIds, l.EID)
		default:
			log.ZError(context.Background(), "invalid version status found", errors.New("dirty database data"), "objID", v.ID.Hex(), "did", v.DID, "elem", l)
		}
	}
	return
}
