package db

import (
	"Open_IM/pkg/common/config"
	server_api_params "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

const cExtendMsgSet = "extend_msgs"
const MaxNum = 100

type ExtendMsgSet struct {
	SourceID         string               `bson:"source_id" json:"sourceID"`
	SessionType      int32                `bson:"session_type" json:"sessionType"`
	ExtendMsgs       map[string]ExtendMsg `bson:"extend_msgs" json:"extendMsgs"`
	ExtendMsgNum     int32                `bson:"extend_msg_num" json:"extendMsgNum"`
	CreateTime       int64                `bson:"create_time" json:"createTime"`               // this block's create time
	MaxMsgUpdateTime int64                `bson:"max_msg_update_time" json:"maxMsgUpdateTime"` // index find msg
}

type KeyValue struct {
	TypeKey          string `bson:"type_key" json:"typeKey"`
	Value            string `bson:"value" json:"value"`
	LatestUpdateTime int64  `bson:"latest_update_time" json:"latestUpdateTime"`
}

type ExtendMsg struct {
	ReactionExtensionList map[string]KeyValue `bson:"reaction_extension_list" json:"reactionExtensionList"`
	ClientMsgID           string              `bson:"client_msg_id" json:"clientMsgID"`
	MsgFirstModifyTime    int64               `bson:"msg_first_modify_time" json:"msgFirstModifyTime"` // this extendMsg create time
	AttachedInfo          string              `bson:"attached_info" json:"attachedInfo"`
	Ex                    string              `bson:"ex" json:"ex"`
}

func GetExtendMsgMaxNum() int32 {
	return MaxNum
}

func GetExtendMsgSourceID(ID string, index int32) string {
	return ID + ":" + strconv.Itoa(int(index))
}

func SplitSourceIDAndGetIndex(sourceID string) int32 {
	l := strings.Split(sourceID, ":")
	index, _ := strconv.Atoi(l[len(l)-1])
	return int32(index)
}

func (d *DataBases) CreateExtendMsgSet(set *ExtendMsgSet) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cExtendMsgSet)
	_, err := c.InsertOne(ctx, set)
	return err
}

type GetAllExtendMsgSetOpts struct {
	ExcludeExtendMsgs bool
}

func (d *DataBases) GetAllExtendMsgSet(ID string, opts *GetAllExtendMsgSetOpts) (sets []*ExtendMsgSet, err error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cExtendMsgSet)
	regex := fmt.Sprintf("^%s", ID)
	var findOpts *options.FindOptions
	if opts != nil {
		if opts.ExcludeExtendMsgs {
			findOpts = &options.FindOptions{}
			findOpts.SetProjection(bson.M{"extend_msgs": 0})
		}
	}
	cursor, err := c.Find(ctx, bson.M{"uid": primitive.Regex{Pattern: regex}}, findOpts)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	err = cursor.All(context.Background(), &sets)
	if err != nil {
		return nil, utils.Wrap(err, fmt.Sprintf("cursor is %s", cursor.Current.String()))
	}
	return sets, nil
}

type GetExtendMsgSetOpts struct {
	ExcludeExtendMsgs bool
}

// first modify msg
func (d *DataBases) InsertExtendMsg(sourceID string, sessionType int32, msg *ExtendMsg) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cExtendMsgSet)
	regex := fmt.Sprintf("^%s", sourceID)
	var err error
	findOpts := options.Find().SetLimit(1).SetSkip(0).SetSort(bson.M{"source_id": -1}).SetProjection(bson.M{"extend_msgs": 0})
	// update newest
	result, err := c.Find(ctx, bson.M{"source_id": primitive.Regex{Pattern: regex}, "session_type": sessionType}, findOpts)
	if err != nil {
		return utils.Wrap(err, "")
	}
	var setList []ExtendMsgSet
	if err := result.All(ctx, &setList); err != nil {
		return utils.Wrap(err, "")
	}
	if len(setList) == 0 || setList[0].ExtendMsgNum >= GetExtendMsgMaxNum() {
		var index int32
		if len(setList) > 0 {
			index = SplitSourceIDAndGetIndex(setList[0].SourceID)
		}
		err = d.CreateExtendMsgSet(&ExtendMsgSet{
			SourceID:         GetExtendMsgSourceID(sourceID, index),
			SessionType:      sessionType,
			ExtendMsgs:       map[string]ExtendMsg{msg.ClientMsgID: *msg},
			ExtendMsgNum:     1,
			CreateTime:       msg.MsgFirstModifyTime,
			MaxMsgUpdateTime: msg.MsgFirstModifyTime,
		})
	} else {
		_, err = c.UpdateOne(ctx, bson.M{"source_id": setList[0].SourceID, "session_type": sessionType}, bson.M{"$set": bson.M{"max_msg_update_time": msg.MsgFirstModifyTime, "$inc": bson.M{"extend_msg_num": 1}, fmt.Sprintf("extend_msgs.%s", msg.ClientMsgID): msg}})
	}
	return utils.Wrap(err, "")
}

