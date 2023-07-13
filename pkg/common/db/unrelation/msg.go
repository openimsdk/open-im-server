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

package unrelation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/proto"

	table "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
)

var ErrMsgListNotExist = errors.New("user not have msg in mongoDB")

type MsgMongoDriver struct {
	MsgCollection *mongo.Collection
	model         table.MsgDocModel
}

func NewMsgMongoDriver(database *mongo.Database) table.MsgDocModelInterface {
	collection := database.Collection(table.MsgDocModel{}.TableName())
	return &MsgMongoDriver{MsgCollection: collection}
}

func (m *MsgMongoDriver) PushMsgsToDoc(ctx context.Context, docID string, msgsToMongo []table.MsgInfoModel) error {
	return m.MsgCollection.FindOneAndUpdate(ctx, bson.M{"doc_id": docID}, bson.M{"$push": bson.M{"msgs": bson.M{"$each": msgsToMongo}}}).
		Err()
}

func (m *MsgMongoDriver) Create(ctx context.Context, model *table.MsgDocModel) error {
	_, err := m.MsgCollection.InsertOne(ctx, model)
	return err
}

func (m *MsgMongoDriver) UpdateMsg(
	ctx context.Context,
	docID string,
	index int64,
	key string,
	value any,
) (*mongo.UpdateResult, error) {
	var field string
	if key == "" {
		field = fmt.Sprintf("msgs.%d", index)
	} else {
		field = fmt.Sprintf("msgs.%d.%s", index, key)
	}
	filter := bson.M{"doc_id": docID}
	update := bson.M{"$set": bson.M{field: value}}
	res, err := m.MsgCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	return res, nil
}

// PushUnique value must slice
func (m *MsgMongoDriver) PushUnique(
	ctx context.Context,
	docID string,
	index int64,
	key string,
	value any,
) (*mongo.UpdateResult, error) {
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
	res, err := m.MsgCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	return res, nil
}

func (m *MsgMongoDriver) UpdateMsgContent(ctx context.Context, docID string, index int64, msg []byte) error {
	_, err := m.MsgCollection.UpdateOne(
		ctx,
		bson.M{"doc_id": docID},
		bson.M{"$set": bson.M{fmt.Sprintf("msgs.%d.msg", index): msg}},
	)
	if err != nil {
		return utils.Wrap(err, "")
	}
	return nil
}

func (m *MsgMongoDriver) UpdateMsgStatusByIndexInOneDoc(
	ctx context.Context,
	docID string,
	msg *sdkws.MsgData,
	seqIndex int,
	status int32,
) error {
	msg.Status = status
	bytes, err := proto.Marshal(msg)
	if err != nil {
		return utils.Wrap(err, "")
	}
	_, err = m.MsgCollection.UpdateOne(
		ctx,
		bson.M{"doc_id": docID},
		bson.M{"$set": bson.M{fmt.Sprintf("msgs.%d.msg", seqIndex): bytes}},
	)
	if err != nil {
		return utils.Wrap(err, "")
	}
	return nil
}

func (m *MsgMongoDriver) FindOneByDocID(ctx context.Context, docID string) (*table.MsgDocModel, error) {
	doc := &table.MsgDocModel{}
	err := m.MsgCollection.FindOne(ctx, bson.M{"doc_id": docID}).Decode(doc)
	return doc, err
}

func (m *MsgMongoDriver) GetMsgDocModelByIndex(
	ctx context.Context,
	conversationID string,
	index, sort int64,
) (*table.MsgDocModel, error) {
	if sort != 1 && sort != -1 {
		return nil, errs.ErrArgs.Wrap("mongo sort must be 1 or -1")
	}
	findOpts := options.Find().SetLimit(1).SetSkip(index).SetSort(bson.M{"doc_id": sort})
	cursor, err := m.MsgCollection.Find(
		ctx,
		bson.M{"doc_id": primitive.Regex{Pattern: fmt.Sprintf("^%s:", conversationID)}},
		findOpts,
	)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var msgs []table.MsgDocModel
	err = cursor.All(ctx, &msgs)
	if err != nil {
		return nil, utils.Wrap(err, fmt.Sprintf("cursor is %s", cursor.Current.String()))
	}
	if len(msgs) > 0 {
		return &msgs[0], nil
	}
	return nil, ErrMsgListNotExist
}

