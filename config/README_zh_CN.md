# 						OpenIM配置文件说明以及常用配置修改说明

## 配置文件说明

| Configuration File              | Description                                                  |
| ------------------------------- | ------------------------------------------------------------ |
| **kafka.yml**                   | Kafka用户名、密码、地址等配置                                |
| **redis.yml**                   | Redis密码、地址等配置                                        |
| **minio.yml**                   | MinIO用户名、密码、地址及外网IP域名等配置；未修改外网IP或域名可能导致图片文件发送失败 |
| **zookeeper.yml**               | ZooKeeper用户、密码、地址等配置                              |
| **mongodb.yml**                 | MongoDB用户名、密码、地址等配置                              |
| **log.yml**                     | 日志级别及存储目录等配置                                     |
| **notification.yml**            | 添加好友、创建群组等事件通知配置                             |
| **share.yml**                   | OpenIM各服务所需的公共配置，如secret等                       |
| **webhooks.yml**                | Webhook中URL等配置                                           |
| **local-cache.yml**             | 本地缓存配置                                                 |
| **openim-rpc-third.yml**        | openim-rpc-third服务的监听IP、端口及图片视频对象存储配置     |
| **openim-rpc-user.yml**         | openim-rpc-user服务的监听IP、端口配置                        |
| **openim-api.yml**              | openim-api服务的监听IP、端口等配置项                         |
| **openim-crontask.yml**         | openim-crontask服务配置                                      |
| **openim-msggateway.yml**       | openim-msggateway服务的监听IP、端口等配置                    |
| **openim-msgtransfer.yml**      | openim-msgtransfer服务配置                                   |
| **openim-push.yml**             | openim-push服务的监听IP、端口及离线推送配置                  |
| **openim-rpc-auth.yml**         | openim-rpc-auth服务的监听IP、端口及token有效期等配置         |
| **openim-rpc-conversation.yml** | openim-rpc-conversation服务的监听IP、端口等配置              |
| **openim-rpc-friend.yml**       | openim-rpc-friend服务的监听IP、端口等配置                    |
| **openim-rpc-group.yml**        | openim-rpc-group服务的监听IP、端口等配置                     |
| **openim-rpc-msg.yml**          | openim-rpc-msg服务的监听IP、端口及消息发送是否验证好友关系等配置 |

## 常用配置修改

| 修改配置项                                      | 配置文件                |
| ----------------------------------------------- | ----------------------- |
| 使用minio作为图片视频文件对象存储               | `minio.yml`             |
| 生产环境日志调整                                | `log.yml`               |
| 发送消息是否验证好友关系                        | `openim-rpc-msg.yml`    |
| 修改secret                                      | `share.yml`             |
| 使用oss, cos, aws, kodo作为图片视频文件对象存储 | `openim-rpc-third.yml`  |
| 设置多端互踢策略                                | `openim-msggateway.yml` |
| 设置离线推送                                    | `openim-push.yml`       |

## 启动某个OpenIM服务的多个实例

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



