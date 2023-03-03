package cache

import (
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/tracelog"
	pbMsg "OpenIM/pkg/proto/msg"
	"OpenIM/pkg/proto/sdkws"
	"OpenIM/pkg/utils"
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	userIncrSeq             = "REDIS_USER_INCR_SEQ:" // user incr seq
	appleDeviceToken        = "DEVICE_TOKEN"
	userMinSeq              = "REDIS_USER_MIN_SEQ:"
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
	uidPidToken             = "UID_PID_TOKEN_STATUS:"
)

type Model interface {
	IncrUserSeq(ctx context.Context, userID string) (int64, error)
	GetUserMaxSeq(ctx context.Context, userID string) (int64, error)
	SetUserMaxSeq(ctx context.Context, userID string, maxSeq int64) error
	SetUserMinSeq(ctx context.Context, userID string, minSeq int64) (err error)
	GetUserMinSeq(ctx context.Context, userID string) (int64, error)
	SetGroupUserMinSeq(ctx context.Context, groupID, userID string, minSeq int64) (err error)
	GetGroupUserMinSeq(ctx context.Context, groupID, userID string) (int64, error)
	GetGroupMaxSeq(ctx context.Context, groupID string) (int64, error)
	GetGroupMinSeq(ctx context.Context, groupID string) (int64, error)
	IncrGroupMaxSeq(ctx context.Context, groupID string) (int64, error)
	SetGroupMaxSeq(ctx context.Context, groupID string, maxSeq int64) error
	SetGroupMinSeq(ctx context.Context, groupID string, minSeq int64) error
	AddTokenFlag(ctx context.Context, userID string, platformID int, token string, flag int) error
	GetTokensWithoutError(ctx context.Context, userID, platformID string) (map[string]int, error)
	SetTokenMapByUidPid(ctx context.Context, userID string, platform string, m map[string]int) error
	DeleteTokenByUidPid(ctx context.Context, userID string, platform string, fields []string) error
	GetMessagesBySeq(ctx context.Context, userID string, seqList []int64) (seqMsg []*sdkws.MsgData, failedSeqList []int64, err error)
	SetMessageToCache(ctx context.Context, userID string, msgList []*pbMsg.MsgDataToMQ) (int, error)
	DeleteMessageFromCache(ctx context.Context, userID string, msgList []*pbMsg.MsgDataToMQ) error
	CleanUpOneUserAllMsg(ctx context.Context, userID string) error
	HandleSignalInvite(ctx context.Context, msg *sdkws.MsgData, pushToUserID string) (isSend bool, err error)
	GetSignalInvitationInfoByClientMsgID(ctx context.Context, clientMsgID string) (invitationInfo *sdkws.SignalInviteReq, err error)
	GetAvailableSignalInvitationInfo(ctx context.Context, userID string) (invitationInfo *sdkws.SignalInviteReq, err error)
	DelUserSignalList(ctx context.Context, userID string) error
	DelMsgFromCache(ctx context.Context, userID string, seqList []int64) error
	SetGetuiToken(ctx context.Context, token string, expireTime int64) error
	GetGetuiToken(ctx context.Context) (string, error)
	SetGetuiTaskID(ctx context.Context, taskID string, expireTime int64) error
	GetGetuiTaskID(ctx context.Context) (string, error)
	SetSendMsgStatus(ctx context.Context, id string, status int32) error
	GetSendMsgStatus(ctx context.Context, id string) (int32, error)
	SetFcmToken(ctx context.Context, account string, platformID int, fcmToken string, expireTime int64) (err error)
	GetFcmToken(ctx context.Context, account string, platformID int) (string, error)
	DelFcmToken(ctx context.Context, account string, platformID int) error
	IncrUserBadgeUnreadCountSum(ctx context.Context, userID string) (int, error)
	SetUserBadgeUnreadCountSum(ctx context.Context, userID string, value int) error
	GetUserBadgeUnreadCountSum(ctx context.Context, userID string) (int, error)
	JudgeMessageReactionExist(ctx context.Context, clientMsgID string, sessionType int32) (bool, error)
	GetOneMessageAllReactionList(ctx context.Context, clientMsgID string, sessionType int32) (map[string]string, error)
	DeleteOneMessageKey(ctx context.Context, clientMsgID string, sessionType int32, subKey string) error
	SetMessageReactionExpire(ctx context.Context, clientMsgID string, sessionType int32, expiration time.Duration) (bool, error)
	GetMessageTypeKeyValue(ctx context.Context, clientMsgID string, sessionType int32, typeKey string) (string, error)
	SetMessageTypeKeyValue(ctx context.Context, clientMsgID string, sessionType int32, typeKey, value string) error
	LockMessageTypeKey(ctx context.Context, clientMsgID string, TypeKey string) error
	UnLockMessageTypeKey(ctx context.Context, clientMsgID string, TypeKey string) error
}

