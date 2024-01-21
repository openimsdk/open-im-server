// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package push

import (
	"context"
	"encoding/json"
	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/OpenIMSDK/tools/mcontext"
	"github.com/openimsdk/open-im-server/v3/internal/push/offlinepush"
	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"
	"github.com/openimsdk/open-im-server/v3/pkg/msgprocessor"
	"github.com/openimsdk/open-im-server/v3/pkg/rpccache"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	"github.com/redis/go-redis/v9"

	"github.com/IBM/sarama"
	"google.golang.org/protobuf/proto"

	"github.com/OpenIMSDK/protocol/constant"
	pbchat "github.com/OpenIMSDK/protocol/msg"
	pbpush "github.com/OpenIMSDK/protocol/push"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/utils"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	kfk "github.com/openimsdk/open-im-server/v3/pkg/common/kafka"
)

type ConsumerHandler struct {
	pushConsumerGroup      *kfk.MConsumerGroup
	offlinePusher          offlinepush.OfflinePusher
	onlinePusher           OnlinePusher
	groupLocalCache        *rpccache.GroupLocalCache
	conversationLocalCache *rpccache.ConversationLocalCache
	msgRpcClient           rpcclient.MessageRpcClient
	conversationRpcClient  rpcclient.ConversationRpcClient
	groupRpcClient         rpcclient.GroupRpcClient
}

func NewConsumerHandler(offlinePusher offlinepush.OfflinePusher,
	rdb redis.UniversalClient, disCov discoveryregistry.SvcDiscoveryRegistry) *ConsumerHandler {
	var consumerHandler ConsumerHandler
	consumerHandler.pushConsumerGroup = kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{
		KafkaVersion:   sarama.V2_0_0_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false,
	}, []string{config.Config.Kafka.MsgToPush.Topic}, config.Config.Kafka.Addr,
		config.Config.Kafka.ConsumerGroupID.MsgToPush)
	consumerHandler.offlinePusher = offlinePusher
	consumerHandler.onlinePusher = NewOnlinePusher(disCov)
	consumerHandler.groupRpcClient = rpcclient.NewGroupRpcClient(disCov)
	consumerHandler.groupLocalCache = rpccache.NewGroupLocalCache(consumerHandler.groupRpcClient, rdb)
	consumerHandler.msgRpcClient = rpcclient.NewMessageRpcClient(disCov)
	consumerHandler.conversationRpcClient = rpcclient.NewConversationRpcClient(disCov)
	consumerHandler.conversationLocalCache = rpccache.NewConversationLocalCache(consumerHandler.conversationRpcClient, rdb)
	return &consumerHandler
}

