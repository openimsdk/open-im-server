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
var ErrMsgNotFound = errors.New("msg not found")

type MsgMongoDriver struct {
	MsgCollection *mongo.Collection
	msg           table.MsgDocModel
}

func NewMsgMongoDriver(database *mongo.Database) table.MsgDocModelInterface {
	return &MsgMongoDriver{MsgCollection: database.Collection(table.MsgDocModel{}.TableName())}
}

func (m *MsgMongoDriver) PushMsgsToDoc(ctx context.Context, docID string, msgsToMongo []table.MsgInfoModel) error {
	return m.MsgCollection.FindOneAndUpdate(ctx, bson.M{"doc_id": docID}, bson.M{"$push": bson.M{"msgs": bson.M{"$each": msgsToMongo}}}).Err()
}

func (m *MsgMongoDriver) Create(ctx context.Context, model *table.MsgDocModel) error {
	_, err := m.MsgCollection.InsertOne(ctx, model)
	return err
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

func (m *MsgMongoDriver) GetMsgsByIndex(ctx context.Context, conversationID string, index int64) (*table.MsgDocModel, error) {
	findOpts := options.Find().SetLimit(1).SetSkip(index).SetSort(bson.M{"doc_id": 1})
	cursor, err := m.MsgCollection.Find(ctx, bson.M{"doc_id": primitive.Regex{Pattern: fmt.Sprintf("^%s:", conversationID)}}, findOpts)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var msgs []table.MsgDocModel
	err = cursor.All(context.Background(), &msgs)
	if err != nil {
		return nil, utils.Wrap(err, fmt.Sprintf("cursor is %s", cursor.Current.String()))
	}
	if len(msgs) > 0 {
		return &msgs[0], nil
	}
	return nil, ErrMsgListNotExist
}

func (m *MsgMongoDriver) GetNewestMsg(ctx context.Context, conversationID string) (*table.MsgInfoModel, error) {
	var msgDocs []table.MsgDocModel
	cursor, err := m.MsgCollection.Find(ctx, bson.M{"doc_id": bson.M{"$regex": fmt.Sprintf("^%s:", conversationID)}}, options.Find().SetLimit(1).SetSort(bson.M{"doc_id": -1}))
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	err = cursor.All(ctx, &msgDocs)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	if len(msgDocs) > 0 {
		if len(msgDocs[0].Msg) > 0 {
			return &msgDocs[0].Msg[len(msgDocs[0].Msg)-1], nil
		}
		return nil, errs.ErrRecordNotFound.Wrap("len(msgDocs[0].Msgs) < 0")
	}
	return nil, ErrMsgNotFound
}

func (m *MsgMongoDriver) GetOldestMsg(ctx context.Context, conversationID string) (*table.MsgInfoModel, error) {
	var msgDocs []table.MsgDocModel
	cursor, err := m.MsgCollection.Find(ctx, bson.M{"doc_id": bson.M{"$regex": fmt.Sprintf("^%s:", conversationID)}}, options.Find().SetLimit(1).SetSort(bson.M{"doc_id": 1}))
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &msgDocs)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var oldestMsg table.MsgInfoModel
	if len(msgDocs) > 0 {
		for _, v := range msgDocs[0].Msg {
			if v.SendTime != 0 {
				oldestMsg = v
				break
			}
		}
		if len(oldestMsg.Msg) == 0 {
			if len(msgDocs[0].Msg) > 0 {
				oldestMsg = msgDocs[0].Msg[0]
			}
		}
		return &oldestMsg, nil
	}
	return nil, ErrMsgNotFound
}

func (m *MsgMongoDriver) Delete(ctx context.Context, docIDs []string) error {
	if docIDs == nil {
		return nil
	}
	_, err := m.MsgCollection.DeleteMany(ctx, bson.M{"doc_id": bson.M{"$in": docIDs}})
	return err
}

func (m *MsgMongoDriver) UpdateOneDoc(ctx context.Context, msg *table.MsgDocModel) error {
	_, err := m.MsgCollection.UpdateOne(ctx, bson.M{"doc_id": msg.DocID}, bson.M{"$set": bson.M{"msgs": msg.Msg}})
	return err
}

func (m *MsgMongoDriver) GetMsgBySeqIndexIn1Doc(ctx context.Context, docID string, beginSeq, endSeq int64) (msgs []*sdkws.MsgData, seqs []int64, err error) {
	beginIndex := m.msg.GetMsgIndex(beginSeq)
	num := endSeq - beginSeq + 1
	log.ZDebug(ctx, "info", "beginIndex", beginIndex, "num", num)
	pipeline := bson.A{
		bson.M{
			"$match": bson.M{"doc_id": docID},
		},
		bson.M{
			"$project": bson.M{
				"msgs": bson.M{
					"$slice": []interface{}{"$msgs", beginIndex, num},
				},
			},
		},
	}
	cursor, err := m.MsgCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, nil, errs.Wrap(err)
	}
	var msgInfos []table.MsgInfoModel
	if err := cursor.All(ctx, &msgInfos); err != nil {
		return nil, nil, err
	}
	if len(msgInfos) < 1 {
		return nil, nil, errs.ErrRecordNotFound.Wrap("mongo GetMsgBySeqIndex failed, len is 0")
	}
	log.ZDebug(ctx, "msgInfos", "num", len(msgInfos))
	for _, v := range msgInfos {
		var msg sdkws.MsgData
		if err := proto.Unmarshal(v.Msg, &msg); err != nil {
			return nil, nil, err
		}
		if msg.Seq >= beginSeq && msg.Seq <= endSeq {
			log.ZDebug(ctx, "find msg", "msg", &msg)
			msgs = append(msgs, &msg)
			seqs = append(seqs, msg.Seq)
		} else {
			log.ZWarn(ctx, "this msg is at wrong position", nil, "msg", &msg)
		}
	}
	return msgs, seqs, nil
}
