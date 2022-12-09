package db

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	log2 "Open_IM/pkg/common/log"
	pbChat "Open_IM/pkg/proto/msg"
	pbRtc "Open_IM/pkg/proto/rtc"
	pbCommon "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	go_redis "github.com/go-redis/redis/v8"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

const (
	accountTempCode               = "ACCOUNT_TEMP_CODE"
	resetPwdTempCode              = "RESET_PWD_TEMP_CODE"
	userIncrSeq                   = "REDIS_USER_INCR_SEQ:" // user incr seq
	appleDeviceToken              = "DEVICE_TOKEN"
	userMinSeq                    = "REDIS_USER_MIN_SEQ:"
	uidPidToken                   = "UID_PID_TOKEN_STATUS:"
	conversationReceiveMessageOpt = "CON_RECV_MSG_OPT:"
	getuiToken                    = "GETUI_TOKEN"
	messageCache                  = "MESSAGE_CACHE:"
	SignalCache                   = "SIGNAL_CACHE:"
	SignalListCache               = "SIGNAL_LIST_CACHE:"
	GlobalMsgRecvOpt              = "GLOBAL_MSG_RECV_OPT"
	FcmToken                      = "FCM_TOKEN:"
	groupUserMinSeq               = "GROUP_USER_MIN_SEQ:"
	groupMaxSeq                   = "GROUP_MAX_SEQ:"
	groupMinSeq                   = "GROUP_MIN_SEQ:"
	sendMsgFailedFlag             = "SEND_MSG_FAILED_FLAG:"
	userBadgeUnreadCountSum       = "USER_BADGE_UNREAD_COUNT_SUM:"
	exTypeKeyLocker               = "EX_LOCK:"
)

func (d *DataBases) JudgeAccountEXISTS(account string) (bool, error) {
	key := accountTempCode + account
	n, err := d.RDB.Exists(context.Background(), key).Result()
	if n > 0 {
		return true, err
	} else {
		return false, err
	}
}
func (d *DataBases) SetAccountCode(account string, code, ttl int) (err error) {
	key := accountTempCode + account
	return d.RDB.Set(context.Background(), key, code, time.Duration(ttl)*time.Second).Err()
}
func (d *DataBases) GetAccountCode(account string) (string, error) {
	key := accountTempCode + account
	return d.RDB.Get(context.Background(), key).Result()
}

//Perform seq auto-increment operation of user messages
func (d *DataBases) IncrUserSeq(uid string) (uint64, error) {
	key := userIncrSeq + uid
	seq, err := d.RDB.Incr(context.Background(), key).Result()
	return uint64(seq), err
}

//Get the largest Seq
func (d *DataBases) GetUserMaxSeq(uid string) (uint64, error) {
	key := userIncrSeq + uid
	seq, err := d.RDB.Get(context.Background(), key).Result()
	return uint64(utils.StringToInt(seq)), err
}

//set the largest Seq
func (d *DataBases) SetUserMaxSeq(uid string, maxSeq uint64) error {
	key := userIncrSeq + uid
	return d.RDB.Set(context.Background(), key, maxSeq, 0).Err()
}

//Set the user's minimum seq
func (d *DataBases) SetUserMinSeq(uid string, minSeq uint32) (err error) {
	key := userMinSeq + uid
	return d.RDB.Set(context.Background(), key, minSeq, 0).Err()
}

//Get the smallest Seq
func (d *DataBases) GetUserMinSeq(uid string) (uint64, error) {
	key := userMinSeq + uid
	seq, err := d.RDB.Get(context.Background(), key).Result()
	return uint64(utils.StringToInt(seq)), err
}

func (d *DataBases) SetGroupUserMinSeq(groupID, userID string, minSeq uint64) (err error) {
	key := groupUserMinSeq + "g:" + groupID + "u:" + userID
	return d.RDB.Set(context.Background(), key, minSeq, 0).Err()
}
func (d *DataBases) GetGroupUserMinSeq(groupID, userID string) (uint64, error) {
	key := groupUserMinSeq + "g:" + groupID + "u:" + userID
	seq, err := d.RDB.Get(context.Background(), key).Result()
	return uint64(utils.StringToInt(seq)), err
}

func (d *DataBases) GetGroupMaxSeq(groupID string) (uint64, error) {
	key := groupMaxSeq + groupID
	seq, err := d.RDB.Get(context.Background(), key).Result()
	return uint64(utils.StringToInt(seq)), err
}

