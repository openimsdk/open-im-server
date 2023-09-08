<p align="center">
    <a href="https://www.openim.online">
        <img src="./assets/logo-gif/openim-logo.gif" width="60%" height="30%"/>
    </a>
</p>

<h3 align="center" style="border-bottom: none">
    ⭐️  Open source Instant Messaging Server ⭐️ <br>
<h3>


<p align=center>
<a href="https://goreportcard.com/report/github.com/openimsdk/open-im-server"><img src="https://goreportcard.com/badge/github.com/openimsdk/open-im-server" alt="A+"></a>
<a href="https://github.com/openimsdk/open-im-server/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc+label%3A%22good+first+issue%22"><img src="https://img.shields.io/github/issues/openimsdk/open-im-server/good%20first%20issue?logo=%22github%22" alt="good first"></a>
<a href="https://github.com/openimsdk/open-im-server"><img src="https://img.shields.io/github/stars/openimsdk/open-im-server.svg?style=flat&logo=github&colorB=deeppink&label=stars"></a>
<a href="https://join.slack.com/t/openimsdk/shared_invite/zt-22720d66b-o_FvKxMTGXtcnnnHiMqe9Q"><img src="https://img.shields.io/badge/Slack-300%2B-blueviolet?logo=slack&amp;logoColor=white"></a>
<a href="https://github.com/openimsdk/open-im-server/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-Apache--2.0-green"></a>
<a href="https://golang.org/"><img src="https://img.shields.io/badge/Language-Go-blue.svg"></a>
</p>

</p>

<p align="center">
    <a href="./README.md"><b> English </b></a> •
    <a href="./README-zh_CN.md"><b> 简体中文 </b></a> •
    <a href="https://www.openim.online/en"><b> Docs </b></a>
</p>


</p>

## ✨ 关于 OpenIM

Open-IM-Server 是使用纯 Golang 精心制作的强大的即时消息服务器。其通过 JSON over WebSocket 进行通信的独特方法将每次交互都视为消息。这简化了定制并消除了修改服务器代码的需求。通过利用微服务架构，服务器可以通过集群部署，保证出色的性能和可伸缩性。

Open-IM-Server 不仅仅是一个即时消息服务器；它是将实时网络集成到您的应用程序中的强大工具，定位为您集成的首选选择！🚀

