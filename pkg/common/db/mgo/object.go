package mgo

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/mgo/mtool"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewS3Mongo(db *mongo.Database) (relation.ObjectInfoModelInterface, error) {
	coll := db.Collection("s3")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.M{"name": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
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
	return mtool.UpdateOne(ctx, o.coll, filter, bson.M{"$set": update}, false, options.Update().SetUpsert(true))
}

func (o *S3Mongo) Take(ctx context.Context, engine string, name string) (*relation.ObjectModel, error) {
	if engine == "" {
		return mtool.FindOne[*relation.ObjectModel](ctx, o.coll, bson.M{"name": name})
	}
	return mtool.FindOne[*relation.ObjectModel](ctx, o.coll, bson.M{"name": name, "engine": engine})
}

func (o *S3Mongo) Delete(ctx context.Context, engine string, name string) error {
	return mtool.DeleteOne(ctx, o.coll, bson.M{"name": name, "engine": engine})
}
