package db

import (
	"Open_IM/pkg/common/config"
	"context"
	"strconv"
	"time"
)

const cExtendMsgSet = "extend_msg_set"

type ExtendMsgSet struct {
	ID               string       `bson:"id" json:"ID"`
	ExtendMsgs       []*ExtendMsg `bson:"extend_msg" json:"extendMsg"`
	LatestUpdateTime int32        `bson:"latest_update_time" json:"latestUpdateTime"`
	AttachedInfo     string       `bson:"attached_info" json:"attachedInfo"`
	Ex               string       `bson:"ex" json:"ex"`
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

//type Vote struct {
//	Content      string     `bson:"content" json:"content"`
//	AttachedInfo string     `bson:"attached_info" json:"attachedInfo"`
//	Ex           string     `bson:"ex" json:"ex"`
//	Options      []*Options `bson:"options" json:"options"`
//}
//
//type Options struct {
//	Content        string   `bson:"content" json:"content"`
//	AttachedInfo   string   `bson:"attached_info" json:"attachedInfo"`
//	Ex             string   `bson:"ex" json:"ex"`
//	VoteUserIDList []string `bson:"vote_user_id_list" json:"voteUserIDList"`
//}
//
//type ExtendMsgComment struct {
//	UserID         string `bson:"user_id" json:"userID"`
//	ReplyUserID    string `bson:"reply_user_id" json:"replyUserID"`
//	ReplyContentID string `bson:"reply_content_id" json:"replyContentID"`
//	ContentID      string `bson:"content_id" json:"contentID"`
//	Content        string `bson:"content" json:"content"`
//	CreateTime     int32  `bson:"create_time" json:"createTime"`
//	AttachedInfo   string `bson:"attached_info" json:"attachedInfo"`
//	Ex             string `bson:"ex" json:"ex"`
//}

func GetExtendMsgSetID(ID string, index int32) string {
	return ID + ":" + strconv.Itoa(int(index))
}

func (d *DataBases) CreateExtendMsgSet(set *ExtendMsgSet) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cExtendMsgSet)
	_, err := c.InsertOne(ctx, set)
	return err
}

func (d *DataBases) GetAllExtendMsgSet(ID string) ([]*ExtendMsgSet, error) {
	//ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	//c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cExtendMsgSet)

}

type GetExtendMsgSetOpts struct {
	IncludeExtendMsgs bool
}

func (d *DataBases) GetExtendMsgSet(ID string, index int32, opts *GetExtendMsgSetOpts) (*ExtendMsgSet, error) {
	//ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	//c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cExtendMsgSet)
}

func (d *DataBases) InsertExtendMsg(ID string, msg *ExtendMsg) error {
	//ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	//c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cExtendMsgSet)
	return nil
}

func (d *DataBases) UpdateOneExtendMsgSet(ID string, index, MsgIndex int32, msg *ExtendMsg, msgSet *ExtendMsgSet) error {
	//ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	//c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cExtendMsgSet)
	return nil
}

func (d *DataBases) GetExtendMsgList(ID string, index, msgStartIndex, msgEndIndex int32) ([]*ExtendMsgSet, error) {
	//ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	//c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cExtendMsgSet)
	return nil, nil
}
