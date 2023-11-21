package newmgo

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/newmgo/mgotool"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewS3Mongo(db *mongo.Database) (relation.ObjectInfoModelInterface, error) {
	return &S3Mongo{
		coll: db.Collection("s3"),
	}, nil
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
	return mgotool.UpdateOne(ctx, o.coll, filter, bson.M{"$set": update}, false, options.Update().SetUpsert(true))
}

func (o *S3Mongo) Take(ctx context.Context, engine string, name string) (*relation.ObjectModel, error) {
	if engine == "" {
		return mgotool.FindOne[*relation.ObjectModel](ctx, o.coll, bson.M{"name": name})
	}
	return mgotool.FindOne[*relation.ObjectModel](ctx, o.coll, bson.M{"name": name, "engine": engine})
}

func (o *S3Mongo) Delete(ctx context.Context, engine string, name string) error {
	return mgotool.DeleteOne(ctx, o.coll, bson.M{"name": name, "engine": engine})
}