func NewCacheModel(client redis.UniversalClient) Model {
	return &cache{rdb: client}
}

type cache struct {
	rdb redis.UniversalClient
}

func (c *cache) IncrUserSeq(ctx context.Context, userID string) (int64, error) {
	return utils.Wrap2(c.rdb.Get(ctx, userIncrSeq+userID).Int64())
}

func (c *cache) GetUserMaxSeq(ctx context.Context, userID string) (int64, error) {
	return utils.Wrap2(c.rdb.Get(ctx, userIncrSeq+userID).Int64())
}

func (c *cache) SetUserMaxSeq(ctx context.Context, userID string, maxSeq int64) error {
	return utils.Wrap1(c.rdb.Set(ctx, userIncrSeq+userID, maxSeq, 0).Err())
}

func (c *cache) SetUserMinSeq(ctx context.Context, userID string, minSeq int64) (err error) {
	return utils.Wrap1(c.rdb.Set(ctx, userMinSeq+userID, minSeq, 0).Err())
}

func (c *cache) GetUserMinSeq(ctx context.Context, userID string) (int64, error) {
	return utils.Wrap2(c.rdb.Get(ctx, userMinSeq+userID).Int64())
}

func (c *cache) SetGroupUserMinSeq(ctx context.Context, groupID, userID string, minSeq int64) (err error) {
	key := groupUserMinSeq + "g:" + groupID + "u:" + userID
	return utils.Wrap1(c.rdb.Set(ctx, key, minSeq, 0).Err())
}

func (c *cache) GetGroupUserMinSeq(ctx context.Context, groupID, userID string) (int64, error) {
	key := groupUserMinSeq + "g:" + groupID + "u:" + userID
	return utils.Wrap2(c.rdb.Get(ctx, key).Int64())
}

func (c *cache) GetGroupMaxSeq(ctx context.Context, groupID string) (int64, error) {
	return utils.Wrap2(c.rdb.Get(ctx, groupMaxSeq+groupID).Int64())
}

func (c *cache) GetGroupMinSeq(ctx context.Context, groupID string) (int64, error) {
	return utils.Wrap2(c.rdb.Get(ctx, groupMinSeq+groupID).Int64())
}

func (c *cache) IncrGroupMaxSeq(ctx context.Context, groupID string) (int64, error) {
	key := groupMaxSeq + groupID
	seq, err := c.rdb.Incr(ctx, key).Uint64()
	return int64(seq), utils.Wrap1(err)
}

func (c *cache) SetGroupMaxSeq(ctx context.Context, groupID string, maxSeq int64) error {
	key := groupMaxSeq + groupID
	return utils.Wrap1(c.rdb.Set(ctx, key, maxSeq, 0).Err())
}

func (c *cache) SetGroupMinSeq(ctx context.Context, groupID string, minSeq int64) error {
	key := groupMinSeq + groupID
	return utils.Wrap1(c.rdb.Set(ctx, key, minSeq, 0).Err())
}

func (c *cache) AddTokenFlag(ctx context.Context, userID string, platformID int, token string, flag int) error {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
	return utils.Wrap1(c.rdb.HSet(ctx, key, token, flag).Err())
}

func (c *cache) GetTokensWithoutError(ctx context.Context, userID, platformID string) (map[string]int, error) {
	key := uidPidToken + userID + ":" + platformID
	m, err := c.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, utils.Wrap1(err)
	}
	mm := make(map[string]int)
	for k, v := range m {
		mm[k] = utils.StringToInt(v)
	}
	return mm, nil
}

func (c *cache) SetTokenMapByUidPid(ctx context.Context, userID string, platform string, m map[string]int) error {
	key := uidPidToken + userID + ":" + platform
	mm := make(map[string]interface{})
	for k, v := range m {
		mm[k] = v
	}
	return utils.Wrap1(c.rdb.HSet(ctx, key, mm).Err())
}