请注意，Open-IM-Server 不作为独立产品运行，也不提供内置的帐户注册或登录服务。为了简化您的实施过程，我们已开源了 [chat repository](https://github.com/OpenIMSDK/chat)，其中包括这些功能。与 Open-IM-Server 一起部署此聊天业务服务器可加快全面的聊天产品的设置。👥

为了进一步增强您的体验，我们还提供了 SDK 客户端，在其中实现了大多数复杂逻辑。可以在 [此链接](https://github.com/OpenIMSDK/openim-sdk-core) 找到 [SDK repository](https://github.com/OpenIMSDK/openim-sdk-core)。[chat repository](https://github.com/OpenIMSDK/chat) 是我们的业务服务器，而 'core' 代表 SDK 的高级封装，它们协同工作以提供卓越的结果。✨

## :star2: 为什么选择 OpenIM

**🔍 功能截图显示**

<div align="center">

|                      💻🔄📱 多终端同步 🔄🖥️                       |                        📅⚡ 高效会议 🚀💼                        |
| :----------------------------------------------------------: | :----------------------------------------------------------: |
| ![multiple-message](./assets/demo/multi-terminal-synchronization.png) | ![efficient-meetings](./assets/demo/efficient-meetings.png) |
|                    📲🔄 **一对一和群聊** 👥🗣️                    |               🎁💻 **特殊功能 - 自定义消息** ✉️🎨                |
| ![group-chat](./assets/demo/group-chat.png) | ![special-function](./assets/demo/special-function.png) |

</div>

1. **全面的消息类型支持 :speech_balloon:**

   ✅ 支持几乎所有类型的消息，包括文本、图片、表情符号、语音、视频、地理位置、文件、报价、名片、系统通知、自定义消息等

   ✅ 支持一对一和多人音视频通话

   ✅ 为 iOS、Android、Flutter、uni-app、ReactNative、Electron、Web、H5 等多个平台提供终端支持

2. **随时随地的高效会议 :earth_americas:**

   ✅ 基于具有 100% 可靠强制信令功能的 IM (Instant Messaging)，为与聊天应用程序深度集成的 IM 系统铺平了道路

   ✅ 支持单次会议中的数百人，订阅人数达到数千，以及服务器端音视频录制

3. **适用于各种社交场景的一对一和群聊 :busts_in_silhouette:**

   ✅ OpenIM 有四种角色：应用程序管理员、群主、群管理员和普通成员

   ✅ 强大的群特性，如静音、群公告、群验证、无限群成员和根据需要加载群消息

4. **独特的功能 :star2:**

   ✅ 支持读取并烧毁私人聊天，可自定义时长

   ✅ 消息编辑功能扩大了社交场景，使即时通讯变得更加多样化和有趣

5. **开源 :open_hands:**

   ✅ OpenIM 的代码是开源的，数据自控，旨在构建一个全球领先的 IM 开源社区，包括客户端 SDK 和服务器

   ✅ 基于开源服务器，已经开发了许多出色的开源项目，例如 [OpenKF](https://github.com/OpenIMSDK/OpenKF) (开源 AI 客户服务系统)

6. **易于扩展 :wrench:**

   ✅ OpenIM 服务器是用 Golang 实现的，引入了创新的 "一切都是消息" 通信模型，简化了自定义消息和扩展功能的实现

7. **高性能 :racing_car:**

   ✅ OpenIM 支持集群中的分层治理架构，经过大量用户的测试，并抽象了在线消息、离线消息和历史消息的存储模型

8. **全平台支持 :tv:**

   ✅ 支持原生 iOS、Android；跨平台 Flutter、uni-app、ReactNative；主要的 Web 前端框架如 React、Vue；小程序和 Electron 支持的 PC 平台

9. **终极部署体验 🤖**

   ✅ 支持 [集群部署](https://github.com/openimsdk/open-im-server/edit/main/deployments/README.md)

   ✅ 支持多架构镜像，我们的 Docker 镜像不仅托管在 GitHub 上，而且还在阿里云和 Docker Hub 上支持多个架构。请访问 [我们的 GitHub packages](https://github.com/orgs/OpenIMSDK/packages?repo_name=Open-IM-Server) 并阅读我们的 [版本管理文档](https://github.com/openimsdk/open-im-server/blob/main/docs/conversions/version.md) 以获取更多信息。

10. **开源社区的大生态系统 🤲**

    ✅ 我们有数万用户和许多解决方案来解决问题。 

    ✅  我们有一个大型的开源社区叫 [OpenIMSDK](https://github.com/OpenIMSDK)，它运行核心模块，我们还有一个开源社区叫 [openim-sigs](https://github.com/openim-sigs) 以探索更多基于 IM 的基础设施产品。

## :rocket: 快速开始

<details>   <summary>使用 Docker Compose 部署</summary>

1. 克隆项目

```
# 选择您需要的
BRANCH=release-v3.1
git clone -b $BRANCH https://github.com/openimsdk/open-im-server openim && export openim=$(pwd)/openim && cd $openim && make build
```

> **注意** 阅读我们的发布策略：https://github.com/openimsdk/open-im-server/blob/main/docs/conversions/version.md

1. 修改 `.env`

```
USER=root #无需修改
PASSWORD=openIM123  #8位或更多数字和字母的组合，此密码适用于redis、mysql、mongo，以及config/config.yaml中的accessSecret
ENDPOINT=http://127.0.0.1:10005 #minio的外部服务IP和端口，或使用域名storage.xx.xx，应用程序必须能够访问此IP和端口或域名，
API_URL=http://127.0.0.1:10002/object/ #应用程序必须能够访问此IP和端口或域名，
DATA_DIR=./  #指定大磁盘目录
```

1. 部署并启动

> **注意** 此命令只能执行一次。它会基于 `.env` 中的 `PASSWORD` 变量修改 docker-compose 中的组件密码，并修改 `config/config.yaml` 中的组件密码。如果 `.env` 中的密码发生变化，您需要首先执行 `docker-compose down`；`rm components -rf` 然后执行此命令。

```

make install
```

1. 检查服务

```

make check
```

![https://github.com/openimsdk/open-im-server/blob/main/docs/images/docker_build.png](https://github.com/openimsdk/open-im-server/blob/main/docs/images/docker_build.png)

</details>  <details>  <summary>从源码编译</summary>

您需要 `Go 1.18` 或更高版本，以及 `make`。

版本详情：https://github.com/openimsdk/open-im-server/blob/main/docs/conversions/version.md

```
# 选择您需要的
BRANCH=release-v3.1
git clone -b $BRANCH https://github.com/openimsdk/open-im-server openim && export openim=$(pwd)/openim && cd $openim && make build
```

阅读关于 [OpenIM 版本策略](https://github.com/openimsdk/open-im-server/blob/main/docs/conversions/version.md)

使用 `make help` 来查看 OpenIM 支持的指令。

如图所示，所有服务已成功构建

![成功编译](https://github.com/openimsdk/open-im-server/blob/main/docs/images/build.png)

</details> <details>  <summary>组件配置说明</summary>

config/config.yaml 文件为存储组件提供了详细的配置说明。

- Zookeeper

  - 用于 RPC 服务发现和注册，支持集群。

    ```
    zookeeper:
      schema: openim                          #不建议修改
      address: [ 127.0.0.1:2181 ]             #地址
      username:                               #用户名
      password:                               #密码
    ```

- MySQL

  - 用于存储用户、关系和群组，支持主从数据库。

    ```
    mysql:
      address: [ 127.0.0.1:13306 ]            #地址
      username: root                          #用户名
      password: openIM123                     #密码
      database: openIM_v2                     #不建议修改
      maxOpenConn: 1000                       #最大连接
      maxIdleConn: 100                        #最大空闲连接
      maxLifeTime: 60                         #连接可重用的最大时间(秒)
      logLevel: 4                             #日志级别 1=静音 2=错误 3=警告 4=信息
      slowThreshold: 500                      #慢语句阈值(毫秒)
    ```

- Mongo

  - 用于存储离线消息，支持 mongo 分片集群。

    ```
    mongo:
      uri:                                    #如果不为空，则直接使用此值
      address: [ 127.0.0.1:37017 ]            #地址
      database: openIM                        #默认 mongo 数据库
      username: root                          #用户名
      password: openIM123                     #密码
      maxPoolSize: 100                        #最大连接数
    ```

- Redis

  - 用于存储消息序列号、最新消息、用户令牌和 mysql 缓存，支持集群部署。

    ```
    redis:
      address: [ 127.0.0.1:16379 ]            #地址
      username:                               #用户名
      password: openIM123                     #密码
    ```

- Kafka

  - 用于消息队列，用于消息解耦，支持集群部署。

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

</details> <details>  <summary>启动和停止服务</summary>

启动服务

```

./scripts/start-all.sh;
```

检查服务

```

./scripts/check-all.sh
```

停止服务

```

./scripts/stop-all.sh
```

</details>

<details>  <summary>开放 IM 端口</summary>

| TCP 端口  | 描述                                                | 操作                                    |
| --------- | --------------------------------------------------- | --------------------------------------- |
| TCP:10001 | ws 协议，消息端口如消息发送、推送等，用于客户端 SDK | 端口释放或 nginx 反向代理，并关闭防火墙 |
| TCP:10002 | api 端口，如用户、朋友、组、消息接口。              | 端口释放或 nginx 反向代理，并关闭防火墙 |
| TCP:10005 | 选择 minio 存储时所需 (openIM 默认使用 minio 存储)  | 端口释放或 nginx 反向代理，并关闭防火墙 |

</details>  <details>  <summary>开放聊天端口</summary>

- 聊天仓库: https://github.com/OpenIMSDK/chat

| TCP 端口  | 描述                     | 操作                                    |
| --------- | ------------------------ | --------------------------------------- |
| TCP:10008 | 业务系统，如注册、登录等 | 端口释放或 nginx 反向代理，并关闭防火墙 |
| TCP:10009 | 管理后台，如统计、封禁等 | 端口释放或 nginx 反向代理，并关闭防火墙 |

</details>

## :link: APP 和 OpenIM 之间的关系

OpenIM 不仅仅是一个开源的即时消息组件，它是您的应用程序生态系统的一个不可分割的部分。查看此图表以了解 AppServer、AppClient、Open-IM-Server 和 Open-IM-SDK 如何互动。

![App-OpenIM 关系](https://github.com/openimsdk/open-im-server/blob/main/docs/images/open-im-server.png)

## :building_construction: 总体架构

深入了解 Open-IM-Server 的功能与我们的架构图。

![总体架构](https://github.com/openimsdk/open-im-server/blob/main/docs/images/Architecture.jpg)

## :hammer_and_wrench: 开始开发 OpenIM

OpenIM 我们的目标是建立一个顶级的开源社区。我们有一套标准，在 [Community repository](https://github.com/OpenIMSDK/community) 中。

如果您想为这个 Open-IM-Server 仓库做贡献，请阅读我们的 [贡献者文档](https://github.com/openimsdk/open-im-server/blob/main/CONTRIBUTING.md)。

在您开始之前，请确保您的更改是需要的。最好的方法是创建一个 [新的讨论](https://github.com/openimsdk/open-im-server/discussions/new/choose) 或 [Slack 通讯](https://join.slack.com/t/openimsdk/shared_invite/zt-22720d66b-o_FvKxMTGXtcnnnHiMqe9Q)，或者如果您发现一个问题，首先 [报告它](https://github.com/openimsdk/open-im-server/issues/new/choose)。

- [代码标准](https://github.com/openimsdk/open-im-server/blob/main/docs/conversions/go_code.md)
- [Docker 图像标准](https://github.com/openimsdk/open-im-server/blob/main/docs/conversions/images.md)
- [目录标准](https://github.com/openimsdk/open-im-server/blob/main/docs/conversions/directory.md)
- [提交标准](https://github.com/openimsdk/open-im-server/blob/main/docs/conversions/commit.md)
- [版本控制标准](https://github.com/openimsdk/open-im-server/blob/main/docs/conversions/version.md)
- [接口标准](https://github.com/openimsdk/open-im-server/blob/main/docs/conversions/api.md)
- [日志标准](https://github.com/openimsdk/open-im-server/blob/main/docs/conversions/)
- [错误代码标准](https://github.com/openimsdk/open-im-server/blob/main/docs/conversions/error_code.md)

## :busts_in_silhouette: 社区

- 📚 [OpenIM 社区](https://github.com/OpenIMSDK/community)
- 💕 [OpenIM 兴趣小组](https://github.com/Openim-sigs)
- 🚀 [加入我们的 Slack 社区](https://join.slack.com/t/openimsdk/shared_invite/zt-22720d66b-o_FvKxMTGXtcnnnHiMqe9Q)
- :eyes: [加入我们的微信群 (微信群)](https://openim-1253691595.cos.ap-nanjing.myqcloud.com/WechatIMG20.jpeg)

## :calendar: 社区会议

我们希望任何人都可以参与我们的社区并贡献代码，我们提供礼物和奖励，欢迎您每周四晚上加入我们。

我们的会议在 [OpenIM Slack](https://join.slack.com/t/openimsdk/shared_invite/zt-22720d66b-o_FvKxMTGXtcnnnHiMqe9Q) 🎯，然后您可以搜索 Open-IM-Server 管道加入。

我们在 [GitHub 讨论](https://github.com/openimsdk/open-im-server/discussions/categories/meeting) 中记下每次 [双周会议](https://github.com/orgs/OpenIMSDK/discussions/categories/meeting) 的笔记，我们的历史会议记录以及会议回放都可在 [Google Docs :bookmark_tabs:](https://docs.google.com/document/d/1nx8MDpuG74NASx081JcCpxPgDITNTpIIos0DS6Vr9GU/edit?usp=sharing) 中找到。

## :eyes: 谁在使用 OpenIM

查看我们的 [用户案例研究](https://github.com/OpenIMSDK/community/blob/main/ADOPTERS.md) 页面以获取项目用户列表。不要犹豫，留下一个 [📝评论](https://github.com/openimsdk/open-im-server/issues/379) 并分享您的使用案例。

## :page_facing_up: 许可证

OpenIM 根据 Apache 2.0 许可证授权。请查看 [LICENSE](https://github.com/openimsdk/open-im-server/tree/main/LICENSE) 以获取完整的许可证文本。

OpenIM logo，包括其变体和动画版本，在此存储库 [OpenIM](https://github.com/openimsdk/open-im-server) 下的 [assets/logo](./assets/logo) 和 [assets/logo-gif](./assets/logo-gif) 目录中显示，受版权法保护。

## 🔮 感谢我们的贡献者！

<a href="https://github.com/openimsdk/open-im-server/graphs/contributors">   <img src="https://contrib.rocks/image?repo=openimsdk/open-im-server" /> </a>