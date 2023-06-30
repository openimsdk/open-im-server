package unrelation

import (
	"context"
	"errors"
	"fmt"

	unRelationTb "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ExtendMsgSetMongoDriver struct {
	mgoDB                  *mongo.Database
	ExtendMsgSetCollection *mongo.Collection
}

func NewExtendMsgSetMongoDriver(mgoDB *mongo.Database) unRelationTb.ExtendMsgSetModelInterface {
	return &ExtendMsgSetMongoDriver{mgoDB: mgoDB, ExtendMsgSetCollection: mgoDB.Collection(unRelationTb.CExtendMsgSet)}
}

func (e *ExtendMsgSetMongoDriver) CreateExtendMsgSet(ctx context.Context, set *unRelationTb.ExtendMsgSetModel) error {
	_, err := e.ExtendMsgSetCollection.InsertOne(ctx, set)
	return err
}

func (e *ExtendMsgSetMongoDriver) GetAllExtendMsgSet(ctx context.Context, ID string, opts *unRelationTb.GetAllExtendMsgSetOpts) (sets []*unRelationTb.ExtendMsgSetModel, err error) {
	regex := fmt.Sprintf("^%s", ID)
	var findOpts *options.FindOptions
	if opts != nil {
		if opts.ExcludeExtendMsgs {
			findOpts = &options.FindOptions{}
			findOpts.SetProjection(bson.M{"extend_msgs": 0})
		}
	}
	cursor, err := e.ExtendMsgSetCollection.Find(ctx, bson.M{"doc_id": primitive.Regex{Pattern: regex}}, findOpts)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	err = cursor.All(ctx, &sets)
	if err != nil {
		return nil, utils.Wrap(err, fmt.Sprintf("cursor is %s", cursor.Current.String()))
	}
	return sets, nil
}

func (e *ExtendMsgSetMongoDriver) GetExtendMsgSet(ctx context.Context, conversationID string, sessionType int32, maxMsgUpdateTime int64) (*unRelationTb.ExtendMsgSetModel, error) {
	var err error
	findOpts := options.Find().SetLimit(1).SetSkip(0).SetSort(bson.M{"source_id": -1}).SetProjection(bson.M{"extend_msgs": 0})
	// update newest
	find := bson.M{"source_id": primitive.Regex{Pattern: fmt.Sprintf("^%s", conversationID)}, "session_type": sessionType}
	if maxMsgUpdateTime > 0 {
		find["max_msg_update_time"] = maxMsgUpdateTime
	}
	result, err := e.ExtendMsgSetCollection.Find(ctx, find, findOpts)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var setList []unRelationTb.ExtendMsgSetModel
	if err := result.All(ctx, &setList); err != nil {
		return nil, utils.Wrap(err, "")
	}
	if len(setList) == 0 {
		return nil, nil
	}
	return &setList[0], nil
}

// first modify msg
func (e *ExtendMsgSetMongoDriver) InsertExtendMsg(ctx context.Context, conversationID string, sessionType int32, msg *unRelationTb.ExtendMsgModel) error {
	set, err := e.GetExtendMsgSet(ctx, conversationID, sessionType, 0)
	if err != nil {
		return utils.Wrap(err, "")
	}
	if set == nil || set.ExtendMsgNum >= set.GetExtendMsgMaxNum() {
		var index int32
		if set != nil {
			index = set.SplitConversationIDAndGetIndex()
		}
		err = e.CreateExtendMsgSet(ctx, &unRelationTb.ExtendMsgSetModel{
			ConversationID:   set.GetConversationID(conversationID, index),
			SessionType:      sessionType,
			ExtendMsgs:       map[string]unRelationTb.ExtendMsgModel{msg.ClientMsgID: *msg},
			ExtendMsgNum:     1,
			CreateTime:       msg.MsgFirstModifyTime,
			MaxMsgUpdateTime: msg.MsgFirstModifyTime,
		})
	} else {
		_, err = e.ExtendMsgSetCollection.UpdateOne(ctx, bson.M{"conversation_id": set.ConversationID, "session_type": sessionType}, bson.M{"$set": bson.M{"max_msg_update_time": msg.MsgFirstModifyTime, "$inc": bson.M{"extend_msg_num": 1}, fmt.Sprintf("extend_msgs.%s", msg.ClientMsgID): msg}})
	}
	return utils.Wrap(err, "")
}

