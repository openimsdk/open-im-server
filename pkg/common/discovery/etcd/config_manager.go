package etcd

import (
	"context"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"syscall"

	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	ShutDowns []func() error
)

func RegisterShutDown(shutDown ...func() error) {
	ShutDowns = append(ShutDowns, shutDown...)
}

type ConfigManager struct {
	client           *clientv3.Client
	watchConfigNames []string
	lock             sync.Mutex
}

func BuildKey(s string) string {
	return ConfigKeyPrefix + s
}

func NewConfigManager(client *clientv3.Client, configNames []string) *ConfigManager {
	return &ConfigManager{
		client:           client,
		watchConfigNames: datautil.Batch(func(s string) string { return BuildKey(s) }, append(configNames, RestartKey))}
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
						c.lock.Lock()
						err := restartServer(ctx)
						if err != nil {
							log.ZError(ctx, "restart server err", err)
						}
						c.lock.Unlock()
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
	log.ZInfo(ctx, "cmd start over")

	os.Exit(0)
	return nil
}
