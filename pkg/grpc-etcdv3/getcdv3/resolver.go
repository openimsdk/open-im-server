package getcdv3

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	//"google.golang.org/genproto/googleapis/ads/googleads/v1/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
	"strings"
	"sync"
	"time"
)

type Resolver struct {
	cc                 resolver.ClientConn
	serviceName        string
	grpcClientConn     *grpc.ClientConn
	cli                *clientv3.Client
	schema             string
	etcdAddr           string
	watchStartRevision int64
}

var (
	nameResolver        = make(map[string]*Resolver)
	rwNameResolverMutex sync.RWMutex
)

func NewResolver(schema, etcdAddr, serviceName string) (*Resolver, error) {
	etcdCli, err := clientv3.New(clientv3.Config{
		Endpoints: strings.Split(etcdAddr, ","),
	})
	if err != nil {
		return nil, err
	}

	var r Resolver
	r.serviceName = serviceName
	r.cli = etcdCli
	r.schema = schema
	r.etcdAddr = etcdAddr
	resolver.Register(&r)

	conn, err := grpc.Dial(
		GetPrefix(schema, serviceName),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithInsecure(),
		grpc.WithTimeout(time.Duration(5)*time.Second),
	)
	if err == nil {
		r.grpcClientConn = conn
	}
	return &r, err
}

func (r1 *Resolver) ResolveNow(rn resolver.ResolveNowOptions) {
}

func (r1 *Resolver) Close() {
}

func GetConn(schema, etcdaddr, serviceName string) *grpc.ClientConn {
	rwNameResolverMutex.RLock()
	r, ok := nameResolver[schema+serviceName]
	rwNameResolverMutex.RUnlock()
	if ok {
		return r.grpcClientConn
	}

	rwNameResolverMutex.Lock()
	r, ok = nameResolver[schema+serviceName]
	if ok {
		rwNameResolverMutex.Unlock()
		return r.grpcClientConn
	}

	r, err := NewResolver(schema, etcdaddr, serviceName)
	if err != nil {
		rwNameResolverMutex.Unlock()
		return nil
	}

	nameResolver[schema+serviceName] = r
	rwNameResolverMutex.Unlock()
	return r.grpcClientConn
}

func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	if r.cli == nil {
		return nil, fmt.Errorf("etcd clientv3 client failed, etcd:%s", target)
	}
	r.cc = cc

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//     "%s:///%s"
	prefix := GetPrefix(r.schema, r.serviceName)
	// get key first
	resp, err := r.cli.Get(ctx, prefix, clientv3.WithPrefix())
	if err == nil {
		var addrList []resolver.Address
		for i := range resp.Kvs {
			//log.Debug("", "init addr: ", string(resp.Kvs[i].Value))
			addrList = append(addrList, resolver.Address{Addr: string(resp.Kvs[i].Value)})
		}
		r.cc.UpdateState(resolver.State{Addresses: addrList})
		r.watchStartRevision = resp.Header.Revision + 1
		go r.watch(prefix, addrList)
	} else {
		return nil, fmt.Errorf("etcd get failed, prefix: %s", prefix)
	}

	return r, nil
}

func (r *Resolver) Scheme() string {
	return r.schema
}

func exists(addrList []resolver.Address, addr string) bool {
	for _, v := range addrList {
		if v.Addr == addr {
			return true
		}
	}
	return false
}

func remove(s []resolver.Address, addr string) ([]resolver.Address, bool) {
	for i := range s {
		if s[i].Addr == addr {
			s[i] = s[len(s)-1]
			return s[:len(s)-1], true
		}
	}
	return nil, false
}

func (r *Resolver) watch(prefix string, addrList []resolver.Address) {
	rch := r.cli.Watch(context.Background(), prefix, clientv3.WithPrefix(), clientv3.WithPrefix())
	for n := range rch {
		flag := 0
		for _, ev := range n.Events {
			switch ev.Type {
			case mvccpb.PUT:
				if !exists(addrList, string(ev.Kv.Value)) {
					flag = 1
					addrList = append(addrList, resolver.Address{Addr: string(ev.Kv.Value)})
					//log.Debug("", "after add, new list: ", addrList)
				}
			case mvccpb.DELETE:
				//log.Debug("remove addr key: ", string(ev.Kv.Key), "value:", string(ev.Kv.Value))
				i := strings.LastIndexAny(string(ev.Kv.Key), "/")
				if i < 0 {
					return
				}
				t := string(ev.Kv.Key)[i+1:]
				//log.Debug("remove addr key: ", string(ev.Kv.Key), "value:", string(ev.Kv.Value), "addr:", t)
				if s, ok := remove(addrList, t); ok {
					flag = 1
					addrList = s
					//log.Debug("after remove, new list: ", addrList)
				}
			}
		}

		if flag == 1 {
			r.cc.UpdateState(resolver.State{Addresses: addrList})
			//log.Debug("update: ", addrList)
		}
	}
}

func GetConn4Unique(schema, etcdaddr, servicename string) []*grpc.ClientConn {
	gEtcdCli, err := clientv3.New(clientv3.Config{Endpoints: strings.Split(etcdaddr, ",")})
	if err != nil {
		//log.Error("clientv3.New failed", err.Error())
		return nil
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//     "%s:///%s"
	prefix := GetPrefix4Unique(schema, servicename)

	resp, err := gEtcdCli.Get(ctx, prefix, clientv3.WithPrefix())
	//  "%s:///%s:ip:port"   -> %s:ip:port
	allService := make([]string, 0)
	if err == nil {
		for i := range resp.Kvs {
			k := string(resp.Kvs[i].Key)

			b := strings.LastIndex(k, "///")
			k1 := k[b+len("///"):]

			e := strings.Index(k1, "/")
			k2 := k1[:e]
			allService = append(allService, k2)
		}
	} else {
		gEtcdCli.Close()
		//log.Error("gEtcdCli.Get failed", err.Error())
		return nil
	}
	gEtcdCli.Close()

	allConn := make([]*grpc.ClientConn, 0)
	for _, v := range allService {
		r := GetConn(schema, etcdaddr, v)
		allConn = append(allConn, r)
	}

	return allConn
}

var (
	service2pool   = make(map[string]*Pool)
	service2poolMu sync.Mutex
)

func GetconnFactory(schema, etcdaddr, servicename string) (*grpc.ClientConn, error) {
	c := GetConn(schema, etcdaddr, servicename)
	if c != nil {
		return c, nil
	} else {
		return c, fmt.Errorf("GetConn failed")
	}
}

func GetConnPool(schema, etcdaddr, servicename string) (*ClientConn, error) {
	//get pool
	p := NewPool(schema, etcdaddr, servicename)
	//poo->get

	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(1000*time.Millisecond))

	c, err := p.Get(ctx)
	//log.Info("", "Get ", err)
	return c, err

}

func NewPool(schema, etcdaddr, servicename string) *Pool {

	if _, ok := service2pool[schema+servicename]; !ok {
		//
		service2poolMu.Lock()
		if _, ok1 := service2pool[schema+servicename]; !ok1 {
			p, err := New(GetconnFactory, schema, etcdaddr, servicename, 5, 10, 1)
			if err == nil {
				service2pool[schema+servicename] = p
			}
		}
		service2poolMu.Unlock()
	}

	return service2pool[schema+servicename]
}
func GetGrpcConn(schema, etcdaddr, servicename string) *grpc.ClientConn {
	return nameResolver[schema+servicename].grpcClientConn
}