func (m *MsgMongoDriver) GetNewestMsg(ctx context.Context, conversationID string) (*table.MsgInfoModel, error) {
	var skip int64 = 0
	for {
		msgDocModel, err := m.GetMsgDocModelByIndex(ctx, conversationID, skip, -1)
		if err != nil {
			return nil, err
		}
		for i := len(msgDocModel.Msg) - 1; i >= 0; i-- {
			if msgDocModel.Msg[i].Msg != nil {
				return msgDocModel.Msg[i], nil
			}
		}
		skip++
	}
}

func (m *MsgMongoDriver) GetOldestMsg(ctx context.Context, conversationID string) (*table.MsgInfoModel, error) {
	var skip int64 = 0
	for {
		msgDocModel, err := m.GetMsgDocModelByIndex(ctx, conversationID, skip, 1)
		if err != nil {
			return nil, err
		}
		for i, v := range msgDocModel.Msg {
			if v.Msg != nil {
				return msgDocModel.Msg[i], nil
			}
		}
		skip++
	}
}

func (m *MsgMongoDriver) DeleteMsgsInOneDocByIndex(ctx context.Context, docID string, indexes []int) error {
	updates := bson.M{
		"$set": bson.M{},
	}
	for _, index := range indexes {
		updates["$set"].(bson.M)[fmt.Sprintf("msgs.%d", index)] = bson.M{
			"msg": nil,
		}
	}
	_, err := m.MsgCollection.UpdateMany(ctx, bson.M{"doc_id": docID}, updates)
	if err != nil {
		return utils.Wrap(err, "")
	}
	return nil
}

func (m *MsgMongoDriver) DeleteDocs(ctx context.Context, docIDs []string) error {
	if docIDs == nil {
		return nil
	}
	_, err := m.MsgCollection.DeleteMany(ctx, bson.M{"doc_id": bson.M{"$in": docIDs}})
	return err
}

