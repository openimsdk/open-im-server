package push

import (
	"context"
	"math/rand"
	"strconv"

	"github.com/openimsdk/open-im-server/v3/internal/push/offlinepush"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/redis"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/dbbuild"
	"github.com/openimsdk/open-im-server/v3/pkg/mqbuild"
	pbpush "github.com/openimsdk/protocol/push"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"google.golang.org/grpc"
)

type pushServer struct {
	pbpush.UnimplementedPushMsgServiceServer
	database      controller.PushDatabase
	disCov        discovery.Conn
	offlinePusher offlinepush.OfflinePusher
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
	FcmConfigPath      config.Path
}

func (p pushServer) DelUserPushToken(ctx context.Context,
	req *pbpush.DelUserPushTokenReq) (resp *pbpush.DelUserPushTokenResp, err error) {
	if err = p.database.DelFcmToken(ctx, req.UserID, int(req.PlatformID)); err != nil {
		return nil, err
	}
	return &pbpush.DelUserPushTokenResp{}, nil
}

func Start(ctx context.Context, config *Config, client discovery.Conn, server grpc.ServiceRegistrar) error {
	dbb := dbbuild.NewBuilder(nil, &config.RedisConfig)
	rdb, err := dbb.Redis(ctx)
	if err != nil {
		return err
	}
	cacheModel := redis.NewThirdCache(rdb)
	offlinePusher, err := offlinepush.NewOfflinePusher(&config.RpcConfig, cacheModel, string(config.FcmConfigPath))
	if err != nil {
		return err
	}
	builder := mqbuild.NewBuilder(&config.KafkaConfig)

	offlinePushProducer, err := builder.GetTopicProducer(ctx, config.KafkaConfig.ToOfflinePushTopic)
	if err != nil {
		return err
	}
	database := controller.NewPushDatabase(cacheModel, offlinePushProducer)

	pushConsumer, err := builder.GetTopicConsumer(ctx, config.KafkaConfig.ToPushTopic)
	if err != nil {
		return err
	}
	offlinePushConsumer, err := builder.GetTopicConsumer(ctx, config.KafkaConfig.ToOfflinePushTopic)
	if err != nil {
		return err
	}

	pushHandler, err := NewConsumerHandler(ctx, config, database, offlinePusher, rdb, client)
	if err != nil {
		return err
	}

	offlineHandler := NewOfflinePushConsumerHandler(offlinePusher)

	pbpush.RegisterPushMsgServiceServer(server, &pushServer{
		database:      database,
		disCov:        client,
		offlinePusher: offlinePusher,
	})

	go func() {
		pushHandler.WaitCache()
		fn := func(ctx context.Context, key string, value []byte) error {
			pushHandler.HandleMs2PsChat(ctx, value)
			return nil
		}
		consumerCtx := mcontext.SetOperationID(context.Background(), "push_"+strconv.Itoa(int(rand.Uint32())))
		log.ZInfo(consumerCtx, "begin consume messages")
		for {
			if err := pushConsumer.Subscribe(consumerCtx, fn); err != nil {
				log.ZError(consumerCtx, "subscribe err", err)
				return
			}
		}
	}()

	go func() {
		fn := func(ctx context.Context, key string, value []byte) error {
			offlineHandler.HandleMsg2OfflinePush(ctx, value)
			return nil
		}
		consumerCtx := mcontext.SetOperationID(context.Background(), "push_"+strconv.Itoa(int(rand.Uint32())))
		log.ZInfo(consumerCtx, "begin consume messages")
		for {
			if err := offlinePushConsumer.Subscribe(consumerCtx, fn); err != nil {
				log.ZError(consumerCtx, "subscribe err", err)
				return
			}
		}
	}()

	return nil
}
