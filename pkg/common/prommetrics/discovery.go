package prommetrics

import (
	"context"
	"fmt"

	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/jsonutil"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	APIKeyName             = "api"
	MessageTransferKeyName = "message-transfer"
)

type Target struct {
	Target string            `json:"target"`
	Labels map[string]string `json:"labels"`
}

type RespTarget struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

func BuildDiscoveryKeyPrefix(name string) string {
	return fmt.Sprintf("%s/%s/%s/", "openim", "prometheus_discovery", name)
}

func BuildDiscoveryKey(name string, host string, port int) string {
	return fmt.Sprintf("%s/%s/%s/%s:%d", "openim", "prometheus_discovery", name, host, port)
}

func BuildDefaultTarget(host string, ip int) Target {
	return Target{
		Target: fmt.Sprintf("%s:%d", host, ip),
		Labels: map[string]string{
			"namespace": "default",
		},
	}
}

func Register(ctx context.Context, etcdClient *clientv3.Client, rpcRegisterName string, registerIP string, prometheusPort int) error {
	// create lease
	leaseResp, err := etcdClient.Grant(ctx, 30)
	if err != nil {
		return errs.WrapMsg(err, "failed to create lease in etcd")
	}
	// release
	keepAliveChan, err := etcdClient.KeepAlive(ctx, leaseResp.ID)
	if err != nil {
		return errs.WrapMsg(err, "failed to keep alive lease")
	}
	// release resp
	go func() {
		for range keepAliveChan {
		}
	}()
	putOpts := []clientv3.OpOption{}
	if leaseResp != nil {
		putOpts = append(putOpts, clientv3.WithLease(leaseResp.ID))
	}
	_, err = etcdClient.Put(ctx, BuildDiscoveryKey(rpcRegisterName, registerIP, prometheusPort), jsonutil.StructToJsonString(BuildDefaultTarget(registerIP, prometheusPort)), putOpts...)
	if err != nil {
		return errs.WrapMsg(err, "etcd put err")
	}
	return nil
}