func (c *ConsumerHandler) handleMs2PsChat(ctx context.Context, msg []byte) {
	msgFromMQ := pbchat.PushMsgDataToMQ{}
	if err := proto.Unmarshal(msg, &msgFromMQ); err != nil {
		log.ZError(ctx, "push Unmarshal msg err", err, "msg", string(msg))
		return
	}
	pbData := &pbpush.PushMsgReq{
		MsgData:        msgFromMQ.MsgData,
		ConversationID: msgFromMQ.ConversationID,
	}
	sec := msgFromMQ.MsgData.SendTime / 1000
	nowSec := utils.GetCurrentTimestampBySecond()
	log.ZDebug(ctx, "push msg", "msg", pbData.String(), "sec", sec, "nowSec", nowSec)
	if nowSec-sec > 10 {
		return
	}
	var err error
	switch msgFromMQ.MsgData.SessionType {
	case constant.SuperGroupChatType:
		err = c.Push2SuperGroup(ctx, pbData.MsgData.GroupID, pbData.MsgData)
	default:
		var pushUserIDList []string
		isSenderSync := utils.GetSwitchFromOptions(pbData.MsgData.Options, constant.IsSenderSync)
		if !isSenderSync || pbData.MsgData.SendID == pbData.MsgData.RecvID {
			pushUserIDList = append(pushUserIDList, pbData.MsgData.RecvID)
		} else {
			pushUserIDList = append(pushUserIDList, pbData.MsgData.RecvID, pbData.MsgData.SendID)
		}
		err = c.Push2User(ctx, pushUserIDList, pbData.MsgData)
	}
	if err != nil {
		log.ZError(ctx, "push failed", err, "msg", pbData.String())
	}
}
func (*ConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (*ConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (c *ConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for msg := range claim.Messages() {
		ctx := c.pushConsumerGroup.GetContextFromMsg(msg)
		c.handleMs2PsChat(ctx, msg.Value)
		sess.MarkMessage(msg, "")
	}
	return nil
}

// Push2User Suitable for two types of conversations, one is SingleChatType and the other is NotificationChatType.
func (c *ConsumerHandler) Push2User(ctx context.Context, userIDs []string, msg *sdkws.MsgData) error {
	log.ZDebug(ctx, "Get msg from msg_transfer And push msg", "userIDs", userIDs, "msg", msg.String())
	if err := callbackOnlinePush(ctx, userIDs, msg); err != nil {
		return err
	}

	wsResults, err := c.onlinePusher.GetConnsAndOnlinePush(ctx, msg, userIDs)
	if err != nil {
		return err
	}

	log.ZDebug(ctx, "single and notification push result", "result", wsResults, "msg", msg, "push_to_userID", userIDs)

	if !c.shouldPushOffline(ctx, msg) {
		return nil
	}

	for _, v := range wsResults {
		//message sender do not need offline push
		if msg.SendID == v.UserID {
			continue
		}
		//receiver online push success
		if v.OnlinePush {
			return nil
		}
	}
	offlinePUshUserID := []string{msg.RecvID}
	//receiver offline push
	if err = callbackOfflinePush(ctx, offlinePUshUserID, msg, nil); err != nil {
		return err
	}

	err = c.offlinePushMsg(ctx, msg, offlinePUshUserID)
	if err != nil {
		return err
	}

	return nil
}

func (c *ConsumerHandler) Push2SuperGroup(ctx context.Context, groupID string, msg *sdkws.MsgData) (err error) {
	log.ZDebug(ctx, "Get super group msg from msg_transfer and push msg", "msg", msg.String(), "groupID", groupID)
	var pushToUserIDs []string
	if err = callbackBeforeSuperGroupOnlinePush(ctx, groupID, msg, &pushToUserIDs); err != nil {
		return err
	}

	err = c.groupMessagesHandler(ctx, groupID, &pushToUserIDs, msg)
	if err != nil {
		return err
	}

	wsResults, err := c.onlinePusher.GetConnsAndOnlinePush(ctx, msg, pushToUserIDs)
	if err != nil {
		return err
	}

	log.ZDebug(ctx, "group push result", "result", wsResults, "msg", msg)

	if !c.shouldPushOffline(ctx, msg) {
		return nil
	}
	needOfflinePushUserIDs := c.onlinePusher.GetOnlinePushFailedUserIDs(ctx, msg, wsResults, &pushToUserIDs)

	//filter some user, like don not disturb or don't need offline push etc.
	needOfflinePushUserIDs, err = c.filterGroupMessageOfflinePush(ctx, groupID, msg, needOfflinePushUserIDs)
	if err != nil {
		return err
	}
	// Use offline push messaging
	if len(needOfflinePushUserIDs) > 0 {
		var offlinePushUserIDs []string
		err = callbackOfflinePush(ctx, needOfflinePushUserIDs, msg, &offlinePushUserIDs)
		if err != nil {
			return err
		}

		if len(offlinePushUserIDs) > 0 {
			needOfflinePushUserIDs = offlinePushUserIDs
		}

		err = c.offlinePushMsg(ctx, msg, needOfflinePushUserIDs)
		if err != nil {
			log.ZError(ctx, "offlinePushMsg failed", err, "groupID", groupID, "msg", msg)
			return err
		}

	}

	return nil
}

func (c *ConsumerHandler) offlinePushMsg(ctx context.Context, msg *sdkws.MsgData, offlinePushUserIDs []string) error {
	title, content, opts, err := c.getOfflinePushInfos(msg)
	if err != nil {
		return err
	}
	err = c.offlinePusher.Push(ctx, offlinePushUserIDs, title, content, opts)
	if err != nil {
		prommetrics.MsgOfflinePushFailedCounter.Inc()
		return err
	}
	return nil
}

func (c *ConsumerHandler) filterGroupMessageOfflinePush(ctx context.Context, groupID string, msg *sdkws.MsgData,
	offlinePushUserIDs []string) (userIDs []string, err error) {

	//todo local cache Obtain the difference set through local comparison.
	needOfflinePushUserIDs, err := c.conversationRpcClient.GetConversationOfflinePushUserIDs(
		ctx, utils.GenGroupConversationID(groupID), offlinePushUserIDs)
	if err != nil {
		return nil, err
	}
	return needOfflinePushUserIDs, nil
}

func (c *ConsumerHandler) getOfflinePushInfos(msg *sdkws.MsgData) (title, content string, opts *offlinepush.Opts, err error) {
	type AtTextElem struct {
		Text       string   `json:"text,omitempty"`
		AtUserList []string `json:"atUserList,omitempty"`
		IsAtSelf   bool     `json:"isAtSelf"`
	}

	opts = &offlinepush.Opts{Signal: &offlinepush.Signal{}}
	if msg.OfflinePushInfo != nil {
		opts.IOSBadgeCount = msg.OfflinePushInfo.IOSBadgeCount
		opts.IOSPushSound = msg.OfflinePushInfo.IOSPushSound
		opts.Ex = msg.OfflinePushInfo.Ex
	}

	if msg.OfflinePushInfo != nil {
		title = msg.OfflinePushInfo.Title
		content = msg.OfflinePushInfo.Desc
	}
	if title == "" {
		switch msg.ContentType {
		case constant.Text:
			fallthrough
		case constant.Picture:
			fallthrough
		case constant.Voice:
			fallthrough
		case constant.Video:
			fallthrough
		case constant.File:
			title = constant.ContentType2PushContent[int64(msg.ContentType)]
		case constant.AtText:
			ac := AtTextElem{}
			_ = utils.JsonStringToStruct(string(msg.Content), &ac)
		case constant.SignalingNotification:
			title = constant.ContentType2PushContent[constant.SignalMsg]
		default:
			title = constant.ContentType2PushContent[constant.Common]
		}
	}
	if content == "" {
		content = title
	}
	return
}
func (c *ConsumerHandler) groupMessagesHandler(ctx context.Context, groupID string, pushToUserIDs *[]string, msg *sdkws.MsgData) (err error) {
	if len(*pushToUserIDs) == 0 {
		*pushToUserIDs, err = c.groupLocalCache.GetGroupMemberIDs(ctx, groupID)
		if err != nil {
			return err
		}
		switch msg.ContentType {
		case constant.MemberQuitNotification:
			var tips sdkws.MemberQuitTips
			if unmarshalNotificationElem(msg.Content, &tips) != nil {
				return err
			}
			if err = c.DeleteMemberAndSetConversationSeq(ctx, groupID, []string{tips.QuitUser.UserID}); err != nil {
				log.ZError(ctx, "MemberQuitNotification DeleteMemberAndSetConversationSeq", err, "groupID", groupID, "userID", tips.QuitUser.UserID)
			}
			*pushToUserIDs = append(*pushToUserIDs, tips.QuitUser.UserID)
		case constant.MemberKickedNotification:
			var tips sdkws.MemberKickedTips
			if unmarshalNotificationElem(msg.Content, &tips) != nil {
				return err
			}
			kickedUsers := utils.Slice(tips.KickedUserList, func(e *sdkws.GroupMemberFullInfo) string { return e.UserID })
			if err = c.DeleteMemberAndSetConversationSeq(ctx, groupID, kickedUsers); err != nil {
				log.ZError(ctx, "MemberKickedNotification DeleteMemberAndSetConversationSeq", err, "groupID", groupID, "userIDs", kickedUsers)
			}

			*pushToUserIDs = append(*pushToUserIDs, kickedUsers...)
		case constant.GroupDismissedNotification:
			if msgprocessor.IsNotification(msgprocessor.GetConversationIDByMsg(msg)) { // 消息先到,通知后到
				var tips sdkws.GroupDismissedTips
				if unmarshalNotificationElem(msg.Content, &tips) != nil {
					return err
				}
				log.ZInfo(ctx, "GroupDismissedNotificationInfo****", "groupID", groupID, "num", len(*pushToUserIDs), "list", pushToUserIDs)
				if len(config.Config.Manager.UserID) > 0 {
					ctx = mcontext.WithOpUserIDContext(ctx, config.Config.Manager.UserID[0])
				}
				defer func(groupID string) {
					if err = c.groupRpcClient.DismissGroup(ctx, groupID); err != nil {
						log.ZError(ctx, "DismissGroup Notification clear members", err, "groupID", groupID)
					}
				}(groupID)
			}
		}
	}
	return err
}

func (c *ConsumerHandler) DeleteMemberAndSetConversationSeq(ctx context.Context, groupID string, userIDs []string) error {
	conversationID := msgprocessor.GetConversationIDBySessionType(constant.SuperGroupChatType, groupID)
	maxSeq, err := c.msgRpcClient.GetConversationMaxSeq(ctx, conversationID)
	if err != nil {
		return err
	}
	return c.conversationRpcClient.SetConversationMaxSeq(ctx, userIDs, conversationID, maxSeq)
}

func unmarshalNotificationElem(bytes []byte, t any) error {
	var notification sdkws.NotificationElem
	if err := json.Unmarshal(bytes, &notification); err != nil {
		return err
	}

	return json.Unmarshal([]byte(notification.Detail), t)
}

func (c *ConsumerHandler) shouldPushOffline(_ context.Context, msg *sdkws.MsgData) bool {
	isOfflinePush := utils.GetSwitchFromOptions(msg.Options, constant.IsOfflinePush)
	if !isOfflinePush {
		return false
	}
	if msg.ContentType == constant.SignalingNotification {
		return false
	}
	return true
}
