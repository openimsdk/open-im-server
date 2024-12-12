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
	"github.com/openimsdk/tools/db/pagination"
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

func (o *S3Mongo) Delete(ctx context.Context, engine string, name string) error {
	return mongoutil.DeleteOne(ctx, o.coll, bson.M{"name": name, "engine": engine})
}

// Find Expires object
func (o *S3Mongo) FindNeedDeleteObjectByDB(ctx context.Context, duration time.Time, needDelType []string, pagination pagination.Pagination) (total int64, objects []*model.Object, err error) {
	return mongoutil.FindPage[*model.Object](ctx, o.coll, bson.M{
		"create_time": bson.M{"$lt": duration},
		"group":       bson.M{"$in": needDelType},
	}, pagination)
}

// Find object by key
func (o *S3Mongo) FindModelsByKey(ctx context.Context, key string) (objects []*model.Object, err error) {
	return mongoutil.Find[*model.Object](ctx, o.coll, bson.M{
		"key": key,
	})
}
