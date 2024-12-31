package api

import (
	"context"
	"encoding/json"
	"reflect"
	"strconv"
	"time"

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

const (
	// wait for Restart http call return
	waitHttp = time.Millisecond * 200
)

type ConfigManager struct {
	imAdminUserID []string
	config        *config.AllConfig
	client        *clientv3.Client

	configPath string
	runtimeEnv string
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
		err = compareAndSave[config.Discovery](c, cm.config.Name2Config(req.ConfigName), &req, cm)
	case cm.config.Kafka.GetConfigFileName():
		err = compareAndSave[config.Kafka](c, cm.config.Name2Config(req.ConfigName), &req, cm)
	case cm.config.LocalCache.GetConfigFileName():
		err = compareAndSave[config.LocalCache](c, cm.config.Name2Config(req.ConfigName), &req, cm)
	case cm.config.Log.GetConfigFileName():
		err = compareAndSave[config.Log](c, cm.config.Name2Config(req.ConfigName), &req, cm)
	case cm.config.Minio.GetConfigFileName():
		err = compareAndSave[config.Minio](c, cm.config.Name2Config(req.ConfigName), &req, cm)
	case cm.config.Mongo.GetConfigFileName():
		err = compareAndSave[config.Mongo](c, cm.config.Name2Config(req.ConfigName), &req, cm)
	case cm.config.Notification.GetConfigFileName():
		err = compareAndSave[config.Notification](c, cm.config.Name2Config(req.ConfigName), &req, cm)
	case cm.config.API.GetConfigFileName():
		err = compareAndSave[config.API](c, cm.config.Name2Config(req.ConfigName), &req, cm)
	case cm.config.CronTask.GetConfigFileName():
		err = compareAndSave[config.CronTask](c, cm.config.Name2Config(req.ConfigName), &req, cm)
	case cm.config.MsgGateway.GetConfigFileName():
		err = compareAndSave[config.MsgGateway](c, cm.config.Name2Config(req.ConfigName), &req, cm)
	case cm.config.MsgTransfer.GetConfigFileName():
		err = compareAndSave[config.MsgTransfer](c, cm.config.Name2Config(req.ConfigName), &req, cm)
	case cm.config.Push.GetConfigFileName():
		err = compareAndSave[config.Push](c, cm.config.Name2Config(req.ConfigName), &req, cm)
	case cm.config.Auth.GetConfigFileName():
		err = compareAndSave[config.Auth](c, cm.config.Name2Config(req.ConfigName), &req, cm)
	case cm.config.Conversation.GetConfigFileName():
		err = compareAndSave[config.Conversation](c, cm.config.Name2Config(req.ConfigName), &req, cm)
	case cm.config.Friend.GetConfigFileName():
		err = compareAndSave[config.Friend](c, cm.config.Name2Config(req.ConfigName), &req, cm)
	case cm.config.Group.GetConfigFileName():
		err = compareAndSave[config.Group](c, cm.config.Name2Config(req.ConfigName), &req, cm)
	case cm.config.Msg.GetConfigFileName():
		err = compareAndSave[config.Msg](c, cm.config.Name2Config(req.ConfigName), &req, cm)
	case cm.config.Third.GetConfigFileName():
		err = compareAndSave[config.Third](c, cm.config.Name2Config(req.ConfigName), &req, cm)
	case cm.config.User.GetConfigFileName():
		err = compareAndSave[config.User](c, cm.config.Name2Config(req.ConfigName), &req, cm)
	case cm.config.Redis.GetConfigFileName():
		err = compareAndSave[config.Redis](c, cm.config.Name2Config(req.ConfigName), &req, cm)
	case cm.config.Share.GetConfigFileName():
		err = compareAndSave[config.Share](c, cm.config.Name2Config(req.ConfigName), &req, cm)
	case cm.config.Webhooks.GetConfigFileName():
		err = compareAndSave[config.Webhooks](c, cm.config.Name2Config(req.ConfigName), &req, cm)
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

func compareAndSave[T any](c *gin.Context, old any, req *apistruct.SetConfigReq, cm *ConfigManager) error {
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
	_, err = cm.client.Put(c, etcd.BuildKey(req.ConfigName), string(data))
	if err != nil {
		return errs.WrapMsg(err, "save to etcd failed")
	}
	return nil
}

func (cm *ConfigManager) ResetConfig(c *gin.Context) {
	go cm.resetConfig(c)
	apiresp.GinSuccess(c, nil)
}

func (cm *ConfigManager) resetConfig(c *gin.Context) {
	txn := cm.client.Txn(c)
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
	}

	ops := make([]clientv3.Op, 0)
	for _, k := range changedKeys {
		data, err := json.Marshal(configMap[k].new)
		if err != nil {
			log.ZError(c, "marshal config failed", err)
			continue
		}
		ops = append(ops, clientv3.OpPut(etcd.BuildKey(k), string(data)))
	}
	if len(ops) > 0 {
		txn.Then(ops...)
		_, err := txn.Commit()
		if err != nil {
			log.ZError(c, "commit etcd txn failed", err)
			return
		}
	}
}

func (cm *ConfigManager) Restart(c *gin.Context) {
	go cm.restart(c)
	apiresp.GinSuccess(c, nil)
}

func (cm *ConfigManager) restart(c *gin.Context) {
	time.Sleep(waitHttp) // wait for Restart http call return
	t := time.Now().Unix()
	_, err := cm.client.Put(c, etcd.BuildKey(etcd.RestartKey), strconv.Itoa(int(t)))
	if err != nil {
		log.ZError(c, "restart etcd put key failed", err)
	}
}

func (cm *ConfigManager) SetEnableConfigManager(c *gin.Context) {
	var req apistruct.SetEnableConfigManagerReq
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap())
		return
	}
	var enableStr string
	if req.Enable {
		enableStr = etcd.Enable
	} else {
		enableStr = etcd.Disable
	}
	resp, err := cm.client.Get(c, etcd.BuildKey(etcd.EnableConfigCenterKey))
	if err != nil {
		apiresp.GinError(c, errs.WrapMsg(err, "getEnableConfigManager failed"))
		return
	}
	if !(resp.Count > 0 && string(resp.Kvs[0].Value) == etcd.Enable) && req.Enable {
		go func() {
			time.Sleep(waitHttp) // wait for Restart http call return
			err := cm.writeAllConfig(c, clientv3.OpPut(etcd.BuildKey(etcd.EnableConfigCenterKey), enableStr))
			if err != nil {
				log.ZError(c, "writeAllConfig failed", err)
			}
		}()
	} else {
		_, err = cm.client.Put(c, etcd.BuildKey(etcd.EnableConfigCenterKey), enableStr)
		if err != nil {
			apiresp.GinError(c, errs.WrapMsg(err, "setEnableConfigManager failed"))
			return
		}
	}

	apiresp.GinSuccess(c, nil)
}

