package push

import (
	"context"
	"github.com/OpenIMSDK/protocol/msggateway"
	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/utils"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"os"
	"sync"
)

const (
	ENVNAME    = "ENVS_DISCOVERY"
	KUBERNETES = "k8s"
	ZOOKEEPER  = "zookeeper"
)

type OnlinePusher interface {
	GetConnsAndOnlinePush(ctx context.Context, msg *sdkws.MsgData,
		pushToUserIDs []string) (wsResults []*msggateway.SingleMsgToUserResults, err error)
	GetOnlinePushFailedUserIDs(ctx context.Context, msg *sdkws.MsgData, wsResults []*msggateway.SingleMsgToUserResults,
		pushToUserIDs *[]string) []string
}

type emptyOnlinePUsher struct{}

func newEmptyOnlinePUsher() *emptyOnlinePUsher {
	return &emptyOnlinePUsher{}
}

func (emptyOnlinePUsher) GetConnsAndOnlinePush(ctx context.Context, msg *sdkws.MsgData,
	pushToUserIDs []string) (wsResults []*msggateway.SingleMsgToUserResults, err error) {
	log.ZWarn(ctx, "emptyOnlinePUsher GetConnsAndOnlinePush", nil)
	return nil, nil
}
func (u emptyOnlinePUsher) GetOnlinePushFailedUserIDs(ctx context.Context, msg *sdkws.MsgData,
	wsResults []*msggateway.SingleMsgToUserResults, pushToUserIDs *[]string) []string {
	log.ZWarn(ctx, "emptyOnlinePUsher GetOnlinePushFailedUserIDs", nil)
	return nil
}

func NewOnlinePusher(disCov discoveryregistry.SvcDiscoveryRegistry) OnlinePusher {
	var envType string
	if value := os.Getenv(ENVNAME); value != "" {
		envType = os.Getenv(ENVNAME)
	} else {
		envType = config.Config.Envs.Discovery
	}
	switch envType {
	case KUBERNETES:
		return NewK8sStaticConsistentHash(disCov)
	case ZOOKEEPER:
		return NewDefaultAllNode(disCov)
	default:
		return newEmptyOnlinePUsher()
	}
}

type DefaultAllNode struct {
	disCov discoveryregistry.SvcDiscoveryRegistry
}

func NewDefaultAllNode(disCov discoveryregistry.SvcDiscoveryRegistry) *DefaultAllNode {
	return &DefaultAllNode{disCov: disCov}
}

func (d *DefaultAllNode) GetConnsAndOnlinePush(ctx context.Context, msg *sdkws.MsgData,
	pushToUserIDs []string) (wsResults []*msggateway.SingleMsgToUserResults, err error) {
	conns, err := d.disCov.GetConns(ctx, config.Config.RpcRegisterName.OpenImMessageGatewayName)
	log.ZDebug(ctx, "get gateway conn", "conn length", len(conns))
	if err != nil {
		return nil, err
	}

	var (
		mu         sync.Mutex
		wg         = errgroup.Group{}
		input      = &msggateway.OnlineBatchPushOneMsgReq{MsgData: msg, PushToUserIDs: pushToUserIDs}
		maxWorkers = config.Config.Push.MaxConcurrentWorkers
	)

	if maxWorkers < 3 {
		maxWorkers = 3
	}

	wg.SetLimit(maxWorkers)

	// Online push message
	for _, conn := range conns {
		conn := conn // loop var safe
		wg.Go(func() error {
			msgClient := msggateway.NewMsgGatewayClient(conn)
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

	return utils.SliceSub(*pushToUserIDs, onlineSuccessUserIDs)
}

type K8sStaticConsistentHash struct {
	disCov discoveryregistry.SvcDiscoveryRegistry
}

func NewK8sStaticConsistentHash(disCov discoveryregistry.SvcDiscoveryRegistry) *K8sStaticConsistentHash {
	return &K8sStaticConsistentHash{disCov: disCov}
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
		maxWorkers = config.Config.Push.MaxConcurrentWorkers
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
