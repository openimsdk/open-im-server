package mgo

import (
	"context"
	"errors"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/datautil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func NewVersionLog(coll *mongo.Collection) (database.VersionLog, error) {
	lm := &VersionLogMgo{coll: coll}
	if lm.initIndex(context.Background()) != nil {
		return nil, errs.ErrInternalServer.WrapMsg("init index failed", "coll", coll.Name())
	}
	return lm, nil
}

type VersionLogMgo struct {
	coll *mongo.Collection
}

func (l *VersionLogMgo) initIndex(ctx context.Context) error {
	_, err := l.coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"d_id": 1,
		},
	})
	return err
}

func (l *VersionLogMgo) IncrVersion(ctx context.Context, dId string, eIds []string, state int32) error {
	if len(eIds) == 0 {
		return errs.ErrArgs.WrapMsg("elem id is empty", "dId", dId)
	}
	if datautil.Duplicate(eIds) {
		return errs.ErrArgs.WrapMsg("elem id is duplicate", "dId", dId, "eIds", eIds)
	}
	now := time.Now()
	res, err := l.writeLogBatch(ctx, dId, eIds, state, now)
	if err != nil {
		return err
	}
	if res.MatchedCount > 0 {
		return nil
	}
	if _, err := l.initDoc(ctx, dId, eIds, state, now); err == nil {
		return nil
	} else if !mongo.IsDuplicateKeyError(err) {
		return err
	}
	if res, err := l.writeLogBatch(ctx, dId, eIds, state, now); err != nil {
		return err
	} else if res.MatchedCount == 0 {
		return errs.ErrInternalServer.WrapMsg("mongodb return value that should not occur", "coll", l.coll.Name(), "dId", dId, "eIds", eIds)
	}
	return nil
}

func (l *VersionLogMgo) initDoc(ctx context.Context, dId string, eIds []string, state int32, now time.Time) (*model.VersionLogTable, error) {
	wl := model.VersionLogTable{
		ID:         primitive.NewObjectID(),
		DID:        dId,
		Logs:       make([]model.VersionLogElem, 0, len(eIds)),
		Version:    database.FirstVersion,
		Deleted:    database.DefaultDeleteVersion,
		LastUpdate: now,
	}
	for _, eId := range eIds {
		wl.Logs = append(wl.Logs, model.VersionLogElem{
			EID:        eId,
			State:      state,
			Version:    database.FirstVersion,
			LastUpdate: now,
		})
	}
	_, err := l.coll.InsertOne(ctx, &wl)
	return &wl, err
}

func (l *VersionLogMgo) writeLogBatch(ctx context.Context, dId string, eIds []string, state int32, now time.Time) (*mongo.UpdateResult, error) {
	if eIds == nil {
		eIds = []string{}
	}
	filter := bson.M{
		"d_id": dId,
	}
	elems := make([]bson.M, 0, len(eIds))
	for _, eId := range eIds {
		elems = append(elems, bson.M{
			"e_id":        eId,
			"version":     "$version",
			"state":       state,
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

func (l *VersionLogMgo) findDoc(ctx context.Context, dId string) (*model.VersionLog, error) {
	vl, err := mongoutil.FindOne[*model.VersionLogTable](ctx, l.coll, bson.M{"d_id": dId}, options.FindOne().SetProjection(bson.M{"logs": 0}))
	if err != nil {
		return nil, err
	}
	return vl.VersionLog(), nil
}

func (l *VersionLogMgo) FindChangeLog(ctx context.Context, dId string, version uint, limit int) (*model.VersionLog, error) {
	if wl, err := l.findChangeLog(ctx, dId, version, limit); err == nil {
		return wl, nil
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	if res, err := l.initDoc(ctx, dId, nil, 0, time.Now()); err == nil {
		return res.VersionLog(), nil
	} else if mongo.IsDuplicateKeyError(err) {
		return l.findChangeLog(ctx, dId, version, limit)
	} else {
		return nil, err
	}
}

func (l *VersionLogMgo) findChangeLog(ctx context.Context, dId string, version uint, limit int) (*model.VersionLog, error) {
	if version == 0 && limit == 0 {
		return l.findDoc(ctx, dId)
	}
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
	vl, err := mongoutil.Aggregate[*model.VersionLog](ctx, l.coll, pipeline)
	if err != nil {
		return nil, err
	}
	if len(vl) == 0 {
		return nil, mongo.ErrNoDocuments
	}
	return vl[0], nil
}

func (l *VersionLogMgo) DeleteAfterUnchangedLog(ctx context.Context, deadline time.Time) error {
	return mongoutil.DeleteMany(ctx, l.coll, bson.M{
		"last_update": bson.M{
			"$lt": deadline,
		},
	})
}
