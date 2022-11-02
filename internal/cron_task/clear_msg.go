package cronTask

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	server_api_params "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	goRedis "github.com/go-redis/redis/v8"
	"github.com/golang/protobuf/proto"
	"math"
)

const oldestList = 0
const newestList = -1

func ResetUserGroupMinSeq(operationID, groupID string, userIDList []string) error {
	var delMsgIDList [][2]interface{}
	minSeq, err := deleteMongoMsg(operationID, groupID, oldestList, &delMsgIDList)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), groupID, "deleteMongoMsg failed")
		return utils.Wrap(err, "")
	}
	if minSeq == 0 {
		return nil
	}
	log.NewDebug(operationID, utils.GetSelfFuncName(), "delMsgIDList:", delMsgIDList, "minSeq", minSeq)
	for _, userID := range userIDList {
		userMinSeq, err := db.DB.GetGroupUserMinSeq(groupID, userID)
		if err != nil && err != goRedis.Nil {
			log.NewError(operationID, utils.GetSelfFuncName(), "GetGroupUserMinSeq failed", groupID, userID, err.Error())
			continue
		}
		if userMinSeq > uint64(minSeq) {
			err = db.DB.SetGroupUserMinSeq(groupID, userID, userMinSeq)
		} else {
			err = db.DB.SetGroupUserMinSeq(groupID, userID, uint64(minSeq))
		}
		if err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), groupID, userID, userMinSeq, minSeq)
		}
	}
	return nil
}

func DeleteMongoMsgAndResetRedisSeq(operationID, userID string) error {
	var delMsgIDList [][2]interface{}
	minSeq, err := deleteMongoMsg(operationID, userID, oldestList, &delMsgIDList)
	if err != nil {
		return utils.Wrap(err, "")
	}
	if minSeq == 0 {
		return nil
	}
	log.NewDebug(operationID, utils.GetSelfFuncName(), "delMsgIDMap: ", delMsgIDList, "minSeq", minSeq)
	err = db.DB.SetUserMinSeq(userID, minSeq)
	return err
}

func delMongoMsgs(operationID string, delMsgIDList *[][2]interface{}) error {
	if len(*delMsgIDList) > 0 {
		var IDList []string
		for _, v := range *delMsgIDList {
			IDList = append(IDList, v[0].(string))
		}
		err := db.DB.DelMongoMsgs(IDList)
		if err != nil {
			return utils.Wrap(err, "DelMongoMsgs failed")
		}
	}
	return nil
}

// recursion
func deleteMongoMsg(operationID string, ID string, index int64, delMsgIDList *[][2]interface{}) (uint32, error) {
	// 从最旧的列表开始找
	msgs, err := db.DB.GetUserMsgListByIndex(ID, index)
	if err != nil || msgs.UID == "" {
		if err != nil {
			if err == db.ErrMsgListNotExist {
				log.NewDebug(operationID, utils.GetSelfFuncName(), ID, index, err.Error())
			} else {
				log.NewError(operationID, utils.GetSelfFuncName(), "GetUserMsgListByIndex failed", err.Error(), index, ID)
			}
		}
		return getDelMaxSeqByIDList(*delMsgIDList), delMongoMsgs(operationID, delMsgIDList)
	}
	if index == 0 && msgs == nil {
		return 0, nil
	}
	log.NewDebug(operationID, "ID:", ID, "index:", index, "uid:", msgs.UID, "len:", len(msgs.Msg))
	if len(msgs.Msg) > db.GetSingleGocMsgNum() {
		log.NewWarn(operationID, utils.GetSelfFuncName(), "msgs too large", len(msgs.Msg), msgs.UID)
	}
	log.NewDebug(operationID, utils.GetSelfFuncName(), "get msgs: ", msgs.UID)
	for i, msg := range msgs.Msg {
		// 找到列表中不需要删除的消息了
		if utils.GetCurrentTimestampByMill() < msg.SendTime+int64(config.Config.Mongo.DBRetainChatRecords)*24*60*60*1000 {
			if err := delMongoMsgs(operationID, delMsgIDList); err != nil {
				return 0, err
			}
			minSeq := getDelMaxSeqByIDList(*delMsgIDList)
			if i > 0 {
				msgPb := &server_api_params.MsgData{}
				err = proto.Unmarshal(msg.Msg, msgPb)
				if err != nil {
					log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), ID, index)
				} else {
					err = db.DB.ReplaceMsgToBlankByIndex(msgs.UID, i-1)
					if err != nil {
						log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), msgs.UID, i)
						return minSeq, nil
					}
					minSeq = msgPb.Seq
				}
			}
			return minSeq, nil
		}
	}
	if len(msgs.Msg) > 0 {
		msgPb := &server_api_params.MsgData{}
		err = proto.Unmarshal(msgs.Msg[len(msgs.Msg)-1].Msg, msgPb)
		if err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), len(msgs.Msg)-1, msgs.UID)
			return 0, utils.Wrap(err, "proto.Unmarshal failed")
		}
		*delMsgIDList = append(*delMsgIDList, [2]interface{}{msgs.UID, msgPb.Seq})
	}
	// 没有找到 代表需要全部删除掉 继续递归查找下一个比较旧的列表
	seq, err := deleteMongoMsg(operationID, utils.GetSelfFuncName(), index+1, delMsgIDList)
	if err != nil {
		return 0, utils.Wrap(err, "deleteMongoMsg failed")
	}
	return seq, nil
}

func getDelMaxSeqByIDList(delMsgIDList [][2]interface{}) uint32 {
	if len(delMsgIDList) == 0 {
		return 0
	}
	return delMsgIDList[len(delMsgIDList)-1][1].(uint32)
}

func checkMaxSeqWithMongo(operationID, ID string, diffusionType int) error {
	var maxSeq uint64
	var err error
	if diffusionType == constant.WriteDiffusion {
		maxSeq, err = db.DB.GetUserMaxSeq(ID)
	} else {
		maxSeq, err = db.DB.GetGroupMaxSeq(ID)
	}
	if err != nil {
		if err == goRedis.Nil {
			return nil
		}
		return utils.Wrap(err, "GetUserMaxSeq failed")
	}
	msg, err := db.DB.GetNewestMsg(ID)
	if err != nil {
		return utils.Wrap(err, "GetNewestMsg failed")
	}
	msgPb := &server_api_params.MsgData{}
	err = proto.Unmarshal(msg.Msg, msgPb)
	if err != nil {
		return utils.Wrap(err, "")
	}
	if math.Abs(float64(msgPb.Seq-uint32(maxSeq))) > 10 {
		log.NewWarn(operationID, utils.GetSelfFuncName(), maxSeq, msgPb.Seq, "redis maxSeq is different with msg.Seq")
	} else {
		log.NewInfo(operationID, utils.GetSelfFuncName(), diffusionType, ID, "seq and msg OK", msgPb.Seq, uint32(maxSeq))
	}
	return nil
}
