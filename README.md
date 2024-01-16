<p align="center">
    <a href="https://openim.io">
        <img src="./assets/logo-gif/openim-logo.gif" width="60%" height="30%"/>
    </a>
</p>

<div align="center">

[![Stars](https://img.shields.io/github/stars/openimsdk/open-im-server?style=for-the-badge&logo=github&colorB=ff69b4)](https://github.com/openimsdk/open-im-server/stargazers)
[![Forks](https://img.shields.io/github/forks/openimsdk/open-im-server?style=for-the-badge&logo=github&colorB=blue)](https://github.com/openimsdk/open-im-server/network/members)
[![Codecov](https://img.shields.io/codecov/c/github/openimsdk/open-im-server?style=for-the-badge&logo=codecov&colorB=orange)](https://app.codecov.io/gh/openimsdk/open-im-server)
[![Go Report Card](https://goreportcard.com/badge/github.com/openimsdk/open-im-server?style=for-the-badge)](https://goreportcard.com/report/github.com/openimsdk/open-im-server)
[![Go Reference](https://img.shields.io/badge/Go%20Reference-blue.svg?style=for-the-badge&logo=go&logoColor=white)](https://pkg.go.dev/github.com/openimsdk/open-im-server/v3)
[![License](https://img.shields.io/badge/license-Apache--2.0-green?style=for-the-badge)](https://github.com/openimsdk/open-im-server/blob/main/LICENSE)
[![Slack](https://img.shields.io/badge/Slack-500%2B-blueviolet?style=for-the-badge&logo=slack&logoColor=white)](https://join.slack.com/t/openimsdk/shared_invite/zt-22720d66b-o_FvKxMTGXtcnnnHiMqe9Q)
[![Best Practices](https://img.shields.io/badge/Best%20Practices-purple?style=for-the-badge)](https://www.bestpractices.dev/projects/8045)
[![Good First Issues](https://img.shields.io/github/issues/openimsdk/open-im-server/good%20first%20issue?style=for-the-badge&logo=github)](https://github.com/openimsdk/open-im-server/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc+label%3A%22good+first+issue%22)
[![Language](https://img.shields.io/badge/Language-Go-blue.svg?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)

[**English**](./README.md) •
[**简体中文**](./README-zh_CN.md) •
[**Docs**](https://openim.io/en)

</div>

</p>

## Ⓜ️ About OpenIM

OpenIM is a service platform specifically designed for integrating chat, audio-video calls, notifications, and AI chatbots into applications. It provides a range of powerful APIs and Webhooks, enabling developers to easily incorporate these interactive features into their applications. OpenIM is not a standalone chat application, but rather serves as a platform to support other applications in achieving rich communication functionalities. The following diagram illustrates the interaction between AppServer, AppClient, OpenIMServer, and OpenIMSDK to explain in detail.

![App-OpenIM Relationship](./docs/images/oepnim-design.png)

## 🚀 About OpenIMSDK

**OpenIMSDK** is an IM SDK designed for **OpenIMServer**, created specifically for embedding in client applications. Its main features and modules are as follows:

+ 🌟 Main Features:

  - 📦 Local storage
  - 🔔 Listener callbacks
  - 🛡️ API wrapping
  - 🌐 Connection management

  ## 📚 Main Modules:

  1. 🚀 Initialization and Login
  2. 👤 User Management
  3. 👫 Friend Management
  4. 🤖 Group Functions
  5. 💬 Conversation Handling

It is built using Golang and supports cross-platform deployment, ensuring a consistent access experience across all platforms.

👉 **[Explore GO SDK](https://github.com/openimsdk/openim-sdk-core)**

## 🌐 About OpenIMServer

+ **OpenIMServer** has the following characteristics:
  - 🌐 Microservice architecture: Supports cluster mode, including a gateway and multiple rpc services.
  - 🚀 Diverse deployment methods: Supports deployment via source code, Kubernetes, or Docker.
  - Support for massive user base: Super large groups with hundreds of thousands of users, tens of millions of users, and billions of messages.

### Enhanced Business Functionality:

+ **REST API**: OpenIMServer offers REST APIs for business systems, aimed at empowering businesses with more functionalities, such as creating groups and sending push messages through backend interfaces.
+ **Webhooks**: OpenIMServer provides callback capabilities to extend more business forms. A callback means that OpenIMServer sends a request to the business server before or after a certain event, like callbacks before or after sending a message.

👉 **[Learn more](https://docs.openim.io/guides/introduction/product)**

## :rocket: Quick Start

🤲 To facilitate user experience, we offer various deployment solutions. You can choose your deployment method from the list below:

+ **[Source Code Deployment Guide](https://docs.openim.io/guides/gettingStarted/imSourceCodeDeployment)**
+ **[Docker Deployment Guide](https://docs.openim.io/guides/gettingStarted/dockerCompose)**
+ **[Kubernetes Deployment Guide](https://docs.openim.io/guides/gettingStarted/k8s-deployment)**

## :hammer_and_wrench: To Start Developing OpenIM

[![Open in Dev Container](https://img.shields.io/static/v1?label=Dev%20Container&message=Open&color=blue&logo=visualstudiocode)](https://vscode.dev/github/openimsdk/open-im-server)

OpenIM Our goal is to build a top-level open source community. We have a set of standards, in the [Community repository](https://github.com/OpenIMSDK/community).

If you'd like to contribute to this Open-IM-Server repository, please read our [contributor documentation](https://github.com/openimsdk/open-im-server/blob/main/CONTRIBUTING.md).

Before you start, please make sure your changes are in demand. The best for that is to create a [new discussion](https://github.com/openimsdk/open-im-server/discussions/new/choose) OR [Slack Communication](https://join.slack.com/t/openimsdk/shared_invite/zt-22720d66b-o_FvKxMTGXtcnnnHiMqe9Q), or if you find an issue, [report it](https://github.com/openimsdk/open-im-server/issues/new/choose) first.

- [OpenIM API Reference](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/api.md)
- [OpenIM Bash Logging](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/bash-log.md)
- [OpenIM CI/CD Actions](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/cicd-actions.md)
- [OpenIM Code Conventions](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/code-conventions.md)
- [OpenIM Commit Guidelines](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/commit.md)
- [OpenIM Development Guide](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/development.md)
- [OpenIM Directory Structure](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/directory.md)
- [OpenIM Environment Setup](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/environment.md)
- [OpenIM Error Code Reference](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/error-code.md)
- [OpenIM Git Workflow](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/git-workflow.md)
- [OpenIM Git Cherry Pick Guide](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/gitcherry-pick.md)
- [OpenIM GitHub Workflow](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/github-workflow.md)
- [OpenIM Go Code Standards](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/go-code.md)
- [OpenIM Image Guidelines](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/images.md)
- [OpenIM Initial Configuration](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/init-config.md)
- [OpenIM Docker Installation Guide](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/install-docker.md)
- [OpenIM OpenIM Linux System Installation](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/install-openim-linux-system.md)
- [OpenIM Linux Development Guide](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/linux-development.md)
- [OpenIM Local Actions Guide](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/local-actions.md)
- [OpenIM Logging Conventions](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/logging.md)
- [OpenIM Offline Deployment](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/offline-deployment.md)
- [OpenIM Protoc Tools](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/protoc-tools.md)
- [OpenIM Testing Guide](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/test.md)
- [OpenIM Utility Go](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/util-go.md)
- [OpenIM Makefile Utilities](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/util-makefile.md)
- [OpenIM Script Utilities](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/util-scripts.md)
- [OpenIM Versioning](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/version.md)
- [Manage backend and monitor deployment](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/prometheus-grafana.md)
- [Mac Developer Deployment Guide for OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/mac-developer-deployment-guide.md)


## :busts_in_silhouette: Community

+ 📚 [OpenIM Community](https://github.com/OpenIMSDK/community)
+ 💕 [OpenIM Interest Group](https://github.com/Openim-sigs)
+ 🚀 [Join our Slack community](https://join.slack.com/t/openimsdk/shared_invite/zt-22720d66b-o_FvKxMTGXtcnnnHiMqe9Q)
+ :eyes: [Join our wechat (微信群)](https://openim-1253691595.cos.ap-nanjing.myqcloud.com/WechatIMG20.jpeg)

## :calendar: Community Meetings

We want anyone to get involved in our community and contributing code, we offer gifts and rewards, and we welcome you to join us every Thursday night.

Our conference is in the [OpenIM Slack](https://join.slack.com/t/openimsdk/shared_invite/zt-22720d66b-o_FvKxMTGXtcnnnHiMqe9Q) 🎯, then you can search the Open-IM-Server pipeline to join

We take notes of each [biweekly meeting](https://github.com/orgs/OpenIMSDK/discussions/categories/meeting) in [GitHub discussions](https://github.com/openimsdk/open-im-server/discussions/categories/meeting), Our historical meeting notes, as well as replays of the meetings are available at [Google Docs :bookmark_tabs:](https://docs.google.com/document/d/1nx8MDpuG74NASx081JcCpxPgDITNTpIIos0DS6Vr9GU/edit?usp=sharing).

## :eyes: Who Are Using OpenIM

Check out our [user case studies](https://github.com/OpenIMSDK/community/blob/main/ADOPTERS.md) page for a list of the project users. Don't hesitate to leave a [📝comment](https://github.com/openimsdk/open-im-server/issues/379) and share your use case.

## :page_facing_up: License

OpenIM is licensed under the Apache 2.0 license. See [LICENSE](https://github.com/openimsdk/open-im-server/tree/main/LICENSE) for the full license text.

The OpenIM logo, including its variations and animated versions, displayed in this repository [OpenIM](https://github.com/openimsdk/open-im-server) under the [assets/logo](./assets/logo) and [assets/logo-gif](assets/logo-gif) directories, are protected by copyright laws.

## 🔮 Thanks to our contributors!

<a href="https://github.com/openimsdk/open-im-server/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=openimsdk/open-im-server" />
</a>
