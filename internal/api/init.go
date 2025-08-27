package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/tools/mw"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/network"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	kdisc "github.com/openimsdk/open-im-server/v3/pkg/common/discoveryregister"
	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/discovery/etcd"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/system/program"
)

type Config struct {
	API       config.API
	Share     config.Share
	Discovery config.Discovery
}

func Start(ctx context.Context, index int, cfg *Config) error {
	apiPort, err := datautil.GetElemByIndex(cfg.API.Api.Ports, index)
	if err != nil {
		return err
	}

	var client discovery.SvcDiscoveryRegistry

	// Determine whether zk is passed according to whether it is a clustered deployment
	client, err = kdisc.NewDiscoveryRegister(&cfg.Discovery, &cfg.Share, []string{
		cfg.Share.RpcRegisterName.MessageGateway,
	})
	if err != nil {
		return errs.WrapMsg(err, "failed to register discovery service")
	}
	client.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "round_robin")))

	var (
		netDone        = make(chan struct{}, 1)
		netErr         error
		prometheusPort int
	)

	router, err := newGinRouter(ctx, client, cfg)
	if err != nil {
		return err
	}
	registerIP, err := network.GetRpcRegisterIP("")
	if err != nil {
		return err
	}

	getAutoPort := func() (net.Listener, int, error) {
		registerAddr := net.JoinHostPort(registerIP, "0")
		listener, err := net.Listen("tcp", registerAddr)
		if err != nil {
			return nil, 0, errs.WrapMsg(err, "listen err", "registerAddr", registerAddr)
		}
		_, portStr, _ := net.SplitHostPort(listener.Addr().String())
		port, _ := strconv.Atoi(portStr)
		return listener, port, nil
	}

	if cfg.API.Prometheus.AutoSetPorts && cfg.Discovery.Enable != config.ETCD {
		return errs.New("only etcd support autoSetPorts", "RegisterName", "api").Wrap()
	}

	if cfg.API.Prometheus.Enable {
		var (
			listener net.Listener
		)

		if cfg.API.Prometheus.AutoSetPorts {
			listener, prometheusPort, err = getAutoPort()
			if err != nil {
				return err
			}

			etcdClient := client.(*etcd.SvcDiscoveryRegistryImpl).GetClient()

			prommetrics.Register(ctx, etcdClient, prommetrics.APIKeyName, registerIP, prometheusPort)
		} else {
			prometheusPort, err = datautil.GetElemByIndex(cfg.API.Prometheus.Ports, index)
			if err != nil {
				return err
			}
			listener, err = net.Listen("tcp", fmt.Sprintf(":%d", prometheusPort))
			if err != nil {
				return errs.WrapMsg(err, "listen err", "addr", fmt.Sprintf(":%d", prometheusPort))
			}
		}

		go func() {
			if err := prommetrics.ApiInit(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
				netErr = errs.WrapMsg(err, fmt.Sprintf("api prometheus start err: %d", prometheusPort))
				netDone <- struct{}{}
			}
		}()

	}
	address := net.JoinHostPort(network.GetListenIP(cfg.API.Api.ListenIP), strconv.Itoa(apiPort))

	server := http.Server{Addr: address, Handler: router}
	log.CInfo(ctx, "API server is initializing", "address", address, "apiPort", apiPort, "prometheusPort", prometheusPort)
	go func() {
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			netErr = errs.WrapMsg(err, fmt.Sprintf("api start err: %s", server.Addr))
			netDone <- struct{}{}

		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	select {
	case <-sigs:
		program.SIGTERMExit()
		err := server.Shutdown(ctx)
		if err != nil {
			return errs.WrapMsg(err, "shutdown err")
		}
	case <-netDone:
		close(netDone)
		return netErr
	}
	return nil
}
