package db

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	log2 "Open_IM/pkg/common/log"
	pbChat "Open_IM/pkg/proto/chat"
	pbRtc "Open_IM/pkg/proto/rtc"
	pbCommon "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"

	//osconfig "google.golang.org/genproto/googleapis/cloud/osconfig/v1alpha"
	"strconv"
)

const (
	accountTempCode               = "ACCOUNT_TEMP_CODE"
	resetPwdTempCode              = "RESET_PWD_TEMP_CODE"
	userIncrSeq                   = "REDIS_USER_INCR_SEQ:" // user incr seq
	appleDeviceToken              = "DEVICE_TOKEN"
	userMinSeq                    = "REDIS_USER_MIN_SEQ:"
	uidPidToken                   = "UID_PID_TOKEN_STATUS:"
	conversationReceiveMessageOpt = "CON_RECV_MSG_OPT:"
	getuiToken                    = "GETUI"
	userInfoCache                 = "USER_INFO_CACHE:"
	friendRelationCache           = "FRIEND_RELATION_CACHE:"
	blackListCache                = "BLACK_LIST_CACHE:"
	groupCache                    = "GROUP_CACHE:"
	messageCache                  = "MESSAGE_CACHE:"
	SignalCache                   = "SIGNAL_CACHE:"
	SignalListCache               = "SIGNAL_LIST_CACHE:"
	GlobalMsgRecvOpt              = "GLOBAL_MSG_RECV_OPT"
)

func (d *DataBases) Exec(cmd string, key interface{}, args ...interface{}) (interface{}, error) {
	con := d.redisPool.Get()
	if err := con.Err(); err != nil {
		log2.Error("", "", "redis cmd = %v, err = %v", cmd, err)
		return nil, err
	}
	defer con.Close()

	params := make([]interface{}, 0)
	params = append(params, key)

	if len(args) > 0 {
		for _, v := range args {
			params = append(params, v)
		}
	}

	return con.Do(cmd, params...)
}
func (d *DataBases) JudgeAccountEXISTS(account string) (bool, error) {
	key := accountTempCode + account
	return redis.Bool(d.Exec("EXISTS", key))
}
func (d *DataBases) SetAccountCode(account string, code, ttl int) (err error) {
	key := accountTempCode + account
	_, err = d.Exec("SET", key, code, "ex", ttl)
	return err
}
func (d *DataBases) GetAccountCode(account string) (string, error) {
	key := accountTempCode + account
	return redis.String(d.Exec("GET", key))
}

//Perform seq auto-increment operation of user messages
func (d *DataBases) IncrUserSeq(uid string) (uint64, error) {
	key := userIncrSeq + uid
	return redis.Uint64(d.Exec("INCR", key))
}

//Get the largest Seq
func (d *DataBases) GetUserMaxSeq(uid string) (uint64, error) {
	key := userIncrSeq + uid
	return redis.Uint64(d.Exec("GET", key))
}

//set the largest Seq
func (d *DataBases) SetUserMaxSeq(uid string, maxSeq uint64) error {
	key := userIncrSeq + uid
	_, err := d.Exec("SET", key, maxSeq)
	return err
}

//Set the user's minimum seq
func (d *DataBases) SetUserMinSeq(uid string, minSeq uint32) (err error) {
	key := userMinSeq + uid
	_, err = d.Exec("SET", key, minSeq)
	return err
}

//Get the smallest Seq
func (d *DataBases) GetUserMinSeq(uid string) (uint64, error) {
	key := userMinSeq + uid
	return redis.Uint64(d.Exec("GET", key))
}

//Store Apple's device token to redis
func (d *DataBases) SetAppleDeviceToken(accountAddress, value string) (err error) {
	key := appleDeviceToken + accountAddress
	_, err = d.Exec("SET", key, value)
	return err
}

//Delete Apple device token
func (d *DataBases) DelAppleDeviceToken(accountAddress string) (err error) {
	key := appleDeviceToken + accountAddress
	_, err = d.Exec("DEL", key)
	return err
}

//Store userid and platform class to redis
func (d *DataBases) AddTokenFlag(userID string, platformID int, token string, flag int) error {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
	log2.NewDebug("", "add token key is ", key)
	_, err1 := d.Exec("HSet", key, token, flag)
	return err1
}

func (d *DataBases) GetTokenMapByUidPid(userID, platformID string) (map[string]int, error) {
	key := uidPidToken + userID + ":" + platformID
	log2.NewDebug("", "get token key is ", key)
	return redis.IntMap(d.Exec("HGETALL", key))
}
func (d *DataBases) SetTokenMapByUidPid(userID string, platformID int, m map[string]int) error {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
	_, err := d.Exec("hmset", key, redis.Args{}.Add().AddFlat(m)...)
	return err
}
func (d *DataBases) DeleteTokenByUidPid(userID string, platformID int, fields []string) error {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
	_, err := d.Exec("HDEL", key, redis.Args{}.Add().AddFlat(fields)...)
	return err
}

