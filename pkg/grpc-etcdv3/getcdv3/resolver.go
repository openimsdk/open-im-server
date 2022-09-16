package getcdv3

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"context"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"

	//"go.etcd.io/etcd/mvcc/mvccpb"
	//"google.golang.org/genproto/googleapis/ads/googleads/v1/services"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
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

func NewResolver(schema, etcdAddr, serviceName string, operationID string) (*Resolver, error) {
	etcdCli, err := clientv3.New(clientv3.Config{
		Endpoints: strings.Split(etcdAddr, ","),
		Username:  config.Config.Etcd.UserName,
		Password:  config.Config.Etcd.Password,
	})
	if err != nil {
		log.Error(operationID, "etcd client v3 failed")
		return nil, utils.Wrap(err, "")
	}

	var r Resolver
	r.serviceName = serviceName
	r.cli = etcdCli
	r.schema = schema
	r.etcdAddr = etcdAddr
	resolver.Register(&r)
	//
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	conn, err := grpc.DialContext(ctx, GetPrefix(schema, serviceName),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithInsecure())
	log.Debug(operationID, "etcd key ", GetPrefix(schema, serviceName))
	if err == nil {
		r.grpcClientConn = conn
	}
	return &r, utils.Wrap(err, "")
}

func (r1 *Resolver) ResolveNow(rn resolver.ResolveNowOptions) {
}

func (r1 *Resolver) Close() {
}

func getConn(schema, etcdaddr, serviceName string, operationID string) *grpc.ClientConn {
	rwNameResolverMutex.RLock()
	r, ok := nameResolver[schema+serviceName]
	rwNameResolverMutex.RUnlock()
	if ok {
		log.Debug(operationID, "etcd key ", schema+serviceName, "value ", *r.grpcClientConn, *r)
		return r.grpcClientConn
	}

	rwNameResolverMutex.Lock()
	r, ok = nameResolver[schema+serviceName]
	if ok {
		rwNameResolverMutex.Unlock()
		log.Debug(operationID, "etcd key ", schema+serviceName, "value ", *r.grpcClientConn, *r)
		return r.grpcClientConn
	}

	r, err := NewResolver(schema, etcdaddr, serviceName, operationID)
	if err != nil {
		log.Error(operationID, "etcd failed ", schema, etcdaddr, serviceName, err.Error())
		rwNameResolverMutex.Unlock()
		return nil
	}

	log.Debug(operationID, "etcd key ", schema+serviceName, "value ", *r.grpcClientConn, *r)
	nameResolver[schema+serviceName] = r
	rwNameResolverMutex.Unlock()
	return r.grpcClientConn
}

func GetConfigConn(serviceName string, operationID string) *grpc.ClientConn {
	rpcRegisterIP := config.Config.RpcRegisterIP
	var err error
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error(operationID, "GetLocalIP failed ", err.Error())
			return nil
		}
	}

	var configPortList []int
	//1
	if config.Config.RpcRegisterName.OpenImUserName == serviceName {
		configPortList = config.Config.RpcPort.OpenImUserPort
	}
	//2
	if config.Config.RpcRegisterName.OpenImFriendName == serviceName {
		configPortList = config.Config.RpcPort.OpenImFriendPort
	}
	//3
	if config.Config.RpcRegisterName.OpenImMsgName == serviceName {
		configPortList = config.Config.RpcPort.OpenImMessagePort
	}
	//4
	if config.Config.RpcRegisterName.OpenImPushName == serviceName {
		configPortList = config.Config.RpcPort.OpenImPushPort
	}
	//5
	if config.Config.RpcRegisterName.OpenImRelayName == serviceName {
		configPortList = config.Config.RpcPort.OpenImMessageGatewayPort
	}
	//6
	if config.Config.RpcRegisterName.OpenImGroupName == serviceName {
		configPortList = config.Config.RpcPort.OpenImGroupPort
	}
	//7
	if config.Config.RpcRegisterName.OpenImAuthName == serviceName {
		configPortList = config.Config.RpcPort.OpenImAuthPort
	}
	//10
	if config.Config.RpcRegisterName.OpenImOfficeName == serviceName {
		configPortList = config.Config.RpcPort.OpenImOfficePort
	}
	//11
	if config.Config.RpcRegisterName.OpenImOrganizationName == serviceName {
		configPortList = config.Config.RpcPort.OpenImOrganizationPort
	}
	//12
	if config.Config.RpcRegisterName.OpenImConversationName == serviceName {
		configPortList = config.Config.RpcPort.OpenImConversationPort
	}
	//13
	if config.Config.RpcRegisterName.OpenImCacheName == serviceName {
		configPortList = config.Config.RpcPort.OpenImCachePort
	}
	//14
	if config.Config.RpcRegisterName.OpenImRealTimeCommName == serviceName {
		configPortList = config.Config.RpcPort.OpenImRealTimeCommPort
	}
	if len(configPortList) == 0 {
		log.Error(operationID, "len(configPortList) == 0  ")
		return nil
	}
	target := rpcRegisterIP + ":" + utils.Int32ToString(int32(configPortList[0]))
	log.Info(operationID, "rpcRegisterIP ", rpcRegisterIP, " port ", configPortList, " grpc target: ", target, " serviceName: ", serviceName)
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		log.Error(operationID, "grpc.Dail failed ", err.Error())
		return nil
	}
	log.NewDebug(operationID, utils.GetSelfFuncName(), serviceName, conn)
	return conn
}

