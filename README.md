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

## What is Open-IM-Server

Instant messaging server. Backend in pure Golang, wire transport protocol is JSON over websocket.
Everything is a message in Open-IM-Server, so you can extend custom messages easily, there is no need to modify the server code.
Using microservice architectures, Open-IM-Server can be deployed using clusters.
By deployment of the Open-IM-Server on the customer's server, developers can integrate instant messaging and real-time network capabilities into their own applications free of charge and quickly, and ensure the security and privacy of business data.

## Features

- Everything in Free
- Scalable architecture
- Easy integration
- Good scalability
- High performance
- Lightweight
- Supports multiple protocols

## Community
- Visit the Chinese official website here: [📚 Open-IM docs](https://www.openim.online/zh)

## Quick start

### Installing Open-IM-Server

> Open-IM relies on five open source high-performance components: ETCD, MySQL, MongoDB, Redis, and Kafka. Privatization deployment Before Open-IM-Server, please make sure that the above five components have been installed. If your server does not have the above components, you must first install Missing components. If you have the above components, it is recommended to use them directly. If not, it is recommended to use Docker-compose, no To install dependencies, one-click deployment, faster and more convenient.

#### Deploy using Docker

1. Install [Go environment](https://golang.org/doc/install). Make sure Go version is at least 1.17

2. Clone the Open-IM project to your server

   ```
   git clone https://github.com/OpenIMSDK/Open-IM-Server.git --recursive
   ```

3. Deploy

    1. Modify env

       ```
       #cd Open-IM-server
       USER=root  
       PASSWORD=openIM123    #Password with more than 8 digits, excluding special characters
       ENDPOINT=http://127.0.0.1:10005 #Replace 127.0.0.1 with Internet IP
       DATA_DIR=./ 
       ```

    2. Deploy && Start
    
       ```
       chmod +x install_im_server.sh;
       ./install_im_server.sh;
       ```
    
    4. Check service
    
       ```
       cd script;
       ./docker_check_service.sh./check_all.sh
       ```
       
       ![OpenIMServersonSystempng](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/Open-IM-Servers-on-System.png)

#### Deploy using source code

1. Go 1.17 or above.
2. Clone

```shell
git clone https://github.com/OpenIMSDK/Open-IM-Server.git --recursive 
cd cmd/Open-IM-SDK-Core
git checkout main
```

1. Set executable permissions

```shell
cd ../../script/
chmod +x *.sh
```

1. build

```shell
./batch_build_all_service.sh
```

all services build success

### CONFIGURATION INSTRUCTIONS

> Open-IM configuration is divided into basic component configuration and business internal service configuration. Developers need to fill in the address of each component as the address of their server component when using the product, and ensure that the internal service port of the business is not occupied

#### Basic Component Configuration Instructions

- ETCD
    - Etcd is used for the discovery and registration of rpc services, etcd Schema is the prefix of the registered name, it is recommended to modify it to your company name, etcd address (ip+port) supports clustered deployment, you can fill in multiple ETCD addresses separated by commas, and also only one etcd address.
- MySQL
    - mysql is used for full storage of messages and user relationships. Cluster deployment is not supported for the time being. Modify addresses and users, passwords, and database names.
- Mongo
    - Mongo is used for offline storage of messages. The default storage is 7 days. Cluster deployment is temporarily not supported. Just modify the address and database name.
- Redis
    - Redis is currently mainly used for message serial number storage and user token information storage. Cluster deployment is temporarily not supported. Just modify the corresponding redis address and password.
- Kafka
    - Kafka is used as a message transfer storage queue to support cluster deployment, just modify the corresponding address

#### Internal Service Configuration Instructions

- credential&&push
    - The Open-IM needs to use the three-party offline push function. Currently, Tencent's three-party push is used. It supports IOS, Android and OSX push. This information is some registration information pushed by Tencent. Developers need to go to Tencent Cloud Mobile Push to register the corresponding information. If you do not fill in the corresponding information, you cannot use the offline message push function
- api&&rpcport&&longconnsvr&&rpcregistername
    - The api port is the http interface, longconnsvr is the websocket listening port, and rpcport is the internal service startup port. Both support cluster deployment. Make sure that these ports are not used. If you want to open multiple services for a single service, fill in multiple ports separated by commas. rpcregistername is the service name registered by each service to the registry etcd, no need to modify
- log&&modulename
    - The log configuration includes the storage path of the log file, and the log is sent to elasticsearch for log viewing. Currently, the log is not supported to be sent to elasticsearch. The configuration does not need to be modified for the time being. The modulename is used to split the log according to the name of the service module. The default configuration is fine.
- multiloginpolicy&&tokenpolicy
    - Open-IM supports multi-terminal login. Currently, there are three multi-terminal login policies. The PC terminal and the mobile terminal are online at the same time by default. When multiple policies are configured to be true, the first policy with true is used by default, and the token policy is the generated token policy. , The developer can customize the expiration time of the token

### SCRIPT DESCRIPTION

> Open-IM script provides service compilation, start, and stop script. There are four Open-IM script start modules, one is the http+rpc service start module, the second is the websocket service start module, then the msg_transfer module, and the last is the push module

- path_info.cfg&&style_info.cfg&&functions.sh
    - Contains the path information of each module, including the path where the source code is located, the name of the service startup, the shell print font style, and some functions for processing shell strings
- build_all_service.sh
    - Compile the module, compile all the source code of Open-IM into a binary file and put it into the bin directory
- start_rpc_api_service.sh&&msg_gateway_start.sh&&msg_transfer_start.sh&&push_start.sh
    - Independent script startup module, followed by api and rpc modules, message gateway module, message transfer module, and push module
- start_all.sh&&stop_all.sh
    - Total script, start all services and close all services

## Authentication Clow Chart 

![avatar](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/open-im-server.png)

## Architecture

![avatar](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/Architecture.jpg)

## To start developing OpenIM
The [community repository](https://github.com/OpenIMSDK/community) hosts all information about building Kubernetes from source, how to contribute code and documentation, who to contact about what, etc.


## Contributing

Contributions to this project are welcome! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for details.

## Community Meetings
We want anyone to get involved in our community, we offer gifts and rewards, and we welcome you to join us every Thursday night.

We take notes of each [biweekly meeting](https://github.com/OpenIMSDK/Open-IM-Server/issues/381) in [GitHub discussions](https://github.com/OpenIMSDK/Open-IM-Server/discussions/categories/meeting), and our minutes are written in [Google Docs](https://docs.google.com/document/d/1nx8MDpuG74NASx081JcCpxPgDITNTpIIos0DS6Vr9GU/edit?usp=sharing).


## Who are using Open-IM-Server
The [user case studies](https://github.com/OpenIMSDK/community/blob/main/ADOPTERS.md) page includes the user list of the project. You can leave a [📝comment](https://github.com/OpenIMSDK/Open-IM-Server/issues/379) to let us know your use case.

## ⚠️ Announcement

Dear users of the OpenIM repository, we are pleased to announce that we are currently undergoing a major overhaul to improve our service quality and user experience. We will be using the [errcode](https://github.com/OpenIMSDK/Open-IM-Server/tree/errcode) branch to make extensive updates and improvements to the main branch, ensuring that our code repository is in optimal condition.

During this time, we will need to pause PR and issue handling on the main branch. We understand that this may cause some inconvenience for you, but we believe it will provide us with better service and a more reliable code repository. If you need to submit code during this period, please submit it to the [errcode](https://github.com/OpenIMSDK/Open-IM-Server/tree/errcode) branch, and we will process your request as soon as possible.

We appreciate your support and trust, as well as your patience and understanding throughout this process. We value your contributions and suggestions, which are the driving force behind our continuous improvement and growth.

We anticipate that this work will be completed soon, and we will do our utmost to minimize any impact on you. Once again, we express our heartfelt thanks and apologies to you.

Thank you!

![avatar](https://github.com/OpenIMSDK/OpenIM-Docs/blob/main/docs/images/WechatIMG20.jpeg)

## License

Open-IM-Server is under the Apache 2.0 license. See the [LICENSE](./LICENSE) file for details
