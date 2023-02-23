package cache

import (
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/tracelog"
	pbChat "OpenIM/pkg/proto/msg"
	pbRtc "OpenIM/pkg/proto/rtc"
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

	SignalListCache = "SIGNAL_LIST_CACHE:"

	SignalCache = "SIGNAL_CACHE:"
)

type MsgCache interface {
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

	SetTokenMapByUidPid(ctx context.Context, userID string, platformID int, m map[string]int) error
	DeleteTokenByUidPid(ctx context.Context, userID string, platformID int, fields []string) error
	GetMessagesBySeq(ctx context.Context, userID string, seqList []int64) (seqMsg []*sdkws.MsgData, failedSeqList []int64, err error)
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

	SetSendMsgStatus(ctx context.Context, id string, status int32) error
	GetSendMsgStatus(ctx context.Context, id string) (int32, error)
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

func NewMsgCache(client redis.UniversalClient) MsgCache {
	return &msgCache{rdb: client}
}

type msgCache struct {
	rdb redis.UniversalClient
}

func (m *msgCache) IncrUserSeq(ctx context.Context, userID string) (int64, error) {
	return utils.Wrap2(m.rdb.Get(ctx, userIncrSeq+userID).Int64())
}

func (m *msgCache) GetUserMaxSeq(ctx context.Context, userID string) (int64, error) {
	return utils.Wrap2(m.rdb.Get(ctx, userIncrSeq+userID).Int64())
}

func (m *msgCache) SetUserMaxSeq(ctx context.Context, userID string, maxSeq int64) error {
	return utils.Wrap1(m.rdb.Set(ctx, userIncrSeq+userID, maxSeq, 0).Err())
}

func (m *msgCache) SetUserMinSeq(ctx context.Context, userID string, minSeq int64) (err error) {
	return utils.Wrap1(m.rdb.Set(ctx, userMinSeq+userID, minSeq, 0).Err())
}

func (m *msgCache) GetUserMinSeq(ctx context.Context, userID string) (int64, error) {
	return utils.Wrap2(m.rdb.Get(ctx, userMinSeq+userID).Int64())
}

func (m *msgCache) SetGroupUserMinSeq(ctx context.Context, groupID, userID string, minSeq int64) (err error) {
	key := groupUserMinSeq + "g:" + groupID + "u:" + userID
	return utils.Wrap1(m.rdb.Set(ctx, key, minSeq, 0).Err())
}

func (m *msgCache) GetGroupUserMinSeq(ctx context.Context, groupID, userID string) (int64, error) {
	return utils.Wrap2(m.rdb.Get(ctx, groupMinSeq+groupID).Int64())
}

func (m *msgCache) GetGroupMaxSeq(ctx context.Context, groupID string) (int64, error) {
	return utils.Wrap2(m.rdb.Get(ctx, groupMaxSeq+groupID).Int64())
}

func (m *msgCache) GetGroupMinSeq(ctx context.Context, groupID string) (int64, error) {
	return utils.Wrap2(m.rdb.Get(ctx, groupMinSeq+groupID).Int64())
}

func (m *msgCache) IncrGroupMaxSeq(ctx context.Context, groupID string) (int64, error) {
	key := groupMaxSeq + groupID
	seq, err := m.rdb.Incr(ctx, key).Uint64()
	return int64(seq), utils.Wrap1(err)
}

func (m *msgCache) SetGroupMaxSeq(ctx context.Context, groupID string, maxSeq int64) error {
	key := groupMaxSeq + groupID
	return utils.Wrap1(m.rdb.Set(ctx, key, maxSeq, 0).Err())
}

func (m *msgCache) SetGroupMinSeq(ctx context.Context, groupID string, minSeq int64) error {
	key := groupMinSeq + groupID
	return utils.Wrap1(m.rdb.Set(ctx, key, minSeq, 0).Err())
}

func (m *msgCache) AddTokenFlag(ctx context.Context, userID string, platformID int, token string, flag int) error {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
	return utils.Wrap1(m.rdb.HSet(ctx, key, token, flag).Err())
}

func (m *msgCache) GetTokensWithoutError(ctx context.Context, userID, platformID string) (map[string]int, error) {
	key := uidPidToken + userID + ":" + platformID
	m, err := m.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, utils.Wrap1(err)
	}
	mm := make(map[string]int)
	for k, v := range m {
		mm[k] = utils.StringToInt(v)
	}
	return mm, nil
}

