<p align="center">
    <a href="https://www.openim.online">
        <img src="./assets/logo-gif/openim-logo.gif" width="60%" height="30%"/>
    </a>
</p>

<h3 align="center" style="border-bottom: none">
    ⭐️  Open source Instant Messaging Server ⭐️ <br>
<h3>


<p align=center>
<a href="https://goreportcard.com/report/github.com/OpenIMSDK/Open-IM-Server"><img src="https://goreportcard.com/badge/github.com/OpenIMSDK/Open-IM-Server" alt="A+"></a>
<a href="https://github.com/OpenIMSDK/Open-IM-Server/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc+label%3A%22good+first+issue%22"><img src="https://img.shields.io/github/issues/OpenIMSDK/Open-IM-Server/good%20first%20issue?logo=%22github%22" alt="good first"></a>
<a href="https://github.com/OpenIMSDK/Open-IM-Server"><img src="https://img.shields.io/github/stars/OpenIMSDK/Open-IM-Server.svg?style=flat&logo=github&colorB=deeppink&label=stars"></a>
<a href="https://join.slack.com/t/openimsdk/shared_invite/zt-1tmoj26uf-_FDy3dowVHBiGvLk9e5Xkg"><img src="https://img.shields.io/badge/Slack-300%2B-blueviolet?logo=slack&amp;logoColor=white"></a>
<a href="https://github.com/OpenIMSDK/Open-IM-Server/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-Apache--2.0-green"></a>
<a href="https://golang.org/"><img src="https://img.shields.io/badge/Language-Go-blue.svg"></a>
</p>

</p>

<p align="center">
    <a href="./README.md"><b> English </b></a> •
    <a href="./README-zh_CN.md"><b> 简体中文 </b></a> •
    <a href="https://www.openim.online/en"><b> Docs </b></a>
</p>


</p>

## ✨ About OpenIM

Open-IM-Server, meticulously developed in pure Golang, is a powerful instant messaging server. Its distinct communication method, employing pb+websocket, views every interaction as a message, streamlining customization without altering server code. Built on a microservice architecture, it offers deployment through clusters for outstanding performance and scalability.

Open-IM-Server is more than an instant messaging server; it's a powerful tool for incorporating real-time networking into your applications, positioning itself as your premier choice for integration! 🚀

Please be aware that Open-IM-Server does not function as a standalone product and does not offer built-in account registration or login services. To ease your implementation process, we've open-sourced the [chat repository](https://github.com/OpenIMSDK/chat), which comprises these features. Deploying this chat business server in conjunction with Open-IM-Server expedites the setup of a comprehensive chat product. 👥

