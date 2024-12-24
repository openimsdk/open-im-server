package api

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/openimsdk/open-im-server/v3/pkg/apistruct"
	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/discovery/etcd"
	"github.com/openimsdk/open-im-server/v3/version"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/runtimeenv"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type ConfigManager struct {
	imAdminUserID []string
	config        *config.AllConfig
	client        *clientv3.Client
	configPath    string
	runtimeEnv    string
}

func NewConfigManager(IMAdminUserID []string, cfg *config.AllConfig, client *clientv3.Client, configPath string, runtimeEnv string) *ConfigManager {
	return &ConfigManager{
		imAdminUserID: IMAdminUserID,
		config:        cfg,
		client:        client,
		configPath:    configPath,
		runtimeEnv:    runtimeEnv,
	}
}

func (cm *ConfigManager) CheckAdmin(c *gin.Context) {
	if err := authverify.CheckAdmin(c, cm.imAdminUserID); err != nil {
		apiresp.GinError(c, err)
		c.Abort()
	}
}

func (cm *ConfigManager) GetConfig(c *gin.Context) {
	var req apistruct.GetConfigReq
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap())
		return
	}
	conf := cm.config.Name2Config(req.ConfigName)
	if conf == nil {
		apiresp.GinError(c, errs.ErrArgs.WithDetail("config name not found").Wrap())
		return
	}
	b, err := json.Marshal(conf)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, string(b))
}

func (cm *ConfigManager) GetConfigList(c *gin.Context) {
	var resp apistruct.GetConfigListResp
	resp.ConfigNames = cm.config.GetConfigNames()
	resp.Environment = runtimeenv.PrintRuntimeEnvironment()
	resp.Version = version.Version

	apiresp.GinSuccess(c, resp)
}

func (cm *ConfigManager) SetConfig(c *gin.Context) {
	if cm.config.Discovery.Enable != config.ETCD {
		apiresp.GinError(c, errs.New("only etcd support set config").Wrap())
		return
	}
	var req apistruct.SetConfigReq
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap())
		return
	}
	var err error
	switch req.ConfigName {
	case cm.config.Discovery.GetConfigFileName():
		err = compareAndSave[config.Discovery](c, cm.config.Name2Config(req.ConfigName), &req, cm.client)
	case cm.config.Kafka.GetConfigFileName():
		err = compareAndSave[config.Kafka](c, cm.config.Name2Config(req.ConfigName), &req, cm.client)
	case cm.config.LocalCache.GetConfigFileName():
		err = compareAndSave[config.LocalCache](c, cm.config.Name2Config(req.ConfigName), &req, cm.client)
	case cm.config.Log.GetConfigFileName():
		err = compareAndSave[config.Log](c, cm.config.Name2Config(req.ConfigName), &req, cm.client)
	case cm.config.Minio.GetConfigFileName():
		err = compareAndSave[config.Minio](c, cm.config.Name2Config(req.ConfigName), &req, cm.client)
	case cm.config.Mongo.GetConfigFileName():
		err = compareAndSave[config.Mongo](c, cm.config.Name2Config(req.ConfigName), &req, cm.client)
	case cm.config.Notification.GetConfigFileName():
		err = compareAndSave[config.NotificationConfig](c, cm.config.Name2Config(req.ConfigName), &req, cm.client)
	case cm.config.API.GetConfigFileName():
		err = compareAndSave[config.API](c, cm.config.Name2Config(req.ConfigName), &req, cm.client)
	case cm.config.CronTask.GetConfigFileName():
		err = compareAndSave[config.CronTask](c, cm.config.Name2Config(req.ConfigName), &req, cm.client)
	case cm.config.MsgGateway.GetConfigFileName():
		err = compareAndSave[config.MsgGateway](c, cm.config.Name2Config(req.ConfigName), &req, cm.client)
	case cm.config.MsgTransfer.GetConfigFileName():
		err = compareAndSave[config.MsgTransfer](c, cm.config.Name2Config(req.ConfigName), &req, cm.client)
	case cm.config.Push.GetConfigFileName():
		err = compareAndSave[config.Push](c, cm.config.Name2Config(req.ConfigName), &req, cm.client)
	case cm.config.Auth.GetConfigFileName():
		err = compareAndSave[config.Auth](c, cm.config.Name2Config(req.ConfigName), &req, cm.client)
	case cm.config.Conversation.GetConfigFileName():
		err = compareAndSave[config.Conversation](c, cm.config.Name2Config(req.ConfigName), &req, cm.client)
	case cm.config.Friend.GetConfigFileName():
		err = compareAndSave[config.Friend](c, cm.config.Name2Config(req.ConfigName), &req, cm.client)
	case cm.config.Group.GetConfigFileName():
		err = compareAndSave[config.Group](c, cm.config.Name2Config(req.ConfigName), &req, cm.client)
	case cm.config.Msg.GetConfigFileName():
		err = compareAndSave[config.Msg](c, cm.config.Name2Config(req.ConfigName), &req, cm.client)
	case cm.config.Third.GetConfigFileName():
		err = compareAndSave[config.Third](c, cm.config.Name2Config(req.ConfigName), &req, cm.client)
	case cm.config.User.GetConfigFileName():
		err = compareAndSave[config.User](c, cm.config.Name2Config(req.ConfigName), &req, cm.client)
	case cm.config.Redis.GetConfigFileName():
		err = compareAndSave[config.Redis](c, cm.config.Name2Config(req.ConfigName), &req, cm.client)
	case cm.config.Share.GetConfigFileName():
		err = compareAndSave[config.Share](c, cm.config.Name2Config(req.ConfigName), &req, cm.client)
	case cm.config.Webhooks.GetConfigFileName():
		err = compareAndSave[config.Webhooks](c, cm.config.Name2Config(req.ConfigName), &req, cm.client)
	default:
		apiresp.GinError(c, errs.ErrArgs.Wrap())
		return
	}
	if err != nil {
		apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap())
		return
	}
	apiresp.GinSuccess(c, nil)
}

