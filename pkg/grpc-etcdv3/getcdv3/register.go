package getcdv3

import (
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"net"
	"strconv"
	"strings"
	"time"
)

type RegEtcd struct {
	cli    *clientv3.Client
	ctx    context.Context
	cancel context.CancelFunc
	key    string
}

var rEtcd *RegEtcd

// "%s:///%s/"
func GetPrefix(schema, serviceName string) string {
	return fmt.Sprintf("%s:///%s/", schema, serviceName)
}

// "%s:///%s"
func GetPrefix4Unique(schema, serviceName string) string {
	return fmt.Sprintf("%s:///%s", schema, serviceName)
}

// "%s:///%s/" ->  "%s:///%s:ip:port"
func RegisterEtcd4Unique(schema, etcdAddr, myHost string, myPort int, serviceName string, ttl int) error {
	serviceName = serviceName + ":" + net.JoinHostPort(myHost, strconv.Itoa(myPort))
	return RegisterEtcd(schema, etcdAddr, myHost, myPort, serviceName, ttl)
}

func GetTarget(schema, myHost string, myPort int, serviceName string) string {
	serviceName = serviceName + ":" + net.JoinHostPort(myHost, strconv.Itoa(myPort))
	return serviceName
}

//etcdAddr separated by commas
func RegisterEtcd(schema, etcdAddr, myHost string, myPort int, serviceName string, ttl int) error {
	operationID := utils.OperationIDGenerator()
	args := schema + " " + etcdAddr + " " + myHost + " " + serviceName + " " + utils.Int32ToString(int32(myPort))
	ttl = ttl * 3
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: strings.Split(etcdAddr, ","), DialTimeout: 5 * time.Second})

	log.Info(operationID, "RegisterEtcd args: ", args, ttl)
	if err != nil {
		log.Error(operationID, "clientv3.New failed ", args, ttl, err.Error())
		return fmt.Errorf("create etcd clientv3 client failed, errmsg:%v, etcd addr:%s", err, etcdAddr)
	}

	//lease
	ctx, cancel := context.WithCancel(context.Background())
	resp, err := cli.Grant(ctx, int64(ttl))
	if err != nil {
		log.Error(operationID, "Grant failed ", err.Error(), ctx, ttl)
		return fmt.Errorf("grant failed")
	}
	log.Info(operationID, "Grant ok, resp ID ", resp.ID)

	//  schema:///serviceName/ip:port ->ip:port
	serviceValue := net.JoinHostPort(myHost, strconv.Itoa(myPort))
	serviceKey := GetPrefix(schema, serviceName) + serviceValue

	//set key->value
	if _, err := cli.Put(ctx, serviceKey, serviceValue, clientv3.WithLease(resp.ID)); err != nil {
		log.Error(operationID, "cli.Put failed ", err.Error(), ctx, args, resp.ID)
		return fmt.Errorf("put failed, errmsg:%v， key:%s, value:%s", err, serviceKey, serviceValue)
	}

	//keepalive
	kresp, err := cli.KeepAlive(ctx, resp.ID)
	if err != nil {
		log.Error(operationID, "KeepAlive failed ", err.Error(), args, resp.ID)
		return fmt.Errorf("keepalive failed, errmsg:%v, lease id:%d", err, resp.ID)
	}
	log.Info(operationID, "RegisterEtcd ok ", args)

	go func() {
		for {
			select {
			case pv, ok := <-kresp:
				if ok == true {
					log.Debug(operationID, "KeepAlive kresp ok", pv, args)
				} else {
					log.Error(operationID, "KeepAlive kresp failed ", pv, args)
					t := time.NewTicker(time.Duration(ttl/2) * time.Second)
					for {
						select {
						case <-t.C:
						}
						ctx, _ := context.WithCancel(context.Background())
						resp, err := cli.Grant(ctx, int64(ttl))
						if err != nil {
							log.Error(operationID, "Grant failed ", err.Error(), args)
							continue
						}

						if _, err := cli.Put(ctx, serviceKey, serviceValue, clientv3.WithLease(resp.ID)); err != nil {
							log.Error(operationID, "etcd Put failed ", err.Error(), args, " resp ID: ", resp.ID)
							continue
						} else {
							log.Info(operationID, "etcd Put ok ", args, " resp ID: ", resp.ID)
						}
					}
				}
			}
		}
	}()

	rEtcd = &RegEtcd{ctx: ctx,
		cli:    cli,
		cancel: cancel,
		key:    serviceKey}

	return nil
}

func UnRegisterEtcd() {
	//delete
	rEtcd.cancel()
	rEtcd.cli.Delete(rEtcd.ctx, rEtcd.key)
}
