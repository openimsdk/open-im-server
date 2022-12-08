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

type ExtendMsgSet struct {
	SourceID         string               `bson:"id" json:"ID"`
	SessionType      string               `bson:"session_type" json:"sessionType"`
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
	ReactionExtensionList map[string]KeyValue `bson:"content" json:"content"`
	ClientMsgID           string              `bson:"client_msg_id" json:"clientMsgID"`
	MsgFirstModifyTime    int64               `bson:"create_time" json:"createTime"` // this extendMsg create time
	AttachedInfo          *string             `bson:"attached_info" json:"attachedInfo"`
	Ex                    *string             `bson:"ex" json:"ex"`
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

func (d *DataBases) GetExtendMsgSet(ID string, index int32, opts *GetExtendMsgSetOpts) (*ExtendMsgSet, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cExtendMsgSet)
	var set ExtendMsgSet
	var findOneOpt *options.FindOneOptions
	if opts != nil {
		if opts.ExcludeExtendMsgs {
			findOneOpt = &options.FindOneOptions{}
			findOneOpt.SetProjection(bson.M{"extend_msgs": 0})
		}
	}
	err := c.FindOne(ctx, bson.M{"uid": GetExtendMsgSetID(ID, index)}, findOneOpt).Decode(&set)
	return &set, err
}

// first modify msg
func (d *DataBases) InsertExtendMsgAndGetIndex(ID string, index int32, msg *ExtendMsg) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cExtendMsgSet)
	_, err := c.UpdateOne(ctx, bson.M{"uid": GetExtendMsgSetID(ID, index)}, bson.M{"$set": bson.M{"latest_update_time": utils.GetCurrentTimestampBySecond(), "$inc": bson.M{"extend_msg_num": 1}, fmt.Sprintf("extend_msgs.%s", msg.ClientMsgID): msg}})
	return err
}

// insert or update
func (d *DataBases) InsertOrUpdateReactionExtendMsgSet(ID string, index int32, clientMsgID, userID, value string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cExtendMsgSet)
	reactionExtendMsgSet := KeyValue{
		TypeKey:          userID,
		Value:            value,
		LatestUpdateTime: utils.GetCurrentTimestampBySecond(),
	}
	upsert := true
	opt := &options.UpdateOptions{
		Upsert: &upsert,
	}
	_, err := c.UpdateOne(ctx, bson.M{"uid": GetExtendMsgSetID(ID, index)}, bson.M{"$set": bson.M{"latest_update_time": utils.GetCurrentTimestampBySecond()}, fmt.Sprintf("extend_msgs.%s.%s", clientMsgID, userID): reactionExtendMsgSet}, opt)
	return err
}

func (d *DataBases) DeleteReactionExtendMsgSet(ID string, index int32, clientMsgID, userID string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cExtendMsgSet)
	_, err := c.UpdateOne(ctx, bson.M{"uid": GetExtendMsgSetID(ID, index)}, bson.M{"$unset": bson.M{fmt.Sprintf("extend_msgs.%s.%s", clientMsgID, userID): ""}})
	return err
}

// by index start end
func (d *DataBases) GetExtendMsgList(ID string, index int32, clientMsgID string) (extendMsg *ExtendMsg, err error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cExtendMsgSet)
	err = c.FindOne(ctx, bson.M{"uid": GetExtendMsgSetID(ID, index), "extend_msgs": bson.M{}}).Decode(&extendMsg)
	return extendMsg, err
}
