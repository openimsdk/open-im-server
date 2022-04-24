package getcdv3

import (
	cfg "Open_IM/pkg/common/config"
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"strings"
	"sync"
	"time"
)

type Watcher struct {
	rwLock     sync.RWMutex
	client     *clientv3.Client
	kv         clientv3.KV
	watcher    clientv3.Watcher
	catalog    string
	kvs        map[string]string
	allService []string
	schema     string
	address    []string
}

func NewWatcher() (w *Watcher) {
	var (
		catalog string
		config  clientv3.Config
		client  *clientv3.Client
		kv      clientv3.KV
		watcher clientv3.Watcher
		err     error
	)
	catalog = cfg.Config.Etcd.EtcdSchema + ":///" + cfg.Config.RpcRegisterName.OpenImOnlineMessageRelayName

	config = clientv3.Config{
		Endpoints:   cfg.Config.Etcd.EtcdAddr,               // 集群地址
		DialTimeout: time.Duration(5000) * time.Millisecond, // 连接超时
	}
	// 1、建立连接
	if client, err = clientv3.New(config); err != nil {
		panic(err.Error())
		return
	}
	// 2、得到KV和观察者
	kv = clientv3.NewKV(client)
	watcher = clientv3.NewWatcher(client)

	w = &Watcher{
		client:     client,
		kv:         kv,
		watcher:    watcher,
		catalog:    catalog,
		kvs:        make(map[string]string),
		allService: make([]string, 0),
		schema:     cfg.Config.Etcd.EtcdSchema,
		address:    cfg.Config.Etcd.EtcdAddr,
	}
	return
}

// 监听变化
func (w *Watcher) Run() (err error) {
	var (
		resp               *clientv3.GetResponse
		kvpair             *mvccpb.KeyValue
		watchStartRevision int64
		watchChan          clientv3.WatchChan
		watchResp          clientv3.WatchResponse
		watchEvent         *clientv3.Event
		key                string
		value              string
	)

	// 1、get目录下的所有键值对，并且获知当前集群的revision
	if resp, err = w.kv.Get(context.TODO(), w.catalog, clientv3.WithPrefix()); err != nil {
		return
	}
	for _, kvpair = range resp.Kvs {
		key = string(kvpair.Key)
		value = string(kvpair.Value)
		w.kvs[key] = value
	}
	w.updateServices()

	// 2、从该revision向后监听变化事件
	go func() {
		// 从GET时刻的后续版本开始监听变化
		watchStartRevision = resp.Header.Revision + 1
		// 监听目录的后续变化
		watchChan = w.watcher.Watch(context.TODO(), w.catalog, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix())
		// 处理监听事件
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT: // 任务保存事件
					w.rwLock.Lock()

					key = string(watchEvent.Kv.Key)
					value = string(watchEvent.Kv.Value)
					w.kvs[key] = value
					w.updateServices()

					w.rwLock.Unlock()
				case mvccpb.DELETE: // 任务被删除了
					w.rwLock.Lock()

					key = string(watchEvent.Kv.Key)
					delete(w.kvs, key)
					w.updateServices()

					w.rwLock.Unlock()
				}
			}
		}
	}()
	return
}

func (w *Watcher) updateServices() {
	var (
		maps        map[string]string
		key         string
		serviceName string
	)
	w.allService = make([]string, 0)
	maps = make(map[string]string)
	for key, _ = range w.kvs {
		serviceName = getServiceName(key)
		if _, ok := maps[serviceName]; ok == true {
			continue
		}
		maps[serviceName] = serviceName
		w.allService = append(w.allService, serviceName)
	}
}

func getServiceName(key string) (name string) {
	var (
		index int
		str   string
	)
	index = strings.LastIndex(key, "///")
	str = key[index+len("///"):]
	index = strings.Index(str, "/")
	name = str[:index]
	return
}

func (w *Watcher) GetAllConns() (conns []*grpc.ClientConn) {
	var (
		services   []string
		service    string
		clientConn *grpc.ClientConn
	)
	w.rwLock.RLock()
	services = w.allService
	w.rwLock.RUnlock()

	conns = make([]*grpc.ClientConn, 0)
	for _, service = range services {
		clientConn = GetConn(w.schema, strings.Join(w.address, ","), service)
		if clientConn == nil {
			continue
		}
		conns = append(conns, clientConn)
	}
	return
}
