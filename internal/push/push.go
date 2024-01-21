package push

import (
	"context"
	pbpush "github.com/OpenIMSDK/protocol/push"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/openimsdk/open-im-server/v3/internal/push/offlinepush"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/controller"
	"google.golang.org/grpc"
)

type pushServer struct {
	database      controller.PushDatabase
	disCov        discoveryregistry.SvcDiscoveryRegistry
	offlinePusher offlinepush.OfflinePusher
	pushCh        *ConsumerHandler
}

func (p pushServer) PushMsg(ctx context.Context, req *pbpush.PushMsgReq) (*pbpush.PushMsgResp, error) {
	//todo reserved Interface
	return nil, nil
}

func (p pushServer) DelUserPushToken(ctx context.Context,
	req *pbpush.DelUserPushTokenReq) (resp *pbpush.DelUserPushTokenResp, err error) {
	if err = p.database.DelFcmToken(ctx, req.UserID, int(req.PlatformID)); err != nil {
		return nil, err
	}
	return &pbpush.DelUserPushTokenResp{}, nil
}

func Start(disCov discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error {
	rdb, err := cache.NewRedis()
	if err != nil {
		return err
	}
	cacheModel := cache.NewMsgCacheModel(rdb)
	offlinePusher := offlinepush.NewOfflinePusher(cacheModel)
	database := controller.NewPushDatabase(cacheModel)

	consumer := NewConsumerHandler(offlinePusher, rdb, disCov)
	pbpush.RegisterPushMsgServiceServer(server, &pushServer{
		database:      database,
		disCov:        disCov,
		offlinePusher: offlinePusher,
		pushCh:        consumer,
	})
	go consumer.pushConsumerGroup.RegisterHandleAndConsumer(consumer)
	return nil
}
