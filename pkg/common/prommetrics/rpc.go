package prommetrics

import (
	gp "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net"
	"strconv"
)

const rpcPath = commonPath

var (
	grpcMetrics *gp.ServerMetrics
	rpcCounter  = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rpc_count",
			Help: "Total number of RPC calls",
		},
		[]string{"name", "path", "code"},
	)
)

func RpcInit(cs []prometheus.Collector, listener net.Listener) error {
	reg := prometheus.NewRegistry()
	cs = append(append(
		baseCollector,
		rpcCounter,
	), cs...)
	return Init(reg, listener, rpcPath, promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}), cs...)
}

func RPCCall(name string, path string, code int) {
	rpcCounter.With(prometheus.Labels{"name": name, "path": path, "code": strconv.Itoa(code)}).Inc()
}

func GetGrpcServerMetrics() *gp.ServerMetrics {
	if grpcMetrics == nil {
		grpcMetrics = gp.NewServerMetrics()
		grpcMetrics.EnableHandlingTimeHistogram()
	}
	return grpcMetrics
}

func GetGrpcCusMetrics(registerName string, share *config.Share) []prometheus.Collector {
	switch registerName {
	case share.RpcRegisterName.MessageGateway:
		return []prometheus.Collector{OnlineUserGauge}
	case share.RpcRegisterName.Msg:
		return []prometheus.Collector{
			SingleChatMsgProcessSuccessCounter,
			SingleChatMsgProcessFailedCounter,
			GroupChatMsgProcessSuccessCounter,
			GroupChatMsgProcessFailedCounter,
		}
	case share.RpcRegisterName.Push:
		return []prometheus.Collector{
			MsgOfflinePushFailedCounter,
			MsgLoneTimePushCounter,
		}
	case share.RpcRegisterName.Auth:
		return []prometheus.Collector{UserLoginCounter}
	case share.RpcRegisterName.User:
		return []prometheus.Collector{UserRegisterCounter}
	default:
		return nil
	}
}
