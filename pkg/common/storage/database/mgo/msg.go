package mgo

import (
	"context"
	"fmt"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/jsonutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMsgMongo(db *mongo.Database) (database.Msg, error) {
	coll := db.Collection(new(model.MsgDocModel).TableName())
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "doc_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &MsgMgo{coll: coll}, nil
}

type MsgMgo struct {
	coll  *mongo.Collection
	model model.MsgDocModel
}

func (m *MsgMgo) Create(ctx context.Context, msg *model.MsgDocModel) error {
	return mongoutil.InsertMany(ctx, m.coll, []*model.MsgDocModel{msg})
}

func (m *MsgMgo) UpdateMsg(ctx context.Context, docID string, index int64, key string, value any) (*mongo.UpdateResult, error) {
	var field string
	if key == "" {
		field = fmt.Sprintf("msgs.%d", index)
	} else {
		field = fmt.Sprintf("msgs.%d.%s", index, key)
	}
	filter := bson.M{"doc_id": docID}
	update := bson.M{"$set": bson.M{field: value}}
	return mongoutil.UpdateOneResult(ctx, m.coll, filter, update)
}

func (m *MsgMgo) PushUnique(ctx context.Context, docID string, index int64, key string, value any) (*mongo.UpdateResult, error) {
	var field string
	if key == "" {
		field = fmt.Sprintf("msgs.%d", index)
	} else {
		field = fmt.Sprintf("msgs.%d.%s", index, key)
	}
	filter := bson.M{"doc_id": docID}
	update := bson.M{
		"$addToSet": bson.M{
			field: bson.M{"$each": value},
		},
	}
	return mongoutil.UpdateOneResult(ctx, m.coll, filter, update)
}

func (m *MsgMgo) FindOneByDocID(ctx context.Context, docID string) (*model.MsgDocModel, error) {
	return mongoutil.FindOne[*model.MsgDocModel](ctx, m.coll, bson.M{"doc_id": docID})
}

func (m *MsgMgo) GetMsgBySeqIndexIn1Doc(ctx context.Context, userID, docID string, seqs []int64) ([]*model.MsgInfoModel, error) {
	msgs, err := m.getMsgBySeqIndexIn1Doc(ctx, userID, docID, seqs)
	if err != nil {
		return nil, err
	}
	if len(msgs) == len(seqs) {
		return msgs, nil
	}
	tmp := make(map[int64]*model.MsgInfoModel)
	for i, val := range msgs {
		tmp[val.Msg.Seq] = msgs[i]
	}
	res := make([]*model.MsgInfoModel, 0, len(seqs))
	for _, seq := range seqs {
		if val, ok := tmp[seq]; ok {
			res = append(res, val)
		} else {
			res = append(res, &model.MsgInfoModel{Msg: &model.MsgDataModel{Seq: seq}})
		}
	}
	return res, nil
}

func (m *MsgMgo) getMsgBySeqIndexIn1Doc(ctx context.Context, userID, docID string, seqs []int64) ([]*model.MsgInfoModel, error) {
	indexes := make([]int64, 0, len(seqs))
	for _, seq := range seqs {
		indexes = append(indexes, m.model.GetMsgIndex(seq))
	}
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "doc_id", Value: docID},
		}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "doc_id", Value: 1},
			{Key: "msgs", Value: bson.D{
				{Key: "$map", Value: bson.D{
					{Key: "input", Value: indexes},
					{Key: "as", Value: "index"},
					{Key: "in", Value: bson.D{
						{Key: "$arrayElemAt", Value: bson.A{"$msgs", "$$index"}},
					}},
				}},
			}},
		}}},
	}
	msgDocModel, err := mongoutil.Aggregate[*model.MsgDocModel](ctx, m.coll, pipeline)
	if err != nil {
		return nil, err
	}
	if len(msgDocModel) == 0 {
		return nil, errs.Wrap(mongo.ErrNoDocuments)
	}
	msgs := make([]*model.MsgInfoModel, 0, len(msgDocModel[0].Msg))
	for i := range msgDocModel[0].Msg {
		msg := msgDocModel[0].Msg[i]
		if msg == nil || msg.Msg == nil {
			continue
		}
		if datautil.Contain(userID, msg.DelList...) {
			msg.Msg.Content = ""
			msg.Msg.Status = constant.MsgDeleted
		}
		if msg.Revoke != nil {
			revokeContent := sdkws.MessageRevokedContent{
				RevokerID:                   msg.Revoke.UserID,
				RevokerRole:                 msg.Revoke.Role,
				ClientMsgID:                 msg.Msg.ClientMsgID,
				RevokerNickname:             msg.Revoke.Nickname,
				RevokeTime:                  msg.Revoke.Time,
				SourceMessageSendTime:       msg.Msg.SendTime,
				SourceMessageSendID:         msg.Msg.SendID,
				SourceMessageSenderNickname: msg.Msg.SenderNickname,
				SessionType:                 msg.Msg.SessionType,
				Seq:                         msg.Msg.Seq,
				Ex:                          msg.Msg.Ex,
			}
			data, err := jsonutil.JsonMarshal(&revokeContent)
			if err != nil {
				return nil, errs.WrapMsg(err, fmt.Sprintf("docID is %s, seqs is %v", docID, seqs))
			}
			elem := sdkws.NotificationElem{
				Detail: string(data),
			}
			content, err := jsonutil.JsonMarshal(&elem)
			if err != nil {
				return nil, errs.WrapMsg(err, fmt.Sprintf("docID is %s, seqs is %v", docID, seqs))
			}
			msg.Msg.ContentType = constant.MsgRevokeNotification
			msg.Msg.Content = string(content)
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

func (m *MsgMgo) GetNewestMsg(ctx context.Context, conversationID string) (*model.MsgInfoModel, error) {
	for skip := int64(0); ; skip++ {
		msgDocModel, err := m.GetMsgDocModelByIndex(ctx, conversationID, skip, -1)
		if err != nil {
			return nil, err
		}
		for i := len(msgDocModel.Msg) - 1; i >= 0; i-- {
			if msgDocModel.Msg[i].Msg != nil {
				return msgDocModel.Msg[i], nil
			}
		}
	}
}

func (m *MsgMgo) GetOldestMsg(ctx context.Context, conversationID string) (*model.MsgInfoModel, error) {
	for skip := int64(0); ; skip++ {
		msgDocModel, err := m.GetMsgDocModelByIndex(ctx, conversationID, skip, 1)
		if err != nil {
			return nil, err
		}
		for i, v := range msgDocModel.Msg {
			if v.Msg != nil {
				return msgDocModel.Msg[i], nil
			}
		}
	}
}

func (m *MsgMgo) GetMsgDocModelByIndex(ctx context.Context, conversationID string, index, sort int64) (*model.MsgDocModel, error) {
	if sort != 1 && sort != -1 {
		return nil, errs.ErrArgs.WrapMsg("mongo sort must be 1 or -1")
	}
	opt := options.Find().SetSkip(index).SetSort(bson.M{"_id": sort}).SetLimit(1)
	filter := bson.M{"doc_id": primitive.Regex{Pattern: fmt.Sprintf("^%s:", conversationID)}}
	msgs, err := mongoutil.Find[*model.MsgDocModel](ctx, m.coll, filter, opt)
	if err != nil {
		return nil, err
	}
	if len(msgs) > 0 {
		return msgs[0], nil
	}
	return nil, errs.Wrap(model.ErrMsgListNotExist)
}

func (m *MsgMgo) DeleteMsgsInOneDocByIndex(ctx context.Context, docID string, indexes []int) error {
	update := bson.M{
		"$set": bson.M{},
	}
	for _, index := range indexes {
		update["$set"].(bson.M)[fmt.Sprintf("msgs.%d", index)] = bson.M{
			"msg": nil,
		}
	}
	_, err := mongoutil.UpdateMany(ctx, m.coll, bson.M{"doc_id": docID}, update)
	return err
}

func (m *MsgMgo) MarkSingleChatMsgsAsRead(ctx context.Context, userID string, docID string, indexes []int64) error {
	var updates []mongo.WriteModel
	for _, index := range indexes {
		filter := bson.M{
			"doc_id": docID,
			fmt.Sprintf("msgs.%d.msg.send_id", index): bson.M{
				"$ne": userID,
			},
		}
		update := bson.M{
			"$set": bson.M{
				fmt.Sprintf("msgs.%d.is_read", index): true,
			},
		}
		updateModel := mongo.NewUpdateManyModel().
			SetFilter(filter).
			SetUpdate(update)
		updates = append(updates, updateModel)
	}
	if _, err := m.coll.BulkWrite(ctx, updates); err != nil {
		return errs.WrapMsg(err, fmt.Sprintf("docID is %s, indexes is %v", docID, indexes))
	}
	return nil
}

type searchMessageIndex struct {
	ID    primitive.ObjectID `bson:"_id"`
	Index []int64            `bson:"index"`
}

func (m *MsgMgo) searchMessageIndex(ctx context.Context, filter any, nextID primitive.ObjectID, limit int) ([]searchMessageIndex, error) {
	var pipeline bson.A
	if !nextID.IsZero() {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"_id": bson.M{"$gt": nextID}}})
	}
	coarseFilter := bson.M{
		"$or": bson.A{
			bson.M{
				"doc_id": primitive.Regex{Pattern: "^sg_"},
			},
			bson.M{
				"doc_id": primitive.Regex{Pattern: "^si_"},
			},
		},
	}
	pipeline = append(pipeline,
		bson.M{"$sort": bson.M{"_id": 1}},
		bson.M{"$match": coarseFilter},
		bson.M{"$match": filter},
		bson.M{"$limit": limit},
		bson.M{
			"$project": bson.M{
				"_id": 1,
				"msgs": bson.M{
					"$map": bson.M{
						"input": "$msgs",
						"as":    "msg",
						"in": bson.M{
							"$mergeObjects": bson.A{
								"$$msg",
								bson.M{
									"_search_temp_index": bson.M{
										"$indexOfArray": bson.A{
											"$msgs", "$$msg",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		bson.M{"$unwind": "$msgs"},
		bson.M{"$match": filter},
		bson.M{
			"$project": bson.M{
				"_id":                     1,
				"msgs._search_temp_index": 1,
			},
		},
		bson.M{
			"$group": bson.M{
				"_id":   "$_id",
				"index": bson.M{"$push": "$msgs._search_temp_index"},
			},
		},
		bson.M{"$sort": bson.M{"_id": 1}},
	)
	return mongoutil.Aggregate[searchMessageIndex](ctx, m.coll, pipeline)
}

func (m *MsgMgo) searchMessage(ctx context.Context, req *msg.SearchMessageReq) (int64, []searchMessageIndex, error) {
	filter := bson.M{
		"msgs.msg": bson.M{
			"$exists": true,
			"$type":   "object",
		},
	}
	if req.RecvID != "" {
		filter["$or"] = bson.A{
			bson.M{"msgs.msg.recv_id": req.RecvID},
			bson.M{"msgs.msg.group_id": req.RecvID},
		}
	}
	if req.SendID != "" {
		filter["msgs.msg.send_id"] = req.SendID
	}
	if req.ContentType != 0 {
		filter["msgs.msg.content_type"] = req.ContentType
	}
	if req.SessionType != 0 {
		filter["msgs.msg.session_type"] = req.SessionType
	}
	if req.SendTime != "" {
		sendTime, err := time.Parse(time.DateOnly, req.SendTime)
		if err != nil {
			return 0, nil, errs.ErrArgs.WrapMsg("invalid sendTime", "req", req.SendTime, "format", time.DateOnly, "cause", err.Error())
		}
		filter["$and"] = bson.A{
			bson.M{"msgs.msg.send_time": bson.M{
				"$gte": sendTime.UnixMilli(),
			}},
			bson.M{
				"msgs.msg.send_time": bson.M{
					"$lt": sendTime.Add(time.Hour * 24).UnixMilli(),
				},
			},
		}
	}

	var (
		nextID    primitive.ObjectID
		count     int
		dataRange []searchMessageIndex
		skip      = int((req.Pagination.GetPageNumber() - 1) * req.Pagination.GetShowNumber())
	)
	_, _ = dataRange, skip
	const maxDoc = 50
	data := make([]searchMessageIndex, 0, req.Pagination.GetShowNumber())
	push := cap(data)
	for i := 0; ; i++ {
		res, err := m.searchMessageIndex(ctx, filter, nextID, maxDoc)
		if err != nil {
			return 0, nil, err
		}
		if len(res) > 0 {
			nextID = res[len(res)-1].ID
		}
		for _, r := range res {
			var dataIndex []int64
			for _, index := range r.Index {
				if push > 0 && count >= skip {
					dataIndex = append(dataIndex, index)
					push--
				}
				count++
			}
			if len(dataIndex) > 0 {
				data = append(data, searchMessageIndex{
					ID:    r.ID,
					Index: dataIndex,
				})
			}
		}
		if push <= 0 {
			push--
		}
		if len(res) < maxDoc || push < -10 {
			return int64(count), data, nil
		}
	}
}

func (m *MsgMgo) SearchMessage(ctx context.Context, req *msg.SearchMessageReq) (int64, []*model.MsgInfoModel, error) {
	count, data, err := m.searchMessage(ctx, req)
	if err != nil {
		return 0, nil, err
	}
	var msgs []*model.MsgInfoModel
	if len(data) > 0 {
		var n int
		for _, d := range data {
			n += len(d.Index)
		}
		msgs = make([]*model.MsgInfoModel, 0, n)
	}
	for _, val := range data {
		res, err := mongoutil.FindOne[*model.MsgDocModel](ctx, m.coll, bson.M{"_id": val.ID})
		if err != nil {
			return 0, nil, err
		}
		for _, i := range val.Index {
			if i >= int64(len(res.Msg)) {
				continue
			}
			msgs = append(msgs, res.Msg[i])
		}
	}
	return count, msgs, nil
}

func (m *MsgMgo) RangeUserSendCount(ctx context.Context, start time.Time, end time.Time, group bool, ase bool, pageNumber int32, showNumber int32) (msgCount int64, userCount int64, users []*model.UserCount, dateCount map[string]int64, err error) {
	var sort int
	if ase {
		sort = 1
	} else {
		sort = -1
	}
	type Result struct {
		MsgCount  int64 `bson:"msg_count"`
		UserCount int64 `bson:"user_count"`
		Users     []struct {
			UserID string `bson:"_id"`
			Count  int64  `bson:"count"`
		} `bson:"users"`
		Dates []struct {
			Date  string `bson:"_id"`
			Count int64  `bson:"count"`
		} `bson:"dates"`
	}
	or := bson.A{
		bson.M{
			"doc_id": bson.M{
				"$regex":   "^si_",
				"$options": "i",
			},
		},
	}
	if group {
		or = append(or,
			bson.M{
				"doc_id": bson.M{
					"$regex":   "^g_",
					"$options": "i",
				},
			},
			bson.M{
				"doc_id": bson.M{
					"$regex":   "^sg_",
					"$options": "i",
				},
			},
		)
	}
	pipeline := bson.A{
		bson.M{
			"$match": bson.M{
				"$and": bson.A{
					bson.M{
						"msgs.msg.send_time": bson.M{
							"$gte": start.UnixMilli(),
							"$lt":  end.UnixMilli(),
						},
					},
					bson.M{
						"$or": or,
					},
				},
			},
		},
		bson.M{
			"$addFields": bson.M{
				"msgs": bson.M{
					"$filter": bson.M{
						"input": "$msgs",
						"as":    "item",
						"cond": bson.M{
							"$and": bson.A{
								bson.M{
									"$gte": bson.A{
										"$$item.msg.send_time", start.UnixMilli(),
									},
								},
								bson.M{
									"$lt": bson.A{
										"$$item.msg.send_time", end.UnixMilli(),
									},
								},
							},
						},
					},
				},
			},
		},
		bson.M{
			"$project": bson.M{
				"_id": 0,
			},
		},
		bson.M{
			"$project": bson.M{
				"result": bson.M{
					"$map": bson.M{
						"input": "$msgs",
						"as":    "item",
						"in": bson.M{
							"user_id": "$$item.msg.send_id",
							"send_date": bson.M{
								"$dateToString": bson.M{
									"format": "%Y-%m-%d",
									"date": bson.M{
										"$toDate": "$$item.msg.send_time", // Millisecond timestamp
									},
								},
							},
						},
					},
				},
			},
		},
		bson.M{
			"$unwind": "$result",
		},
		bson.M{
			"$group": bson.M{
				"_id": "$result.send_date",
				"count": bson.M{
					"$sum": 1,
				},
				"original": bson.M{
					"$push": "$$ROOT",
				},
			},
		},
		bson.M{
			"$addFields": bson.M{
				"dates": "$$ROOT",
			},
		},
		bson.M{
			"$project": bson.M{
				"_id":            0,
				"count":          0,
				"dates.original": 0,
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": nil,
				"count": bson.M{
					"$sum": 1,
				},
				"dates": bson.M{
					"$push": "$dates",
				},
				"original": bson.M{
					"$push": "$original",
				},
			},
		},
		bson.M{
			"$unwind": "$original",
		},
		bson.M{
			"$unwind": "$original",
		},
		bson.M{
			"$group": bson.M{
				"_id": "$original.result.user_id",
				"count": bson.M{
					"$sum": 1,
				},
				"original": bson.M{
					"$push": "$dates",
				},
			},
		},
		bson.M{
			"$addFields": bson.M{
				"dates": bson.M{
					"$arrayElemAt": bson.A{"$original", 0},
				},
			},
		},
		bson.M{
			"$project": bson.M{
				"original": 0,
			},
		},
		bson.M{
			"$sort": bson.M{
				"count": sort,
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": nil,
				"user_count": bson.M{
					"$sum": 1,
				},
				"users": bson.M{
					"$push": "$$ROOT",
				},
			},
		},
		bson.M{
			"$addFields": bson.M{
				"dates": bson.M{
					"$arrayElemAt": bson.A{"$users", 0},
				},
			},
		},
		bson.M{
			"$addFields": bson.M{
				"dates": "$dates.dates",
			},
		},
		bson.M{
			"$project": bson.M{
				"_id":         0,
				"users.dates": 0,
			},
		},
		bson.M{
			"$addFields": bson.M{
				"msg_count": bson.M{
					"$sum": "$users.count",
				},
			},
		},
		bson.M{
			"$addFields": bson.M{
				"users": bson.M{
					"$slice": bson.A{"$users", pageNumber - 1, showNumber},
				},
			},
		},
	}
	result, err := mongoutil.Aggregate[*Result](ctx, m.coll, pipeline, options.Aggregate().SetAllowDiskUse(true))
	if err != nil {
		return 0, 0, nil, nil, err
	}
	if len(result) == 0 {
		return 0, 0, nil, nil, errs.Wrap(err)
	}
	users = make([]*model.UserCount, len(result[0].Users))
	for i, r := range result[0].Users {
		users[i] = &model.UserCount{
			UserID: r.UserID,
			Count:  r.Count,
		}
	}
	dateCount = make(map[string]int64)
	for _, r := range result[0].Dates {
		dateCount[r.Date] = r.Count
	}
	return result[0].MsgCount, result[0].UserCount, users, dateCount, nil
}

func (m *MsgMgo) RangeGroupSendCount(ctx context.Context, start time.Time, end time.Time, ase bool, pageNumber int32, showNumber int32) (msgCount int64, userCount int64, groups []*model.GroupCount, dateCount map[string]int64, err error) {
	var sort int
	if ase {
		sort = 1
	} else {
		sort = -1
	}
	type Result struct {
		MsgCount  int64 `bson:"msg_count"`
		UserCount int64 `bson:"user_count"`
		Groups    []struct {
			GroupID string `bson:"_id"`
			Count   int64  `bson:"count"`
		} `bson:"groups"`
		Dates []struct {
			Date  string `bson:"_id"`
			Count int64  `bson:"count"`
		} `bson:"dates"`
	}
	pipeline := bson.A{
		bson.M{
			"$match": bson.M{
				"$and": bson.A{
					bson.M{
						"msgs.msg.send_time": bson.M{
							"$gte": start.UnixMilli(),
							"$lt":  end.UnixMilli(),
						},
					},
					bson.M{
						"$or": bson.A{
							bson.M{
								"doc_id": bson.M{
									"$regex":   "^g_",
									"$options": "i",
								},
							},
							bson.M{
								"doc_id": bson.M{
									"$regex":   "^sg_",
									"$options": "i",
								},
							},
						},
					},
				},
			},
		},
		bson.M{
			"$addFields": bson.M{
				"msgs": bson.M{
					"$filter": bson.M{
						"input": "$msgs",
						"as":    "item",
						"cond": bson.M{
							"$and": bson.A{
								bson.M{
									"$gte": bson.A{
										"$$item.msg.send_time", start.UnixMilli(),
									},
								},
								bson.M{
									"$lt": bson.A{
										"$$item.msg.send_time", end.UnixMilli(),
									},
								},
							},
						},
					},
				},
			},
		},
		bson.M{
			"$project": bson.M{
				"_id": 0,
			},
		},
		bson.M{
			"$project": bson.M{
				"result": bson.M{
					"$map": bson.M{
						"input": "$msgs",
						"as":    "item",
						"in": bson.M{
							"group_id": "$$item.msg.group_id",
							"send_date": bson.M{
								"$dateToString": bson.M{
									"format": "%Y-%m-%d",
									"date": bson.M{
										"$toDate": "$$item.msg.send_time", // Millisecond timestamp
									},
								},
							},
						},
					},
				},
			},
		},
		bson.M{
			"$unwind": "$result",
		},
		bson.M{
			"$group": bson.M{
				"_id": "$result.send_date",
				"count": bson.M{
					"$sum": 1,
				},
				"original": bson.M{
					"$push": "$$ROOT",
				},
			},
		},
		bson.M{
			"$addFields": bson.M{
				"dates": "$$ROOT",
			},
		},
		bson.M{
			"$project": bson.M{
				"_id":            0,
				"count":          0,
				"dates.original": 0,
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": nil,
				"count": bson.M{
					"$sum": 1,
				},
				"dates": bson.M{
					"$push": "$dates",
				},
				"original": bson.M{
					"$push": "$original",
				},
			},
		},
		bson.M{
			"$unwind": "$original",
		},
		bson.M{
			"$unwind": "$original",
		},
		bson.M{
			"$group": bson.M{
				"_id": "$original.result.group_id",
				"count": bson.M{
					"$sum": 1,
				},
				"original": bson.M{
					"$push": "$dates",
				},
			},
		},
		bson.M{
			"$addFields": bson.M{
				"dates": bson.M{
					"$arrayElemAt": bson.A{"$original", 0},
				},
			},
		},
		bson.M{
			"$project": bson.M{
				"original": 0,
			},
		},
		bson.M{
			"$sort": bson.M{
				"count": sort,
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": nil,
				"user_count": bson.M{
					"$sum": 1,
				},
				"groups": bson.M{
					"$push": "$$ROOT",
				},
			},
		},
		bson.M{
			"$addFields": bson.M{
				"dates": bson.M{
					"$arrayElemAt": bson.A{"$groups", 0},
				},
			},
		},
		bson.M{
			"$addFields": bson.M{
				"dates": "$dates.dates",
			},
		},
		bson.M{
			"$project": bson.M{
				"_id":          0,
				"groups.dates": 0,
			},
		},
		bson.M{
			"$addFields": bson.M{
				"msg_count": bson.M{
					"$sum": "$groups.count",
				},
			},
		},
		bson.M{
			"$addFields": bson.M{
				"groups": bson.M{
					"$slice": bson.A{"$groups", pageNumber - 1, showNumber},
				},
			},
		},
	}
	result, err := mongoutil.Aggregate[*Result](ctx, m.coll, pipeline, options.Aggregate().SetAllowDiskUse(true))
	if err != nil {
		return 0, 0, nil, nil, err
	}
	if len(result) == 0 {
		return 0, 0, nil, nil, errs.Wrap(err)
	}
	groups = make([]*model.GroupCount, len(result[0].Groups))
	for i, r := range result[0].Groups {
		groups[i] = &model.GroupCount{
			GroupID: r.GroupID,
			Count:   r.Count,
		}
	}
	dateCount = make(map[string]int64)
	for _, r := range result[0].Dates {
		dateCount[r.Date] = r.Count
	}
	return result[0].MsgCount, result[0].UserCount, groups, dateCount, nil
}

func (m *MsgMgo) GetRandBeforeMsg(ctx context.Context, ts int64, limit int) ([]*model.MsgDocModel, error) {
	return mongoutil.Aggregate[*model.MsgDocModel](ctx, m.coll, []bson.M{
		{
			"$match": bson.M{
				"msgs": bson.M{
					"$not": bson.M{
						"$elemMatch": bson.M{
							"msg.send_time": bson.M{
								"$gt": ts,
							},
						},
					},
				},
			},
		},
		{
			"$project": bson.M{
				"_id":                0,
				"doc_id":             1,
				"msgs.msg.send_time": 1,
				"msgs.msg.seq":       1,
			},
		},
		{
			"$sample": bson.M{
				"size": limit,
			},
		},
	})
}

func (m *MsgMgo) DeleteDoc(ctx context.Context, docID string) error {
	return mongoutil.DeleteOne(ctx, m.coll, bson.M{"doc_id": docID})
}

func (m *MsgMgo) GetLastMessageSeqByTime(ctx context.Context, conversationID string, time int64) (int64, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"doc_id": bson.M{
					"$regex": fmt.Sprintf("^%s", conversationID),
				},
			},
		},
		{
			"$match": bson.M{
				"msgs.msg.send_time": bson.M{
					"$lte": time,
				},
			},
		},
		{
			"$sort": bson.M{
				"_id": -1,
			},
		},
		{
			"$limit": 1,
		},
		{
			"$project": bson.M{
				"_id":                0,
				"doc_id":             1,
				"msgs.msg.send_time": 1,
				"msgs.msg.seq":       1,
			},
		},
	}
	res, err := mongoutil.Aggregate[*model.MsgDocModel](ctx, m.coll, pipeline)
	if err != nil {
		return 0, err
	}
	if len(res) == 0 {
		return 0, nil
	}
	var seq int64
	for _, v := range res[0].Msg {
		if v.Msg == nil {
			continue
		}
		if v.Msg.SendTime <= time {
			seq = v.Msg.Seq
		}
	}
	return seq, nil
}

func (m *MsgMgo) GetLastMessage(ctx context.Context, conversationID string) (*model.MsgInfoModel, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"doc_id": bson.M{
					"$regex": fmt.Sprintf("^%s", conversationID),
				},
			},
		},
		{
			"$match": bson.M{
				"msgs.msg.status": bson.M{
					"$lt": constant.MsgStatusHasDeleted,
				},
			},
		},
		{
			"$sort": bson.M{
				"_id": -1,
			},
		},
		{
			"$limit": 1,
		},
		{
			"$project": bson.M{
				"_id":    0,
				"doc_id": 0,
			},
		},
		{
			"$unwind": "$msgs",
		},
		{
			"$match": bson.M{
				"msgs.msg.status": bson.M{
					"$lt": constant.MsgStatusHasDeleted,
				},
			},
		},
		{
			"$sort": bson.M{
				"msgs.msg.seq": -1,
			},
		},
		{
			"$limit": 1,
		},
	}
	type Result struct {
		Msgs *model.MsgInfoModel `bson:"msgs"`
	}
	res, err := mongoutil.Aggregate[*Result](ctx, m.coll, pipeline)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errs.Wrap(mongo.ErrNoDocuments)
	}
	return res[0].Msgs, nil
}

func (m *MsgMgo) onlyFindDocIndex(ctx context.Context, docID string, indexes []int64) ([]*model.MsgInfoModel, error) {
	if len(indexes) == 0 {
		return nil, nil
	}
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "doc_id", Value: docID},
		}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "doc_id", Value: 1},
			{Key: "msgs", Value: bson.D{
				{Key: "$map", Value: bson.D{
					{Key: "input", Value: indexes},
					{Key: "as", Value: "index"},
					{Key: "in", Value: bson.D{
						{Key: "$arrayElemAt", Value: bson.A{"$msgs", "$$index"}},
					}},
				}},
			}},
		}}},
	}
	msgDocModel, err := mongoutil.Aggregate[*model.MsgDocModel](ctx, m.coll, pipeline)
	if err != nil {
		return nil, err
	}
	if len(msgDocModel) == 0 {
		return nil, nil
	}
	return msgDocModel[0].Msg, nil
}

