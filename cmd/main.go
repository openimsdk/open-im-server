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
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "/Users/chao/Desktop/code/open-im-server/config", "config path")
	flag.Parse()
	cmd := newCmds(configPath)
	putCmd1(cmd, auth.Start)
	putCmd1(cmd, conversation.Start)
	putCmd1(cmd, relation.Start)
	putCmd1(cmd, group.Start)
	putCmd1(cmd, msg.Start)
	putCmd1(cmd, third.Start)
	putCmd1(cmd, user.Start)
	putCmd1(cmd, push.Start)
	putCmd2(cmd, msggateway.Start)
	putCmd2(cmd, msgtransfer.Start)
	putCmd2(cmd, api.Start)
	ctx := context.Background()
	if err := cmd.run(ctx); err != nil {
		fmt.Println(err)
	}
	fmt.Println("success")
}

func getTypePath(typ reflect.Type) string {
	return path.Join(typ.PkgPath(), typ.Name())
}

func newCmds(confPath string) *cmds {
	return &cmds{confPath: confPath}
}

type cmds struct {
	confPath string
	cmds     []cmdName
	conf     map[string][]byte
}

func (x *cmds) readConfig() error {
	skip := []string{
		config.DiscoveryConfigFilename,
	}
	if x.conf == nil {
		x.conf = make(map[string][]byte)
	}
	vof := reflect.ValueOf(&config.AllConfig{}).Elem()
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
		itemConf := field.Addr().Interface()
		name := itemConf.(interface{ GetConfigFileName() string }).GetConfigFileName()
		if datautil.Contain(name, skip...) {
			x.conf[getTypePath(field.Type())] = nil
			continue
		}
		data, err := os.ReadFile(filepath.Join(x.confPath, name))
		if err != nil {
			return err
		}
		x.conf[getTypePath(field.Type())] = data
	}
	val := config.Discovery{Enable: config.Standalone}
	var buf bytes.Buffer
	if err := yaml.NewEncoder(&buf).Encode(&val); err != nil {
		return err
	}
	x.conf[getTypePath(reflect.TypeOf(val))] = buf.Bytes()
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
		pkt := getTypePath(field.Type())
		confData, ok := x.conf[pkt]
		if !ok {
			if typeField.Name == "FcmConfigPath" && field.Kind() == reflect.String {
				field.SetString(x.confPath)
				continue
			}
			return fmt.Errorf("config field %s %s not found", vof.Type().Name(), typeField.Name)
		}
		if confData == nil {
			continue
		}
		val := field.Addr().Interface()
		v := viper.New()
		v.SetConfigType("yaml")
		if err := v.ReadConfig(bytes.NewReader(confData)); err != nil {
			return err
		}
		fn := func(conf *mapstructure.DecoderConfig) {
			conf.TagName = config.StructTagName
		}
		if err := v.Unmarshal(val, fn); err != nil {
			return err
		}
	}
	return nil
}

func (x *cmds) add(name string, fn func(ctx context.Context) error) {
	x.cmds = append(x.cmds, cmdName{Name: name, Func: fn})
}

func (x *cmds) run(ctx context.Context) error {
	if x.conf == nil {
		if err := x.readConfig(); err != nil {
			return err
		}
	}
	for _, cmd := range x.cmds {
		fmt.Println("start", cmd.Name)
		if err := cmd.Func(ctx); err != nil {
			fmt.Println("start failed", cmd.Name, err)
			return err
		}
		fmt.Println("start ok", cmd.Name)
	}
	return nil
}

type cmdName struct {
	Name string
	Func func(ctx context.Context) error
}

func getFuncPacketName(fn any) string {
	name := path.Base(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name())
	if index := strings.Index(name, "."); index >= 0 {
		name = name[:index]
	}
	return name
}

func putCmd1[C any](cmd *cmds, fn func(ctx context.Context, config *C, client discovery.Conn, server grpc.ServiceRegistrar) error) {
	cmd.add(getFuncPacketName(fn), func(ctx context.Context) error {
		var conf C
		if err := cmd.parseConf(&conf); err != nil {
			return err
		}
		return fn(ctx, &conf, standalone.GetDiscoveryConn(), standalone.GetServiceRegistrar())
	})
}

func putCmd2[C any](cmd *cmds, fn func(ctx context.Context, index int, config *C) error) {
	cmd.add(getFuncPacketName(fn), func(ctx context.Context) error {
		var conf C
		if err := cmd.parseConf(&conf); err != nil {
			return err
		}
		return fn(ctx, 0, &conf)
	})
}
