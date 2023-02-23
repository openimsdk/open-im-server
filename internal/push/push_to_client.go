/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@open-im.io").
** time(2021/3/5 14:31).
 */
package push

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/cache"
	"Open_IM/pkg/common/db/localcache"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/prome"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/discoveryregistry"
	msggateway "Open_IM/pkg/proto/msggateway"
	pbRtc "Open_IM/pkg/proto/rtc"
	"Open_IM/pkg/proto/sdkws"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"github.com/golang/protobuf/proto"
)

type Pusher struct {
	cache                  cache.MsgCache
	client                 discoveryregistry.SvcDiscoveryRegistry
	offlinePusher          OfflinePusher
	groupLocalCache        localcache.GroupLocalCache
	conversationLocalCache localcache.ConversationLocalCache
	successCount           int
}

func NewPusher(cache cache.MsgCache, client discoveryregistry.SvcDiscoveryRegistry, offlinePusher OfflinePusher) *Pusher {
	return &Pusher{
		cache:         cache,
		client:        client,
		offlinePusher: offlinePusher,
	}
}

func (p *Pusher) MsgToUser(ctx context.Context, userID string, msg *sdkws.MsgData) error {
	operationID := tracelog.GetOperationID(ctx)
	var userIDs = []string{userID}
	log.Debug(operationID, "Get msg from msg_transfer And push msg", msg.String(), userID)
	// callback
	if err := callbackOnlinePush(ctx, userIDs, msg); err != nil {
		return err
	}
	// push
	wsResults, err := p.GetConnsAndOnlinePush(ctx, msg, userIDs)
	if err != nil {
		return err
	}
	isOfflinePush := utils.GetSwitchFromOptions(msg.Options, constant.IsOfflinePush)
	log.NewInfo(operationID, "push_result", wsResults, "sendData", msg, "isOfflinePush", isOfflinePush)
	p.successCount++
	if isOfflinePush && userID != msg.SendID {
		// save invitation info for offline push
		for _, v := range wsResults {
			if v.OnlinePush {
				return nil
			}
		}
		if msg.ContentType == constant.SignalingNotification {
			isSend, err := p.cache.HandleSignalInfo(ctx, msg, userID)
			if err != nil {
				return err
			}
			if !isSend {
				return nil
			}
		}

		if err := callbackOfflinePush(ctx, userIDs, msg, &[]string{}); err != nil {
			return err
		}
		err = p.OfflinePushMsg(ctx, userID, msg, userIDs)
		if err != nil {
			log.NewError(operationID, "OfflinePushMsg failed", userID)
			return err
		}
	}
	return nil
}

func (p *Pusher) MsgToSuperGroupUser(ctx context.Context, groupID string, msg *sdkws.MsgData) (err error) {
	operationID := tracelog.GetOperationID(ctx)
	log.Debug(operationID, "Get super group msg from msg_transfer And push msg", msg.String(), groupID)
	var pushToUserIDs []string
	if err := callbackBeforeSuperGroupOnlinePush(ctx, groupID, msg, &pushToUserIDs); err != nil {
		return err
	}
	if len(pushToUserIDs) == 0 {
		pushToUserIDs, err = p.groupLocalCache.GetGroupMemberIDs(ctx, groupID)
		if err != nil {
			return err
		}
	}
	wsResults, err := p.GetConnsAndOnlinePush(ctx, msg, pushToUserIDs)
	if err != nil {
		return err
	}
	log.Debug(operationID, "push_result", wsResults, "sendData", msg)
	p.successCount++
	isOfflinePush := utils.GetSwitchFromOptions(msg.Options, constant.IsOfflinePush)
	if isOfflinePush {
		var onlineSuccessUserIDs []string
		var WebAndPcBackgroundUserIDs []string
		onlineSuccessUserIDs = append(onlineSuccessUserIDs, msg.SendID)
		for _, v := range wsResults {
			if v.OnlinePush && v.UserID != msg.SendID {
				onlineSuccessUserIDs = append(onlineSuccessUserIDs, v.UserID)
			}
			if !v.OnlinePush {
				if len(v.Resp) != 0 {
					for _, singleResult := range v.Resp {
						if singleResult.ResultCode == -2 {
							if constant.PlatformIDToName(int(singleResult.RecvPlatFormID)) == constant.TerminalPC ||
								singleResult.RecvPlatFormID == constant.WebPlatformID {
								WebAndPcBackgroundUserIDs = append(WebAndPcBackgroundUserIDs, v.UserID)
							}
						}
					}
				}
			}
		}
		needOfflinePushUserIDs := utils.DifferenceString(onlineSuccessUserIDs, pushToUserIDs)
		if msg.ContentType != constant.SignalingNotification {
			notNotificationUserIDs, err := p.conversationLocalCache.GetRecvMsgNotNotifyUserIDs(ctx, groupID)
			if err != nil {
				log.Error(operationID, utils.GetSelfFuncName(), "GetRecvMsgNotNotifyUserIDs failed", groupID)
				return err
			}
			needOfflinePushUserIDs = utils.DifferenceString(notNotificationUserIDs, needOfflinePushUserIDs)
		}
		//Use offline push messaging
		if len(needOfflinePushUserIDs) > 0 {
			var offlinePushUserIDs []string
			err = callbackOfflinePush(ctx, needOfflinePushUserIDs, msg, &offlinePushUserIDs)
			if err != nil {
				return err
			}
			if len(offlinePushUserIDs) > 0 {
				needOfflinePushUserIDs = offlinePushUserIDs
			}
			err = p.OfflinePushMsg(ctx, groupID, msg, offlinePushUserIDs)
			if err != nil {
				log.NewError(operationID, "OfflinePushMsg failed", groupID)
				return err
			}
			_, err := p.GetConnsAndOnlinePush(ctx, msg, utils.IntersectString(needOfflinePushUserIDs, WebAndPcBackgroundUserIDs))
			if err != nil {
				log.NewError(operationID, "OfflinePushMsg failed", groupID)
				return err
			}
		}
	}
	return nil
}

