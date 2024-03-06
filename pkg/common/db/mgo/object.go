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
	"github.com/OpenIMSDK/tools/errs"

	"github.com/OpenIMSDK/tools/mgoutil"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewS3Mongo(db *mongo.Database) (relation.ObjectInfoModelInterface, error) {
	coll := db.Collection("s3")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &S3Mongo{coll: coll}, nil
}

type S3Mongo struct {
	coll *mongo.Collection
}

func (o *S3Mongo) SetObject(ctx context.Context, obj *relation.ObjectModel) error {
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
	return mgoutil.UpdateOne(ctx, o.coll, filter, bson.M{"$set": update}, false, options.Update().SetUpsert(true))
}

func (o *S3Mongo) Take(ctx context.Context, engine string, name string) (*relation.ObjectModel, error) {
	if engine == "" {
		return mgoutil.FindOne[*relation.ObjectModel](ctx, o.coll, bson.M{"name": name})
	}
	return mgoutil.FindOne[*relation.ObjectModel](ctx, o.coll, bson.M{"name": name, "engine": engine})
}

func (o *S3Mongo) Delete(ctx context.Context, engine string, name string) error {
	return mgoutil.DeleteOne(ctx, o.coll, bson.M{"name": name, "engine": engine})
}
