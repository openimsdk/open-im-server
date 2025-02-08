package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/mitchellh/mapstructure"
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
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/discovery/standalone"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "", "config path")
	flag.Parse()
	if configPath == "" {
		if runtime.GOOS == "linux" {
			configPath = "/root/dt/open-im-server/config"
		} else {
			configPath = "/Users/chao/Desktop/code/open-im-server/config"
		}
	}
	cmd := newCmds(configPath)
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
	ctx := context.Background()
	if err := cmd.run(ctx); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("exit")
}

func newCmds(confPath string) *cmds {
	return &cmds{confPath: confPath}
}

type cmdName struct {
	Name  string
	Func  func(ctx context.Context) error
	Block bool
}
type cmds struct {
	confPath string
	cmds     []cmdName
	config   config.AllConfig
	conf     map[string]reflect.Value
}

func (x *cmds) getTypePath(typ reflect.Type) string {
	return path.Join(typ.PkgPath(), typ.Name())
}

func (x *cmds) initDiscovery() {
	x.config.Discovery.Enable = config.Standalone
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
		name := val.(interface{ GetConfigFileName() string }).GetConfigFileName()
		confData, err := os.ReadFile(filepath.Join(x.confPath, name))
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}
		v := viper.New()
		v.SetConfigType("yaml")
		if err := v.ReadConfig(bytes.NewReader(confData)); err != nil {
			return err
		}
		opt := func(conf *mapstructure.DecoderConfig) {
			conf.TagName = config.StructTagName
		}
		if err := v.Unmarshal(val, opt); err != nil {
			return err
		}
	}
	x.initDiscovery()
	return nil
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

func (x *cmds) run(ctx context.Context) error {
	if len(x.cmds) == 0 {
		return fmt.Errorf("no command to run")
	}
	if err := x.initAllConfig(); err != nil {
		return err
	}
	ctx, cancel := context.WithCancelCause(ctx)
	for i := range x.cmds {
		cmd := x.cmds[i]
		go func() {
			if err := cmd.Func(ctx); err != nil {
				cancel(fmt.Errorf("server %s exit %w", cmd.Name, err))
				return
			}
			if cmd.Block {
				cancel(fmt.Errorf("server %s exit", cmd.Name))
			}
		}()
	}
	<-ctx.Done()
	return context.Cause(ctx)
}

func putCmd[C any](cmd *cmds, block bool, fn func(ctx context.Context, config *C, client discovery.Conn, server grpc.ServiceRegistrar) error) {
	name := path.Base(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name())
	if index := strings.Index(name, "."); index >= 0 {
		name = name[:index]
	}
	cmd.add(name, block, func(ctx context.Context) error {
		var conf C
		if err := cmd.parseConf(&conf); err != nil {
			return err
		}
		return fn(ctx, &conf, standalone.GetDiscoveryConn(), standalone.GetServiceRegistrar())
	})
}
