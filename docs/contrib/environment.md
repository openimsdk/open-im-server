# OpenIM ENVIRONMENT CONFIGURATION

<!-- vscode-markdown-toc -->
* 1. [OpenIM Deployment Guide](#OpenIMDeploymentGuide)
	* 1.1. [Deployment Strategies](#DeploymentStrategies)
	* 1.2. [Source Code Deployment](#SourceCodeDeployment)
	* 1.3. [Docker Compose Deployment](#DockerComposeDeployment)
	* 1.4. [Environment Variable Configuration](#EnvironmentVariableConfiguration)
		* 1.4.1. [Recommended using environment variables](#Recommendedusingenvironmentvariables)
		* 1.4.2. [Additional Configuration](#AdditionalConfiguration)
		* 1.4.3. [Security Considerations](#SecurityConsiderations)
		* 1.4.4. [Data Management](#DataManagement)
		* 1.4.5. [Monitoring and Logging](#MonitoringandLogging)
		* 1.4.6. [Troubleshooting](#Troubleshooting)
		* 1.4.7. [Conclusion](#Conclusion)
		* 1.4.8. [Additional Resources](#AdditionalResources)
* 2. [Further Configuration](#FurtherConfiguration)
	* 2.1. [Image Registry Configuration](#ImageRegistryConfiguration)
	* 2.2. [OpenIM Docker Network Configuration](#OpenIMDockerNetworkConfiguration)
	* 2.3. [OpenIM Configuration](#OpenIMConfiguration)
	* 2.4. [OpenIM Chat Configuration](#OpenIMChatConfiguration)
	* 2.5. [Zookeeper Configuration](#ZookeeperConfiguration)
	* 2.6. [MySQL Configuration](#MySQLConfiguration)
	* 2.7. [MongoDB Configuration](#MongoDBConfiguration)
	* 2.8. [Tencent Cloud COS Configuration](#TencentCloudCOSConfiguration)
	* 2.9. [Alibaba Cloud OSS Configuration](#AlibabaCloudOSSConfiguration)
	* 2.10. [Redis Configuration](#RedisConfiguration)
	* 2.11. [Kafka Configuration](#KafkaConfiguration)
	* 2.12. [OpenIM Web Configuration](#OpenIMWebConfiguration)
	* 2.13. [RPC Configuration](#RPCConfiguration)
	* 2.14. [Prometheus Configuration](#PrometheusConfiguration)
	* 2.15. [Grafana Configuration](#GrafanaConfiguration)
	* 2.16. [RPC Port Configuration Variables](#RPCPortConfigurationVariables)
	* 2.17. [RPC Register Name Configuration](#RPCRegisterNameConfiguration)
	* 2.18. [Log Configuration](#LogConfiguration)
	* 2.19. [Additional Configuration Variables](#AdditionalConfigurationVariables)
	* 2.20. [Prometheus Configuration](#PrometheusConfiguration-1)
		* 2.20.1. [General Configuration](#GeneralConfiguration)
		* 2.20.2. [Service-Specific Prometheus Ports](#Service-SpecificPrometheusPorts)

## 0. <a name='TableofContents'></a>OpenIM Config File

Ensuring that OpenIM operates smoothly requires clear direction on the configuration file's location. Here's a detailed step-by-step guide on how to provide this essential path to OpenIM:

1. **Using the Command-line Argument**:

   + **For Configuration Path**: When initializing OpenIM, you can specify the path to the configuration file directly using the `-c` or `--config_folder_path` option.

     ```bash
     ❯ _output/bin/platforms/linux/amd64/openim-api --config_folder_path="/your/config/folder/path"
     ```

   + **For Port Specification**: Similarly, if you wish to designate a particular port, utilize the `-p` option followed by the desired port number.

     ```bash
     ❯ _output/bin/platforms/linux/amd64/openim-api -p 1234
     ```

     Note: If the port is not specified here, OpenIM will fetch it from the configuration file. Setting the port via environment variables isn't supported. We recommend consolidating settings in the configuration file for a more consistent and streamlined setup.

2. **Leveraging the Environment Variable**:

   You have the flexibility to determine OpenIM's configuration path by setting an `OPENIMCONFIG` environment variable. This method provides a seamless way to instruct OpenIM without command-line parameters every time.

   ```bash
   export OPENIMCONFIG="/path/to/your/config"
   ```

3. **Relying on the Default Path**:

   In scenarios where neither command-line arguments nor environment variables are provided, OpenIM will intuitively revert to the `config/` directory to locate its configuration.



##  1. <a name='OpenIMDeploymentGuide'></a>OpenIM Deployment Guide

Welcome to the OpenIM Deployment Guide! OpenIM offers a versatile and robust instant messaging server, and deploying it can be achieved through various methods. This guide will walk you through the primary deployment strategies, ensuring you can set up OpenIM in a way that best suits your needs.

###  1.1. <a name='DeploymentStrategies'></a>Deployment Strategies

OpenIM provides multiple deployment methods, each tailored to different use cases and technical preferences:

1. **[Source Code Deployment Guide](https://doc.rentsoft.cn/guides/gettingStarted/imSourceCodeDeployment)**
2. **[Docker Deployment Guide](https://doc.rentsoft.cn/guides/gettingStarted/dockerCompose)**
3. **[Kubernetes Deployment Guide](https://github.com/openimsdk/open-im-server/tree/main/deployments)**

While the first two methods will be our main focus, it's worth noting that the third method, Kubernetes deployment, is also viable and can be rendered via the `environment.sh` script variables.

###  1.2. <a name='SourceCodeDeployment'></a>Source Code Deployment

In the source code deployment method, the configuration generation process involves executing `make init`, which fundamentally runs the script `./scripts/init-config.sh`. This script utilizes variables defined in the [`environment.sh`](https://github.com/openimsdk/open-im-server/blob/main/scripts/install/environment.sh) script to render the [`openim.yaml`](https://github.com/openimsdk/open-im-server/blob/main/deployments/templates/openim.yaml) template file, subsequently generating the [`config.yaml`](https://github.com/openimsdk/open-im-server/blob/main/config/config.yaml) configuration file.

###  1.3. <a name='DockerComposeDeployment'></a>Docker Compose Deployment

Docker deployment offers a slightly more intricate template. Within the [openim-server](https://github.com/openimsdk/openim-docker/tree/main/openim-server) directory, multiple subdirectories correspond to various versions, each aligning with `openim-chat` as illustrated below:

| openim-server                                                | openim-chat                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ |
| [main](https://github.com/openimsdk/openim-docker/tree/main/openim-server/main) | [main](https://github.com/openimsdk/openim-docker/tree/main/openim-chat/main) |
| [release-v3.2](https://github.com/openimsdk/openim-docker/tree/main/openim-server/release-v3.3) | [release-v3.2](https://github.com/openimsdk/openim-docker/tree/main/openim-chat/release-v1.3) |
| [release-v3.2](https://github.com/openimsdk/openim-docker/tree/main/openim-server/release-v3.2) | [release-v3.2](https://github.com/openimsdk/openim-docker/tree/main/openim-chat/release-v1.2) |

Configuration file modifications can be made by specifying corresponding environment variables, for instance:

```bash
export CHAT_BRANCH="main"   
export SERVER_BRANCH="main" 
```

These variables are stored within the [`environment.sh`](https://github.com/OpenIMSDK/openim-docker/blob/main/scripts/install/environment.sh) configuration:

```bash
readonly CHAT_BRANCH=${CHAT_BRANCH:-'main'}
readonly SERVER_BRANCH=${SERVER_BRANCH:-'main'}
```

Setting a variable, e.g., `export CHAT_BRANCH="release-v1.3"`, will prioritize `CHAT_BRANCH="release-v1.3"` as the variable value. Ultimately, the chosen image version is determined, and rendering is achieved through `make init` (or `./scripts/init-config.sh`).

> Note: Direct modifications to the `config.yaml` file are also permissible without utilizing `make init`.

###  1.4. <a name='EnvironmentVariableConfiguration'></a>Environment Variable Configuration

For convenience, configuration through modifying environment variables is recommended:

####  1.4.1. <a name='Recommendedusingenvironmentvariables'></a>Recommended using environment variables

+ PASSWORD

  + **Description**: Password for mysql, mongodb, redis, and minio.
  + **Default**: `openIM123`
  + Notes:
    + Minimum password length: 8 characters.
    + Special characters are not allowed.

  ```bash
  export PASSWORD="openIM123"
  ```

+ OPENIM_USER

  + **Description**: Username for mysql, mongodb, redis, and minio.
  + **Default**: `root`

  ```bash
  export OPENIM_USER="root"
  ```

+ API_URL

  + **Description**: API address.
  + **Note**: If the server has an external IP, it will be automatically obtained. For internal networks, set this variable to the IP serving internally.

  ```
  export API_URL="http://ip:10002"
  ```

+ DATA_DIR

  + **Description**: Data mount directory for components.
  + **Default**: `/data/openim`

  ```bash
  export DATA_DIR="/data/openim"
  ```

####  1.4.2. <a name='AdditionalConfiguration'></a>Additional Configuration

##### MinIO Access and Secret Key

To secure your MinIO server, you should set up an access key and secret key. These credentials are used to authenticate requests to your MinIO server.

```bash
export MINIO_ACCESS_KEY="YourAccessKey"
export MINIO_SECRET_KEY="YourSecretKey"
```

##### MinIO Browser

MinIO comes with an embedded web-based object browser. You can control the availability of the MinIO browser by setting the `MINIO_BROWSER` environment variable.

```bash
export MINIO_BROWSER="on"
```

####  1.4.3. <a name='SecurityConsiderations'></a>Security Considerations

##### TLS/SSL Configuration

For secure communication, it's recommended to enable TLS/SSL for your MinIO server. You can do this by providing the path to the SSL certificate and key files.

```bash
export MINIO_CERTS_DIR="/path/to/certs/directory"
```

####  1.4.4. <a name='DataManagement'></a>Data Management

##### Data Retention Policy

You may want to set up a data retention policy to automatically delete objects after a specified period.

```bash
export MINIO_RETENTION_DAYS="30"
```

####  1.4.5. <a name='MonitoringandLogging'></a>Monitoring and Logging

##### [Audit Logging](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/environment.md#audit-logging)

Enable audit logging to keep track of access and changes to your data.

```bash
export MINIO_AUDIT="on"
```

####  1.4.6. <a name='Troubleshooting'></a>Troubleshooting

##### Debug Mode

In case of issues, you may enable debug mode to get more detailed logs to assist in troubleshooting.

```bash
export MINIO_DEBUG="on"
```

####  1.4.7. <a name='Conclusion'></a>Conclusion

With the environment variables configured as per your requirements, your MinIO server should be ready to securely store and manage your object data. Ensure to verify the setup and monitor the logs for any unusual activities or errors. Regularly update the MinIO server and review your configuration to adapt to any changes or improvements in the MinIO system.

####  1.4.8. <a name='AdditionalResources'></a>Additional Resources

+ [MinIO Client Quickstart Guide](https://docs.min.io/docs/minio-client-quickstart-guide)
+ [MinIO Admin Complete Guide](https://docs.min.io/docs/minio-admin-complete-guide)
+ [MinIO Docker Quickstart Guide](https://docs.min.io/docs/minio-docker-quickstart-guide)

Feel free to explore the MinIO documentation for more advanced configurations and usage scenarios.



##  2. <a name='FurtherConfiguration'></a>Further Configuration

###  2.1. <a name='ImageRegistryConfiguration'></a>Image Registry Configuration

**Description**: The image registry configuration allows users to select an image address for use. The default is set to use GITHUB images, but users can opt for Docker Hub or Ali Cloud, especially beneficial for Chinese users due to its local proximity.

| Parameter        | Default Value         | Description                                                  |
| ---------------- | --------------------- | ------------------------------------------------------------ |
| `IMAGE_REGISTRY` | `"ghcr.io/openimsdk"` | The registry from which Docker images will be pulled. Other options include `"openim"` and `"registry.cn-hangzhou.aliyuncs.com/openimsdk"`. |

###  2.2. <a name='OpenIMDockerNetworkConfiguration'></a>OpenIM Docker Network Configuration

**Description**: This section configures the Docker network subnet and generates IP addresses for various services within the defined subnet.

| Parameter                   | Example Value     | Description                                                  |
| --------------------------- | ----------------- | ------------------------------------------------------------ |
| `DOCKER_BRIDGE_SUBNET`      | `'172.28.0.0/16'` | The subnet for the Docker network.                           |
| `DOCKER_BRIDGE_GATEWAY`     | Generated IP      | The gateway IP address within the Docker subnet.             |
| `[SERVICE]_NETWORK_ADDRESS` | Generated IP      | The network IP address for a specific service (e.g., MYSQL, MONGO, REDIS, etc.) within the Docker subnet. |

###  2.3. <a name='OpenIMConfiguration'></a>OpenIM Configuration

**Description**: OpenIM configuration involves setting up directories for data, installation, configuration, and logs. It also involves configuring the OpenIM server address and ports for WebSocket and API.

| Parameter               | Default Value            | Description                               |
| ----------------------- | ------------------------ | ----------------------------------------- |
| `OPENIM_DATA_DIR`       | `"/data/openim"`         | Directory for OpenIM data.                |
| `OPENIM_INSTALL_DIR`    | `"/opt/openim"`          | Directory where OpenIM is installed.      |
| `OPENIM_CONFIG_DIR`     | `"/etc/openim"`          | Directory for OpenIM configuration files. |
| `OPENIM_LOG_DIR`        | `"/var/log/openim"`      | Directory for OpenIM logs.                |
| `OPENIM_SERVER_ADDRESS` | Docker Bridge Gateway IP | OpenIM server address.                    |
| `OPENIM_WS_PORT`        | `'10001'`                | Port for OpenIM WebSocket.                |
| `API_OPENIM_PORT`       | `'10002'`                | Port for OpenIM API.                      |

###  2.4. <a name='OpenIMChatConfiguration'></a>OpenIM Chat Configuration

**Description**: Configuration for OpenIM chat, including data directory, server address, and ports for API and chat functionalities.

| Parameter               | Example Value              | Description                     |
| ----------------------- | -------------------------- | ------------------------------- |
| `OPENIM_CHAT_DATA_DIR`  | `"./openim-chat/[BRANCH]"` | Directory for OpenIM chat data. |
| `OPENIM_CHAT_ADDRESS`   | Docker Bridge Gateway IP   | OpenIM chat service address.    |
| `OPENIM_CHAT_API_PORT`  | `"10008"`                  | Port for OpenIM chat API.       |
| `OPENIM_ADMIN_API_PORT` | `"10009"`                  | Port for OpenIM Admin API.      |
| `OPENIM_ADMIN_PORT`     | `"30200"`                  | Port for OpenIM chat Admin.     |
| `OPENIM_CHAT_PORT`      | `"30300"`                  | Port for OpenIM chat.           |

###  2.5. <a name='ZookeeperConfiguration'></a>Zookeeper Configuration

**Description**: Configuration for Zookeeper, including schema, port, address, and credentials.

| Parameter            | Example Value            | Description             |
| -------------------- | ------------------------ | ----------------------- |
| `ZOOKEEPER_SCHEMA`   | `"openim"`               | Schema for Zookeeper.   |
| `ZOOKEEPER_PORT`     | `"12181"`                | Port for Zookeeper.     |
| `ZOOKEEPER_ADDRESS`  | Docker Bridge Gateway IP | Address for Zookeeper.  |
| `ZOOKEEPER_USERNAME` | `""`                     | Username for Zookeeper. |
| `ZOOKEEPER_PASSWORD` | `""`                     | Password for Zookeeper. |

###  2.6. <a name='MySQLConfiguration'></a>MySQL Configuration

**Description**: Configuration for MySQL, including port, address, and credentials.

| Parameter        | Example Value            | Description         |
| ---------------- | ------------------------ | ------------------- |
| `MYSQL_PORT`     | `"13306"`                | Port for MySQL.     |
| `MYSQL_ADDRESS`  | Docker Bridge Gateway IP | Address for MySQL.  |
| `MYSQL_USERNAME` | User-defined             | Username for MySQL. |
| `MYSQL_PASSWORD` | User-defined             | Password for MySQL. |

Note: The configurations for other services (e.g., MONGO, REDIS, KAFKA, etc.) follow a similar pattern to MySQL and can be documented in a similar manner.

###  2.7. <a name='MongoDBConfiguration'></a>MongoDB Configuration

This section involves setting up MongoDB, including its port, address, and credentials.

| Parameter      | Example Value  | Description             |
| -------------- | -------------- | ----------------------- |
| MONGO_PORT     | "27017"        | Port used by MongoDB.   |
| MONGO_ADDRESS  | [Generated IP] | IP address for MongoDB. |
| MONGO_USERNAME | [User Defined] | Username for MongoDB.   |
| MONGO_PASSWORD | [User Defined] | Password for MongoDB.   |

###  2.8. <a name='TencentCloudCOSConfiguration'></a>Tencent Cloud COS Configuration

This section involves setting up Tencent Cloud COS, including its bucket URL and credentials.

| Parameter         | Example Value                                                | Description                          |
| ----------------- | ------------------------------------------------------------ | ------------------------------------ |
| COS_BUCKET_URL    | "[https://temp-1252357374.cos.ap-chengdu.myqcloud.com](https://temp-1252357374.cos.ap-chengdu.myqcloud.com/)" | Tencent Cloud COS bucket URL.        |
| COS_SECRET_ID     | [User Defined]                                               | Secret ID for Tencent Cloud COS.     |
| COS_SECRET_KEY    | [User Defined]                                               | Secret key for Tencent Cloud COS.    |
| COS_SESSION_TOKEN | [User Defined]                                               | Session token for Tencent Cloud COS. |
| COS_PUBLIC_READ   | "false"                                                      | Public read access.                  |

###  2.9. <a name='AlibabaCloudOSSConfiguration'></a>Alibaba Cloud OSS Configuration

This section involves setting up Alibaba Cloud OSS, including its endpoint, bucket name, and credentials.

| Parameter             | Example Value                                                | Description                              |
| --------------------- | ------------------------------------------------------------ | ---------------------------------------- |
| OSS_ENDPOINT          | "[https://oss-cn-chengdu.aliyuncs.com](https://oss-cn-chengdu.aliyuncs.com/)" | Endpoint URL for Alibaba Cloud OSS.      |
| OSS_BUCKET            | "demo-9999999"                                               | Bucket name for Alibaba Cloud OSS.       |
| OSS_BUCKET_URL        | "[https://demo-9999999.oss-cn-chengdu.aliyuncs.com](https://demo-9999999.oss-cn-chengdu.aliyuncs.com/)" | Bucket URL for Alibaba Cloud OSS.        |
| OSS_ACCESS_KEY_ID     | [User Defined]                                               | Access key ID for Alibaba Cloud OSS.     |
| OSS_ACCESS_KEY_SECRET | [User Defined]                                               | Access key secret for Alibaba Cloud OSS. |
| OSS_SESSION_TOKEN     | [User Defined]                                               | Session token for Alibaba Cloud OSS.     |
| OSS_PUBLIC_READ       | "false"                                                      | Public read access.                      |

###  2.10. <a name='RedisConfiguration'></a>Redis Configuration

This section involves setting up Redis, including its port, address, and credentials.

| Parameter      | Example Value              | Description           |
| -------------- | -------------------------- | --------------------- |
| REDIS_PORT     | "16379"                    | Port used by Redis.   |
| REDIS_ADDRESS  | "${DOCKER_BRIDGE_GATEWAY}" | IP address for Redis. |
| REDIS_USERNAME | [User Defined]             | Username for Redis.   |
| REDIS_PASSWORD | "${PASSWORD}"              | Password for Redis.   |

###  2.11. <a name='KafkaConfiguration'></a>Kafka Configuration

This section involves setting up Kafka, including its port, address, credentials, and topics.

| Parameter                    | Example Value              | Description                         |
| ---------------------------- | -------------------------- | ----------------------------------- |
| KAFKA_USERNAME               | [User Defined]             | Username for Kafka.                 |
| KAFKA_PASSWORD               | [User Defined]             | Password for Kafka.                 |
| KAFKA_PORT                   | "19094"                    | Port used by Kafka.                 |
| KAFKA_ADDRESS                | "${DOCKER_BRIDGE_GATEWAY}" | IP address for Kafka.               |
| KAFKA_LATESTMSG_REDIS_TOPIC  | "latestMsgToRedis"         | Topic for latest message to Redis.  |
| KAFKA_OFFLINEMSG_MONGO_TOPIC | "offlineMsgToMongoMysql"   | Topic for offline message to Mongo. |
| KAFKA_MSG_PUSH_TOPIC         | "msgToPush"                | Topic for message to push.          |
| KAFKA_CONSUMERGROUPID_REDIS  | "redis"                    | Consumer group ID to Redis.         |
| KAFKA_CONSUMERGROUPID_MONGO  | "mongo"                    | Consumer group ID to Mongo.         |
| KAFKA_CONSUMERGROUPID_MYSQL  | "mysql"                    | Consumer group ID to MySQL.         |
| KAFKA_CONSUMERGROUPID_PUSH   | "push"                     | Consumer group ID to push.          |

Note: Ensure to replace placeholder values (like [User Defined], `${DOCKER_BRIDGE_GATEWAY}`, and `${PASSWORD}`) with actual values before deploying the configuration.



###  2.12. <a name='OpenIMWebConfiguration'></a>OpenIM Web Configuration

This section involves setting up OpenIM Web, including its port, address, and dist path.

| Parameter            | Example Value              | Description               |
| -------------------- | -------------------------- | ------------------------- |
| OPENIM_WEB_PORT      | "11001"                    | Port used by OpenIM Web.  |
| OPENIM_WEB_ADDRESS   | "${DOCKER_BRIDGE_GATEWAY}" | Address for OpenIM Web.   |
| OPENIM_WEB_DIST_PATH | "/app/dist"                | Dist path for OpenIM Web. |

###  2.13. <a name='RPCConfiguration'></a>RPC Configuration

Configuration for RPC, including the register and listen IP.

| Parameter       | Example Value  | Description          |
| --------------- | -------------- | -------------------- |
| RPC_REGISTER_IP | [User Defined] | Register IP for RPC. |
| RPC_LISTEN_IP   | "0.0.0.0"      | Listen IP for RPC.   |

###  2.14. <a name='PrometheusConfiguration'></a>Prometheus Configuration

Setting up Prometheus, including its port and address.

| Parameter          | Example Value              | Description              |
| ------------------ | -------------------------- | ------------------------ |
| PROMETHEUS_PORT    | "19090"                    | Port used by Prometheus. |
| PROMETHEUS_ADDRESS | "${DOCKER_BRIDGE_GATEWAY}" | Address for Prometheus.  |

###  2.15. <a name='GrafanaConfiguration'></a>Grafana Configuration

Configuration for Grafana, including its port and address.

| Parameter       | Example Value              | Description           |
| --------------- | -------------------------- | --------------------- |
| GRAFANA_PORT    | "3000"                     | Port used by Grafana. |
| GRAFANA_ADDRESS | "${DOCKER_BRIDGE_GATEWAY}" | Address for Grafana.  |

###  2.16. <a name='RPCPortConfigurationVariables'></a>RPC Port Configuration Variables

Configuration for various RPC ports. Note: For launching multiple programs, just fill in multiple ports separated by commas. Try not to have spaces.

| Parameter                   | Example Value | Description                         |
| --------------------------- | ------------- | ----------------------------------- |
| OPENIM_USER_PORT            | '10110'       | OpenIM User Service Port.           |
| OPENIM_FRIEND_PORT          | '10120'       | OpenIM Friend Service Port.         |
| OPENIM_MESSAGE_PORT         | '10130'       | OpenIM Message Service Port.        |
| OPENIM_MESSAGE_GATEWAY_PORT | '10140'       | OpenIM Message Gateway Service Port |
| OPENIM_GROUP_PORT           | '10150'       | OpenIM Group Service Port.          |
| OPENIM_AUTH_PORT            | '10160'       | OpenIM Authorization Service Port.  |
| OPENIM_PUSH_PORT            | '10170'       | OpenIM Push Service Port.           |
| OPENIM_CONVERSATION_PORT    | '10180'       | OpenIM Conversation Service Port.   |
| OPENIM_THIRD_PORT           | '10190'       | OpenIM Third-Party Service Port.    |

###  2.17. <a name='RPCRegisterNameConfiguration'></a>RPC Register Name Configuration

This section involves setting up the RPC Register Names for various OpenIM services.

| Parameter                   | Example Value    | Description                         |
| --------------------------- | ---------------- | ----------------------------------- |
| OPENIM_USER_NAME            | "User"           | OpenIM User Service Name            |
| OPENIM_FRIEND_NAME          | "Friend"         | OpenIM Friend Service Name          |
| OPENIM_MSG_NAME             | "Msg"            | OpenIM Message Service Name         |
| OPENIM_PUSH_NAME            | "Push"           | OpenIM Push Service Name            |
| OPENIM_MESSAGE_GATEWAY_NAME | "MessageGateway" | OpenIM Message Gateway Service Name |
| OPENIM_GROUP_NAME           | "Group"          | OpenIM Group Service Name           |
| OPENIM_AUTH_NAME            | "Auth"           | OpenIM Authorization Service Name   |
| OPENIM_CONVERSATION_NAME    | "Conversation"   | OpenIM Conversation Service Name    |
| OPENIM_THIRD_NAME           | "Third"          | OpenIM Third-Party Service Name     |

###  2.18. <a name='LogConfiguration'></a>Log Configuration

This section involves configuring the log settings, including storage location, rotation time, and log level.

| Parameter                 | Example Value            | Description                       |
| ------------------------- | ------------------------ | --------------------------------- |
| LOG_STORAGE_LOCATION      | ""${OPENIM_ROOT}"/logs/" | Location for storing logs         |
| LOG_ROTATION_TIME         | "24"                     | Log rotation time (in hours)      |
| LOG_REMAIN_ROTATION_COUNT | "2"                      | Number of log rotations to retain |
| LOG_REMAIN_LOG_LEVEL      | "6"                      | Log level to retain               |
| LOG_IS_STDOUT             | "false"                  | Output log to standard output     |
| LOG_IS_JSON               | "false"                  | Log in JSON format                |
| LOG_WITH_STACK            | "false"                  | Include stack info in logs        |

###  2.19. <a name='AdditionalConfigurationVariables'></a>Additional Configuration Variables

This section involves setting up additional configuration variables for Websocket, Push Notifications, and Chat.

| Parameter               | Example Value     | Description                        |
| ----------------------- | ----------------- | ---------------------------------- |
| WEBSOCKET_MAX_CONN_NUM  | "100000"          | Maximum Websocket connections      |
| WEBSOCKET_MAX_MSG_LEN   | "4096"            | Maximum Websocket message length   |
| WEBSOCKET_TIMEOUT       | "10"              | Websocket timeout                  |
| PUSH_ENABLE             | "getui"           | Push notification enable status    |
| GETUI_PUSH_URL          | [Generated URL]   | GeTui Push Notification URL        |
| GETUI_MASTER_SECRET     | [User Defined]    | GeTui Master Secret                |
| GETUI_APP_KEY           | [User Defined]    | GeTui Application Key              |
| GETUI_INTENT            | [User Defined]    | GeTui Push Intent                  |
| GETUI_CHANNEL_ID        | [User Defined]    | GeTui Channel ID                   |
| GETUI_CHANNEL_NAME      | [User Defined]    | GeTui Channel Name                 |
| FCM_SERVICE_ACCOUNT     | "x.json"          | FCM Service Account                |
| JPNS_APP_KEY            | [User Defined]    | JPNS Application Key               |
| JPNS_MASTER_SECRET      | [User Defined]    | JPNS Master Secret                 |
| JPNS_PUSH_URL           | [User Defined]    | JPNS Push Notification URL         |
| JPNS_PUSH_INTENT        | [User Defined]    | JPNS Push Intent                   |
| MANAGER_USERID_1        | "openIM123456"    | Administrator ID 1                 |
| MANAGER_USERID_2        | "openIM654321"    | Administrator ID 2                 |
| MANAGER_USERID_3        | "openIMAdmin"     | Administrator ID 3                 |
| NICKNAME_1              | "system1"         | Nickname 1                         |
| NICKNAME_2              | "system2"         | Nickname 2                         |
| NICKNAME_3              | "system3"         | Nickname 3                         |
| MULTILOGIN_POLICY       | "1"               | Multi-login Policy                 |
| CHAT_PERSISTENCE_MYSQL  | "true"            | Chat Persistence in MySQL          |
| MSG_CACHE_TIMEOUT       | "86400"           | Message Cache Timeout              |
| GROUP_MSG_READ_RECEIPT  | "true"            | Group Message Read Receipt Enable  |
| SINGLE_MSG_READ_RECEIPT | "true"            | Single Message Read Receipt Enable |
| RETAIN_CHAT_RECORDS     | "365"             | Retain Chat Records (in days)      |
| CHAT_RECORDS_CLEAR_TIME | [Cron Expression] | Chat Records Clear Time            |
| MSG_DESTRUCT_TIME       | [Cron Expression] | Message Destruct Time              |
| SECRET                  | "${PASSWORD}"     | Secret Key                         |
| TOKEN_EXPIRE            | "90"              | Token Expiry Time                  |
| FRIEND_VERIFY           | "false"           | Friend Verification Enable         |
| IOS_PUSH_SOUND          | "xxx"             | iOS                                |



###  2.20. <a name='PrometheusConfiguration-1'></a>Prometheus Configuration

This section involves configuring Prometheus, including enabling/disabling it and setting up ports for various services.

####  2.20.1. <a name='GeneralConfiguration'></a>General Configuration

| Parameter           | Example Value | Description                   |
| ------------------- | ------------- | ----------------------------- |
| `PROMETHEUS_ENABLE` | "false"       | Whether to enable Prometheus. |

####  2.20.2. <a name='Service-SpecificPrometheusPorts'></a>Service-Specific Prometheus Ports

| Service                  | Parameter                | Default Port Value           | Description                                        |
| ------------------------ | ------------------------ | ---------------------------- | -------------------------------------------------- |
| User Service             | `USER_PROM_PORT`         | '20110'                      | Prometheus port for the User service.              |
| Friend Service           | `FRIEND_PROM_PORT`       | '20120'                      | Prometheus port for the Friend service.            |
| Message Service          | `MESSAGE_PROM_PORT`      | '20130'                      | Prometheus port for the Message service.           |
| Message Gateway          | `MSG_GATEWAY_PROM_PORT`  | '20140'                      | Prometheus port for the Message Gateway.           |
| Group Service            | `GROUP_PROM_PORT`        | '20150'                      | Prometheus port for the Group service.             |
| Auth Service             | `AUTH_PROM_PORT`         | '20160'                      | Prometheus port for the Auth service.              |
| Push Service             | `PUSH_PROM_PORT`         | '20170'                      | Prometheus port for the Push service.              |
| Conversation Service     | `CONVERSATION_PROM_PORT` | '20230'                      | Prometheus port for the Conversation service.      |
| RTC Service              | `RTC_PROM_PORT`          | '21300'                      | Prometheus port for the RTC service.               |
| Third Service            | `THIRD_PROM_PORT`        | '21301'                      | Prometheus port for the Third service.             |
| Message Transfer Service | `MSG_TRANSFER_PROM_PORT` | '21400, 21401, 21402, 21403' | Prometheus ports for the Message Transfer service. |