func GetDefaultConn(schema, etcdaddr, serviceName string, operationID string) *grpc.ClientConn {
	con := getConn(schema, etcdaddr, serviceName, operationID)
	if con != nil {
		return con
	}
	log.NewWarn(operationID, utils.GetSelfFuncName(), "conn is nil !!!!!", schema, etcdaddr, serviceName, operationID)
	con = GetConfigConn(serviceName, operationID)
	return con
}

func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	if r.cli == nil {
		return nil, fmt.Errorf("etcd clientv3 client failed, etcd:%s", target)
	}
	r.cc = cc
	log.Debug("", "Build..")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//     "%s:///%s"
	prefix := GetPrefix(r.schema, r.serviceName)
	// get key first
	resp, err := r.cli.Get(ctx, prefix, clientv3.WithPrefix())
	if err == nil {
		var addrList []resolver.Address
		for i := range resp.Kvs {
			log.Debug("", "etcd init addr: ", string(resp.Kvs[i].Value))
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
					log.Debug("", "after add, new list: ", addrList)
				}
			case mvccpb.DELETE:
				log.Debug("remove addr key: ", string(ev.Kv.Key), "value:", string(ev.Kv.Value))
				i := strings.LastIndexAny(string(ev.Kv.Key), "/")
				if i < 0 {
					return
				}
				t := string(ev.Kv.Key)[i+1:]
				log.Debug("remove addr key: ", string(ev.Kv.Key), "value:", string(ev.Kv.Value), "addr:", t)
				if s, ok := remove(addrList, t); ok {
					flag = 1
					addrList = s
					log.Debug("after remove, new list: ", addrList)
				}
			}
		}

		if flag == 1 {
			r.cc.UpdateState(resolver.State{Addresses: addrList})
			log.Debug("update: ", addrList)
		}
	}
}

var Conn4UniqueList []*grpc.ClientConn
var Conn4UniqueListMtx sync.RWMutex
var IsUpdateStart bool
var IsUpdateStartMtx sync.RWMutex

func GetDefaultGatewayConn4Unique(schema, etcdaddr, operationID string) []*grpc.ClientConn {
	IsUpdateStartMtx.Lock()
	if IsUpdateStart == false {
		Conn4UniqueList = getConn4Unique(schema, etcdaddr, config.Config.RpcRegisterName.OpenImRelayName)
		go func() {
			for {
				select {
				case <-time.After(time.Second * time.Duration(30)):
					Conn4UniqueListMtx.Lock()
					Conn4UniqueList = getConn4Unique(schema, etcdaddr, config.Config.RpcRegisterName.OpenImRelayName)
					Conn4UniqueListMtx.Unlock()
				}
			}
		}()
	}
	IsUpdateStart = true
	IsUpdateStartMtx.Unlock()

	Conn4UniqueListMtx.Lock()
	var clientConnList []*grpc.ClientConn
	for _, v := range Conn4UniqueList {
		clientConnList = append(clientConnList, v)
	}
	Conn4UniqueListMtx.Unlock()

	//grpcConns := getConn4Unique(schema, etcdaddr, config.Config.RpcRegisterName.OpenImRelayName)
	grpcConns := clientConnList
	if len(grpcConns) > 0 {
		return grpcConns
	}
	log.NewWarn(operationID, utils.GetSelfFuncName(), " len(grpcConns) == 0 ", schema, etcdaddr, config.Config.RpcRegisterName.OpenImRelayName)
	grpcConns = GetDefaultGatewayConn4UniqueFromcfg(operationID)
	log.NewDebug(operationID, utils.GetSelfFuncName(), config.Config.RpcRegisterName.OpenImRelayName, grpcConns)
	return grpcConns
}

func GetDefaultGatewayConn4UniqueFromcfg(operationID string) []*grpc.ClientConn {
	rpcRegisterIP := config.Config.RpcRegisterIP
	var err error
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
			return nil
		}
	}
	var conns []*grpc.ClientConn
	configPortList := config.Config.RpcPort.OpenImMessageGatewayPort
	for _, port := range configPortList {
		target := rpcRegisterIP + ":" + utils.Int32ToString(int32(port))
		log.Info(operationID, "rpcRegisterIP ", rpcRegisterIP, " port ", configPortList, " grpc target: ", target, " serviceName: ", "msgGateway")
		conn, err := grpc.Dial(target, grpc.WithInsecure())
		if err != nil {
			log.Error(operationID, "grpc.Dail failed ", err.Error())
			continue
		}
		conns = append(conns, conn)

	}
	return conns

}

func getConn4Unique(schema, etcdaddr, servicename string) []*grpc.ClientConn {
	gEtcdCli, err := clientv3.New(clientv3.Config{Endpoints: strings.Split(etcdaddr, ",")})
	if err != nil {
		log.Error("clientv3.New failed", err.Error())
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
		r := getConn(schema, etcdaddr, v, "0")
		allConn = append(allConn, r)
	}

	return allConn
}

var (
	service2pool   = make(map[string]*Pool)
	service2poolMu sync.Mutex
)

func GetconnFactory(schema, etcdaddr, servicename string) (*grpc.ClientConn, error) {
	c := getConn(schema, etcdaddr, servicename, "0")
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
