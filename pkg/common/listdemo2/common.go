package listdemo

import (
	"context"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var (
	ErrListNotFound = errors.New("list not found")
	ErrElemExist    = errors.New("elem exist")
	ErrNeedFull     = errors.New("need full")
	ErrNotFound     = mongo.ErrNoDocuments
)

const (
	FirstVersion         = 1
	DefaultDeleteVersion = 0
)

type Elem struct {
	ID      string
	Version uint
}

type ChangeLog struct {
	ChangeIDs []Elem
	DeleteIDs []Elem
}

type WriteLog struct {
	DID        string    `bson:"d_id"`
	Logs       []LogElem `bson:"logs"`
	Version    uint      `bson:"version"`
	Deleted    uint      `bson:"deleted"`
	LastUpdate time.Time `bson:"last_update"`
}

type WriteLogLen struct {
	DID        string    `bson:"d_id"`
	Logs       []LogElem `bson:"logs"`
	Version    uint      `bson:"version"`
	Deleted    uint      `bson:"deleted"`
	LastUpdate time.Time `bson:"last_update"`
	LogLen     int       `bson:"log_len"`
}

type LogElem struct {
	EID        string    `bson:"e_id"`
	Deleted    bool      `bson:"deleted"`
	Version    uint      `bson:"version"`
	LastUpdate time.Time `bson:"last_update"`
}

type LogModel struct {
	coll *mongo.Collection
}

func (l *LogModel) InitIndex(ctx context.Context) error {
	_, err := l.coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"d_id": 1,
		},
	})
	return err
}

func (l *LogModel) writeLog(ctx context.Context, dId string, eId string, deleted bool, now time.Time) (*mongo.UpdateResult, error) {
	filter := bson.M{
		"d_id": dId,
	}
	elem := bson.M{
		"e_id":        eId,
		"version":     "$version",
		"deleted":     deleted,
		"last_update": now,
	}
	pipeline := []bson.M{
		{
			"$addFields": bson.M{
				"elem_index": bson.M{
					"$indexOfArray": []any{"$logs.e_id", eId},
				},
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
					"$cond": bson.M{
						"if": bson.M{
							"$lt": []any{"$elem_index", 0},
						},
						"then": bson.M{
							"$concatArrays": []any{
								"$logs",
								[]bson.M{
									elem,
								},
							},
						},
						"else": bson.M{
							"$map": bson.M{
								"input": bson.M{
									"$range": []any{0, bson.M{"$size": "$logs"}},
								},
								"as": "i",
								"in": bson.M{
									"$cond": bson.M{
										"if": bson.M{
											"$eq": []any{"$$i", "$elem_index"},
										},
										"then": elem,
										"else": bson.M{
											"$arrayElemAt": []any{
												"$logs",
												"$$i",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			"$unset": "elem_index",
		},
	}
	return mongoutil.UpdateMany(ctx, l.coll, filter, pipeline)
}

func (l *LogModel) WriteLogBatch(ctx context.Context, dId string, eIds []string, deleted bool) error {
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
	wl := WriteLog{
		DID:        dId,
		Logs:       make([]LogElem, 0, len(eIds)),
		Version:    FirstVersion,
		Deleted:    DefaultDeleteVersion,
		LastUpdate: now,
	}
	for _, eId := range eIds {
		wl.Logs = append(wl.Logs, LogElem{
			EID:        eId,
			Deleted:    deleted,
			Version:    FirstVersion,
			LastUpdate: now,
		})
	}
	if _, err := l.coll.InsertOne(ctx, &wl); err == nil {
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

func (l *LogModel) writeLogBatch(ctx context.Context, dId string, eIds []string, deleted bool, now time.Time) (*mongo.UpdateResult, error) {
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

func (l *LogModel) FindChangeLog(ctx context.Context, did string, version uint, limit int) (*WriteLogLen, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"d_id": did,
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
	res, err := mongoutil.Aggregate[*WriteLogLen](ctx, l.coll, pipeline)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return &WriteLogLen{}, nil
	}
	return res[0], nil
}

func (l *LogModel) DeleteAfterUnchangedLog(ctx context.Context, deadline time.Time) error {
	return mongoutil.DeleteMany(ctx, l.coll, bson.M{
		"last_update": bson.M{
			"$lt": deadline,
		},
	})
}
