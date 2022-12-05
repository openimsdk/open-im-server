package db

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/utils"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

const cExtendMsgSet = "extend_msg_set"

type ExtendMsgSet struct {
	ID               string       `bson:"id" json:"ID"`
	ExtendMsgs       []*ExtendMsg `bson:"extend_msgs" json:"extendMsgs"`
	LatestUpdateTime int32        `bson:"latest_update_time" json:"latestUpdateTime"`
	AttachedInfo     *string      `bson:"attached_info" json:"attachedInfo"`
	Ex               *string      `bson:"ex" json:"ex"`
	ExtendMsgNum     int32        `bson:"extend_msg_num" json:"extendMsgNum"`
	CreateTime       int32        `bson:"create_time" json:"createTime"`
}

type ReactionExtendMsgSet struct {
	TypeKey string `bson:"type_key" json:"typeKey"`
	Value   string `bson:"value" json:"value"`
}

type ExtendMsg struct {
	Content     []*ReactionExtendMsgSet `bson:"content" json:"content"`
	ClientMsgID string                  `bson:"client_msg_id" json:"clientMsgID"`
	CreateTime  int32                   `bson:"create_time" json:"createTime"`
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

func (d *DataBases) GetAllExtendMsgSet(ID string) (sets []*ExtendMsgSet, err error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cExtendMsgSet)
	regex := fmt.Sprintf("^%s", ID)
	cursor, err := c.Find(ctx, bson.M{"uid": primitive.Regex{Pattern: regex}})
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
	IncludeExtendMsgs bool
}

func (d *DataBases) GetExtendMsgSet(ID string, index int32, opts *GetExtendMsgSetOpts) (*ExtendMsgSet, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cExtendMsgSet)
	var set ExtendMsgSet
	err := c.FindOne(ctx, bson.M{"uid": GetExtendMsgSetID(ID, index)}).Decode(&set)
	return &set, err
}

func (d *DataBases) InsertExtendMsg(ID string, index int32, msg *ExtendMsg) (msgIndex int32, err error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cExtendMsgSet)
	result := c.FindOneAndUpdate(ctx, bson.M{"uid": GetExtendMsgSetID(ID, index)}, bson.M{"$set": bson.M{"create_time": utils.GetCurrentTimestampBySecond(), "$inc": bson.M{"extend_msg_num": 1}, "$push": bson.M{"extend_msgs": msg}}})
	set := &ExtendMsgSet{}
	err = result.Decode(set)
	return set.ExtendMsgNum, err
}

func (d *DataBases) UpdateOneExtendMsgSet(ID string, index, MsgIndex int32, msg *ExtendMsg, msgSet *ExtendMsgSet) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cExtendMsgSet)
	_, err := c.UpdateOne(ctx, bson.M{"uid": GetExtendMsgSetID(ID, index)}, bson.M{})
	return err
}

func (d *DataBases) GetExtendMsgList(ID string, index, msgStartIndex, msgEndIndex int32) (extendMsgList []*ExtendMsg, err error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cExtendMsgSet)
	err = c.FindOne(ctx, bson.M{"uid": GetExtendMsgSetID(ID, index)}).Decode(&extendMsgList)
	return extendMsgList, err
}
