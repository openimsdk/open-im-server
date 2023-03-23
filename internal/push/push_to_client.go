package push

import (
	"context"
	"errors"
	"github.com/OpenIMSDK/Open-IM-Server/internal/push/offlinepush"
	"github.com/OpenIMSDK/Open-IM-Server/internal/push/offlinepush/fcm"
	"github.com/OpenIMSDK/Open-IM-Server/internal/push/offlinepush/getui"
	"github.com/OpenIMSDK/Open-IM-Server/internal/push/offlinepush/jpush"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/controller"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/localcache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/prome"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msggateway"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/golang/protobuf/proto"
)

type Pusher struct {
	database               controller.PushDatabase
	client                 discoveryregistry.SvcDiscoveryRegistry
	offlinePusher          offlinepush.OfflinePusher
	groupLocalCache        *localcache.GroupLocalCache
	conversationLocalCache *localcache.ConversationLocalCache
	successCount           int
}

func NewPusher(client discoveryregistry.SvcDiscoveryRegistry, offlinePusher offlinepush.OfflinePusher, database controller.PushDatabase,
	groupLocalCache *localcache.GroupLocalCache, conversationLocalCache *localcache.ConversationLocalCache) *Pusher {
	return &Pusher{
		database:               database,
		client:                 client,
		offlinePusher:          offlinePusher,
		groupLocalCache:        groupLocalCache,
		conversationLocalCache: conversationLocalCache,
	}
}

func NewOfflinePusher(cache cache.Model) offlinepush.OfflinePusher {
	var offlinePusher offlinepush.OfflinePusher
	if config.Config.Push.Getui.Enable {
		offlinePusher = getui.NewClient(cache)
	}
	if config.Config.Push.Fcm.Enable {
		offlinePusher = fcm.NewClient(cache)
	}
	if config.Config.Push.Jpns.Enable {
		offlinePusher = jpush.NewClient()
	}
	return offlinePusher
}

func (p *Pusher) MsgToUser(ctx context.Context, userID string, msg *sdkws.MsgData) error {
	var userIDs = []string{userID}
	log.ZDebug(ctx, "Get msg from msg_transfer And push msg", "userID", userID, "msg", msg.String())
	// callback
	if err := callbackOnlinePush(ctx, userIDs, msg); err != nil && err != errs.ErrCallbackContinue {
		return err
	}
	// push
	wsResults, err := p.GetConnsAndOnlinePush(ctx, msg, userIDs)
	if err != nil {
		return err
	}
	isOfflinePush := utils.GetSwitchFromOptions(msg.Options, constant.IsOfflinePush)
	//log.NewInfo(operationID, "push_result", wsResults, "sendData", msg, "isOfflinePush", isOfflinePush)
	log.ZDebug(ctx, "push_result", "ws push result", wsResults, "sendData", msg, "isOfflinePush", isOfflinePush)
	p.successCount++
	if isOfflinePush && userID != msg.SendID {
		// save invitation info for offline push
		for _, v := range wsResults {
			if v.OnlinePush {
				return nil
			}
		}
		if msg.ContentType == constant.SignalingNotification {
			isSend, err := p.database.HandleSignalInvite(ctx, msg, userID)
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
		err = p.offlinePushMsg(ctx, userID, msg, userIDs)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Pusher) MsgToSuperGroupUser(ctx context.Context, groupID string, msg *sdkws.MsgData) (err error) {
	operationID := mcontext.GetOperationID(ctx)
	log.Debug(operationID, "Get super group msg from msg_transfer And push msg", msg.String(), groupID)
	var pushToUserIDs []string
	if err := callbackBeforeSuperGroupOnlinePush(ctx, groupID, msg, &pushToUserIDs); err != nil && err != errs.ErrCallbackContinue {
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
			err = p.offlinePushMsg(ctx, groupID, msg, offlinePushUserIDs)
			if err != nil {
				log.NewError(operationID, "offlinePushMsg failed", groupID)
				return err
			}
			_, err := p.GetConnsAndOnlinePush(ctx, msg, utils.IntersectString(needOfflinePushUserIDs, WebAndPcBackgroundUserIDs))
			if err != nil {
				log.NewError(operationID, "offlinePushMsg failed", groupID)
				return err
			}
		}
	}
	return nil
}

func (p *Pusher) GetConnsAndOnlinePush(ctx context.Context, msg *sdkws.MsgData, pushToUserIDs []string) (wsResults []*msggateway.SingleMsgToUserResults, err error) {
	conns, err := p.client.GetConns(config.Config.RpcRegisterName.OpenImMessageGatewayName)
	if err != nil {
		return nil, err
	}
	//Online push message
	for _, v := range conns {
		msgClient := msggateway.NewMsgGatewayClient(v)
		reply, err := msgClient.SuperGroupOnlineBatchPushOneMsg(ctx, &msggateway.OnlineBatchPushOneMsgReq{MsgData: msg, PushToUserIDs: pushToUserIDs})
		if err != nil {

			continue
		}
		if reply != nil && reply.SinglePushResult != nil {
			wsResults = append(wsResults, reply.SinglePushResult...)
		}
	}
	return wsResults, nil
}

func (p *Pusher) offlinePushMsg(ctx context.Context, sourceID string, msg *sdkws.MsgData, offlinePushUserIDs []string) error {
	title, content, opts, err := p.getOfflinePushInfos(sourceID, msg)
	if err != nil {
		return err
	}
	err = p.offlinePusher.Push(ctx, offlinePushUserIDs, title, content, opts)
	if err != nil {
		prome.Inc(prome.MsgOfflinePushFailedCounter)
		return err
	}
	prome.Inc(prome.MsgOfflinePushSuccessCounter)
	return nil
}

func (p *Pusher) GetOfflinePushOpts(msg *sdkws.MsgData) (opts *offlinepush.Opts, err error) {
	opts = &offlinepush.Opts{}
	if msg.ContentType > constant.SignalingNotificationBegin && msg.ContentType < constant.SignalingNotificationEnd {
		req := &sdkws.SignalReq{}
		if err := proto.Unmarshal(msg.Content, req); err != nil {
			return nil, utils.Wrap(err, "")
		}
		switch req.Payload.(type) {
		case *sdkws.SignalReq_Invite, *sdkws.SignalReq_InviteInGroup:
			opts.Signal = &offlinepush.Signal{ClientMsgID: msg.ClientMsgID}
		}
	}
	if msg.OfflinePushInfo != nil {
		opts.IOSBadgeCount = msg.OfflinePushInfo.IOSBadgeCount
		opts.IOSPushSound = msg.OfflinePushInfo.IOSPushSound
		opts.Ex = msg.OfflinePushInfo.Ex
	}
	return opts, nil
}

func (p *Pusher) getOfflinePushInfos(sourceID string, msg *sdkws.MsgData) (title, content string, opts *offlinepush.Opts, err error) {
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
