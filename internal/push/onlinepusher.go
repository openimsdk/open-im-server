package push

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"

	conf "github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/protocol/msggateway"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
)

type OnlinePusher interface {
	GetConnsAndOnlinePush(ctx context.Context, msg *sdkws.MsgData,
		pushToUserIDs []string) (wsResults []*msggateway.SingleMsgToUserResults, err error)
	GetOnlinePushFailedUserIDs(ctx context.Context, msg *sdkws.MsgData, wsResults []*msggateway.SingleMsgToUserResults,
		pushToUserIDs *[]string) []string
}

type emptyOnlinePusher struct{}

func newEmptyOnlinePusher() *emptyOnlinePusher {
	return &emptyOnlinePusher{}
}

func (emptyOnlinePusher) GetConnsAndOnlinePush(ctx context.Context, msg *sdkws.MsgData,
	pushToUserIDs []string) (wsResults []*msggateway.SingleMsgToUserResults, err error) {
	log.ZInfo(ctx, "emptyOnlinePusher GetConnsAndOnlinePush", nil)
	return nil, nil
}
func (u emptyOnlinePusher) GetOnlinePushFailedUserIDs(ctx context.Context, msg *sdkws.MsgData,
	wsResults []*msggateway.SingleMsgToUserResults, pushToUserIDs *[]string) []string {
	log.ZInfo(ctx, "emptyOnlinePusher GetOnlinePushFailedUserIDs", nil)
	return nil
}

func NewOnlinePusher(disCov discovery.SvcDiscoveryRegistry, config *Config) OnlinePusher {
	switch config.Discovery.Enable {
	case conf.KUBERNETES:
		return NewDefaultAllNode(disCov, config)
	case "zookeeper":
		return NewDefaultAllNode(disCov, config)
	case conf.ETCD:
		return NewDefaultAllNode(disCov, config)
	default:
		return newEmptyOnlinePusher()
	}
}

type DefaultAllNode struct {
	disCov discovery.SvcDiscoveryRegistry
	config *Config
}

func NewDefaultAllNode(disCov discovery.SvcDiscoveryRegistry, config *Config) *DefaultAllNode {
	return &DefaultAllNode{disCov: disCov, config: config}
}

func (d *DefaultAllNode) GetConnsAndOnlinePush(ctx context.Context, msg *sdkws.MsgData,
	pushToUserIDs []string) (wsResults []*msggateway.SingleMsgToUserResults, err error) {
	conns, err := d.disCov.GetConns(ctx, d.config.Share.RpcRegisterName.MessageGateway)
	if len(conns) == 0 {
		log.ZWarn(ctx, "get gateway conn 0 ", nil)
	} else {
		log.ZDebug(ctx, "get gateway conn", "conn length", len(conns))
	}

	if err != nil {
		return nil, err
	}

	var (
		mu         sync.Mutex
		wg         = errgroup.Group{}
		input      = &msggateway.OnlineBatchPushOneMsgReq{MsgData: msg, PushToUserIDs: pushToUserIDs}
		maxWorkers = d.config.RpcConfig.MaxConcurrentWorkers
	)

	if maxWorkers < 3 {
		maxWorkers = 3
	}

	wg.SetLimit(maxWorkers)

	// Online push message
	for _, conn := range conns {
		conn := conn // loop var safe
		ctx := ctx
		wg.Go(func() error {
			msgClient := msggateway.NewMsgGatewayClient(conn)
			reply, err := msgClient.SuperGroupOnlineBatchPushOneMsg(ctx, input)
			if err != nil {
				log.ZError(ctx, "SuperGroupOnlineBatchPushOneMsg ", err, "req:", input.String())
				return nil
			}

			log.ZDebug(ctx, "push result", "reply", reply)
			if reply != nil && reply.SinglePushResult != nil {
				mu.Lock()
				wsResults = append(wsResults, reply.SinglePushResult...)
				mu.Unlock()
			}

			return nil
		})
	}

	_ = wg.Wait()

	// always return nil
	return wsResults, nil
}

func (d *DefaultAllNode) GetOnlinePushFailedUserIDs(_ context.Context, msg *sdkws.MsgData,
	wsResults []*msggateway.SingleMsgToUserResults, pushToUserIDs *[]string) []string {

	onlineSuccessUserIDs := []string{msg.SendID}
	for _, v := range wsResults {
		//message sender do not need offline push
		if msg.SendID == v.UserID {
			continue
		}
		// mobile online push success
		if v.OnlinePush {
			onlineSuccessUserIDs = append(onlineSuccessUserIDs, v.UserID)
		}

	}

	return datautil.SliceSub(*pushToUserIDs, onlineSuccessUserIDs)
}
