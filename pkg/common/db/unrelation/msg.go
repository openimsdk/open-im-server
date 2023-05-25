package unrelation

import (
	"context"
	"errors"
	"fmt"

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
var ErrMsgNotFound = errors.New("msg not found")

type MsgMongoDriver struct {
	MsgCollection *mongo.Collection
	msg           table.MsgDocModel
}

func NewMsgMongoDriver(database *mongo.Database) table.MsgDocModelInterface {
	collection := database.Collection(table.MsgDocModel{}.TableName())
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"doc_id": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		panic(err)
	}
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

func (m *MsgMongoDriver) GetMsgAndIndexBySeqsInOneDoc(ctx context.Context, docID string, seqs []int64) (seqMsgs []*sdkws.MsgData, indexes []int, unExistSeqs []int64, err error) {
	//doc, err := m.FindOneByDocID(ctx, docID)
	//if err != nil {
	//	return nil, nil, nil, err
	//}
	//singleCount := 0
	//var hasSeqList []int64
	//for i := 0; i < len(doc.Msg); i++ {
	//	var msg sdkws.MsgData
	//	if err := proto.Unmarshal(doc.Msg[i].Msg, &msg); err != nil {
	//		return nil, nil, nil, err
	//	}
	//	if utils.Contain(msg.Seq, seqs...) {
	//		indexes = append(indexes, i)
	//		seqMsgs = append(seqMsgs, &msg)
	//		hasSeqList = append(hasSeqList, msg.Seq)
	//		singleCount++
	//		if singleCount == len(seqs) {
	//			break
	//		}
	//	}
	//}
	//for _, i := range seqs {
	//	if utils.Contain(i, hasSeqList...) {
	//		continue
	//	}
	//	unExistSeqs = append(unExistSeqs, i)
	//}
	return seqMsgs, indexes, unExistSeqs, nil
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
	//var msgDocs []table.MsgDocModel
	//cursor, err := m.MsgCollection.Find(ctx, bson.M{"doc_id": bson.M{"$regex": fmt.Sprintf("^%s:", conversationID)}}, options.Find().SetLimit(1).SetSort(bson.M{"doc_id": 1}))
	//if err != nil {
	//	return nil, err
	//}
	//err = cursor.All(ctx, &msgDocs)
	//if err != nil {
	//	return nil, utils.Wrap(err, "")
	//}
	//var oldestMsg table.MsgInfoModel
	//if len(msgDocs) > 0 {
	//	for _, v := range msgDocs[0].Msg {
	//		if v.SendTime != 0 {
	//			oldestMsg = v
	//			break
	//		}
	//	}
	//	if len(oldestMsg.Msg) == 0 {
	//		if len(msgDocs[0].Msg) > 0 {
	//			oldestMsg = msgDocs[0].Msg[0]
	//		}
	//	}
	//	return &oldestMsg, nil
	//}
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

func (m *MsgMongoDriver) GetMsgBySeqIndexIn1Doc(ctx context.Context, docID string, seqs []int64) (msgs []*sdkws.MsgData, err error) {
	//beginSeq, endSeq := utils.GetSeqsBeginEnd(seqs)
	//beginIndex := m.msg.GetMsgIndex(beginSeq)
	//num := endSeq - beginSeq + 1
	//pipeline := bson.A{
	//	bson.M{
	//		"$match": bson.M{"doc_id": docID},
	//	},
	//	bson.M{
	//		"$project": bson.M{
	//			"msgs": bson.M{
	//				"$slice": bson.A{"$msgs", beginIndex, num},
	//			},
	//		},
	//	},
	//}
	//cursor, err := m.MsgCollection.Aggregate(ctx, pipeline)
	//if err != nil {
	//	return nil, errs.Wrap(err)
	//}
	//defer cursor.Close(ctx)
	//var doc table.MsgDocModel
	//i := 0
	//for cursor.Next(ctx) {
	//	err := cursor.Decode(&doc)
	//	if err != nil {
	//		return nil, err
	//	}
	//	if i == 0 {
	//		break
	//	}
	//}
	//log.ZDebug(ctx, "msgInfos", "num", len(doc.Msg), "docID", docID)
	//for _, v := range doc.Msg {
	//	var msg sdkws.MsgData
	//	if err := proto.Unmarshal(v.Msg, &msg); err != nil {
	//		return nil, err
	//	}
	//	if msg.Seq >= beginSeq && msg.Seq <= endSeq {
	//		log.ZDebug(ctx, "find msg", "msg", &msg)
	//		msgs = append(msgs, &msg)
	//	} else {
	//		log.ZWarn(ctx, "this msg is at wrong position", nil, "msg", &msg)
	//	}
	//}
	return msgs, nil
}

func (m *MsgMongoDriver) IsExistDocID(ctx context.Context, docID string) (bool, error) {
	count, err := m.MsgCollection.CountDocuments(ctx, bson.M{"doc_id": docID})
	if err != nil {
		return false, errs.Wrap(err)
	}
	return count > 0, nil
}
