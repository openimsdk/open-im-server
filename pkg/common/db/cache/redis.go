package cache

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	pbChat "Open_IM/pkg/proto/msg"
	pbRtc "Open_IM/pkg/proto/rtc"
	"Open_IM/pkg/proto/sdkws"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

const (
	userIncrSeq      = "REDIS_USER_INCR_SEQ:" // user incr seq
	appleDeviceToken = "DEVICE_TOKEN"
	userMinSeq       = "REDIS_USER_MIN_SEQ:"

	getuiToken              = "GETUI_TOKEN"
	getuiTaskID             = "GETUI_TASK_ID"
	messageCache            = "MESSAGE_CACHE:"
	signalCache             = "SIGNAL_CACHE:"
	signalListCache         = "SIGNAL_LIST_CACHE:"
	FcmToken                = "FCM_TOKEN:"
	groupUserMinSeq         = "GROUP_USER_MIN_SEQ:"
	groupMaxSeq             = "GROUP_MAX_SEQ:"
	groupMinSeq             = "GROUP_MIN_SEQ:"
	sendMsgFailedFlag       = "SEND_MSG_FAILED_FLAG:"
	userBadgeUnreadCountSum = "USER_BADGE_UNREAD_COUNT_SUM:"
	exTypeKeyLocker         = "EX_LOCK:"

	uidPidToken = "UID_PID_TOKEN_STATUS:"
)

type Cache interface {
	IncrUserSeq(ctx context.Context, userID string) (int64, error)
	GetUserMaxSeq(ctx context.Context, userID string) (int64, error)
	SetUserMaxSeq(ctx context.Context, userID string, maxSeq int64) error
	SetUserMinSeq(ctx context.Context, userID string, minSeq int64) (err error)
	GetUserMinSeq(ctx context.Context, userID string) (int64, error)

	SetGroupUserMinSeq(ctx context.Context, groupID, userID string, minSeq int64) (err error)
	GetGroupUserMinSeq(ctx context.Context, groupID, userID string) (int64, error)
	GetGroupMaxSeq(ctx context.Context, groupID string) (int64, error)
	IncrGroupMaxSeq(ctx context.Context, groupID string) (int64, error)
	SetGroupMaxSeq(ctx context.Context, groupID string, maxSeq int64) error
	SetGroupMinSeq(ctx context.Context, groupID string, minSeq int64) error

	AddTokenFlag(ctx context.Context, userID string, platformID int, token string, flag int) error

	GetTokensWithoutError(ctx context.Context, userID, platformID string) (map[string]int, error)

	SetTokenMapByUidPid(ctx context.Context, userID string, platformID int, m map[string]int) error
	DeleteTokenByUidPid(ctx context.Context, userID string, platformID int, fields []string) error
	GetMessageListBySeq(ctx context.Context, userID string, seqList []int64) (seqMsg []*sdkws.MsgData, failedSeqList []int64, err error)
	SetMessageToCache(ctx context.Context, userID string, msgList []*pbChat.MsgDataToMQ) (int, error)
	DeleteMessageFromCache(ctx context.Context, userID string, msgList []*pbChat.MsgDataToMQ) error
	CleanUpOneUserAllMsg(ctx context.Context, userID string) error
	HandleSignalInfo(ctx context.Context, msg *sdkws.MsgData, pushToUserID string) (isSend bool, err error)
	GetSignalInfoFromCacheByClientMsgID(ctx context.Context, clientMsgID string) (invitationInfo *pbRtc.SignalInviteReq, err error)
	GetAvailableSignalInvitationInfo(ctx context.Context, userID string) (invitationInfo *pbRtc.SignalInviteReq, err error)
	DelUserSignalList(ctx context.Context, userID string) error
	DelMsgFromCache(ctx context.Context, userID string, seqList []int64) error

	SetGetuiToken(ctx context.Context, token string, expireTime int64) error
	GetGetuiToken(ctx context.Context) (string, error)
	SetGetuiTaskID(ctx context.Context, taskID string, expireTime int64) error
	GetGetuiTaskID(ctx context.Context) (string, error)

	SetSendMsgStatus(ctx context.Context, status int32) error
	GetSendMsgStatus(ctx context.Context) (int, error)
	SetFcmToken(ctx context.Context, account string, platformID int, fcmToken string, expireTime int64) (err error)
	GetFcmToken(ctx context.Context, account string, platformID int) (string, error)
	DelFcmToken(ctx context.Context, account string, platformID int) error
	IncrUserBadgeUnreadCountSum(ctx context.Context, userID string) (int, error)
	SetUserBadgeUnreadCountSum(ctx context.Context, userID string, value int) error
	GetUserBadgeUnreadCountSum(ctx context.Context, userID string) (int, error)
	JudgeMessageReactionEXISTS(ctx context.Context, clientMsgID string, sessionType int32) (bool, error)
	GetOneMessageAllReactionList(ctx context.Context, clientMsgID string, sessionType int32) (map[string]string, error)
	DeleteOneMessageKey(ctx context.Context, clientMsgID string, sessionType int32, subKey string) error
	SetMessageReactionExpire(ctx context.Context, clientMsgID string, sessionType int32, expiration time.Duration) (bool, error)
	GetMessageTypeKeyValue(ctx context.Context, clientMsgID string, sessionType int32, typeKey string) (string, error)
	SetMessageTypeKeyValue(ctx context.Context, clientMsgID string, sessionType int32, typeKey, value string) error
	LockMessageTypeKey(ctx context.Context, clientMsgID string, TypeKey string) error
	UnLockMessageTypeKey(ctx context.Context, clientMsgID string, TypeKey string) error
}

