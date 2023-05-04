package cache

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/gogo/protobuf/jsonpb"

	"google.golang.org/protobuf/proto"

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
	fcmToken                = "FCM_TOKEN:"
	groupUserMinSeq         = "GROUP_USER_MIN_SEQ:"
	groupMaxSeq             = "GROUP_MAX_SEQ:"
	groupMinSeq             = "GROUP_MIN_SEQ:"
	sendMsgFailedFlag       = "SEND_MSG_FAILED_FLAG:"
	userBadgeUnreadCountSum = "USER_BADGE_UNREAD_COUNT_SUM:"
	exTypeKeyLocker         = "EX_LOCK:"
	uidPidToken             = "UID_PID_TOKEN_STATUS:"
	userNotificationSeq     = "USER_NOTIFICATION_SEQ:"
	userMinNotificationSeq  = "USER_MIN_NOTIFICATION_SEQ:"
	groupNotificationSeq    = "GROUP_NOTIFICATION_SEQ:"
	groupMinNotificationSeq = "GROUP_MIN_NOTIFICATION_SEQ:"
)

type MsgModel interface {
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
	SetMessageToCache(ctx context.Context, userID string, msgList []*sdkws.MsgData) (int, error)
	DeleteMessageFromCache(ctx context.Context, userID string, msgList []*sdkws.MsgData) error
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
	// notificatio

	// 	IncrUserNotificationSeq(ctx context.Context, userID string) (int64, error)
	// 	GetUserNotificationMaxSeq(ctx context.Context, userID string) (int64, error)
	// 	SetUserNotificationMaxSeq(ctx context.Context, userID string, maxSeq int64) error
	// 	SetUserNotificationMinSeq(ctx context.Context, userID string, minSeq int64) (err error)
	// 	GetUserNotificationMinSeq(ctx context.Context, userID string) (int64, error)
	// 	SetGroupNotificationUserMinSeq(ctx context.Context, groupID, userID string, minSeq int64) (err error)
	// 	GetGroupNotificationUserMinSeq(ctx context.Context, groupID, userID string) (int64, error)
	// 	GetGroupNotificationMaxSeq(ctx context.Context, groupID string) (int64, error)
	// 	GetGroupNotificationMinSeq(ctx context.Context, groupID string) (int64, error)
	// 	IncrGroupNotificationMaxSeq(ctx context.Context, groupID string) (int64, error)
	// 	SetGroupNotificationMaxSeq(ctx context.Context, groupID string, maxSeq int64) error
	// 	SetGroupNotificationMinSeq(ctx context.Context, groupID string, minSeq int64) error
}

func NewMsgCacheModel(client redis.UniversalClient) MsgModel {
	return &msgCache{rdb: client}
}

type msgCache struct {
	rdb redis.UniversalClient
}

// 兼容老版本调用
func (c *msgCache) DelKeys() {
	for _, key := range []string{"GROUP_CACHE:", "FRIEND_RELATION_CACHE:", "BLACK_LIST_CACHE:", "USER_INFO_CACHE:", "GROUP_INFO_CACHE:", "JOINED_GROUP_LIST_CACHE:",
		"GROUP_MEMBER_INFO_CACHE:", "GROUP_ALL_MEMBER_INFO_CACHE:", "ALL_FRIEND_INFO_CACHE:"} {
		fName := utils.GetSelfFuncName()
		var cursor uint64
		var n int
		for {
			var keys []string
			var err error
			keys, cursor, err = c.rdb.Scan(context.Background(), cursor, key+"*", scanCount).Result()
			if err != nil {
				panic(err.Error())
			}
			n += len(keys)
			// for each for redis cluster
			for _, key := range keys {
				if err = c.rdb.Del(context.Background(), key).Err(); err != nil {
					log.NewError("", fName, key, err.Error())
					err = c.rdb.Del(context.Background(), key).Err()
					if err != nil {
						panic(err.Error())
					}
				}
			}
			if cursor == 0 {
				break
			}
		}
	}
}