func (cm *ConfigManager) GetEnableConfigManager(c *gin.Context) {
	resp, err := cm.client.Get(c, etcd.BuildKey(etcd.EnableConfigCenterKey))
	if err != nil {
		apiresp.GinError(c, errs.WrapMsg(err, "getEnableConfigManager failed"))
		return
	}
	var enable bool
	if resp.Count > 0 && string(resp.Kvs[0].Value) == etcd.Enable {
		enable = true
	}
	apiresp.GinSuccess(c, &apistruct.GetEnableConfigManagerResp{Enable: enable})
}

func (cm *ConfigManager) writeAllConfig(ctx context.Context, ops ...clientv3.Op) error {
	getWriteConfigOp(ctx, cm.config.Discovery.GetConfigFileName(), cm.config.Discovery, &ops)
	getWriteConfigOp(ctx, cm.config.Kafka.GetConfigFileName(), cm.config.Kafka, &ops)
	getWriteConfigOp(ctx, cm.config.LocalCache.GetConfigFileName(), cm.config.LocalCache, &ops)
	getWriteConfigOp(ctx, cm.config.Log.GetConfigFileName(), cm.config.Log, &ops)
	getWriteConfigOp(ctx, cm.config.Minio.GetConfigFileName(), cm.config.Minio, &ops)
	getWriteConfigOp(ctx, cm.config.Mongo.GetConfigFileName(), cm.config.Mongo, &ops)
	getWriteConfigOp(ctx, cm.config.Notification.GetConfigFileName(), cm.config.Notification, &ops)
	getWriteConfigOp(ctx, cm.config.API.GetConfigFileName(), cm.config.API, &ops)
	getWriteConfigOp(ctx, cm.config.CronTask.GetConfigFileName(), cm.config.CronTask, &ops)
	getWriteConfigOp(ctx, cm.config.MsgGateway.GetConfigFileName(), cm.config.MsgGateway, &ops)
	getWriteConfigOp(ctx, cm.config.MsgTransfer.GetConfigFileName(), cm.config.MsgTransfer, &ops)
	getWriteConfigOp(ctx, cm.config.Push.GetConfigFileName(), cm.config.Push, &ops)
	getWriteConfigOp(ctx, cm.config.Auth.GetConfigFileName(), cm.config.Auth, &ops)
	getWriteConfigOp(ctx, cm.config.Conversation.GetConfigFileName(), cm.config.Conversation, &ops)
	getWriteConfigOp(ctx, cm.config.Friend.GetConfigFileName(), cm.config.Friend, &ops)
	getWriteConfigOp(ctx, cm.config.Group.GetConfigFileName(), cm.config.Group, &ops)
	getWriteConfigOp(ctx, cm.config.Msg.GetConfigFileName(), cm.config.Msg, &ops)
	getWriteConfigOp(ctx, cm.config.Third.GetConfigFileName(), cm.config.Third, &ops)
	getWriteConfigOp(ctx, cm.config.User.GetConfigFileName(), cm.config.User, &ops)
	getWriteConfigOp(ctx, cm.config.Redis.GetConfigFileName(), cm.config.Redis, &ops)
	getWriteConfigOp(ctx, cm.config.Share.GetConfigFileName(), cm.config.Share, &ops)
	getWriteConfigOp(ctx, cm.config.Webhooks.GetConfigFileName(), cm.config.Webhooks, &ops)
	txn := cm.client.Txn(ctx)
	txn.Then(ops...)
	_, err := txn.Commit()
	if err != nil {
		return errs.WrapMsg(err, "writeAllConfig failed commit")
	}
	return nil
}

func getWriteConfigOp[T any](ctx context.Context, key string, config T, ops *[]clientv3.Op) {
	data, err := json.Marshal(config)
	if err != nil {
		log.ZError(ctx, "marshal config failed", err)
		return
	}
	*ops = append(*ops, clientv3.OpPut(key, string(data)))
	return
}
