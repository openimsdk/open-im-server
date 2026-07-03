package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/openimsdk/open-im-server/v3/internal/api"
	"github.com/openimsdk/open-im-server/v3/internal/msggateway"
	"github.com/openimsdk/open-im-server/v3/internal/msgtransfer"
	"github.com/openimsdk/open-im-server/v3/internal/push"
	"github.com/openimsdk/open-im-server/v3/internal/rpc/auth"
	"github.com/openimsdk/open-im-server/v3/internal/rpc/conversation"
	"github.com/openimsdk/open-im-server/v3/internal/rpc/group"
	"github.com/openimsdk/open-im-server/v3/internal/rpc/msg"
	"github.com/openimsdk/open-im-server/v3/internal/rpc/relation"
	"github.com/openimsdk/open-im-server/v3/internal/rpc/third"
	"github.com/openimsdk/open-im-server/v3/internal/rpc/user"
	"github.com/openimsdk/open-im-server/v3/internal/tools/cron"
	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/redis"
	"github.com/openimsdk/open-im-server/v3/pkg/dbbuild"
	"github.com/openimsdk/open-im-server/v3/version"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/discovery/inprocess"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/system/program"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/network"
)

func init() {
	config.SetStandalone()
	prommetrics.RegistryAll()
	inprocess.SetContextAdminFunc(authverify.IsAdmin)
}

func main() {
	var (
		configPath string
		index      int
	)
	flag.StringVar(&configPath, "c", "", "config path")
	flag.IntVar(&index, "i", 0, "start index")
	flag.Parse()
	if configPath == "" {
		_, _ = fmt.Fprintln(os.Stderr, "config path is empty")
		os.Exit(1)
		return
	}
	cmd := newCmds(configPath, index)
	putCmd(cmd, false, auth.Start)
	putCmd(cmd, false, conversation.Start)
	putCmd(cmd, false, relation.Start)
	putCmd(cmd, false, group.Start)
	putCmd(cmd, false, msg.Start)
	putCmd(cmd, false, third.Start)
	putCmd(cmd, false, user.Start)
	putCmd(cmd, false, push.Start)
	putCmd(cmd, true, msggateway.Start)
	putCmd(cmd, true, msgtransfer.Start)
	putCmd(cmd, true, api.Start)
	putCmd(cmd, true, cron.Start)
	putCmd(cmd, true, startRedisServerRegister)
	ctx := context.Background()
	if err := cmd.run(ctx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "server exit %s", err)
		os.Exit(1)
		return
	}
}

func newCmds(confPath string, index int) *cmds {
	return &cmds{confPath: confPath, index: index}
}

type cmdName struct {
	Name  string
	Func  func(ctx context.Context) error
	Block bool
}
type cmds struct {
	confPath string
	index    int
	cmds     []cmdName
	config   config.AllConfig
	conf     map[string]reflect.Value
}

func (x *cmds) getTypePath(typ reflect.Type) string {
	return path.Join(typ.PkgPath(), typ.Name())
}

func (x *cmds) initDiscovery() {
	x.config.Discovery.Enable = "standalone"
	vof := reflect.ValueOf(&x.config.Discovery.RpcService).Elem()
	tof := reflect.TypeOf(&x.config.Discovery.RpcService).Elem()
	num := tof.NumField()
	for i := 0; i < num; i++ {
		field := tof.Field(i)
		if !field.IsExported() {
			continue
		}
		if field.Type.Kind() != reflect.String {
			continue
		}
		vof.Field(i).SetString(field.Name)
	}
}

func (x *cmds) initAllConfig() error {
	x.conf = make(map[string]reflect.Value)
	if err := x.loadShareConfig(); err != nil {
		return err
	}
	vof := reflect.ValueOf(&x.config).Elem()
	num := vof.NumField()
	for i := 0; i < num; i++ {
		field := vof.Field(i)
		for ptr := true; ptr; {
			if field.Kind() == reflect.Ptr {
				field = field.Elem()
			} else {
				ptr = false
			}
		}
		x.conf[x.getTypePath(field.Type())] = field
		val := field.Addr().Interface()
		// Fields without a config file (e.g. Timer) are still registered above
		// for parseConf distribution, but there is nothing to load for them.
		fc, ok := val.(config.FileConfig)
		if !ok {
			continue
		}
		name := fc.GetConfigFileName()
		if name == config.ShareFileName || x.skipKafkaConfig(name) {
			continue
		}
		if err := x.loadFileConfig(name, val); err != nil {
			return err
		}
	}
	x.initDiscovery()
	x.config.Redis.Disable = false
	x.config.LocalCache = config.LocalCache{}
	config.InitNotification(&x.config.Notification)
	return nil
}

func (x *cmds) loadShareConfig() error {
	return x.loadFileConfig(config.ShareFileName, &x.config.Share)
}

