package controller

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/protocol/push"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mq"
	"google.golang.org/protobuf/proto"
)

type PushDatabase interface {
	DelFcmToken(ctx context.Context, userID string, platformID int) error
	MsgToOfflinePushMQ(ctx context.Context, key string, userIDs []string, msg2mq *sdkws.MsgData) error
}

type pushDataBase struct {
	cache                 cache.ThirdCache
	producerToOfflinePush mq.Producer
}

func NewPushDatabase(cache cache.ThirdCache, offlinePushProducer mq.Producer) PushDatabase {
	return &pushDataBase{
		cache:                 cache,
		producerToOfflinePush: offlinePushProducer,
	}
}

func (p *pushDataBase) DelFcmToken(ctx context.Context, userID string, platformID int) error {
	return p.cache.DelFcmToken(ctx, userID, platformID)
}

func (p *pushDataBase) MsgToOfflinePushMQ(ctx context.Context, key string, userIDs []string, msg2mq *sdkws.MsgData) error {
	data, err := proto.Marshal(&push.PushMsgReq{MsgData: msg2mq, UserIDs: userIDs})
	if err != nil {
		return err
	}
	if err := p.producerToOfflinePush.SendMessage(ctx, key, data); err != nil {
		log.ZError(ctx, "message is push to offlinePush topic", err, "key", key, "userIDs", userIDs, "msg", msg2mq.String())
	}
	return err
}
