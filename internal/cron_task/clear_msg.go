package cronTask

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	server_api_params "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"math"
	"strconv"
	"strings"

	goRedis "github.com/go-redis/redis/v8"
	"github.com/golang/protobuf/proto"
)

const oldestList = 0
const newestList = -1

func ResetUserGroupMinSeq(operationID, groupID string, userIDList []string) error {
	var delStruct delMsgRecursionStruct
	minSeq, err := deleteMongoMsg(operationID, groupID, oldestList, &delStruct)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), groupID, "deleteMongoMsg failed")
	}
	if minSeq == 0 {
		return nil
	}
	log.NewDebug(operationID, utils.GetSelfFuncName(), "delMsgIDList:", delStruct, "minSeq", minSeq)
	for _, userID := range userIDList {
		userMinSeq, err := db.DB.GetGroupUserMinSeq(groupID, userID)
		if err != nil && err != goRedis.Nil {
			log.NewError(operationID, utils.GetSelfFuncName(), "GetGroupUserMinSeq failed", groupID, userID, err.Error())
			continue
		}
		if userMinSeq < uint64(minSeq) {
			err = db.DB.SetGroupUserMinSeq(groupID, userID, uint64(minSeq))
		}
		if err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), groupID, userID, userMinSeq, minSeq)
		}
	}
	return nil
}

func DeleteMongoMsgAndResetRedisSeq(operationID, userID string) error {
	var delStruct delMsgRecursionStruct
	minSeq, err := deleteMongoMsg(operationID, userID, oldestList, &delStruct)
	if err != nil {
		return utils.Wrap(err, "")
	}
	if minSeq == 0 {
		return nil
	}
	log.NewDebug(operationID, utils.GetSelfFuncName(), "delMsgIDStruct: ", delStruct, "minSeq", minSeq)

	userCurrentMinSeq, err := db.DB.GetUserMinSeq(userID)
	if err != nil && err != goRedis.Nil {
		return err
	}
	userCurrentMaxSeq, err := db.DB.GetUserMaxSeq(userID)
	if err != nil && err != goRedis.Nil {
		return err
	}
	if userCurrentMinSeq > userCurrentMaxSeq {
		minSeq = uint32(userCurrentMaxSeq)
	}

	err = db.DB.SetUserMinSeq(userID, minSeq)
	return utils.Wrap(err, "")
}

// del list
func delMongoMsgsPhysical(uidList []string) error {
	if len(uidList) > 0 {
		err := db.DB.DelMongoMsgs(uidList)
		if err != nil {
			return utils.Wrap(err, "DelMongoMsgs failed")
		}
	}
	return nil
}

type delMsgRecursionStruct struct {
	minSeq     uint32
	delUidList []string
}

func (d *delMsgRecursionStruct) getSetMinSeq() uint32 {
	return d.minSeq
}