func (m *msgCache) SetTokenMapByUidPid(ctx context.Context, userID string, platformID int, m map[string]int) error {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
	mm := make(map[string]interface{})
	for k, v := range m {
		mm[k] = v
	}
	return utils.Wrap1(m.rdb.HSet(ctx, key, mm).Err())
}

func (m *msgCache) DeleteTokenByUidPid(ctx context.Context, userID string, platformID int, fields []string) error {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
	return utils.Wrap1(m.rdb.HDel(ctx, key, fields...).Err())
}

func (m *msgCache) GetMessagesBySeq(ctx context.Context, userID string, seqList []int64) (seqMsg []*sdkws.MsgData, failedSeqList []int64, err error) {
	var errResult error
	for _, v := range seqList {
		//MESSAGE_CACHE:169.254.225.224_reliability1653387820_0_1
		key := messageCache + userID + "_" + strconv.Itoa(int(v))
		result, err := m.rdb.Get(ctx, key).Result()
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

func (m *msgCache) SetMessageToCache(ctx context.Context, userID string, msgList []*pbChat.MsgDataToMQ) (int, error) {
	pipe := m.rdb.Pipeline()
	var failedList []pbChat.MsgDataToMQ
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
	if len(failedList) != 0 {
		return len(failedList), errors.New(fmt.Sprintf("set msg to cache failed, failed lists: %q,%s", failedList, tracelog.GetOperationID(ctx)))
	}
	_, err := pipe.Exec(ctx)
	return 0, err
}

func (m *msgCache) DeleteMessageFromCache(ctx context.Context, userID string, msgList []*pbChat.MsgDataToMQ) error {
	for _, msg := range msgList {
		if err := m.rdb.Del(ctx, messageCache+userID+"_"+strconv.Itoa(int(msg.MsgData.Seq))).Err(); err != nil {
			return utils.Wrap1(err)
		}
	}
	return nil
}

func (m *msgCache) CleanUpOneUserAllMsg(ctx context.Context, userID string) error {
	key := messageCache + userID + "_" + "*"
	vals, err := m.rdb.Keys(ctx, key).Result()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return utils.Wrap1(err)
	}
	for _, v := range vals {
		if err := m.rdb.Del(ctx, v).Err(); err != nil {
			return utils.Wrap1(err)
		}
	}
	return nil
}

func (m *msgCache) HandleSignalInfo(ctx context.Context, msg *sdkws.MsgData, pushToUserID string) (isSend bool, err error) {
	req := &pbRtc.SignalReq{}
	if err := proto.Unmarshal(msg.Content, req); err != nil {
		return false, utils.Wrap1(err)
	}
	var inviteeUserIDList []string
	var isInviteSignal bool
	switch signalInfo := req.Payload.(type) {
	case *pbRtc.SignalReq_Invite:
		inviteeUserIDList = signalInfo.Invite.Invitation.InviteeUserIDList
		isInviteSignal = true
	case *pbRtc.SignalReq_InviteInGroup:
		inviteeUserIDList = signalInfo.InviteInGroup.Invitation.InviteeUserIDList
		isInviteSignal = true
		if !utils.Contain(pushToUserID, inviteeUserIDList...) {
			return false, nil
		}
	case *pbRtc.SignalReq_HungUp, *pbRtc.SignalReq_Cancel, *pbRtc.SignalReq_Reject, *pbRtc.SignalReq_Accept:
		return false, utils.Wrap1(errors.New("signalInfo do not need offlinePush"))
	default:
		return false, nil
	}
	if isInviteSignal {
		for _, userID := range inviteeUserIDList {
			timeout, err := strconv.Atoi(config.Config.Rtc.SignalTimeout)
			if err != nil {
				return false, utils.Wrap1(err)
			}
			keyList := SignalListCache + userID
			err = m.rdb.LPush(ctx, keyList, msg.ClientMsgID).Err()
			if err != nil {
				return false, utils.Wrap1(err)
			}
			err = m.rdb.Expire(ctx, keyList, time.Duration(timeout)*time.Second).Err()
			if err != nil {
				return false, utils.Wrap1(err)
			}
			key := SignalCache + msg.ClientMsgID
			err = m.rdb.Set(ctx, key, msg.Content, time.Duration(timeout)*time.Second).Err()
			if err != nil {
				return false, utils.Wrap1(err)
			}
		}
	}
	return true, nil
}

func (m *msgCache) GetSignalInfoFromCacheByClientMsgID(ctx context.Context, clientMsgID string) (invitationInfo *pbRtc.SignalInviteReq, err error) {
	bytes, err := m.rdb.Get(ctx, SignalCache+clientMsgID).Bytes()
	if err != nil {
		return nil, utils.Wrap1(err)
	}
	req := &pbRtc.SignalReq{}
	if err = proto.Unmarshal(bytes, req); err != nil {
		return nil, utils.Wrap1(err)
	}
	invitationInfo = &pbRtc.SignalInviteReq{}
	switch req2 := req.Payload.(type) {
	case *pbRtc.SignalReq_Invite:
		invitationInfo.Invitation = req2.Invite.Invitation
		invitationInfo.OpUserID = req2.Invite.OpUserID
	case *pbRtc.SignalReq_InviteInGroup:
		invitationInfo.Invitation = req2.InviteInGroup.Invitation
		invitationInfo.OpUserID = req2.InviteInGroup.OpUserID
	}
	return invitationInfo, nil
}

func (m *msgCache) GetAvailableSignalInvitationInfo(ctx context.Context, userID string) (invitationInfo *pbRtc.SignalInviteReq, err error) {
	key, err := m.rdb.LPop(ctx, SignalListCache+userID).Result()
	if err != nil {
		return nil, utils.Wrap1(err)
	}
	invitationInfo, err = m.GetSignalInfoFromCacheByClientMsgID(ctx, key)
	if err != nil {
		return nil, err
	}
	return invitationInfo, m.DelUserSignalList(ctx, userID)
}

func (m *msgCache) DelUserSignalList(ctx context.Context, userID string) error {
	return utils.Wrap1(m.rdb.Del(ctx, SignalListCache+userID).Err())
}

func (m *msgCache) DelMsgFromCache(ctx context.Context, userID string, seqList []int64) error {
	for _, seq := range seqList {
		key := messageCache + userID + "_" + strconv.Itoa(int(seq))
		result, err := m.rdb.Get(ctx, key).Result()
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
		if err := m.rdb.Set(ctx, key, s, time.Duration(config.Config.MsgCacheTimeout)*time.Second).Err(); err != nil {
			return utils.Wrap1(err)
		}
	}
	return nil
}

func (m *msgCache) SetGetuiToken(ctx context.Context, token string, expireTime int64) error {
	return utils.Wrap1(m.rdb.Set(ctx, getuiToken, token, time.Duration(expireTime)*time.Second).Err())
}

func (m *msgCache) GetGetuiToken(ctx context.Context) (string, error) {
	return utils.Wrap2(m.rdb.Get(ctx, getuiToken).Result())
}

func (m *msgCache) SetGetuiTaskID(ctx context.Context, taskID string, expireTime int64) error {
	return utils.Wrap1(m.rdb.Set(ctx, getuiTaskID, taskID, time.Duration(expireTime)*time.Second).Err())
}

func (m *msgCache) GetGetuiTaskID(ctx context.Context) (string, error) {
	return utils.Wrap2(m.rdb.Get(ctx, getuiTaskID).Result())
}

func (m *msgCache) SetSendMsgStatus(ctx context.Context, id string, status int32) error {
	return utils.Wrap1(m.rdb.Set(ctx, sendMsgFailedFlag+id, status, time.Hour*24).Err())
}

func (m *msgCache) GetSendMsgStatus(ctx context.Context, id string) (int32, error) {
	result, err := m.rdb.Get(ctx, sendMsgFailedFlag+id).Int()
	return int32(result), utils.Wrap1(err)
}

func (m *msgCache) SetFcmToken(ctx context.Context, account string, platformID int, fcmToken string, expireTime int64) (err error) {
	return utils.Wrap1(m.rdb.Set(ctx, FcmToken+account+":"+strconv.Itoa(platformID), fcmToken, time.Duration(expireTime)*time.Second).Err())
}

func (m *msgCache) GetFcmToken(ctx context.Context, account string, platformID int) (string, error) {
	return utils.Wrap2(m.rdb.Get(ctx, FcmToken+account+":"+strconv.Itoa(platformID)).Result())
}

func (m *msgCache) DelFcmToken(ctx context.Context, account string, platformID int) error {
	return utils.Wrap1(m.rdb.Del(ctx, FcmToken+account+":"+strconv.Itoa(platformID)).Err())
}

func (m *msgCache) IncrUserBadgeUnreadCountSum(ctx context.Context, userID string) (int, error) {
	seq, err := m.rdb.Incr(ctx, userBadgeUnreadCountSum+userID).Result()
	return int(seq), utils.Wrap1(err)
}

func (m *msgCache) SetUserBadgeUnreadCountSum(ctx context.Context, userID string, value int) error {
	return utils.Wrap1(m.rdb.Set(ctx, userBadgeUnreadCountSum+userID, value, 0).Err())
}

func (m *msgCache) GetUserBadgeUnreadCountSum(ctx context.Context, userID string) (int, error) {
	return utils.Wrap2(m.rdb.Get(ctx, userBadgeUnreadCountSum+userID).Int())
}

func (m *msgCache) LockMessageTypeKey(ctx context.Context, clientMsgID string, TypeKey string) error {
	key := exTypeKeyLocker + clientMsgID + "_" + TypeKey
	return utils.Wrap1(m.rdb.SetNX(ctx, key, 1, time.Minute).Err())
}

func (m *msgCache) UnLockMessageTypeKey(ctx context.Context, clientMsgID string, TypeKey string) error {
	key := exTypeKeyLocker + clientMsgID + "_" + TypeKey
	return utils.Wrap1(m.rdb.Del(ctx, key).Err())
}

func (m *msgCache) getMessageReactionExPrefix(clientMsgID string, sessionType int32) string {
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

func (m *msgCache) JudgeMessageReactionEXISTS(ctx context.Context, clientMsgID string, sessionType int32) (bool, error) {
	n, err := m.rdb.Exists(ctx, m.getMessageReactionExPrefix(clientMsgID, sessionType)).Result()
	if err != nil {
		return false, utils.Wrap(err, "")
	}
	return n > 0, nil
}

func (m *msgCache) SetMessageTypeKeyValue(ctx context.Context, clientMsgID string, sessionType int32, typeKey, value string) error {
	return utils.Wrap1(m.rdb.HSet(ctx, m.getMessageReactionExPrefix(clientMsgID, sessionType), typeKey, value).Err())
}

func (m *msgCache) SetMessageReactionExpire(ctx context.Context, clientMsgID string, sessionType int32, expiration time.Duration) (bool, error) {
	return utils.Wrap2(m.rdb.Expire(ctx, m.getMessageReactionExPrefix(clientMsgID, sessionType), expiration).Result())
}

func (m *msgCache) GetMessageTypeKeyValue(ctx context.Context, clientMsgID string, sessionType int32, typeKey string) (string, error) {
	return utils.Wrap2(m.rdb.HGet(ctx, m.getMessageReactionExPrefix(clientMsgID, sessionType), typeKey).Result())
}

func (m *msgCache) GetOneMessageAllReactionList(ctx context.Context, clientMsgID string, sessionType int32) (map[string]string, error) {
	return utils.Wrap2(m.rdb.HGetAll(ctx, m.getMessageReactionExPrefix(clientMsgID, sessionType)).Result())
}

func (m *msgCache) DeleteOneMessageKey(ctx context.Context, clientMsgID string, sessionType int32, subKey string) error {
	return utils.Wrap1(m.rdb.HDel(ctx, m.getMessageReactionExPrefix(clientMsgID, sessionType), subKey).Err())
}