Further enhancing your experience, we also provide an SDK client, wherein most complex logics are implemented. The [SDK repository](https://github.com/OpenIMSDK/openim-sdk-core) can be found at [this link](https://github.com/OpenIMSDK/openim-sdk-core). The [chat repository](https://github.com/OpenIMSDK/chat) is our business server while the 'core' represents the high-level encapsulation of the SDK, synergistically working together to deliver superior results. ✨

## :star2: Why OpenIM

**🔍 Function screenshot display**

<div align="center">

|            💻🔄📱 Multi Terminal Synchronization 🔄🖥️             |                  📅⚡ Efficient Meetings 🚀💼      |
| :----------------------------------------------------------: | :---------------------------------------------------------: |
| ![multiple-message](./assets/demo/multi-terminal-synchronization.png) | ![efficient-meetings](./assets/demo/efficient-meetings.png) |
|             📲🔄 **One-to-one and Group Chats** 👥🗣️             |        🎁💻 **Special Features - Custom Messages** ✉️🎨|
|         ![group-chat](./assets/demo/group-chat.png)          |   ![special-function](./assets/demo/special-function.png)      |

</div>

1. **Comprehensive Message Type Support :speech_balloon:**

   ✅ Supports almost all types of messages, including text, images, emojis, voice, video, geographical location, files, quotes, business cards, system notifications, custom messages and more

   ✅ Supports one-on-one and multi-person audio and video calls

   ✅ Provides terminal support for multiple platforms such as iOS, Android, Flutter, uni-app, ReactNative, Electron, Web, H5 

2. **Efficient Meetings Anytime, Anywhere :earth_americas:**

   ✅ Based on IM (Instant Messaging) with 100% reliable forced signaling capabilities, it paves the way for IM systems, deeply integrated with chat applications 

   ✅ Supports hundreds of people in a single meeting, with subscription numbers reaching thousands, and server-side audio and video recording 

3. **One-on-one and Group Chats for Various Social Scenarios :busts_in_silhouette:**

   ✅ OpenIM has four roles: application administrator, group owner, group administrator, and regular member 

   ✅ Powerful group features such as muting, group announcements, group validation, unlimited group members, and loading group messages as needed 

4. **Unique Features :star2:**

   ✅ Supports read-and-burn private chats, customizable duration

   ✅ Message editing function broadens social scenarios, making instant communication more diverse and interesting

5. **Open Source :open_hands:**

   ✅ The code of OpenIM is open source, self-controlled data, aimed at building a globally leading IM open source community, including client SDK and server 

   ✅ Based on open source Server, many excellent open source projects have been developed, such as [OpenKF](https://github.com/OpenIMSDK/OpenKF) (Open source AI customer service system) 

6. **Easy to Expand :wrench:**

   ✅ The OpenIM server is implemented in Golang, introducing an innovative "everything is a message" communication model, simplifying the implementation of custom messages and extended features 

7. **High Performance :racing_car:**

   ✅ OpenIM supports a hierarchical governance architecture in the cluster, tested by a large number of users, and abstracts the storage model of online messages, offline messages, and historical messages 

8. **Full Platform Support :tv:**

   ✅ Supports native iOS, Android; cross-platform Flutter, uni-app, ReactNative; major web front-end frameworks such as React, Vue; applets; and PC platforms supported by Electron

9. **The ultimate deployment experience 🤖** 

   ✅  Supports [cluster deployment](https://github.com/OpenIMSDK/Open-IM-Server/edit/main/deployments/README.md)

   ✅  Supports multi-architecture mirroring, our Docker images are hosted not only on GitHub but also on Alibaba Cloud and Docker Hub supporting multiple architectures. Visit [our GitHub packages](https://github.com/orgs/OpenIMSDK/packages?repo_name=Open-IM-Server) and read our [version management document](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/conversions/version.md) for more information.

11. **A large ecosystem of open source communities 🤲**

    ✅  We have tens of thousands of users and many solutions to problems.

    ✅  We have a large open source community called [OpenIMSDK](https://github.com/OpenIMSDK) that runs the core modules, we have an open source community called [openim-sigs](https://github.com/openim-sigs) to explore more IM-based infrastructure products.

## :rocket: Quick Start

You can quickly learn OpenIM engineering solutions,  all it takes is one simple command:

```bash
$ make demo
```

🤲 In order to facilitate the user experience, we have provided a variety of deployment solutions, you can choose your own deployment method according to the list below:

<details>   <summary>Deploying with Docker Compose</summary>


> docker compose will not be maintained in future versions, but it is still the easiest and most convenient way to organize docker compose deployments into a separate project https://github.com/openim-sigs/openim-docker to maintain.

**1. Clone the project**


```bash
# choose what you need, We take branch 3.2 as an example
$ BRANCH=release-v3.2
$ git clone -b $BRANCH https://github.com/OpenIMSDK/Open-IM-Server openim/openim-server && export openim_dir=$(pwd)/openim && cd ${openim_dir}/openim-server
```

> **Note**
> If you don't know OpenIM's versioning policy, 📚Read our release policy: https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/conversions/version.md



**2. Configure the config file**

If you tried to get started quickly with `make demo`, then you know that our config file is generated by automation.

You can use `make init` to quickly initialize a configuration file

```bash
$ make init
$ git diff
```

Then feel free to modify your current config file,  you can also modify `/scripts/install/environment.sh` document template, `make init` is essentially rendering `environment.sh` template, and then through the `make init` to automatically generate a new configuration.

If you only need to change the config file for a short time, or if you don't want to make any major changes in the future, you can modify the `.env file directly

```bash
USER=root 								#no need to modify
PASSWORD=openIM123  					#A combination of 8 or more numbers and letters, this password applies to redis, mysql, mongo, as well as accessSecret in config/config.yaml
ENDPOINT=http://127.0.0.1:10005 		#minio's external service IP and port, or use the domain name storage.xx.xx, the app must be able to access this IP and port or domain,
API_URL=http://127.0.0.1:10002/object/ 	#the app must be able to access this IP and port or domain,
DATA_DIR=./  							#designate large disk directory
```

**3. Deploy and start**

> **Note**
>
> You can deploy either directly with `make install` or with `docker compose up`, the logic is the same

```bash
$ make install
# OR
$ docker-compose up
```

4. Check the service

```bash
$ make check
```

Looking at the command line at this point, there are two items in the output, checking for the start of the component port that OpenIM depends on, and the start of the OpenIM core component

</details> 

<details>  <summary>Compile from Source</summary>


Ur need `Go 1.18` or higher version, and `make`.


```bash
go version && make --version || echo "Error: One of the commands failed."
```

Version Details: https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/conversions/version.md

```bash
# choose what you need
$ BRANCH=release-v3.1
$ git clone -b $BRANCH https://github.com/OpenIMSDK/Open-IM-Server openim && export openim=$(pwd)/openim && cd $openim && make start
```

> `make help` to help you see the instructions supported by OpenIM.

Use `make check` to check all component starts

```bash
$ make check
```

You can use the `make help-all` see OpenIM in action.

</details>

<details>  <summary>Component Configuration Instructions</summary>
The `config/config.yaml` file has detailed configuration instructions for the storage components.


The config file is available via [environment.sh](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/scripts/install/environment.sh) configuration [openim.yaml](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/deployments/templates/openim.yaml) template, and then through the `make init` to automatically generate a new configuration.

- Zookeeper

  - Used for RPC service discovery and registration, cluster support.

    ```bash
    zookeeper:
      schema: openim                          #Not recommended to modify
      address: [ 127.0.0.1:2181 ]             #address
      username:                               #username
      password:                               #password
    ```

- MySQL

  - Used for storing users, relationships, and groups, supports master-slave database.

    ```bash
    mysql:
      address: [ 127.0.0.1:13306 ]            #address
      username: root                          #username
      password: openIM123                     #password
      database: openIM_v2                     #Not recommended to modify
      maxOpenConn: 1000                       #maximum connection
      maxIdleConn: 100                        #maximum idle connection
      maxLifeTime: 60                         #maximum time a connection can be reused (seconds)
      logLevel: 4                             #log level 1=slient 2=error 3=warn 4=info
      slowThreshold: 500                      #slow statement threshold (milliseconds)
    ```

- Mongo

  - Used for storing offline messages, supports mongo sharded clusters.

    ```bash
    mongo:
      uri:                                    #Use this value directly if not empty
      address: [ 127.0.0.1:37017 ]            #address
      database: openIM                        #default mongo db
      username: root                          #username
      password: openIM123                     #password
      maxPoolSize: 100                        #maximum connections
    ```

- Redis

  - Used for storing message sequence numbers, latest messages, user tokens, and mysql cache, supports cluster deployment.

    ```bash
    redis:
      address: [ 127.0.0.1:16379 ]            #address
      username:                               #username
      password: openIM123                     #password
    ```

- Kafka

  - Used for message queues, for message decoupling, supports cluster deployment.

    ```bash
    kafka:
      username:                               #username
      password:                               #password
      addr: [ 127.0.0.1:9092 ]                #address
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

</details>

<details>  <summary>Deployed with kubernetes</summary>


read: https://github.com/OpenIMSDK/Open-IM-Server/blob/main/deployments/README.md

</details> 

<details>  <summary>Open IM Ports</summary>


| TCP Port  | Description                                                  | Operation                                             |
| --------- | ------------------------------------------------------------ | ----------------------------------------------------- |
| TCP:10001 | ws protocol, message port such as message sending, pushing etc, used for client SDK | Port release or nginx reverse proxy, and firewall off |
| TCP:10002 | api port, such as user, friend, group, message interfaces.   | Port release or nginx reverse proxy, and firewall off |
| TCP:10005 | Required when choosing minio storage (openIM uses minio storage by default) | Port release or nginx reverse proxy, and firewall off |

</details> 

<details>  <summary>Open Chat Ports</summary>


+ chat warehouse: https://github.com/OpenIMSDK/chat 

| TCP Port  | Description                                         | Operation                                             |
| --------- | --------------------------------------------------- | ----------------------------------------------------- |
| TCP:10008 | Business system, such as registration, login etc    | Port release or nginx reverse proxy, and firewall off |
| TCP:10009 | Management backend, such as statistics, banning etc | Port release or nginx reverse proxy, and firewall off |

</details>


## :link: Relationship Between APP and OpenIM

OpenIM isn't just an open-source instant messaging component, it's an integral part of your application ecosystem. Check out this diagram to understand how AppServer, AppClient, Open-IM-Server, and Open-IM-SDK interact.

![App-OpenIM Relationship](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/images/open-im-server.png)

## :building_construction: Overall Architecture

Delve into the heart of Open-IM-Server's functionality with our architecture diagram.

![Overall Architecture](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/images/Architecture.jpg)

## :hammer_and_wrench: To start developing OpenIM

OpenIM Our goal is to build a top-level open source community. We have a set of standards, in the [Community repository](https://github.com/OpenIMSDK/community).

If you'd like to contribute to this Open-IM-Server repository, please read our [contributor documentation](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/CONTRIBUTING.md).

Before you start, please make sure your changes are in demand. The best for that is to create a [new discussion](https://github.com/OpenIMSDK/Open-IM-Server/discussions/new/choose) OR [Slack Communication](https://join.slack.com/t/openimsdk/shared_invite/zt-1tmoj26uf-_FDy3dowVHBiGvLk9e5Xkg), or if you find an issue, [report it](https://github.com/OpenIMSDK/Open-IM-Server/issues/new/choose) first.

- [Code Standards](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/conversions/go_code.md)
- [Docker Images Standards](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/conversions/images.md)
- [Directory Standards](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/conversions/directory.md)
- [Commit Standards](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/conversions/commit.md)
- [Versioning Standards](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/conversions/version.md)
- [Interface Standards](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/conversions/api.md)
- [Log Standards](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/conversions/log.md)
- [Error Code Standards](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/conversions/error_code.md)

## :busts_in_silhouette: Community

+ 📚 [OpenIM Community](https://github.com/OpenIMSDK/community)
+ 💕 [OpenIM Interest Group](https://github.com/Openim-sigs)
+ 🚀 [Join our Slack community](https://join.slack.com/t/openimsdk/shared_invite/zt-1tmoj26uf-_FDy3dowVHBiGvLk9e5Xkg)
+ :eyes: [Join our wechat (微信群)](https://openim-1253691595.cos.ap-nanjing.myqcloud.com/WechatIMG20.jpeg)

## :calendar: Community Meetings

We want anyone to get involved in our community and contributing code, we offer gifts and rewards, and we welcome you to join us every Thursday night.

Our conference is in the [OpenIM Slack](https://join.slack.com/t/openimsdk/shared_invite/zt-1tmoj26uf-_FDy3dowVHBiGvLk9e5Xkg) 🎯, then you can search the Open-IM-Server pipeline to join

We take notes of each [biweekly meeting](https://github.com/orgs/OpenIMSDK/discussions/categories/meeting) in [GitHub discussions](https://github.com/OpenIMSDK/Open-IM-Server/discussions/categories/meeting), Our historical meeting notes, as well as replays of the meetings are available at [Google Docs :bookmark_tabs:](https://docs.google.com/document/d/1nx8MDpuG74NASx081JcCpxPgDITNTpIIos0DS6Vr9GU/edit?usp=sharing).

## :eyes: Who are using OpenIM

Check out our [user case studies](https://github.com/OpenIMSDK/community/blob/main/ADOPTERS.md) page for a list of the project users. Don't hesitate to leave a [📝comment](https://github.com/OpenIMSDK/Open-IM-Server/issues/379) and share your use case.

## :page_facing_up: License

OpenIM is licensed under the Apache 2.0 license. See [LICENSE](https://github.com/OpenIMSDK/Open-IM-Server/tree/main/LICENSE) for the full license text.

The OpenIM logo, including its variations and animated versions, displayed in this repository [OpenIM](https://github.com/OpenIMSDK/Open-IM-Server) under the [assets/logo](./assets/logo) and [assets/logo-gif](assets/logo-gif) directories, are protected by copyright laws.

## 🔮 Thanks to our contributors!

<a href="https://github.com/OpenIMSDK/Open-IM-Server/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=OpenIMSDK/Open-IM-Server" />
</a>
