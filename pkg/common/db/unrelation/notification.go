package unrelation

import (
	"context"
	"errors"
	"fmt"

	table "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/proto"
)

var ErrNotificationListNotExist = errors.New("user not have msg in mongoDB")
var ErrNotificationNotFound = errors.New("msg not found")

type NotificationMongoDriver struct {
	MsgCollection *mongo.Collection
	msg           table.NotificationDocModel
}

func NewNotificationMongoDriver(database *mongo.Database) table.NotificationDocModelInterface {
	return &NotificationMongoDriver{MsgCollection: database.Collection(table.NotificationDocModel{}.TableName())}
}

func (m *NotificationMongoDriver) PushMsgsToDoc(ctx context.Context, docID string, msgsToMongo []table.NotificationInfoModel) error {
	filter := bson.M{"doc_id": docID}
	return m.MsgCollection.FindOneAndUpdate(ctx, filter, bson.M{"$push": bson.M{"msg": bson.M{"$each": msgsToMongo}}}).Err()
}

func (m *NotificationMongoDriver) Create(ctx context.Context, model *table.NotificationDocModel) error {
	_, err := m.MsgCollection.InsertOne(ctx, model)
	return err
}

func (m *NotificationMongoDriver) UpdateMsgStatusByIndexInOneDoc(ctx context.Context, docID string, msg *sdkws.MsgData, seqIndex int, status int32) error {
	msg.Status = status
	bytes, err := proto.Marshal(msg)
	if err != nil {
		return utils.Wrap(err, "")
	}
	_, err = m.MsgCollection.UpdateOne(ctx, bson.M{"doc_id": docID}, bson.M{"$set": bson.M{fmt.Sprintf("msg.%d.msg", seqIndex): bytes}})
	if err != nil {
		return utils.Wrap(err, "")
	}
	return nil
}

func (m *NotificationMongoDriver) FindOneByDocID(ctx context.Context, docID string) (*table.NotificationDocModel, error) {
	doc := &table.NotificationDocModel{}
	err := m.MsgCollection.FindOne(ctx, bson.M{"doc_id": docID}).Decode(doc)
	return doc, err
}

func (m *NotificationMongoDriver) GetMsgsByIndex(ctx context.Context, conversationID string, index int64) (*table.NotificationDocModel, error) {
	findOpts := options.Find().SetLimit(1).SetSkip(index).SetSort(bson.M{"doc_id": 1})
	cursor, err := m.MsgCollection.Find(ctx, bson.M{"doc_id": primitive.Regex{Pattern: fmt.Sprintf("^%s:", conversationID)}}, findOpts)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var msgs []table.NotificationDocModel
	err = cursor.All(context.Background(), &msgs)
	if err != nil {
		return nil, utils.Wrap(err, fmt.Sprintf("cursor is %s", cursor.Current.String()))
	}
	if len(msgs) > 0 {
		return &msgs[0], nil
	}
	return nil, ErrMsgListNotExist
}

func (m *NotificationMongoDriver) GetNewestMsg(ctx context.Context, conversationID string) (*table.NotificationInfoModel, error) {
	var msgDocs []table.NotificationDocModel
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
		return nil, errors.New("len(msgDocs[0].Msg) < 0")
	}
	return nil, ErrMsgNotFound
}

func (m *NotificationMongoDriver) GetOldestMsg(ctx context.Context, conversationID string) (*table.NotificationInfoModel, error) {
	var msgDocs []table.NotificationDocModel
	cursor, err := m.MsgCollection.Find(ctx, bson.M{"doc_id": bson.M{"$regex": fmt.Sprintf("^%s:", conversationID)}}, options.Find().SetLimit(1).SetSort(bson.M{"doc_id": 1}))
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &msgDocs)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var oldestMsg table.NotificationInfoModel
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

func (m *NotificationMongoDriver) Delete(ctx context.Context, docIDs []string) error {
	if docIDs == nil {
		return nil
	}
	_, err := m.MsgCollection.DeleteMany(ctx, bson.M{"doc_id": bson.M{"$in": docIDs}})
	return err
}

func (m *NotificationMongoDriver) UpdateOneDoc(ctx context.Context, msg *table.NotificationDocModel) error {
	_, err := m.MsgCollection.UpdateOne(ctx, bson.M{"doc_id": msg.DocID}, bson.M{"$set": bson.M{"msg": msg.Msg}})
	return err
}