func (m *MsgMongoDriver) GetMsgBySeqIndexIn1Doc(
	ctx context.Context,
	userID string,
	docID string,
	seqs []int64,
) (msgs []*table.MsgInfoModel, err error) {
	indexs := make([]int64, 0, len(seqs))
	for _, seq := range seqs {
		indexs = append(indexs, m.model.GetMsgIndex(seq))
	}
	pipeline := mongo.Pipeline{
		{
			{"$match", bson.D{
				{"doc_id", docID},
			}},
		},
		{
			{"$project", bson.D{
				{"_id", 0},
				{"doc_id", 1},
				{"msgs", bson.D{
					{"$map", bson.D{
						{"input", indexs},
						{"as", "index"},
						{"in", bson.D{
							{"$let", bson.D{
								{"vars", bson.D{
									{"currentMsg", bson.D{
										{"$arrayElemAt", []string{"$msgs", "$$index"}},
									}},
								}},
								{"in", bson.D{
									{"$cond", bson.D{
										{"if", bson.D{
											{"$in", []string{userID, "$$currentMsg.del_list"}},
										}},
										{"then", nil},
										{"else", "$$currentMsg"},
									}},
								}},
							}},
						}},
					}},
				}},
			}},
		},
		{
			{"$project", bson.D{
				{"msgs.del_list", 0},
			}},
		},
	}
	cur, err := m.MsgCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	defer cur.Close(ctx)
	var msgDocModel []table.MsgDocModel
	if err := cur.All(ctx, &msgDocModel); err != nil {
		return nil, errs.Wrap(err)
	}
	if len(msgDocModel) == 0 {
		return nil, errs.Wrap(mongo.ErrNoDocuments)
	}
	msgs = make([]*table.MsgInfoModel, 0, len(msgDocModel[0].Msg))
	for i := range msgDocModel[0].Msg {
		msg := msgDocModel[0].Msg[i]
		if msg == nil || msg.Msg == nil {
			continue
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
			data, err := json.Marshal(&revokeContent)
			if err != nil {
				return nil, err
			}
			elem := sdkws.NotificationElem{
				Detail: string(data),
			}
			content, err := json.Marshal(&elem)
			if err != nil {
				return nil, err
			}
			msg.Msg.ContentType = constant.MsgRevokeNotification
			msg.Msg.Content = string(content)
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

func (m *MsgMongoDriver) IsExistDocID(ctx context.Context, docID string) (bool, error) {
	count, err := m.MsgCollection.CountDocuments(ctx, bson.M{"doc_id": docID})
	if err != nil {
		return false, errs.Wrap(err)
	}
	return count > 0, nil
}

func (m *MsgMongoDriver) MarkSingleChatMsgsAsRead(
	ctx context.Context,
	userID string,
	docID string,
	indexes []int64,
) error {
	updates := []mongo.WriteModel{}
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
	_, err := m.MsgCollection.BulkWrite(ctx, updates)
	return err
}

// RangeUserSendCount
// db.msg.aggregate([
//
//	{
//	    $match: {
//	        "msgs.msg.send_time": {
//	            "$gte": 0,
//	            "$lt": 1788122092317
//	        }
//	    }
//	},
//	{
//	    "$addFields": {
//	        "msgs": {
//	            "$filter": {
//	                "input": "$msgs",
//	                "as": "item",
//	                "cond": {
//	                    "$and": [
//	                        {
//	                            $gte: ["$$item.msg.send_time", 0]
//	                        },
//	                        {
//	                            $lt: ["$$item.msg.send_time", 1788122092317]
//	                        }
//	                    ]
//	                }
//	            }
//	        }
//	    }
//	},
//	{
//	    "$project": {
//	        "_id": 0,
//
//	    },
//
//	},
//	{
//	    "$project": {
//	        "result": {
//	            "$map": {
//	                "input": "$msgs",
//	                "as": "item",
//	                "in": {
//	                    user_id: "$$item.msg.send_id",
//	                    send_date: {
//	                        $dateToString: {
//	                            format: "%Y-%m-%d",
//	                            date: {
//	                                $toDate: "$$item.msg.send_time"
//	                            }
//	                        }
//	                    }
//	                }
//	            }
//	        }
//	    },
//
//	},
//	{
//	    "$unwind": "$result"
//	},
//	{
//	    "$group": {
//	        _id: "$result.send_date",
//	        count: {
//	            $sum: 1
//	        },
//	        original: {
//	            $push: "$$ROOT"
//	        }
//	    }
//	},
//	{
//	    "$addFields": {
//	        "dates": "$$ROOT"
//	    }
//	},
//	{
//	    "$project": {
//	        "_id": 0,
//	        "count": 0,
//	        "dates.original": 0,
//
//	    },
//
//	},
//	{
//	    "$group": {
//	        _id: null,
//	        count: {
//	            $sum: 1
//	        },
//	        dates: {
//	            $push: "$dates"
//	        },
//	        original: {
//	            $push: "$original"
//	        },
//
//	    }
//	},
//	{
//	    "$unwind": "$original"
//	},
//	{
//	    "$unwind": "$original"
//	},
//	{
//	    "$group": {
//	        _id: "$original.result.user_id",
//	        count: {
//	            $sum: 1
//	        },
//	        original: {
//	            $push: "$dates"
//	        },
//
//	    }
//	},
//	{
//	    "$addFields": {
//	        "dates": {
//	            $arrayElemAt: ["$original", 0]
//	        }
//	    }
//	},
//	{
//	    "$project": {
//	        original: 0
//	    }
//	},
//	{
//	    $sort: {
//	        count: - 1
//	    }
//	},
//	{
//	    "$group": {
//	        _id: null,
//	        user_count: {
//	            $sum: 1
//	        },
//	        users: {
//	            $push: "$$ROOT"
//	        },
//
//	    }
//	},
//	{
//	    "$addFields": {
//	        "dates": {
//	            $arrayElemAt: ["$users", 0]
//	        }
//	    }
//	},
//	{
//	    "$addFields": {
//	        "dates": "$dates.dates"
//	    }
//	},
//	{
//	    "$project": {
//	        _id: 0,
//	        "users.dates": 0,
//
//	    }
//	},
//	{
//	    "$addFields": {
//	        "msg_count": {
//	            $sum: "$users.count"
//	        }
//	    }
//	},
//	{
//	    "$addFields": {
//	        users: {
//	            $slice: ["$users", 0, 10]
//	        }
//	    }
//	}
//
// ])
func (m *MsgMongoDriver) RangeUserSendCount(
	ctx context.Context,
	start time.Time,
	end time.Time,
	group bool,
	ase bool,
	pageNumber int32,
	showNumber int32,
) (msgCount int64, userCount int64, users []*table.UserCount, dateCount map[string]int64, err error) {
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
										"$toDate": "$$item.msg.send_time", // 毫秒时间戳
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
	cur, err := m.MsgCollection.Aggregate(ctx, pipeline, options.Aggregate().SetAllowDiskUse(true))
	if err != nil {
		return 0, 0, nil, nil, errs.Wrap(err)
	}
	defer cur.Close(ctx)
	var result []Result
	if err := cur.All(ctx, &result); err != nil {
		return 0, 0, nil, nil, errs.Wrap(err)
	}
	if len(result) == 0 {
		return 0, 0, nil, nil, errs.Wrap(err)
	}
	users = make([]*table.UserCount, len(result[0].Users))
	for i, r := range result[0].Users {
		users[i] = &table.UserCount{
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

func (m *MsgMongoDriver) RangeGroupSendCount(
	ctx context.Context,
	start time.Time,
	end time.Time,
	ase bool,
	pageNumber int32,
	showNumber int32,
) (msgCount int64, userCount int64, groups []*table.GroupCount, dateCount map[string]int64, err error) {
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
										"$toDate": "$$item.msg.send_time", // 毫秒时间戳
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
	cur, err := m.MsgCollection.Aggregate(ctx, pipeline, options.Aggregate().SetAllowDiskUse(true))
	if err != nil {
		return 0, 0, nil, nil, errs.Wrap(err)
	}
	defer cur.Close(ctx)
	var result []Result
	if err := cur.All(ctx, &result); err != nil {
		return 0, 0, nil, nil, errs.Wrap(err)
	}
	if len(result) == 0 {
		return 0, 0, nil, nil, errs.Wrap(err)
	}
	groups = make([]*table.GroupCount, len(result[0].Groups))
	for i, r := range result[0].Groups {
		groups[i] = &table.GroupCount{
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

func (m *MsgMongoDriver) SearchMessage(ctx context.Context, req *msg.SearchMessageReq) ([]*table.MsgInfoModel, error) {
	msgs, err := m.searchMessage(ctx, req)
	if err != nil {
		return nil, err
	}
	for _, msg1 := range msgs {
		if msg1.IsRead {
			msg1.Msg.IsRead = true
		}
	}
	return msgs, nil
}

func (m *MsgMongoDriver) searchMessage(ctx context.Context, req *msg.SearchMessageReq) ([]*table.MsgInfoModel, error) {
	var pipe mongo.Pipeline
	conditon := bson.A{}
	if req.SendTime != "" {
		conditon = append(conditon, bson.M{"$eq": bson.A{bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": bson.M{"$toDate": "$$item.msg.send_time"}}}, req.SendTime}})
	}
	if req.MsgType != 0 {
		conditon = append(conditon, bson.M{"$eq": bson.A{"$$item.msg.content_type", req.MsgType}})
	}
	if req.SessionType != 0 {
		conditon = append(conditon, bson.M{"$eq": bson.A{"$$item.msg.session_type", req.SessionType}})
	}
	if req.RecvID != "" {
		conditon = append(conditon, bson.M{"$regexFind": bson.M{"input": "$$item.msg.recv_id", "regex": req.RecvID}})
	}
	if req.SendID != "" {
		conditon = append(conditon, bson.M{"$regexFind": bson.M{"input": "$$item.msg.send_id", "regex": req.SendID}})
	}

	or := bson.A{
		bson.M{
			"doc_id": bson.M{
				"$regex":   "^si_",
				"$options": "i",
			},
		},
	}
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

	pipe = mongo.Pipeline{
		{
			{"$match", bson.D{
				{
					"$or", or,
				},
			}},
		},
		{
			{"$project", bson.D{
				{"msgs", bson.D{
					{"$filter", bson.D{
						{"input", "$msgs"},
						{"as", "item"},
						{"cond", bson.D{
							{"$and", conditon},
						},
						}},
					}},
				},
				{"doc_id", 1},
			}},
		},
	}
	cursor, err := m.MsgCollection.Aggregate(ctx, pipe)
	if err != nil {
		return nil, err
	}

	var msgsDocs []table.MsgDocModel
	err = cursor.All(ctx, &msgsDocs)
	if err != nil {
		return nil, err
	}
	if len(msgsDocs) == 0 {
		return nil, errs.Wrap(mongo.ErrNoDocuments)
	}
	msgs := make([]*table.MsgInfoModel, 0)
	for index, _ := range msgsDocs {
		for i := range msgsDocs[index].Msg {
			msg := msgsDocs[index].Msg[i]
			if msg == nil || msg.Msg == nil {
				continue
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
				data, err := json.Marshal(&revokeContent)
				if err != nil {
					return nil, err
				}
				elem := sdkws.NotificationElem{
					Detail: string(data),
				}
				content, err := json.Marshal(&elem)
				if err != nil {
					return nil, err
				}
				msg.Msg.ContentType = constant.MsgRevokeNotification
				msg.Msg.Content = string(content)
			}
			msgs = append(msgs, msg)
		}
	}
	return msgs, nil
}
