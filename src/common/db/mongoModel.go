package db

import (
	"Open_IM/src/common/config"
	"Open_IM/src/common/constant"
	"Open_IM/src/common/log"
	pbMsg "Open_IM/src/proto/chat"
	"errors"
	"github.com/golang/protobuf/proto"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const cChat = "chat"
const cGroup = "group"

type MsgInfo struct {
	SendTime int64
	Msg      []byte
}

type UserChat struct {
	UID string
	Msg []MsgInfo
	BId bson.ObjectId `bson:"bid"`
}

type GroupMember struct {
	GroupID string
	UIDList []string
}

func (d *DataBases) GetMsgBySeqRange(uid string, seqBegin, seqEnd int64) (SingleMsg []*pbMsg.MsgFormat, GroupMsg []*pbMsg.MsgFormat, MaxSeq int64, MinSeq int64, err error) {
	var count int64
	session := d.mgoSession.Clone()
	if session == nil {
		return nil, nil, MaxSeq, MinSeq, errors.New("session == nil")
	}
	defer session.Close()

	c := session.DB(config.Config.Mongo.DBDatabase).C(cChat)

	sChat := UserChat{}
	if err = c.Find(bson.M{"uid": uid}).One(&sChat); err != nil {
		return nil, nil, MaxSeq, MinSeq, err
	}
	pChat := pbMsg.MsgSvrToPushSvrChatMsg{}
	for i := 0; i < len(sChat.Msg); i++ {
		temp := new(pbMsg.MsgFormat)
		if err = proto.Unmarshal(sChat.Msg[i].Msg, &pChat); err != nil {
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
			temp.ClientMsgID = pChat.ClientMsgID
			temp.SenderFaceURL = pChat.SenderFaceURL
			temp.SenderNickName = pChat.SenderNickName
			if pChat.RecvSeq > MaxSeq {
				MaxSeq = pChat.RecvSeq
			}
			if count == 0 {
				MinSeq = pChat.RecvSeq
			}
			if pChat.RecvSeq < MinSeq {
				MinSeq = pChat.RecvSeq
			}
			if pChat.SessionType == constant.SingleChatType {
				SingleMsg = append(SingleMsg, temp)
			} else {
				GroupMsg = append(GroupMsg, temp)
			}
			count++
			if count == (seqEnd - seqBegin + 1) {
				break
			}
		}
	}

	return SingleMsg, GroupMsg, MaxSeq, MinSeq, nil
}
func (d *DataBases) GetMsgBySeqList(uid string, seqList []int64) (SingleMsg []*pbMsg.MsgFormat, GroupMsg []*pbMsg.MsgFormat, MaxSeq int64, MinSeq int64, err error) {
	count := 0
	session := d.mgoSession.Clone()
	if session == nil {
		return nil, nil, MaxSeq, MinSeq, errors.New("session == nil")
	}
	defer session.Close()

	c := session.DB(config.Config.Mongo.DBDatabase).C(cChat)

	sChat := UserChat{}
	if err = c.Find(bson.M{"uid": uid}).One(&sChat); err != nil {
		return nil, nil, MaxSeq, MinSeq, err
	}
	pChat := pbMsg.MsgSvrToPushSvrChatMsg{}
	for i := 0; i < len(sChat.Msg); i++ {
		temp := new(pbMsg.MsgFormat)
		if err = proto.Unmarshal(sChat.Msg[i].Msg, &pChat); err != nil {
			return nil, nil, MaxSeq, MinSeq, err
		}
		if isContainInt64(pChat.RecvSeq, seqList) {
			temp.SendID = pChat.SendID
			temp.RecvID = pChat.RecvID
			temp.MsgFrom = pChat.MsgFrom
			temp.Seq = pChat.RecvSeq
			temp.ServerMsgID = pChat.MsgID
			temp.SendTime = pChat.SendTime
			temp.Content = pChat.Content
			temp.ContentType = pChat.ContentType
			temp.SenderPlatformID = pChat.PlatformID
			temp.ClientMsgID = pChat.ClientMsgID
			temp.SenderFaceURL = pChat.SenderFaceURL
			temp.SenderNickName = pChat.SenderNickName
			if pChat.RecvSeq > MaxSeq {
				MaxSeq = pChat.RecvSeq
			}
			if count == 0 {
				MinSeq = pChat.RecvSeq
			}
			if pChat.RecvSeq < MinSeq {
				MinSeq = pChat.RecvSeq
			}
			if pChat.SessionType == constant.SingleChatType {
				SingleMsg = append(SingleMsg, temp)
			} else {
				GroupMsg = append(GroupMsg, temp)
			}
			count++
			if count == len(seqList) {
				break
			}
		}
	}

	return SingleMsg, GroupMsg, MaxSeq, MinSeq, nil
}
func (d *DataBases) SaveUserChat(uid string, sendTime int64, m proto.Message) error {
	newTime := getCurrentTimestampByMill()
	session := d.mgoSession.Clone()
	if session == nil {
		return errors.New("session == nil")
	}
	defer func() {
		if session != nil {
			session.Close()
		}
	}()

	log.NewInfo("", "get mgoSession cost time", getCurrentTimestampByMill()-newTime)
	c := session.DB(config.Config.Mongo.DBDatabase).C(cChat)
	err := c.EnsureIndexKey("uid")
	if err != nil {
		log.NewError("", "EnsureIndexKey uid failed ", err.Error())
	}
	n, err := c.Find(bson.M{"uid": uid}).Count()
	if err != nil {
		return err
	}
	log.NewInfo("", "find and create index mgo uid cost time", getCurrentTimestampByMill()-newTime)
	sMsg := MsgInfo{}
	sMsg.SendTime = sendTime
	if sMsg.Msg, err = proto.Marshal(m); err != nil {
		return err
	}

	if n == 0 {
		sChat := UserChat{}
		sChat.UID = uid
		sChat.Msg = append(sChat.Msg, sMsg)
		err = c.Insert(&sChat)
		if err != nil {
			return err
		}
	} else {
		err = c.Update(bson.M{"uid": uid}, bson.M{"$push": bson.M{"msg": sMsg}})
		if err != nil {
			return err
		}
	}
	log.NewInfo("", "insert mgo data cost time", getCurrentTimestampByMill()-newTime)
	return nil
}

func (d *DataBases) DelUserChat(uid string) error {
	session := d.mgoSession.Clone()
	if session == nil {
		return errors.New("session == nil")
	}
	defer session.Close()

	c := session.DB(config.Config.Mongo.DBDatabase).C(cChat)

	delTime := time.Now().Unix() - int64(config.Config.Mongo.DBRetainChatRecords)*24*3600
	if err := c.Update(bson.M{"uid": uid}, bson.M{"$pull": bson.M{"msg": bson.M{"sendtime": bson.M{"$lte": delTime}}}}); err != nil {
		return err
	}

	return nil
}

func (d *DataBases) MgoUserCount() (int, error) {
	session := d.mgoSession.Clone()
	if session == nil {
		return 0, errors.New("session == nil")
	}
	defer session.Close()

	c := session.DB(config.Config.Mongo.DBDatabase).C(cChat)

	return c.Find(nil).Count()
}

func (d *DataBases) MgoSkipUID(count int) (string, error) {
	session := d.mgoSession.Clone()
	if session == nil {
		return "", errors.New("session == nil")
	}
	defer session.Close()

	c := session.DB(config.Config.Mongo.DBDatabase).C(cChat)

	sChat := UserChat{}
	c.Find(nil).Skip(count).Limit(1).One(&sChat)
	return sChat.UID, nil
}

func (d *DataBases) GetGroupMember(groupID string) []string {
	groupInfo := GroupMember{}
	groupInfo.GroupID = groupID
	groupInfo.UIDList = make([]string, 0)

	session := d.mgoSession.Clone()
	if session == nil {
		return groupInfo.UIDList
	}
	defer session.Close()

	c := session.DB(config.Config.Mongo.DBDatabase).C(cGroup)

	if err := c.Find(bson.M{"groupid": groupInfo.GroupID}).One(&groupInfo); err != nil {
		return groupInfo.UIDList
	}

	return groupInfo.UIDList
}

func (d *DataBases) AddGroupMember(groupID, uid string) error {
	session := d.mgoSession.Clone()
	if session == nil {
		return errors.New("session == nil")
	}
	defer session.Close()

	c := session.DB(config.Config.Mongo.DBDatabase).C(cGroup)

	n, err := c.Find(bson.M{"groupid": groupID}).Count()
	if err != nil {
		return err
	}

	if n == 0 {
		groupInfo := GroupMember{}
		groupInfo.GroupID = groupID
		groupInfo.UIDList = append(groupInfo.UIDList, uid)
		err = c.Insert(&groupInfo)
		if err != nil {
			return err
		}
	} else {
		err = c.Update(bson.M{"groupid": groupID}, bson.M{"$addToSet": bson.M{"uidlist": uid}})
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *DataBases) DelGroupMember(groupID, uid string) error {
	session := d.mgoSession.Clone()
	if session == nil {
		return errors.New("session == nil")
	}
	defer session.Close()

	c := session.DB(config.Config.Mongo.DBDatabase).C(cGroup)

	if err := c.Update(bson.M{"groupid": groupID}, bson.M{"$pull": bson.M{"uidlist": uid}}); err != nil {
		return err
	}

	return nil
}
func isContainInt64(target int64, List []int64) bool {

	for _, element := range List {

		if target == element {
			return true
		}
	}
	return false

}
func getCurrentTimestampByMill() int64 {
	return time.Now().UnixNano() / 1e6
}
