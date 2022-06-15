package db

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	log2 "Open_IM/pkg/common/log"
	pbChat "Open_IM/pkg/proto/chat"
	pbRtc "Open_IM/pkg/proto/rtc"
	pbCommon "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	go_redis "github.com/go-redis/redis/v8"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"strconv"
	"time"
)

//func  (d *  DataBases)pubMessage(channel, msg string) {
//   d.rdb.Publish(context.Background(),channel,msg)
//}
//func  (d *  DataBases)pubMessage(channel, msg string) {
//	d.rdb.Publish(context.Background(),channel,msg)
//}
func (d *DataBases) JudgeAccountEXISTS(account string) (bool, error) {
	key := accountTempCode + account
	n, err := d.rdb.Exists(context.Background(), key).Result()
	if n > 0 {
		return true, err
	} else {
		return false, err
	}
}
func (d *DataBases) SetAccountCode(account string, code, ttl int) (err error) {
	key := accountTempCode + account
	return d.rdb.Set(context.Background(), key, code, time.Duration(ttl)*time.Second).Err()
}
func (d *DataBases) GetAccountCode(account string) (string, error) {
	key := accountTempCode + account
	return d.rdb.Get(context.Background(), key).Result()
}

//Perform seq auto-increment operation of user messages
func (d *DataBases) IncrUserSeq(uid string) (uint64, error) {
	key := userIncrSeq + uid
	seq, err := d.rdb.Incr(context.Background(), key).Result()
	return uint64(seq), err
}

//Get the largest Seq
func (d *DataBases) GetUserMaxSeq(uid string) (uint64, error) {
	key := userIncrSeq + uid
	seq, err := d.rdb.Get(context.Background(), key).Result()
	return uint64(utils.StringToInt(seq)), err
}

//set the largest Seq
func (d *DataBases) SetUserMaxSeq(uid string, maxSeq uint64) error {
	key := userIncrSeq + uid
	return d.rdb.Set(context.Background(), key, maxSeq, 0).Err()
}

//Set the user's minimum seq
func (d *DataBases) SetUserMinSeq(uid string, minSeq uint32) (err error) {
	key := userMinSeq + uid
	return d.rdb.Set(context.Background(), key, minSeq, 0).Err()
}

//Get the smallest Seq
func (d *DataBases) GetUserMinSeq(uid string) (uint64, error) {
	key := userMinSeq + uid
	seq, err := d.rdb.Get(context.Background(), key).Result()
	return uint64(utils.StringToInt(seq)), err
}

//Store userid and platform class to redis
func (d *DataBases) AddTokenFlag(userID string, platformID int, token string, flag int) error {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
	log2.NewDebug("", "add token key is ", key)
	return d.rdb.HSet(context.Background(), key, token, flag).Err()
}

func (d *DataBases) GetTokenMapByUidPid(userID, platformID string) (map[string]int, error) {
	key := uidPidToken + userID + ":" + platformID
	log2.NewDebug("", "get token key is ", key)
	m, err := d.rdb.HGetAll(context.Background(), key).Result()
	mm := make(map[string]int)
	for k, v := range m {
		mm[k] = utils.StringToInt(v)
	}
	return mm, err
}
func (d *DataBases) SetTokenMapByUidPid(userID string, platformID int, m map[string]int) error {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
	return d.rdb.HMSet(context.Background(), key, m).Err()
}
func (d *DataBases) DeleteTokenByUidPid(userID string, platformID int, fields []string) error {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
	return d.rdb.HDel(context.Background(), key, fields...).Err()
}
func (d *DataBases) SetSingleConversationRecvMsgOpt(userID, conversationID string, opt int32) error {
	key := conversationReceiveMessageOpt + userID
	return d.rdb.HSet(context.Background(), key, conversationID, opt).Err()
}

