package unrelation

import (
	"Open_IM/pkg/common/db/table"
	server_api_params "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ExtendMsgSetMongoDriver struct {
	mgoDB                  *mongo.Database
	ExtendMsgSetCollection *mongo.Collection
}

func NewExtendMsgSetMongoDriver(mgoDB *mongo.Database) *ExtendMsgSetMongoDriver {
	return &ExtendMsgSetMongoDriver{mgoDB: mgoDB, ExtendMsgSetCollection: mgoDB.Collection(table.CExtendMsgSet)}
}

func (e *ExtendMsgSetMongoDriver) CreateExtendMsgSet(ctx context.Context, set *table.ExtendMsgSet) error {
	_, err := e.ExtendMsgSetCollection.InsertOne(ctx, set)
	return err
}

type GetAllExtendMsgSetOpts struct {
	ExcludeExtendMsgs bool
}

func (e *ExtendMsgSetMongoDriver) GetAllExtendMsgSet(ctx context.Context, ID string, opts *GetAllExtendMsgSetOpts) (sets []*table.ExtendMsgSet, err error) {
	regex := fmt.Sprintf("^%s", ID)
	var findOpts *options.FindOptions
	if opts != nil {
		if opts.ExcludeExtendMsgs {
			findOpts = &options.FindOptions{}
			findOpts.SetProjection(bson.M{"extend_msgs": 0})
		}
	}
	cursor, err := e.ExtendMsgSetCollection.Find(ctx, bson.M{"uid": primitive.Regex{Pattern: regex}}, findOpts)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	err = cursor.All(context.Background(), &sets)
	if err != nil {
		return nil, utils.Wrap(err, fmt.Sprintf("cursor is %s", cursor.Current.String()))
	}
	return sets, nil
}

func (e *ExtendMsgSetMongoDriver) GetExtendMsgSet(ctx context.Context, sourceID string, sessionType int32, maxMsgUpdateTime int64) (*table.ExtendMsgSet, error) {
	var err error
	findOpts := options.Find().SetLimit(1).SetSkip(0).SetSort(bson.M{"source_id": -1}).SetProjection(bson.M{"extend_msgs": 0})
	// update newest
	find := bson.M{"source_id": primitive.Regex{Pattern: fmt.Sprintf("^%s", sourceID)}, "session_type": sessionType}
	if maxMsgUpdateTime > 0 {
		find["max_msg_update_time"] = maxMsgUpdateTime
	}
	result, err := e.ExtendMsgSetCollection.Find(ctx, find, findOpts)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var setList []table.ExtendMsgSet
	if err := result.All(ctx, &setList); err != nil {
		return nil, utils.Wrap(err, "")
	}
	if len(setList) == 0 {
		return nil, nil
	}
	return &setList[0], nil
}

// first modify msg
func (e *ExtendMsgSetMongoDriver) InsertExtendMsg(ctx context.Context, sourceID string, sessionType int32, msg *table.ExtendMsg) error {
	set, err := e.GetExtendMsgSet(ctx, sourceID, sessionType, 0)
	if err != nil {
		return utils.Wrap(err, "")
	}
	if set == nil || set.ExtendMsgNum >= set.GetExtendMsgMaxNum() {
		var index int32
		if set != nil {
			index = set.SplitSourceIDAndGetIndex()
		}
		err = e.CreateExtendMsgSet(ctx, &table.ExtendMsgSet{
			SourceID:         set.GetSourceID(sourceID, index),
			SessionType:      sessionType,
			ExtendMsgs:       map[string]table.ExtendMsg{msg.ClientMsgID: *msg},
			ExtendMsgNum:     1,
			CreateTime:       msg.MsgFirstModifyTime,
			MaxMsgUpdateTime: msg.MsgFirstModifyTime,
		})
	} else {
		_, err = e.ExtendMsgSetCollection.UpdateOne(ctx, bson.M{"source_id": set.SourceID, "session_type": sessionType}, bson.M{"$set": bson.M{"max_msg_update_time": msg.MsgFirstModifyTime, "$inc": bson.M{"extend_msg_num": 1}, fmt.Sprintf("extend_msgs.%s", msg.ClientMsgID): msg}})
	}
	return utils.Wrap(err, "")
}

// insert or update
func (e *ExtendMsgSetMongoDriver) InsertOrUpdateReactionExtendMsgSet(ctx context.Context, sourceID string, sessionType int32, clientMsgID string, msgFirstModifyTime int64, reactionExtensionList map[string]*server_api_params.KeyValue) error {
	var updateBson = bson.M{}
	for _, v := range reactionExtensionList {
		updateBson[fmt.Sprintf("extend_msgs.%s.%s", clientMsgID, v.TypeKey)] = v
	}
	upsert := true
	opt := &options.UpdateOptions{
		Upsert: &upsert,
	}
	set, err := e.GetExtendMsgSet(ctx, sourceID, sessionType, msgFirstModifyTime)
	if err != nil {
		return utils.Wrap(err, "")
	}
	if set == nil {
		return errors.New(fmt.Sprintf("sourceID %s has no set", sourceID))
	}
	_, err = e.ExtendMsgSetCollection.UpdateOne(ctx, bson.M{"source_id": set.SourceID, "session_type": sessionType}, bson.M{"$set": updateBson}, opt)
	return utils.Wrap(err, "")
}

// delete TypeKey
func (e *ExtendMsgSetMongoDriver) DeleteReactionExtendMsgSet(ctx context.Context, sourceID string, sessionType int32, clientMsgID string, msgFirstModifyTime int64, reactionExtensionList map[string]*server_api_params.KeyValue) error {
	var updateBson = bson.M{}
	for _, v := range reactionExtensionList {
		updateBson[fmt.Sprintf("extend_msgs.%s.%s", clientMsgID, v.TypeKey)] = ""
	}
	set, err := e.GetExtendMsgSet(ctx, sourceID, sessionType, msgFirstModifyTime)
	if err != nil {
		return utils.Wrap(err, "")
	}
	if set == nil {
		return errors.New(fmt.Sprintf("sourceID %s has no set", sourceID))
	}
	_, err = e.ExtendMsgSetCollection.UpdateOne(ctx, bson.M{"source_id": set.SourceID, "session_type": sessionType}, bson.M{"$unset": updateBson})
	return err
}

func (e *ExtendMsgSetMongoDriver) GetExtendMsg(ctx context.Context, sourceID string, sessionType int32, clientMsgID string, maxMsgUpdateTime int64) (extendMsg *table.ExtendMsg, err error) {
	findOpts := options.Find().SetLimit(1).SetSkip(0).SetSort(bson.M{"source_id": -1}).SetProjection(bson.M{fmt.Sprintf("extend_msgs.%s", clientMsgID): 1})
	regex := fmt.Sprintf("^%s", sourceID)
	result, err := e.ExtendMsgSetCollection.Find(ctx, bson.M{"source_id": primitive.Regex{Pattern: regex}, "session_type": sessionType, "max_msg_update_time": bson.M{"$lte": maxMsgUpdateTime}}, findOpts)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var setList []table.ExtendMsgSet
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