func (c *msgCache) IncrUserSeq(ctx context.Context, userID string) (int64, error) {
	return utils.Wrap2(c.rdb.Get(ctx, userIncrSeq+userID).Int64())
}

func (c *msgCache) GetUserMaxSeq(ctx context.Context, userID string) (int64, error) {
	return utils.Wrap2(c.rdb.Get(ctx, userIncrSeq+userID).Int64())
}

func (c *msgCache) SetUserMaxSeq(ctx context.Context, userID string, maxSeq int64) error {
	return errs.Wrap(c.rdb.Set(ctx, userIncrSeq+userID, maxSeq, 0).Err())
}

func (c *msgCache) SetUserMinSeq(ctx context.Context, userID string, minSeq int64) (err error) {
	return errs.Wrap(c.rdb.Set(ctx, userMinSeq+userID, minSeq, 0).Err())
}

func (c *msgCache) GetUserMinSeq(ctx context.Context, userID string) (int64, error) {
	return utils.Wrap2(c.rdb.Get(ctx, userMinSeq+userID).Int64())
}

func (c *msgCache) SetGroupUserMinSeq(ctx context.Context, groupID, userID string, minSeq int64) (err error) {
	key := groupUserMinSeq + "g:" + groupID + "u:" + userID
	return errs.Wrap(c.rdb.Set(ctx, key, minSeq, 0).Err())
}

func (c *msgCache) GetGroupUserMinSeq(ctx context.Context, groupID, userID string) (int64, error) {
	key := groupUserMinSeq + "g:" + groupID + "u:" + userID
	return utils.Wrap2(c.rdb.Get(ctx, key).Int64())
}

func (c *msgCache) GetGroupMaxSeq(ctx context.Context, groupID string) (int64, error) {
	return utils.Wrap2(c.rdb.Get(ctx, groupMaxSeq+groupID).Int64())
}

func (c *msgCache) GetGroupMinSeq(ctx context.Context, groupID string) (int64, error) {
	return utils.Wrap2(c.rdb.Get(ctx, groupMinSeq+groupID).Int64())
}

func (c *msgCache) IncrGroupMaxSeq(ctx context.Context, groupID string) (int64, error) {
	key := groupMaxSeq + groupID
	seq, err := c.rdb.Incr(ctx, key).Uint64()
	return int64(seq), errs.Wrap(err)
}

func (c *msgCache) SetGroupMaxSeq(ctx context.Context, groupID string, maxSeq int64) error {
	key := groupMaxSeq + groupID
	return errs.Wrap(c.rdb.Set(ctx, key, maxSeq, 0).Err())
}

func (c *msgCache) SetGroupMinSeq(ctx context.Context, groupID string, minSeq int64) error {
	key := groupMinSeq + groupID
	return errs.Wrap(c.rdb.Set(ctx, key, minSeq, 0).Err())
}

func (c *msgCache) AddTokenFlag(ctx context.Context, userID string, platformID int, token string, flag int) error {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
	return errs.Wrap(c.rdb.HSet(ctx, key, token, flag).Err())
}

func (c *msgCache) GetTokensWithoutError(ctx context.Context, userID, platformID string) (map[string]int, error) {
	key := uidPidToken + userID + ":" + platformID
	m, err := c.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	mm := make(map[string]int)
	for k, v := range m {
		mm[k] = utils.StringToInt(v)
	}
	return mm, nil
}

func (c *msgCache) SetTokenMapByUidPid(ctx context.Context, userID string, platform string, m map[string]int) error {
	key := uidPidToken + userID + ":" + platform
	mm := make(map[string]interface{})
	for k, v := range m {
		mm[k] = v
	}
	return errs.Wrap(c.rdb.HSet(ctx, key, mm).Err())
}

func (c *msgCache) DeleteTokenByUidPid(ctx context.Context, userID string, platform string, fields []string) error {
	key := uidPidToken + userID + ":" + platform
	return errs.Wrap(c.rdb.HDel(ctx, key, fields...).Err())
}

