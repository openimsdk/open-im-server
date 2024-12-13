package push

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/internal/push/offlinepush"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/redis"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/controller"
	pbpush "github.com/openimsdk/protocol/push"
	"github.com/openimsdk/tools/db/redisutil"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/utils/runtimeenv"
	"google.golang.org/grpc"
)

type pushServer struct {
	pbpush.UnimplementedPushMsgServiceServer
	database      controller.PushDatabase
	disCov        discovery.SvcDiscoveryRegistry
	offlinePusher offlinepush.OfflinePusher
	pushCh        *ConsumerHandler
	offlinePushCh *OfflinePushConsumerHandler
}

type Config struct {
	RpcConfig          config.Push
	RedisConfig        config.Redis
	KafkaConfig        config.Kafka
	NotificationConfig config.Notification
	Share              config.Share
	WebhooksConfig     config.Webhooks
	LocalCacheConfig   config.LocalCache
	Discovery          config.Discovery
	FcmConfigPath      string
	
	runTimeEnv         string
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

func Start(ctx context.Context, config *Config, client discovery.SvcDiscoveryRegistry, server *grpc.Server) error {
	config.runTimeEnv = runtimeenv.PrintRuntimeEnvironment()

	rdb, err := redisutil.NewRedisClient(ctx, config.RedisConfig.Build())
	if err != nil {
		return err
	}
	cacheModel := redis.NewThirdCache(rdb)
	offlinePusher, err := offlinepush.NewOfflinePusher(&config.RpcConfig, cacheModel, config.FcmConfigPath)
	if err != nil {
		return err
	}

	database := controller.NewPushDatabase(cacheModel, &config.KafkaConfig)

	consumer, err := NewConsumerHandler(config, database, offlinePusher, rdb, client)
	if err != nil {
		return err
	}

	offlinePushConsumer, err := NewOfflinePushConsumerHandler(config, offlinePusher)
	if err != nil {
		return err
	}

	pbpush.RegisterPushMsgServiceServer(server, &pushServer{
		database:      database,
		disCov:        client,
		offlinePusher: offlinePusher,
		pushCh:        consumer,
		offlinePushCh: offlinePushConsumer,
	})

	go consumer.pushConsumerGroup.RegisterHandleAndConsumer(ctx, consumer)

	go offlinePushConsumer.OfflinePushConsumerGroup.RegisterHandleAndConsumer(ctx, offlinePushConsumer)

	return nil
}
