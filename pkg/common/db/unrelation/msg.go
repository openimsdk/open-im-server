package unrelation

import (
	"context"
	"errors"
	"fmt"
	table "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/golang/protobuf/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	filter := bson.M{"uid": docID}
	return m.MsgCollection.FindOneAndUpdate(ctx, filter, bson.M{"$push": bson.M{"msg": bson.M{"$each": msgsToMongo}}}).Err()
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
	_, err = m.MsgCollection.UpdateOne(ctx, bson.M{"uid": docID}, bson.M{"$set": bson.M{fmt.Sprintf("msg.%d.msg", seqIndex): bytes}})
	if err != nil {
		return utils.Wrap(err, "")
	}
	return nil
}

func (m *MsgMongoDriver) FindOneByDocID(ctx context.Context, docID string) (*table.MsgDocModel, error) {
	doc := &table.MsgDocModel{}
	err := m.MsgCollection.FindOne(ctx, bson.M{"uid": docID}).Decode(doc)
	return doc, err
}

func (m *MsgMongoDriver) GetMsgsByIndex(ctx context.Context, sourceID string, index int64) (*table.MsgDocModel, error) {
	findOpts := options.Find().SetLimit(1).SetSkip(index).SetSort(bson.M{"uid": 1})
	cursor, err := m.MsgCollection.Find(ctx, bson.M{"uid": primitive.Regex{Pattern: fmt.Sprintf("^%s:", sourceID)}}, findOpts)
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

func (m *MsgMongoDriver) GetNewestMsg(ctx context.Context, sourceID string) (*table.MsgInfoModel, error) {
	var msgDocs []table.MsgDocModel
	cursor, err := m.MsgCollection.Find(ctx, bson.M{"uid": bson.M{"$regex": fmt.Sprintf("^%s:", sourceID)}}, options.Find().SetLimit(1).SetSort(bson.M{"uid": -1}))
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
		return nil, errors.New("len(msgDocs[0].Msg) < 0")
	}
	return nil, ErrMsgNotFound
}

func (m *MsgMongoDriver) GetOldestMsg(ctx context.Context, sourceID string) (*table.MsgInfoModel, error) {
	var msgDocs []table.MsgDocModel
	cursor, err := m.MsgCollection.Find(ctx, bson.M{"uid": bson.M{"$regex": fmt.Sprintf("^%s:", sourceID)}}, options.Find().SetLimit(1).SetSort(bson.M{"uid": 1}))
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
	_, err := m.MsgCollection.DeleteMany(ctx, bson.M{"uid": bson.M{"$in": docIDs}})
	return err
}

func (m *MsgMongoDriver) UpdateOneDoc(ctx context.Context, msg *table.MsgDocModel) error {
	_, err := m.MsgCollection.UpdateOne(ctx, bson.M{"uid": msg.DocID}, bson.M{"$set": bson.M{"msg": msg.Msg}})
	return err
}
