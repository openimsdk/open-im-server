package push

import (
	"context"
	"sync"

	"google.golang.org/grpc"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/controller"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/localcache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	pbPush "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/push"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
)

type pushServer struct {
	pusher *Pusher
}

func Start(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error {
	rdb, err := cache.NewRedis()
	if err != nil {
		return err
	}
	cacheModel := cache.NewMsgCacheModel(rdb)
	offlinePusher := NewOfflinePusher(cacheModel)
	database := controller.NewPushDatabase(cacheModel)
	groupRpcClient := rpcclient.NewGroupRpcClient(client)
	conversationRpcClient := rpcclient.NewConversationRpcClient(client)
	msgRpcClient := rpcclient.NewMessageRpcClient(client)
	pusher := NewPusher(
		client,
		offlinePusher,
		database,
		localcache.NewGroupLocalCache(&groupRpcClient),
		localcache.NewConversationLocalCache(&conversationRpcClient),
		&conversationRpcClient,
		&groupRpcClient,
		&msgRpcClient,
	)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		pbPush.RegisterPushMsgServiceServer(server, &pushServer{
			pusher: pusher,
		})
	}()
	go func() {
		defer wg.Done()
		consumer := NewConsumer(pusher)
		consumer.initPrometheus()
		consumer.Start()
	}()
	wg.Wait()
	return nil
}

func (r *pushServer) PushMsg(ctx context.Context, pbData *pbPush.PushMsgReq) (resp *pbPush.PushMsgResp, err error) {
	switch pbData.MsgData.SessionType {
	case constant.SuperGroupChatType:
		err = r.pusher.Push2SuperGroup(ctx, pbData.MsgData.GroupID, pbData.MsgData)
	default:
		err = r.pusher.Push2User(ctx, []string{pbData.MsgData.RecvID, pbData.MsgData.SendID}, pbData.MsgData)
	}
	if err != nil {
		if err != errNoOfflinePusher {
			return nil, err
		} else {
			log.ZWarn(ctx, "offline push failed", err, "msg", pbData.String())
		}
	}
	return &pbPush.PushMsgResp{}, nil
}

func (r *pushServer) DelUserPushToken(
	ctx context.Context,
	req *pbPush.DelUserPushTokenReq,
) (resp *pbPush.DelUserPushTokenResp, err error) {
	if err = r.pusher.database.DelFcmToken(ctx, req.UserID, int(req.PlatformID)); err != nil {
		return nil, err
	}
	return &pbPush.DelUserPushTokenResp{}, nil
}
