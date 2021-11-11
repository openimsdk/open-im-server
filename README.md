# Open-IM-Server 

![avatar](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/open-im-logo.png)

[![LICENSE](https://img.shields.io/badge/license-Apache--2.0-green)](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/LICENSE) [![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)

## Open-IM-Server: Open source Instant Messaging Server

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

- Join the Telegram-OpenIM group: https://t.me/joinchat/zSJLPaHBNLZmODI1
- 中文官网访问这里：[Open-IM中文开发文档](https://doc.rentsoft.cn/)

## Quick start

### Installing Open-IM-Server

> Open-IM relies on five open source high-performance components: ETCD, MySQL, MongoDB, Redis, and Kafka. Privatization deployment Before Open-IM-Server, please make sure that the above five components have been installed. If your server does not have the above components, you must first install Missing components. If you have the above components, it is recommended to use them directly. If not, it is recommended to use Docker-compose, no To install dependencies, one-click deployment, faster and more convenient.

#### Source code deployment

1. Install [Go environment](https://golang.org/doc/install). Make sure Go version is at least 1.15.

2. Clone the Open-IM project to your server.

   ```
   git clone https://github.com/OpenIMSDK/Open-IM-Server.git --recursive
   ```

3. Build and start Service.

    1. Shell authorization

       ```
       #cd Open-IM-server/script
       
       chmod +x *.sh
       ```

    2. Execute the build shell

       ```
       ./build_all_service.sh
       ```

    3. Start service

       ```
       ./start_all.sh
       ```

    4. Check service

       ```
       ./check_all.sh
       ```

       ![OpenIMServersonSystempng](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/Open-IM-Servers-on-System.png)

#### Docker deployment

All images are available at https://hub.docker.com/r/lyt1123/open_im_server

1. [Install Docker](https://docs.docker.com/install/) 1.13 or above.

2. [Install Docker Compose](https://docs.docker.com/compose/install/) 1.22 or above.

3. Clone the Open-IM project to your server.

   ```
   git clone https://github.com/OpenIMSDK/Open-IM-Server.git --recursive
   ```

4. Start docker-compose with one click(Docker automatically pulls all images)

   ```
   docker-compose up -d
   ```

5. Check service

   ```
   ./docker_check_service.sh 
   ```

   ![OpenIMServersondockerpng](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/Open-IM-Servers-on-docker.png)

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

> Open-IM script provides service compilation, start, and stop scripts. There are four Open-IM script start modules, one is the http+rpc service start module, the second is the websocket service start module, then the msg_transfer module, and the last is the push module

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

## License

Open-IM-Server is under the Apache 2.0 license. See the [LICENSE](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/LICENSE) file for details
