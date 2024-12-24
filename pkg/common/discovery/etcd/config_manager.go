package etcd

import (
	"context"
	"os"
	"os/exec"
	"runtime"
	"syscall"

	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	ConfigKeyPrefix = "/open-im/config/"
)

var (
	ShutDowns  []func() error
	CanRestart chan struct{}
)

func init() {
	CanRestart = make(chan struct{}, 1)
	CanRestart <- struct{}{}
}

func RegisterShutDown(shutDown ...func() error) {
	ShutDowns = append(ShutDowns, shutDown...)
}

type ConfigManager struct {
	client           *clientv3.Client
	watchConfigNames []string
}

func BuildKey(s string) string {
	return ConfigKeyPrefix + s
}

func NewConfigManager(client *clientv3.Client, configNames []string) *ConfigManager {
	return &ConfigManager{
		client:           client,
		watchConfigNames: datautil.Batch(func(s string) string { return BuildKey(s) }, configNames)}
}

func (c *ConfigManager) Watch(ctx context.Context) {
	chans := make([]clientv3.WatchChan, 0, len(c.watchConfigNames))
	for _, name := range c.watchConfigNames {
		chans = append(chans, c.client.Watch(ctx, name, clientv3.WithPrefix()))
	}

	doWatch := func(watchChan clientv3.WatchChan) {
		for watchResp := range watchChan {
			if watchResp.Err() != nil {
				log.ZError(ctx, "watch err", errs.Wrap(watchResp.Err()))
				continue
			}
			for _, event := range watchResp.Events {
				if event.IsModify() {
					if datautil.Contain(string(event.Kv.Key), c.watchConfigNames...) {
						<-CanRestart
						err := restartServer(ctx)
						if err != nil {
							log.ZError(ctx, "restart server err", err)
							CanRestart <- struct{}{}
						}
					}
				}
			}
		}
	}
	for _, ch := range chans {
		go doWatch(ch)
	}
}

func restartServer(ctx context.Context) error {
	exePath, err := os.Executable()
	if err != nil {
		return errs.New("get executable path fail").Wrap()
	}

	args := os.Args
	env := os.Environ()

	cmd := exec.Command(exePath, args[1:]...)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if runtime.GOOS != "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	log.ZInfo(ctx, "shutdown server")
	for _, f := range ShutDowns {
		if err = f(); err != nil {
			log.ZError(ctx, "shutdown fail", err)
		}
	}

	log.ZInfo(ctx, "restart server")
	err = cmd.Start()
	if err != nil {
		return errs.New("restart server fail").Wrap()
	}

	os.Exit(0)
	return nil
}
