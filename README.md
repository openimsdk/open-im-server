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
<a href="https://join.slack.com/t/openimsdk/shared_invite/zt-1tmoj26uf-_FDy3dowVHBiGvLk9e5Xkg"><img src="https://img.shields.io/badge/Slack-100%2B-blueviolet?logo=slack&amp;logoColor=white"></a>
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

Open-IM-Server is an instant messaging server developed using pure Golang, adopting JSON over WebSocket as the communication protocol. In Open-IM-Server, everything is a message, so you can easily extend custom messages without modifying the server code. With a microservice architecture, Open-IM-Server can be deployed using clusters. By deploying Open-IM-Server on a server, developers can quickly integrate instant messaging and real-time networking features into their applications, ensuring the security and privacy of business data.

Open-IM-Server is not a standalone product and does not include account registration and login services. For your convenience, we have open-sourced the [chat repository](https://github.com/OpenIMSDK/chat) which includes login and registration functionality. By deploying the chat business server alongside Open-IM-Server, a chat product can be set up.

## Features

- Open source
- Easy to integrate
- Excellent scalability
- High performance
- Lightweight
- Supports multiple protocols

## Community
- Visit the official website: [OpenIM  Developer Documentation](https://www.openim.online/)

## Quick Start

### Deploying with docker-compose

1. Clone the project

```
clone https://github.com/OpenIMSDK/Open-IM-Server 
cd Open-IM-Server
git checkout release-v3.0 #or other release branch
```

1. Modify .env

```
USER=root #no need to modify
PASSWORD=openIM123  #A combination of 8 or more numbers and letters, this password applies to redis, mysql, mongo, as well as accessSecret in config/config.yaml
ENDPOINT=http://127.0.0.1:10005 #minio's external service IP and port, or use the domain name storage.xx.xx, the app must be able to access this IP and port or domain,
API_URL=http://127.0.0.1:10002/object/ #the app must be able to access this IP and port or domain,
DATA_DIR=./  #designate large disk directory
```

1. Deploy and start

> **Note**: This command can only be executed once. It will modify the component passwords in docker-compose based on the PASSWORD variable in .env, and modify the component passwords in config/config.yaml. If the password in .env changes, you need to first execute docker-compose down; rm components -rf and then execute this command.

```
chmod +x install_im_server.sh;
./install_im_server.sh;
```

1. Check the service

```
cd scripts;
./docker_check_service.sh
```

![https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/images/docker_build.png](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/images/docker_build.png)



### Compile from source

1. Go 1.18 or higher version.

2. Clone

   ```
   git clone https://github.com/OpenIMSDK/Open-IM-Server 
   cd Open-IM-Server
   git checkout release-v3.0 #or other release branch
   ```

3. Compile

   ```
   cd Open-IM-server/scripts
   chmod +x *.sh
   ./build_all_service.sh
   ```

All services have been successfully built as shown in the figure

![Successful Compilation](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/images/build.png)

### Component Configuration Instructions

The config/config.yaml file has detailed configuration instructions for the storage components.

- Zookeeper

  - Used for RPC service discovery and registration, cluster support.

    ```
    zookeeper:
      schema: openim                          #Not recommended to modify
      address: [ 127.0.0.1:2181 ]             #address
      username:                               #username
      password:                               #password
    ```

- MySQL

  - Used for storing users, relationships, and groups, supports master-slave database.

    ```
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

    ```
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

    ```
    redis:
      address: [ 127.0.0.1:16379 ]            #address
      username:                               #username
      password: openIM123                     #password
    ```

- Kafka

  - Used for message queues, for message decoupling, supports cluster deployment.

    ```
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

### Start and Stop Services

Start services

```
./start_all.sh;
```

Check services

```
./check_all.sh
```

Stop services

```
./stop_all.sh
```

### Open IM Ports

| TCP Port  | Description                                                  | Operation                                             |
| --------- | ------------------------------------------------------------ | ----------------------------------------------------- |
| TCP:10001 | ws protocol, message port such as message sending, pushing etc, used for client SDK | Port release or nginx reverse proxy, and firewall off |
| TCP:10002 | api port, such as user, friend, group, message interfaces.   | Port release or nginx reverse proxy, and firewall off |
| TCP:10005 | Required when choosing minio storage (openIM uses minio storage by default) | Port release or nginx reverse proxy, and firewall off |

### Open Chat Ports

| TCP Port  | Description                                         | Operation                                             |
| --------- | --------------------------------------------------- | ----------------------------------------------------- |
| TCP:10008 | Business system, such as registration, login etc    | Port release or nginx reverse proxy, and firewall off |
| TCP:10009 | Management backend, such as statistics, banning etc | Port release or nginx reverse proxy, and firewall off |

## Relationship Between APP and OpenIM

OpenIM is an open source instant messaging component, it is not an independent product. This image shows the relationship between AppServer, AppClient, Open-IM-Server and Open-IM-SDK.

![https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/open-im-server.png](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/images/open-im-server.png)

## Overall Architecture

![https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/Architecture.jpg](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/images/Architecture.jpg)

## To start developing OpenIM
The [community repository](https://github.com/OpenIMSDK/community) hosts all information about building Kubernetes from source, how to contribute code and documentation, who to contact about what, etc.


## Contributing

Contributions to this project are welcome! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for details.

## Community Meetings
We want anyone to get involved in our community, we offer gifts and rewards, and we welcome you to join us every Thursday night.

We take notes of each [biweekly meeting](https://github.com/OpenIMSDK/Open-IM-Server/issues/381) in [GitHub discussions](https://github.com/OpenIMSDK/Open-IM-Server/discussions/categories/meeting), and our minutes are written in [Google Docs](https://docs.google.com/document/d/1nx8MDpuG74NASx081JcCpxPgDITNTpIIos0DS6Vr9GU/edit?usp=sharing).


## Who are using Open-IM-Server
The [user case studies](https://github.com/OpenIMSDK/community/blob/main/ADOPTERS.md) page includes the user list of the project. You can leave a [📝comment](https://github.com/OpenIMSDK/Open-IM-Server/issues/379) to let us know your use case.

![avatar](https://openim-1253691595.cos.ap-nanjing.myqcloud.com/WechatIMG20.jpeg)

## License

Open-IM-Server is under the Apache 2.0 license. See the [LICENSE](./LICENSE) file for details
