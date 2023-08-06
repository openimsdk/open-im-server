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
    <a href="./README-zh_CN.md"><b> 中文 </b></a>
</p>

</p>

## What is Open-IM-Server

Open-IM-Server is a robust instant messaging server engineered using pure Golang, leveraging JSON over WebSocket for communication. The server treats everything as a message, facilitating straightforward customization without modifying the server code. Its microservice architecture enables deployment using clusters, ensuring high performance and scalability.

Whether you're looking to integrate instant messaging or real-time networking into your applications, Open-IM-Server is your go-to solution! :rocket:

It's important to note that Open-IM-Server isn't a standalone product, and it doesn't include account registration and login services. However, we've made your life easier by open-sourcing the [chat repository](https://github.com/OpenIMSDK/chat) that includes login and registration features. Deploying the chat business server alongside Open-IM-Server quickly sets up a comprehensive chat product. :busts_in_silhouette:

## :star2: Why OpenIM

1. Comprehensive Message Type Support :speech_balloon:
   ✅  Supports almost all types of messages, including text, images, emojis, voice, video, geographical location, files, quotes, business cards, system notifications, custom messages and more :mailbox_with_mail:

   ✅ Supports one-on-one and multi-person audio and video calls :telephone_receiver:

   ✅ Provides terminal support for multiple platforms such as iOS, Android, Flutter, uni-app, ReactNative, Electron, Web, H5 :iphone:

2. Efficient Meetings Anytime, Anywhere :earth_americas:

   ✅ Based on IM (Instant Messaging) with 100% reliable forced signaling capabilities, it paves the way for IM systems, deeply integrated with chat applications :link:

   ✅ Supports hundreds of people in a single meeting, with subscription numbers reaching thousands, and server-side audio and video recording :video_camera:

3. One-on-one and Group Chats for Various Social Scenarios :busts_in_silhouette:

   ✅ OpenIM has four roles: application administrator, group owner, group administrator, and regular member :man_teacher:

   ✅ Powerful group features such as muting, group announcements, group validation, unlimited group members, and loading group messages as needed :loudspeaker:

4. Unique Features :star2:

   ✅ Supports read-and-burn private chats, customizable duration :fire:

   ✅ Message editing function broadens social scenarios, making instant communication more diverse and interesting :pencil2:

5. Open Source :open_hands:

   ✅ The code of OpenIM is open source, self-controlled data, aimed at building a globally leading IM open source community, including client SDK and server :globe_with_meridians:

   ✅ Based on open source Server, many excellent open source projects have been developed, such as [OpenKF](https://github.com/OpenIMSDK/OpenKF) (Open source AI customer service system) ✨ 

6. Easy to Expand :wrench:

   ✅ The OpenIM server is implemented in Golang, introducing an innovative "everything is a message" communication model, simplifying the implementation of custom messages and extended features :computer:

7. High Performance :racing_car:

   ✅ OpenIM supports a hierarchical governance architecture in the cluster, tested by a large number of users, and abstracts the storage model of online messages, offline messages, and historical messages :rocket:

8. Full Platform Support :tv:

   ✅ Supports native iOS, Android; cross-platform Flutter, uni-app, ReactNative; major web front-end frameworks such as React, Vue; applets; and PC platforms supported by Electron :desktop_computer:

9. The ultimate deployment experience 🤖 

   ✅  Supports cluster deployment

   ✅  Supports multi-architecture mirroring

10. A large ecosystem of open source communities 🤲


## :busts_in_silhouette: Community

Explore our [OpenIM Developer Documentation](https://www.openim.online/) for more details.

## :rocket: Quick Start

<details>   <summary>Deploying with Docker Compose</summary>

1. Clone the project

```bash
# choose what you need
BRANCH=release-v3.1
git clone -b $BRANCH https://github.com/OpenIMSDK/Open-IM-Server openim && export openim=$(pwd)/openim && cd $openim && make build
```

> **Note**
>
> Read our release policy: https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/conversions/version.md

2. Modify .env

```bash
USER=root #no need to modify
PASSWORD=openIM123  #A combination of 8 or more numbers and letters, this password applies to redis, mysql, mongo, as well as accessSecret in config/config.yaml
ENDPOINT=http://127.0.0.1:10005 #minio's external service IP and port, or use the domain name storage.xx.xx, the app must be able to access this IP and port or domain,
API_URL=http://127.0.0.1:10002/object/ #the app must be able to access this IP and port or domain,
DATA_DIR=./  #designate large disk directory
```

3. Deploy and start

> **Note**
> This command can only be executed once. It will modify the component passwords in docker-compose based on the `PASSWORD` variable in `.env`, and modify the component passwords in `config/config.yaml`. If the password in `.env` changes, you need to first execute `docker-compose down`; `rm components -rf` and then execute this command.

```
make install
```

4. Check the service

```bash
make check
```

![https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/images/docker_build.png](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/images/docker_build.png)

</details> 

<details>  <summary>Compile from Source</summary>

Ur need `Go 1.18` or higher version, and `make`.

Version Details: https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/conversions/version.md

```bash
# choose what you need
BRANCH=release-v3.1
git clone -b $BRANCH https://github.com/OpenIMSDK/Open-IM-Server openim && export openim=$(pwd)/openim && cd $openim && make build
```

Read about the [OpenIM Version Policy](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/conversions/version.md)

`make help` to help you see the instructions supported by OpenIM.

All services have been successfully built as shown in the figure

![Successful Compilation](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/images/build.png)

</details>

<details>  <summary>Component Configuration Instructions</summary>

The config/config.yaml file has detailed configuration instructions for the storage components.

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

<details>  <summary>Start and Stop Services</summary>

Start services

```
./scripts/start_all.sh;
```

Check services

```
./scripts/check_all.sh
```

Stop services

```
./scripts/stop_all.sh
```

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

The [community repository](https://github.com/OpenIMSDK/community) hosts all the information you need to start contributing to OpenIM.

## :heart: Contributing

We welcome all contributions to the project! For more details, please see [CONTRIBUTING.md](./CONTRIBUTING.md).

## :calendar: Community Meetings

We love community involvement! Join us for our biweekly meetings every Thursday night. We even offer gifts and rewards! :gift:

You can find the meeting minutes on [GitHub discussions](https://github.com/OpenIMSDK/Open-IM-Server/discussions/categories/meeting) and our [Google Docs](https://docs.google.com/document/d/1nx8MDpuG74NASx081JcCpxPgDITNTpIIos0DS6Vr9GU/edit?usp=sharing).

## :eyes: Who are using Open-IM-Server

Check out our [user case studies](https://github.com/OpenIMSDK/community/blob/main/ADOPTERS.md) page for a list of the project users. Don't hesitate to leave a [📝comment](https://github.com/OpenIMSDK/Open-IM-Server/issues/379) and share your use case.

![avatar](https://openim-1253691595.cos.ap-nanjing.myqcloud.com/WechatIMG20.jpeg)

## :page_facing_up: License

Open-IM-Server is licensed under the Apache 2.0 license. See the [LICENSE](./LICENSE) file for details.

## :cloud: Docker Images

Our Docker images are hosted not only on GitHub but also on Alibaba Cloud and Docker Hub supporting multiple architectures. Visit [our GitHub packages](https://github.com/orgs/OpenIMSDK/packages?repo_name=Open-IM-Server) and read our [version management document](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/conversions/version.md) for more information.
