package mgotool

import (
	"context"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/openimsdk/open-im-server/v3/pkg/common/pagination"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Anys[T any](ts []T) []any {
	val := make([]any, len(ts))
	for i := range ts {
		val[i] = ts[i]
	}
	return val
}

func findOptionToCountOption(opts []*options.FindOptions) *options.CountOptions {
	countOpt := options.Count()
	for _, opt := range opts {
		if opt.Skip != nil {
			countOpt.SetSkip(*opt.Skip)
		}
		if opt.Limit != nil {
			countOpt.SetLimit(*opt.Limit)
		}
	}
	return countOpt
}

func InsertMany[T any](ctx context.Context, coll *mongo.Collection, val []T, opts ...*options.InsertManyOptions) error {
	_, err := coll.InsertMany(ctx, Anys(val), opts...)
	if err != nil {
		return errs.Wrap(err)
	}
	return nil
}

func UpdateOne(ctx context.Context, coll *mongo.Collection, filter any, update any, notMatchedErr bool, opts ...*options.UpdateOptions) error {
	res, err := coll.UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		return errs.Wrap(err)
	}
	if notMatchedErr && res.MatchedCount == 0 {
		return errs.Wrap(mongo.ErrNoDocuments)
	}
	return nil
}

func Find[T any](ctx context.Context, coll *mongo.Collection, filter any, opts ...*options.FindOptions) ([]T, error) {
	cur, err := coll.Find(ctx, filter, opts...)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	defer cur.Close(ctx)
	var res []T
	if err := cur.All(ctx, &res); err != nil {
		return nil, errs.Wrap(err)
	}
	return res, nil
}

func FindOne[T any](ctx context.Context, coll *mongo.Collection, filter any, opts ...*options.FindOneOptions) (res T, err error) {
	cur := coll.FindOne(ctx, filter, opts...)
	if err := cur.Err(); err != nil {
		return res, errs.Wrap(err)
	}
	if err := cur.Decode(&res); err != nil {
		return res, errs.Wrap(err)
	}
	return res, nil
}

func FindPage[T any](ctx context.Context, coll *mongo.Collection, filter any, pagination pagination.Pagination, opts ...*options.FindOptions) (int64, []T, error) {
	count, err := Count(ctx, coll, filter, findOptionToCountOption(opts))
	if err != nil {
		return 0, nil, err
	}
	if count == 0 || pagination == nil {
		return count, nil, nil
	}
	skip := int64(pagination.GetPageNumber()-1) * int64(pagination.GetShowNumber())
	if skip < 0 || skip >= count || pagination.GetShowNumber() <= 0 {
		return count, nil, nil
	}
	opt := options.Find().SetSkip(skip).SetLimit(int64(pagination.GetShowNumber()))
	res, err := Find[T](ctx, coll, filter, append(opts, opt)...)
	if err != nil {
		return 0, nil, err
	}
	return count, res, nil
}

func Count(ctx context.Context, coll *mongo.Collection, filter any, opts ...*options.CountOptions) (int64, error) {
	return coll.CountDocuments(ctx, filter, opts...)
}

func Exist(ctx context.Context, coll *mongo.Collection, filter any, opts ...*options.CountOptions) (bool, error) {
	opts = append(opts, options.Count().SetLimit(1))
	count, err := Count(ctx, coll, filter, opts...)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func DeleteOne(ctx context.Context, coll *mongo.Collection, filter any, opts ...*options.DeleteOptions) error {
	if _, err := coll.DeleteOne(ctx, filter, opts...); err != nil {
		return errs.Wrap(err)
	}
	return nil
}

func DeleteMany(ctx context.Context, coll *mongo.Collection, filter any, opts ...*options.DeleteOptions) error {
	if _, err := coll.DeleteMany(ctx, filter, opts...); err != nil {
		return errs.Wrap(err)
	}
	return nil
}

//func Upsert[T any](ctx context.Context, coll *mongo.Collection, val *T, opts ...*options.InsertManyOptions) error {
//	return nil
//}
