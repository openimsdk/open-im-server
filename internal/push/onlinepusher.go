package push

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

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
	case "k8s":
		return NewK8sStaticConsistentHash(disCov, config)
	case "zookeeper":
		return NewDefaultAllNode(disCov, config)
	case "etcd":
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

type K8sStaticConsistentHash struct {
	disCov discovery.SvcDiscoveryRegistry
	config *Config
}

func NewK8sStaticConsistentHash(disCov discovery.SvcDiscoveryRegistry, config *Config) *K8sStaticConsistentHash {
	return &K8sStaticConsistentHash{disCov: disCov, config: config}
}

func (k *K8sStaticConsistentHash) GetConnsAndOnlinePush(ctx context.Context, msg *sdkws.MsgData,
	pushToUserIDs []string) (wsResults []*msggateway.SingleMsgToUserResults, err error) {

	var usersHost = make(map[string][]string)
	for _, v := range pushToUserIDs {
		tHost, err := k.disCov.GetUserIdHashGatewayHost(ctx, v)
		if err != nil {
			log.ZError(ctx, "get msg gateway hash error", err)
			return nil, err
		}
		tUsers, tbl := usersHost[tHost]
		if tbl {
			tUsers = append(tUsers, v)
			usersHost[tHost] = tUsers
		} else {
			usersHost[tHost] = []string{v}
		}
	}
	log.ZDebug(ctx, "genUsers send hosts struct:", "usersHost", usersHost)
	var usersConns = make(map[*grpc.ClientConn][]string)
	for host, userIds := range usersHost {
		tconn, _ := k.disCov.GetConn(ctx, host)
		usersConns[tconn] = userIds
	}
	var (
		mu         sync.Mutex
		wg         = errgroup.Group{}
		maxWorkers = k.config.RpcConfig.MaxConcurrentWorkers
	)
	if maxWorkers < 3 {
		maxWorkers = 3
	}
	wg.SetLimit(maxWorkers)
	for conn, userIds := range usersConns {
		tcon := conn
		tuserIds := userIds
		wg.Go(func() error {
			input := &msggateway.OnlineBatchPushOneMsgReq{MsgData: msg, PushToUserIDs: tuserIds}
			msgClient := msggateway.NewMsgGatewayClient(tcon)
			reply, err := msgClient.SuperGroupOnlineBatchPushOneMsg(ctx, input)
			if err != nil {
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
	return wsResults, nil
}
func (k *K8sStaticConsistentHash) GetOnlinePushFailedUserIDs(_ context.Context, _ *sdkws.MsgData,
	wsResults []*msggateway.SingleMsgToUserResults, _ *[]string) []string {
	var needOfflinePushUserIDs []string
	for _, v := range wsResults {
		if !v.OnlinePush {
			needOfflinePushUserIDs = append(needOfflinePushUserIDs, v.UserID)
		}
	}
	return needOfflinePushUserIDs
}
