package mgo

import (
	"context"

	"github.com/OpenIMSDK/tools/mgoutil"
	"github.com/OpenIMSDK/tools/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
)

func NewBlackMongo(db *mongo.Database) (relation.BlackModelInterface, error) {
	coll := db.Collection("black")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "owner_user_id", Value: 1},
			{Key: "block_user_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
	}
	return &BlackMgo{coll: coll}, nil
}

type BlackMgo struct {
	coll *mongo.Collection
}

func (b *BlackMgo) blackFilter(ownerUserID, blockUserID string) bson.M {
	return bson.M{
		"owner_user_id": ownerUserID,
		"block_user_id": blockUserID,
	}
}

func (b *BlackMgo) blacksFilter(blacks []*relation.BlackModel) bson.M {
	if len(blacks) == 0 {
		return nil
	}
	or := make(bson.A, 0, len(blacks))
	for _, black := range blacks {
		or = append(or, b.blackFilter(black.OwnerUserID, black.BlockUserID))
	}
	return bson.M{"$or": or}
}

func (b *BlackMgo) Create(ctx context.Context, blacks []*relation.BlackModel) (err error) {
	return mgoutil.InsertMany(ctx, b.coll, blacks)
}

func (b *BlackMgo) Delete(ctx context.Context, blacks []*relation.BlackModel) (err error) {
	if len(blacks) == 0 {
		return nil
	}
	return mgoutil.DeleteMany(ctx, b.coll, b.blacksFilter(blacks))
}

func (b *BlackMgo) UpdateByMap(ctx context.Context, ownerUserID, blockUserID string, args map[string]any) (err error) {
	if len(args) == 0 {
		return nil
	}
	return mgoutil.UpdateOne(ctx, b.coll, b.blackFilter(ownerUserID, blockUserID), bson.M{"$set": args}, false)
}

func (b *BlackMgo) Find(ctx context.Context, blacks []*relation.BlackModel) (blackList []*relation.BlackModel, err error) {
	return mgoutil.Find[*relation.BlackModel](ctx, b.coll, b.blacksFilter(blacks))
}

func (b *BlackMgo) Take(ctx context.Context, ownerUserID, blockUserID string) (black *relation.BlackModel, err error) {
	return mgoutil.FindOne[*relation.BlackModel](ctx, b.coll, b.blackFilter(ownerUserID, blockUserID))
}

func (b *BlackMgo) FindOwnerBlacks(ctx context.Context, ownerUserID string, pagination pagination.Pagination) (total int64, blacks []*relation.BlackModel, err error) {
	return mgoutil.FindPage[*relation.BlackModel](ctx, b.coll, bson.M{"owner_user_id": ownerUserID}, pagination)
}

func (b *BlackMgo) FindOwnerBlackInfos(ctx context.Context, ownerUserID string, userIDs []string) (blacks []*relation.BlackModel, err error) {
	if len(userIDs) == 0 {
		return mgoutil.Find[*relation.BlackModel](ctx, b.coll, bson.M{"owner_user_id": ownerUserID})
	}
	return mgoutil.Find[*relation.BlackModel](ctx, b.coll, bson.M{"owner_user_id": ownerUserID, "block_user_id": bson.M{"$in": userIDs}})
}

func (b *BlackMgo) FindBlackUserIDs(ctx context.Context, ownerUserID string) (blackUserIDs []string, err error) {
	return mgoutil.Find[string](ctx, b.coll, bson.M{"owner_user_id": ownerUserID}, options.Find().SetProjection(bson.M{"_id": 0, "block_user_id": 1}))
}
