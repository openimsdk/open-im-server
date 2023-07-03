package unrelation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"

	table "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/proto"
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
	return m.MsgCollection.FindOneAndUpdate(ctx, bson.M{"doc_id": docID}, bson.M{"$push": bson.M{"msgs": bson.M{"$each": msgsToMongo}}}).Err()
}

func (m *MsgMongoDriver) Create(ctx context.Context, model *table.MsgDocModel) error {
	_, err := m.MsgCollection.InsertOne(ctx, model)
	return err
}

func (m *MsgMongoDriver) UpdateMsg(ctx context.Context, docID string, index int64, key string, value any) (*mongo.UpdateResult, error) {
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
func (m *MsgMongoDriver) PushUnique(ctx context.Context, docID string, index int64, key string, value any) (*mongo.UpdateResult, error) {
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
	_, err := m.MsgCollection.UpdateOne(ctx, bson.M{"doc_id": docID}, bson.M{"$set": bson.M{fmt.Sprintf("msgs.%d.msg", index): msg}})
	if err != nil {
		return utils.Wrap(err, "")
	}
	return nil
}

func (m *MsgMongoDriver) UpdateMsgStatusByIndexInOneDoc(ctx context.Context, docID string, msg *sdkws.MsgData, seqIndex int, status int32) error {
	msg.Status = status
	bytes, err := proto.Marshal(msg)
	if err != nil {
		return utils.Wrap(err, "")
	}
	_, err = m.MsgCollection.UpdateOne(ctx, bson.M{"doc_id": docID}, bson.M{"$set": bson.M{fmt.Sprintf("msgs.%d.msg", seqIndex): bytes}})
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

func (m *MsgMongoDriver) GetMsgDocModelByIndex(ctx context.Context, conversationID string, index, sort int64) (*table.MsgDocModel, error) {
	if sort != 1 && sort != -1 {
		return nil, errs.ErrArgs.Wrap("mongo sort must be 1 or -1")
	}
	findOpts := options.Find().SetLimit(1).SetSkip(index).SetSort(bson.M{"doc_id": sort})
	cursor, err := m.MsgCollection.Find(ctx, bson.M{"doc_id": primitive.Regex{Pattern: fmt.Sprintf("^%s:", conversationID)}}, findOpts)
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

func (m *MsgMongoDriver) GetMsgBySeqIndexIn1Doc(ctx context.Context, userID string, docID string, seqs []int64) (msgs []*table.MsgInfoModel, err error) {
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

func (m *MsgMongoDriver) MarkSingleChatMsgsAsRead(ctx context.Context, userID string, docID string, indexes []int64) error {
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

func (m *MsgMongoDriver) RangeUserSendCount(ctx context.Context, start time.Time, end time.Time, ase bool, pageNumber int32, showNumber int32) (msgCount int64, userCount int64, users []*table.UserCount, err error) {
	var sort int
	if ase {
		sort = -1
	} else {
		sort = 1
	}
	type Result struct {
		MsgCount  int64 `bson:"msg_count"`
		UserCount int64 `bson:"user_count"`
		Result    []struct {
			UserID string `bson:"_id"`
			Count  int64  `bson:"count"`
		}
	}
	pipeline := bson.A{
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
			"project": bson.M{
				"_id":    0,
				"doc_id": 0,
			},
		},
		bson.M{
			"$project": bson.M{
				"msgs": bson.M{
					"$map": bson.M{
						"input": "$msgs",
						"as":    "item",
						"in":    "$$item.msg.send_id",
					},
				},
			},
		},
		bson.M{
			"$unwind": "$msgs",
		},
		bson.M{
			"$sortByCount": "$msgs",
		},
		bson.M{
			"sort": bson.M{
				"count": sort,
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": nil,
				"result": bson.M{
					"$push": "$$ROOT",
				},
			},
		},
		bson.M{
			"addFields": bson.M{
				"user_count": bson.M{
					"$size": "$result",
				},
				"msg_count": bson.M{
					"sum": "$result.count",
				},
			},
		},
		bson.M{
			"$addFields": bson.M{
				"result": bson.M{
					"$slice": bson.A{
						"$result", pageNumber - 1, showNumber,
					},
				},
			},
		},
	}
	cur, err := m.MsgCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, 0, nil, errs.Wrap(err)
	}
	defer cur.Close(ctx)
	var result []Result
	if err := cur.All(ctx, &result); err != nil {
		return 0, 0, nil, err
	}
	if len(result) == 0 {
		return 0, 0, nil, nil
	}
	res := make([]*table.UserCount, len(result[0].Result))
	for i, r := range result[0].Result {
		res[i] = &table.UserCount{
			UserID: r.UserID,
			Count:  r.Count,
		}
	}
	return result[0].MsgCount, result[0].UserCount, res, nil
}
