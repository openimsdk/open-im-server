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
    <a href="./README_zh.md"><b>中文</b></a>
</p>

</p>

## Open-IM-Server 是什么

Open-IM-Server 是一款即时通讯服务器，使用纯 Golang 开发，采用 JSON over WebSocket 传输协议。在 Open-IM-Server 中，所有东西都是消息，因此您可以轻松扩展自定义消息，而无需修改服务器代码。使用微服务架构，Open-IM-Server 可以使用集群进行部署。通过在客户端服务器上部署 Open-IM-Server，开发人员可以免费快速地将即时通讯和实时网络功能集成到自己的应用程序中，并确保业务数据的安全性和隐私性。

## 特点

+ 免费
+ 可扩展架构
+ 易于集成
+ 良好的可扩展性
+ 高性能
+ 轻量级
+ 支持多种协议

## 社区

+ 访问中文官方网站：[Open-IM中文开发文档](https://doc.rentsoft.cn/)

## 快速开始

### 安装Open-IM-Server

> Open-IM 依赖于五个开源高性能组件：ETCD、MySQL、MongoDB、Redis 和 Kafka。在部署 Open-IM-Server 之前，请确保已安装上述五个组件。如果您的服务器没有上述组件，则必须首先安装缺失组件。如果您已经拥有上述组件，则建议直接使用它们。如果没有，则建议使用 Docker-compose，无需安装依赖项，一键部署，更快更方便。

### 使用 Docker 部署

1. 安装 [Go 环境](https://golang.org/doc/install)。确保 Go 版本至少为 1.17。

2. 克隆 Open-IM 项目到您的服务器

   `git clone <https://github.com/OpenIMSDK/Open-IM-Server.git> --recursive`

3. 部署

   1. 修改 env

      ```
      #cd Open-IM-server
      USER=root
      PASSWORD=openIM123    #密码至少8位数字，不包括特殊字符
      ENDPOINT=http://127.0.0.1:10005 #请用互联网IP替换127.0.0.1
      DATA_DIR=./
      ```

   2. 部署和启动

      ```
      chmod +x install_im_server.sh;
      ./install_im_server.sh;
      ```

   3. 检查服务

      ```
      cd script;
      ./docker_check_service.sh
      ./check_all.sh
      ```

      ![https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/Open-IM-Servers-on-System.png](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/Open-IM-Servers-on-System.png)

### 使用源代码部署

1. Go 1.17 或更高版本。

2. 克隆

   ```
   git clone <https://github.com/OpenIMSDK/Open-IM-Server.git> --recursive
   cd cmd/openim-sdk-core
   git checkout main
   ```

3. 设置可执行权限

   ```
   cd ../../script/
   chmod +x *.sh
   ```

4. 构建

   ```
   ./batch_build_all_service.sh
   ```

所有服务已成功构建

### 配置说明

> Open-IM 配置分为基本组件配置和业务内部服务配置。当使用产品时，开发人员需要将每个组件的地址填写为其服务器组件的地址，并确保业务内部服务端口未被占用。

### 基本组件配置说明

+ ETCD
  + Etcd 用于 RPC 服务的发现和注册，Etcd Schema 是注册名称的前缀，建议将其修改为公司名称，Etcd 地址(ip+port)支持集群部署，可以填写多个 ETCD 地址，也可以只有一个 etcd 地址。
+ MySQL
  + MySQL 用于消息和用户关系的全存储，暂时不支持集群部署。修改地址、用户、密码和数据库名称。
+ Mongo
  + Mongo 用于消息的离线存储，默认存储 7 天。暂时不支持集群部署。只需修改地址和数据库名称即可。
+ Redis
  + Redis 目前主要用于消息序列号存储和用户令牌信息存储。暂时不支持集群部署。只需修改相应的 Redis 地址和密码即可。
+ Kafka
  + Kafka 用作消息传输存储队列，支持集群部署，只需修改相应的地址。

### 内部服务配置说明

+ credential&&push
  + Open-IM 需要使用三方离线推送功能。目前使用的是腾讯的三方推送，支持 IOS、Android 和 OSX 推送。这些信息是腾讯推送的一些注册信息，开发人员需要去腾讯云移动推送注册相应的信息。如果您没有填写相应的信息，则无法使用离线消息推送功能。
+ api&&rpcport&&longconnsvr&&rpcregistername
  + API 端口是 HTTP 接口，longconnsvr 是 WebSocket 监听端口，rpcport 是内部服务启动端口。两者都支持集群部署。请确保这些端口未被使用。如果要为单个服务打开多个服务，请填写多个以逗号分隔的端口。rpcregistername 是每个服务在注册表 Etcd 中注册的服务名称，无需修改。
+ log&&modulename
  + 日志配置包括日志文件的存储路径，日志发送到 Elasticsearch 进行日志查看。目前不支持将日志发送到 Elasticsearch。暂时不需要修改配置。modulename 用于根据服务模块的名称拆分日志。默认配置可以。

### 脚本说明

> Open-IM 脚本提供服务编译、启动和停止脚本。有四个 Open-IM 脚本启动模块，一个是 http+rpc 服务启动模块，第二个是 WebSocket 服务启动模块，然后是 msg_transfer 模块，最后是 push 模块。

+ path_info.cfg&&style_info.cfg&&

  functions.sh

  + 包含每个模块的路径信息，包括源代码所在的路径、服务启动名称、shell 打印字体样式以及一些用于处理 shell 字符串的函数。

+ build_all_service.sh

  + 编译模块，将 Open-IM 的所有源代码编译为二进制文件并放入 bin 目录。

+ start_rpc_api_service.sh&&msg_gateway_start.sh&&msg_transfer_start.sh&&push_start.sh

  + 独立脚本启动模块，后跟 API 和 RPC 模块、消息网关模块、消息传输模块和推送模块。

+ start_all.sh&&stop_all.sh

  + 总脚本，启动所有服务和关闭所有服务。

## 认证流程图

![https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/open-im-server.png](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/open-im-server.png)

## 架构

![https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/Architecture.jpg](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/Architecture.jpg)

## 开始开发 OpenIM

[社区存储库](https://github.com/OpenIMSDK/community)包含有关从源代码构建 Kubernetes、如何贡献代码和文档、有关什么的联系人等所有信息。

## 贡献

欢迎对该项目进行贡献！请参见 [CONTRIBUTING.md](http://CONTRIBUTING.md) 了解详细信息。

## 社区会议

我们希望任何人都能参与我们的社区，我们提供礼品和奖励，并欢迎您每周四晚上加入我们。

我们在 [GitHub 讨论](https://github.com/OpenIMSDK/Open-IM-Server/discussions/categories/meeting) 中记录每个 [两周会议](https://github.com/OpenIMSDK/Open-IM-Server/issues/381)，我们的记录写在 [Google 文档](https://docs.google.com/document/d/1nx8MDpuG74NASx081JcCpxPgDITNTpIIos0DS6Vr9GU/edit?usp=sharing) 中。

## 谁在使用 Open-IM-Server

[用户案例研究](https://github.com/OpenIMSDK/community/blob/main/ADOPTERS.md) 页面包括该项目的用户列表。您可以留下 [📝评论](https://github.com/OpenIMSDK/Open-IM-Server/issues/379) 让我们知道您的用例。

## ⚠️公告

尊敬的 OpenIM 存储库用户，我们很高兴地宣布，我们正在进行重大改进，以改善我们的服务质量和用户体验。我们将使用 [errcode](https://github.com/OpenIMSDK/Open-IM-Server/tree/errcode) 分支对主分支进行广泛的更新和改进，确保我们的代码存储库处于最佳状态。

在此期间，我们需要暂停主分支上的 PR 和问题处理。我们理解这可能会给您带来一些不便，但我们相信它将为我们提供更好的服务和更可靠的代码存储库。如果您需要在此期间提交代码，请将其提交到 [errcode](https://github.com/OpenIMSDK/Open-IM-Server/tree/errcode) 分支，我们将尽快处理您的请求。

我们感谢您的支持和信任，以及您在整个过程中的耐心和理解。我们重视您的贡献和建议，这是我们持续改进和成长的动力。

我们预计这项工作将很快完成，并尽最大努力将对您的影响最小化。再次感谢您的支持和谅解。

谢谢！

![https://github.com/OpenIMSDK/OpenIM-Docs/blob/main/docs/images/WechatIMG20.jpeg](https://github.com/OpenIMSDK/OpenIM-Docs/blob/main/docs/images/WechatIMG20.jpeg)

## 许可证

Open-IM-Server 使用 Apache 2.0 许可证。有关详情，请参阅 LICENSE 文件。