func (d *DataBases) IncrGroupMaxSeq(groupID string) (uint64, error) {
	key := groupMaxSeq + groupID
	seq, err := d.RDB.Incr(context.Background(), key).Result()
	return uint64(seq), err
}

func (d *DataBases) SetGroupMaxSeq(groupID string, maxSeq uint64) error {
	key := groupMaxSeq + groupID
	return d.RDB.Set(context.Background(), key, maxSeq, 0).Err()
}

func (d *DataBases) SetGroupMinSeq(groupID string, minSeq uint32) error {
	key := groupMinSeq + groupID
	return d.RDB.Set(context.Background(), key, minSeq, 0).Err()
}

//Store userid and platform class to redis
func (d *DataBases) AddTokenFlag(userID string, platformID int, token string, flag int) error {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
	log2.NewDebug("", "add token key is ", key)
	return d.RDB.HSet(context.Background(), key, token, flag).Err()
}

func (d *DataBases) GetTokenMapByUidPid(userID, platformID string) (map[string]int, error) {
	key := uidPidToken + userID + ":" + platformID
	log2.NewDebug("", "get token key is ", key)
	m, err := d.RDB.HGetAll(context.Background(), key).Result()
	mm := make(map[string]int)
	for k, v := range m {
		mm[k] = utils.StringToInt(v)
	}
	return mm, err
}
func (d *DataBases) SetTokenMapByUidPid(userID string, platformID int, m map[string]int) error {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
	mm := make(map[string]interface{})
	for k, v := range m {
		mm[k] = v
	}
	return d.RDB.HSet(context.Background(), key, mm).Err()
}
func (d *DataBases) DeleteTokenByUidPid(userID string, platformID int, fields []string) error {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
	return d.RDB.HDel(context.Background(), key, fields...).Err()
}
func (d *DataBases) SetSingleConversationRecvMsgOpt(userID, conversationID string, opt int32) error {
	key := conversationReceiveMessageOpt + userID
	return d.RDB.HSet(context.Background(), key, conversationID, opt).Err()
}

func (d *DataBases) GetSingleConversationRecvMsgOpt(userID, conversationID string) (int, error) {
	key := conversationReceiveMessageOpt + userID
	result, err := d.RDB.HGet(context.Background(), key, conversationID).Result()
	return utils.StringToInt(result), err
}
func (d *DataBases) SetUserGlobalMsgRecvOpt(userID string, opt int32) error {
	key := conversationReceiveMessageOpt + userID
	return d.RDB.HSet(context.Background(), key, GlobalMsgRecvOpt, opt).Err()
}
func (d *DataBases) GetUserGlobalMsgRecvOpt(userID string) (int, error) {
	key := conversationReceiveMessageOpt + userID
	result, err := d.RDB.HGet(context.Background(), key, GlobalMsgRecvOpt).Result()
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
	for _, v := range seqList {
		//MESSAGE_CACHE:169.254.225.224_reliability1653387820_0_1
		key := messageCache + userID + "_" + strconv.Itoa(int(v))
		result, err := d.RDB.Get(context.Background(), key).Result()
		if err != nil {
			errResult = err
			failedSeqList = append(failedSeqList, v)
			log2.Debug(operationID, "redis get message error: ", err.Error(), v)
		} else {
			msg := pbCommon.MsgData{}
			err = jsonpb.UnmarshalString(result, &msg)
			if err != nil {
				errResult = err
				failedSeqList = append(failedSeqList, v)
				log2.NewWarn(operationID, "Unmarshal err ", result, err.Error())
			} else {
				log2.NewDebug(operationID, "redis get msg is ", msg.String())
				seqMsg = append(seqMsg, &msg)
			}

		}
	}
	return seqMsg, failedSeqList, errResult
}

func (d *DataBases) SetMessageToCache(msgList []*pbChat.MsgDataToMQ, uid string, operationID string) (error, int) {
	ctx := context.Background()
	pipe := d.RDB.Pipeline()
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
		return errors.New(fmt.Sprintf("set msg to cache failed, failed lists: %q,%s", failedList, operationID)), len(failedList)
	}
	_, err := pipe.Exec(ctx)
	return err, 0
}
func (d *DataBases) DeleteMessageFromCache(msgList []*pbChat.MsgDataToMQ, uid string, operationID string) error {
	ctx := context.Background()
	for _, msg := range msgList {
		key := messageCache + uid + "_" + strconv.Itoa(int(msg.MsgData.Seq))
		err := d.RDB.Del(ctx, key).Err()
		if err != nil {
			log2.NewWarn(operationID, utils.GetSelfFuncName(), "redis failed", "args:", key, uid, err.Error(), msgList)
		}
	}
	return nil
}

