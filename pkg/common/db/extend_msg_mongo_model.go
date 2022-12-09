package db

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/utils"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

const cExtendMsgSet = "extend_msgs"
const MaxNum = 100

type ExtendMsgSet struct {
	SourceID         string               `bson:"source_id" json:"ID"`
	SessionType      int32                `bson:"session_type" json:"sessionType"`
	ExtendMsgs       map[string]ExtendMsg `bson:"extend_msgs" json:"extendMsgs"`
	ExtendMsgNum     int32                `bson:"extend_msg_num" json:"extendMsgNum"`
	CreateTime       int64                `bson:"create_time" json:"createTime"` // this block's create time
	MaxMsgUpdateTime int64                `bson:"max_msg_update_time"`           // index find msg
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

func GetExtendMsgSetID(ID string, index int32) string {
	return ID + ":" + strconv.Itoa(int(index))
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
	result, err := c.UpdateOne(ctx, bson.M{"source_id": sourceID, "session_type": sessionType}, bson.M{"$set": bson.M{"max_msg_update_time": msg.MsgFirstModifyTime, "$inc": bson.M{"extend_msg_num": 1}, fmt.Sprintf("extend_msgs.%s", msg.ClientMsgID): msg}})
	if err != nil {
		return utils.Wrap(err, "")
	}
	if result.UpsertedCount == 0 {
		if err := d.CreateExtendMsgSet(&ExtendMsgSet{
			SourceID:         sourceID,
			SessionType:      sessionType,
			ExtendMsgs:       map[string]ExtendMsg{msg.ClientMsgID: *msg},
			ExtendMsgNum:     1,
			CreateTime:       msg.MsgFirstModifyTime,
			MaxMsgUpdateTime: msg.MsgFirstModifyTime,
		}); err != nil {
			return err
		}
	}
	return nil
}

// insert or update
func (d *DataBases) InsertOrUpdateReactionExtendMsgSet(sourceID string, sessionType int32, clientMsgID, typeKey, value string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cExtendMsgSet)
	reactionExtendMsgSet := KeyValue{
		TypeKey:          typeKey,
		Value:            value,
		LatestUpdateTime: utils.GetCurrentTimestampBySecond(),
	}
	upsert := true
	opt := &options.UpdateOptions{
		Upsert: &upsert,
	}
	_, err := c.UpdateOne(ctx, bson.M{"source_id": sourceID, "session_type": sessionType}, bson.M{"$set": bson.M{fmt.Sprintf("extend_msgs.%s.%s", clientMsgID, typeKey): reactionExtendMsgSet}}, opt)
	return err
}

// delete TypeKey
func (d *DataBases) DeleteReactionExtendMsgSet(sourceID string, sessionType int32, clientMsgID []string, userID string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cExtendMsgSet)
	_, err := c.UpdateOne(ctx, bson.M{"source_id": sourceID, "session_type": sessionType}, bson.M{"$unset": bson.M{fmt.Sprintf("extend_msgs.%s.%s", clientMsgID, userID): ""}})
	return err
}

func (d *DataBases) GetExtendMsg(sourceID string, sessionType int32, clientMsgID string, maxMsgUpdateTime int64) error {
}