func (p *Pusher) GetConnsAndOnlinePush(ctx context.Context, msg *sdkws.MsgData, pushToUserIDs []string) (wsResults []*msggateway.SingelMsgToUserResultList, err error) {
	conns, err := p.client.GetConns(config.Config.RpcRegisterName.OpenImMessageGatewayName)
	if err != nil {
		return nil, err
	}
	//Online push message
	for _, v := range conns {
		msgClient := msggateway.NewRelayClient(v)
		reply, err := msgClient.SuperGroupOnlineBatchPushOneMsg(ctx, &msggateway.OnlineBatchPushOneMsgReq{OperationID: tracelog.GetOperationID(ctx), MsgData: msg, PushToUserIDList: pushToUserIDs})
		if err != nil {
			log.NewError(tracelog.GetOperationID(ctx), msg, len(pushToUserIDs), "err", err)
			continue
		}
		if reply != nil && reply.SinglePushResult != nil {
			wsResults = append(wsResults, reply.SinglePushResult...)
		}
	}
	return wsResults, nil
}

func (p *Pusher) OfflinePushMsg(ctx context.Context, sourceID string, msg *sdkws.MsgData, offlinePushUserIDs []string) error {
	title, content, opts, err := p.GetOfflinePushInfos(sourceID, msg)
	if err != nil {
		return err
	}
	err = p.offlinePusher.Push(ctx, offlinePushUserIDs, title, content, opts)
	if err != nil {
		prome.PromeInc(prome.MsgOfflinePushFailedCounter)
		return err
	}
	prome.PromeInc(prome.MsgOfflinePushSuccessCounter)
	return nil
}

func (p *Pusher) GetOfflinePushOpts(msg *sdkws.MsgData) (opts *Opts, err error) {
	opts = &Opts{}
	if msg.ContentType > constant.SignalingNotificationBegin && msg.ContentType < constant.SignalingNotificationEnd {
		req := &pbRtc.SignalReq{}
		if err := proto.Unmarshal(msg.Content, req); err != nil {
			return nil, utils.Wrap(err, "")
		}
		switch req.Payload.(type) {
		case *pbRtc.SignalReq_Invite, *pbRtc.SignalReq_InviteInGroup:
			opts.Signal = &Signal{ClientMsgID: msg.ClientMsgID}
		}
	}
	if msg.OfflinePushInfo != nil {
		opts.IOSBadgeCount = msg.OfflinePushInfo.IOSBadgeCount
		opts.IOSPushSound = msg.OfflinePushInfo.IOSPushSound
		opts.Ex = msg.OfflinePushInfo.Ex
	}
	return opts, nil
}

func (p *Pusher) GetOfflinePushInfos(sourceID string, msg *sdkws.MsgData) (title, content string, opts *Opts, err error) {
	if p.offlinePusher == nil {
		err = errors.New("no offlinePusher is configured")
		return
	}
	type AtContent struct {
		Text       string   `json:"text"`
		AtUserList []string `json:"atUserList"`
		IsAtSelf   bool     `json:"isAtSelf"`
	}
	opts, err = p.GetOfflinePushOpts(msg)
	if err != nil {
		return
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
			a := AtContent{}
			_ = utils.JsonStringToStruct(string(msg.Content), &a)
			if utils.IsContain(sourceID, a.AtUserList) {
				title = constant.ContentType2PushContent[constant.AtText] + constant.ContentType2PushContent[constant.Common]
			} else {
				title = constant.ContentType2PushContent[constant.GroupMsg]
			}
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