// native redis operate

//func NewRedis() *RedisClient {
//	o := &RedisClient{}
//	o.InitRedis()
//	return o
//}

func NewRedis() (*RedisClient, error) {
	var rdb redis.UniversalClient
	if config.Config.Redis.EnableCluster {
		rdb = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    config.Config.Redis.DBAddress,
			Username: config.Config.Redis.DBUserName,
			Password: config.Config.Redis.DBPassWord, // no password set
			PoolSize: 50,
		})
		//if err := rdb.Ping(ctx).Err();err != nil {
		//	return nil, fmt.Errorf("redis ping %w", err)
		//}
		//return &RedisClient{rdb: rdb}, nil
	} else {
		rdb = redis.NewClient(&redis.Options{
			Addr:     config.Config.Redis.DBAddress[0],
			Username: config.Config.Redis.DBUserName,
			Password: config.Config.Redis.DBPassWord, // no password set
			DB:       0,                              // use default DB
			PoolSize: 100,                            // 连接池大小
		})
		//err := rdb.Ping(ctx).Err()
		//if err != nil {
		//	panic(err.Error() + " redis " + config.Config.Redis.DBAddress[0] + config.Config.Redis.DBUserName + config.Config.Redis.DBPassWord)
		//}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err := rdb.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("redis ping %w", err)
	}
	return &RedisClient{rdb: rdb}, nil
}

type RedisClient struct {
	rdb redis.UniversalClient
}

func NewRedisClient(rdb redis.UniversalClient) *RedisClient {
	return &RedisClient{rdb: rdb}
}

func (r *RedisClient) GetClient() redis.UniversalClient {
	return r.rdb
}

// Perform seq auto-increment operation of user messages
func (r *RedisClient) IncrUserSeq(ctx context.Context, uid string) (int64, error) {
	key := userIncrSeq + uid
	seq, err := r.rdb.Incr(context.Background(), key).Result()
	return seq, err
}

// Get the largest Seq
func (r *RedisClient) GetUserMaxSeq(ctx context.Context, uid string) (int64, error) {
	key := userIncrSeq + uid
	seq, err := r.rdb.Get(context.Background(), key).Result()
	return int64(utils.StringToInt(seq)), err
}