func (d *DataBases) SetSingleConversationRecvMsgOpt(userID, conversationID string, opt int32) error {
	key := conversationReceiveMessageOpt + userID
	_, err := d.Exec("HSet", key, conversationID, opt)
	return err
}

func (d *DataBases) GetSingleConversationRecvMsgOpt(userID, conversationID string) (int, error) {
	key := conversationReceiveMessageOpt + userID
	return redis.Int(d.Exec("HGet", key, conversationID))
}
func (d *DataBases) GetAllConversationMsgOpt(userID string) (map[string]int, error) {
	key := conversationReceiveMessageOpt + userID
	return redis.IntMap(d.Exec("HGETALL", key))
}
func (d *DataBases) SetMultiConversationMsgOpt(userID string, m map[string]int) error {
	key := conversationReceiveMessageOpt + userID
	_, err := d.Exec("hmset", key, redis.Args{}.Add().AddFlat(m)...)
	return err
}
func (d *DataBases) GetMultiConversationMsgOpt(userID string, conversationIDs []string) (m map[string]int, err error) {
	m = make(map[string]int)
	key := conversationReceiveMessageOpt + userID
	i, err := redis.Ints(d.Exec("hmget", key, redis.Args{}.Add().AddFlat(conversationIDs)...))
	if err != nil {
		return m, err
	}
	for k, v := range conversationIDs {
		m[v] = i[k]
	}
	return m, nil

}

func (d *DataBases) SetGetuiToken(token string, expireTime int64) error {
	_, err := d.Exec("SET", getuiToken, token, "ex", expireTime)
	return err
}

func (d *DataBases) GetGetuiToken() (string, error) {
	result, err := redis.String(d.Exec("GET", getuiToken))
	return result, err
}

func (d *DataBases) SearchContentType() {

}

func (d *DataBases) SetUserInfoToCache(userID string, m map[string]interface{}) error {
	_, err := d.Exec("hmset", userInfoCache+userID, redis.Args{}.Add().AddFlat(m)...)
	return err
}

func (d *DataBases) GetUserInfoFromCache(userID string) (*pbCommon.UserInfo, error) {
	result, err := redis.String(d.Exec("hgetall", userInfoCache+userID))
	log2.NewInfo("", result)
	if err != nil {
		return nil, err
	}
	userInfo := &pbCommon.UserInfo{}
	err = json.Unmarshal([]byte(result), userInfo)
	return userInfo, err
}

func (d *DataBases) AddFriendToCache(userID string, friendIDList ...string) error {
	var IDList []interface{}
	for _, id := range friendIDList {
		IDList = append(IDList, id)
	}
	_, err := d.Exec("SADD", friendRelationCache+userID, IDList...)
	return err
}

func (d *DataBases) ReduceFriendToCache(userID string, friendIDList ...string) error {
	var IDList []interface{}
	for _, id := range friendIDList {
		IDList = append(IDList, id)
	}
	_, err := d.Exec("SREM", friendRelationCache+userID, IDList...)
	return err
}

func (d *DataBases) GetFriendIDListFromCache(userID string) ([]string, error) {
	result, err := redis.Strings(d.Exec("SMEMBERS", friendRelationCache+userID))
	return result, err
}

func (d *DataBases) AddBlackUserToCache(userID string, blackList ...string) error {
	var IDList []interface{}
	for _, id := range blackList {
		IDList = append(IDList, id)
	}
	_, err := d.Exec("SADD", blackListCache+userID, IDList...)
	return err
}

func (d *DataBases) ReduceBlackUserFromCache(userID string, blackList ...string) error {
	var IDList []interface{}
	for _, id := range blackList {
		IDList = append(IDList, id)
	}
	_, err := d.Exec("SREM", blackListCache+userID, IDList...)
	return err
}

func (d *DataBases) GetBlackListFromCache(userID string) ([]string, error) {
	result, err := redis.Strings(d.Exec("SMEMBERS", blackListCache+userID))
	return result, err
}

func (d *DataBases) AddGroupMemberToCache(groupID string, userIDList ...string) error {
	var IDList []interface{}
	for _, id := range userIDList {
		IDList = append(IDList, id)
	}
	_, err := d.Exec("SADD", groupCache+groupID, IDList...)
	return err
}

func (d *DataBases) ReduceGroupMemberFromCache(groupID string, userIDList ...string) error {
	var IDList []interface{}
	for _, id := range userIDList {
		IDList = append(IDList, id)
	}
	_, err := d.Exec("SREM", groupCache+groupID, IDList...)
	return err
}