// index 0....19(del) 20...69
// seq 70
// set minSeq 21
// recursion 删除list并且返回设置的最小seq
func deleteMongoMsg(operationID string, ID string, index int64, delStruct *delMsgRecursionStruct) (uint32, error) {
	// find from oldest list
	msgs, err := db.DB.GetUserMsgListByIndex(ID, index)
	if err != nil || msgs.UID == "" {
		if err != nil {
			if err == db.ErrMsgListNotExist {
				log.NewInfo(operationID, utils.GetSelfFuncName(), "ID:", ID, "index:", index, err.Error())
			} else {
				log.NewError(operationID, utils.GetSelfFuncName(), "GetUserMsgListByIndex failed", err.Error(), index, ID)
			}
		}
		// 获取报错，或者获取不到了，物理删除并且返回seq
		err = delMongoMsgsPhysical(delStruct.delUidList)
		if err != nil {
			return 0, err
		}
		return delStruct.getSetMinSeq(), nil
	}
	log.NewDebug(operationID, "ID:", ID, "index:", index, "uid:", msgs.UID, "len:", len(msgs.Msg))
	if len(msgs.Msg) > db.GetSingleGocMsgNum() {
		log.NewWarn(operationID, utils.GetSelfFuncName(), "msgs too large", len(msgs.Msg), msgs.UID)
	}
	if msgs.Msg[len(msgs.Msg)-1].SendTime+(int64(config.Config.Mongo.DBRetainChatRecords)*24*60*60*1000) < utils.GetCurrentTimestampByMill() && msgListIsFull(msgs) {
		delStruct.delUidList = append(delStruct.delUidList, msgs.UID)
		lastMsgPb := &server_api_params.MsgData{}
		err = proto.Unmarshal(msgs.Msg[len(msgs.Msg)-1].Msg, lastMsgPb)
		if err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), len(msgs.Msg)-1, msgs.UID)
			return 0, utils.Wrap(err, "proto.Unmarshal failed")
		}
		delStruct.minSeq = lastMsgPb.Seq + 1
		log.NewDebug(operationID, utils.GetSelfFuncName(), msgs.UID, "add to delUidList", "minSeq", lastMsgPb.Seq+1)
	} else {
		var hasMarkDelFlag bool
		for index, msg := range msgs.Msg {
			if msg.SendTime == 0 {
				continue
			}
			msgPb := &server_api_params.MsgData{}
			err = proto.Unmarshal(msg.Msg, msgPb)
			if err != nil {
				log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), len(msgs.Msg)-1, msgs.UID)
				return 0, utils.Wrap(err, "proto.Unmarshal failed")
			}
			if utils.GetCurrentTimestampByMill() > msg.SendTime+(int64(config.Config.Mongo.DBRetainChatRecords)*24*60*60*1000) {
				msgPb.Status = constant.MsgDeleted
				bytes, _ := proto.Marshal(msgPb)
				msgs.Msg[index].Msg = bytes
				msgs.Msg[index].SendTime = 0
				hasMarkDelFlag = true
			} else {
				if err := delMongoMsgsPhysical(delStruct.delUidList); err != nil {
					return 0, err
				}
				if hasMarkDelFlag {
					log.NewInfo(operationID, ID, "hasMarkDelFlag", "index:", index, "msgPb:", msgPb, msgs.UID)
					if err := db.DB.UpdateOneMsgList(msgs); err != nil {
						return delStruct.getSetMinSeq(), utils.Wrap(err, "")
					}
				}
				return msgPb.Seq, nil
			}
		}
	}
	log.NewDebug(operationID, ID, "continue to", delStruct)
	//  继续递归 index+1
	seq, err := deleteMongoMsg(operationID, ID, index+1, delStruct)
	return seq, utils.Wrap(err, "deleteMongoMsg failed")
}

func msgListIsFull(chat *db.UserChat) bool {
	index, _ := strconv.Atoi(strings.Split(chat.UID, ":")[1])
	if index == 0 {
		if len(chat.Msg) >= 4999 {
			return true
		}
	}
	if len(chat.Msg) >= 5000 {
		return true
	}
	return false
}

func checkMaxSeqWithMongo(operationID, ID string, diffusionType int) error {
	var seqRedis uint64
	var err error
	if diffusionType == constant.WriteDiffusion {
		seqRedis, err = db.DB.GetUserMaxSeq(ID)
	} else {
		seqRedis, err = db.DB.GetGroupMaxSeq(ID)
	}
	if err != nil {
		if err == goRedis.Nil {

		} else {
			return utils.Wrap(err, "GetUserMaxSeq failed")
		}
	}
	msg, err := db.DB.GetNewestMsg(ID)
	if err != nil {
		return utils.Wrap(err, "GetNewestMsg failed")
	}
	if msg == nil {
		return nil
	}
	if math.Abs(float64(msg.Seq-uint32(seqRedis))) > 10 {
		log.NewWarn(operationID, utils.GetSelfFuncName(), "seqMongo, seqRedis", msg.Seq, seqRedis, ID, "redis maxSeq is different with msg.Seq > 10", "status: ", msg.Status, msg.SendTime)
	} else {
		log.NewInfo(operationID, utils.GetSelfFuncName(), "seqMongo, seqRedis", msg.Seq, seqRedis, ID, "seq and msg OK", "status:", msg.Status, msg.SendTime)
	}
	return nil
}