func (c *msgCache) getMessageCacheKey(sourceID string, seq int64) string {
	return messageCache + sourceID + "_" + strconv.Itoa(int(seq))
}

func (c *msgCache) allMessageCacheKey(sourceID string) string {
	return messageCache + sourceID + "_*"
}

func (c *msgCache) GetMessagesBySeq(ctx context.Context, userID string, seqs []int64) (seqMsgs []*sdkws.MsgData, failedSeqs []int64, err error) {
	pipe := c.rdb.Pipeline()
	for _, v := range seqs {
		//MESSAGE_CACHE:169.254.225.224_reliability1653387820_0_1
		key := c.getMessageCacheKey(userID, v)
		if err := pipe.Get(ctx, key).Err(); err != nil && err != redis.Nil {
			return nil, nil, err
		}
	}
	result, err := pipe.Exec(ctx)
	for i, v := range result {
		if v.Err() != nil {
			failedSeqs = append(failedSeqs, seqs[i])
		} else {
			msg := sdkws.MsgData{}
			err = jsonpb.UnmarshalString(v.String(), &msg)
			if err != nil {
				failedSeqs = append(failedSeqs, seqs[i])
			} else {
				seqMsgs = append(seqMsgs, &msg)
			}
		}
	}
	return seqMsgs, failedSeqs, err
}

func (c *msgCache) SetMessageToCache(ctx context.Context, userID string, msgList []*sdkws.MsgData) (int, error) {
	pipe := c.rdb.Pipeline()
	var failedMsgs []sdkws.MsgData
	for _, msg := range msgList {
		key := c.getMessageCacheKey(userID, msg.Seq)
		s, err := utils.Pb2String(msg)
		if err != nil {
			return 0, errs.Wrap(err)
		}
		err = pipe.Set(ctx, key, s, time.Duration(config.Config.MsgCacheTimeout)*time.Second).Err()
		if err != nil {
			return 0, errs.Wrap(err)
		}
	}
	if len(failedMsgs) != 0 {
		return len(failedMsgs), fmt.Errorf("set msg to msgCache failed, failed lists: %v, %s", failedMsgs, userID)
	}
	_, err := pipe.Exec(ctx)
	return 0, err
}

func (c *msgCache) DeleteMessageFromCache(ctx context.Context, userID string, msgList []*sdkws.MsgData) error {
	pipe := c.rdb.Pipeline()
	for _, v := range msgList {
		if err := pipe.Del(ctx, c.getMessageCacheKey(userID, v.Seq)).Err(); err != nil {
			return errs.Wrap(err)
		}
	}
	_, err := pipe.Exec(ctx)
	return errs.Wrap(err)
}

func (c *msgCache) CleanUpOneUserAllMsg(ctx context.Context, userID string) error {
	vals, err := c.rdb.Keys(ctx, c.allMessageCacheKey(userID)).Result()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return errs.Wrap(err)
	}
	pipe := c.rdb.Pipeline()
	for _, v := range vals {
		if err := pipe.Del(ctx, v).Err(); err != nil {
			return errs.Wrap(err)
		}
	}
	_, err = pipe.Exec(ctx)
	return errs.Wrap(err)
}

func (c *msgCache) HandleSignalInvite(ctx context.Context, msg *sdkws.MsgData, pushToUserID string) (isSend bool, err error) {
	req := &sdkws.SignalReq{}
	if err := proto.Unmarshal(msg.Content, req); err != nil {
		return false, errs.Wrap(err)
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
		return false, errs.Wrap(errors.New("signalInfo do not need offlinePush"))
	default:
		return false, nil
	}
	if isInviteSignal {
		pipe := c.rdb.Pipeline()
		for _, userID := range inviteeUserIDs {
			timeout, err := strconv.Atoi(config.Config.Rtc.SignalTimeout)
			if err != nil {
				return false, errs.Wrap(err)
			}
			keys := signalListCache + userID
			err = pipe.LPush(ctx, keys, msg.ClientMsgID).Err()
			if err != nil {
				return false, errs.Wrap(err)
			}
			err = pipe.Expire(ctx, keys, time.Duration(timeout)*time.Second).Err()
			if err != nil {
				return false, errs.Wrap(err)
			}
			key := signalCache + msg.ClientMsgID
			err = pipe.Set(ctx, key, msg.Content, time.Duration(timeout)*time.Second).Err()
			if err != nil {
				return false, errs.Wrap(err)
			}
		}
		_, err := pipe.Exec(ctx)
		if err != nil {
			return false, errs.Wrap(err)
		}
	}
	return true, nil
}