func (d *DataBases) CleanUpOneUserAllMsgFromRedis(userID string, operationID string) error {
	ctx := context.Background()
	key := messageCache + userID + "_" + "*"
	vals, err := d.RDB.Keys(ctx, key).Result()
	log2.Debug(operationID, "vals: ", vals)
	if err == go_redis.Nil {
		return nil
	}
	if err != nil {
		return utils.Wrap(err, "")
	}
	for _, v := range vals {
		err = d.RDB.Del(ctx, v).Err()
	}
	return nil
}

func (d *DataBases) HandleSignalInfo(operationID string, msg *pbCommon.MsgData, pushToUserID string) (isSend bool, err error) {
	req := &pbRtc.SignalReq{}
	if err := proto.Unmarshal(msg.Content, req); err != nil {
		return false, err
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
		if !utils.IsContain(pushToUserID, inviteeUserIDList) {
			return false, nil
		}
	case *pbRtc.SignalReq_HungUp, *pbRtc.SignalReq_Cancel, *pbRtc.SignalReq_Reject, *pbRtc.SignalReq_Accept:
		return false, errors.New("signalInfo do not need offlinePush")
	default:
		log2.NewDebug(operationID, utils.GetSelfFuncName(), "req invalid type", string(msg.Content))
		return false, nil
	}
	if isInviteSignal {
		log2.NewDebug(operationID, utils.GetSelfFuncName(), "invite userID list:", inviteeUserIDList)
		for _, userID := range inviteeUserIDList {
			log2.NewInfo(operationID, utils.GetSelfFuncName(), "invite userID:", userID)
			timeout, err := strconv.Atoi(config.Config.Rtc.SignalTimeout)
			if err != nil {
				return false, err
			}
			keyList := SignalListCache + userID
			err = d.RDB.LPush(context.Background(), keyList, msg.ClientMsgID).Err()
			if err != nil {
				return false, err
			}
			err = d.RDB.Expire(context.Background(), keyList, time.Duration(timeout)*time.Second).Err()
			if err != nil {
				return false, err
			}
			key := SignalCache + msg.ClientMsgID
			err = d.RDB.Set(context.Background(), key, msg.Content, time.Duration(timeout)*time.Second).Err()
			if err != nil {
				return false, err
			}
		}
	}
	return true, nil
}

func (d *DataBases) GetSignalInfoFromCacheByClientMsgID(clientMsgID string) (invitationInfo *pbRtc.SignalInviteReq, err error) {
	key := SignalCache + clientMsgID
	invitationInfo = &pbRtc.SignalInviteReq{}
	bytes, err := d.RDB.Get(context.Background(), key).Bytes()
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
	result := d.RDB.LPop(context.Background(), keyList)
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
	err := d.RDB.Del(context.Background(), keyList).Err()
	return err
}

func (d *DataBases) DelMsgFromCache(uid string, seqList []uint32, operationID string) {
	for _, seq := range seqList {
		key := messageCache + uid + "_" + strconv.Itoa(int(seq))
		result, err := d.RDB.Get(context.Background(), key).Result()
		if err != nil {
			if err == go_redis.Nil {
				log2.NewDebug(operationID, utils.GetSelfFuncName(), err.Error(), "redis nil")
			} else {
				log2.NewError(operationID, utils.GetSelfFuncName(), err.Error(), key)
			}
			continue
		}
		var msg pbCommon.MsgData
		if err := utils.String2Pb(result, &msg); err != nil {
			log2.Error(operationID, utils.GetSelfFuncName(), "String2Pb failed", msg, result, key, err.Error())
			continue
		}
		msg.Status = constant.MsgDeleted
		s, err := utils.Pb2String(&msg)
		if err != nil {
			log2.Error(operationID, utils.GetSelfFuncName(), "Pb2String failed", msg, err.Error())
			continue
		}
		if err := d.RDB.Set(context.Background(), key, s, time.Duration(config.Config.MsgCacheTimeout)*time.Second).Err(); err != nil {
			log2.Error(operationID, utils.GetSelfFuncName(), "Set failed", err.Error())
		}
	}
}