func compareAndSave[T any](c *gin.Context, old any, req *apistruct.SetConfigReq, client *clientv3.Client) error {
	conf := new(T)
	err := json.Unmarshal([]byte(req.Data), &conf)
	if err != nil {
		return errs.ErrArgs.WithDetail(err.Error()).Wrap()
	}
	eq := reflect.DeepEqual(old, conf)
	if eq {
		return nil
	}
	data, err := json.Marshal(conf)
	if err != nil {
		return errs.ErrArgs.WithDetail(err.Error()).Wrap()
	}
	_, err = client.Put(c, etcd.BuildKey(req.ConfigName), string(data))
	if err != nil {
		return errs.WrapMsg(err, "save to etcd failed")
	}
	return nil
}

func (cm *ConfigManager) ResetConfig(c *gin.Context) {
	go cm.restart(c)
	apiresp.GinSuccess(c, nil)
}

func (cm *ConfigManager) restart(c *gin.Context) {
	<-etcd.CanRestart
	type initConf struct {
		old       any
		new       any
		isChanged bool
	}
	configMap := map[string]*initConf{
		cm.config.Discovery.GetConfigFileName():    {old: &cm.config.Discovery, new: new(config.Discovery)},
		cm.config.Kafka.GetConfigFileName():        {old: &cm.config.Kafka, new: new(config.Kafka)},
		cm.config.LocalCache.GetConfigFileName():   {old: &cm.config.LocalCache, new: new(config.LocalCache)},
		cm.config.Log.GetConfigFileName():          {old: &cm.config.Log, new: new(config.Log)},
		cm.config.Minio.GetConfigFileName():        {old: &cm.config.Minio, new: new(config.Minio)},
		cm.config.Mongo.GetConfigFileName():        {old: &cm.config.Mongo, new: new(config.Mongo)},
		cm.config.Notification.GetConfigFileName(): {old: &cm.config.Notification, new: new(config.Notification)},
		cm.config.API.GetConfigFileName():          {old: &cm.config.API, new: new(config.API)},
		cm.config.CronTask.GetConfigFileName():     {old: &cm.config.CronTask, new: new(config.CronTask)},
		cm.config.MsgGateway.GetConfigFileName():   {old: &cm.config.MsgGateway, new: new(config.MsgGateway)},
		cm.config.MsgTransfer.GetConfigFileName():  {old: &cm.config.MsgTransfer, new: new(config.MsgTransfer)},
		cm.config.Push.GetConfigFileName():         {old: &cm.config.Push, new: new(config.Push)},
		cm.config.Auth.GetConfigFileName():         {old: &cm.config.Auth, new: new(config.Auth)},
		cm.config.Conversation.GetConfigFileName(): {old: &cm.config.Conversation, new: new(config.Conversation)},
		cm.config.Friend.GetConfigFileName():       {old: &cm.config.Friend, new: new(config.Friend)},
		cm.config.Group.GetConfigFileName():        {old: &cm.config.Group, new: new(config.Group)},
		cm.config.Msg.GetConfigFileName():          {old: &cm.config.Msg, new: new(config.Msg)},
		cm.config.Third.GetConfigFileName():        {old: &cm.config.Third, new: new(config.Third)},
		cm.config.User.GetConfigFileName():         {old: &cm.config.User, new: new(config.User)},
		cm.config.Redis.GetConfigFileName():        {old: &cm.config.Redis, new: new(config.Redis)},
		cm.config.Share.GetConfigFileName():        {old: &cm.config.Share, new: new(config.Share)},
		cm.config.Webhooks.GetConfigFileName():     {old: &cm.config.Webhooks, new: new(config.Webhooks)},
	}

	changedKeys := make([]string, 0, len(configMap))
	for k, v := range configMap {
		err := config.Load(
			cm.configPath,
			k,
			config.EnvPrefixMap[k],
			cm.runtimeEnv,
			v.new,
		)
		if err != nil {
			log.ZError(c, "load config failed", err)
			continue
		}
		v.isChanged = reflect.DeepEqual(v.old, v.new)
		if !v.isChanged {
			changedKeys = append(changedKeys, k)
		}

		if k == cm.config.Push.GetConfigFileName() {
			fmt.Println("JJJJJJJJJJJJJJJJJJJ")
			fmt.Println(v.old)
			fmt.Println(v.new)
			fmt.Println(v.isChanged)
			fmt.Println("SSSSSSSSSSSSSSSSSSSSS")
		}
	}

	for _, k := range changedKeys {
		data, err := json.Marshal(configMap[k].new)
		if err != nil {
			log.ZError(c, "marshal config failed", err)
			continue
		}
		_, err = cm.client.Put(c, etcd.BuildKey(k), string(data))
		if err != nil {
			log.ZError(c, "save to etcd failed", err)
			continue
		}
	}
	etcd.CanRestart <- struct{}{}
}