// insert or update
func (e *ExtendMsgSetMongoDriver) InsertOrUpdateReactionExtendMsgSet(ctx context.Context, conversationID string, sessionType int32, clientMsgID string, msgFirstModifyTime int64, reactionExtensionList map[string]*unRelationTb.KeyValueModel) error {
	var updateBson = bson.M{}
	for _, v := range reactionExtensionList {
		updateBson[fmt.Sprintf("extend_msgs.%s.%s", clientMsgID, v.TypeKey)] = v
	}
	upsert := true
	opt := &options.UpdateOptions{
		Upsert: &upsert,
	}
	set, err := e.GetExtendMsgSet(ctx, conversationID, sessionType, msgFirstModifyTime)
	if err != nil {
		return utils.Wrap(err, "")
	}
	if set == nil {
		return errors.New(fmt.Sprintf("conversationID %s has no set", conversationID))
	}
	_, err = e.ExtendMsgSetCollection.UpdateOne(ctx, bson.M{"source_id": set.ConversationID, "session_type": sessionType}, bson.M{"$set": updateBson}, opt)
	return utils.Wrap(err, "")
}

// delete TypeKey
func (e *ExtendMsgSetMongoDriver) DeleteReactionExtendMsgSet(ctx context.Context, conversationID string, sessionType int32, clientMsgID string, msgFirstModifyTime int64, reactionExtensionList map[string]*unRelationTb.KeyValueModel) error {
	var updateBson = bson.M{}
	for _, v := range reactionExtensionList {
		updateBson[fmt.Sprintf("extend_msgs.%s.%s", clientMsgID, v.TypeKey)] = ""
	}
	set, err := e.GetExtendMsgSet(ctx, conversationID, sessionType, msgFirstModifyTime)
	if err != nil {
		return utils.Wrap(err, "")
	}
	if set == nil {
		return errors.New(fmt.Sprintf("conversationID %s has no set", conversationID))
	}
	_, err = e.ExtendMsgSetCollection.UpdateOne(ctx, bson.M{"source_id": set.ConversationID, "session_type": sessionType}, bson.M{"$unset": updateBson})
	return err
}

func (e *ExtendMsgSetMongoDriver) TakeExtendMsg(ctx context.Context, conversationID string, sessionType int32, clientMsgID string, maxMsgUpdateTime int64) (extendMsg *unRelationTb.ExtendMsgModel, err error) {
	findOpts := options.Find().SetLimit(1).SetSkip(0).SetSort(bson.M{"source_id": -1}).SetProjection(bson.M{fmt.Sprintf("extend_msgs.%s", clientMsgID): 1})
	regex := fmt.Sprintf("^%s", conversationID)
	result, err := e.ExtendMsgSetCollection.Find(ctx, bson.M{"source_id": primitive.Regex{Pattern: regex}, "session_type": sessionType, "max_msg_update_time": bson.M{"$lte": maxMsgUpdateTime}}, findOpts)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var setList []unRelationTb.ExtendMsgSetModel
	if err := result.All(ctx, &setList); err != nil {
		return nil, utils.Wrap(err, "")
	}
	if len(setList) == 0 {
		return nil, utils.Wrap(errors.New("GetExtendMsg failed, len(setList) == 0"), "")
	}
	if v, ok := setList[0].ExtendMsgs[clientMsgID]; ok {
		return &v, nil
	}
	return nil, errors.New(fmt.Sprintf("cant find client msg id: %s", clientMsgID))
}
