# 						OpenIM Configuration File Descriptions and Common Configuration Modifications

## External Component Configurations

| Configuration File | Description                                                 |
| ------------------ |-------------------------------------------------------------|
| **kafka.yml**      | Configuration for Kafka username, password, address, etc.   |
| **redis.yml**      | Configuration for Redis password, address, etc.             |
| **minio.yml**      | Configuration for MinIO username, password, address, etc.   |
| **mongodb.yml**    | Configuration for MongoDB username, password, address, etc. |
| **discovery.yml**  | Service discovery and etcd credentials and address.         |

## OpenIMServer Related Configurations
| Configuration File              | Description                                    |
| ------------------------------- | ---------------------------------------------- |
| **log.yml**                     | Configuration for logging levels and storage directory                   |
| **notification.yml**            | Event notification settings (e.g., add friend, create group)           |
| **share.yml**                   | Common settings for all services (e.g., secrets)            |
| **webhooks.yml**                | Webhook URLs and related settings                           |
| **local-cache.yml**             | Local cache settings (generally do not modify)                 |
| **openim-rpc-third.yml**        | openim-rpc-third listen IP, port, and object storage settings  |
| **openim-rpc-user.yml**         | openim-rpc-user listen IP and port settings              |
| **openim-api.yml**              | openim-api listen IP, port, and other settings               |
| **openim-crontask.yml**         | openim-crontask scheduled task settings                   |
| **openim-msggateway.yml**       | openim-msggateway listen IP, port, and other settings           |
| **openim-msgtransfer.yml**      | Settings for openim-msgtransfer service                   |
| **openim-push.yml**             | openim-push listen IP, port, and offline push settings        |
| **openim-rpc-auth.yml**         | openim-rpc-auth listen IP, port, token validity settings |
| **openim-rpc-conversation.yml** | openim-rpc-conversation listen IP and port settings     |
| **openim-rpc-friend.yml**       | openim-rpc-friend listen IP and port settings           |
| **openim-rpc-group.yml**        | openim-rpc-group listen IP and port settings           |
| **openim-rpc-msg.yml**          | openim-rpc-msg listen IP and port settings         |


## Monitoring and Alerting Related Configurations
| Configuration File             | Description     |
| ------------------------------ | --------------- |
| **prometheus.yml**             | Prometheus configuration |
| **instance-down-rules.yml**    | Alert rules       |
| **alertmanager.yml**           | Alertmanager configuration   |
| **email.tmpl**                 | Email alert template   |
| **grefana-template/Demo.json** | Default Grafana dashboard |

## Common Configuration Modifications
| Configuration Item                                              | Configuration File                |
| -------------------------------------------------------- | ----------------------- |
| Configure MinIO as object storage (focus on the externalAddress field) | `minio.yml`             |
| Adjust log level and number of log files                              | `log.yml`               |
| Enable or disable friend verification when sending messages                           | `openim-rpc-msg.yml`    |
| OpenIMServer secret                                         | `share.yml`             |
| Configure OSS, COS, AWS, or Kodo as object storage               | `openim-rpc-third.yml`  |
| Multi-end mutual kick strategy and max concurrent connections per gateway                 | `openim-msggateway.yml` |
| Offline message push configuration                                            | `openim-push.yml`       |
| Configure webhooks for callback notifications (e.g., before/after message send)         | `webhooks.yml`          |
| Whether new group members can view historical messages                          | `openim-rpc-group.yml`  |
| Token expiration time settings                                      | `openim-rpc-auth.yml`     |
| Scheduled task settings (e.g., how long to retain messages)                      | `openim-crontask.yml`   |

## Starting Multiple Instances of a Service and Maximum File Descriptors


To start multiple instances of an OpenIM service, simply add the corresponding port numbers and modify the `start-config.yml` file in the project’s root directory, 
then restart the service. For example, to start 2 instances of `openim-rpc-user`:

```yaml
rpc:
  registerIP: ''
  listenIP: 0.0.0.0
  ports: [ 10110, 10111 ]

prometheus:
  enable: true
  ports: [ 20100, 20101 ]
```

Modify`start-config.yml`:

```yaml
serviceBinaries:
  openim-rpc-user: 2
```

To set the maximum number of open file descriptors (typically one per online user):

```
maxFileDescriptors: 10000
```