func (d *DataBases) GetGroupMemberIDListFromCache(groupID string) ([]string, error) {
	result, err := redis.Strings(d.Exec("SMEMBERS", groupCache+groupID))
	return result, err
}
func (d *DataBases) GetMessageListBySeq(userID string, seqList []uint32, operationID string) (seqMsg []*pbCommon.MsgData, failedSeqList []uint32, errResult error) {
	for _, v := range seqList {
		//MESSAGE_CACHE:169.254.225.224_reliability1653387820_0_1
		key := messageCache + userID + "_" + strconv.Itoa(int(v))

		result, err := redis.String(d.Exec("GET", key))
		if err != nil {
			errResult = err
			failedSeqList = append(failedSeqList, v)
			log2.NewWarn(operationID, "redis get message error:", err.Error(), v)
		} else {
			msg := pbCommon.MsgData{}
			err = jsonpb.UnmarshalString(result, &msg)
			if err != nil {
				errResult = err
				failedSeqList = append(failedSeqList, v)
				log2.NewWarn(operationID, "Unmarshal err", result, err.Error())
			} else {
				log2.NewDebug(operationID, "redis get msg is ", msg.String())
				seqMsg = append(seqMsg, &msg)
			}

		}
	}
	return seqMsg, failedSeqList, errResult
}

func (d *DataBases) SetMessageToCache(msgList []*pbChat.MsgDataToMQ, uid string, operationID string) error {
	var failedList []pbChat.MsgDataToMQ
	for _, msg := range msgList {
		key := messageCache + uid + "_" + strconv.Itoa(int(msg.MsgData.Seq))
		s, err := utils.Pb2String(msg.MsgData)
		if err != nil {
			log2.NewWarn(operationID, utils.GetSelfFuncName(), "Pb2String failed", msg.MsgData.String(), uid, err.Error())
			continue
		}
		log2.NewDebug(operationID, "convert string is ", s)
		_, err = d.Exec("SET", key, s, "ex", config.Config.MsgCacheTimeout)
		if err != nil {
			log2.NewWarn(operationID, utils.GetSelfFuncName(), "redis failed", "args:", key, *msg, uid, s)
			failedList = append(failedList, *msg)
		}
	}
	if len(failedList) != 0 {
		return errors.New(fmt.Sprintf("set msg to cache failed, failed lists: %q", failedList))
	}
	return nil
}

func (d *DataBases) DelMsgFromCache(uid string, seqList []uint32, operationID string) {
	for _, seq := range seqList {
		key := messageCache + uid + "_" + strconv.Itoa(int(seq))
		result, err := redis.String(d.Exec("GET", key))
		if err != nil {
			log2.NewWarn(operationID, utils.GetSelfFuncName(), "redis failed", "args:", key)
			continue
		}
		log2.Debug(operationID, utils.GetSelfFuncName(), "del result", result)
		var msg pbCommon.MsgData
		err = utils.String2Pb(result, &msg)
		log2.NewDebug(operationID, utils.GetSelfFuncName(), "msg", msg)
		if err != nil {
			log2.NewWarn(operationID, utils.GetSelfFuncName(), "string2Pb failed", msg, err.Error())
			continue
		}
		msg.Status = constant.MsgDeleted
		s, err := utils.Pb2String(&msg)
		if err != nil {
			log2.NewWarn(operationID, utils.GetSelfFuncName(), "Pb2String failed", msg, err.Error())
			continue
		}
		_, err = d.Exec("SET", key, s)
		if err != nil {
			log2.NewWarn(operationID, utils.GetSelfFuncName(), "redis failed", "args:", key, msg, s)
		}
	}
}

func (d *DataBases) CacheSignalInfo(msg *pbCommon.MsgData) error {
	key := SignalCache + msg.ClientMsgID
	_, err := d.Exec("SET", key, msg.Content, "ex", config.Config.Rtc.SignalTimeout)
	return err
}

func (d *DataBases) GetSignalInfoFromCache(clientMsgID string) (invitationInfo *pbRtc.SignalInviteReq, err error) {
	key := SignalCache + clientMsgID
	result, err := redis.Bytes(d.Exec("GET", key))
	log2.NewDebug("", utils.GetSelfFuncName(), clientMsgID, result, string(result))
	if err != nil {
		return nil, err
	}
	req := &pbRtc.SignalReq{}
	if err = proto.Unmarshal(result, req); err != nil {
		return nil, err
	}
	req2 := req.Payload.(*pbRtc.SignalReq_Invite)
	invitationInfo = req2.Invite
	return invitationInfo, err
}
