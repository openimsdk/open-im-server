<p align="center">
    <a href="https://openim.io">
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
    <a href="https://openim.io/en"><b> Docs </b></a>
</p>


</p>

## 🟢 扫描微信进群交流
<img src="./docs/images/Wechat.jpg" width="300">


## Ⓜ️ 关于 OpenIM

OpenIM 是一个专门设计用于在应用程序中集成聊天、音视频通话、通知以及AI聊天机器人等通信功能的服务平台。它通过提供一系列强大的API和Webhooks，使开发者可以轻松地在他们的应用中加入这些交互特性。OpenIM 本身并不是一个独立运行的聊天应用，而是作为一个平台，为其他应用提供支持，实现丰富的通信功能。下图展示 AppServer、AppClient、OpenIMServer 和 OpenIMSDK 之间的交互关系来具体说明。





![App-OpenIM 关系](./docs/images/oepnim-design.png)

## 🚀 关于 OpenIMSDK

**OpenIMSDK** 是为 **OpenIMServer** 设计的IM SDK，专为嵌入客户端应用而生。其主要功能及模块如下：

+ 🌟 主要功能：

  - 📦 本地存储
  - 🔔 监听器回调
  - 🛡️ API封装
  - 🌐 连接管理

  ## 📚 主要模块：

  1. 🚀 初始化及登录
  2. 👤 用户管理
  3. 👫 好友管理
  4. 🤖 群组功能
  5. 💬 会话处理

它使用 Golang 构建，并支持跨平台部署，确保在所有平台上提供一致的接入体验。

👉 **[探索 GO SDK](https://github.com/openimsdk/openim-sdk-core)**

## 🌐 关于 OpenIMServer

+ **OpenIMServer** 具有以下特点：
  - 🌐 微服务架构：支持集群模式，包括网关(gateway)和多个rpc服务。
  - 🚀 部署方式多样：支持源代码、kubernetes或docker部署。
  - 海量用户支持：十万超级大群，千万用户，及百亿消息

### 增强的业务功能：

+ **REST API**：OpenIMServer 提供了REST API供业务系统使用，旨在赋予业务更多功能，例如通过后台接口建立群组、发送推送消息等。
+ **Webhooks**：OpenIMServer提供了回调能力以扩展更多的业务形态，所谓回调，即OpenIMServer会在某一事件发生之前或者之后，向业务服务器发送请求，如发送消息之前或之后的回调。

👉 **[了解更多](https://docs.openim.io/guides/introduction/product)**

## :rocket: 快速开始

在线体验iOS/Android/H5/PC/Web：

👉 **[OpenIM online demo](https://www.openim.io/zh/commercial)**

🤲 为了方便用户体验，我们提供了多种部署解决方案，您可以根据下面的列表选择自己的部署方法：

+ **[源代码部署指南](https://docs.openim.io/guides/gettingStarted/imSourceCodeDeployment)**
+ **[Docker 部署指南](https://docs.openim.io/guides/gettingStarted/dockerCompose)**
+ **[Kubernetes 部署指南](https://docs.openim.io/guides/gettingStarted/k8s-deployment)**

## :hammer_and_wrench: 开始开发 OpenIM

OpenIM 我们的目标是建立一个顶级的开源社区。我们有一套标准，在[社区仓库](https://github.com/OpenIMSDK/community)中。

如果你想为这个 Open-IM-Server 仓库做贡献，请阅读我们的[贡献者文档](https://github.com/openimsdk/open-im-server/blob/main/CONTRIBUTING.md)。

在开始之前，请确保你的更改是有需求的。最好的方法是创建一个[新的讨论](https://github.com/openimsdk/open-im-server/discussions/new/choose) 或 [Slack 通信](https://join.slack.com/t/openimsdk/shared_invite/zt-22720d66b-o_FvKxMTGXtcnnnHiMqe9Q)，或者如果你发现一个问题，首先[报告它](https://github.com/openimsdk/open-im-server/issues/new/choose)。

+ [代码标准](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/go-code.md)

+ [Docker 镜像标准](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/images.md)

+ [目录标准](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/directory.md)

+ [提交标准](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/commit.md)

+ [版本控制标准](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/version.md)

+ [接口标准](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/interface.md)

+ [OpenIM配置和环境变量设置](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/environment.md)

> **Note**
> 针对中国的用户，阅读我们的 [Docker 镜像标准](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/images.md) 以便使用国内 aliyun 的镜像地址。OpenIM 也有针对中国的 gitee 同步仓库，你可以在 [gitee.com](https://gitee.com/openimsdk) 上找到它。

## :link: 链接

  + **[完整文档](https://doc.rentsoft.cn/)**
  + **[更新日志](https://github.com/openimsdk/open-im-server/blob/main/CHANGELOG.md)**
  + **[FAQ](https://github.com/openimsdk/open-im-server/blob/main/FAQ.md)**
  + **[代码示例](https://github.com/openimsdk/open-im-server/blob/main/examples)**

## :handshake: 社区

  + **[GitHub Discussions](https://github.com/openimsdk/open-im-server/discussions)**
  + **[Slack 通信](https://join.slack.com/t/openimsdk/shared_invite/zt-22720d66b-o_FvKxMTGXtcnnnHiMqe9Q)**
  + **[GitHub Issues](https://github.com/openimsdk/open-im-server/issues)**

  您可以加入这些平台，讨论问题，提出建议，或分享您的成功故事！

## :writing_hand: 贡献

  我们欢迎任何形式的贡献！请确保在提交 Pull Request 之前阅读我们的[贡献者文档](https://github.com/openimsdk/open-im-server/blob/main/CONTRIBUTING.md)。

  + **[报告 Bug](https://github.com/openimsdk/open-im-server/issues/new?assignees=&labels=bug&template=bug_report.md&title=)**
  + **[提出新特性](https://github.com/openimsdk/open-im-server/issues/new?assignees=&labels=enhancement&template=feature_request.md&title=)**
  + **[提交 Pull Request](https://github.com/openimsdk/open-im-server/pulls)**

  感谢您的贡献，我们一起打造一个强大的即时通信解决方案！

## :closed_book: 许可证

  OpenIMSDK 在 Apache License 2.0 许可下可用。查看[LICENSE 文件](https://github.com/openimsdk/open-im-server/blob/main/LICENSE)了解更多信息。

## 🔮 Thanks to our contributors!

<a href="https://github.com/openimsdk/open-im-server/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=openimsdk/open-im-server" />
</a>
