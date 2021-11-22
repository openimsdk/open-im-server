package db

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	pbMsg "Open_IM/pkg/proto/chat"
	"errors"
	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

const cChat = "chat"
const cGroup = "group"
const cGroupMemberModel = "groupMemberModel"
const singleGocMsgNum = 5000

type MsgInfo struct {
	SendTime int64
	Msg      []byte
}

type UserChat struct {
	UID string
	Msg []MsgInfo
}

type GroupMember struct {
	GroupID string
	UIDList []string
}

type GroupMemberModel struct {
	GroupId            string    `bson:"group_id"`
	Uid                string    `bson:"uid"`
	NickName           string    `bson:"nickname"`
	AdministratorLevel int32     `bson:"administrator_level"`
	JoinTime           time.Time `bson:"join_time"`
	UserGroupFaceUrl   string    `bson:"user_group_face_url"`
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
func (d *DataBases) GetMinSeqFromMongo(uid string) (MinSeq int64, err error) {
	var i int64
	var seqUid string
	session := d.mgoSession.Clone()
	if session == nil {
		return MinSeq, errors.New("session == nil")
	}
	defer session.Close()
	c := session.DB(config.Config.Mongo.DBDatabase).C(cChat)
	MaxSeq, err := d.GetUserMaxSeq(uid)
	if err != nil && err != redis.ErrNil {
		return MinSeq, err
	}
	NB := MaxSeq / singleGocMsgNum
	for i = 0; i <= NB; i++ {
		seqUid = indexGen(uid, i)
		n, err := c.Find(bson.M{"uid": seqUid}).Count()
		if err == nil && n != 0 {
			if i == 0 {
				MinSeq = 1
			} else {
				MinSeq = i * singleGocMsgNum
			}
			break
		}
	}
	return MinSeq, nil
}
func (d *DataBases) GetMsgBySeqList(uid string, seqList []int64) (SingleMsg []*pbMsg.MsgFormat, GroupMsg []*pbMsg.MsgFormat, MaxSeq int64, MinSeq int64, err error) {
	allCount := 0
	singleCount := 0
	session := d.mgoSession.Clone()
	if session == nil {
		return nil, nil, MaxSeq, MinSeq, errors.New("session == nil")
	}
	defer session.Close()

	c := session.DB(config.Config.Mongo.DBDatabase).C(cChat)
	m := func(uid string, seqList []int64) map[string][]int64 {
		t := make(map[string][]int64)
		for i := 0; i < len(seqList); i++ {
			seqUid := getSeqUid(uid, seqList[i])
			if value, ok := t[seqUid]; !ok {
				var temp []int64
				t[seqUid] = append(temp, seqList[i])
			} else {
				t[seqUid] = append(value, seqList[i])
			}
		}
		return t
	}(uid, seqList)
	sChat := UserChat{}
	pChat := pbMsg.MsgSvrToPushSvrChatMsg{}
	for seqUid, value := range m {
		if err = c.Find(bson.M{"uid": seqUid}).One(&sChat); err != nil {
			log.NewError("", "not find seqUid", seqUid, value, uid, seqList)
			continue
		}
		singleCount = 0
		for i := 0; i < len(sChat.Msg); i++ {
			temp := new(pbMsg.MsgFormat)
			if err = proto.Unmarshal(sChat.Msg[i].Msg, &pChat); err != nil {
				log.NewError("", "not find seqUid", seqUid, value, uid, seqList)
				return nil, nil, MaxSeq, MinSeq, err
			}
			if isContainInt64(pChat.RecvSeq, value) {
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
				if allCount == 0 {
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
				allCount++
				singleCount++
				if singleCount == len(value) {
					break
				}
			}
		}
	}
	return SingleMsg, GroupMsg, MaxSeq, MinSeq, nil
}
func (d *DataBases) SaveUserChat(uid string, sendTime int64, m *pbMsg.MsgSvrToPushSvrChatMsg) error {
	var seqUid string
	newTime := getCurrentTimestampByMill()
	session := d.mgoSession.Clone()
	if session == nil {
		return errors.New("session == nil")
	}
	defer session.Close()
	log.NewInfo("", "get mgoSession cost time", getCurrentTimestampByMill()-newTime)
	c := session.DB(config.Config.Mongo.DBDatabase).C(cChat)
	seqUid = getSeqUid(uid, m.RecvSeq)
	n, err := c.Find(bson.M{"uid": seqUid}).Count()
	if err != nil {
		return err
	}
	log.NewInfo("", "find mgo uid cost time", getCurrentTimestampByMill()-newTime)
	sMsg := MsgInfo{}
	sMsg.SendTime = sendTime
	if sMsg.Msg, err = proto.Marshal(m); err != nil {
		return err
	}
	if n == 0 {
		sChat := UserChat{}
		sChat.UID = seqUid
		sChat.Msg = append(sChat.Msg, sMsg)
		err = c.Insert(&sChat)
		if err != nil {
			return err
		}
	} else {
		err = c.Update(bson.M{"uid": seqUid}, bson.M{"$push": bson.M{"msg": sMsg}})
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

func getCurrentTimestampByMill() int64 {
	return time.Now().UnixNano() / 1e6
}
func getSeqUid(uid string, seq int64) string {
	seqSuffix := seq / singleGocMsgNum
	return indexGen(uid, seqSuffix)
}
func isContainInt64(target int64, List []int64) bool {

	for _, element := range List {

		if target == element {
			return true
		}
	}
	return false

}
func indexGen(uid string, seqSuffix int64) string {
	return uid + ":" + strconv.FormatInt(seqSuffix, 10)
}

func (d *DataBases) InsertIntoGroupMember(groupId, uid, nickName, userGroupFaceUrl string, administratorLevel int32) error {
	session := d.mgoSession.Clone()
	if session == nil {
		return errors.New("session == nil")
	}
	defer session.Close()

	c := session.DB(config.Config.Mongo.DBDatabase).C(cGroupMemberModel)

	n, err := c.Find(bson.M{"group_id": groupId, "uid": uid}).Count()
	if err != nil {
		return err
	}

	if n == 0 {

		groupMsgInfo := GroupMemberModel{}
		groupMsgInfo.GroupId = groupId
		groupMsgInfo.Uid = uid
		groupMsgInfo.NickName = nickName
		groupMsgInfo.AdministratorLevel = administratorLevel
		groupMsgInfo.JoinTime = time.Now()
		groupMsgInfo.UserGroupFaceUrl = userGroupFaceUrl

		err = c.Insert(&groupMsgInfo)
		if err != nil {
			return err
		}
	} else {
		err = c.Update(bson.M{"group_id": groupId, "uid": uid}, bson.M{"$set": bson.M{"nickname": nickName,
			"administrator_level": administratorLevel, "join_time": time.Now(), "user_group_face_url": userGroupFaceUrl}})
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *DataBases) FindGroupMemberListByUserId(uid string) ([]GroupMemberModel, error) {
	session := d.mgoSession.Clone()
	if session == nil {
		return nil, errors.New("session == nil")
	}
	defer session.Close()

	c := session.DB(config.Config.Mongo.DBDatabase).C(cGroupMemberModel)

	var sGroupMemberModel []GroupMemberModel
	if err := c.Find(bson.M{"uid": uid}).All(&sGroupMemberModel); err != nil {
		return nil, err
	}

	return sGroupMemberModel, nil

}

func (d *DataBases) FindGroupMemberListByGroupId(groupId string) ([]GroupMemberModel, error) {
	session := d.mgoSession.Clone()
	if session == nil {
		return nil, errors.New("session == nil")
	}
	defer session.Close()

	c := session.DB(config.Config.Mongo.DBDatabase).C(cGroupMemberModel)

	var sGroupMemberModel []GroupMemberModel
	if err := c.Find(bson.M{"group_id": groupId}).All(&sGroupMemberModel); err != nil {
		return nil, err
	}

	return sGroupMemberModel, nil
}

func (d *DataBases) FindGroupMemberListByGroupIdAndFilterInfo(groupId string, filter int32) ([]GroupMemberModel, error) {
	session := d.mgoSession.Clone()
	if session == nil {
		return nil, errors.New("session == nil")
	}
	defer session.Close()

	c := session.DB(config.Config.Mongo.DBDatabase).C(cGroupMemberModel)

	var sGroupMemberModel []GroupMemberModel
	if err := c.Find(bson.M{"group_id": groupId, "administrator_level": filter}).All(&sGroupMemberModel); err != nil {
		return nil, err
	}

	return sGroupMemberModel, nil
}

func (d *DataBases) FindGroupMemberInfoByGroupIdAndUserId(groupId, uid string) (*GroupMemberModel, error) {
	session := d.mgoSession.Clone()
	if session == nil {
		return nil, errors.New("session == nil")
	}
	defer session.Close()

	c := session.DB(config.Config.Mongo.DBDatabase).C(cGroupMemberModel)

	sGroupMemberModel := GroupMemberModel{}
	if err := c.Find(bson.M{"group_id": groupId, "uid": uid}).One(&sGroupMemberModel); err != nil {
		return nil, err
	}

	return &sGroupMemberModel, nil
}

func (d *DataBases) DeleteGroupMemberByGroupIdAndUserId(groupId, uid string) error {
	session := d.mgoSession.Clone()
	if session == nil {
		return errors.New("session == nil")
	}
	defer session.Close()

	c := session.DB(config.Config.Mongo.DBDatabase).C(cGroupMemberModel)

	if err := c.Remove(bson.M{"group_id": groupId, "uid": uid}); err != nil {
		return err
	}

	return nil
}

func (d *DataBases) UpdateOwnerGroupNickName(groupId, userId, groupNickName string) error {
	session := d.mgoSession.Clone()
	if session == nil {
		return errors.New("session == nil")
	}
	defer session.Close()

	c := session.DB(config.Config.Mongo.DBDatabase).C(cGroupMemberModel)

	if err := c.Update(bson.M{"group_id": groupId, "uid": userId}, bson.M{"$set": bson.M{"nickname": groupNickName}}); err != nil {
		return err
	}

	return nil
}

func (d *DataBases) SelectGroupList(groupId string) ([]string, error) {
	session := d.mgoSession.Clone()
	if session == nil {
		return nil, errors.New("session == nil")
	}
	defer session.Close()

	c := session.DB(config.Config.Mongo.DBDatabase).C(cGroupMemberModel)

	var sGroupMemberModel []GroupMemberModel
	if err := c.Find(bson.M{"group_id": groupId}).All(&sGroupMemberModel); err != nil {
		return nil, err
	}

	var groupUserID string
	var groupList []string
	for _, GroupMemberModel := range sGroupMemberModel {
		groupUserID = GroupMemberModel.Uid
		groupList = append(groupList, groupUserID)
	}

	return groupList, nil
}

func (d *DataBases) UpdateTheUserAdministratorLevel(groupId, uid string, administratorLevel int64) error {
	session := d.mgoSession.Clone()
	if session == nil {
		return errors.New("session == nil")
	}
	defer session.Close()

	c := session.DB(config.Config.Mongo.DBDatabase).C(cGroupMemberModel)

	if err := c.Update(bson.M{"group_id": groupId, "uid": uid}, bson.M{"$set": bson.M{"administrator_level": administratorLevel}}); err != nil {
		return err
	}

	return nil
}

func (d *DataBases) GetOwnerManagerByGroupId(groupId string) ([]GroupMemberModel, error) {
	session := d.mgoSession.Clone()
	if session == nil {
		return nil, errors.New("session == nil")
	}
	defer session.Close()

	c := session.DB(config.Config.Mongo.DBDatabase).C(cGroupMemberModel)

	var sGroupMemberModel []GroupMemberModel
	if err := c.Find(bson.M{"group_id": groupId, "administrator_level": bson.M{"$gt": 0}}).All(&sGroupMemberModel); err != nil {
		return nil, err
	}

	return sGroupMemberModel, nil

}

func (d *DataBases) IsExistGroupMember(groupId, uid string) bool {
	session := d.mgoSession.Clone()
	if session == nil {
		return false
	}
	defer session.Close()

	c := session.DB(config.Config.Mongo.DBDatabase).C(cGroupMemberModel)

	number, err := c.Find(bson.M{"group_id": groupId, "uid": uid}).Count()
	if err != nil {
		return false
	}

	if number != 1 {
		return false
	}

	return true
}

func (d *DataBases) RemoveGroupMember(groupId string, memberId string) error {
	return d.DeleteGroupMemberByGroupIdAndUserId(groupId, memberId)
}

func (d *DataBases) GetMemberInfoById(groupId string, memberId string) (*GroupMemberModel, error) {
	return d.FindGroupMemberInfoByGroupIdAndUserId(groupId, memberId)
}

func (d *DataBases) GetGroupMemberByGroupId(groupId string, filter int32, begin int32, maxNumber int32) ([]GroupMemberModel, error) {
	memberList, err := d.FindGroupMemberListByGroupId(groupId) //sorted by join time
	if err != nil {
		return nil, err
	}
	if begin >= int32(len(memberList)) {
		return nil, nil
	}

	var end int32
	if begin+int32(maxNumber) < int32(len(memberList)) {
		end = begin + maxNumber
	} else {
		end = int32(len(memberList))
	}
	return memberList[begin:end], nil
}

func (d *DataBases) GetJoinedGroupIdListByMemberId(memberId string) ([]GroupMemberModel, error) {
	return d.FindGroupMemberListByUserId(memberId)
}

func (d *DataBases) GetGroupMemberNumByGroupId(groupId string) int32 {
	session := d.mgoSession.Clone()
	if session == nil {
		return 0
	}
	defer session.Close()

	c := session.DB(config.Config.Mongo.DBDatabase).C(cGroupMemberModel)

	var number int
	number, err := c.Find(bson.M{"group_id": groupId}).Count()
	if err != nil {
		return 0
	}

	if number != 1 {
		return 0
	}

	return int32(number)
}

func (d *DataBases) GetGroupOwnerByGroupId(groupId string) string {
	omList, err := d.GetOwnerManagerByGroupId(groupId)
	if err != nil {
		return ""
	}
	for _, v := range omList {
		if v.AdministratorLevel == 1 {
			return v.Uid
		}
	}
	return ""
}

func (d *DataBases) InsertGroupMember(groupId, userId, nickName, userFaceUrl string, role int32) error {
	return d.InsertIntoGroupMember(groupId, userId, nickName, userFaceUrl, role)
}