func (x *cmds) skipKafkaConfig(name string) bool {
	return name == config.KafkaConfigFileName &&
		config.NormalizeQueueEngine(x.config.Share.Queue.Engine) != config.QueueEngineKafka
}

func (x *cmds) loadFileConfig(name string, val any) error {
	confData, err := os.ReadFile(filepath.Join(x.confPath, name))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	v := viper.New()
	v.SetConfigType(config.StructTagName)
	if err := v.ReadConfig(bytes.NewReader(confData)); err != nil {
		return err
	}
	return v.Unmarshal(val, config.ApplyDecoderConfig)
}

func (x *cmds) parseConf(conf any) error {
	vof := reflect.ValueOf(conf)
	for {
		if vof.Kind() == reflect.Ptr {
			vof = vof.Elem()
		} else {
			break
		}
	}
	tof := vof.Type()
	numField := vof.NumField()
	for i := 0; i < numField; i++ {
		typeField := tof.Field(i)
		if !typeField.IsExported() {
			continue
		}
		field := vof.Field(i)
		pkt := x.getTypePath(field.Type())
		val, ok := x.conf[pkt]
		if !ok {
			switch field.Interface().(type) {
			case config.Index:
				field.Set(reflect.ValueOf(config.Index(x.index)))
			case config.Path:
				field.SetString(x.confPath)
			case config.AllConfig:
				field.Set(reflect.ValueOf(x.config))
			case *config.AllConfig:
				field.Set(reflect.ValueOf(&x.config))
			default:
				return fmt.Errorf("config field %s %s not found", vof.Type().Name(), typeField.Name)
			}
			continue
		}
		field.Set(val)
	}
	return nil
}

func (x *cmds) add(name string, block bool, fn func(ctx context.Context) error) {
	x.cmds = append(x.cmds, cmdName{Name: name, Block: block, Func: fn})
}

func (x *cmds) initLog() error {
	conf := x.config.Log
	if err := log.InitLoggerFromConfig(
		"openim-service-log",
		program.GetProcessName(),
		"", "",
		conf.RemainLogLevel,
		conf.IsStdout,
		conf.IsJson,
		conf.StorageLocation,
		conf.RemainRotationCount,
		conf.RotationTime,
		strings.TrimSpace(version.Version),
		conf.IsSimplify,
	); err != nil {
		return err
	}
	return nil

}

func (x *cmds) run(ctx context.Context) error {
	if len(x.cmds) == 0 {
		return fmt.Errorf("no command to run")
	}
	if err := x.initAllConfig(); err != nil {
		return err
	}
	if err := x.initLog(); err != nil {
		return err
	}

	ctx, cancel := context.WithCancelCause(ctx)

	go func() {
		<-ctx.Done()
		log.ZError(ctx, "context server exit cause", context.Cause(ctx))
	}()

	if prometheus := x.config.API.Prometheus; prometheus.Enable {
		var (
			port int
			err  error
		)
		if !prometheus.AutoSetPorts {
			port, err = datautil.GetElemByIndex(prometheus.Ports, 0)
			if err != nil {
				return err
			}
		}
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			return fmt.Errorf("prometheus listen %d error %w", port, err)
		}
		defer listener.Close()
		log.ZDebug(ctx, "prometheus start", "addr", listener.Addr())
		go func() {
			err := prommetrics.Start(listener)
			if err == nil {
				err = fmt.Errorf("http done")
			}
			cancel(fmt.Errorf("prometheus %w", err))
		}()
	}

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
		select {
		case <-ctx.Done():
			return
		case val := <-sigs:
			log.ZDebug(ctx, "recv signal", "signal", val.String())
			cancel(fmt.Errorf("signal %s", val.String()))
		}
	}()

	for i := range x.cmds {
		cmd := x.cmds[i]
		if cmd.Block {
			continue
		}
		if err := cmd.Func(ctx); err != nil {
			cancel(fmt.Errorf("server %s exit %w", cmd.Name, err))
			return err
		}
		go func() {
			if cmd.Block {
				cancel(fmt.Errorf("server %s exit", cmd.Name))
			}
		}()
	}

	var wait cmdManger
	for i := range x.cmds {
		cmd := x.cmds[i]
		if !cmd.Block {
			continue
		}
		wait.Start(cmd.Name)
		go func() {
			defer wait.Shutdown(cmd.Name)
			if err := cmd.Func(ctx); err != nil {
				cancel(fmt.Errorf("server %s exit %w", cmd.Name, err))
				return
			}
			cancel(fmt.Errorf("server %s exit", cmd.Name))
		}()
	}
	<-ctx.Done()
	exitCause := context.Cause(ctx)
	log.ZWarn(ctx, "notification of service closure", exitCause)
	done := wait.Wait()
	timeout := time.NewTimer(time.Second * 10)
	defer timeout.Stop()
	for {
		select {
		case <-timeout.C:
			log.ZWarn(ctx, "server exit timeout", nil, "running", wait.Running())
			return exitCause
		case _, ok := <-done:
			if ok {
				log.ZWarn(ctx, "waiting for the service to exit", nil, "running", wait.Running())
			} else {
				log.ZInfo(ctx, "all server exit done")
				return exitCause
			}
		}
	}
}

