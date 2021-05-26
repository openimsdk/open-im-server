# Open-IM-Server
![avatar](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/Open-IM.png)

[![LICENSE](https://img.shields.io/badge/license-Apache--2.0-green)](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/LICENSE)
[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)

## Open-IM-Server: Open source  Instant Messaging Server

Instant messaging server. Backend in pure Golang, wire transport protocol is JSON over websocket.

Everything is a message in Open-IM-Server, so you can extend custom messages easily, there is no need to modify the
server code.

Using microservice architectures, Open-IM-Server can be deployed using clusters.

By deployment of the Open-IM-Server on the customer's server, developers can integrate instant messaging and real-time
network capabilities into their own applications free of charge and quickly, and ensure the security and privacy of
business data.

## Features

* Everything in Free

* Scalable architecture

* Easy integration

* Good scalability

* High performance

* Lightweight

* Supports multiple protocols

## Community

* Join the Telegram-OpenIM group: https://t.me/joinchat/zSJLPaHBNLZmODI1
* 扫码加入微信群: https://github.com/Open-IM-IM/opim_admin/blob/main/docs/Wechat.jpg

## Quick start

### Installing Open-IM-Server

#### Building from Source

> Open-IM relies on five open source high-performance components: **ETCD**, **MySQL**, **MongoDB**, **Redis**, **Kafka**. Before you deploy Open-IM-Server privately, please make sure that you have installed the above five components and **check the component  connection parameters ** in the configuration file. you must install the missing components first,If your server does not have the above components. **It is recommended to use it directly, if you have the above components, if not, Docker installation is recommended, which is faster and more convenient**.

1. Install [Go environment](https://golang.org/doc/install). Make sure Go version is at least 1.15.

3. Git clone  Open-IM project

   ```
   https://github.com/OpenIMSDK/Open-IM-Server.git
   ```

4. Open [config.yaml](https://github.com/Open-IM-IM/opim_admin/blob/main/config/config.yaml),then modify the following parameters.

    - Check that the ETCD connection parameters.

      ```
      etcd:
        etcdAddr: [ 127.0.0.1:2379]
      ```

    - Check or modify database(MySQL) connection parameters are correct for your database.

      ```
      mysql:
        dbAddress: [ 127.0.0.1:3306]
        dbUserName: xxx
        dbPassword: xxx
      ```

    - Check or modify database(MongoDB) connection parameters are correct for your database.

      ```
      mongo:
      dbAddress: [ 127.0.0.1:27017 ]
        dbUserName:
        dbPassword:
      ```

    - Check or modify Redis connection parameters.

      ```
      redis:
        dbAddress: [ 127.0.0.1:6379]
        dbPassWord: 
      ```

    - Check or modify Kafka connection parameters.

      ```
      kafka:
        ws2mschat:
          addr: [ 127.0.0.1:9092 ]
        ms2pschat:
          addr: [ 127.0.0.1:9092 ]
      ```

5. Build Open-IM server and database initializer:

    - **MySQL**

      ```
      need to add
      ```

    - **MongoDB**

      ```
      need to add
      ```

6. Enter the script directory and execute the script according to the steps。

    1. Shell authorization

       ```
       chmod +x *.sh
       ```

    2. Execute build shell

       ```
       ./build_all_service.sh
       ```

    3. Start service

       ```
       ./start_all.sh
       ```

#### Using Docker to run Open-IM-Server

> Open-IM relies on five open source high-performance components: **ETCD**, **MySQL**, **MongoDB**, **Redis**, **Kafka**. Before you deploy Open-IM-Server privately, please make sure that you have installed the above five components and **check the component  connection parameters ** in the configuration file. you must install the missing components first,If your server does not have the above components. **It is recommended to use it directly, if you have the above components, if not, Docker installation is recommended, which is faster and more convenient**.

All images are available at https://hub.docker.com/r/lyt1123/open_im_server

1. [Install Docker](https://docs.docker.com/install/) 1.13 or above.

3. Pull Open_IM_Server Image from docker

   ```
   docker pull docker.io/lyt1123/open_im_server:[tag]
   #eg
   docker pull docker.io/lyt1123/open_im_server:1.0
   ```

4. External config file,the container comes with a built-in config file which can be customized with values from the environment variables .**If changes are extensive it may be more convenient to replace the built-in config file with a custom one**. In that case map the config file located on your host.

    - Create configuration folder directory

      ```
      mkdir -p open_im_server/config
      ```

    - Download the [config.yaml](https://github.com/Open-IM-IM/opim_admin/blob/main/config/config.yaml) file from github, then modify the following parameters

        - Check or modify the Etcd connection parameters.

      ```
      etcd:
        etcdAddr: [ 127.0.0.1:2379]
      ```

        - Check or modify  database(MySQL) connection parameters are correct for your database.

      ```
      mysql:
        dbAddress: [ 127.0.0.1:3306]
        dbUserName: xxx
        dbPassword: xxx
      ```

        - Check or modify  database(MongoDB) connection parameters are correct for your database.

      ```
      mongo:
      dbAddress: [ 127.0.0.1:27017 ]
        dbUserName:
        dbPassword:
      ```

        - Check or modify  the Redis connection parameters.

      ```
      redis:
        dbAddress: [ 127.0.0.1:6379]
        dbPassWord: 
      ```

        - Check or modify  the Kafka connection parameters.

      ```
      kafka:
        ws2mschat:
          addr: [ 127.0.0.1:9092 ]
        ms2pschat:
          addr: [ 127.0.0.1:9092 ]
      ```

5. Start Open-IM-Server Service

   ```
   docker run -p 10000:10000 -p 7777:7777 --name open_im_server -v /home/open_im_server/logs:/home/open_im_server/logs -v /home/open_im_server/config/config.yaml:/home/open_im_server/config/config.yaml --restart always -d docker.io/lyt1123/open_im_server:[tag]
   ```

    - -p 10000:10000	The container port maps the host 10000 port, provides api service.
    - -p  7777:7777    The container port maps the host  7777 port, provides message services.
    - --name open_im_server   Container service name
    - -v /home/open_im_server/logs:/home/open_im_server/logs    The container log directory maps the host directory
    - -v /home/open_im_server/config/config.yaml:/home/open_im_server/config/config.yaml    The container configuration file maps the host configuration file
    - --restart always    Automatically start when the container is closed abnormally
    - -d  Running service in the background

### CONFIGURATION INSTRUCTIONS
>Open-IM configuration is divided into basic component configuration and business internal service configuration. Developers need to fill in the address of each component as the address of their server component when using the product, and ensure that the internal service port of the business is not occupied

#### Basic Component Configuration Instructions
* **ETCD**
  
    * Etcd is used for the discovery and registration of rpc services, etcd Schema is the prefix of the registered name, it is recommended to modify it to  your company name, etcd address (ip+port) supports clustered deployment, you can fill in multiple ETCD addresses separated by commas, and also only one etcd address.
* **MySQL**
    * mysql is used for full storage of messages and user relationships. Cluster deployment is not supported for the time being. Modify addresses and users, passwords, and database names.
* **Mongo**
    * Mongo is used for offline storage of messages. The default storage is 7 days. Cluster deployment is temporarily not supported. Just modify the address and database name.
* **Redis**
    * Redis is currently mainly used for message serial number storage and user token information storage. Cluster deployment is temporarily not supported. Just modify the corresponding redis address and password.
* **Kafka**
    * Kafka is used as a message transfer storage queue to support cluster deployment, just modify the corresponding address
#### Internal Service Configuration Instructions
* **credential&&push**
    * The Open-IM  needs to use the three-party offline push function. Currently, Tencent's three-party push is used. It supports IOS, Android and OSX push. This information is some registration information pushed by Tencent. Developers need to go to Tencent Cloud Mobile Push to register the corresponding information. If you do not fill in the corresponding information, you cannot use the offline message push function
* **api&&rpcport&&longconnsvr&&rpcregistername**
    * The api port is the http interface, longconnsvr is the websocket listening port, and rpcport is the internal service startup port. Both support cluster deployment. Make sure that these ports are not used. If you want to open multiple services for a single service, fill in multiple ports separated by commas. rpcregistername is the service name registered by each service to the registry etcd, no need to modify
* **log&&modulename**
    
    * The log configuration includes the storage path of the log file, and the log is sent to elasticsearch for log viewing. Currently, the log is not supported to be sent to elasticsearch. The configuration does not need to be modified for the time being. The modulename is used to split the log according to the name of the service module. The default configuration is fine.
* **multiloginpolicy&&tokenpolicy**
    * Open-IM supports multi-terminal login. Currently, there are three multi-terminal login policies. The PC terminal and the mobile terminal are online at the same time by default. When multiple policies are configured to be true, the first policy with true is used by default, and the token policy is the generated token policy. , The developer can customize the expiration time of the token

### SCRIPT DESCRIPTION
>Open-IM script provides service compilation, start, and stop scripts. There are four Open-IM script start modules, one is the http+rpc service start module, the second is the websocket service start module, then the msg_transfer module, and the last is the push module
* **path_info.cfg&&style_info.cfg&&functions.sh**
    * Contains the path information of each module, including the path where the source code is located, the name of the service startup, the shell print font style, and some functions for processing shell strings
* **build_all_service.sh**
    * Compile the module, compile all the source code of Open-IM into a binary file and put it into the bin directory
* **start_rpc_api_service.sh&&msg_gateway_start.sh&&msg_transfer_start.sh&&push_start.sh**
    * Independent script startup module, followed by api and rpc modules, message gateway module, message transfer module, and push module
* **start_all.sh&&stop_all.sh**
    * Total script, start all services and close all services

### Server-side authentication api graphic explanation of the login authentication process
   
- **User Register**
    - **Request URL**
       ```
       http://x.x.x.x:10000/auth/user_register
      ```
    - **Request method**
      ```
      POST
      ```
    - **Parameter**

      | parameter name | required | Type   | Description                                                  |
      | -------------- | -------- | ------ | ------------------------------------------------------------ |
      | secret         | Y        | string | The secret key used by the app server to connect to the sdk server. The maximum length is 32 characters. It must be ensured that the secret keys of the app server and the sdk server are the same. There is a risk of secret leakage, and it is best to save it on the user server. |
      | platform       | Y        | int    | Platform type iOS 1, Android 2, Windows 3, OSX 4, WEB 5, applet 6, linux 7 |
      | uid            | Y        | string | User ID, with a maximum length of 64 characters, must be unique within an APP |
      | name           | Y        | string | User nickname, the maximum length is 64 characters, can be set as an empty string |
      | icon           | N        | string | User avatar, the maximum length is 1024 bytes, can be set as an empty string |
      | gender         | N        | int    | User gender, 0 means unknown, 1 means male, 2 female means female, others will report parameter errors |
      | mobile         | N        | string | User mobile, the maximum length is 32 characters, non-Mainland China mobile phone numbers need to fill in the country code (such as the United States: +1-xxxxxxxxxx) or the area code (such as Hong Kong: +852-xxxxxxxx), which can be set as an empty string |
      | birth          | N        | string | The birthday of the user, the maximum length is 16 characters, can be set as an empty string |
      | email          | N        | string | User email, the maximum length is 64 characters, can be set as an empty string |
      | ex             | N        | string | User business card extension field, the maximum length is 1024 characters, users can extend it by themselves, it is recommended to encapsulate it into a JSON string, or set it to an empty string |

    - **Return Parameter**
      ```
      {
         "errCode": 0,
         "errMsg": "",
         "data":{
            "uid": "",
            "token": "",
            "expiredTime": 0,
         }
      }
      ```

- **Refresh Token**
    - **Request URL**
       ```
       http://x.x.x.x:10000/auth/user_token
      ```
    - **Request method**
      ```
      POST
      ```
    - **Parameter**

      | parameter name | required | Type   | Description                                                  |
      | -------------- | -------- | ------ | ------------------------------------------------------------ |
      | secret         | Y        | string | The secret key used by the app server to connect to the sdk server. The maximum length is 32 characters. It must be ensured that the secret keys of the app server and the sdk server are the same. There is a risk of secret leakage, and it is best to save it on the user server. |
      | platform       | Y        | int    | Platform type iOS 1, Android 2, Windows 3, OSX 4, WEB 5, applet 6, linux 7 |
      | uid            | Y        | string | User ID, with a maximum length of 64 characters, must be unique within an APP |
    
    - **Return Parameter**
      ```
      {
         "errCode": 0,
         "errMsg": "",
         "data":{
            "uid": "",
            "token": "",
            "expiredTime": 0,
         }
      }
      ```

- **API call description**
   
   ```
   app：app client
   app-server：app server
   open-im-sdk：Tuoyun's open source sdk
   open-im-server：Tuoyun's open source sdk service 
   ```

- **Authentication Clow Chart**

![avatar](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/open-im-server.png)
  
## Architecture

![avatar](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/Architecture.jpg)

## License

Open-IM-Server is under the Apache 2.0 license. See
the [LICENSE](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/LICENSE) file for details.