//func (m *MsgMgo) FindSeqs(ctx context.Context, conversationID string, seqs []int64) ([]*model.MsgInfoModel, error) {
//	if len(seqs) == 0 {
//		return nil, nil
//	}
//	result := make([]*model.MsgInfoModel, 0, len(seqs))
//	for docID, seqs := range m.model.GetDocIDSeqsMap(conversationID, seqs) {
//		res, err := m.onlyFindDocIndex(ctx, docID, datautil.Slice(seqs, m.model.GetMsgIndex))
//		if err != nil {
//			return nil, err
//		}
//		for i, re := range res {
//			if re == nil || re.Msg == nil {
//				continue
//			}
//			result = append(result, res[i])
//		}
//	}
//	return result, nil
//}

func (m *MsgMgo) findBeforeDocSendTime(ctx context.Context, docID string, limit int64) (int64, int64, error) {
	if limit == 0 {
		return 0, 0, nil
	}
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"doc_id": docID,
			},
		},
		{
			"$project": bson.M{
				"_id":    0,
				"doc_id": 0,
			},
		},
		{
			"$unwind": "$msgs",
		},
		{
			"$project": bson.M{
				//"_id":                0,
				//"doc_id":             0,
				"msgs.msg.send_time": 1,
				"msgs.msg.seq":       1,
			},
		},
	}
	if limit > 0 {
		pipeline = append(pipeline, bson.M{"$limit": limit})
	}
	type Result struct {
		Msgs *model.MsgInfoModel `bson:"msgs"`
	}
	res, err := mongoutil.Aggregate[Result](ctx, m.coll, pipeline)
	if err != nil {
		return 0, 0, err
	}
	for i := len(res) - 1; i > 0; i-- {
		v := res[i]
		if v.Msgs != nil && v.Msgs.Msg != nil && v.Msgs.Msg.SendTime > 0 {
			return v.Msgs.Msg.Seq, v.Msgs.Msg.SendTime, nil
		}
	}
	return 0, 0, nil
}

