package mgo

import (
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewMsgDelMongo(db *mongo.Database) *MsgDelMgo {
	return &MsgDelMgo{coll: db.Collection(new(relation.MsgDocModel).TableName())}
}

type MsgDelMgo struct {
	coll  *mongo.Collection
	model relation.MsgDocModel
}

//func (m *MsgDelMgo) getEmptyMsg(ctx context.Context, limit int) ([]string, error) {
//	return mongoutil.Aggregate[string](ctx, m.coll, []bson.M{
//		{
//			"$match": bson.M{
//				"msgs": bson.M{
//					"$exists": true,
//				},
//			},
//		},
//		{
//			"$project": bson.M{
//				"_id":    0,
//				"doc_id": 1,
//				"all_null_msgs": bson.M{
//					"$not": []bson.M{
//						{
//							"$anyElementTrue": bson.M{
//								"$map": bson.M{
//									"input": "$msgs",
//									"as":    "item",
//									"in":    "$$item.msg",
//								},
//							},
//						},
//					},
//				},
//			},
//		},
//		{
//			"$project": bson.M{
//				"doc_id": 1,
//			},
//		},
//		{
//			"$limit": limit,
//		},
//	})
//}
//
//func (m *MsgDelMgo) deleteEmptyMsgs(ctx context.Context) error {
//	for {
//		docIDs, err := m.getEmptyMsg(ctx, 100)
//		if err != nil {
//			return err
//		}
//		if len(docIDs) == 0 {
//			return nil
//		}
//		for _, docID := range docIDs {
//			if err := m.deleteEmptyMsg(ctx, docID); err != nil {
//				return err
//			}
//		}
//	}
//}
