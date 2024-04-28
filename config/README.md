---
title: 'OpenIM Configuration Files and Common Configuration Item Modifications Guide'

## Configuration Files Explanation

| Configuration File              | Description                                                  |
| ------------------------------- | ------------------------------------------------------------ |
| **kafka.yml**                   | Configurations for Kafka username, password, address, etc.   |
| **redis.yml**                   | Configurations for Redis password, address, etc.             |
| **minio.yml**                   | Configurations for MinIO username, password, address, and external IP/domain; failing to modify external IP or domain may cause image file sending failures |
| **zookeeper.yml**               | Configurations for ZooKeeper user, password, address, etc.   |
| **mongodb.yml**                 | Configurations for MongoDB username, password, address, etc. |
| **log.yml**                     | Configurations for log level and storage directory.          |
| **notification.yml**            | Configurations for events like adding friends, creating groups, etc. |
| **share.yml**                   | Common configurations needed by various OpenIM services, such as secret. |
| **webhooks.yml**                | Configurations for URLs in Webhook.                          |
| **local-cache.yml**             | Local cache configurations.                                  |
| **openim-rpc-third.yml**        | Configurations for listening IP, port, and storage settings for images and videos in openim-rpc-third service. |
| **openim-rpc-user.yml**         | Configurations for listening IP and port in openim-rpc-user service. |
| **openim-api.yml**              | Configurations for listening IP, port, etc., in openim-api service. |
| **openim-crontask.yml**         | Configurations for openim-crontask service.                  |
| **openim-msggateway.yml**       | Configurations for listening IP, port, etc., in openim-msggateway service. |
| **openim-msgtransfer.yml**      | Configurations for openim-msgtransfer service.               |
| **openim-push.yml**             | Configurations for listening IP, port, and offline push settings in openim-push service. |
| **openim-rpc-auth.yml**         | Configurations for listening IP, port, and token expiration settings in openim-rpc-auth service. |
| **openim-rpc-conversation.yml** | Configurations for listening IP, port, etc., in openim-rpc-conversation service. |
| **openim-rpc-friend.yml**       | Configurations for listening IP, port, etc., in openim-rpc-friend service. |
| **openim-rpc-group.yml**        | Configurations for listening IP, port, etc., in openim-rpc-group service. |
| **openim-rpc-msg.yml**          | Configurations for listening IP, port, and whether to verify friendship before sending messages in openim-rpc-msg service. |

## Common Configuration Item Modifications

| Configuration Item Modification                       | Configuration File      |
| ----------------------------------------------------- | ----------------------- |
| Using MinIO for image and video file object storage   | `minio.yml`             |
| Adjusting production environment logs                 | `log.yml`               |
| Verifying friendship before sending messages          | `openim-rpc-msg.yml`    |
| Modifying secret                                      | `share.yml`             |
| Using OSS, COS, AWS, Kodo for image and video storage | `openim-rpc-third.yml`  |
| Setting multiple login policy                         | `openim-msggateway.yml` |
| Setting up offline push                               | `openim-push.yml`       |

## Starting Multiple Instances of an OpenIM Service

To start multiple instances of an OpenIM service, simply increase the corresponding port numbers and modify the `start-config.yml` file in the project root directory. Restart the service to take effect. For example, the configuration to start 2 instances of `openim-rpc-user` is as follows:

```yaml
rpc:
  registerIP: ''
  listenIP: 0.0.0.0
  ports: [ 10110, 10111 ]

prometheus:
  enable: true
  ports: [ 20100, 20101 ]
```

Modify `start-config.yml`:

```yaml
serviceBinaries:
  openim-rpc-user: 2
```