func (c *cache) DeleteTokenByUidPid(ctx context.Context, userID string, platform string, fields []string) error {
	key := uidPidToken + userID + ":" + platform
	return utils.Wrap1(c.rdb.HDel(ctx, key, fields...).Err())
}

func (c *cache) GetMessagesBySeq(ctx context.Context, userID string, seqList []int64) (seqMsg []*sdkws.MsgData, failedSeqList []int64, err error) {
	var errResult error
	for _, v := range seqList {
		//MESSAGE_CACHE:169.254.225.224_reliability1653387820_0_1
		key := messageCache + userID + "_" + strconv.Itoa(int(v))
		result, err := c.rdb.Get(ctx, key).Result()
		if err != nil {
			errResult = err
			failedSeqList = append(failedSeqList, v)
		} else {
			msg := sdkws.MsgData{}
			err = jsonpb.UnmarshalString(result, &msg)
			if err != nil {
				errResult = err
				failedSeqList = append(failedSeqList, v)
			} else {
				seqMsg = append(seqMsg, &msg)
			}

		}
	}
	return seqMsg, failedSeqList, errResult
}

func (c *cache) SetMessageToCache(ctx context.Context, userID string, msgList []*pbMsg.MsgDataToMQ) (int, error) {
	pipe := c.rdb.Pipeline()
	var failedMsgs []pbMsg.MsgDataToMQ
	for _, msg := range msgList {
		key := messageCache + userID + "_" + strconv.Itoa(int(msg.MsgData.Seq))
		s, err := utils.Pb2String(msg.MsgData)
		if err != nil {
			return 0, utils.Wrap1(err)
		}
		err = pipe.Set(ctx, key, s, time.Duration(config.Config.MsgCacheTimeout)*time.Second).Err()
		if err != nil {
			return 0, utils.Wrap1(err)
		}
	}
	if len(failedMsgs) != 0 {
		return len(failedMsgs), errors.New(fmt.Sprintf("set msg to cache failed, failed lists: %q,%s", failedMsgs, tracelog.GetOperationID(ctx)))
	}
	_, err := pipe.Exec(ctx)
	return 0, err
}

func (c *cache) DeleteMessageFromCache(ctx context.Context, userID string, msgList []*pbMsg.MsgDataToMQ) error {
	for _, v := range msgList {
		if err := c.rdb.Del(ctx, messageCache+userID+"_"+strconv.Itoa(int(v.MsgData.Seq))).Err(); err != nil {
			return utils.Wrap1(err)
		}
	}
	return nil
}

func (c *cache) CleanUpOneUserAllMsg(ctx context.Context, userID string) error {
	key := messageCache + userID + "_" + "*"
	vals, err := c.rdb.Keys(ctx, key).Result()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return utils.Wrap1(err)
	}
	for _, v := range vals {
		if err := c.rdb.Del(ctx, v).Err(); err != nil {
			return utils.Wrap1(err)
		}
	}
	return nil
}

func (c *cache) HandleSignalInvite(ctx context.Context, msg *sdkws.MsgData, pushToUserID string) (isSend bool, err error) {
	req := &sdkws.SignalReq{}
	if err := proto.Unmarshal(msg.Content, req); err != nil {
		return false, utils.Wrap1(err)
	}
	var inviteeUserIDs []string
	var isInviteSignal bool
	switch signalInfo := req.Payload.(type) {
	case *sdkws.SignalReq_Invite:
		inviteeUserIDs = signalInfo.Invite.Invitation.InviteeUserIDList
		isInviteSignal = true
	case *sdkws.SignalReq_InviteInGroup:
		inviteeUserIDs = signalInfo.InviteInGroup.Invitation.InviteeUserIDList
		isInviteSignal = true
		if !utils.Contain(pushToUserID, inviteeUserIDs...) {
			return false, nil
		}
	case *sdkws.SignalReq_HungUp, *sdkws.SignalReq_Cancel, *sdkws.SignalReq_Reject, *sdkws.SignalReq_Accept:
		return false, utils.Wrap1(errors.New("signalInfo do not need offlinePush"))
	default:
		return false, nil
	}
	if isInviteSignal {
		for _, userID := range inviteeUserIDs {
			timeout, err := strconv.Atoi(config.Config.Rtc.SignalTimeout)
			if err != nil {
				return false, utils.Wrap1(err)
			}
			keyList := signalListCache + userID
			err = c.rdb.LPush(ctx, keyList, msg.ClientMsgID).Err()
			if err != nil {
				return false, utils.Wrap1(err)
			}
			err = c.rdb.Expire(ctx, keyList, time.Duration(timeout)*time.Second).Err()
			if err != nil {
				return false, utils.Wrap1(err)
			}
			key := signalCache + msg.ClientMsgID
			err = c.rdb.Set(ctx, key, msg.Content, time.Duration(timeout)*time.Second).Err()
			if err != nil {
				return false, utils.Wrap1(err)
			}
		}
	}
	return true, nil
}