// set the largest Seq
func (r *RedisClient) SetUserMaxSeq(ctx context.Context, uid string, maxSeq int64) error {
	key := userIncrSeq + uid
	return r.rdb.Set(context.Background(), key, maxSeq, 0).Err()
}

// Set the user's minimum seq
func (r *RedisClient) SetUserMinSeq(ctx context.Context, uid string, minSeq int64) (err error) {
	key := userMinSeq + uid
	return r.rdb.Set(context.Background(), key, minSeq, 0).Err()
}

// Get the smallest Seq
func (r *RedisClient) GetUserMinSeq(ctx context.Context, uid string) (int64, error) {
	key := userMinSeq + uid
	seq, err := r.rdb.Get(context.Background(), key).Result()
	return int64(utils.StringToInt(seq)), err
}

func (r *RedisClient) SetGroupUserMinSeq(ctx context.Context, groupID, userID string, minSeq int64) (err error) {
	key := groupUserMinSeq + "g:" + groupID + "u:" + userID
	return r.rdb.Set(context.Background(), key, minSeq, 0).Err()
}
func (r *RedisClient) GetGroupUserMinSeq(ctx context.Context, groupID, userID string) (int64, error) {
	key := groupUserMinSeq + "g:" + groupID + "u:" + userID
	seq, err := r.rdb.Get(context.Background(), key).Result()
	return int64(utils.StringToInt(seq)), err
}

func (r *RedisClient) GetGroupMaxSeq(ctx context.Context, groupID string) (int64, error) {
	key := groupMaxSeq + groupID
	seq, err := r.rdb.Get(context.Background(), key).Result()
	return int64(utils.StringToInt(seq)), err
}

func (r *RedisClient) IncrGroupMaxSeq(ctx context.Context, groupID string) (int64, error) {
	key := groupMaxSeq + groupID
	seq, err := r.rdb.Incr(context.Background(), key).Result()
	return seq, err
}

func (r *RedisClient) SetGroupMaxSeq(ctx context.Context, groupID string, maxSeq int64) error {
	key := groupMaxSeq + groupID
	return r.rdb.Set(context.Background(), key, maxSeq, 0).Err()
}

func (r *RedisClient) SetGroupMinSeq(ctx context.Context, groupID string, minSeq int64) error {
	key := groupMinSeq + groupID
	return r.rdb.Set(context.Background(), key, minSeq, 0).Err()
}

// Store userid and platform class to redis
func (r *RedisClient) AddTokenFlag(ctx context.Context, userID string, platform string, token string, flag int) error {
	key := uidPidToken + userID + ":" + platform
	return r.rdb.HSet(context.Background(), key, token, flag).Err()
}

//key:userID+platform-> <token, flag>
func (r *RedisClient) GetTokenMapByUidPid(ctx context.Context, userID, platformID string) (map[string]int, error) {
	key := uidPidToken + userID + ":" + platformID
	m, err := r.rdb.HGetAll(context.Background(), key).Result()
	mm := make(map[string]int)
	for k, v := range m {
		mm[k] = utils.StringToInt(v)
	}
	return mm, err
}

func (r *RedisClient) GetTokensWithoutError(ctx context.Context, userID, platform string) (map[string]int, error) {
	key := uidPidToken + userID + ":" + platform
	m, err := r.rdb.HGetAll(context.Background(), key).Result()
	if err != nil && err == redis.Nil {
		return nil, nil
	}
	mm := make(map[string]int)
	for k, v := range m {
		mm[k] = utils.StringToInt(v)
	}
	return mm, utils.Wrap(err, "")
}

func (r *RedisClient) SetTokenMapByUidPid(ctx context.Context, userID string, platform string, m map[string]int) error {
	key := uidPidToken + userID + ":" + platform
	mm := make(map[string]interface{})
	for k, v := range m {
		mm[k] = v
	}
	return r.rdb.HSet(context.Background(), key, mm).Err()
}

