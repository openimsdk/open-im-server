package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	conf "github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/network"
	"github.com/openimsdk/tools/utils/runtimeenv"
	"google.golang.org/grpc"
)

type Config struct {
	conf.AllConfig

	ConfigPath conf.Path
	Index      conf.Index
}

func Start(ctx context.Context, config *Config, client discovery.SvcDiscoveryRegistry, service grpc.ServiceRegistrar) error {
	apiPort, err := datautil.GetElemByIndex(config.API.Api.Ports, int(config.Index))
	if err != nil {
		return err
	}

	router, err := newGinRouter(ctx, client, config)
	if err != nil {
		return err
	}

	apiCtx, apiCancel := context.WithCancelCause(context.Background())
	done := make(chan struct{})
	go func() {
		httpServer := &http.Server{
			Handler: router,
			Addr:    net.JoinHostPort(network.GetListenIP(config.API.Api.ListenIP), strconv.Itoa(apiPort)),
		}
		go func() {
			defer close(done)
			select {
			case <-ctx.Done():
				apiCancel(fmt.Errorf("recv ctx %w", context.Cause(ctx)))
			case <-apiCtx.Done():
			}
			log.ZDebug(ctx, "api server is shutting down")
			if err := httpServer.Shutdown(context.Background()); err != nil {
				log.ZWarn(ctx, "api server shutdown err", err)
			}
		}()
		log.CInfo(ctx, "api server is init", "runtimeEnv", runtimeenv.RuntimeEnvironment(), "address", httpServer.Addr, "apiPort", apiPort)
		err := httpServer.ListenAndServe()
		if err == nil {
			err = errors.New("api done")
		}
		apiCancel(err)
	}()

	<-apiCtx.Done()
	exitCause := context.Cause(apiCtx)
	log.ZWarn(ctx, "api server exit", exitCause)
	timer := time.NewTimer(time.Second * 15)
	defer timer.Stop()
	select {
	case <-timer.C:
		log.ZWarn(ctx, "api server graceful stop timeout", nil)
	case <-done:
		log.ZDebug(ctx, "api server graceful stop done")
	}
	return exitCause
}