func (d *DataBases) GetSingleConversationRecvMsgOpt(userID, conversationID string) (int, error) {
	key := conversationReceiveMessageOpt + userID
	result, err := d.rdb.HGet(context.Background(), key, conversationID).Result()
	return utils.StringToInt(result), err
}
func (d *DataBases) SetUserGlobalMsgRecvOpt(userID string, opt int32) error {
	key := conversationReceiveMessageOpt + userID
	return d.rdb.HSet(context.Background(), key, GlobalMsgRecvOpt, opt).Err()
}
func (d *DataBases) GetUserGlobalMsgRecvOpt(userID string) (int, error) {
	key := conversationReceiveMessageOpt + userID
	result, err := d.rdb.HGet(context.Background(), key, GlobalMsgRecvOpt).Result()
	if err != nil {
		if err == go_redis.Nil {
			return 0, nil
		} else {
			return 0, err
		}
	}
	return utils.StringToInt(result), err
}
func (d *DataBases) GetMessageListBySeq(userID string, seqList []uint32, operationID string) (seqMsg []*pbCommon.MsgData, failedSeqList []uint32, errResult error) {
	var keys []string
	for _, v := range seqList {
		//MESSAGE_CACHE:169.254.225.224_reliability1653387820_0_1
		key := messageCache + userID + "_" + strconv.Itoa(int(v))
		keys = append(keys, key)
	}
	result, err := d.rdb.MGet(context.Background(), keys...).Result()
	if err != nil {
		errResult = err
		failedSeqList = seqList
		log2.NewWarn(operationID, "redis get message error:", err.Error(), seqList)
	} else {
		for _, v := range result {
			msg := pbCommon.MsgData{}
			err = jsonpb.UnmarshalString(v.(string), &msg)
			if err != nil {
				errResult = err
				failedSeqList = seqList
				log2.NewWarn(operationID, "Unmarshal err", result, err.Error())
				break
			} else {
				log2.NewDebug(operationID, "redis get msg is ", msg.String())
				seqMsg = append(seqMsg, &msg)
			}
		}
	}
	return seqMsg, failedSeqList, errResult
}
func (d *DataBases) SetMessageToCache(msgList []*pbChat.MsgDataToMQ, uid string, operationID string) error {
	ctx := context.Background()
	pipe := d.rdb.Pipeline()
	var failedList []pbChat.MsgDataToMQ
	for _, msg := range msgList {
		key := messageCache + uid + "_" + strconv.Itoa(int(msg.MsgData.Seq))
		s, err := utils.Pb2String(msg.MsgData)
		if err != nil {
			log2.NewWarn(operationID, utils.GetSelfFuncName(), "Pb2String failed", msg.MsgData.String(), uid, err.Error())
			continue
		}
		log2.NewDebug(operationID, "convert string is ", s)
		err = pipe.Set(ctx, key, s, time.Duration(config.Config.MsgCacheTimeout)*time.Second).Err()
		//err = d.rdb.HMSet(context.Background(), "12", map[string]interface{}{"1": 2, "343": false}).Err()
		if err != nil {
			log2.NewWarn(operationID, utils.GetSelfFuncName(), "redis failed", "args:", key, *msg, uid, s, err.Error())
			failedList = append(failedList, *msg)
		}
	}
	if len(failedList) != 0 {
		return errors.New(fmt.Sprintf("set msg to cache failed, failed lists: %q,%s", failedList, operationID))
	}
	_, err := pipe.Exec(ctx)
	return err
}

func (d *DataBases) CleanUpOneUserAllMsgFromRedis(userID string, operationID string) error {
	ctx := context.Background()
	key := messageCache + userID + "_" + "*"
	vals, err := d.rdb.Keys(ctx, key).Result()
	log2.Debug(operationID, "vals: ", vals)
	if err == go_redis.Nil {
		return nil
	}
	if err != nil {
		return utils.Wrap(err, "")
	}
	if err = d.rdb.Del(ctx, vals...).Err(); err != nil {
		return utils.Wrap(err, "")
	}
	return nil
}

