# OpenIM enviroment


## How to change the configuration


**Modify the configuration files:**

Three ways to modify the configuration:

#### **1. Recommended using environment variables:**

```bash
export PASSWORD="openIM123" # Set password
export USER="root" # Set username
# Choose chat version and server version https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/conversions/images.md, eg: main, release-v*.*
export CHAT_BRANCH="main"
export SERVER_BRANCH="main"
#... Other environment variables
# MONGO_USERNAME: This sets the MongoDB username
# MONGO_PASSWORD: Set the MongoDB password
# MONGO_DATABASE: Sets the MongoDB database name
# MINIO_ENDPOINT: set the MinIO service address
# API_URL: under network environment, set OpenIM Server API address
export API_URL="http://127.0.0.1:10002"
```

Next, update the configuration using `make init`:

```bash
make init
```

#### **2. Modify the automation script:**

```bash
scripts/install/environment.sh
```

Next, update the configuration using `make init`:

```bash
make init
```

#### 3. Modify `config.yaml` and `.env` files (but will be overwritten when using `make init` again).

The `config/config.yaml` file has detailed configuration instructions for the storage components.


The config file is available via [environment.sh](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/scripts/install/environment.sh) configuration [openim.yaml](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/deployments/templates/openim.yaml) template, and then through the `make init` to automatically generate a new configuration.




## Configuration Details

###### Zookeeper

- **Purpose**: Used for RPC service discovery and registration, cluster support.
  
    ```bash
    zookeeper:
      schema: openim                          # Not recommended to modify
      address: [ 127.0.0.1:2181 ]            # Address
      username:                               # Username
      password:                               # Password
    ```

###### MySQL

- **Purpose**: Used for storing users, relationships, and groups. Supports master-slave database.

    ```bash
    mysql:
      address: [ 127.0.0.1:13306 ]            # Address
      username: root                          # Username
      password: openIM123                     # Password
      database: openIM_v2                     # Not recommended to modify
      maxOpenConn: 1000                       # Maximum connection
      maxIdleConn: 100                        # Maximum idle connection
      maxLifeTime: 60                         # Max time a connection can be reused (seconds)
      logLevel: 4                             # Log level (1=silent, 2=error, 3=warn, 4=info)
      slowThreshold: 500                      # Slow statement threshold (milliseconds)
    ```

###### Mongo

- **Purpose**: Used for storing offline messages. Supports mongo sharded clusters.

    ```bash
    mongo:
      uri:                                    # Use this value directly if not empty
      address: [ 127.0.0.1:37017 ]            # Address
      database: openIM                        # Default mongo db
      username: root                          # Username
      password: openIM123                     # Password
      maxPoolSize: 100                        # Maximum connections
    ```

###### Redis

- **Purpose**: Used for storing message sequence numbers, latest messages, user tokens, and MySQL cache. Supports cluster deployment.

    ```bash
    redis:
      address: [ 127.0.0.1:16379 ]            # Address
      username:                               # Username
      password: openIM123                     # Password
    ```

###### Kafka

- **Purpose**: Used for message queues for decoupling. Supports cluster deployment.

    ```bash
    kafka:
      username:                               # Username
      password:                               # Password
      addr: [ 127.0.0.1:9092 ]                # Address
      topics:
        latestMsgToRedis: "latestMsgToRedis"
        offlineMsgToMongo: "offlineMsgToMongoMysql"
        msgToPush: "msgToPush"
        msgToModify: "msgToModify"
      consumerGroupID:
        msgToRedis: redis
        msgToMongo: mongo
        msgToMySql: mysql
        msgToPush: push
        msgToModify: modify
    ```



## Config options

...