func (d *DataBases) SetGetuiToken(token string, expireTime int64) error {
	return d.RDB.Set(context.Background(), getuiToken, token, time.Duration(expireTime)*time.Second).Err()
}

func (d *DataBases) GetGetuiToken() (string, error) {
	result, err := d.RDB.Get(context.Background(), getuiToken).Result()
	return result, err
}

func (d *DataBases) SetSendMsgStatus(status int32, operationID string) error {
	return d.RDB.Set(context.Background(), sendMsgFailedFlag+operationID, status, time.Hour*24).Err()
}

func (d *DataBases) GetSendMsgStatus(operationID string) (int, error) {
	result, err := d.RDB.Get(context.Background(), sendMsgFailedFlag+operationID).Result()
	if err != nil {
		return 0, err
	}
	status, err := strconv.Atoi(result)
	return status, err
}

func (d *DataBases) SetFcmToken(account string, platformID int, fcmToken string, expireTime int64) (err error) {
	key := FcmToken + account + ":" + strconv.Itoa(platformID)
	return d.RDB.Set(context.Background(), key, fcmToken, time.Duration(expireTime)*time.Second).Err()
}

func (d *DataBases) GetFcmToken(account string, platformID int) (string, error) {
	key := FcmToken + account + ":" + strconv.Itoa(platformID)
	return d.RDB.Get(context.Background(), key).Result()
}
func (d *DataBases) DelFcmToken(account string, platformID int) error {
	key := FcmToken + account + ":" + strconv.Itoa(platformID)
	return d.RDB.Del(context.Background(), key).Err()
}
func (d *DataBases) IncrUserBadgeUnreadCountSum(uid string) (int, error) {
	key := userBadgeUnreadCountSum + uid
	seq, err := d.RDB.Incr(context.Background(), key).Result()
	return int(seq), err
}
func (d *DataBases) SetUserBadgeUnreadCountSum(uid string, value int) error {
	key := userBadgeUnreadCountSum + uid
	return d.RDB.Set(context.Background(), key, value, 0).Err()
}
func (d *DataBases) GetUserBadgeUnreadCountSum(uid string) (int, error) {
	key := userBadgeUnreadCountSum + uid
	seq, err := d.RDB.Get(context.Background(), key).Result()
	return utils.StringToInt(seq), err
}
func (d *DataBases) JudgeMessageReactionEXISTS(clientMsgID string, sessionType int32) (bool, error) {
	key := getMessageReactionExPrefix(clientMsgID, sessionType)
	n, err := d.RDB.Exists(context.Background(), key).Result()
	if n > 0 {
		return true, err
	} else {
		return false, err
	}
}
func (d *DataBases) GetMessageTypeKeyValue(clientMsgID string, sessionType int32, typeKey string) (string, error) {
	key := getMessageReactionExPrefix(clientMsgID, sessionType)
	result, err := d.RDB.HGet(context.Background(), key, typeKey).Result()
	return result, err

}
func (d *DataBases) SetMessageTypeKeyValue(clientMsgID string, sessionType int32, typeKey, value string) error {
	key := getMessageReactionExPrefix(clientMsgID, sessionType)
	return d.RDB.HSet(context.Background(), key, typeKey, value).Err()

}
func (d *DataBases) LockMessageTypeKey(clientMsgID string, TypeKey string) error {
	key := exTypeKeyLocker + clientMsgID + "_" + TypeKey
	return d.RDB.SetNX(context.Background(), key, 1, time.Minute).Err()
}
func (d *DataBases) UnLockMessageTypeKey(clientMsgID string, TypeKey string) error {
	key := exTypeKeyLocker + clientMsgID + "_" + TypeKey
	return d.RDB.Del(context.Background(), key).Err()

}

func getMessageReactionExPrefix(clientMsgID string, sessionType int32) string {
	switch sessionType {
	case constant.SingleChatType:
		return "EX_SINGLE_" + clientMsgID
	case constant.GroupChatType:
		return "EX_GROUP_" + clientMsgID
	case constant.SuperGroupChatType:
		return "EX_SUPER_GROUP_" + clientMsgID
	case constant.NotificationChatType:
		return "EX_NOTIFICATION" + clientMsgID
	}
	return ""
}
