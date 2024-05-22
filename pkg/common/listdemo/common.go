package listdemo

import (
	"context"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/pagination"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrListNotFound = errors.New("list not found")
	ErrElemExist    = errors.New("elem exist")
	ErrNotFound     = mongo.ErrNoDocuments
)

type ListDoc interface {
	IDName() string               // 外层业务id字段名字 user_id
	ElemsName() string            // 外层列表名字 friends
	VersionName() string          // 外层版本号 version
	DeleteVersion() string        // 删除版本号
	BuildDoc(lid any, e Elem) any // 返回一个组装的doc文档
}

type Elem interface {
	IDName() string        // 业务id名字 friend_user_id
	IDValue() any          // 业务id值 userID -> "100000000"
	VersionName() string   // 版本号
	DeletedName() string   // 删除字段名字
	ToMap() map[string]any // 把结构体转换为map
}

type List[D any, E Elem] struct {
	coll *mongo.Collection
	lf   ListDoc
}

func (l *List[D, E]) zeroE() E {
	var t E
	return t
}

func (l *List[D, E]) FindElem(ctx context.Context, lid any, eid any) (E, error) {
	res, err := l.FindElems(ctx, lid, []any{eid})
	if err != nil {
		return l.zeroE(), err
	}
	if len(res) == 0 {
		return l.zeroE(), ErrNotFound
	}
	return res[0], nil
}

// FindElems 查询Elems
func (l *List[D, E]) FindElems(ctx context.Context, lid any, eids []any) ([]E, error) {
	//pipeline := []bson.M{
	//	{
	//		"$match": bson.M{
	//			l.lf.IDName(): lid,
	//			l.lf.IDName() + "." + l.lf.ElemsID(): bson.M{
	//				"$in": eids,
	//			},
	//		},
	//	},
	//	{
	//		"$unwind": "$" + l.lf.ElemsName(),
	//	},
	//	{
	//		"$match": bson.M{
	//			l.lf.IDName() + "." + l.lf.ElemsID(): bson.M{
	//				"$in": eids,
	//			},
	//		},
	//	},
	//}
	panic("todo")
}

func (l *List[D, E]) Find(ctx context.Context, filter any, opts ...*options.FindOptions) ([]E, error) {
	return nil, nil
}

func (l *List[D, E]) Count(ctx context.Context, filter any, opts ...*options.CountOptions) (int64, error) {
	return 0, nil
}

func (l *List[D, E]) Update(ctx context.Context, lid any, eid any) (*mongo.UpdateResult, error) {

	return nil, nil
}

func (l *List[D, E]) Delete(ctx context.Context, lid any, eids any) (*mongo.UpdateResult, error) {

	return nil, nil
}

func (l *List[D, E]) Page(ctx context.Context, filter any, pagination pagination.Pagination, opts ...*options.FindOptions) (int64, []E, error) {
	return 0, nil, nil
}

func (l *List[D, E]) ElemIDs(ctx context.Context, filter any, opts ...*options.FindOptions) ([]E, error) {

	return nil, nil
}

// InsertElem 插入一个
func (l *List[D, E]) InsertElem(ctx context.Context, lid any, e Elem) error {
	if err := l.insertElem(ctx, lid, e); err == nil {
		return nil
	} else if !errors.Is(err, ErrListNotFound) {
		return err
	}
	if _, err := l.coll.InsertOne(ctx, l.lf.BuildDoc(lid, e)); err == nil {
		return nil
	} else if mongo.IsDuplicateKeyError(err) {
		return l.insertElem(ctx, lid, e)
	} else {
		return err
	}
}

func (l *List[D, E]) insertElem(ctx context.Context, lid any, e Elem) error {
	data := e.ToMap()
	data[e.VersionName()] = "$max_version"
	filter := bson.M{
		l.lf.IDName(): lid,
	}
	pipeline := []bson.M{
		{
			"$addFields": bson.M{
				"found_elem": bson.M{
					"$in": bson.A{e.IDValue(), l.lf.ElemsName() + "." + e.IDName()},
				},
			},
		},
		{
			"$set": bson.M{
				"max_version": bson.M{
					"$cond": bson.M{
						"if":   "$found_elem",
						"then": "$max_version",
						"else": bson.M{"$add": bson.A{"max_version", 1}},
					},
				},
			},
		},
		{
			"$set": bson.M{
				l.lf.ElemsName(): bson.M{
					"$cond": bson.M{
						"if":   "$found_elem",
						"then": "$" + l.lf.ElemsName(),
						"else": bson.M{
							"$concatArrays": bson.A{
								"$" + l.lf.ElemsName(),
								bson.A{
									data,
								},
							},
						},
					},
				},
			},
		},
		{
			"$unset": "found_elem",
		},
	}
	res, err := mongoutil.UpdateMany(ctx, l.coll, filter, pipeline)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrListNotFound
	}
	if res.ModifiedCount == 0 {
		return ErrElemExist
	}
	return nil
}