func putCmd[C any](cmd *cmds, block bool, fn func(ctx context.Context, config *C, client discovery.SvcDiscoveryRegistry, server grpc.ServiceRegistrar) error) {
	name := path.Base(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name())
	if index := strings.Index(name, "."); index >= 0 {
		name = name[:index]
	}
	cmd.add(name, block, func(ctx context.Context) error {
		var conf C
		if err := cmd.parseConf(&conf); err != nil {
			return err
		}
		return fn(ctx, &conf, inprocess.GetSvcDiscoveryRegistry(), inprocess.GetServiceRegistrar())
	})
}

type cmdManger struct {
	lock  sync.Mutex
	done  chan struct{}
	count int
	names map[string]struct{}
}

func (x *cmdManger) Start(name string) {
	x.lock.Lock()
	defer x.lock.Unlock()
	if x.names == nil {
		x.names = make(map[string]struct{})
	}
	if x.done == nil {
		x.done = make(chan struct{}, 1)
	}
	if _, ok := x.names[name]; ok {
		panic(fmt.Errorf("cmd %s already exists", name))
	}
	x.count++
	x.names[name] = struct{}{}
}

func (x *cmdManger) Shutdown(name string) {
	x.lock.Lock()
	defer x.lock.Unlock()
	if _, ok := x.names[name]; !ok {
		panic(fmt.Errorf("cmd %s not exists", name))
	}
	delete(x.names, name)
	x.count--
	if x.count == 0 {
		close(x.done)
	} else {
		select {
		case x.done <- struct{}{}:
		default:
		}
	}
}

func (x *cmdManger) Wait() <-chan struct{} {
	x.lock.Lock()
	defer x.lock.Unlock()
	if x.count == 0 || x.done == nil {
		tmp := make(chan struct{})
		close(tmp)
		return tmp
	}
	return x.done
}

func (x *cmdManger) Running() []string {
	x.lock.Lock()
	defer x.lock.Unlock()
	names := make([]string, 0, len(x.names))
	for name := range x.names {
		names = append(names, name)
	}
	return names
}

type serverConfig struct {
	API         config.API
	Share       config.Share
	RedisConfig config.Redis
	Index       config.Index
}

func startRedisServerRegister(ctx context.Context, cfg *serverConfig, client discovery.SvcDiscoveryRegistry, service grpc.ServiceRegistrar) error {
	apiPort, err := datautil.GetElemByIndex(cfg.API.Api.Ports, int(cfg.Index))
	if err != nil {
		return err
	}
	if apiPort <= 0 {
		return errors.New("standalone api port is 0")
	}
	registerIP, err := network.GetRpcRegisterIP(cfg.API.Api.RegisterIP)
	if err != nil {
		return err
	}
	const validTime = time.Second * 10
	dbb := dbbuild.NewBuilder(nil, &cfg.RedisConfig)
	rdb, err := dbb.Redis(ctx)
	if err != nil {
		return err
	}
	gateway := redis.NewStandaloneGatewayRedis(rdb, validTime)
	selfAddr := net.JoinHostPort(registerIP, strconv.Itoa(apiPort))
	inprocess.SetLocalTarget(selfAddr)
	inprocess.SetBroadcastAddress(cfg.Share.Secret, func(ctx context.Context) ([]string, error) {
		address, err := gateway.GetGatewayAddrs(ctx)
		if err != nil {
			return nil, err
		}
		notSelf := make([]string, 0, len(address))
		for _, addr := range address {
			if addr != selfAddr {
				notSelf = append(notSelf, addr)
			}
		}
		return notSelf, nil
	})
	register := func() {
		ctx, cancel := context.WithTimeout(ctx, validTime/2)
		defer cancel()
		if err := gateway.RegisterGateway(ctx, selfAddr); err != nil {
			log.ZWarn(ctx, "gateway register failed", err, "address", selfAddr)
		}
	}
	timer := time.NewTicker(validTime / 2)
	defer func() {
		timer.Stop()
		ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), time.Second)
		defer cancel()
		if err := gateway.UnregisterGateway(ctx, selfAddr); err != nil {
			log.ZWarn(ctx, "gateway unregister failed", err, "address", selfAddr)
		}
	}()
	register()
	for {
		select {
		case <-timer.C:
			register()
		case <-ctx.Done():
			return context.Cause(ctx)
		}
	}
}
