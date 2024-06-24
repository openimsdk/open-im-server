package mgo

import (
	"context"
	"errors"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/versionctx"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
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
	_, err := l.IncrVersionResult(ctx, dId, eIds, state)
	return err
}

func (l *VersionLogMgo) IncrVersionResult(ctx context.Context, dId string, eIds []string, state int32) (*model.VersionLog, error) {
	vl, err := l.incrVersionResult(ctx, dId, eIds, state)
	if err != nil {
		return nil, err
	}
	versionctx.GetVersionLog(ctx).Append(versionctx.Collection{
		Name: l.coll.Name(),
		Doc:  vl,
	})
	return vl, nil
}

func (l *VersionLogMgo) incrVersionResult(ctx context.Context, dId string, eIds []string, state int32) (*model.VersionLog, error) {
	if len(eIds) == 0 {
		return nil, errs.ErrArgs.WrapMsg("elem id is empty", "dId", dId)
	}
	now := time.Now()
	if res, err := l.writeLogBatch2(ctx, dId, eIds, state, now); err == nil {
		return res, nil
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	if res, err := l.initDoc(ctx, dId, eIds, state, now); err == nil {
		return res, nil
	} else if !mongo.IsDuplicateKeyError(err) {
		return nil, err
	}
	return l.writeLogBatch2(ctx, dId, eIds, state, now)
}

func (l *VersionLogMgo) initDoc(ctx context.Context, dId string, eIds []string, state int32, now time.Time) (*model.VersionLog, error) {
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
	if _, err := l.coll.InsertOne(ctx, &wl); err != nil {
		return nil, err
	}
	return wl.VersionLog(), nil
}

func (l *VersionLogMgo) writeLogBatch2(ctx context.Context, dId string, eIds []string, state int32, now time.Time) (*model.VersionLog, error) {
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
	opt := options.FindOneAndUpdate().SetUpsert(false).SetReturnDocument(options.After).SetProjection(bson.M{"logs": 0})
	return mongoutil.FindOneAndUpdate[*model.VersionLog](ctx, l.coll, filter, pipeline, opt)
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
	log.ZDebug(ctx, "init doc", "dId", dId)
	if res, err := l.initDoc(ctx, dId, nil, 0, time.Now()); err == nil {
		log.ZDebug(ctx, "init doc success", "dId", dId)
		return res, nil
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

func (l *VersionLogMgo) Delete(ctx context.Context, dId string) error {
	return mongoutil.DeleteOne(ctx, l.coll, bson.M{"d_id": dId})
}