func (d *DataBases) HandleSignalInfo(operationID string, msg *pbCommon.MsgData) error {
	req := &pbRtc.SignalReq{}
	if err := proto.Unmarshal(msg.Content, req); err != nil {
		return err
	}
	//log.NewDebug(pushMsg.OperationID, utils.GetSelfFuncName(), "SignalReq: ", req.String())
	var inviteeUserIDList []string
	var isInviteSignal bool
	switch signalInfo := req.Payload.(type) {
	case *pbRtc.SignalReq_Invite:
		inviteeUserIDList = signalInfo.Invite.Invitation.InviteeUserIDList
		isInviteSignal = true
	case *pbRtc.SignalReq_InviteInGroup:
		inviteeUserIDList = signalInfo.InviteInGroup.Invitation.InviteeUserIDList
		isInviteSignal = true
	case *pbRtc.SignalReq_HungUp, *pbRtc.SignalReq_Cancel, *pbRtc.SignalReq_Reject, *pbRtc.SignalReq_Accept:
		return errors.New("signalInfo do not need offlinePush")
	default:
		log2.NewDebug(operationID, utils.GetSelfFuncName(), "req invalid type", string(msg.Content))
		return nil
	}
	if isInviteSignal {
		log2.NewInfo(operationID, utils.GetSelfFuncName(), "invite userID list:", inviteeUserIDList)
		for _, userID := range inviteeUserIDList {
			log2.NewInfo(operationID, utils.GetSelfFuncName(), "invite userID:", userID)
			timeout, err := strconv.Atoi(config.Config.Rtc.SignalTimeout)
			if err != nil {
				return err
			}
			keyList := SignalListCache + userID
			err = d.rdb.LPush(context.Background(), keyList, msg.ClientMsgID).Err()
			if err != nil {
				return err
			}
			err = d.rdb.Expire(context.Background(), keyList, time.Duration(timeout)*time.Second).Err()
			if err != nil {
				return err
			}
			key := SignalCache + msg.ClientMsgID
			err = d.rdb.Set(context.Background(), key, msg.Content, time.Duration(timeout)*time.Second).Err()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *DataBases) GetSignalInfoFromCacheByClientMsgID(clientMsgID string) (invitationInfo *pbRtc.SignalInviteReq, err error) {
	key := SignalCache + clientMsgID
	invitationInfo = &pbRtc.SignalInviteReq{}
	bytes, err := d.rdb.Get(context.Background(), key).Bytes()
	if err != nil {
		return nil, err
	}
	req := &pbRtc.SignalReq{}
	if err = proto.Unmarshal(bytes, req); err != nil {
		return nil, err
	}
	switch req2 := req.Payload.(type) {
	case *pbRtc.SignalReq_Invite:
		invitationInfo.Invitation = req2.Invite.Invitation
		invitationInfo.OpUserID = req2.Invite.OpUserID
	case *pbRtc.SignalReq_InviteInGroup:
		invitationInfo.Invitation = req2.InviteInGroup.Invitation
		invitationInfo.OpUserID = req2.InviteInGroup.OpUserID
	}
	return invitationInfo, err
}

func (d *DataBases) GetAvailableSignalInvitationInfo(userID string) (invitationInfo *pbRtc.SignalInviteReq, err error) {
	keyList := SignalListCache + userID
	result := d.rdb.LPop(context.Background(), keyList)
	if err = result.Err(); err != nil {
		return nil, utils.Wrap(err, "GetAvailableSignalInvitationInfo failed")
	}
	key, err := result.Result()
	if err != nil {
		return nil, utils.Wrap(err, "GetAvailableSignalInvitationInfo failed")
	}
	log2.NewDebug("", utils.GetSelfFuncName(), result, result.String())
	invitationInfo, err = d.GetSignalInfoFromCacheByClientMsgID(key)
	if err != nil {
		return nil, utils.Wrap(err, "GetSignalInfoFromCacheByClientMsgID")
	}
	err = d.DelUserSignalList(userID)
	if err != nil {
		return nil, utils.Wrap(err, "GetSignalInfoFromCacheByClientMsgID")
	}
	return invitationInfo, nil
}

func (d *DataBases) DelUserSignalList(userID string) error {
	keyList := SignalListCache + userID
	err := d.rdb.Del(context.Background(), keyList).Err()
	return err
}

func (d *DataBases) DelMsgFromCache(uid string, seqList []uint32, operationID string) {
	for _, seq := range seqList {
		key := messageCache + uid + "_" + strconv.Itoa(int(seq))
		result := d.rdb.Get(context.Background(), key).String()
		var msg pbCommon.MsgData
		if err := utils.String2Pb(result, &msg); err != nil {
			log2.Error(operationID, utils.GetSelfFuncName(), "String2Pb failed", msg, err.Error())
			continue
		}
		msg.Status = constant.MsgDeleted
		s, err := utils.Pb2String(&msg)
		if err != nil {
			log2.Error(operationID, utils.GetSelfFuncName(), "Pb2String failed", msg, err.Error())
			continue
		}
		if err := d.rdb.Set(context.Background(), key, s, time.Duration(config.Config.MsgCacheTimeout)*time.Second).Err(); err != nil {
			log2.Error(operationID, utils.GetSelfFuncName(), "Set failed", err.Error())
		}
	}
}

func (d *DataBases) SetGetuiToken(token string, expireTime int64) error {
	return d.rdb.Set(context.Background(), getuiToken, token, time.Duration(expireTime)*time.Second).Err()
}

func (d *DataBases) GetGetuiToken() (string, error) {
	result := d.rdb.Get(context.Background(), getuiToken)
	return result.String(), result.Err()
}

func (d *DataBases) AddFriendToCache(userID string, friendIDList ...string) error {
	var IDList []interface{}
	for _, id := range friendIDList {
		IDList = append(IDList, id)
	}
	return d.rdb.SAdd(context.Background(), friendRelationCache+userID, IDList...).Err()
}

func (d *DataBases) ReduceFriendToCache(userID string, friendIDList ...string) error {
	var IDList []interface{}
	for _, id := range friendIDList {
		IDList = append(IDList, id)
	}
	return d.rdb.SRem(context.Background(), friendRelationCache+userID, IDList...).Err()
}

func (d *DataBases) GetFriendIDListFromCache(userID string) ([]string, error) {
	result := d.rdb.SMembers(context.Background(), friendRelationCache+userID)
	return result.Result()
}

func (d *DataBases) AddBlackUserToCache(userID string, blackList ...string) error {
	var IDList []interface{}
	for _, id := range blackList {
		IDList = append(IDList, id)
	}
	return d.rdb.SAdd(context.Background(), blackListCache+userID, IDList...).Err()
}

func (d *DataBases) ReduceBlackUserFromCache(userID string, blackList ...string) error {
	var IDList []interface{}
	for _, id := range blackList {
		IDList = append(IDList, id)
	}
	return d.rdb.SRem(context.Background(), blackListCache+userID, IDList...).Err()
}

func (d *DataBases) GetBlackListFromCache(userID string) ([]string, error) {
	result := d.rdb.SMembers(context.Background(), blackListCache+userID)
	return result.Result()
}

func (d *DataBases) AddGroupMemberToCache(groupID string, userIDList ...string) error {
	var IDList []interface{}
	for _, id := range userIDList {
		IDList = append(IDList, id)
	}
	return d.rdb.SAdd(context.Background(), groupCache+groupID, IDList...).Err()
}

func (d *DataBases) ReduceGroupMemberFromCache(groupID string, userIDList ...string) error {
	var IDList []interface{}
	for _, id := range userIDList {
		IDList = append(IDList, id)
	}
	return d.rdb.SRem(context.Background(), groupCache+groupID, IDList...).Err()
}

func (d *DataBases) GetGroupMemberIDListFromCache(groupID string) ([]string, error) {
	result := d.rdb.SMembers(context.Background(), groupCache+groupID)
	return result.Result()
}

func (d *DataBases) SetUserInfoToCache(userID string, m map[string]interface{}) error {
	return d.rdb.HSet(context.Background(), userInfoCache+userID, m).Err()
}

func (d *DataBases) GetUserInfoFromCache(userID string) (*pbCommon.UserInfo, error) {
	result, err := d.rdb.HGetAll(context.Background(), userInfoCache+userID).Result()
	bytes, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	userInfo := &pbCommon.UserInfo{}
	if err := proto.Unmarshal(bytes, userInfo); err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, userInfo)
	return userInfo, err
}