func (r *RedisClient) DeleteTokenByUidPid(ctx context.Context, userID string, platform string, fields []string) error {
	key := uidPidToken + userID + ":" + platform
	return r.rdb.HDel(context.Background(), key, fields...).Err()
}

func (r *RedisClient) GetMessageListBySeq(ctx context.Context, userID string, seqList []int64, operationID string) (seqMsg []*sdkws.MsgData, failedSeqList []int64, err2 error) {
	for _, v := range seqList {
		//MESSAGE_CACHE:169.254.225.224_reliability1653387820_0_1
		key := messageCache + userID + "_" + strconv.Itoa(int(v))
		result, err := r.rdb.Get(context.Background(), key).Result()
		if err != nil {
			if err != redis.Nil {
				err2 = err
			}
			failedSeqList = append(failedSeqList, v)
		} else {
			msg := sdkws.MsgData{}
			err = jsonpb.UnmarshalString(result, &msg)
			if err != nil {
				err2 = err
				failedSeqList = append(failedSeqList, v)
			} else {
				seqMsg = append(seqMsg, &msg)
			}
		}
	}
	return seqMsg, failedSeqList, err2
}

func (r *RedisClient) SetMessageToCache(ctx context.Context, userID string, msgList []*pbChat.MsgDataToMQ, uid string) (int, error) {
	pipe := r.rdb.Pipeline()
	var failedList []pbChat.MsgDataToMQ
	for _, msg := range msgList {
		key := messageCache + uid + "_" + strconv.Itoa(int(msg.MsgData.Seq))
		s, err := utils.Pb2String(msg.MsgData)
		if err != nil {
			continue
		}
		err = pipe.Set(ctx, key, s, time.Duration(config.Config.MsgCacheTimeout)*time.Second).Err()
		//err = r.rdb.HMSet(context.Background(), "12", map[string]interface{}{"1": 2, "343": false}).Err()
		if err != nil {
			failedList = append(failedList, *msg)
		}
	}
	if len(failedList) != 0 {
		return len(failedList), errors.New(fmt.Sprintf("set msg to cache failed, failed lists: %q,%s", failedList))
	}
	_, err := pipe.Exec(ctx)
	return 0, err
}
func (r *RedisClient) DeleteMessageFromCache(ctx context.Context, userID string, msgList []*pbChat.MsgDataToMQ) error {
	for _, msg := range msgList {
		key := messageCache + userID + "_" + strconv.Itoa(int(msg.MsgData.Seq))
		err := r.rdb.Del(ctx, key).Err()
		if err != nil {
		}
	}
	return nil
}

func (r *RedisClient) CleanUpOneUserAllMsg(ctx context.Context, userID string) error {
	key := messageCache + userID + "_" + "*"
	vals, err := r.rdb.Keys(ctx, key).Result()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return utils.Wrap(err, "")
	}
	for _, v := range vals {
		err = r.rdb.Del(ctx, v).Err()
	}
	return nil
}

func (r *RedisClient) HandleSignalInfo(ctx context.Context, operationID string, msg *sdkws.MsgData, pushToUserID string) (isSend bool, err error) {
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
		return false, nil
	default:
		return false, nil
	}
	if isInviteSignal {
		for _, userID := range inviteeUserIDList {
			timeout, err := strconv.Atoi(config.Config.Rtc.SignalTimeout)
			if err != nil {
				return false, err
			}
			keyList := signalListCache + userID
			err = r.rdb.LPush(context.Background(), keyList, msg.ClientMsgID).Err()
			if err != nil {
				return false, err
			}
			err = r.rdb.Expire(context.Background(), keyList, time.Duration(timeout)*time.Second).Err()
			if err != nil {
				return false, err
			}
			key := signalCache + msg.ClientMsgID
			err = r.rdb.Set(context.Background(), key, msg.Content, time.Duration(timeout)*time.Second).Err()
			if err != nil {
				return false, err
			}
		}
	}
	return true, nil
}

