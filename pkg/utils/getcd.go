package utils

import (
	"Open_IM/pkg/common/config"
	"context"
	"fmt"
	"github.com/OpenIMSDK/getcdv3"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"

	"gopkg.in/yaml.v3"
	"strings"
)

func registerConf(key, conf string) {
	etcdAddr := strings.Join(config.Config.Etcd.EtcdAddr, ",")
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: strings.Split(etcdAddr, ","), DialTimeout: 5 * time.Second})

	if err != nil {
		panic(err.Error())
	}
	//lease
	if _, err := cli.Put(context.Background(), key, conf); err != nil {
		fmt.Println("panic, params: ")
		panic(err.Error())
	}
}

func RegisterConf() {
	bytes, err := yaml.Marshal(config.Config)
	if err != nil {
		panic(err.Error())
	}
	secretMD5 := Md5(config.Config.Etcd.Secret)
	confBytes, err := AesEncrypt(bytes, []byte(secretMD5[0:16]))
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("start register", secretMD5, getcdv3.GetPrefix(config.Config.Etcd.EtcdSchema, config.ConfName))
	registerConf(getcdv3.GetPrefix(config.Config.Etcd.EtcdSchema, config.ConfName), string(confBytes))
	fmt.Println("etcd register conf ok")
}
