package dataver

import (
	"context"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/datautil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	FirstVersion         = 1
	DefaultDeleteVersion = 0
)

type WriteLog struct {
	DID        string    `bson:"d_id"`
	Logs       []Elem    `bson:"logs"`
	Version    uint      `bson:"version"`
	Deleted    uint      `bson:"deleted"`
	LastUpdate time.Time `bson:"last_update"`
	LogLen     int       `bson:"log_len"`
}

func (w *WriteLog) Full() bool {
	if w.Version == 0 {
		return true
	}
	return len(w.Logs) != w.LogLen
}

func (w *WriteLog) DeleteEId() []string {
	var eIds []string
	for _, l := range w.Logs {
		if l.Deleted {
			eIds = append(eIds, l.EID)
		}
	}
	return eIds
}

type Elem struct {
	EID        string    `bson:"e_id"`
	Deleted    bool      `bson:"deleted"`
	Version    uint      `bson:"version"`
	LastUpdate time.Time `bson:"last_update"`
}

type DataLog interface {
	WriteLog(ctx context.Context, dId string, eIds []string, deleted bool) error
	FindChangeLog(ctx context.Context, dId string, version uint, limit int) (*WriteLog, error)
	DeleteAfterUnchangedLog(ctx context.Context, deadline time.Time) error
}

func NewDataLog(coll *mongo.Collection) (DataLog, error) {
	lm := &logModel{coll: coll}
	if lm.initIndex(context.Background()) != nil {
		return nil, errs.ErrInternalServer.WrapMsg("init index failed", "coll", coll.Name())
	}
	return lm, nil
}

type logModel struct {
	coll *mongo.Collection
}

func (l *logModel) initIndex(ctx context.Context) error {
	_, err := l.coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"d_id": 1,
		},
	})
	return err
}

func (l *logModel) WriteLog(ctx context.Context, dId string, eIds []string, deleted bool) error {
	if len(eIds) == 0 {
		return errs.ErrArgs.WrapMsg("elem id is empty", "dId", dId)
	}
	if datautil.Duplicate(eIds) {
		return errs.ErrArgs.WrapMsg("elem id is duplicate", "dId", dId, "eIds", eIds)
	}
	now := time.Now()
	res, err := l.writeLogBatch(ctx, dId, eIds, deleted, now)
	if err != nil {
		return err
	}
	if res.MatchedCount > 0 {
		return nil
	}
	if err := l.initDoc(ctx, dId, eIds, deleted, now); err == nil {
		return nil
	} else if !mongo.IsDuplicateKeyError(err) {
		return err
	}
	if res, err := l.writeLogBatch(ctx, dId, eIds, deleted, now); err != nil {
		return err
	} else if res.ModifiedCount == 0 {
		return errs.ErrInternalServer.WrapMsg("mongodb return value that should not occur", "coll", l.coll.Name(), "dId", dId, "eIds", eIds)
	}
	return nil
}

func (l *logModel) initDoc(ctx context.Context, dId string, eIds []string, deleted bool, now time.Time) error {
	type tableWriteLog struct {
		DID        string    `bson:"d_id"`
		Logs       []Elem    `bson:"logs"`
		Version    uint      `bson:"version"`
		Deleted    uint      `bson:"deleted"`
		LastUpdate time.Time `bson:"last_update"`
	}
	wl := tableWriteLog{
		DID:        dId,
		Logs:       make([]Elem, 0, len(eIds)),
		Version:    FirstVersion,
		Deleted:    DefaultDeleteVersion,
		LastUpdate: now,
	}
	for _, eId := range eIds {
		wl.Logs = append(wl.Logs, Elem{
			EID:        eId,
			Deleted:    deleted,
			Version:    FirstVersion,
			LastUpdate: now,
		})
	}
	_, err := l.coll.InsertOne(ctx, &wl)
	return err
}

func (l *logModel) writeLogBatch(ctx context.Context, dId string, eIds []string, deleted bool, now time.Time) (*mongo.UpdateResult, error) {
	if len(eIds) == 0 {
		return nil, errs.ErrArgs.WrapMsg("elem id is empty", "dId", dId)
	}
	filter := bson.M{
		"d_id": dId,
	}
	elems := make([]bson.M, 0, len(eIds))
	for _, eId := range eIds {
		elems = append(elems, bson.M{
			"e_id":        eId,
			"version":     "$version",
			"deleted":     deleted,
			"last_update": now,
		})
	}
	pipeline := []bson.M{
		{
			"$addFields": bson.M{
				"delete_e_ids": eIds,
			},
		},
		{
			"$set": bson.M{
				"version":     bson.M{"$add": []any{"$version", 1}},
				"last_update": now,
			},
		},
		{
			"$set": bson.M{
				"logs": bson.M{
					"$filter": bson.M{
						"input": "$logs",
						"as":    "log",
						"cond": bson.M{
							"$not": bson.M{
								"$in": []any{"$$log.e_id", "$delete_e_ids"},
							},
						},
					},
				},
			},
		},
		{
			"$set": bson.M{
				"logs": bson.M{
					"$concatArrays": []any{
						"$logs",
						elems,
					},
				},
			},
		},
		{
			"$unset": "delete_e_ids",
		},
	}
	return mongoutil.UpdateMany(ctx, l.coll, filter, pipeline)
}

func (l *logModel) FindChangeLog(ctx context.Context, dId string, version uint, limit int) (*WriteLog, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"d_id": dId,
			},
		},
		{
			"$addFields": bson.M{
				"logs": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$or": []bson.M{
								{"$lt": []any{"$version", version}},
								{"$gte": []any{"$deleted", version}},
							},
						},
						"then": []any{},
						"else": "$logs",
					},
				},
			},
		},
		{
			"$addFields": bson.M{
				"logs": bson.M{
					"$filter": bson.M{
						"input": "$logs",
						"as":    "l",
						"cond": bson.M{
							"$gt": []any{"$$l.version", version},
						},
					},
				},
			},
		},
		{
			"$addFields": bson.M{
				"log_len": bson.M{"$size": "$logs"},
			},
		},
		{
			"$addFields": bson.M{
				"logs": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$gt": []any{"$log_len", limit},
						},
						"then": []any{},
						"else": "$logs",
					},
				},
			},
		},
	}
	if limit <= 0 {
		pipeline = pipeline[:len(pipeline)-1]
	}
	res, err := mongoutil.Aggregate[*WriteLog](ctx, l.coll, pipeline)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errs.Wrap(mongo.ErrNoDocuments)
	}
	return res[0], nil
}

func (l *logModel) DeleteAfterUnchangedLog(ctx context.Context, deadline time.Time) error {
	return mongoutil.DeleteMany(ctx, l.coll, bson.M{
		"last_update": bson.M{
			"$lt": deadline,
		},
	})
}