func (r *RedisClient) GetSignalInfoFromCacheByClientMsgID(ctx context.Context, clientMsgID string) (invitationInfo *pbRtc.SignalInviteReq, err error) {
	key := signalCache + clientMsgID
	invitationInfo = &pbRtc.SignalInviteReq{}
	bytes, err := r.rdb.Get(context.Background(), key).Bytes()
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

func (r *RedisClient) GetAvailableSignalInvitationInfo(ctx context.Context, userID string) (invitationInfo *pbRtc.SignalInviteReq, err error) {
	keyList := signalListCache + userID
	result := r.rdb.LPop(context.Background(), keyList)
	if err = result.Err(); err != nil {
		return nil, utils.Wrap(err, "GetAvailableSignalInvitationInfo failed")
	}
	key, err := result.Result()
	if err != nil {
		return nil, utils.Wrap(err, "GetAvailableSignalInvitationInfo failed")
	}
	invitationInfo, err = r.GetSignalInfoFromCacheByClientMsgID(ctx, key)
	if err != nil {
		return nil, utils.Wrap(err, "GetSignalInfoFromCacheByClientMsgID")
	}
	err = r.DelUserSignalList(ctx, userID)
	if err != nil {
		return nil, utils.Wrap(err, "GetSignalInfoFromCacheByClientMsgID")
	}
	return invitationInfo, nil
}

func (r *RedisClient) DelUserSignalList(ctx context.Context, userID string) error {
	keyList := signalListCache + userID
	err := r.rdb.Del(context.Background(), keyList).Err()
	return err
}

func (r *RedisClient) DelMsgFromCache(ctx context.Context, uid string, seqList []int64, operationID string) {
	for _, seq := range seqList {
		key := messageCache + uid + "_" + strconv.Itoa(int(seq))
		result, err := r.rdb.Get(context.Background(), key).Result()
		if err != nil {
			if err == redis.Nil {
			} else {
			}
			continue
		}
		var msg sdkws.MsgData
		if err := utils.String2Pb(result, &msg); err != nil {
			continue
		}
		msg.Status = constant.MsgDeleted
		s, err := utils.Pb2String(&msg)
		if err != nil {
			continue
		}
		if err := r.rdb.Set(context.Background(), key, s, time.Duration(config.Config.MsgCacheTimeout)*time.Second).Err(); err != nil {
		}
	}
}

func (r *RedisClient) SetGetuiToken(ctx context.Context, token string, expireTime int64) error {
	return r.rdb.Set(context.Background(), getuiToken, token, time.Duration(expireTime)*time.Second).Err()
}

func (r *RedisClient) GetGetuiToken(ctx context.Context) (string, error) {
	result, err := r.rdb.Get(context.Background(), getuiToken).Result()
	return result, err
}

func (r *RedisClient) SetGetuiTaskID(ctx context.Context, taskID string, expireTime int64) error {
	return r.rdb.Set(context.Background(), getuiTaskID, taskID, time.Duration(expireTime)*time.Second).Err()
}

func (r *RedisClient) GetGetuiTaskID(ctx context.Context) (string, error) {
	result, err := r.rdb.Get(context.Background(), getuiTaskID).Result()
	return result, err
}

func (r *RedisClient) SetSendMsgStatus(ctx context.Context, status int32, operationID string) error {
	return r.rdb.Set(context.Background(), sendMsgFailedFlag+operationID, status, time.Hour*24).Err()
}

func (r *RedisClient) GetSendMsgStatus(ctx context.Context, operationID string) (int, error) {
	result, err := r.rdb.Get(context.Background(), sendMsgFailedFlag+operationID).Result()
	if err != nil {
		return 0, err
	}
	status, err := strconv.Atoi(result)
	return status, err
}

func (r *RedisClient) SetFcmToken(ctx context.Context, account string, platformID int, fcmToken string, expireTime int64) (err error) {
	key := FcmToken + account + ":" + strconv.Itoa(platformID)
	return r.rdb.Set(context.Background(), key, fcmToken, time.Duration(expireTime)*time.Second).Err()
}

func (r *RedisClient) GetFcmToken(ctx context.Context, account string, platformID int) (string, error) {
	key := FcmToken + account + ":" + strconv.Itoa(platformID)
	return r.rdb.Get(context.Background(), key).Result()
}
func (r *RedisClient) DelFcmToken(ctx context.Context, account string, platformID int) error {
	key := FcmToken + account + ":" + strconv.Itoa(platformID)
	return r.rdb.Del(context.Background(), key).Err()
}
func (r *RedisClient) IncrUserBadgeUnreadCountSum(ctx context.Context, uid string) (int, error) {
	key := userBadgeUnreadCountSum + uid
	seq, err := r.rdb.Incr(context.Background(), key).Result()
	return int(seq), err
}
func (r *RedisClient) SetUserBadgeUnreadCountSum(ctx context.Context, uid string, value int) error {
	key := userBadgeUnreadCountSum + uid
	return r.rdb.Set(context.Background(), key, value, 0).Err()
}
func (r *RedisClient) GetUserBadgeUnreadCountSum(ctx context.Context, uid string) (int, error) {
	key := userBadgeUnreadCountSum + uid
	seq, err := r.rdb.Get(context.Background(), key).Result()
	return utils.StringToInt(seq), err
}
func (r *RedisClient) JudgeMessageReactionEXISTS(ctx context.Context, clientMsgID string, sessionType int32) (bool, error) {
	key := r.getMessageReactionExPrefix(clientMsgID, sessionType)
	n, err := r.rdb.Exists(context.Background(), key).Result()
	return n > 0, err
}

func (r *RedisClient) GetOneMessageAllReactionList(ctx context.Context, clientMsgID string, sessionType int32) (map[string]string, error) {
	key := r.getMessageReactionExPrefix(clientMsgID, sessionType)
	return r.rdb.HGetAll(context.Background(), key).Result()

}
func (r *RedisClient) DeleteOneMessageKey(ctx context.Context, clientMsgID string, sessionType int32, subKey string) error {
	key := r.getMessageReactionExPrefix(clientMsgID, sessionType)
	return r.rdb.HDel(context.Background(), key, subKey).Err()

}
func (r *RedisClient) SetMessageReactionExpire(ctx context.Context, clientMsgID string, sessionType int32, expiration time.Duration) (bool, error) {
	key := r.getMessageReactionExPrefix(clientMsgID, sessionType)
	return r.rdb.Expire(context.Background(), key, expiration).Result()
}

func (r *RedisClient) GetMessageTypeKeyValue(ctx context.Context, clientMsgID string, sessionType int32, typeKey string) (string, error) {
	key := r.getMessageReactionExPrefix(clientMsgID, sessionType)
	result, err := r.rdb.HGet(context.Background(), key, typeKey).Result()
	return result, err
}

func (r *RedisClient) SetMessageTypeKeyValue(ctx context.Context, clientMsgID string, sessionType int32, typeKey, value string) error {
	key := r.getMessageReactionExPrefix(clientMsgID, sessionType)
	return r.rdb.HSet(context.Background(), key, typeKey, value).Err()

}

func (r *RedisClient) LockMessageTypeKey(ctx context.Context, clientMsgID string, TypeKey string) error {
	key := exTypeKeyLocker + clientMsgID + "_" + TypeKey
	return r.rdb.SetNX(context.Background(), key, 1, time.Minute).Err()
}
func (r *RedisClient) UnLockMessageTypeKey(ctx context.Context, clientMsgID string, TypeKey string) error {
	key := exTypeKeyLocker + clientMsgID + "_" + TypeKey
	return r.rdb.Del(context.Background(), key).Err()

}

func (r *RedisClient) getMessageReactionExPrefix(clientMsgID string, sessionType int32) string {
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
