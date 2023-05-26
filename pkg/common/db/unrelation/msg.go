package unrelation

import (
	"context"
	"errors"
	"fmt"

	table "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
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

func (m *MsgMongoDriver) GetMsgBySeqIndexIn1Doc(ctx context.Context, docID string, seqs []int64) (msgs []*table.MsgInfoModel, err error) {
	beginSeq, endSeq := utils.GetSeqsBeginEnd(seqs)
	beginIndex := m.model.GetMsgIndex(beginSeq)
	num := endSeq - beginSeq + 1
	pipeline := bson.A{
		bson.M{
			"$match": bson.M{"doc_id": docID},
		},
		bson.M{
			"$project": bson.M{
				"msgs": bson.M{
					"$slice": bson.A{"$msgs", beginIndex, num},
				},
			},
		},
	}
	cursor, err := m.MsgCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	defer cursor.Close(ctx)
	var doc table.MsgDocModel
	i := 0
	for cursor.Next(ctx) {
		err := cursor.Decode(&doc)
		if err != nil {
			return nil, err
		}
		if i == 0 {
			break
		}
	}
	log.ZDebug(ctx, "msgInfos", "num", len(doc.Msg), "docID", docID)
	for _, v := range doc.Msg {
		if v.Msg.Seq >= beginSeq && v.Msg.Seq <= endSeq {
			log.ZDebug(ctx, "find msg", "msg", v.Msg)
			msgs = append(msgs, v)
		} else {
			log.ZWarn(ctx, "this msg is at wrong position", nil, "msg", v.Msg)
		}
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