// insert or update
func (d *DataBases) InsertOrUpdateReactionExtendMsgSet(sourceID string, sessionType int32, clientMsgID string, msgFirstModifyTime int64, reactionExtensionList map[string]*server_api_params.KeyValue) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cExtendMsgSet)
	var updateBson = bson.M{}
	for _, v := range reactionExtensionList {
		updateBson[fmt.Sprintf("extend_msgs.%s.%s", clientMsgID, v.TypeKey)] = v
	}
	upsert := true
	opt := &options.UpdateOptions{
		Upsert: &upsert,
	}
	findOpts := options.Find().SetLimit(1).SetSkip(0).SetSort(bson.M{"source_id": -1}).SetProjection(bson.M{"extend_msgs": 0})
	regex := fmt.Sprintf("^%s", sourceID)
	result, err := c.Find(ctx, bson.M{"source_id": primitive.Regex{Pattern: regex}, "session_type": sessionType, "max_msg_update_time": bson.M{"$lte": msgFirstModifyTime}}, findOpts)
	if err != nil {
		return utils.Wrap(err, "")
	}
	var setList []ExtendMsgSet
	if err := result.All(ctx, &setList); err != nil {
		return utils.Wrap(err, "")
	}
	if len(setList) == 0 {
		return utils.Wrap(errors.New("InsertOrUpdateReactionExtendMsgSet failed, len(setList) == 0"), "")
	}

	_, err = c.UpdateOne(ctx, bson.M{"source_id": setList[0].SourceID, "session_type": sessionType}, bson.M{"$set": updateBson}, opt)
	return utils.Wrap(err, "")
}

// delete TypeKey
func (d *DataBases) DeleteReactionExtendMsgSet(sourceID string, sessionType int32, clientMsgID string, msgFirstModifyTime int64, reactionExtensionList []*server_api_params.KeyValue) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cExtendMsgSet)
	var updateBson = bson.M{}
	for _, v := range reactionExtensionList {
		updateBson[fmt.Sprintf("extend_msgs.%s.%s", clientMsgID, v.TypeKey)] = ""
	}

	findOpts := options.Find().SetLimit(1).SetSkip(0).SetSort(bson.M{"source_id": -1}).SetProjection(bson.M{"extend_msgs": 0})
	regex := fmt.Sprintf("^%s", sourceID)
	result, err := c.Find(ctx, bson.M{"source_id": primitive.Regex{Pattern: regex}, "session_type": sessionType, "max_msg_update_time": bson.M{"$lte": msgFirstModifyTime}}, findOpts)
	if err != nil {
		return utils.Wrap(err, "")
	}
	var setList []ExtendMsgSet
	if err := result.All(ctx, &setList); err != nil {
		return utils.Wrap(err, "")
	}
	if len(setList) == 0 {
		return utils.Wrap(errors.New("InsertOrUpdateReactionExtendMsgSet failed, len(setList) == 0"), "")
	}
	_, err = c.UpdateOne(ctx, bson.M{"source_id": setList[0].SourceID, "session_type": sessionType}, bson.M{"$unset": updateBson})
	return err
}

func (d *DataBases) GetExtendMsg(sourceID string, sessionType int32, clientMsgID string, maxMsgUpdateTime int64) (extendMsg *ExtendMsg, err error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cExtendMsgSet)
	findOpts := options.Find().SetLimit(1).SetSkip(0).SetSort(bson.M{"source_id": -1}).SetProjection(bson.M{fmt.Sprintf("extend_msgs.%s", clientMsgID): 1})
	regex := fmt.Sprintf("^%s", sourceID)
	result, err := c.Find(ctx, bson.M{"source_id": primitive.Regex{Pattern: regex}, "session_type": sessionType, "msgFirstModifyTime": bson.M{"$lte": maxMsgUpdateTime}}, findOpts)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var extendMsgList []ExtendMsg
	if err := result.All(ctx, &extendMsgList); err != nil {
		return nil, utils.Wrap(err, "")
	}
	if len(extendMsgList) == 0 {
		return nil, utils.Wrap(errors.New("GetExtendMsg failed, len(setList) == 0"), "")
	}
	return &extendMsgList[0], nil
}