func (c *msgCache) GetSignalInvitationInfoByClientMsgID(ctx context.Context, clientMsgID string) (signalInviteReq *sdkws.SignalInviteReq, err error) {
	bytes, err := c.rdb.Get(ctx, signalCache+clientMsgID).Bytes()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	signalReq := &sdkws.SignalReq{}
	if err = proto.Unmarshal(bytes, signalReq); err != nil {
		return nil, errs.Wrap(err)
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

func (c *msgCache) GetAvailableSignalInvitationInfo(ctx context.Context, userID string) (invitationInfo *sdkws.SignalInviteReq, err error) {
	key, err := c.rdb.LPop(ctx, signalListCache+userID).Result()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	invitationInfo, err = c.GetSignalInvitationInfoByClientMsgID(ctx, key)
	if err != nil {
		return nil, err
	}
	return invitationInfo, errs.Wrap(c.DelUserSignalList(ctx, userID))
}

func (c *msgCache) DelUserSignalList(ctx context.Context, userID string) error {
	return errs.Wrap(c.rdb.Del(ctx, signalListCache+userID).Err())
}

func (c *msgCache) DelMsgFromCache(ctx context.Context, userID string, seqs []int64) error {
	for _, seq := range seqs {
		key := c.getMessageCacheKey(userID, seq)
		result, err := c.rdb.Get(ctx, key).Result()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			return errs.Wrap(err)
		}
		var msg sdkws.MsgData
		if err := jsonpb.UnmarshalString(result, &msg); err != nil {
			return err
		}
		msg.Status = constant.MsgDeleted
		s, err := utils.Pb2String(&msg)
		if err != nil {
			return errs.Wrap(err)
		}
		if err := c.rdb.Set(ctx, key, s, time.Duration(config.Config.MsgCacheTimeout)*time.Second).Err(); err != nil {
			return errs.Wrap(err)
		}
	}
	return nil
}

func (c *msgCache) SetGetuiToken(ctx context.Context, token string, expireTime int64) error {
	return errs.Wrap(c.rdb.Set(ctx, getuiToken, token, time.Duration(expireTime)*time.Second).Err())
}

func (c *msgCache) GetGetuiToken(ctx context.Context) (string, error) {
	return utils.Wrap2(c.rdb.Get(ctx, getuiToken).Result())
}

func (c *msgCache) SetGetuiTaskID(ctx context.Context, taskID string, expireTime int64) error {
	return errs.Wrap(c.rdb.Set(ctx, getuiTaskID, taskID, time.Duration(expireTime)*time.Second).Err())
}

func (c *msgCache) GetGetuiTaskID(ctx context.Context) (string, error) {
	return utils.Wrap2(c.rdb.Get(ctx, getuiTaskID).Result())
}

func (c *msgCache) SetSendMsgStatus(ctx context.Context, id string, status int32) error {
	return errs.Wrap(c.rdb.Set(ctx, sendMsgFailedFlag+id, status, time.Hour*24).Err())
}

func (c *msgCache) GetSendMsgStatus(ctx context.Context, id string) (int32, error) {
	result, err := c.rdb.Get(ctx, sendMsgFailedFlag+id).Int()
	return int32(result), errs.Wrap(err)
}

func (c *msgCache) SetFcmToken(ctx context.Context, account string, platformID int, fcmToken string, expireTime int64) (err error) {
	return errs.Wrap(c.rdb.Set(ctx, fcmToken+account+":"+strconv.Itoa(platformID), fcmToken, time.Duration(expireTime)*time.Second).Err())
}

func (c *msgCache) GetFcmToken(ctx context.Context, account string, platformID int) (string, error) {
	return utils.Wrap2(c.rdb.Get(ctx, fcmToken+account+":"+strconv.Itoa(platformID)).Result())
}

func (c *msgCache) DelFcmToken(ctx context.Context, account string, platformID int) error {
	return errs.Wrap(c.rdb.Del(ctx, fcmToken+account+":"+strconv.Itoa(platformID)).Err())
}

func (c *msgCache) IncrUserBadgeUnreadCountSum(ctx context.Context, userID string) (int, error) {
	seq, err := c.rdb.Incr(ctx, userBadgeUnreadCountSum+userID).Result()
	return int(seq), errs.Wrap(err)
}

func (c *msgCache) SetUserBadgeUnreadCountSum(ctx context.Context, userID string, value int) error {
	return errs.Wrap(c.rdb.Set(ctx, userBadgeUnreadCountSum+userID, value, 0).Err())
}

func (c *msgCache) GetUserBadgeUnreadCountSum(ctx context.Context, userID string) (int, error) {
	return utils.Wrap2(c.rdb.Get(ctx, userBadgeUnreadCountSum+userID).Int())
}

func (c *msgCache) LockMessageTypeKey(ctx context.Context, clientMsgID string, TypeKey string) error {
	key := exTypeKeyLocker + clientMsgID + "_" + TypeKey
	return errs.Wrap(c.rdb.SetNX(ctx, key, 1, time.Minute).Err())
}

func (c *msgCache) UnLockMessageTypeKey(ctx context.Context, clientMsgID string, TypeKey string) error {
	key := exTypeKeyLocker + clientMsgID + "_" + TypeKey
	return errs.Wrap(c.rdb.Del(ctx, key).Err())
}

func (c *msgCache) getMessageReactionExPrefix(clientMsgID string, sessionType int32) string {
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

func (c *msgCache) JudgeMessageReactionExist(ctx context.Context, clientMsgID string, sessionType int32) (bool, error) {
	n, err := c.rdb.Exists(ctx, c.getMessageReactionExPrefix(clientMsgID, sessionType)).Result()
	if err != nil {
		return false, utils.Wrap(err, "")
	}
	return n > 0, nil
}

func (c *msgCache) SetMessageTypeKeyValue(ctx context.Context, clientMsgID string, sessionType int32, typeKey, value string) error {
	return errs.Wrap(c.rdb.HSet(ctx, c.getMessageReactionExPrefix(clientMsgID, sessionType), typeKey, value).Err())
}

func (c *msgCache) SetMessageReactionExpire(ctx context.Context, clientMsgID string, sessionType int32, expiration time.Duration) (bool, error) {
	return utils.Wrap2(c.rdb.Expire(ctx, c.getMessageReactionExPrefix(clientMsgID, sessionType), expiration).Result())
}

func (c *msgCache) GetMessageTypeKeyValue(ctx context.Context, clientMsgID string, sessionType int32, typeKey string) (string, error) {
	return utils.Wrap2(c.rdb.HGet(ctx, c.getMessageReactionExPrefix(clientMsgID, sessionType), typeKey).Result())
}

func (c *msgCache) GetOneMessageAllReactionList(ctx context.Context, clientMsgID string, sessionType int32) (map[string]string, error) {
	return utils.Wrap2(c.rdb.HGetAll(ctx, c.getMessageReactionExPrefix(clientMsgID, sessionType)).Result())
}

func (c *msgCache) DeleteOneMessageKey(ctx context.Context, clientMsgID string, sessionType int32, subKey string) error {
	return errs.Wrap(c.rdb.HDel(ctx, c.getMessageReactionExPrefix(clientMsgID, sessionType), subKey).Err())
}