func (c *cache) GetSignalInvitationInfoByClientMsgID(ctx context.Context, clientMsgID string) (signalInviteReq *sdkws.SignalInviteReq, err error) {
	bytes, err := c.rdb.Get(ctx, signalCache+clientMsgID).Bytes()
	if err != nil {
		return nil, utils.Wrap1(err)
	}
	signalReq := &sdkws.SignalReq{}
	if err = proto.Unmarshal(bytes, signalReq); err != nil {
		return nil, utils.Wrap1(err)
	}
	signalInviteReq = &sdkws.SignalInviteReq{}
	switch req := signalReq.Payload.(type) {
	case *sdkws.SignalReq_Invite:
		signalInviteReq.Invitation = req.Invite.Invitation
		signalInviteReq.OpUserID = req.Invite.OpUserID
	case *sdkws.SignalReq_InviteInGroup:
		signalInviteReq.Invitation = req.InviteInGroup.Invitation
		signalInviteReq.OpUserID = req.InviteInGroup.OpUserID
	}
	return signalInviteReq, nil
}

func (c *cache) GetAvailableSignalInvitationInfo(ctx context.Context, userID string) (invitationInfo *sdkws.SignalInviteReq, err error) {
	key, err := c.rdb.LPop(ctx, signalListCache+userID).Result()
	if err != nil {
		return nil, utils.Wrap1(err)
	}
	invitationInfo, err = c.GetSignalInvitationInfoByClientMsgID(ctx, key)
	if err != nil {
		return nil, err
	}
	return invitationInfo, utils.Wrap1(c.DelUserSignalList(ctx, userID))
}

func (c *cache) DelUserSignalList(ctx context.Context, userID string) error {
	return utils.Wrap1(c.rdb.Del(ctx, signalListCache+userID).Err())
}

func (c *cache) DelMsgFromCache(ctx context.Context, userID string, seqList []int64) error {
	for _, seq := range seqList {
		key := messageCache + userID + "_" + strconv.Itoa(int(seq))
		result, err := c.rdb.Get(ctx, key).Result()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			return utils.Wrap1(err)
		}
		var msg sdkws.MsgData
		if err := jsonpb.UnmarshalString(result, &msg); err != nil {
			return err
		}
		msg.Status = constant.MsgDeleted
		s, err := utils.Pb2String(&msg)
		if err != nil {
			return utils.Wrap1(err)
		}
		if err := c.rdb.Set(ctx, key, s, time.Duration(config.Config.MsgCacheTimeout)*time.Second).Err(); err != nil {
			return utils.Wrap1(err)
		}
	}
	return nil
}

func (c *cache) SetGetuiToken(ctx context.Context, token string, expireTime int64) error {
	return utils.Wrap1(c.rdb.Set(ctx, getuiToken, token, time.Duration(expireTime)*time.Second).Err())
}

func (c *cache) GetGetuiToken(ctx context.Context) (string, error) {
	return utils.Wrap2(c.rdb.Get(ctx, getuiToken).Result())
}

func (c *cache) SetGetuiTaskID(ctx context.Context, taskID string, expireTime int64) error {
	return utils.Wrap1(c.rdb.Set(ctx, getuiTaskID, taskID, time.Duration(expireTime)*time.Second).Err())
}

func (c *cache) GetGetuiTaskID(ctx context.Context) (string, error) {
	return utils.Wrap2(c.rdb.Get(ctx, getuiTaskID).Result())
}

func (c *cache) SetSendMsgStatus(ctx context.Context, id string, status int32) error {
	return utils.Wrap1(c.rdb.Set(ctx, sendMsgFailedFlag+id, status, time.Hour*24).Err())
}

