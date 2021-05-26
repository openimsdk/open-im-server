package db

import (
	"Open_IM/src/common/config"
	"Open_IM/src/common/constant"
	pbMsg "Open_IM/src/proto/chat"
	"errors"
	"github.com/golang/protobuf/proto"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type UserChat struct {
	UID string
	Msg [][]byte
}

func (d *DataBases) GetUserChat(uid string, seqBegin, seqEnd int64) (SingleMsg []*pbMsg.MsgFormat, GroupMsg []*pbMsg.MsgFormat, MaxSeq int64, MinSeq int64, err error) {
	session := d.session(config.Config.Mongo.DBDatabase[0]).Clone()
	if session == nil {
		return nil, nil, MaxSeq, MinSeq, errors.New("session == nil")
	}
	defer session.Close()

	c := session.DB(config.Config.Mongo.DBDatabase[0]).C("chat")

	sChat := UserChat{}
	if err = c.Find(bson.M{"uid": uid}).One(&sChat); err != nil {
		return nil, nil, MaxSeq, MinSeq, err
	}
	pChat := pbMsg.MsgSvrToPushSvrChatMsg{}
	for i := 0; i < len(sChat.Msg); i++ {
		//每次产生新的指针
		temp := new(pbMsg.MsgFormat)
		if err = proto.Unmarshal(sChat.Msg[i], &pChat); err != nil {
			return nil, nil, MaxSeq, MinSeq, err
		}
		if pChat.RecvSeq >= seqBegin && pChat.RecvSeq <= seqEnd {
			temp.SendID = pChat.SendID
			temp.RecvID = pChat.RecvID
			temp.MsgFrom = pChat.MsgFrom
			temp.Seq = pChat.RecvSeq
			temp.ServerMsgID = pChat.MsgID
			temp.SendTime = pChat.SendTime
			temp.Content = pChat.Content
			temp.ContentType = pChat.ContentType
			temp.SenderPlatformID = pChat.PlatformID
			if pChat.RecvSeq > MaxSeq {
				MaxSeq = pChat.RecvSeq
			}
			if i == 0 {
				MinSeq = pChat.RecvSeq
			}
			if pChat.RecvSeq < MinSeq {
				MinSeq = pChat.RecvSeq
			}
			//单聊消息
			if pChat.SessionType == constant.SingleChatType {
				SingleMsg = append(SingleMsg, temp)
			} else {
				GroupMsg = append(GroupMsg, temp)
			}
		}
	}

	//d.DelUserChat(&sChat)

	return SingleMsg, GroupMsg, MaxSeq, MinSeq, nil
}

func (d *DataBases) SaveUserChat(uid string, m proto.Message) error {
	session := d.session(config.Config.Mongo.DBDatabase[0]).Clone()
	if session == nil {
		return errors.New("session == nil")
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(config.Config.Mongo.DBDatabase[0]).C("chat")

	n, err := c.Find(bson.M{"uid": uid}).Count()
	if err != nil {
		return err
	}

	if n == 0 {
		sChat := UserChat{}
		sChat.UID = uid
		bMsg, _ := proto.Marshal(m)
		sChat.Msg = append(sChat.Msg, bMsg)

		err = c.Insert(&sChat)
		if err != nil {
			return err
		}
	} else {
		bMsg, err := proto.Marshal(m)
		err = c.Update(bson.M{"uid": uid}, bson.M{"$addToSet": bson.M{"msg": bMsg}})
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *DataBases) DelUserChat(uc *UserChat) {
	delMaxIndex := 0
	pbData := pbMsg.WSToMsgSvrChatMsg{}
	for i := 0; i < len(uc.Msg); i++ {
		if err := proto.Unmarshal(uc.Msg[i], &pbData); err != nil {
			delMaxIndex = i
		} else {
			if time.Now().Unix()-pbData.SendTime > 7*24*3600 {
				delMaxIndex = i
			} else {
				break
			}
		}
	}

	if delMaxIndex > 0 {
		uc.Msg = uc.Msg[delMaxIndex:]

		session := d.session(config.Config.Mongo.DBDatabase[0]).Clone()
		if session == nil {
			return
		}
		defer session.Close()

		c := session.DB(config.Config.Mongo.DBDatabase[0]).C("chat")
		if err := c.Update(bson.M{"uid": uc.UID}, bson.M{"msg": uc.Msg}); err != nil {
			return
		}
	}
}

func (d *DataBases) DelHistoryChat(days int64, ids []string) error {
	session := d.session(config.Config.Mongo.DBDatabase[0]).Clone()
	if session == nil {
		return errors.New("mgo session == nil")
	}
	defer session.Close()

	c := session.DB(config.Config.Mongo.DBDatabase[0]).C("chat")

	for i := 0; i < len(ids); i++ {
		d.delHistoryUserChat(c, days, ids[i])
		//time.Sleep(1 * time.Millisecond)
	}

	return nil
}

func (d *DataBases) delHistoryUserChat(c *mgo.Collection, days int64, id string) error {
	sChat := UserChat{}
	if err := c.Find(bson.M{"uid": id}).One(&sChat); err != nil {
		return err
	}

	delMaxIndex := 0
	pbData := pbMsg.WSToMsgSvrChatMsg{}
	for i := 0; i < len(sChat.Msg); i++ {
		if err := proto.Unmarshal(sChat.Msg[i], &pbData); err != nil {
			delMaxIndex = i
		} else {
			if time.Now().Unix()-pbData.SendTime > int64(days)*24*3600 {
				delMaxIndex = i
			} else {
				break
			}
		}
	}

	if delMaxIndex > 0 {
		if delMaxIndex < len(sChat.Msg) {
			sChat.Msg = sChat.Msg[delMaxIndex:]
		} else {
			sChat.Msg = sChat.Msg[0:0]
		}

		if err := c.Update(bson.M{"uid": sChat.UID}, bson.M{"msg": sChat.Msg}); err != nil {
			return err
		}
	}

	return nil
}