func (m *MsgMgo) findBeforeSendTime(ctx context.Context, conversationID string, seq int64) (int64, int64, error) {
	first := true
	for i := m.model.GetDocIndex(seq); i >= 0; i-- {
		limit := int64(-1)
		if first {
			first = false
			limit = m.model.GetMsgIndex(seq)
		}
		docID := m.model.BuildDocIDByIndex(conversationID, i)
		msgSeq, msgSendTime, err := m.findBeforeDocSendTime(ctx, docID, limit)
		if err != nil {
			return 0, 0, err
		}
		if msgSendTime > 0 {
			return msgSeq, msgSendTime, nil
		}
	}
	return 0, 0, nil
}

func (m *MsgMgo) FindSeqs(ctx context.Context, conversationID string, seqs []int64) ([]*model.MsgInfoModel, error) {
	if len(seqs) == 0 {
		return nil, nil
	}
	var abnormalSeq []int64
	result := make([]*model.MsgInfoModel, 0, len(seqs))
	for docID, docSeqs := range m.model.GetDocIDSeqsMap(conversationID, seqs) {
		res, err := m.onlyFindDocIndex(ctx, docID, datautil.Slice(docSeqs, m.model.GetMsgIndex))
		if err != nil {
			return nil, err
		}
		if len(res) == 0 {
			abnormalSeq = append(abnormalSeq, docSeqs...)
			continue
		}
		for i, re := range res {
			if re == nil || re.Msg == nil || re.Msg.SendTime == 0 {
				abnormalSeq = append(abnormalSeq, docSeqs[i])
				continue
			}
			result = append(result, res[i])
		}
	}
	if len(abnormalSeq) > 0 {
		datautil.Sort(abnormalSeq, false)
		sendTime := make(map[int64]int64)
		var (
			lastSeq      int64
			lastSendTime int64
		)
		for _, seq := range abnormalSeq {
			if lastSendTime > 0 && lastSeq <= seq {
				sendTime[seq] = lastSendTime
				continue
			}
			msgSeq, msgSendTime, err := m.findBeforeSendTime(ctx, conversationID, seq)
			if err != nil {
				return nil, err
			}
			if msgSendTime <= 0 {
				break
			}
			sendTime[seq] = msgSendTime
			lastSeq = msgSeq
			lastSendTime = msgSendTime
		}
		for _, seq := range abnormalSeq {
			result = append(result, &model.MsgInfoModel{
				Msg: &model.MsgDataModel{
					Seq:      seq,
					Status:   constant.MsgStatusHasDeleted,
					SendTime: sendTime[seq],
				},
			})
		}
	}
	return result, nil
}
