# 						OpenIM配置文件说明以及常用配置修改说明

## 外部组件相关配置

| Configuration File | Description                        |
| ------------------ | ---------------------------------- |
| **kafka.yml**      | Kafka用户名、密码、地址等配置      |
| **redis.yml**      | Redis密码、地址等配置              |
| **minio.yml**      | MinIO用户名、密码、地址等配置      |
| **mongodb.yml**    | MongoDB用户名、密码、地址等配置    |
| **discovery.yml**  | 服务发现以及etcd用户名、密码、地址 |

## OpenIMServer相关配置
| Configuration File              | Description                                    |
| ------------------------------- | ---------------------------------------------- |
| **log.yml**                     | 日志级别及存储目录等配置                       |
| **notification.yml**            | 添加好友、创建群组等事件通知配置               |
| **share.yml**                   | 各服务所需的公共配置，如secret等               |
| **webhooks.yml**                | Webhook中URL等配置                             |
| **local-cache.yml**             | 本地缓存配置，一般不用修改                     |
| **openim-rpc-third.yml**        | openim-rpc-third监听IP、端口及对象存储配置     |
| **openim-rpc-user.yml**         | openim-rpc-user监听IP、端口配置                |
| **openim-api.yml**              | openim-api监听IP、端口等配置                   |
| **openim-crontask.yml**         | openim-crontask定时任务配置                    |
| **openim-msggateway.yml**       | openim-msggateway监听IP、端口等配置            |
| **openim-msgtransfer.yml**      | openim-msgtransfer服务配置                     |
| **openim-push.yml**             | openim-push监听IP、端口及离线推送配置          |
| **openim-rpc-auth.yml**         | openim-rpc-auth监听IP、端口及token有效期等配置 |
| **openim-rpc-conversation.yml** | openim-rpc-conversation监听IP、端口等配置      |
| **openim-rpc-friend.yml**       | openim-rpc-friend监听IP、端口等配置            |
| **openim-rpc-group.yml**        | openim-rpc-group监听IP、端口等配置             |
| **openim-rpc-msg.yml**          | openim-rpc-msg服务的监听IP、端口等配置         |


## 监控告警相关配置
| Configuration File             | Description     |
| ------------------------------ | --------------- |
| **prometheus.yml**             | prometheus配置  |
| **instance-down-rules.yml**    | 告警规则        |
| **alertmanager.yml**           | 告警管理配置    |
| **email.tmpl**                 | 邮件告警模版    |
| **grefana-template/Demo.json** | 默认的dashboard |

## 常用配置修改
| 修改配置项                                               | 配置文件                |
| -------------------------------------------------------- | ----------------------- |
| 使用minio作为对象存储时配置，重点关注externalAddress字段 | `minio.yml`             |
| 日志级别及日志文件数量调整                               | `log.yml`               |
| 发送消息是否需要验证好友关系                             | `openim-rpc-msg.yml`    |
| OpenIMServer秘钥                                         | `share.yml`             |
| 使用oss, cos, aws, kodo作为对象存储时配置                | `openim-rpc-third.yml`  |
| 多端互踢策略，单个gateway同时最大连接数                  | `openim-msggateway.yml` |
| 消息离线推送                                             | `openim-push.yml`       |
| 配置webhook来通知回调服务器，如消息发送前后回调          | `webhooks.yml`          |
| 新入群用户是否可以查看历史消息                           | `openim-rpc-group.yml`  |
| token 过期时间设置                                       | `openim-rpc-auth.yml`     |
| 定时任务设置，例如消息保存多长时间                       | `openim-crontask.yml`   |

## 启动某个服务的多个实例和最大文件句柄数


若要启动某个OpenIM的多个实例，只需增加对应的端口数，并修改项目根目录下的`start-config.yml`文件，重启服务即可生效。例如，启动2个`openim-rpc-user`实例的配置如下：

```yaml
rpc:
  registerIP: ''
  listenIP: 0.0.0.0
  ports: [ 10110, 10111 ]

prometheus:
  enable: true
  ports: [ 20100, 20101 ]
```

修改`start-config.yml`:

```yaml
serviceBinaries:
  openim-rpc-user: 2
```

修改最大同时打开的文件句柄数，一般是每个在线用户占用一个

```
maxFileDescriptors: 10000
```
