package listdemo

import (
	"context"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/errs"
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

type Elem struct {
	ID      string
	Version uint
}

type ChangeLog struct {
	ChangeIDs []Elem
	DeleteIDs []Elem
}

type WriteLog struct {
	DID           string    `bson:"d_id"`
	Logs          []LogElem `bson:"logs"`
	Version       uint      `bson:"version"`
	LastUpdate    time.Time `bson:"last_update"`
	DeleteVersion uint      `bson:"delete_version"`
}

type LogElem struct {
	EID        string    `bson:"e_id"`
	Deleted    bool      `bson:"deleted"`
	Version    uint      `bson:"version"`
	UpdateTime time.Time `bson:"update_time"`
}

type LogModel struct {
	coll *mongo.Collection
}

func (l *LogModel) InitIndex(ctx context.Context) error {
	return nil
}

func (l *LogModel) WriteLog(ctx context.Context, dId string, eId string, deleted bool) error {
	now := time.Now()
	res, err := l.writeLog(ctx, dId, eId, deleted, now)
	if err != nil {
		return err
	}
	if res.MatchedCount > 0 {
		return nil
	}
	wl := WriteLog{
		DID: dId,
		Logs: []LogElem{
			{
				EID:        eId,
				Deleted:    deleted,
				Version:    1,
				UpdateTime: now,
			},
		},
		Version:       1,
		LastUpdate:    now,
		DeleteVersion: 0,
	}
	if _, err := l.coll.InsertOne(ctx, &wl); err == nil {
		return nil
	} else if !mongo.IsDuplicateKeyError(err) {
		return err
	}
	if res, err := l.writeLog(ctx, dId, eId, deleted, now); err != nil {
		return err
	} else if res.ModifiedCount == 0 {
		return errs.ErrInternalServer.WrapMsg("mongodb return value that should not occur", "coll", l.coll.Name(), "dId", dId, "eId", eId)
	}
	return nil
}

func (l *LogModel) writeLog(ctx context.Context, dId string, eId string, deleted bool, now time.Time) (*mongo.UpdateResult, error) {
	filter := bson.M{
		"d_id": dId,
	}
	elem := bson.M{
		"e_id":        eId,
		"version":     "$version",
		"deleted":     deleted,
		"update_time": now,
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
				"update_time": now,
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

func (l *LogModel) FindChangeLog(ctx context.Context, did string, version uint) (*ChangeLog, error) {
	return nil, nil
}
