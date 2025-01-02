// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mgo

import (
	"context"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"

	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewS3Mongo(db *mongo.Database) (database.ObjectInfo, error) {
	coll := db.Collection(database.ObjectName)

	// Create index for name
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}

	// Create index for create_time
	_, err = coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "create_time", Value: 1},
		},
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}

	// Create index for key
	_, err = coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "key", Value: 1},
		},
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}

	return &S3Mongo{coll: coll}, nil
}

type S3Mongo struct {
	coll *mongo.Collection
}

func (o *S3Mongo) SetObject(ctx context.Context, obj *model.Object) error {
	filter := bson.M{"name": obj.Name, "engine": obj.Engine}
	update := bson.M{
		"name":         obj.Name,
		"engine":       obj.Engine,
		"key":          obj.Key,
		"size":         obj.Size,
		"content_type": obj.ContentType,
		"group":        obj.Group,
		"create_time":  obj.CreateTime,
	}
	return mongoutil.UpdateOne(ctx, o.coll, filter, bson.M{"$set": update}, false, options.Update().SetUpsert(true))
}

func (o *S3Mongo) Take(ctx context.Context, engine string, name string) (*model.Object, error) {
	if engine == "" {
		return mongoutil.FindOne[*model.Object](ctx, o.coll, bson.M{"name": name})
	}
	return mongoutil.FindOne[*model.Object](ctx, o.coll, bson.M{"name": name, "engine": engine})
}

func (o *S3Mongo) Delete(ctx context.Context, engine string, name []string) error {
	if len(name) == 0 {
		return nil
	}
	return mongoutil.DeleteOne(ctx, o.coll, bson.M{"engine": engine, "name": bson.M{"$in": name}})
}

func (o *S3Mongo) FindExpirationObject(ctx context.Context, engine string, expiration time.Time, needDelType []string, count int64) ([]*model.Object, error) {
	opt := options.Find()
	if count > 0 {
		opt.SetLimit(count)
	}
	return mongoutil.Find[*model.Object](ctx, o.coll, bson.M{
		"engine":      engine,
		"create_time": bson.M{"$lt": expiration},
		"group":       bson.M{"$in": needDelType},
	}, opt)
}

func (o *S3Mongo) GetKeyCount(ctx context.Context, engine string, key string) (int64, error) {
	return mongoutil.Count(ctx, o.coll, bson.M{"engine": engine, "key": key})
}

func (o *S3Mongo) GetEngineCount(ctx context.Context, engine string) (int64, error) {
	return mongoutil.Count(ctx, o.coll, bson.M{"engine": engine})
}

func (o *S3Mongo) GetEngineInfo(ctx context.Context, engine string, limit int, skip int) ([]*model.Object, error) {
	return mongoutil.Find[*model.Object](ctx, o.coll, bson.M{"engine": engine}, options.Find().SetLimit(int64(limit)).SetSkip(int64(skip)))
}

func (o *S3Mongo) UpdateEngine(ctx context.Context, oldEngine, oldName string, newEngine string) error {
	return mongoutil.UpdateOne(ctx, o.coll, bson.M{"engine": oldEngine, "name": oldName}, bson.M{"$set": bson.M{"engine": newEngine}}, false)
}
