<h1 align="center" style="border-bottom: none">
    <b>
        <a href="https://doc.rentsoft.cn/">Open IM Server</a><br>
    </b>
    ⭐️  Open source Instant Messaging Server  ⭐️ <br>
</h1>


<p align=center>
<a href="https://goreportcard.com/report/github.com/OpenIMSDK/Open-IM-Server"><img src="https://goreportcard.com/badge/github.com/OpenIMSDK/Open-IM-Server" alt="A+"></a>
<a href="https://github.com/OpenIMSDK/Open-IM-Server/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc+label%3A%22good+first+issue%22"><img src="https://img.shields.io/github/issues/OpenIMSDK/Open-IM-Server/good%20first%20issue?logo=%22github%22" alt="good first"></a>
<a href="https://github.com/OpenIMSDK/Open-IM-Server"><img src="https://img.shields.io/github/stars/OpenIMSDK/Open-IM-Server.svg?style=flat&logo=github&colorB=deeppink&label=stars"></a>
<a href="https://join.slack.com/t/openimsdk/shared_invite/zt-1tmoj26uf-_FDy3dowVHBiGvLk9e5Xkg"><img src="https://img.shields.io/badge/Slack-100%2B-blueviolet?logo=slack&amp;logoColor=white"></a>
<a href="https://github.com/OpenIMSDK/Open-IM-Server/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-Apache--2.0-green"></a>
<a href="https://golang.org/"><img src="https://img.shields.io/badge/Language-Go-blue.svg"></a>
</p>

</p>

<p align="center">
    <a href="./README.md"><b>English</b></a> •
    <a href="./README_zh-CN.md"><b>中文</b></a>
</p>

</p>

## Open-IM-Server 是什么

Open-IM-Server 是一款即时通讯服务器，使用纯 Golang 开发，采用 JSON over WebSocket 传输协议。在 Open-IM-Server 中，所有东西都是消息，因此您可以轻松扩展自定义消息，而无需修改服务器代码。使用微服务架构，Open-IM-Server 可以使用集群进行部署。通过在服务器上部署 Open-IM-Server，开发人员可以快速地将即时通讯和实时网络功能集成到自己的应用程序中，并确保业务数据的安全性和隐私性。