func (c *cache) GetSendMsgStatus(ctx context.Context, id string) (int32, error) {
	result, err := c.rdb.Get(ctx, sendMsgFailedFlag+id).Int()
	return int32(result), utils.Wrap1(err)
}

func (c *cache) SetFcmToken(ctx context.Context, account string, platformID int, fcmToken string, expireTime int64) (err error) {
	return utils.Wrap1(c.rdb.Set(ctx, FcmToken+account+":"+strconv.Itoa(platformID), fcmToken, time.Duration(expireTime)*time.Second).Err())
}

func (c *cache) GetFcmToken(ctx context.Context, account string, platformID int) (string, error) {
	return utils.Wrap2(c.rdb.Get(ctx, FcmToken+account+":"+strconv.Itoa(platformID)).Result())
}

func (c *cache) DelFcmToken(ctx context.Context, account string, platformID int) error {
	return utils.Wrap1(c.rdb.Del(ctx, FcmToken+account+":"+strconv.Itoa(platformID)).Err())
}

func (c *cache) IncrUserBadgeUnreadCountSum(ctx context.Context, userID string) (int, error) {
	seq, err := c.rdb.Incr(ctx, userBadgeUnreadCountSum+userID).Result()
	return int(seq), utils.Wrap1(err)
}

func (c *cache) SetUserBadgeUnreadCountSum(ctx context.Context, userID string, value int) error {
	return utils.Wrap1(c.rdb.Set(ctx, userBadgeUnreadCountSum+userID, value, 0).Err())
}

func (c *cache) GetUserBadgeUnreadCountSum(ctx context.Context, userID string) (int, error) {
	return utils.Wrap2(c.rdb.Get(ctx, userBadgeUnreadCountSum+userID).Int())
}

func (c *cache) LockMessageTypeKey(ctx context.Context, clientMsgID string, TypeKey string) error {
	key := exTypeKeyLocker + clientMsgID + "_" + TypeKey
	return utils.Wrap1(c.rdb.SetNX(ctx, key, 1, time.Minute).Err())
}

func (c *cache) UnLockMessageTypeKey(ctx context.Context, clientMsgID string, TypeKey string) error {
	key := exTypeKeyLocker + clientMsgID + "_" + TypeKey
	return utils.Wrap1(c.rdb.Del(ctx, key).Err())
}

func (c *cache) getMessageReactionExPrefix(clientMsgID string, sessionType int32) string {
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

func (c *cache) JudgeMessageReactionExist(ctx context.Context, clientMsgID string, sessionType int32) (bool, error) {
	n, err := c.rdb.Exists(ctx, c.getMessageReactionExPrefix(clientMsgID, sessionType)).Result()
	if err != nil {
		return false, utils.Wrap(err, "")
	}
	return n > 0, nil
}

func (c *cache) SetMessageTypeKeyValue(ctx context.Context, clientMsgID string, sessionType int32, typeKey, value string) error {
	return utils.Wrap1(c.rdb.HSet(ctx, c.getMessageReactionExPrefix(clientMsgID, sessionType), typeKey, value).Err())
}

func (c *cache) SetMessageReactionExpire(ctx context.Context, clientMsgID string, sessionType int32, expiration time.Duration) (bool, error) {
	return utils.Wrap2(c.rdb.Expire(ctx, c.getMessageReactionExPrefix(clientMsgID, sessionType), expiration).Result())
}

func (c *cache) GetMessageTypeKeyValue(ctx context.Context, clientMsgID string, sessionType int32, typeKey string) (string, error) {
	return utils.Wrap2(c.rdb.HGet(ctx, c.getMessageReactionExPrefix(clientMsgID, sessionType), typeKey).Result())
}

func (c *cache) GetOneMessageAllReactionList(ctx context.Context, clientMsgID string, sessionType int32) (map[string]string, error) {
	return utils.Wrap2(c.rdb.HGetAll(ctx, c.getMessageReactionExPrefix(clientMsgID, sessionType)).Result())
}

func (c *cache) DeleteOneMessageKey(ctx context.Context, clientMsgID string, sessionType int32, subKey string) error {
	return utils.Wrap1(c.rdb.HDel(ctx, c.getMessageReactionExPrefix(clientMsgID, sessionType), subKey).Err())
}