Open-IM-Server并不是一个独立的产品，本身不包含账号的注册和登录服务。
为方便大家测试，我们开源了包括登录注册功能的 [chat 仓库](https://github.com/OpenIMSDK/chat)，chat 业务服务端和 Open-IM-Server 一起部署，即可搭建一个聊天产品。

## 特点

+ 开源
+ 易于集成
+ 良好的可扩展性
+ 高性能
+ 轻量级
+ 支持多种协议

## 社区

+ 访问中文官方网站：[OpenIM中文开发文档](https://doc.rentsoft.cn/)

## 快速开始

### 安装Open-IM-Server

> Open-IM-Server依赖于五个开源组件：Zookeeper、MySQL、MongoDB、Redis 和 Kafka。在部署 Open-IM-Server 之前，请确保已安装上述五个组件。如果没有，则建议使用 docker-compose，一键部署，方便快捷。

### 使用 docker-compose 部署

1. 隆项目

```
git clone https://github.com/OpenIMSDK/Open-IM-Server 
cd Open-IM-Server
git checkout release-v3.0 #or other release branch
```

2. 修改 env

```
此处主要修改相关组件密码
USER=root #不用修改
PASSWORD=openIM123  #8位以上的数字和字母组合密码，密码对redis mysql mongo生效，以及config/config.yaml中的accessSecret
ENDPOINT=http://127.0.0.1:10005 #minio对外服务的ip和端口，或用域名storage.xx.xx，app要能访问到此ip和端口或域名，
API_URL=http://127.0.0.1:10002/object/ #app要能访问到此ip和端口或域名，
DATA_DIR=./  #指定大磁盘目录
```

3. 部署和启动

注意：此命令只能执行一次，它会根据.env 中的 PASSWORD 变量修改 docker-compose 中组件密码，并修改 config/config.yaml 中的组件密码
如果.env 中的密码变了，需要先 docker-compose down ; rm components -rf 后再执行此命令。

```
chmod +x install_im_server.sh;
./install_im_server.sh;
```

4. 检查服务

```
cd scripts;
./docker_check_service.sh
```

![https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/images/docker_build.png](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/images/docker_build.png)

5. 开放 IM 端口

| TCP 端口  | 说明                                                  | 操作                                    |
| --------- | ----------------------------------------------------- | --------------------------------------- |
| TCP:10001 | ws 协议，消息端口，如消息发送、推送等，用于客户端 SDK | 端口放行或 nginx 反向代理，并关闭防火墙 |
| TCP:10002 | api 端口，如用户、好友、群组、消息等接口。            | 端口放行或 nginx 反向代理，并关闭防火墙 |
| TCP:10005 | 选择 minio 存储时需要(openIM 默认使用 minio 存储)     | 端口放行或 nginx 反向代理，并关闭防火墙 |

6. 开放 Chat 端口

| TCP 端口  | 说明                     | 操作                                    |
| --------- | ------------------------ | --------------------------------------- |
| TCP:10008 | 业务系统，如注册、登录等 | 端口放行或 nginx 反向代理，并关闭防火墙 |
| TCP:10009 | 管理后台，如统计、封号等 | 端口放行或 nginx 反向代理，并关闭防火墙 |

### 使用源代码部署

1. Go 1.18或更高版本。

2. 克隆

   ```
   git clone https://github.com/OpenIMSDK/Open-IM-Server 
   cd Open-IM-Server
   git checkout release-v3.0 #or other release branch
   ```

3. 编译

   ```
   cd Open-IM-server/scripts
   chmod +x *.sh
   ./build_all_service.sh
   ```

所有服务已成功构建如图所示

![编译成功](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/images/build.png)

> 

### 组件配置说明

config/config.yaml中针对存储组件有详细的配置说明

+ Zookeeper
  + 用于RPC 服务发现和注册，支持集群。
  
    ````
    ```
    zookeeper:
      schema: openim                          #不建议修改
      address: [ 127.0.0.1:2181 ]             #地址
      username:                               #用户名
      password:                               #密码
    ```
    ````
  
    
  
+ MySQL
  
  + 用于存储用户、关系链、群组，支持数据库主备。
  
    ```
    mysql:
      address: [ 127.0.0.1:13306 ]            #地址
      username: root                          #用户名
      password: openIM123                     #密码
      database: openIM_v2                     #不建议修改
      maxOpenConn: 1000                       #最大连接数
      maxIdleConn: 100                        #最大空闲连接数
      maxLifeTime: 60                         #连接可以重复使用的最长时间（秒）
      logLevel: 4                             #日志级别 1=slient 2=error 3=warn 4=info
      slowThreshold: 500                      #慢语句阈值 （毫秒）
    ```
  
    
  
+ Mongo
  + 用于存储离线消息，支持mongo分片集群。
  
    ```
    mongo:
      uri:                                    #不为空则直接使用该值
      address: [ 127.0.0.1:37017 ]            #地址
      database: openIM                        #mongo db 默认即可
      username: root                          #用户名
      password: openIM123                     #密码
      maxPoolSize: 100                        #最大连接数
    ```
  
+ Redis
  + 用于存储消息序列号、最新消息、用户token及mysql缓存，支持集群部署。
  
    ```
    redis:
      address: [ 127.0.0.1:16379 ]            #地址
      username:                               #用户名
      password: openIM123                     #密码
    ```
  
+ Kafka
  + 用于消息队列，用于消息解耦，支持集群部署。
  
    ```
    kafka:
      username:                               #用户名
      password:                               #密码
      addr: [ 127.0.0.1:9092 ]                #地址
      latestMsgToRedis:
        topic: "latestMsgToRedis"
      offlineMsgToMongo:
        topic: "offlineMsgToMongoMysql"
      msgToPush:
        topic: "msqToPush"
      msgToModify:
        topic: "msgToModify"
      consumerGroupID:
        msgToRedis: redis
        msgToMongo: mongo
        msgToMySql: mysql
        msgToPush: push
        msgToModify: modify
    ```

## APP和OpenIM关系

OpenIM 是开源的即时通讯组件，它并不是一个独立的产品，此图展示了AppServer、AppClient、Open-IM-Server以及Open-IM-SDK之间的关系

![https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/open-im-server.png](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/images/open-im-server.png)

## 整体架构

![https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/Architecture.jpg](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/images/Architecture.jpg)

## 开始开发 OpenIM

[社区存储库](https://github.com/OpenIMSDK/community)包含有关从源代码构建 Kubernetes、如何贡献代码和文档。

## 贡献

欢迎对该项目进行贡献！请参见 [CONTRIBUTING.md](http://CONTRIBUTING.md) 了解详细信息。

## 社区会议

我们希望任何人都能参与我们的社区，我们提供礼品和奖励，并欢迎您每周四晚上加入我们。

我们在 [GitHub 讨论](https://github.com/OpenIMSDK/Open-IM-Server/discussions/categories/meeting) 中记录每个 [两周会议](https://github.com/OpenIMSDK/Open-IM-Server/issues/381)，我们的记录写在 [Google 文档](https://docs.google.com/document/d/1nx8MDpuG74NASx081JcCpxPgDITNTpIIos0DS6Vr9GU/edit?usp=sharing) 中。

## 谁在使用 Open-IM-Server

[用户案例研究](https://github.com/OpenIMSDK/community/blob/main/ADOPTERS.md) 页面包括该项目的用户列表。您可以留下 [📝评论](https://github.com/OpenIMSDK/Open-IM-Server/issues/379) 让我们知道您的用例。

![https://github.com/OpenIMSDK/OpenIM-Docs/blob/main/docs/images/WechatIMG20.jpeg](https://github.com/OpenIMSDK/OpenIM-Docs/blob/main/docs/images/WechatIMG20.jpeg)

## 许可证

Open-IM-Server 使用 Apache 2.0 许可证。有关详情，请参阅 LICENSE 文件。
